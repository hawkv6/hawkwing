package client

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/vishvananda/netlink"
)

type EbpfClient struct {
	iface     netlink.Link
	xdpLinker *linker.XdpLinker
	tcLinker  *linker.TcLinker
	wg        *sync.WaitGroup
	mainErrCh chan error
}

func NewEbpfClient(interfaceName string, mainErrCh chan error) (*EbpfClient, *maps.ClientMap, error) {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return nil, nil, fmt.Errorf("could not lookup network iface %q: %s", interfaceName, err)
	}

	realClientBpfReader := &client.RealClientBpfReader{}
	clientObjs, err := realClientBpfReader.ReadClientBpfObjects()
	if err != nil {
		return nil, nil, fmt.Errorf("could not load client BPF objects: %s", err)
	}

	xdpLinker := linker.NewXdpLinker(iface, clientObjs.ClientIngress)
	tcLinker := linker.NewTcLinker(iface, clientObjs.ClientEgress, "egress")

	err = bpf.Mount()
	if err != nil {
		return nil, nil, fmt.Errorf("could not mount BPF filesystem: %s", err)
	}

	realBpf := &bpf.RealBpf{}
	clientMap, err := maps.NewClientMap(realBpf, realClientBpfReader)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create client map: %s", err)
	}
	err = clientMap.Create()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create client lookup map: %s", err)
	}

	return &EbpfClient{
		iface:     iface,
		xdpLinker: xdpLinker,
		tcLinker:  tcLinker,
		wg:        &sync.WaitGroup{},
		mainErrCh: mainErrCh,
	}, clientMap, nil
}

func (ec *EbpfClient) Start() {
	ec.wg.Add(2)

	go func() {
		defer ec.wg.Done()
		if err := ec.xdpLinker.Attach(); err != nil {
			ec.mainErrCh <- fmt.Errorf("could not attach client XDP program: %s", err)
		}
	}()

	go func() {
		defer ec.wg.Done()
		if err := ec.tcLinker.Attach(); err != nil {
			ec.mainErrCh <- fmt.Errorf("could not attach client TC program: %s", err)
		}
	}()

	fmt.Printf("Client started on interface %q\n", ec.iface.Attrs().Name)
	fmt.Println("Press Ctrl-C to exit and remove the program")
}

func (ec *EbpfClient) Stop() {
	if err := ec.xdpLinker.Detach(); err != nil {
		ec.mainErrCh <- fmt.Errorf("could not detach client XDP program: %s", err)
	}
	if err := ec.tcLinker.Detach(); err != nil {
		ec.mainErrCh <- fmt.Errorf("could not detach client TC program: %s", err)
	}
	ec.wg.Wait()
}
