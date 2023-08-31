#ifndef __DNS_H
#define __DNS_H

#define MAX_DNS_NAME_LENGTH 256

struct dns_header {
    __u16 id;
    __u16 flags;
    __u16 qdcount;
    __u16 ancount;
    __u16 nscount;
    __u16 arcount;
};

struct dns_query {
    __u16 record_type;
    __u16 record_class;
    char name[MAX_DNS_NAME_LENGTH];
};

struct dns_response {
    __u16 query_pointer;
    __u16 record_type;
    __u16 record_class;
    __u32 ttl;
    __u16 data_length;
};

struct a_record {
    struct in6_addr ipv6_addr;
    __u32 ttl;
};

#endif