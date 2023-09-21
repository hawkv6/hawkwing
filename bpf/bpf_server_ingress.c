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
#include "lib/ipv6.h"
#include "lib/csum.h"
#include "lib/consts.h"

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

    int srv6_len = 0;
    srv6_len = srh->hdr_ext_len * 8 + 8;
    int sidlist_size = 0;
    sidlist_size = (srv6_len - sizeof(struct srv6_hdr)) / 16;
    __u8 srh_next_hdr = srh->next_hdr;
    if ((void *)(srh + 1) > data_end || (void *)(srh + 1) + srv6_len > data_end)
        return XDP_DROP;

    struct ethhdr eth_copy;
    struct ipv6hdr ipv6_copy;
    struct srv6_hdr srh_copy;
    memcpy(&eth_copy, eth, sizeof(struct ethhdr));
    memcpy(&ipv6_copy, ipv6, sizeof(struct ipv6hdr));
    memcpy(&srh_copy, srh, sizeof(struct srv6_hdr));

    struct in6_addr srv6_sidlist[MAX_SEGMENTLIST_ENTRIES] = {0};
    struct in6_addr *sidlist_start = (struct in6_addr *)(srh + 1);
    // memcpy(srv6_sidlist, srh + 1, sizeof(struct in6_addr) * sidlist_size);

    if (sidlist_size > MAX_SEGMENTLIST_ENTRIES)
        return XDP_DROP;

    sidlist_size = 3;
    for (int i = 0; i < sidlist_size; i++) {
        if ((void *)(sidlist_start + i + 1) > data_end) {
            return XDP_DROP; // or some other handling
        }
        
        // Copy manually to satisfy the verifier
        for(int j = 0; j < 4; ++j) {
            srv6_sidlist[i].s6_addr32[j] = sidlist_start[i].s6_addr32[j];
        }
    }

    
    struct in6_addr old_dst;
    if (1 > MAX_SEGMENTLIST_ENTRIES)
        return XDP_DROP;
    
    old_dst = srv6_sidlist[2];


    ipv6->daddr = old_dst;
    __wsum original_pseudo_chk = ipv6_pseudohdr_checksum(ipv6, IPPROTO_ROUTING, bpf_ntohs(ipv6->payload_len), 0);


    struct in6_addr new_dst = ipv6->daddr;



    if (bpf_xdp_adjust_head(ctx, srv6_len) < 0)
        return XDP_DROP;

    data = (void *)(long)ctx->data;
    data_end = (void *)(long)ctx->data_end;

    eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_DROP;
    memcpy(eth, &eth_copy, sizeof(struct ethhdr));

    ipv6 = (struct ipv6hdr *)(eth + 1);
    if ((void *)(ipv6 + 1) > data_end)
        return XDP_DROP;
    memcpy(ipv6, &ipv6_copy, sizeof(struct ipv6hdr));

    ipv6->payload_len = bpf_htons(bpf_ntohs(ipv6->payload_len) - srv6_len);

    if (srh_next_hdr == IPPROTO_TCP)
        ipv6->nexthdr = IPPROTO_TCP;
    else if (srh_next_hdr == IPPROTO_UDP)
        ipv6->nexthdr = IPPROTO_UDP;
    else
        return XDP_DROP;


    struct tcphdr *tcp = (struct tcphdr *)(ipv6 + 1);

    if ((void *)(tcp + 1) > data_end)
        return XDP_DROP;


    ipv6_copy.daddr = old_dst;
    ipv6_copy.payload_len = bpf_htons(bpf_ntohs(ipv6_copy.payload_len) + srv6_len);
    // Calculate pseudo-header checksum for original packet (with SRH)
    // __wsum original_pseudo_chk = ipv6_pseudohdr_checksum(&ipv6_copy, IPPROTO_ROUTING, bpf_ntohs(ipv6_copy.payload_len), 0);

    // Calculate pseudo-header checksum for modified packet (without SRH)
    __wsum new_pseudo_chk = ipv6_pseudohdr_checksum(ipv6, IPPROTO_TCP, bpf_ntohs(ipv6->payload_len), 0);

    // Compute the checksum difference
    __wsum csum_diff = bpf_csum_diff(&original_pseudo_chk, sizeof(original_pseudo_chk),
                                     &new_pseudo_chk, sizeof(new_pseudo_chk), 0);

    // Update the original TCP checksum
    // tcp->check = ~csum_fold(csum_add(~csum_unfold(tcp->check), csum_diff));

    // tcp->check = 0;



    bpf_printk("[xdp | server] packet processed\n");

    return XDP_PASS;
}