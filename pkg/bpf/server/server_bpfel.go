// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64

package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type serverServerLookupKey struct {
	Addr struct{ In6U struct{ U6Addr8 [16]uint8 } }
	Port uint16
}

type serverSidlistData struct {
	SidlistSize uint32
	Sidlist     [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}

// loadServer returns the embedded CollectionSpec for server.
func loadServer() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_ServerBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load server: %w", err)
	}

	return spec, err
}

// loadServerObjects loads server and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*serverObjects
//	*serverPrograms
//	*serverMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadServerObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadServer()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// serverSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type serverSpecs struct {
	serverProgramSpecs
	serverMapSpecs
}

// serverSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type serverProgramSpecs struct {
	ServerEgress  *ebpf.ProgramSpec `ebpf:"server_egress"`
	ServerIngress *ebpf.ProgramSpec `ebpf:"server_ingress"`
}

// serverMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type serverMapSpecs struct {
	PercpuSidlistMap   *ebpf.MapSpec `ebpf:"percpu_sidlist_map"`
	ServerLookupMap    *ebpf.MapSpec `ebpf:"server_lookup_map"`
	ServerTempSidMap   *ebpf.MapSpec `ebpf:"server_temp_sid_map"`
	ServerTempValueMap *ebpf.MapSpec `ebpf:"server_temp_value_map"`
}

// serverObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadServerObjects or ebpf.CollectionSpec.LoadAndAssign.
type serverObjects struct {
	serverPrograms
	serverMaps
}

func (o *serverObjects) Close() error {
	return _ServerClose(
		&o.serverPrograms,
		&o.serverMaps,
	)
}

// serverMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadServerObjects or ebpf.CollectionSpec.LoadAndAssign.
type serverMaps struct {
	PercpuSidlistMap   *ebpf.Map `ebpf:"percpu_sidlist_map"`
	ServerLookupMap    *ebpf.Map `ebpf:"server_lookup_map"`
	ServerTempSidMap   *ebpf.Map `ebpf:"server_temp_sid_map"`
	ServerTempValueMap *ebpf.Map `ebpf:"server_temp_value_map"`
}

func (m *serverMaps) Close() error {
	return _ServerClose(
		m.PercpuSidlistMap,
		m.ServerLookupMap,
		m.ServerTempSidMap,
		m.ServerTempValueMap,
	)
}

// serverPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadServerObjects or ebpf.CollectionSpec.LoadAndAssign.
type serverPrograms struct {
	ServerEgress  *ebpf.Program `ebpf:"server_egress"`
	ServerIngress *ebpf.Program `ebpf:"server_ingress"`
}

func (p *serverPrograms) Close() error {
	return _ServerClose(
		p.ServerEgress,
		p.ServerIngress,
	)
}

func _ServerClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed server_bpfel.o
var _ServerBytes []byte
