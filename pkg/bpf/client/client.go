package client

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS client ../../../bpf/bpf_client.c -- -I../../../bpf

func ReadClientBpfObjects() (*clientObjects, error) {
	obj := &clientObjects{}
	ops := &ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			PinPath: bpf.BpffsRoot,
		},
	}
	err := loadClientObjects(obj, ops)
	if err != nil {
		return nil, fmt.Errorf("could not load client BPF objects: %s", err)
	}
	return obj, nil
}
