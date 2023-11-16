package maps

import (
	"fmt"
	"net"
	"strings"
)

type SidListData struct {
	SidlistSize uint32
	Sidlist     [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
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

	if domain == "." {
		return result, fmt.Errorf("root domain not allowed")
	}

	labels := strings.Split(domain, ".")
	offset := 0

	if domain == labels[0] {
		return result, fmt.Errorf("domain name have to consist of at least two labels")
	}

	for _, label := range labels {
		result[offset] = byte(len(label))
		offset++

		copy(result[offset:], label)
		offset += len(label)
	}

	result[offset] = 0

	return result, nil
}

// Ipv6ToInet6 takes an IPv6 address in string format and returns a struct containing
// the IPv6 address in byte format. The returned struct is intended for use as a key
// in an eBPF map.
//
// Parameters:
//   - ipv6Addr: An IPv6 address in string format.
//
// Returns:
//   - A struct containing a 128-bit IPv6 address in byte format.
func Ipv6ToInet6(ipv6Addr string) struct{ In6U struct{ U6Addr8 [16]uint8 } } {
	ipv6 := net.ParseIP(ipv6Addr)
	var result struct{ In6U struct{ U6Addr8 [16]uint8 } }
	copy(result.In6U.U6Addr8[:], ipv6.To16())
	return result
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
func SidToInet6Sid(sidList []string) ([10]struct{ In6U struct{ U6Addr8 [16]uint8 } }, error) {
	if len(sidList) > 9 {
		return [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{}, fmt.Errorf("too many SIDs")
	}

	var result [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }

	// Leave [0] empty, start from 1
	for i, sid := range sidList {
		ipv6 := net.ParseIP(sid)
		// Reverse the order of the input list while inserting into the result
		copy(result[len(sidList)-i].In6U.U6Addr8[:], ipv6.To16())
	}
	return result, nil
}

// GenerateSidLookupValue takes a slice of IPv6 Segment IDs (SIDs) in string format and returns
// a SidListData struct containing the IPv6 SIDs in reversed order. The returned struct is intended
// for use as a value in an eBPF map.
//
// Parameters:
//   - sidList: A slice of IPv6 SIDs in string format.
//
// Returns:
//   - A SidListData struct containing the IPv6 SIDs in reversed order.
func GenerateSidLookupValue(sidList []string) (SidListData, error) {
	if len(sidList) == 0 {
		return SidListData{
			SidlistSize: 0,
		}, nil
	}
	convertedSidList, err := SidToInet6Sid(sidList)
	if err != nil {
		return SidListData{}, err
	}
	result := SidListData{
		SidlistSize: uint32(len(sidList) + 1), // +1 for the empty one
		Sidlist:     convertedSidList,
	}
	return result, nil
}
