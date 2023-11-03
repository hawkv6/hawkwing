package maps

import (
	"fmt"

	"github.com/hawkv6/hawkwing/pkg/bpf/server"
)

type ServerMap struct {
	Lookup *Map
}

func NewServerMap() (*ServerMap, error) {
	collSpec, err := server.ReadServerBpfSpecs()
	if err != nil {
		return nil, fmt.Errorf("could not load server BPF specs: %s", err)
	}

	lookupMapSpec := collSpec.Maps["server_lookup_map"]
	return &ServerMap{
		Lookup: NewMap(lookupMapSpec),
	}, nil
}

func (sm *ServerMap) Create() error {
	if err := sm.Lookup.OpenOrCreate(); err != nil {
		return fmt.Errorf("could not create server lookup map: %s", err)
	}
	return nil
}
