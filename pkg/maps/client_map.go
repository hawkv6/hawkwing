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

func NewClientMap(clientElfSpec *ebpf.CollectionSpec) *ClientMap {
	outerMapSpec := clientElfSpec.Maps["client_outer_map"]
	return &ClientMap{
		outerMapSpec: outerMapSpec,
	}
}

func (cm *ClientMap) clientInnerMapSpecs() map[string]*ebpf.MapSpec {
	innerMapSpecs := make(map[string]*ebpf.MapSpec)
	for key, services := range config.Params.Services {
		innerMapSpec := cm.outerMapSpec.InnerMap.Copy()
		innerMapSpec.Name = fmt.Sprintf("%s_%s", cm.outerMapSpec.InnerMap.Name, key)
		innerMapSpec.Contents = make([]ebpf.MapKV, len(services))
		for i, service := range services {
			value := GenerateSidLookupValue(service.Sid)
			innerMapSpec.Contents[uint32(i)] = ebpf.MapKV{
				Key:   uint16(service.Port),
				Value: value,
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
