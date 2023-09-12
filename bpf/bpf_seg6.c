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
    Check if there is an entry for the dstaddr in the reverse map. 
    If there is, we can continue.
    If there isn't, we can let the packet pass.
    */
    char *domain_name;
    domain_name = bpf_map_lookup_elem(&client_reverse_map, &ipv6->daddr);
	if (!domain_name) {
#ifdef DEBUG
		bpf_printk("[tc-egress] no entry in reverse map\n");
#endif
		return TC_ACT_OK;
	}
    bpf_printk("[tc-egress] domain_name: %s\n", domain_name);

    /*
    Check if there is an entry for the dstport in the inner map.
    If there is, we can continue.
    If there isn't, we can let the packet pass.
    */ 
	__u32 *inner_map_pointer;
	inner_map_pointer = bpf_map_lookup_elem(&client_outer_map, domain_name);
	if (!inner_map_pointer) {
#ifdef DEBUG
		bpf_printk("[tc-egress] no entry in inner map\n");
#endif
		return TC_ACT_OK;
	}

	/*
	Check if the destination port of tcp is available as key in the inner map.
	*/
	// int mapfd;
	// mapfd = bpf_obj_get(inner_map_pointer);
	struct bpf_map_def *inner_map;
	inner_map = bpf_map_lookup_elem(&client_outer_map, domain_name);
	if (!inner_map) {
#ifdef DEBUG
		bpf_printk("[tc-egress] no entry in inner map\n");
#endif
		return TC_ACT_OK;
	}

// 	__u16 dstport = 80;
// 	__u16 *key = &dstport;
// 	struct in6_addr (*segments)[MAX_SEGMENTLIST_ENTRIES];
// 	segments = bpf_map_lookup_elem(&inner_map, key);
// 	if (!segments) {
// #ifdef DEBUG
// 		bpf_printk("[tc-egress] no entry in inner map\n");
// #endif
// 		return TC_ACT_OK;
// 	}


	bpf_printk("[tc-egress] FINISHED\n");
	return TC_ACT_OK;
}

char _license[] SEC("license") = "GPL";