#ifndef __LIB_UDP_H
#define __LIB_UDP_H

#include <linux/udp.h>

#define memcpy __builtin_memcpy

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

#endif