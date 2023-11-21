package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

type InnerMap struct {
	Map
	ID int
}

func NewInnerMap(bpf bpf.Bpf, spec *ebpf.MapSpec) *InnerMap {
	return &InnerMap{
		Map: Map{
			bpf:  bpf,
			spec: spec,
		},
	}
}

func (im *InnerMap) Build() error {
	mapInstance, err := im.bpf.CreateMap(im.spec, "")
	if err != nil {
		return fmt.Errorf("could not create inner map: %s", err)
	}
	im.m = mapInstance
	return nil
}

type OuterMap struct {
	Map
}

func NewOuterMap(bpf bpf.Bpf, spec *ebpf.MapSpec) *OuterMap {
	return &OuterMap{
		Map: Map{
			bpf:  bpf,
			spec: spec,
		},
	}
}

func (om *OuterMap) BuildWith(inners map[string]*InnerMap) error {
	outerContents := []ebpf.MapKV{}
	for _, inner := range inners {
		// innerMap, err := ebpf.NewMap(inner.spec)
		innerMap, err := inner.Map.bpf.CreateMap(inner.spec, "")
		if err != nil {
			return fmt.Errorf("could not create inner map: %s", err)
		}
		outerContents = append(outerContents, ebpf.MapKV{
			Key:   uint32(inner.ID),
			Value: innerMap,
		})
	}
	om.spec.Contents = outerContents
	err := om.Map.OpenOrCreate()
	if err != nil {
		return fmt.Errorf("could not create outer map: %s", err)
	}
	return nil
}

type LookupMap struct {
	Map
}

func NewLookupMap(bpf bpf.Bpf, spec *ebpf.MapSpec) *LookupMap {
	return &LookupMap{
		Map: Map{
			bpf:  bpf,
			spec: spec,
		},
	}
}

func (lm *LookupMap) BuildWith(inners map[string]*InnerMap) error {
	lookupContents := []ebpf.MapKV{}
	for key, inner := range inners {
		if config.Params.Services[key].DomainName != "" {
			formattedDnsName, err := FormatDNSName(config.Params.Services[key].DomainName)
			if err != nil {
				return fmt.Errorf("could not format DNS name: %s", err)
			}
			lookupContents = append(lookupContents, ebpf.MapKV{
				Key:   formattedDnsName,
				Value: uint32(inner.ID),
			})
		}
	}

	lm.Map.spec.Contents = lookupContents
	err := lm.Map.OpenOrCreate()
	if err != nil {
		return fmt.Errorf("could not create lookup map: %s", err)
	}
	return nil
}

type ReverseMap struct {
	Map
}

func NewReverseMap(bpf bpf.Bpf, spec *ebpf.MapSpec) *ReverseMap {
	return &ReverseMap{
		Map: Map{
			bpf:  bpf,
			spec: spec,
		},
	}
}

func (rm *ReverseMap) BuildWith(inners map[string]*InnerMap) error {
	reverseMapContents := []ebpf.MapKV{}
	for key, inner := range inners {
		if config.Params.Services[key].DomainName == "" && len(config.Params.Services[key].Ipv6Addresses) > 0 {
			for _, ipv6Address := range config.Params.Services[key].Ipv6Addresses {
				reverseMapContents = append(reverseMapContents, ebpf.MapKV{
					Key:   Ipv6ToInet6(ipv6Address),
					Value: uint32(inner.ID),
				})
			}
		}
	}

	rm.Map.spec.Contents = reverseMapContents
	err := rm.Map.OpenOrCreate()
	if err != nil {
		return fmt.Errorf("could not create reverse map: %s", err)
	}
	return nil
}
