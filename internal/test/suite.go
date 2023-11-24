package test

import (
	"log"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/config"
)

var (
	MockClientInnerMapSpec = ebpf.MapSpec{
		Name:       "client_inner_map",
		Type:       ebpf.Hash,
		KeySize:    2,
		ValueSize:  164,
		MaxEntries: 1,
	}
	MockClientOuterMapSpec = ebpf.MapSpec{
		Name:       "client_outer_map",
		Type:       ebpf.HashOfMaps,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 100,
		Pinning:    ebpf.PinByName,
		Contents:   make([]ebpf.MapKV, 1),
		InnerMap:   &MockClientInnerMapSpec,
	}
	MockClientLookupMapSpec = ebpf.MapSpec{
		Name:       "client_lookup_map",
		Type:       ebpf.Hash,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 100,
	}
	MockClientReverseMapSpec = ebpf.MapSpec{
		Name:       "client_reverse_map",
		Type:       ebpf.Hash,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 100,
	}
	MockClientCollectionSpec = ebpf.CollectionSpec{
		Maps: map[string]*ebpf.MapSpec{
			"client_outer_map":   &MockClientOuterMapSpec,
			"client_lookup_map":  &MockClientLookupMapSpec,
			"client_reverse_map": &MockClientReverseMapSpec,
		},
	}
)

const testConfigPath = "../../test_assets/test_config.yaml"

func SetupTestConfig(tb testing.TB) {
	config.GetInstance().SetConfigFile(testConfigPath)
	if err := config.Parse(); err != nil {
		log.Fatalln(err)
	}
}
