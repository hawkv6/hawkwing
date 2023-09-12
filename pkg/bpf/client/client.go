package client

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS xdp ../../../bpf/bpf_dns.c -- -I../../../bpf
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS tc ../../../bpf/bpf_seg6.c -- -I../../../bpf

type ClientData struct {
	DstPort  uint16
	_        [2]byte  // padding have to be added manually
	Dstaddr  [16]byte // empty for initializing
	Segments [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}

func ReadClientXdpObjects() (*xdpObjects, error) {
	obj := &xdpObjects{}
	ops := &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: bpf.BpffsRoot,
		},
	}
	err := loadXdpObjects(obj, ops)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func ReadClientTcObjects() (*tcObjects, error) {
	obj := &tcObjects{}
	ops := &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: bpf.BpffsRoot,
		},
	}
	err := loadTcObjects(obj, ops)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

type InnerMapData struct {
	DstPort uint16
	// _        [2]byte // padding have to be added manually
	Segments [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}

// func CreateMaps(outerMap *ebpf.Map, innerMapSpec *ebpf.MapSpec) (*ebpf.Map, *ebpf.Map, error) {
// 	outer := bpf.NewOuterMap(outerMap, innerMapSpec)
// 	innerMap, err := outer.CreateAndInsertInnerMap("wb.hawk.net")
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	innerMap2, err := outer.CreateAndInsertInnerMap("wc.hawk.net")
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return innerMap.InnerMap, innerMap2.InnerMap, nil
// }

func CreateInnerMaps() (*ebpf.Map, *ebpf.Map, error) {
	spec, err := ebpf.LoadCollectionSpec("./pkg/bpf/client/xdp_bpfel.o")
	if err != nil {
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}
	innerMap1Specs := spec.Maps["client_inner_map"]
	innerMap1Specs.Pinning = ebpf.PinByName
	innerMap1Specs.Name = "client_inner_map1"
	innerMap1, err := ebpf.NewMapWithOptions(innerMap1Specs, ebpf.MapOptions{
		PinPath: bpf.BpffsRoot,
	})
	if err != nil {
		return nil, nil, err
	}
	innerMap2Specs := spec.Maps["client_inner_map"]
	innerMap2Specs.Pinning = ebpf.PinByName
	innerMap2Specs.Name = "client_inner_map2"
	innerMap2, err := ebpf.NewMapWithOptions(innerMap2Specs, ebpf.MapOptions{
		PinPath: bpf.BpffsRoot,
	})
	if err != nil {
		return nil, nil, err
	}
	err = innerMap1.Pin("/sys/fs/bpf/client_inner_map1")
	if err != nil {
		log.Fatalf("Could not pin map: %s", err)
		return nil, nil, err
	}
	err = innerMap2.Pin("/sys/fs/bpf/client_inner_map2")
	if err != nil {
		log.Fatalf("Could not pin map: %s", err)
		return nil, nil, err
	}
	err = InitializeInnerMap(innerMap1, innerMap2)
	if err != nil {
		return nil, nil, err
	}
	return innerMap1, innerMap2, nil
}

func InitializeInnerMap(wb *ebpf.Map, wc *ebpf.Map) error {
	innerMapData := InnerMapData{}
	innerMapData.DstPort = uint16(80)
	segments := []string{"fcbb:bb00:1::1", "fcbb:bb00:2::1"}
	for i, segment := range segments {
		ipv6Segment := net.ParseIP(segment)
		copy(innerMapData.Segments[i].In6U.U6Addr8[:], ipv6Segment.To16())
	}
	err := wb.Put(innerMapData.DstPort, innerMapData.Segments)
	if err != nil {
		return err
	}

	innerMapData2 := InnerMapData{}
	innerMapData2.DstPort = uint16(8080)
	segments2 := []string{"fcbb:bb00:1::1", "fcbb:bb00:2::1"}
	for i, segment := range segments2 {
		ipv6Segment := net.ParseIP(segment)
		copy(innerMapData2.Segments[i].In6U.U6Addr8[:], ipv6Segment.To16())
	}
	err = wc.Put(innerMapData2.DstPort, innerMapData2.Segments)
	if err != nil {
		return err
	}
	return nil
}

func InitializeOuterMap(outer *ebpf.Map, innerWB *ebpf.Map, innerWC *ebpf.Map) error {
	key, err := FormatDNSName("wb.hawk.net")
	if err != nil {
		return err
	}
	err = outer.Put(key, uint32(innerWB.FD()))
	if err != nil {
		return err
	}
	key2, err := FormatDNSName("wc.hawk.net")
	if err != nil {
		return err
	}
	err = outer.Put(key2, uint32(innerWC.FD()))
	if err != nil {
		return err
	}
	return nil
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

// func domainToKey(domain string) [256]byte {
// 	var key [256]byte
// 	copy(key[:], []byte(domain))
// 	return key
// }

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

func Hope() *ebpf.Map {
	innerMapData := InnerMapData{}
	innerMapData.DstPort = uint16(80)
	segments := []string{"fcbb:bb00:1::1", "fcbb:bb00:2::1"}
	for i, segment := range segments {
		ipv6Segment := net.ParseIP(segment)
		copy(innerMapData.Segments[i].In6U.U6Addr8[:], ipv6Segment.To16())
	}

	innerMapData2 := InnerMapData{}
	innerMapData2.DstPort = uint16(8080)
	segments2 := []string{"fcbb:bb00:1::1", "fcbb:bb00:2::1"}
	for i, segment := range segments2 {
		ipv6Segment := net.ParseIP(segment)
		copy(innerMapData2.Segments[i].In6U.U6Addr8[:], ipv6Segment.To16())
	}

	innerMapDataList := []InnerMapData{innerMapData, innerMapData2}
	// innerMapKeyList := []string{"wb.hawk.net", "wc.hawk.net"}

	outerMapSpec := ebpf.MapSpec{
		Name:       "client_outer_map",
		Type:       ebpf.HashOfMaps,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 1024,
		Pinning:    ebpf.PinByName,
		Contents:   make([]ebpf.MapKV, 2),
		InnerMap: &ebpf.MapSpec{
			Name:      "client_inner_map",
			Type:      ebpf.LRUHash,
			KeySize:   2,   // sizeof(uint16)
			ValueSize: 160, // sizeof(struct in6_addr) * MAX_SEGMENTLIST_ENTRIES
			// Flags:      0x1000,
			MaxEntries: 1,
		},
	}

	for i := uint32(0); i < 2; i++ {
		innerMapSpec := outerMapSpec.InnerMap.Copy()
		innerMapSpec.Name = fmt.Sprintf("%s_%d", outerMapSpec.InnerMap.Name, i)
		innerMapSpec.MaxEntries = 1025 // TODO change this
		innerMapSpec.Contents = make([]ebpf.MapKV, 1)

		for j := range innerMapSpec.Contents {
			innerMapSpec.Contents[uint32(j)] = ebpf.MapKV{
				Key:   innerMapDataList[i].DstPort,
				Value: innerMapDataList[i].Segments,
			}
		}

		innerMap, err := ebpf.NewMap(innerMapSpec)
		if err != nil {
			log.Fatalf("Could not create inner-map: %s", err)
		}

		// key, err := FormatDNSName(innerMapKeyList[i])
		if err != nil {
			log.Fatalf("Could not format DNS name: %s", err)
		}
		outerMapSpec.Contents[i] = ebpf.MapKV{
			Key:   i,
			Value: innerMap,
		}
	}
	outerMap, err := ebpf.NewMapWithOptions(&outerMapSpec, ebpf.MapOptions{
		PinPath: bpf.BpffsRoot,
	})
	if err != nil {
		log.Fatalf("Could not create outer-map: %s", err)
	}
	return outerMap

}
