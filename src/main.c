#include <linux/bpf.h>
#include <stdbool.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/udp.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/seg6.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

#include "xdp_map.h"
int counter = 0;

SEC("xdp")
int intercept_dns(struct xdp_md *ctx) {
    // bpf_printk("Hello World %d", counter);
    // counter++; 
    // return XDP_PASS;

    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    // __u32 probe_key = XDP_PASS;
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;

    if (eth->h_proto != bpf_htons(ETH_P_IPV6))
        return XDP_PASS;

    struct ipv6hdr *ipv6 = (void *)(eth + 1);
    if ((void *)(ipv6 + 1) > data_end)
        return XDP_PASS;

    // Validate IPv6 header: Check if the Next Header is UDP (17)
    if (ipv6->nexthdr != 17)
        bpf_printk("udp parsed");
        return XDP_PASS;
        
    bpf_printk("working until here");

    // struct udphdr *udp = (void *)(ipv6 + 1);
    // if ((void *)(udp + 1) > data_end)
    //     return XDP_PASS;

    // // Validate UDP header: Check if the UDP destination port is 53 (DNS)
    // if (udp->dest != bpf_htons(53))
    //     bpf_printk("dns parsed");
    //     return XDP_PASS;

    // // Define the DNS header structure
    // struct dns_header {
    //     __be16 id;
    //     __be16 flags;
    //     __be16 qdcount;
    //     __be16 ancount;
    //     __be16 nscount;
    //     __be16 arcount;
    // };

    // struct dns_header *dns = (void *)(udp + 1);
    // if ((void *)(dns + 1) > data_end)
    //     return XDP_PASS;

    // // Only proceed if there's at least one answer
    // if (dns->ancount == 0)
    //     return XDP_PASS;

    // // Here, you would normally skip the question section to reach the answer section
    // // For this example, we're skipping directly to where we'd expect an IPv6 address.
    // // This is a simplification.
    // __u8 *dns_data = (__u8 *)(dns + 1);
    // if ((void *)(dns_data + 16) > data_end)  // Check if enough data exists
    //     return XDP_PASS;

    // // Assuming the IPv6 address is directly in the DNS answer section
    // // This is a gross simplification; you'll need proper parsing logic here.
    // struct in6_addr *ipv6_address = (struct in6_addr *)(dns_data);

    // // Do something with the IPv6 address, e.g., store it in an eBPF map
    // // ...

    // bpf_printk("IPv6 address: %pI6\n", ipv6_address);

    return XDP_PASS;
}

// SEC("lwt_encap")
// int encap_sr(struct __sk_buff *skb) {
//     return BPF_OK;
// }

// SEC("xdp")
// int reverse_sid(struct xdp_md *ctx) {
//     return XDP_PASS;
// }

char _license[] SEC("license") = "GPL";