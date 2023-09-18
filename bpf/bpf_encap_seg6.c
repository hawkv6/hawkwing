
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
#include "lib/srv6.h"

#define memcpy __builtin_memcpy

static int encap_seg6(struct __sk_buff *skb, struct ipv6hdr *ipv6,
					  struct in6_addr *sids, __u8 num_sids)
{
	void *data_end = (void *)(long)skb->data_end;

	struct srv6_hdr srv6 = {};
	// fill up the SRH
	srv6.next_hdr = SRV6_NEXT_HDR;
	srv6.hdr_ext_len = SRV6_HDR_EXT_LEN;
	srv6.routing_type = SRV6_ROUTING_TYPE;
	srv6.segments_left = num_sids - 1;
	srv6.last_entry = num_sids - 1;

	// copy ipv6 dstaddr to sids[0]
	memcpy(&sids[0], &ipv6->daddr, sizeof(struct in6_addr));

	// ipv6 destination address is set to the last entry in sid list
	memcpy(&ipv6->daddr, &sids[num_sids - 1], sizeof(struct in6_addr));

	// adjust room for the SRH
	if (bpf_skb_adjust_room(skb,
							sizeof(struct srv6_hdr) +
								sizeof(struct in6_addr) * (num_sids + 1),
							BPF_ADJ_ROOM_NET, 0) < 0) {
		return -1;
	}

	// store the SRH after the IPv6 header
	if (bpf_skb_store_bytes(
			skb, sizeof(struct ethhdr) + sizeof(struct ipv6hdr), &srv6,
			sizeof(struct srv6_hdr) + sizeof(struct in6_addr) + (num_sids + 1),
			0) < 0) {
		return -1;
	}

	// forward the packet
	return 0;
}

/*
[0] destination address of original ipv6 packet
[1]
...
[n] first SID which has to be the destination address of the encapsulated packet

*/