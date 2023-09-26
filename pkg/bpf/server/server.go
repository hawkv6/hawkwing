package server

import (
	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS xdp ../../../bpf/bpf_server_ingress.c -- -I../../../bpf

// // ReadServerTcObjects reads the TC objects from the BPF filesystem.
// func ReadServerTcObjects() (*tcObjects, error) {
// 	obj := &tcObjects{}
// 	ops := &ebpf.CollectionOptions{
// 		Maps: ebpf.MapOptions{
// 			PinPath: bpf.BpffsRoot,
// 		},
// 	}
// 	err := loadTcObjects(obj, ops)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

func ReadServerXdpObjects() (*xdpObjects, error) {
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
