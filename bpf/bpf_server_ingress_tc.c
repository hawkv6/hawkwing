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

#include "lib/srv6.h"

#define memcpy __builtin_memcpy

char _license[] SEC("license") = "GPL";

static __always_inline int pull_and_recast_data(struct __sk_buff *skb, void **data, void **data_end, __u32 size) 
{
    if (bpf_skb_pull_data(skb, size))
        return -1;

    *data_end = (void *)(long)skb->data_end;
    *data = (void *)(long)skb->data;
    return 0;
}

SEC("tc-ingress")
int filter_ingress(struct __sk_buff *skb)
{
	void *data_end = (void *)(long)skb->data_end;
	void *data = (void *)(long)skb->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);
    struct srv6_hdr *srh = (struct srv6_hdr *)(ipv6 + 1);

	if ((void *)(eth + 1) > data_end)
		return TC_ACT_OK;

	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		return TC_ACT_OK;

	if ((void *)(ipv6 + 1) > data_end) {
        if (pull_and_recast_data(skb, &data, &data_end, sizeof(struct ethhdr) + sizeof(struct ipv6hdr)))
            return TC_ACT_OK;
        
        eth = data;
        ipv6 = (struct ipv6hdr *)(eth + 1);

        if ((void *)(eth + 1) > data_end)
            return TC_ACT_OK;
        
        if ((void *)(ipv6 + 1) > data_end)
            return TC_ACT_OK;
        
    }

    if (ipv6->nexthdr != IPPROTO_ROUTING)
        return TC_ACT_OK;

    if ((void *)(srh + 1) > data_end) {
        if (pull_and_recast_data(skb, &data, &data_end, sizeof(struct ethhdr) + sizeof(struct ipv6hdr) + sizeof(struct srv6_hdr)))
            return TC_ACT_OK;

        eth = data;
        ipv6 = (struct ipv6hdr *)(eth + 1);
        srh = (struct srv6_hdr *)(ipv6 + 1);

        if ((void *)(eth + 1) > data_end)
            return TC_ACT_OK;
        
        if ((void *)(ipv6 + 1) > data_end)
            return TC_ACT_OK;

        if ((void *)(srh + 1) > data_end)
            return TC_ACT_OK;
    }


    // verify srh and srv6_len bytes of data are present
    // if ((void *)(srh + 1) + srv6_len > data_end) {
    //     if (pull_and_recast_data(skb, &data, &data_end, sizeof(struct ethhdr) + sizeof(struct ipv6hdr) + sizeof(struct srv6_hdr) + srv6_len))
    //         return TC_ACT_OK;
        
    //     eth = data;
    //     ipv6 = (struct ipv6hdr *)(eth + 1);
    //     srh = (struct srv6_hdr *)(ipv6 + 1);

    //     if ((void *)(eth + 1) > data_end)
    //         return TC_ACT_OK;
        
    //     if ((void *)(ipv6 + 1) > data_end)
    //         return TC_ACT_OK;

    //     if ((void *)(srh + 1) > data_end)
    //         return TC_ACT_OK;

    //     if ((void *)(srh + 1) + srv6_len > data_end)
    //         return TC_ACT_OK;
    // }

    // ipv6->nexthdr = IPPROTO_TCP;

    // if (bpf_skb_adjust_room(skb, -srv6_len, BPF_ADJ_ROOM_NET, 0) < 0) {
    //     bpf_printk("[server | ingress] bpf_skb_adjust_room failed\n");
    //     return TC_ACT_SHOT;
    // }

    bpf_printk("[server | ingress] srv6 packet received\n");

    return TC_ACT_OK;
}