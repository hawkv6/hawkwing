package maps

import (
	"fmt"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/server"
)

type ServerMap struct {
	bpf    bpf.Bpf
	Lookup *Map
}

func NewServerMap(bpf bpf.Bpf) (*ServerMap, error) {
	collSpec, err := server.ReadServerBpfSpecs()
	if err != nil {
		return nil, fmt.Errorf("could not load server BPF specs: %s", err)
	}

	lookupMapSpec := collSpec.Maps["server_lookup_map"]
	return &ServerMap{
		bpf:    bpf,
		Lookup: NewMap(bpf, lookupMapSpec),
	}, nil
}

func (sm *ServerMap) Create() error {
	if err := sm.Lookup.OpenOrCreate(); err != nil {
		return fmt.Errorf("could not create server lookup map: %s", err)
	}
	return nil
}
