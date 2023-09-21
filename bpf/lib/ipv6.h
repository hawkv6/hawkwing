#ifndef __LIB_IPV6_H
#define __LIB_IPV6_H

#include <linux/ipv6.h>

#define NEXTHDR_ROUTING 43 /* Routing header. */
#define NEXTHDR_NONE 59	   /* No next header. */

#define IPV6_SADDR_OFF offsetof(struct ipv6hdr, saddr)
#define IPV6_DADDR_OFF offsetof(struct ipv6hdr, daddr)

static __always_inline __be32 ipv6_pseudohdr_checksum(struct ipv6hdr *hdr,
						      __u8 next_hdr,
						      __u16 payload_len, __be32 sum)
{
	__be32 len = bpf_htonl((__u32)payload_len);
	__be32 nexthdr = bpf_htonl((__u32)next_hdr);

	sum = bpf_csum_diff(NULL, 0, &hdr->saddr, sizeof(struct in6_addr), sum);
	sum = bpf_csum_diff(NULL, 0, &hdr->daddr, sizeof(struct in6_addr), sum);
	sum = bpf_csum_diff(NULL, 0, &len, sizeof(len), sum);
	sum = bpf_csum_diff(NULL, 0, &nexthdr, sizeof(nexthdr), sum);

	return sum;
}

#endif