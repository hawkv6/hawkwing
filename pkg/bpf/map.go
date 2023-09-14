package bpf

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/config"
)

type InnerMapData struct {
	DstPort  uint16
	Segments [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}

type ClientMap struct {
	outerMapSpec *ebpf.MapSpec
}

func (cm *ClientMap) CreateClientDataMaps() error {
	lookupMapSpec := cm.clientLookupMapSpec()
	innerMapSpecs := cm.clientInnerMapSpecs()
	i := 0
	for key, innerMapSpec := range innerMapSpecs {
		innerMap, err := ebpf.NewMap(innerMapSpec)
		if err != nil {
			return fmt.Errorf("could not create inner map: %s", err)
		}

		formatedDnsName, err := FormatDNSName(key)
		if err != nil {
			return fmt.Errorf("could not format DNS name: %s", err)
		}

		cm.outerMapSpec.Contents[i] = ebpf.MapKV{
			Key:   uint32(i),
			Value: innerMap,
		}

		lookupMapSpec.Contents[i] = ebpf.MapKV{
			Key:   formatedDnsName,
			Value: uint32(i),
		}

		i++
	}

	_, err := ebpf.NewMapWithOptions(lookupMapSpec, ebpf.MapOptions{
		PinPath: BpffsRoot,
	})
	if err != nil {
		return fmt.Errorf("could not create lookup map: %s", err)
	}

	_, err = ebpf.NewMapWithOptions(cm.outerMapSpec, ebpf.MapOptions{
		PinPath: BpffsRoot,
	})
	if err != nil {
		return fmt.Errorf("could not create outer map: %s", err)
	}

	return nil
}

func NewClientMap() *ClientMap {
	outerMapSpec := &ebpf.MapSpec{
		Name:       "client_outer_map",
		Type:       ebpf.HashOfMaps,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 1024,
		Pinning:    ebpf.PinByName,
		Contents:   make([]ebpf.MapKV, len(config.Params.Services)),
		InnerMap: &ebpf.MapSpec{
			Name:       "client_inner_map",
			Type:       ebpf.LRUHash,
			KeySize:    2,
			ValueSize:  160,
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
			innerMapSpec.Contents[uint32(i)] = ebpf.MapKV{
				Key:   uint16(service.Port),
				Value: SidToInet6Sid(service.Sid),
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
