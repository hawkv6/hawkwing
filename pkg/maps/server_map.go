package maps

import (
	"github.com/cilium/ebpf"
)

type ServerMap struct {
	lookupMap *ebpf.MapSpec
}

func (sm *ServerMap) CreateServerDataMaps() error {
	lookupBuilder := NewEbpfMapBuilder(sm.lookupMap, pinnedMapOptions)
	if err := lookupBuilder.Build(); err != nil {
		return err
	}
	return nil
}

func NewServerMap() *ServerMap {
	lookupMapSpec := &ebpf.MapSpec{
		Name:       "server_lookup_map",
		Type:       ebpf.LRUHash,
		KeySize:    18,  // 16 bytes for IPv6 address + 2 bytes for port
		ValueSize:  160, // 10 * 16 bytes for IPv6 address
		MaxEntries: 1024,
		Pinning:    ebpf.PinByName,
		Contents:   make([]ebpf.MapKV, 1),
	}
	return &ServerMap{
		lookupMap: lookupMapSpec,
	}
}
