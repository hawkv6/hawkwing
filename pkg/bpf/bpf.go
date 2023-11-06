package bpf

import (
	"fmt"
	"os"

	"github.com/cilium/ebpf"
)

func createMap(spec *ebpf.MapSpec, opts *ebpf.MapOptions) (*ebpf.Map, error) {
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
func CreateMap(spec *ebpf.MapSpec, pinDir string) (*ebpf.Map, error) {
	var opts ebpf.MapOptions

	if spec.Pinning != 0 {
		if pinDir == "" {
			return nil, fmt.Errorf("map requires pinning, but no pinDir specified")
		}
		if spec.Name == "" {
			return nil, fmt.Errorf("map requires pinning, but no name specified")
		}
		if err := os.MkdirAll(pinDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create pin directory: %w", err)
		}
		opts.PinPath = pinDir
	}

	m, err := createMap(spec, &opts)

	return m, err
}
