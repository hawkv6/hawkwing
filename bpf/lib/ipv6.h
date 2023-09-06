#ifndef __LIB_IPV6_H
#define __LIB_IPV6_H

#include <linux/ipv6.h>

#define NEXTHDR_ROUTING     43      /* Routing header. */
#define NEXTHDR_NONE        59      /* No next header. */

#define IPV6_SADDR_OFF offsetof(struct ipv6hdr, saddr)
#define IPV6_DADDR_OFF offsetof(struct ipv6hdr, daddr)

#endif