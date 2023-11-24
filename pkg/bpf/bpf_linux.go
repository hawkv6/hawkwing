package bpf

import (
	"fmt"
	"os"

	"github.com/cilium/ebpf"
)

type RealBpf struct{}

func (r *RealBpf) createMap(spec *ebpf.MapSpec, opts *ebpf.MapOptions) (*ebpf.Map, error) {
	if opts == nil {
		opts = &ebpf.MapOptions{}
	}

	m, err := ebpf.NewMapWithOptions(spec, *opts)

	return m, err
}

// CreateMap creates a new eBPF map.
//
// Parameters:
//   - spec: The specification of the map to create.
//   - pinDir: The directory to pin the map to.
//
// Returns:
//   - The created map.
//   - An error if the map could not be created.
func (r *RealBpf) CreateMap(spec *ebpf.MapSpec, pinDir string) (*ebpf.Map, error) {
	var opts ebpf.MapOptions

	if spec.Pinning != 0 {
		if pinDir == "" {
			return nil, fmt.Errorf("map requires pinning, but no pinDir specified")
		}
		if err := os.MkdirAll(pinDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create pin directory: %w", err)
		}
		opts.PinPath = pinDir
	}

	m, err := r.createMap(spec, &opts)

	return m, err
}

// LoadPinnedMap loads a pinned eBPF map.
//
// Parameters:
//   - pinPath: The path to the pinned map.
//
// Returns:
//   - The loaded map.
//   - An error if the map could not be loaded.
func (r *RealBpf) LoadPinnedMap(pinPath string) (*ebpf.Map, error) {
	m, err := ebpf.LoadPinnedMap(pinPath, nil)

	return m, err
}

// LookupMap looks up a value in a map.
//
// Parameters:
//   - m: The map to lookup the value in.
//   - key: The key to lookup.
//   - value: The value to store the result in.
//
// Returns:
//   - An error if the value could not be looked up.
func (r *RealBpf) LookupMap(m *ebpf.Map, key interface{}, value interface{}) error {
	return m.Lookup(key, value)
}

// LoadMapFromId loads a map from its ID.
//
// Parameters:
//   - id: The ID of the map to load.
//
// Returns:
//   - The loaded map.
//   - An error if the map could not be loaded.
func (r *RealBpf) LoadMapFromId(id ebpf.MapID) (*ebpf.Map, error) {
	m, err := ebpf.NewMapFromID(id)
	return m, err
}

// PutMap puts a value in a map.
//
// Parameters:
//   - m: The map to put the value in.
//   - key: The key to put the value for.
//   - value: The value to put.
//
// Returns:
//   - An error if the value could not be put.
func (r *RealBpf) PutMap(m *ebpf.Map, key interface{}, value interface{}) error {
	return m.Put(key, value)
}
