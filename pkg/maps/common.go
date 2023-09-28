package maps

import (
	"encoding/binary"
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

// TODO - check of remove
type In6Addr struct {
	In6U struct{ U6Addr8 [16]uint8 }
}

// TODO - check of remove
type SidLookupValue struct {
	SidlistSize uint32
	Sidlist     [10]In6Addr
}

// TODO - check of remove
func (s *SidLookupValue) MarshalBinary() ([]byte, error) {
	const ipv6Size = 16
	const maxEntries = 10
	size := maxEntries*ipv6Size + 4 // 4 bytes for sidListSize
	buf := make([]byte, size)

	offset := 0

	for _, sid := range s.Sidlist {
		for _, b := range sid.In6U.U6Addr8 {
			buf[offset] = b
			offset++
		}
	}

	binary.LittleEndian.PutUint32(buf[offset:], s.SidlistSize)

	return buf, nil
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
// func SidToInet6Sid(sidList []string) [10]struct{ In6U struct{ U6Addr8 [16]uint8 } } {
// 	var result [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
// 	// Leave [0] empty, start from 1
// 	for i, sid := range sidList {
// 		if i >= 9 {
// 			break // Max 9 addresses plus the empty one
// 		}
// 		ipv6 := net.ParseIP(sid)
// 		// Reverse the order of the input list while inserting into the result
// 		copy(result[len(sidList)-i].In6U.U6Addr8[:], ipv6.To16())
// 	}
// 	return result
// }

func SidToInet6Sid(sidList []string) [10]In6Addr {
	var result [10]In6Addr
	// Leave [0] empty, start from 1
	for i, sid := range sidList {
		if i >= 9 {
			break // Max 9 addresses plus the empty one
		}
		ipv6 := net.ParseIP(sid)
		// Reverse the order of the input list while inserting into the result
		copy(result[len(sidList)-i].In6U.U6Addr8[:], ipv6.To16())
	}
	return result
}

// TODO - check of remove
func GenerateSidLookupValue(sidList []string) *SidLookupValue {
	// var result SidLookupValue
	result := SidLookupValue{
		SidlistSize: uint32(len(sidList)),
		Sidlist:     SidToInet6Sid(sidList),
	}
	// result.SidlistSize = uint32(len(sidList))
	// result.Sidlist = SidToInet6Sid(sidList)
	return &result
}
