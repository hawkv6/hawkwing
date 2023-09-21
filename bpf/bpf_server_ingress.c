#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/pkt_cls.h>
#include <linux/tcp.h>
#include <linux/udp.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#include "lib/consts.h"
#include "lib/srv6.h"

#define memcpy __builtin_memcpy

char _license[] SEC("license") = "GPL";

SEC("xdp")
int filter_ingress(struct xdp_md *ctx)
{
	void *data_end = (void *)(long)ctx->data_end;
	void *data = (void *)(long)ctx->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);
	struct srv6_hdr *srh = (struct srv6_hdr *)(ipv6 + 1);

	if ((void *)(eth + 1) > data_end || eth->h_proto != bpf_htons(ETH_P_IPV6))
		return XDP_PASS;

	if ((void *)(ipv6 + 1) > data_end || ipv6->nexthdr != IPPROTO_ROUTING)
		return XDP_PASS;

	if (srv6_check_boundaries(srh, data_end) < 0)
		return XDP_DROP;

	int sidlist_size = srv6_get_segment_list_len(srh);
	struct in6_addr srv6_sidlist[sidlist_size];
	memcpy(srv6_sidlist, (struct in6_addr *)(srh + 1),
		   sizeof(struct in6_addr) * sidlist_size);

	if (srv6_remove_srh(ctx, data, data_end, srh) < 0)
		return XDP_DROP;

	bpf_printk("[xdp | server] packet processed\n");

	return XDP_PASS;
}