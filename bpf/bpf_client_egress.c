#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/pkt_cls.h>
#include <linux/tcp.h>
#include <linux/types.h>
#include <linux/udp.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#include "lib/client_maps.h"
#include "lib/consts.h"
#include "lib/srv6.h"
#include "lib/map_helpers.h"
// #include "lib/tcp.h"
// #include "lib/udp.h"

#define memcpy __builtin_memcpy
#define memset __builtin_memset

char _license[] SEC("license") = "GPL";

SEC("tc-egress")
int client_egress(struct __sk_buff *skb)
{
	void *data_end = (void *)(long)skb->data_end;
	void *data = (void *)(long)skb->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);
	struct in6_addr *segment_list;

	if ((void *)(eth + 1) > data_end)
		goto pass;
	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		goto pass;
	if ((void *)(ipv6 + 1) > data_end)
		goto pass;

	switch(ipv6->nexthdr) {
		case IPPROTO_UDP:
			goto handle_srh;
		case IPPROTO_TCP:
			goto handle_srh;
		default:
			goto pass;
	}

handle_srh:
	if (client_get_sid(skb, ipv6, &segment_list) < 0)
		goto pass;

	// TODO make it somehow dynamic
	__u8 num_sids = 3;

	if (add_srh(skb, data, data_end, segment_list, num_sids) < 0)
		goto drop;

	bpf_printk("[client-egress] srv6 packet send\n");
	goto pass;

pass:
	return TC_ACT_OK;

drop:
	return TC_ACT_SHOT;
}