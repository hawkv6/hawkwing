#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/seg6.h>
#include <linux/udp.h>
#include <stdbool.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#include "dns.h"
#include "xdp_map.h"

#define memcpy __builtin_memcpy

static int parse_dns_query(struct xdp_md *ctx, void *dns_query_start,
						   struct dns_query *query);
static int parse_dns_answer(struct xdp_md *ctx, struct dns_hdr *dns_hdr,
							int query_length, struct dns_answer *a);

SEC("xdp")
int intercept_dns(struct xdp_md *ctx)
{
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
	struct dns_hdr *dns = (void *)(udp + 1);

	// Check DNS header validity
	if ((void *)(dns + 1) > data_end)
		return XDP_PASS;

	// Check if there is at least one answer
	if (dns->ans_count == 0)
		return XDP_PASS;

	// Get a pointer to the start of the DNS query
	void *query_start = (void *)(dns + 1);
	struct dns_query query;
	int query_length = 0;

	// Parse the DNS query
	query_length = parse_dns_query(ctx, query_start, &query);
	if (query_length < 1) {
#ifdef DEBUG
		bpf_printk("Error parsing DNS query");
#endif
		return XDP_PASS;
	}

	// Check if the query is an AAAA query (IPv6 query) dns record type 28
	if (query.record_type != 28) {
#ifdef DEBUG
		bpf_printk("Not an AAAA query");
#endif
		return XDP_PASS;
	}

	// Parse the DNS answer
	struct dns_answer dns_answer;
	int dns_answer_result =
		parse_dns_answer(ctx, dns, query_length, &dns_answer);
	if (dns_answer_result < 0) {
#ifdef DEBUG
		bpf_printk("Error parsing DNS answer");
#endif
		return XDP_PASS;
	}

	// Check if there is an entry with that domain name in the map
	// If there is no entry, we stop here
	struct client_data *client_data;
	client_data = bpf_map_lookup_elem(&client_map, query.name);
	if (!client_data) {
#ifdef DEBUG
		bpf_printk("No entry found for %s\n", query.name);
#endif
		return XDP_PASS;
	}

#ifdef DEBUG
	bpf_printk("Found entry for %s\n", query.name);
#endif

	client_data->dstaddr = dns_answer.ipv6_address;
	bpf_map_update_elem(&client_map, query.name, client_data, BPF_EXIST);

#ifdef DEBUG
	bpf_printk("Updated IPv6 address for domain name: %s\n", query.name);
#endif

	return XDP_PASS;
}

static int parse_dns_answer(struct xdp_md *ctx, struct dns_hdr *dns_hdr,
							int query_length, struct dns_answer *a)
{
	void *data_end = (void *)(long)ctx->data_end;
	struct dns_answer *temp_a;

	// Calculate the pointer where the DNS answer begins in the packet
	temp_a = (struct dns_answer *)((__u8 *)dns_hdr + sizeof(struct dns_hdr) +
								   query_length);
	// Boundary check
	if ((void *)(temp_a) + sizeof(struct dns_answer) > data_end) {
#ifdef DEBUG
		bpf_printk("Error: boundary exceeded while parsing DNS answer");
#endif
		return -1;
	}
	memcpy(a, temp_a, sizeof(struct dns_answer));
	a->query_pointer = bpf_ntohs(a->query_pointer);
	a->record_type = bpf_ntohs(a->record_type);
	a->record_class = bpf_ntohs(a->record_class);
	a->ttl = bpf_ntohl(a->ttl);
	a->data_length = bpf_ntohs(a->data_length);

	return 0;
}

static int parse_dns_query(struct xdp_md *ctx, void *query_start,
						   struct dns_query *q)
{
	void *data_end = (void *)(long)ctx->data_end;
	__u16 i;
	void *cursor = query_start;
	int namepos = 0;

	// Fill dns_query.name with zero bytes
	// Not doing so will make the verifier complain when dns_query is used as a
	// key in bpf_map_lookup
	for (i = 0; i < MAX_DNS_NAME_LENGTH; i++) {
		q->name[i] = 0;
	}
	// Fill record_type and class with default values to satisfy verifier
	q->record_type = 0;
	q->record_class = 0;

	for (i = 0; i < MAX_DNS_NAME_LENGTH; i++) {
		// Boundary check of cursor. Verifier requires a +1 here.
		// Probably because we are advancing the pointer at the end of the loop
		if (cursor + 1 > data_end) {
#ifdef DEBUG
			bpf_printk("Error: boundary exceeded while parsing DNS query name");
#endif
			break;
		}

		// If separator is zero we've reached the end of the domain query
		if (*(char *)(cursor) == 0) {

			// We've reached the end of the query name.
			// This will be followed by 2x 2 bytes: the dns type and dns class.
			if (cursor + 5 > data_end) {
#ifdef DEBUG
				bpf_printk("Error: boundary exceeded while retrieving DNS "
						   "record type and class");
#endif
			} else {
				q->record_type = bpf_htons(*(__u16 *)(cursor + 1));
				q->record_class = bpf_htons(*(__u16 *)(cursor + 3));
			}

			// Return the bytecount of (namepos + current '0' byte + dns type +
			// dns class) as the query length.
			return namepos + 1 + 2 + 2;
		}

		// Read and fill data into struct
		q->name[namepos] = *(char *)(cursor);
		namepos++;
		cursor++;
	}

	return -1;
}

char _license[] SEC("license") = "GPL";