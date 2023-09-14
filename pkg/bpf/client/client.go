package client

import (
	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS xdp ../../../bpf/bpf_dns.c -- -I../../../bpf
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS tc ../../../bpf/bpf_seg6.c -- -I../../../bpf

// ReadClientXdpObjects reads the XDP objects from the BPF filesystem.
func ReadClientXdpObjects() (*xdpObjects, error) {
	obj := &xdpObjects{}
	ops := &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: bpf.BpffsRoot,
		},
	}
	err := loadXdpObjects(obj, ops)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// ReadClientTcObjects reads the TC objects from the BPF filesystem.
func ReadClientTcObjects() (*tcObjects, error) {
	obj := &tcObjects{}
	ops := &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: bpf.BpffsRoot,
		},
	}
	err := loadTcObjects(obj, ops)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
