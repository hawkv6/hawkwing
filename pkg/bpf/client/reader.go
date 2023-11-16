package client

import "github.com/cilium/ebpf"

type ClientBpfReader interface {
	ReadClientBpfObjects() (*clientObjects, error)
	ReadClientBpfSpecs() (*ebpf.CollectionSpec, error)
}
