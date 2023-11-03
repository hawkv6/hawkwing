package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
)

type InnerMap struct {
	Map
	ID int
}

func NewInnerMap(spec *ebpf.MapSpec) *InnerMap {
	return &InnerMap{
		Map: Map{
			spec: spec,
		},
	}
}

func (im *InnerMap) Build() error {
	mapInstance, err := ebpf.NewMap(im.spec)
	if err != nil {
		return fmt.Errorf("could not create inner map: %s", err)
	}
	im.m = mapInstance
	return nil
}

type OuterMap struct {
	Map
}

func NewOuterMap(spec *ebpf.MapSpec) *OuterMap {
	return &OuterMap{Map: Map{spec: spec}}
}

func (om *OuterMap) BuildWith(inners map[string]*InnerMap) error {
	outerContents := []ebpf.MapKV{}
	for _, inner := range inners {
		innerMap, err := ebpf.NewMap(inner.spec)
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

func NewLookupMap(spec *ebpf.MapSpec) *LookupMap {
	return &LookupMap{
		Map: Map{spec: spec},
	}
}

func (lm *LookupMap) BuildWith(inners map[string]*InnerMap) error {
	lookupContents := []ebpf.MapKV{}
	for key, inner := range inners {
		formattedDnsName, err := FormatDNSName(key)
		if err != nil {
			return fmt.Errorf("could not format DNS name: %s", err)
		}
		lookupContents = append(lookupContents, ebpf.MapKV{
			Key:   formattedDnsName,
			Value: uint32(inner.ID),
		})
	}

	lm.Map.spec.Contents = lookupContents
	err := lm.Map.OpenOrCreate()
	if err != nil {
		return fmt.Errorf("could not create lookup map: %s", err)
	}
	return nil
}
