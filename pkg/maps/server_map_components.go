package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
)

type ServerLookupMap struct {
	spec *ebpf.MapSpec
}

func NewServerLookupMap(spec *ebpf.MapSpec) *ServerLookupMap {
	return &ServerLookupMap{spec: spec}
}

func (slm *ServerLookupMap) Build() error {
	_, err := ebpf.NewMapWithOptions(slm.spec, pinnedMapOptions)
	if err != nil {
		return fmt.Errorf("could not create server lookup map: %s", err)
	}
	return nil
}
