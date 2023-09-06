package client

import (
	"log"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/vishvananda/netlink"
)

func NewClient(interfaceName string) {
	objs, err := bpf.ReadBpfObjects(nil)
	if err != nil {
		log.Fatalf("Could not load XDP program: %s", err)
	}
	defer objs.Close()

	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		log.Fatalf("Could not lookup network iface %q: %s", interfaceName, err)
	}
	xdpLinker := linker.NewXdpLinker(link, objs.InterceptDns)
	err = xdpLinker.Attach()
	if err != nil {
		log.Fatalf("Could not attach XDP program: %s", err)
	}

	if err != nil {
		log.Fatalf("Could not attach XDP program: %s", err)
	}
	defer xdpLinker.Detach()

	err = bpf.InitializeBpfMap(objs.ClientMap)
	if err != nil {
		log.Fatalf("Could not initialize BPF map: %s", err)
	}

	log.Printf("Attached XDP program to iface %q", link.Attrs().Name)
	log.Printf("Press Ctrl-C to exit and remove the program")
	select {}
}

func NewTcClient(interfaceName string) {
	objs, err := bpf.ReadTcObjects(nil)
	if err != nil {
		log.Fatalf("Could not load TC program: %s", err)
	}
	defer objs.Close()
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		log.Fatalf("Could not lookup network iface %q: %s", interfaceName, err)
	}

	tcLinker := linker.NewTcLinker(link, objs.EncapEgress)
	err = tcLinker.Attach()
	if err != nil {
		log.Fatalf("Could not attach TC program: %s", err)
	}
	defer tcLinker.Detach()
	select {}
}
