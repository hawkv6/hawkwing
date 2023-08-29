#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>
#include <bpf/bpf_helpers.h>

SEC("filter")
int parse_dns(struct __sk_buff *skb) {

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";