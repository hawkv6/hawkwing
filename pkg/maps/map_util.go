package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
)

// BuildClientDataMap constructs and initializes eBPF maps based on the provided specifications.
// It creates a lookup map and an outer map, populating them with the contents generated from
// the inner map specifications. The function also formats the DNS names and sets them as keys
// for the lookup map. If the maps are successfully created, they are built using the provided
// builder instances.
//
// Parameters:
//   - lookupSpec: Specification for the lookup eBPF map.
//   - outerSpec: Specification for the outer eBPF map.
//   - innerSpecs: A map where the key is a string (usually DNS name) and the value is the
//     specification for the inner eBPF map.
//
// Returns:
//   - error: An error is returned if any step of the map creation or building fails.
func BuildClientDataMap(lookupSpec, outerSpec *ebpf.MapSpec, innerSpecs map[string]*ebpf.MapSpec) error {
	lookupBuilder := NewEbpfMapBuilder(lookupSpec, pinnedMapOptions)
	outerBuilder := NewEbpfMapBuilder(outerSpec, pinnedMapOptions)

	lookupContents := []ebpf.MapKV{}
	outerContents := []ebpf.MapKV{}
	i := 0
	for key, innerSpec := range innerSpecs {
		innerMap, err := ebpf.NewMap(innerSpec)
		if err != nil {
			return fmt.Errorf("could not create inner map: %s", err)
		}

		formatedDnsName, err := FormatDNSName(key)
		if err != nil {
			return fmt.Errorf("could not format DNS name: %s", err)
		}

		outerContents = append(outerContents, ebpf.MapKV{
			Key:   uint32(i),
			Value: innerMap,
		})

		lookupContents = append(lookupContents, ebpf.MapKV{
			Key:   formatedDnsName,
			Value: uint32(i),
		})

		i++
	}

	// Setting contents before building
	lookupBuilder.SetContents(lookupContents)
	outerBuilder.SetContents(outerContents)

	if err := lookupBuilder.Build(); err != nil {
		return fmt.Errorf("could not build lookup map: %s", err)
	}

	if err := outerBuilder.Build(); err != nil {
		return fmt.Errorf("could not build outer map: %s", err)
	}

	return nil
}
