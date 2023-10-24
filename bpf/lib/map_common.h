#ifndef __LIB_MAP_COMMON_H
#define __LIB_MAP_COMMON_H

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ipv6.h>
#include <linux/types.h>

#include "consts.h"

struct sidlist_data {
	__u32 sidlist_size;
	struct in6_addr sidlist[MAX_SEGMENTLIST_ENTRIES];
} __attribute__((packed));

struct percpu_sidlist_map {
	__uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
	__uint(max_entries, 1);
	__type(key, __u32);
	__type(value, struct in6_addr[MAX_SEGMENTLIST_ENTRIES]);
} percpu_sidlist_map SEC(".maps");

#endif