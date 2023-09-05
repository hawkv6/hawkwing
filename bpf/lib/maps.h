#ifndef __LIB_MAPS_H
#define __LIB_MAPS_H
#include "client.h"
#include "consts.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/types.h>
struct {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, char[MAX_DNS_NAME_LEN]);
	__type(value, struct client_data);
} client_map SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, int);
	__type(value, struct in6_addr);
} test_map SEC(".maps");

#endif