package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/config"
)

type ClientMap struct {
	outerMapSpec *ebpf.MapSpec
}

func (cm *ClientMap) CreateClientDataMaps() error {
	lookupMapSpec := cm.clientLookupMapSpec()
	innerMapSpecs := cm.clientInnerMapSpecs()

	err := BuildClientDataMap(lookupMapSpec, cm.outerMapSpec, innerMapSpecs)
	if err != nil {
		return fmt.Errorf("could not build client data maps: %s", err)
	}

	return nil
}

func NewClientMap() *ClientMap {
	outerMapSpec := &ebpf.MapSpec{
		Name:       "client_outer_map",
		Type:       ebpf.HashOfMaps,
		KeySize:    4, // 2 bytes for domain id
		ValueSize:  4, // 2 bytes for inner map fd
		MaxEntries: 1024,
		Pinning:    ebpf.PinByName,
		Contents:   make([]ebpf.MapKV, len(config.Params.Services)),
		InnerMap: &ebpf.MapSpec{
			Name:       "client_inner_map",
			Type:       ebpf.LRUHash,
			KeySize:    2,
			ValueSize:  164, // TODO use this when storing struct 164
			MaxEntries: 1,
		},
	}
	return &ClientMap{
		outerMapSpec: outerMapSpec,
	}
}

func (cm *ClientMap) clientInnerMapSpecs() map[string]*ebpf.MapSpec {
	innerMapSpecs := make(map[string]*ebpf.MapSpec)
	for key, services := range config.Params.Services {
		innerMapSpec := cm.outerMapSpec.InnerMap.Copy()
		innerMapSpec.Name = fmt.Sprintf("%s_%s", cm.outerMapSpec.InnerMap.Name, key)
		innerMapSpec.MaxEntries = uint32(len(services))
		innerMapSpec.Contents = make([]ebpf.MapKV, len(services))
		for i, service := range services {
			// TODO - check or remove
			// value := GenerateSidLookupValueTest(service.Sid)
			// fmt.Println(unsafe.Sizeof(value))
			// b, err := value.Marshal()
			// fmt.Println(len(b))
			// if err != nil {
			// 	panic(err)
			// }
			// value := SidToInet6Sid(service.Sid)
			innerMapSpec.Contents[uint32(i)] = ebpf.MapKV{
				Key:   uint16(service.Port),
				Value: SidToInet6Sid(service.Sid),
				// Value: b,
			}
		}
		innerMapSpecs[key] = innerMapSpec
	}
	return innerMapSpecs
}

func (cm *ClientMap) clientLookupMapSpec() *ebpf.MapSpec {
	return &ebpf.MapSpec{
		Name:       "client_lookup_map",
		Type:       ebpf.LRUHash,
		KeySize:    256,
		ValueSize:  4,
		MaxEntries: 1024,
		Pinning:    ebpf.PinByName,
		Contents:   make([]ebpf.MapKV, len(config.Params.Services)),
	}
}
