#ifndef __LIB_SERVER_MAPS_H
#define __LIB_SERVER_MAPS_H

#include "consts.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ipv6.h>
#include <linux/types.h>

struct server_lookup_map {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, struct server_lookup_key);
	__type(value, struct server_lookup_value);
	__uint(pinning, LIBBPF_PIN_BY_NAME);
} server_lookup_map SEC(".maps");

struct server_lookup_key {
	struct in6_addr addr;
	__u16 port;
} __attribute__((packed));

struct server_lookup_value {
	struct in6_addr sidlist[MAX_SEGMENTLIST_ENTRIES];
	int sidlist_size;
} __attribute__((packed));

struct server_temp_value_map {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, struct server_lookup_value);
} server_temp_value_map SEC(".maps");

struct server_temp_sid_map {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, struct in6_addr[MAX_SEGMENTLIST_ENTRIES]);
} server_temp_sid_map SEC(".maps");
#endif