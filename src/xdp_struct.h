#ifndef __XDP_STRUCTS_H
#define __XDP_STRUCTS_H

#include <linux/types.h>
#include <linux/in6.h>

struct client_data
{
    __u16 dstport; // load from config file
    struct in6_addr dstaddr; // load from dns reply
    struct in6_addr segments[MAX_SEGMENTLIST_ENTRIES];
};


#endif