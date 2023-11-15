package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
)

type ClientMap struct {
	Outer   *OuterMap
	Inners  map[string]*InnerMap
	Lookup  *LookupMap
	Reverse *ReverseMap
}

func NewClientMap() (*ClientMap, error) {
	clientElfSpec, err := client.ReadClientBpfSpecs()
	if err != nil {
		return nil, fmt.Errorf("could not load client BPF specs: %s", err)
	}
	outerMapSpec := clientElfSpec.Maps["client_outer_map"]
	lookupMapSpec := clientElfSpec.Maps["client_lookup_map"]
	reverseMapSpec := clientElfSpec.Maps["client_reverse_map"]
	return &ClientMap{
		Outer:   NewOuterMap(outerMapSpec),
		Inners:  make(map[string]*InnerMap),
		Lookup:  NewLookupMap(lookupMapSpec),
		Reverse: NewReverseMap(reverseMapSpec),
	}, nil
}

func (cm *ClientMap) Create() error {
	innerSpecs, err := cm.clientInnerMapSpecs()
	if err != nil {
		return fmt.Errorf("could not create inner map specs: %s", err)
	}
	i := 0
	for key, spec := range innerSpecs {
		inner := NewInnerMap(spec)
		inner.ID = i
		cm.Inners[key] = inner
		i++
	}
	return cm.BuildClientDataMap()
}

func (cm *ClientMap) BuildClientDataMap() error {
	if err := cm.Lookup.BuildWith(cm.Inners); err != nil {
		return fmt.Errorf("could not build lookup map: %s", err)
	}
	if err := cm.Outer.BuildWith(cm.Inners); err != nil {
		return fmt.Errorf("could not build outer map: %s", err)
	}
	if err := cm.Reverse.BuildWith(cm.Inners); err != nil {
		return fmt.Errorf("could not build reverse map: %s", err)
	}
	return nil
}

func (cm *ClientMap) clientInnerMapSpecs() (map[string]*ebpf.MapSpec, error) {
	innerMapSpecs := make(map[string]*ebpf.MapSpec)
	for key, serviceCfg := range config.Params.Services {
		innerMapSpec := cm.Outer.spec.InnerMap.Copy()
		innerMapSpec.Name = fmt.Sprintf("%s_%s", "client_inner", key)
		innerMapSpec.Contents = make([]ebpf.MapKV, len(serviceCfg.Applications))
		for i, application := range serviceCfg.Applications {
			value, err := GenerateSidLookupValue(application.Sid)
			if err != nil {
				return nil, fmt.Errorf("could not generate SID lookup value: %s", err)
			}
			innerMapSpec.Contents[uint32(i)] = ebpf.MapKV{
				Key:   uint16(application.Port),
				Value: value,
			}
		}
		innerMapSpecs[key] = innerMapSpec
	}
	return innerMapSpecs, nil
}
