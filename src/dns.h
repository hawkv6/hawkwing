#ifndef __DNS_H
#define __DNS_H

#define MAX_DNS_NAME_LENGTH 256

struct dns_hdr
{
    __u16 transaction_id;
    __u8 rd : 1;      //Recursion desired
    __u8 tc : 1;      //Truncated
    __u8 aa : 1;      //Authoritive answer
    __u8 opcode : 4;  //Opcode
    __u8 qr : 1;      //Query/response flag
    __u8 rcode : 4;   //Response code
    __u8 cd : 1;      //Checking disabled
    __u8 ad : 1;      //Authenticated data
    __u8 z : 1;       //Z reserved bit
    __u8 ra : 1;      //Recursion available
    __u16 q_count;    //Number of questions
    __u16 ans_count;  //Number of answer RRs
    __u16 auth_count; //Number of authority RRs
    __u16 add_count;  //Number of resource RRs
};

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