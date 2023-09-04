package client

import (
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

func NewClient(interfaceName string) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		log.Fatalf("Could not lookup network iface %q: %s", interfaceName, err)
	}

	objs, err := bpf.ReadBpfObjects(nil)
	if err != nil {
		log.Fatalf("Could not load XDP program: %s", err)
	}
	defer objs.Close()

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.InterceptDns,
		Interface: iface.Index,
		Flags:     link.XDPGenericMode,
	})

	if err != nil {
		log.Fatalf("Could not attach XDP program: %s", err)
	}
	defer l.Close()

	log.Printf("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)
	log.Printf("Press Ctrl-C to exit and remove the program")

	err = bpf.InitializeBpfMap(objs.ClientMap)
	if err != nil {
		log.Fatalf("Could not initialize BPF map: %s", err)
	}
	select {}
}
