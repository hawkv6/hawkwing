#ifndef __LIB_CLIENT_MAPS_H
#define __LIB_CLIENT_MAPS_H

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/types.h>

#include "consts.h"
#include "map_common.h"

struct client_inner_map {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, __u16); // dstport
	__type(value, struct sidlist_data);
} client_inner_map SEC(".maps");

struct client_outer_map {
	__uint(type, BPF_MAP_TYPE_HASH_OF_MAPS);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__uint(pinning, LIBBPF_PIN_BY_NAME);
	__type(key, __u32); // domain_name id
	__array(values, struct client_inner_map);
} client_outer_map SEC(".maps"); 

struct client_reverse_map {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, struct in6_addr);
	__type(value, __u32); // domain_name id
	__uint(pinning, LIBBPF_PIN_BY_NAME);
} client_reverse_map SEC(".maps");

struct client_lookup_map {
	__uint(type, BPF_MAP_TYPE_LRU_HASH);
	__uint(max_entries, MAX_MAP_ENTRIES);
	__type(key, char[MAX_DNS_NAME_LEN]);
	__type(value, __u32); // domain_name id
	__uint(pinning, LIBBPF_PIN_BY_NAME);
} client_lookup_map SEC(".maps");

#endif