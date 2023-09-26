// Code generated by bpf2go; DO NOT EDIT.
//go:build arm64be || armbe || mips || mips64 || mips64p32 || ppc64 || s390 || s390x || sparc || sparc64

package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type xdpServerLookupKey struct {
	Addr struct{ In6U struct{ U6Addr8 [16]uint8 } }
	Port uint16
}

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
	FilterIngress *ebpf.ProgramSpec `ebpf:"filter_ingress"`
}

// xdpMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type xdpMapSpecs struct {
	ServerLookupMap *ebpf.MapSpec `ebpf:"server_lookup_map"`
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
	ServerLookupMap *ebpf.Map `ebpf:"server_lookup_map"`
}

func (m *xdpMaps) Close() error {
	return _XdpClose(
		m.ServerLookupMap,
	)
}

// xdpPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadXdpObjects or ebpf.CollectionSpec.LoadAndAssign.
type xdpPrograms struct {
	FilterIngress *ebpf.Program `ebpf:"filter_ingress"`
}

func (p *xdpPrograms) Close() error {
	return _XdpClose(
		p.FilterIngress,
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
//go:embed xdp_bpfeb.o
var _XdpBytes []byte
