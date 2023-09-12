#ifndef __LIB_ENCAP_H
#define __LIB_ENCAP_H

static __always_inline bool needs_encapsulation(__u32 addr, __u16 port)
{
	return false;
}

#endif