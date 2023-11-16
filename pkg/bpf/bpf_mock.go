package bpf

import (
	"fmt"

	"github.com/cilium/ebpf"
)

type MockBpf struct{}

func (m *MockBpf) CreateMap(spec *ebpf.MapSpec, pinDir string) (*ebpf.Map, error) {
	var opts ebpf.MapOptions

	if spec.Pinning != 0 {
		if pinDir == "" {
			return nil, fmt.Errorf("map requires pinning, but no pinDir specified")
		}
		opts.PinPath = pinDir
	}

	mockMap := &ebpf.Map{}
	return mockMap, nil
}

func (m *MockBpf) LoadPinnedMap(pinPath string) (*ebpf.Map, error) {
	if pinPath == "" {
		return nil, fmt.Errorf("pinPath cannot be empty")
	}
	mockMap := &ebpf.Map{}
	return mockMap, nil
}

func (m *MockBpf) LookupMap(em *ebpf.Map, key interface{}, value interface{}) error {
	if em == nil {
		return fmt.Errorf("map cannot be nil")
	}
	if key == nil {
		return fmt.Errorf("key cannot be nil")
	}
	if value == nil {
		return fmt.Errorf("value cannot be nil")
	}
	return nil
}

func (m *MockBpf) LoadMapFromId(id ebpf.MapID) (*ebpf.Map, error) {
	if id == 0 {
		return nil, fmt.Errorf("id cannot be 0")
	}
	mockMap := &ebpf.Map{}
	return mockMap, nil
}

func (m *MockBpf) PutMap(em *ebpf.Map, key interface{}, value interface{}) error {
	if em == nil {
		return fmt.Errorf("map cannot be nil")
	}
	if key == nil {
		return fmt.Errorf("key cannot be nil")
	}
	if value == nil {
		return fmt.Errorf("value cannot be nil")
	}
	return nil
}
