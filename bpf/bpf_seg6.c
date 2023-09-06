#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/pkt_cls.h>
#include <bpf/bpf_helpers.h>

SEC("tc-egress")
int encap_egress(struct __sk_buff *skb)
{
    bpf_printk("[tc] egress got packet\n");

    return TC_ACT_OK;
}

char _license[] SEC("license") = "GPL";