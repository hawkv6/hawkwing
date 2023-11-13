#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/pkt_cls.h>
#include <linux/seg6.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <linux/types.h>
#include <stdbool.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#include "lib/client_maps.h"
#include "lib/consts.h"
#include "lib/dns.h"
#include "lib/map_common.h"
#include "lib/srv6.h"
#include "lib/client_map_helpers.h"

#define memcpy __builtin_memcpy

char _license[] SEC("license") = "GPL";

SEC("xdp")
int client_ingress(struct xdp_md *ctx)
{
	void *data_end = (void *)(long)ctx->data_end;
	void *data = (void *)(long)ctx->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);
	struct udphdr *udp;
	struct dns_hdr *dns;
	struct srh *srh;

	if ((void *)(eth + 1) > data_end)
		goto pass;
	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		goto pass;
	if ((void *)(ipv6 + 1) > data_end)
		goto pass;

	switch (ipv6->nexthdr) {
		case IPPROTO_UDP:
			goto handle_dns;
		case IPPROTO_TCP:
			goto pass;
		case IPPROTO_ROUTING:
			goto handle_srh;
		default:
			goto pass;
	}

handle_dns:
	udp = (void *)(ipv6 + 1);
	dns = (void *)(udp + 1);

	if ((void *)(udp + 1) > data_end)
		goto pass;
	if (udp->source != bpf_htons(UDP_P_DNS))
		goto pass;

	struct dns_query query;
	struct dns_answer dns_answer;
	if (parsing_dns_answer(ctx, dns, &query, &dns_answer, data_end) < 0)
		goto pass;

	if (store_dns_tuple(&query, &dns_answer) < 0)
		goto pass;

	bpf_printk("[client-ingress] handled dns\n");
	goto pass;

handle_srh:
	srh = (struct srh *)(ipv6 + 1);
	if (srh_check_boundaries(srh, data_end) < 0)
		goto drop;
	if (remove_srh(ctx, data, data_end, srh) < 0)
		goto drop;

	bpf_printk("[client-ingress] srv6 packet processed\n");
	goto pass;

pass:
	return XDP_PASS;

drop:
	return XDP_DROP;
}

SEC("tc")
int client_egress(struct __sk_buff *skb)
{
	void *data_end = (void *)(long)skb->data_end;
	void *data = (void *)(long)skb->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);
	struct sidlist_data *sidlist_data;

	if ((void *)(eth + 1) > data_end)
		goto pass;
	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		goto pass;
	if ((void *)(ipv6 + 1) > data_end)
		goto pass;

	switch (ipv6->nexthdr) {
		case IPPROTO_UDP:
			goto handle_srh;
		case IPPROTO_TCP:
			goto handle_srh;
		default:
			goto pass;
	}

handle_srh:
	if (client_get_sid(skb, ipv6, &sidlist_data) < 0)
		goto pass;

	if (sidlist_data->sidlist_size == 0)
		goto pass;

	if (add_srh(skb, data, data_end, sidlist_data) < 0)
		goto drop;

	bpf_printk("[client-egress] srv6 packet send\n");
	goto pass;

pass:
	return TC_ACT_OK;

drop:
	return TC_ACT_SHOT;
}