package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/server"
)

type ServerMap struct {
	collSpec *ebpf.CollectionSpec
}

func NewServerMap() (*ServerMap, error) {
	collSpec, err := server.ReadServerBpfSpecs()
	if err != nil {
		return nil, fmt.Errorf("could not load server BPF specs: %s", err)
	}
	return &ServerMap{
		collSpec: collSpec,
	}, nil
}

func (sm *ServerMap) CreateServerLookupMap() error {
	lookupMap := sm.collSpec.Maps["server_lookup_map"]
	lookupBuilder := NewEbpfMapBuilder(lookupMap, pinnedMapOptions)
	if err := lookupBuilder.Build(); err != nil {
		return fmt.Errorf("could not build server lookup map: %s", err)
	}
	return nil
}
