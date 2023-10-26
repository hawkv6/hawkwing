#ifndef __LIB_CLIENT_MAP_HELPERS_H
#define __LIB_CLIENT_MAP_HELPERS_H

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
#include "client_maps.h"
#include "ipproto.h"

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
	if (parse_ipproto_dstport(skb, ipv6, &dstport) < 0)
		return -1;

	*sid = bpf_map_lookup_elem(inner_map, &dstport);
	if (!*sid)
		return -1;

	return 0;
}

static __always_inline int client_get_sid_test(struct __sk_buff *skb,
										  struct ipv6hdr *ipv6,
										  struct sidlist_data **sidlist_data)
{
	__u32 *domain_id = bpf_map_lookup_elem(&client_reverse_map, &ipv6->daddr);
	if (!domain_id)
		return -1;

	struct bpf_elf_map *inner_map =
		bpf_map_lookup_elem(&client_outer_map, domain_id);
	if (!inner_map)
		return -1;

	__u16 dstport = 0;
	if (parse_ipproto_dstport(skb, ipv6, &dstport) < 0)
		return -1;

	*sidlist_data = bpf_map_lookup_elem(inner_map, &dstport);
	if (!*sidlist_data)
		return -1;

	return 0;
}

#endif