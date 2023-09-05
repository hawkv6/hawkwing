#ifndef __CLIENT_H
#define __CLIENT_H

#include <linux/in6.h>
#include <linux/types.h>

struct client_data {
	__u16 dstport;			 // load from config file
	struct in6_addr dstaddr; // load from dns reply
	struct in6_addr segments[MAX_SEGMENTLIST_ENTRIES];
};

#endif