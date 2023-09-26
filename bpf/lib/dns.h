#ifndef __LIB_DNS_H
#define __LIB_DNS_H

#include <linux/in6.h>

#include "consts.h"

#define memcpy __builtin_memcpy

#define UDP_P_DNS 53

struct dns_hdr {
	__u16 transaction_id;
	__u8 rd : 1;
	__u8 tc : 1;
	__u8 aa : 1;
	__u8 opcode : 4;
	__u8 qr : 1;
	__u8 rcode : 4;
	__u8 cd : 1;
	__u8 ad : 1;
	__u8 z : 1;
	__u8 ra : 1;
	__u16 q_count;
	__u16 ans_count;
	__u16 auth_count;
	__u16 add_count;
} __attribute__((packed));

struct dns_query {
	__u16 record_type;
	__u16 record_class;
	char name[MAX_DNS_NAME_LEN];
} __attribute__((packed));

struct dns_answer {
	__u16 query_pointer;
	__u16 record_type;
	__u16 record_class;
	__u32 ttl;
	__u16 data_length;
	struct in6_addr ipv6_address;
} __attribute__((packed));

static int parse_dns_answer(struct xdp_md *ctx, struct dns_hdr *dns_hdr,
							int query_length, struct dns_answer *a)
{
	void *data_end = (void *)(long)ctx->data_end;
	struct dns_answer *temp_a;

	// Calculate the pointer where the DNS answer begins in the packet
	temp_a = (struct dns_answer *)((__u8 *)dns_hdr + sizeof(struct dns_hdr) +
								   query_length);

	if ((void *)(temp_a) + sizeof(struct dns_answer) > data_end)
		return -1;

	memcpy(a, temp_a, sizeof(struct dns_answer));
	a->query_pointer = bpf_ntohs(a->query_pointer);
	a->record_type = bpf_ntohs(a->record_type);
	a->record_class = bpf_ntohs(a->record_class);
	a->ttl = bpf_ntohl(a->ttl);
	a->data_length = bpf_ntohs(a->data_length);

	return 0;
}

static int parse_dns_query(struct xdp_md *ctx, void *query_start,
						   struct dns_query *q)
{
	void *data_end = (void *)(long)ctx->data_end;
	__u16 i;
	void *cursor = query_start;
	int namepos = 0;

	// Fill dns_query.name with zero bytes
	// Not doing so will make the verifier complain when dns_query is used as a
	// key in bpf_map_lookup
	for (i = 0; i < MAX_DNS_NAME_LEN; i++) {
		q->name[i] = 0;
	}
	// Fill record_type and class with default values to satisfy verifier
	q->record_type = 0;
	q->record_class = 0;

	for (i = 0; i < MAX_DNS_NAME_LEN; i++) {
		// Boundary check of cursor. Verifier requires a +1 here.
		// Probably because we are advancing the pointer at the end of the loop
		if (cursor + 1 > data_end) {
#ifdef DEBUG
			bpf_printk("Error: boundary exceeded while parsing DNS query name");
#endif
			break;
		}

		// If separator is zero we've reached the end of the domain query
		if (*(char *)(cursor) == 0) {

			// We've reached the end of the query name.
			// This will be followed by 2x 2 bytes: the dns type and dns class.
			if (cursor + 5 > data_end) {
#ifdef DEBUG
				bpf_printk("Error: boundary exceeded while retrieving DNS "
						   "record type and class");
#endif
			} else {
				q->record_type = bpf_htons(*(__u16 *)(cursor + 1));
				q->record_class = bpf_htons(*(__u16 *)(cursor + 3));
			}

			// Return the bytecount of (namepos + current '0' byte + dns type +
			// dns class) as the query length.
			return namepos + 1 + 2 + 2;
		}

		// Read and fill data into struct
		q->name[namepos] = *(char *)(cursor);
		namepos++;
		cursor++;
	}

	return -1;
}

static int is_dns_reply(struct xdp_md *ctx, struct dns_hdr *dns_hdr)
{
	void *data_end = (void *)(long)ctx->data_end;

	if (ctx->data + sizeof(struct ethhdr) + sizeof(struct ipv6hdr) +
			sizeof(struct udphdr) >
		data_end)
		return 0;
	
	struct udphdr *udp = ctx->data + sizeof(struct ethhdr) +
						 sizeof(struct ipv6hdr);
	
	if (udp->source != bpf_htons(UDP_P_DNS))
		return -1;

	dns_hdr = ctx->data + sizeof(struct ethhdr) + sizeof(struct ipv6hdr) +
			   sizeof(struct udphdr);

	if ((void *)(dns_hdr) + sizeof(struct dns_hdr) > data_end)
		return -1;

	return 0;
}

#endif