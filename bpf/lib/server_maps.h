#ifndef __LIB_SERVER_MAPS_H
#define __LIB_SERVER_MAPS_H

#include "consts.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/types.h>

struct server_lookup_map {
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, MAX_MAP_ENTRIES);
    __type(key, struct server_lookup_key);
    __type(value, struct in6_addr[MAX_SEGMENTLIST_ENTRIES]);
    __uint(pinning, LIBBPF_PIN_BY_NAME);
} server_lookup_map SEC(".maps");

struct server_lookup_key {
    struct in6_addr addr;
    __u16 port;
} __attribute__((packed));

static __always_inline int store_incoming_triple(struct ipv6hdr *ipv6, struct srv6_hdr *srh, void *trans_hdr)
{
    int sidlist_size = srv6_get_segment_list_len(srh);
    struct in6_addr sidlist[sidlist_size];
    memcpy(sidlist, srh + 1, sidlist_size * sizeof(struct in6_addr));


}

#endif