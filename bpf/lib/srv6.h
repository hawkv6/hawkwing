#ifndef __LIB_SRV6_H
#define __LIB_SRV6_H

#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>

#include "consts.h"
#include "map_common.h"

#define SRV6_NEXT_HDR 43	/* Routing header. */
#define SRV6_HDR_EXT_LEN 0	/* Routing header extension length. */
#define SRV6_ROUTING_TYPE 4 /* SRv6 routing type. */

#define memcpy __builtin_memcpy
#define memset __builtin_memset

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
struct srh {
	__u8 next_hdr;
	__u8 hdr_ext_len;
	__u8 routing_type;
	__u8 segments_left;
	__u8 last_entry;
	__u8 flags;
	__u16 tag;
	struct in6_addr segments[0];
} __attribute__((packed));

static __always_inline int srh_get_hdr_len(struct srh *hdr)
{
	return (hdr->hdr_ext_len + 1) * 8;
}

static __always_inline int srh_get_segment_list_len(struct srh *hdr)
{
	return hdr->last_entry + 1;
}

static __always_inline int srh_check_boundaries(struct srh *hdr, void *end)
{
	if ((void *)hdr + sizeof(struct srh) > end ||
		(void *)hdr + srh_get_hdr_len(hdr) > end)
		return -1;
	return 0;
}

