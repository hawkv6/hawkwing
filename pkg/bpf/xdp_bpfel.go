// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64

package bpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadXdp returns the embedded CollectionSpec for xdp.
func loadXdp() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_XdpBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load xdp: %w", err)
	}

	return spec, err
}

// loadXdpObjects loads xdp and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*xdpObjects
//	*xdpPrograms
//	*xdpMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadXdpObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadXdp()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// xdpSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpSpecs struct {
	xdpProgramSpecs
	xdpMapSpecs
}

// xdpSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpProgramSpecs struct {
	ParseDns *ebpf.ProgramSpec `ebpf:"parse_dns"`
}

// xdpMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpMapSpecs struct {
}

// xdpObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadXdpObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpObjects struct {
	xdpPrograms
	xdpMaps
}

func (o *xdpObjects) Close() error {
	return _XdpClose(
		&o.xdpPrograms,
		&o.xdpMaps,
	)
}

// xdpMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadXdpObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpMaps struct {
}

func (m *xdpMaps) Close() error {
	return _XdpClose()
}

// xdpPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadXdpObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpPrograms struct {
	ParseDns *ebpf.Program `ebpf:"parse_dns"`
}

func (p *xdpPrograms) Close() error {
	return _XdpClose(
		p.ParseDns,
	)
}

func _XdpClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed xdp_bpfel.o
var _XdpBytes []byte
