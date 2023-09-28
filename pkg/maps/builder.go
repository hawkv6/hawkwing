package maps

import (
	"fmt"

	"github.com/cilium/ebpf"
)

type MapBuilder interface {
	Build() error
}

type EbpfMapBuilder struct {
	spec     *ebpf.MapSpec
	opts     ebpf.MapOptions
	contents []ebpf.MapKV
}

func NewEbpfMapBuilder(spec *ebpf.MapSpec, opts ebpf.MapOptions) *EbpfMapBuilder {
	return &EbpfMapBuilder{
		spec: spec,
		opts: opts,
	}
}

func (e *EbpfMapBuilder) SetContents(contents []ebpf.MapKV) {
	e.contents = contents
}

func (e *EbpfMapBuilder) Build() error {
	e.spec.Contents = e.contents
	_, err := ebpf.NewMapWithOptions(e.spec, e.opts)
	if err != nil {
		return fmt.Errorf("could not create map: %s, %s", e.spec.Name, err)
	}
	return nil
}
