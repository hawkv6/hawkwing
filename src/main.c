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
    struct dns_header *dns = (void *)(udp + 1);

    // Check DNS header validity
    if ((void *)(dns + 1) > data_end)
        return XDP_PASS;

    // Only proceed if there's at least one answer
    if (dns->ancount == 0)
        return XDP_PASS;

    // Here, you would normally skip the question section to reach the answer section
    // For this example, we're skipping directly to where we'd expect an IPv6 address.
    // This is a simplification.
    __u8 *dns_data = (__u8 *)(dns + 1);
    if ((void *)(dns_data) > data_end)  // Check if enough data exists
        return XDP_PASS;

    // Get a pointer to the start of the DNS query
    void *query_start = dns_data + 16;
    struct dns_query query;
    int query_length = 0;
    query_length = parse_dns_query(ctx, query_start, &query);
    if (query_length < 1) {
        return XDP_PASS;
    }

    bpf_printk("Query: %s\n", query.name);

    if ((void *)(dns_data + 16) > data_end)  // Check if enough data exists
        return XDP_PASS;

    // Assuming the IPv6 address is directly in the DNS answer section
    // This is a gross simplification; you'll need proper parsing logic here.
    struct in6_addr *ipv6_address = (struct in6_addr *)(dns_data);

    bpf_printk("IPv6 address: %pI6\n", ipv6_address);

    // Check if there is an entry with that domain name in the map
    // If there is no entry, we stop here
    // struct client_data *value;
    // value = bpf_map_lookup_elem(&client_map, &query.name);

    struct client_data *value = bpf_map_lookup_elem(&client_map, &query.name);
    if (!value)
        bpf_printk("No entry found for %s\n", &query.name);
        return XDP_PASS;

    value->dstaddr = *ipv6_address;
    if (bpf_map_update_elem(&client_map, &query.name, value, BPF_EXIST) != 0) {
        bpf_printk("Error updating map\n");
    }

    return XDP_PASS;
}

static int parse_dns_query(struct xdp_md *ctx, void* dns_query_start, struct dns_query *query) {
    void *data_end = (void *)(long)ctx->data_end;
    __u16 i;
    void *cursor = dns_query_start;
    int namepos = 0;

    for (int i = 0; i < sizeof(query->name); i++) {
        query->name[i] = 0;
    }
    query->record_type = 0;
    query->record_class = 0;

    for (i = 0; i < MAX_DNS_NAME_LENGTH; i++) {
        //Boundary check of cursor. Verifier requires a +1 here. 
        //Probably because we are advancing the pointer at the end of the loop
        if (cursor + 1 > data_end) {
            break;
        }

        //If separator is zero we've reached the end of the domain query
        if (*(char *)(cursor) == 0) {

            //We've reached the end of the query name.
            //This will be followed by 2x 2 bytes: the dns type and dns class.
            if (cursor + 5 > data_end) {
                #ifdef DEBUG
                bpf_printk("Error: boundary exceeded while retrieving DNS record type and class");
                #endif
            } else {
                query->record_type = bpf_htons(*(__u16 *)(cursor + 1));
                query->record_class = bpf_htons(*(__u16 *)(cursor + 3));
            }

            //Return the bytecount of (namepos + current '0' byte + dns type + dns class) as the query length.
            return namepos + 1 + 2 + 2;
        }

        //Read and fill data into struct
        query->name[namepos] = *(char *)(cursor);
        namepos++;
        cursor++;
    }

    return -1;
}


char _license[] SEC("license") = "GPL";