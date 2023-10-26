#ifndef __LIB_IPPROTO_H
#define __LIB_IPPROTO_H

#include <linux/bpf.h>
#include <linux/ipv6.h>
#include <linux/tcp.h>
#include <linux/udp.h>

#include "tcp.h"
#include "udp.h"

static int parse_ipproto_dstport(struct __sk_buff *skb, struct ipv6hdr *ipv6,
                                 __u16 *dstport)
{
    void *data_end = (void *)(long)skb->data_end;

    if ((void *)(ipv6 + 1) > data_end) {
        return -1;
    }

    switch (ipv6->nexthdr) {
        case IPPROTO_TCP: {
            struct tcphdr *tcp = (struct tcphdr *)(ipv6 + 1);
            if (parse_tcp_hdr(skb, tcp, dstport) < 0)
                return -1;
            break;
        }
        case IPPROTO_UDP: {
            struct udphdr *udp = (struct udphdr *)(ipv6 + 1);
            if (parse_udp_hdr(skb, udp, dstport) < 0)
                return -1;
            break;
        }
        default:
            return -1;
    }

    return 0;
}

#endif