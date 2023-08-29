#ifndef __XDP_STRUCTS_H
#define __XDP_STRUCTS_H

#include <linux/types.h>
#include <linux/if_ether.h>
#include <linux/in6.h>
#include "xdp_consts.h"

struct intent_service_data
{
    char domain_name[MAX_DOMAIN_NAME_LEN];
    struct in6_addr srcaddr_v6;
    struct in6_addr dstaddr_v6;
    __u16 dstport;
    struct in6_addr segments[MAX_SEGMENTLIST_ENTRIES];
};

#endif