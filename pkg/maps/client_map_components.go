package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
)

type InnerMap struct {
	ID   int
	spec *ebpf.MapSpec
	m    *ebpf.Map
}

func NewInnerMap(spec *ebpf.MapSpec) *InnerMap {
	return &InnerMap{spec: spec}
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
	spec *ebpf.MapSpec
	m    *ebpf.Map
}

func NewOuterMap(spec *ebpf.MapSpec) *OuterMap {
	return &OuterMap{spec: spec}
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
	mapInstance, err := ebpf.NewMapWithOptions(om.spec, pinnedMapOptions)
	if err != nil {
		return fmt.Errorf("could not create outer map: %s", err)
	}
	om.m = mapInstance
	return nil
}

type LookupMap struct {
	spec *ebpf.MapSpec
	m    *ebpf.Map
}

func NewLookupMap(spec *ebpf.MapSpec) *LookupMap {
	return &LookupMap{spec: spec}
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
	lm.spec.Contents = lookupContents
	mapInstance, err := ebpf.NewMapWithOptions(lm.spec, pinnedMapOptions)
	if err != nil {
		return fmt.Errorf("could not create lookup map: %s", err)
	}
	lm.m = mapInstance
	return nil
}
