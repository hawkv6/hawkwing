package bpf

import (
	"fmt"
	"net"
	"strings"

	"github.com/cilium/ebpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -no-global-types -cc $BPF_CLANG -cflags $BPF_CFLAGS xdp ../../bpf/bpf_dns.c -- -I../../bpf
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -no-global-types -cc $BPF_CLANG -cflags $BPF_CFLAGS tc ../../bpf/bpf_seg6.c -- -I../../bpf

type ClientData struct {
	DstPort  uint16
	_        [2]byte  // padding have to be added manually
	Dstaddr  [16]byte // empty for initializing
	Segments [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}

func ReadBpfObjects(ops *ebpf.CollectionOptions) (*xdpObjects, error) {
	obj := &xdpObjects{}
	err := loadXdpObjects(obj, ops)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func ReadTcObjects(ops *ebpf.CollectionOptions) (*tcObjects, error) {
	obj := &tcObjects{}
	err := loadTcObjects(obj, ops)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func InitializeBpfMap(m *ebpf.Map) error {
	clientData := ClientData{}
	clientData.DstPort = uint16(80)
	segments := []string{"fcbb:bb00:1::1", "fcbb:bb00:2::1"}
	for i, segment := range segments {
		ipv6Segment := net.ParseIP(segment)
		copy(clientData.Segments[i].In6U.U6Addr8[:], ipv6Segment.To16())
	}

	// error := m.Put(domainToKey("wbhawknet"), clientData)
	key, err := FormatDNSName("wb.hawk.net")
	if err != nil {
		return err
	}
	error := m.Put(key, clientData)
	if error != nil {
		return error
	}

	return nil
}

func domainToKey(domain string) [256]byte {
	var key [256]byte
	copy(key[:], []byte(domain))
	return key
}

func FormatDNSName(domain string) ([256]byte, error) {
	var result [256]byte
	labels := strings.Split(domain, ".")
	offset := 0

	for _, label := range labels {
		if len(label) == 0 {
			return result, fmt.Errorf("Empty label detected")
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
