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

#include "lib/consts.h"
#include "lib/maps.h"

SEC("tc-egress")
int encap_egress(struct __sk_buff *skb)
{
	void *data_end = (void *)(long)skb->data_end;
	void *data = (void *)(long)skb->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);

	// Validate ethernet header
	if ((void *)(eth + 1) > data_end) {
#ifdef DEBUG
		bpf_printk("[tc-egress] invalid ethernet header\n");
#endif
		return TC_ACT_OK;
	}

	// Validate IPv6 header
	if ((void *)(ipv6 + 1) > data_end) {
#ifdef DEBUG
		bpf_printk("[tc-egress] invalid IPv6 header\n");
#endif
		return TC_ACT_OK;
	}

	// Check if it is an IPv6 packet
	if (eth->h_proto != bpf_htons(ETH_P_IPV6)) {
#ifdef DEBUG
		bpf_printk("[tc-egress] not IPv6 packet\n");
#endif
		return TC_ACT_OK;
	}

	// Check if it is a UDP or TCP packet
	if (ipv6->nexthdr != IPPROTO_UDP && ipv6->nexthdr != IPPROTO_TCP) {
#ifdef DEBUG
		bpf_printk("[tc-egress] not UDP or TCP packet\n");
#endif
		return TC_ACT_OK;
	}

    /*
	Check if there is an entry for the dstaddr in the client_reverse_map
	Value should be an __u32 domain_id
    */
	__u32 *domain_id = bpf_map_lookup_elem(&client_reverse_map, &ipv6->daddr);
	if (!domain_id) {
#ifdef DEBUG
		bpf_printk("[tc-egress] no entry in reverse map\n");
#endif
		return TC_ACT_OK;
	}

	/*
	Check if there is an inner_map for the domain_id in the client_outer_map	
	*/
	struct bpf_elf_map *inner_map = bpf_map_lookup_elem(&client_outer_map, domain_id);
	if (!inner_map) {
#ifdef DEBUG
		bpf_printk("[tc-egress] no inner map for domain_id\n");
#endif
		return TC_ACT_OK;
	}

	/*
	Check if there is an entry for __u16 80 in the innermap
	value is struct in6_addr[MAX_SEGMENT_LIST_ENTRIES]
	print the second in6_addr
	*/
	__u16 port = 80;
	struct in6_addr *segment_list = bpf_map_lookup_elem(inner_map, &port);
	if (!segment_list) {
#ifdef DEBUG
		bpf_printk("[tc-egress] no segment list for port 80\n");
#endif
		return TC_ACT_OK;
	}
	bpf_printk("[tc-egress] segment list for port 80: %pI6\n", segment_list);
	bpf_printk("[tc-egress] segment list for port 80: %pI6\n", segment_list + 1);




	bpf_printk("[tc-egress] FINISHED\n");
	return TC_ACT_OK;
}

char _license[] SEC("license") = "GPL";