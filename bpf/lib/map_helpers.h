#ifndef __LIB_MAP_HELPERS_H
#define __LIB_MAP_HELPERS_H

#include "consts.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/types.h>

#include <linux/ipv6.h>
#include <linux/tcp.h>
#include <linux/udp.h>

#include "map_common.h"
#include "server_maps.h"
#include "client_maps.h"
#include "tcp.h"
#include "udp.h"

#define memset __builtin_memset
#define memcpy __builtin_memcpy

static __always_inline int
store_incoming_triple(struct xdp_md *ctx, struct ipv6hdr *ipv6, struct srh *srh)
{
	void *data_end = (void *)(long)ctx->data_end;
	int srh_len = srh_get_hdr_len(srh);
	int sidlist_size = srh_get_segment_list_len(srh);

	struct sid_lookup_value *value;
	__u32 temp_key = 0;

	value = bpf_map_lookup_elem(&server_temp_value_map, &temp_key);
	if (!value) {
		bpf_printk("ERROR: bpf_map_lookup_elem failed\n");
		return -1;
	}
	memset(value, 0, sizeof(struct server_lookup_value));

	for (int i = 0; i < MAX_SEGMENTLIST_ENTRIES; i++) {
		if (!(i < sidlist_size)) {
			break;
		}

		if ((void *)srh + sizeof(struct srh) +
				sizeof(struct in6_addr) * (i + 1) + 1 >
			data_end) {
			bpf_printk("ERROR: out of bounds\n");
			return -1;
		}

		memcpy(&value->sidlist[i], &srh->segments[i], sizeof(struct in6_addr));
	}
	value->sidlist_size = sidlist_size;

	struct server_lookup_key key = {};
	key.addr = ipv6->saddr;

	struct tcphdr *tcp;
	struct udphdr *udp;

	if (srh->next_hdr == IPPROTO_TCP) {
		if ((void *)srh + srh_len + sizeof(struct tcphdr) > data_end)
			return -1;
		tcp = (struct tcphdr *)((__u8 *)srh + srh_len);
		if ((void *)tcp + sizeof(struct tcphdr) <=
			data_end) { // Explicit bounds check
			key.port = bpf_ntohs(tcp->source);
		} else {
			return -1;
		}
	} else {
		if ((void *)srh + srh_len + sizeof(struct udphdr) > data_end)
			return -1;
		udp = (struct udphdr *)((__u8 *)srh + srh_len);
		if ((void *)udp + sizeof(struct udphdr) <=
			data_end) { // Explicit bounds check
			key.port = bpf_ntohs(udp->source);
		} else {
			return -1;
		}
	}

	if (sidlist_size > MAX_SEGMENTLIST_ENTRIES) {
		bpf_printk("ERROR: sidlist_size is out of bounds\n");
		return -1;
	}

	struct in6_addr reversed[MAX_SEGMENTLIST_ENTRIES];

	for (int i = 0; i < sidlist_size; ++i) {
		int j = sidlist_size - i - 1;
		if (j < 0 || j >= sidlist_size) {
			bpf_printk("ERROR: Reverse index out of bounds\n");
			return -1;
		}
		reversed[j] = value->sidlist[i];
	}
	memset(&reversed[sidlist_size - 1], 0, sizeof(struct in6_addr));

	for (int i = 1; i < sidlist_size + 1; ++i) {
		if (i >= MAX_SEGMENTLIST_ENTRIES) {
			bpf_printk("ERROR: Shift index out of bounds\n");
			return -1;
		}
		value->sidlist[i] = reversed[i - 1];
	}

	value->sidlist[0] = ipv6->saddr;

	if (bpf_map_update_elem(&server_lookup_map, &key, value, BPF_ANY) < 0)
		return -1;

	return 0;
}

static __always_inline int client_get_sid(struct __sk_buff *skb,
										  struct ipv6hdr *ipv6,
										  struct in6_addr **sid)
{
	__u32 *domain_id = bpf_map_lookup_elem(&client_reverse_map, &ipv6->daddr);
	if (!domain_id)
		return -1;

	struct bpf_elf_map *inner_map =
		bpf_map_lookup_elem(&client_outer_map, domain_id);
	if (!inner_map)
		return -1;

	__u16 dstport = 0;
	switch (ipv6->nexthdr) {
		case IPPROTO_TCP: {
			struct tcphdr *tcp = (struct tcphdr *)(ipv6 + 1);
			if (parse_tcp_hdr(skb, tcp, &dstport) < 0)
				return -1;
			break;
		}
		case IPPROTO_UDP: {
			struct udphdr *udp = (struct udphdr *)(ipv6 + 1);
			if (parse_udp_hdr(skb, udp, &dstport) < 0)
				return -1;
			break;
		}
		default:
			return -1;
	}
	*sid = bpf_map_lookup_elem(inner_map, &dstport);
	if (!*sid)
		return -1;

	return 0;
}

#endif