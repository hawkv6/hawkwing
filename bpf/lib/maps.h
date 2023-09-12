#ifndef __LIB_MAPS_H
#define __LIB_MAPS_H
#include "client.h"
#include "consts.h"
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/types.h>

// struct {
// 	__uint(type, BPF_MAP_TYPE_LRU_HASH);
// 	__uint(max_entries, MAX_MAP_ENTRIES);
// 	__type(key, char[MAX_DNS_NAME_LEN]);
// 	__type(value, struct client_data);
// 	__uint(pinning, LIBBPF_PIN_BY_NAME);
// } client_map SEC(".maps");

struct client_inner_map {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, __u16); // dstport
	__type(value, struct in6_addr[MAX_SEGMENTLIST_ENTRIES]);
	// __uint(pinning, LIBBPF_PIN_BY_NAME);
} client_inner_map SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_HASH_OF_MAPS);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, __u32);
	// __type(key, char[MAX_DNS_NAME_LEN]);
	// __type(value, __u32);  // index for client_inner_map fd
	__uint(pinning, LIBBPF_PIN_BY_NAME);
	__array(values, struct client_inner_map);
} client_outer_map SEC(".maps") = {
	.values = { [0] = &client_inner_map },
};

struct {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, struct in6_addr);
	// __type(value, char[MAX_DNS_NAME_LEN]);
	__type(value, __u32); // key for lookup in client_outer_map
	__uint(pinning, LIBBPF_PIN_BY_NAME);
} client_reverse_map SEC(".maps");

struct {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, char[MAX_DNS_NAME_LEN]);
	__type(value, __u32); // key to store in client_reverse_map
} client_lookup_map SEC(".maps");

#endif