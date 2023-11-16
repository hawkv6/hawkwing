package bpf

import "github.com/cilium/ebpf"

type Bpf interface {
	CreateMap(spec *ebpf.MapSpec, pinDir string) (*ebpf.Map, error)
	LoadPinnedMap(pinPath string) (*ebpf.Map, error)
	LookupMap(m *ebpf.Map, key interface{}, value interface{}) error
	LoadMapFromId(id ebpf.MapID) (*ebpf.Map, error)
	PutMap(m *ebpf.Map, key interface{}, value interface{}) error
}
