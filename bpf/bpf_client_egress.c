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

#define memcpy __builtin_memcpy
#define memset __builtin_memset

char _license[] SEC("license") = "GPL";

static int parse_tcp_hdr(struct __sk_buff *skb, struct tcphdr *tcp,
						 __u16 *dstport);
static int parse_udp_hdr(struct __sk_buff *skb, struct udphdr *udp,
						 __u16 *dstport);

SEC("tc-egress")
int client_egress(struct __sk_buff *skb)
{
	void *data_end = (void *)(long)skb->data_end;
	void *data = (void *)(long)skb->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);

	if ((void *)(eth + 1) > data_end)
		goto pass;
	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		goto pass;
	if ((void *)(ipv6 + 1) > data_end)
		goto pass;

	if (ipv6->nexthdr != IPPROTO_UDP && ipv6->nexthdr != IPPROTO_TCP)
		return TC_ACT_OK;

	/*
	Check if there is an entry for the dstaddr in the client_reverse_map
	Value should be an __u32 domain_id
	*/
	__u32 *domain_id = bpf_map_lookup_elem(&client_reverse_map, &ipv6->daddr);
	if (!domain_id)
		return TC_ACT_OK;

	/*
	Check if there is an inner_map for the domain_id in the client_outer_map
	*/
	struct bpf_elf_map *inner_map =
		bpf_map_lookup_elem(&client_outer_map, domain_id);
	if (!inner_map)
		return TC_ACT_OK;

	__u16 dstport = 0;
	if (ipv6->nexthdr == IPPROTO_UDP) {
		struct udphdr *udp = (struct udphdr *)(ipv6 + 1);
		if (parse_udp_hdr(skb, udp, &dstport) < 0) {
			return TC_ACT_OK;
		}
	} else if (ipv6->nexthdr == IPPROTO_TCP) {
		struct tcphdr *tcp = (struct tcphdr *)(ipv6 + 1);
		if (parse_tcp_hdr(skb, tcp, &dstport) < 0) {
			return TC_ACT_OK;
		}
	}

	struct in6_addr *segment_list = bpf_map_lookup_elem(inner_map, &dstport);
	if (!segment_list)
		return TC_ACT_OK;

	// TODO make it somehow dynamic
	__u8 num_sids = 3;

	if (add_srh(skb, data, data_end, segment_list, num_sids) < 0)
		goto drop;
	// return TC_ACT_OK;

	bpf_printk("[tc-egress] srv6 packet send\n");

	return TC_ACT_OK;

pass:
	bpf_printk("[client-egress] pass\n");
	return TC_ACT_OK;

drop:
	bpf_printk("[client-egress] drop\n");
	return TC_ACT_SHOT;
}

static int parse_tcp_hdr(struct __sk_buff *skb, struct tcphdr *tcp,
						 __u16 *dstport)
{
	void *data_end = (void *)(long)skb->data_end;

	if ((void *)(tcp + 1) > data_end) {
		return -1;
	}
	__u16 temp = bpf_ntohs(tcp->dest);
	memcpy(dstport, &temp, sizeof(__u16));

	return 0;
}

static int parse_udp_hdr(struct __sk_buff *skb, struct udphdr *udp,
						 __u16 *dstport)
{
	void *data_end = (void *)(long)skb->data_end;

	if ((void *)(udp + 1) > data_end) {
		return -1;
	}
	__u16 temp = bpf_ntohs(udp->dest);
	memcpy(dstport, &temp, sizeof(__u16));

	return 0;
}
