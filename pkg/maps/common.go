package maps

import (
	"fmt"
	"net"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

var (
	pinnedMapOptions = ebpf.MapOptions{
		PinPath: bpf.BpffsRoot,
	}
)

// InnerMapData is a struct containing the data stored in the inner map of the
// client/server map. The inner map is a LRU hash map with a key size of 2 bytes
// and a value size of 160 bytes. The key is the destination port of the packet
// and the value is a struct containing the IPv6 Segment IDs (SIDs).
type InnerMapData struct {
	DstPort  uint16
	Segments [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}

// FormatDNSName takes a domain name in string format and returns a byte array
// containing the domain name in DNS format. The returned byte array is intended
// for use as a key in an eBPF map.
//
// Parameters:
//   - domain: A domain name in string format.
//
// Returns:
//   - A byte array containing the domain name in DNS format.
//   - An error if the domain name is invalid.
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

// SidToInet6Sid takes a slice of IPv6 Segment IDs (SIDs) in string format and returns an array
// of 10 structs containing the IPv6 SIDs in reversed order. The returned array is intended
// for use in constructing the SID List for the Segment Routing Header (SRH) in IPv6 packets.
//
// Specifically, the function performs the following:
//  1. The first index [0] of the returned array is left empty. This is intended for the last
//     hop, which should be sourced from the IPv6 destination address of the original packet.
//  2. The remaining IPv6 SIDs are reversed, meaning that the last SID in the input slice becomes
//     the first SID in the output array, the second to last becomes the second, and so on.
//
// Parameters:
//   - sidList: A slice of IPv6 SIDs in string format.
//
// Returns:
//   - An array of 10 structs, each containing a 128-bit IPv6 address in byte format.
func SidToInet6Sid(sidList []string) [10]struct{ in6U struct{ u6Addr8 [16]uint8 } } {
	var result [10]struct{ in6U struct{ u6Addr8 [16]uint8 } }
	// Leave [0] empty, start from 1
	for i, sid := range sidList {
		if i >= 9 {
			break // Max 9 addresses plus the empty one
		}
		ipv6 := net.ParseIP(sid)
		// Reverse the order of the input list while inserting into the result
		copy(result[len(sidList)-i].in6U.u6Addr8[:], ipv6.To16())
	}
	return result
}
