#include <linux/bpf.h>
#include <stdbool.h>
#include <linux/if_ether.h>
#include <linux/udp.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/seg6.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#include "xdp_map.h"
#include "dns.h"

static int parse_dns_query(struct xdp_md *ctx, void *dns_query_start, struct dns_query *query);

SEC("xdp")
int intercept_dns(struct xdp_md *ctx) {
    // Initialize data
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;

    // Scanning ethernet header
    struct ethhdr *eth = data;

    // Check ethernet header validity
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;

    // Validate ethernet header: Check if the EtherType is IPv6 (0x86DD)
    if (eth->h_proto != bpf_htons(ETH_P_IPV6))
        return XDP_PASS;

    // Scanning IPv6 header
    struct ipv6hdr *ipv6 = (void *)(eth + 1);

    // Check IPv6 header validity
    if ((void *)(ipv6 + 1) > data_end)
        return XDP_PASS;

    // Validate IPv6 header: Check if the Next Header is UDP (17)
    if (ipv6->nexthdr != IPPROTO_UDP)
        return XDP_PASS;
        

    // Scanning UDP header
    struct udphdr *udp = (void *)(ipv6 + 1);

    // Check UDP header validity
    if ((void *)(udp + 1) > data_end)
        return XDP_PASS;

    // Validate UDP header: Check if the UDP destination port is 53 (DNS)
    if (udp->source != bpf_htons(UDP_P_DNS))
        return XDP_PASS;

    // Scanning DNS header
    // struct dns_header *dns = (void *)(udp + 1);
    struct dns_hdr *dns = (void *)(udp + 1);

    // Check DNS header validity
    if ((void *)(dns + 1) > data_end)
        return XDP_PASS;

    // Only proceed if there's at least one answer
    // if (dns->ancount == 0)
    //     return XDP_PASS;
    if (dns->ans_count == 0)
        return XDP_PASS;

    // Scanning DNS data
    __u8 *dns_data = (__u8 *)(dns + 1);
    if ((void *)(dns_data) > data_end)  // Check if enough data exists
        return XDP_PASS;

    // Get a pointer to the start of the DNS query
    // void *query_start = dns + 32;
    // struct dns_hdr *dns_hdr = (struct dns_hdr *)dns;
    // void *query_start = (void *)dns_hdr + sizeof(struct dns_hdr);
    void *query_start = (void *)dns + sizeof(struct dns_hdr);
    struct dns_query query;
    int query_length = 0;
    // Parse the DNS query
    query_length = parse_dns_query(ctx, query_start, &query);
    if (query_length < 1) {
        bpf_printk("Error parsing DNS query\n");
        return XDP_PASS;
    }

    // Check if the query is a AAAA query (IPv6 query) dns record type 28
    if (query.record_type != 28) {
        bpf_printk("Not an AAAA query\n");
        return XDP_PASS;
    }

    // if ((void *)(dns_data + 16) > data_end)  
    //     bpf_printk("Error: boundary exceeded while retrieving IPv6 address");
    //     return XDP_PASS;

    // Assuming that the first answer is the one we want
    struct in6_addr *ipv6_address = (struct in6_addr *)(dns_data);

    bpf_printk("IPv6 address: %pI6\n", ipv6_address);
    bpf_printk("DNS query name: %s\n", query.name);

    int key = 0;
    bpf_map_update_elem(&test_map, &key, query.name, BPF_ANY);

    // Check if there is an entry with that domain name in the map
    // If there is no entry, we stop here
    struct client_data *client_data;
    client_data = bpf_map_lookup_elem(&client_map, query.name);
    if (!client_data) {
        bpf_printk("No entry found for %s\n", query.name);
        return XDP_PASS;
    }
    bpf_printk("works until here");

    // // Update the destination address of the packet
    if ((void *)(ipv6_address + 1) > data_end) {
        bpf_printk("Error: boundary exceeded while updating IPv6 address");
        return XDP_PASS;
    }
    client_data->dstaddr = *ipv6_address;
    bpf_map_update_elem(&client_map, query.name, client_data, BPF_EXIST);

    return XDP_PASS;
}

//Parse query and return query length
static int parse_dns_query(struct xdp_md *ctx, void *query_start, struct dns_query *q)
{
    void *data_end = (void *)(long)ctx->data_end;
    __u16 i;
    void *cursor = query_start;
    int namepos = 0;

    //Fill dns_query.name with zero bytes
    //Not doing so will make the verifier complain when dns_query is used as a key in bpf_map_lookup
    for (i = 0; i < MAX_DNS_NAME_LENGTH; i++) {
        q->name[i] = 0;
    }
    //Fill record_type and class with default values to satisfy verifier
    q->record_type = 0;
    q->record_class = 0;

    for (i = 0; i < MAX_DNS_NAME_LENGTH; i++) {
        //Boundary check of cursor. Verifier requires a +1 here. 
        //Probably because we are advancing the pointer at the end of the loop
        if (cursor + 1 > data_end) {
            bpf_printk("Error: boundary exceeded while parsing DNS query name");
            break;
        }

        //If separator is zero we've reached the end of the domain query
        if (*(char *)(cursor) == 0) {

            //We've reached the end of the query name.
            //This will be followed by 2x 2 bytes: the dns type and dns class.
            if (cursor + 5 > data_end) {
                bpf_printk("Error: boundary exceeded while retrieving DNS record type and class");
            } else {
                q->record_type = bpf_htons(*(__u16 *)(cursor + 1));
                q->record_class = bpf_htons(*(__u16 *)(cursor + 3));
            }

            //Return the bytecount of (namepos + current '0' byte + dns type + dns class) as the query length.
            return namepos + 1 + 2 + 2;
        }

        //Read and fill data into struct
        q->name[namepos] = *(char *)(cursor);
        namepos++;
        cursor++;
    }

    return -1;
}

char _license[] SEC("license") = "GPL";