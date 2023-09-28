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

#include "lib/client_maps.h"
#include "lib/map_common.h"
#include "lib/consts.h"
#include "lib/dns.h"

#define memcpy __builtin_memcpy

SEC("xdp")
int intercept_dns(struct xdp_md *ctx)
{
	void *data_end = (void *)(long)ctx->data_end;
	void *data = (void *)(long)ctx->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);

	if ((void *)(eth + 1) > data_end)
		return XDP_PASS;

	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		return XDP_PASS;

	if ((void *)(ipv6 + 1) > data_end)
		return XDP_PASS;

	if (ipv6->nexthdr != IPPROTO_UDP)
		return XDP_PASS;

	// struct dns_hdr *dns;
	// if (is_dns_reply(ctx, &dns) < 0)
	// 	bpf_printk("[xdp] not a dns reply\n");
	// 	return XDP_PASS;

	struct udphdr *udp = (void *)(ipv6 + 1);

	if ((void *)(udp + 1) > data_end)
		return XDP_PASS;

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
	if (query_length < 1)
		return XDP_PASS;

	// Parse the DNS answer
	struct dns_answer dns_answer;
	int dns_answer_result =
		parse_dns_answer(ctx, dns, query_length, &dns_answer);
	if (dns_answer_result < 0) 
		return XDP_PASS;

	/*
	Check if there is an entry for the extraced domain in the client_lookup_map
	When yes print the value which is the id.
	When no go ahead
	*/
	__u32 *domain_id;
	domain_id = bpf_map_lookup_elem(&client_lookup_map, query.name);
	if (!domain_id) {
#ifdef DEBUG
		bpf_printk("[xdp] no entry found for %s\n in client_lookup_map",
				   query.name);
#endif
		return XDP_PASS;
	}
#ifdef DEBUG
	bpf_printk("[xdp] found entry for %s\n in client_lookup_map with id: %d\n",
			   query.name, *domain_id);
#endif

	bpf_map_update_elem(&client_reverse_map, &dns_answer.ipv6_address,
						domain_id, BPF_ANY);
#ifdef DEBUG
	bpf_printk(
		"[xdp] updated reverse map with ipv6 address for domain id: %d\n",
		*domain_id);
#endif

	return XDP_PASS;
}


char _license[] SEC("license") = "GPL";