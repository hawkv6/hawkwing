#ifndef __LIB_SRV6_H
#define __LIB_SRV6_H

#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include "consts.h"

#define SRV6_NEXT_HDR 43	/* Routing header. */
#define SRV6_HDR_EXT_LEN 0	/* Routing header extension length. */
#define SRV6_ROUTING_TYPE 4 /* SRv6 routing type. */

#define memcpy __builtin_memcpy

/*
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Next Header   |  Hdr Ext Len  | Routing Type  | Segments Left |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |  Last Entry   |     Flags     |              Tag              |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[0] (128-bit IPv6 address)             |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |                                                               |
							   ...
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[n] (128-bit IPv6 address)             |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 //                                                             //
 //         Optional Type Length Value objects (variable)       //
 //                                                             //
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
struct srv6_hdr {
	__u8 next_hdr;
	__u8 hdr_ext_len;
	__u8 routing_type;
	__u8 segments_left;
	__u8 last_entry;
	__u8 flags;
	__u16 tag;
	// variable length of segment list entries
	// length has to be get from the client_inner_map values
} __attribute__((packed));

static __always_inline int srv6_get_hdr_len(struct srv6_hdr *hdr)
{
	return (hdr->hdr_ext_len + 1) * 8;
}

static __always_inline int srv6_get_segment_list_len(struct srv6_hdr *hdr)
{
	return (srv6_get_hdr_len(hdr) - sizeof(struct srv6_hdr)) / 16;
}

static __always_inline int srv6_check_boundaries(struct srv6_hdr *hdr,
												 void *end)
{
	if ((void *)hdr + sizeof(struct srv6_hdr) > end ||
		(void *)hdr + srv6_get_hdr_len(hdr) > end)
		return -1;
	return 0;
}

static __always_inline int srv6_remove_srh(struct xdp_md *ctx, void *data,
										   void *data_end, struct srv6_hdr *hdr)
{
	struct ethhdr eth_copy;
	struct ipv6hdr ipv6_copy;
	memcpy(&eth_copy, data, sizeof(struct ethhdr));
	memcpy(&ipv6_copy, data + sizeof(struct ethhdr), sizeof(struct ipv6hdr));

	int srh_len = srv6_get_hdr_len(hdr);
	__u8 srh_next_hdr = hdr->next_hdr;

	if (bpf_xdp_adjust_head(ctx, srh_len) < 0)
		return -1;

	data_end = (void *)(long)ctx->data_end;
	data = (void *)(long)ctx->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = data + sizeof(struct ethhdr);

	if (data + sizeof(struct ethhdr) + sizeof(struct ipv6hdr) > data_end)
		return -1;

	memcpy(eth, &eth_copy, sizeof(struct ethhdr));
	memcpy(ipv6, &ipv6_copy, sizeof(struct ipv6hdr));

	ipv6->payload_len = bpf_htons(bpf_ntohs(ipv6->payload_len) - srh_len);

	if (srh_next_hdr == IPPROTO_TCP)
		ipv6->nexthdr = IPPROTO_TCP;
	else if (srh_next_hdr == IPPROTO_UDP)
		ipv6->nexthdr = IPPROTO_UDP;
	else
		ipv6->nexthdr = IPPROTO_NONE;

	return 0;
}

#endif
