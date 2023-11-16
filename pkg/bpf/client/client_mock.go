package client

import "github.com/cilium/ebpf"

type MockClientBpfReader struct{}

func (r *MockClientBpfReader) ReadClientBpfObjects() (*clientObjects, error) {
	return nil, nil
}

func (r *MockClientBpfReader) ReadClientBpfSpecs() (*ebpf.CollectionSpec, error) {
	clientOuterMapspec := &ebpf.MapSpec{
		Name:       "client_outer_map",
		Type:       ebpf.Hash,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 1,
	}
	clientLookupMapspec := &ebpf.MapSpec{
		Name:       "client_lookup_map",
		Type:       ebpf.Hash,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 1,
	}
	clientReverseMapspec := &ebpf.MapSpec{
		Name:       "client_reverse_map",
		Type:       ebpf.Hash,
		KeySize:    4,
		ValueSize:  4,
		MaxEntries: 1,
	}
	return &ebpf.CollectionSpec{
		Maps: map[string]*ebpf.MapSpec{
			"client_outer_map":   clientOuterMapspec,
			"client_lookup_map":  clientLookupMapspec,
			"client_reverse_map": clientReverseMapspec,
		},
		Programs: map[string]*ebpf.ProgramSpec{},
	}, nil
}
