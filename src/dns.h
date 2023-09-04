#ifndef __DNS_H
#define __DNS_H

#include <linux/in6.h>

#define MAX_DNS_NAME_LENGTH 256

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
	char name[MAX_DNS_NAME_LENGTH];
} __attribute__((packed));

struct dns_answer {
	__u16 query_pointer;
	__u16 record_type;
	__u16 record_class;
	__u32 ttl;
	__u16 data_length;
	struct in6_addr ipv6_address;
} __attribute__((packed));

#endif