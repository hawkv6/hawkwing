#ifndef __LIB_SRV6_H
#define __LIB_SRV6_H

#include "consts.h"

#define SRV6_NEXT_HDR 43	/* Routing header. */
#define SRV6_HDR_EXT_LEN 0	/* Routing header extension length. */
#define SRV6_ROUTING_TYPE 4 /* SRv6 routing type. */

/*
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 | Next Header   |  Hdr Ext Len  | Routing Type  | Segments Left |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |  Last Entry   |     Flags     |              Tag              |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[0] (128-bit IPv6 address)             |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |                                                               |
							   ...
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                                                               |
 |            Segment List[n] (128-bit IPv6 address)             |
 |                                                               |
 |                                                               |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 //                                                             //
 //         Optional Type Length Value objects (variable)       //
 //                                                             //
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/
struct srv6_hdr {
	__u8 next_hdr;
	__u8 hdr_ext_len;
	__u8 routing_type;
	__u8 segments_left;
	__u8 last_entry;
	__u8 flags;
	__u16 tag;
	// variable length of segment list entries
	// length has to be get from the client_inner_map values
} __attribute__((packed));

#endif

/*
Steps:
- allocate new space with bpf_skb_adjust_room with the length of the complete
srh + sid list length
- insert srh after ipv6 header before all subsequent headers
- update ipv6 header with new payload length and next header
- recalculate checksums for tcp and udp -> have a look at bpf_l4_csum_replace
- forward packet to next hop
*/