package bpf

import (
	"fmt"
	"net"
	"strings"
)

// FormatDNSName takes a domain name and converts it into a dns name.
// This function is used to convert the domain name into a format that can be used as a key in the BPF map.
func FormatDNSName(domain string) ([256]byte, error) {
	var result [256]byte
	labels := strings.Split(domain, ".")
	offset := 0

	for _, label := range labels {
		if len(label) == 0 {
			return result, fmt.Errorf("empty label detected")
		}

		// Write the label length
		result[offset] = byte(len(label))
		offset++

		// Write the label itself
		for i := 0; i < len(label); i++ {
			result[offset] = label[i]
			offset++
		}
	}

	// Append zero byte to indicate end of the domain name
	result[offset] = 0
	offset++

	return result, nil
}

// SidToInet6Sid takes a list of IPv6 addresses and converts them into a list of inet6 addresses
// This function is used to convert the list of IPv6 addresses into a format that can be used as a value in the BPF map.
func SidToInet6Sid(sidList []string) [10]struct{ in6U struct{ u6Addr8 [16]uint8 } } {
	var result [10]struct{ in6U struct{ u6Addr8 [16]uint8 } }
	for i, sid := range sidList {
		ipv6 := net.ParseIP(sid)
		// var inner struct{ in6U struct{ u6Addr8 [16]uint8 } }
		// copy(inner.in6U.u6Addr8[:], ipv6.To16())
		// result = append(result, inner)
		copy(result[i].in6U.u6Addr8[:], ipv6.To16())
	}
	return result
}
