#ifndef __XDP_MAPS_H
#define __XDP_MAPS_H
#include <linux/types.h>
#include <linux/if_ether.h>
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include "xdp_consts.h"
#include "xdp_struct.h"
struct
{
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, MAX_MAP_ENTRIES);
    __type(key, char[MAX_DOMAIN_NAME_LEN]);
    __type(value, struct client_data);
} client_map SEC(".maps");

struct
{
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    __uint(max_entries, MAX_MAP_ENTRIES);
    __type(key, int);
    __type(value, char[MAX_DOMAIN_NAME_LEN]);
} test_map SEC(".maps");

#endif