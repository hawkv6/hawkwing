#ifndef __LIB_MAP_COMMON_H
#define __LIB_MAP_COMMON_H

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ipv6.h>
#include <linux/types.h>

#include "consts.h"

struct sid_lookup_value {
	__u32 sidlist_size;
	struct in6_addr sidlist[MAX_SEGMENTLIST_ENTRIES];
} __attribute__((packed));

#endif