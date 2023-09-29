#ifndef __LIB_TCP_H
#define __LIB_TCP_H

#include <linux/tcp.h>

#define memcpy __builtin_memcpy

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

#endif