static __always_inline int remove_srh(struct xdp_md *ctx, void *data,
									  void *data_end, struct srh *hdr)
{
	struct ethhdr eth_copy;
	struct ipv6hdr ipv6_copy;
	memcpy(&eth_copy, data, sizeof(struct ethhdr));
	memcpy(&ipv6_copy, data + sizeof(struct ethhdr), sizeof(struct ipv6hdr));

	int srh_len = srh_get_hdr_len(hdr);
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

static __always_inline int add_srh(struct __sk_buff *skb, void *data,
								   void *data_end, struct in6_addr *sids,
								   __u8 sidlist_size)
{
	struct ipv6hdr *ipv6 = data + sizeof(struct ethhdr);
	if (data + sizeof(struct ethhdr) + sizeof(struct ipv6hdr) > data_end)
		return -1;

	int hdr_ext_len =
		(sizeof(struct srh) + sizeof(struct in6_addr) * sidlist_size - 8) / 8;
	struct srh srh = {
		.next_hdr = ipv6->nexthdr,
		.hdr_ext_len = hdr_ext_len,
		.routing_type = SRV6_ROUTING_TYPE,
		.segments_left = sidlist_size - 1,
		.last_entry = sidlist_size - 1,
	};

	// use percpu map to store sidlist temporary
	struct in6_addr *percpu_sidlist;
	__u32 percpu_sidlist_key = 0;
	percpu_sidlist = bpf_map_lookup_elem(&percpu_sidlist_map, &percpu_sidlist_key);
	if (!percpu_sidlist)
		return -1;
	memset(percpu_sidlist, 0, sizeof(struct in6_addr) * MAX_SEGMENTLIST_ENTRIES);

	// memcpy(percpu_sidlist, sids, sizeof(struct in6_addr) * MAX_SEGMENTLIST_ENTRIES);
	for (int i = 0; i < sidlist_size; i++) {
		if (!(i < sidlist_size)) {
			break;
		}
		memcpy(&percpu_sidlist[i], &sids[i], sizeof(struct in6_addr));
	}

	// memcpy(&sids[0], &ipv6->daddr, sizeof(struct in6_addr));
	memcpy(&percpu_sidlist[0], &ipv6->daddr, sizeof(struct in6_addr));
	ipv6->payload_len =
		bpf_htons(bpf_ntohs(ipv6->payload_len) + sizeof(struct srh) +
				  sizeof(struct in6_addr) * sidlist_size);
	ipv6->nexthdr = SRV6_NEXT_HDR;
	// memcpy(&ipv6->daddr, &sids[sidlist_size - 1], sizeof(struct in6_addr));
	memcpy(&ipv6->daddr, &percpu_sidlist[sidlist_size - 1], sizeof(struct in6_addr));

	if (bpf_skb_adjust_room(
			skb, sizeof(struct srh) + sizeof(struct in6_addr) * sidlist_size,
			BPF_ADJ_ROOM_NET, 0) < 0)
		return -1;

	if (bpf_skb_store_bytes(skb, sizeof(struct ethhdr) + sizeof(struct ipv6hdr),
							&srh, sizeof(struct srh), 0) < 0)
		return -1;

	// if (bpf_skb_store_bytes(
	// 		skb,
	// 		sizeof(struct ethhdr) + sizeof(struct ipv6hdr) + sizeof(struct srh),
	// 		sids, sizeof(struct in6_addr) * sidlist_size, 0) < 0)
	// 	return -1;
	if (bpf_skb_store_bytes(
			skb,
			sizeof(struct ethhdr) + sizeof(struct ipv6hdr) + sizeof(struct srh),
			percpu_sidlist, sizeof(struct in6_addr) * sidlist_size, 0) < 0)
		return -1;

	return 0;
}

static __always_inline int add_srh_server(struct __sk_buff *skb, void *data,
								   void *data_end, struct sidlist_data *sidlist_data)
{
	struct ipv6hdr *ipv6 = data + sizeof(struct ethhdr);
	if (data + sizeof(struct ethhdr) + sizeof(struct ipv6hdr) > data_end)
		return -1;

	int hdr_ext_len =
		(sizeof(struct srh) + sizeof(struct in6_addr) * sidlist_data->sidlist_size - 8) / 8;

	struct srh srh = {
		.next_hdr = ipv6->nexthdr,
		.hdr_ext_len = hdr_ext_len,
		.routing_type = SRV6_ROUTING_TYPE,
		.segments_left = sidlist_data->sidlist_size - 1,
		.last_entry = sidlist_data->sidlist_size - 1,
	};

	// use percpu map to store sidlist temporary
	struct in6_addr *percpu_sidlist;
	__u32 percpu_sidlist_key = 0;
	percpu_sidlist = bpf_map_lookup_elem(&percpu_sidlist_map, &percpu_sidlist_key);
	if (!percpu_sidlist)
		return -1;

	memset(percpu_sidlist, 0, sizeof(struct in6_addr) * MAX_SEGMENTLIST_ENTRIES);

	memcpy(percpu_sidlist, sidlist_data->sidlist, sizeof(struct in6_addr) * MAX_SEGMENTLIST_ENTRIES);

	memcpy(&percpu_sidlist[0], &ipv6->daddr, sizeof(struct in6_addr));
	ipv6->payload_len =
		bpf_htons(bpf_ntohs(ipv6->payload_len) + sizeof(struct srh) +
				  sizeof(struct in6_addr) * sidlist_data->sidlist_size);
	ipv6->nexthdr = SRV6_NEXT_HDR;

	__u16 last_index = sidlist_data->sidlist_size - 1;
	if (last_index < MAX_SEGMENTLIST_ENTRIES) {
		memcpy(&ipv6->daddr, &percpu_sidlist[last_index], sizeof(struct in6_addr));
	} else {
		return -1;
	}

	if (bpf_skb_adjust_room(
			skb, sizeof(struct srh) + sizeof(struct in6_addr) * sidlist_data->sidlist_size,
			BPF_ADJ_ROOM_NET, 0) < 0)
		return -1;

	if (bpf_skb_store_bytes(skb, sizeof(struct ethhdr) + sizeof(struct ipv6hdr),
							&srh, sizeof(struct srh), 0) < 0)
		return -1;

	__u16 i;
	__u32 offset = sizeof(struct ethhdr) + sizeof(struct ipv6hdr) + sizeof(struct srh);

	for (i = 0; i < sidlist_data->sidlist_size; i++) {
		// check is needed for verifier
		if (i >= MAX_SEGMENTLIST_ENTRIES) {
			break;
		}
		if (bpf_skb_store_bytes(skb, offset, &percpu_sidlist[i], sizeof(struct in6_addr), 0) < 0) {
			return -1;
		}
		offset += sizeof(struct in6_addr);
	}

	return 0;
}

#endif
