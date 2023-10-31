package server

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS server ../../../bpf/bpf_server.c -- -I../../../bpf

func ReadServerBpfObjects() (*serverObjects, error) {
	obj := &serverObjects{}
	ops := &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: bpf.BpffsRoot,
		},
	}
	err := loadServerObjects(obj, ops)
	if err != nil {
		return nil, fmt.Errorf("could not load server BPF objects: %s", err)
	}
	return obj, nil
}

func ReadServerBpfSpecs() (*ebpf.CollectionSpec, error) {
	specs, err := loadServer()
	if err != nil {
		return nil, fmt.Errorf("could not load server BPF collection: %s", err)
	}
	return specs, nil
}
