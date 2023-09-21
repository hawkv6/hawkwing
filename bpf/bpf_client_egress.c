#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/in6.h>
#include <linux/ip.h>
#include <linux/ipv6.h>
#include <linux/pkt_cls.h>
#include <linux/tcp.h>
#include <linux/types.h>
#include <linux/udp.h>

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>

// #include "bpf_encap_seg6.c"
#include "lib/client_maps.h"
#include "lib/consts.h"
#include "lib/srv6.h"
#include "lib/ipv6.h"
#include "lib/csum.h"

#define memcpy __builtin_memcpy

char _license[] SEC("license") = "GPL";

static int parse_tcp_hdr(struct __sk_buff *skb, struct tcphdr *tcp,
						 __u16 *dstport);
static int parse_udp_hdr(struct __sk_buff *skb, struct udphdr *udp,
						 __u16 *dstport);

SEC("tc-egress")
int filter_egress(struct __sk_buff *skb)
{
	void *data_end = (void *)(long)skb->data_end;
	void *data = (void *)(long)skb->data;
	struct ethhdr *eth = data;
	struct ipv6hdr *ipv6 = (struct ipv6hdr *)(eth + 1);

	if ((void *)(eth + 1) > data_end)
		return TC_ACT_OK;
		
	if (eth->h_proto != bpf_htons(ETH_P_IPV6))
		return TC_ACT_OK;

	if ((void *)(ipv6 + 1) > data_end)
		return TC_ACT_OK;

	if (ipv6->nexthdr != IPPROTO_UDP && ipv6->nexthdr != IPPROTO_TCP)
		return TC_ACT_OK;

	__wsum original_pseudo_chk = ipv6_pseudohdr_checksum(ipv6, IPPROTO_TCP, bpf_ntohs(ipv6->payload_len), 0);
	struct ipv6hdr old_ipv6;
	memcpy(&old_ipv6, ipv6, sizeof(struct ipv6hdr));
	/*
	Check if there is an entry for the dstaddr in the client_reverse_map
	Value should be an __u32 domain_id
	*/
	__u32 *domain_id = bpf_map_lookup_elem(&client_reverse_map, &ipv6->daddr);
	if (!domain_id)
		return TC_ACT_OK;

	/*
	Check if there is an inner_map for the domain_id in the client_outer_map
	*/
	struct bpf_elf_map *inner_map =
		bpf_map_lookup_elem(&client_outer_map, domain_id);
	if (!inner_map)
		return TC_ACT_OK;

	__u16 dstport = 0;
	if (ipv6->nexthdr == IPPROTO_UDP) {
		struct udphdr *udp = (struct udphdr *)(ipv6 + 1);
		if (parse_udp_hdr(skb, udp, &dstport) < 0) {
			return TC_ACT_OK;
		}
	} else if (ipv6->nexthdr == IPPROTO_TCP) {
		struct tcphdr *tcp = (struct tcphdr *)(ipv6 + 1);
		if (parse_tcp_hdr(skb, tcp, &dstport) < 0) {
			return TC_ACT_OK;
		}
	}

	struct in6_addr *segment_list = bpf_map_lookup_elem(inner_map, &dstport);
	if (!segment_list)
		return TC_ACT_OK;

	__u8 num_sids = 3;

	struct srv6_hdr srv6;
	__builtin_memset(&srv6, 0, sizeof(struct srv6_hdr));
	// fill up the SRH
	srv6.next_hdr = 6;
	// int hdr_ext_len = ((16 * num_sids) + 8 - 1) / 8;
	int hdr_ext_len = (sizeof(struct srv6_hdr) + sizeof(struct in6_addr) * num_sids - 8)/8;
	srv6.hdr_ext_len = hdr_ext_len;
	srv6.routing_type = SRV6_ROUTING_TYPE;
	srv6.segments_left = num_sids - 1;
	srv6.last_entry = num_sids - 1;

	// copy ipv6 dstaddr to sids[0]
	memcpy(&segment_list[0], &ipv6->daddr, sizeof(struct in6_addr));

	__u16 old_payload_len = bpf_ntohs(ipv6->payload_len);
	__u16 new_payload_len = old_payload_len + sizeof(struct srv6_hdr) +
							sizeof(struct in6_addr) * num_sids;
	ipv6->payload_len = bpf_htons(new_payload_len);
	// ipv6 next header is set to SRH
	ipv6->nexthdr = SRV6_NEXT_HDR;
	// ipv6 destination address is set to the last entry in sid list
	memcpy(&ipv6->daddr, &segment_list[num_sids - 1], sizeof(struct in6_addr));

	// adjust room for the SRH
	if (bpf_skb_adjust_room(skb,
							sizeof(struct srv6_hdr) +
								sizeof(struct in6_addr) * num_sids,
							BPF_ADJ_ROOM_NET, 0) < 0)
		return TC_ACT_OK;

	// store the SRH after the IPv6 header
	if (bpf_skb_store_bytes(
			skb, sizeof(struct ethhdr) + sizeof(struct ipv6hdr), &srv6,
			sizeof(struct srv6_hdr),
			0) < 0)
		return TC_ACT_OK;

	// store the sids after the SRH
	if (bpf_skb_store_bytes(
			skb, sizeof(struct ethhdr) + sizeof(struct ipv6hdr) +
					 sizeof(struct srv6_hdr),
			segment_list, sizeof(struct in6_addr) * num_sids, 0) < 0) 
		return TC_ACT_OK;

	data_end = (void *)(long)skb->data_end;
	data = (void *)(long)skb->data;
	eth = data;
	ipv6 = (struct ipv6hdr *)(eth + 1);

	if ((void *)(eth + 1) > data_end || (void *)(ipv6 + 1) > data_end)
		return TC_ACT_OK;


	// tcp header starts after ipv6 + srh plus segment list hdr_ext_len
	struct tcphdr *tcp = (struct tcphdr *)(ipv6 + 1) + hdr_ext_len;
	if ((void *)(tcp + 1) > data_end)
		return TC_ACT_OK;

	__wsum new_pseudo_chk = ipv6_pseudohdr_checksum(ipv6, ipv6->nexthdr, bpf_ntohs(ipv6->payload_len), 0);
	__wsum csum_diff = bpf_csum_diff(&original_pseudo_chk, sizeof(original_pseudo_chk), &new_pseudo_chk, sizeof(new_pseudo_chk), 0);
	// tcp->check = ~csum_fold(csum_add(~csum_unfold(tcp->check), csum_diff));

	// if (bpf_l4_csum_replace(skb, sizeof(struct ethhdr) + sizeof(struct ipv6hdr) + (hdr_ext_len + 1) * 8, original_pseudo_chk, new_pseudo_chk, 0) < 0)
	// 	return TC_ACT_OK;

	__wsum wsum;
	wsum = ipv6_pseudohdr_checksum(ipv6, IPPROTO_TCP, ipv6->payload_len, 0);

	// tcp->check = csum_fold(bpf_csum_diff(NULL, 0, &ipv6, sizeof(ipv6),
	// 				bpf_csum_diff(NULL, 0, &old_ipv6,
	// 					sizeof(old_ipv6),
	// 						bpf_csum_diff(NULL, 0, tcp,
	// 							sizeof(tcp), wsum))));
	
	// bpf_l4_csum_replace(skb, offsetof(struct tcphdr, check), &old_ipv6.daddr, &ipv6->daddr, sizeof(ipv6->daddr));
	// bpf_l4_csum_replace(skb, offsetof(struct tcphdr, check), &old_ipv6.nexthdr, &ipv6->nexthdr, sizeof(ipv6->nexthdr));


	bpf_printk("[tc-egress] srv6 packet send\n");

	return TC_ACT_OK;
}

static int parse_tcp_hdr(struct __sk_buff *skb, struct tcphdr *tcp,
						 __u16 *dstport)
{
	void *data_end = (void *)(long)skb->data_end;

	if ((void *)(tcp + 1) > data_end) {
		return -1;
	}
	__u16 temp = bpf_ntohs(tcp->dest);
	memcpy(dstport, &temp, sizeof(__u16));

	return 0;
}

static int parse_udp_hdr(struct __sk_buff *skb, struct udphdr *udp,
						 __u16 *dstport)
{
	void *data_end = (void *)(long)skb->data_end;

	if ((void *)(udp + 1) > data_end) {
		return -1;
	}
	__u16 temp = bpf_ntohs(udp->dest);
	memcpy(dstport, &temp, sizeof(__u16));

	return 0;
}
