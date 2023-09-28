package client

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/vishvananda/netlink"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const (
	Subsystem = "go-client"
)

type Client struct {
	iface     netlink.Link
	xdpLinker *linker.XdpLinker
	tcLinker  *linker.TcLinker
	wg        *sync.WaitGroup
}

func NewClient(interfaceName string) (*Client, error) {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup network iface %q: %s", interfaceName, err)
	}
	xdpObjs, err := client.ReadClientXdpObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load XDP program: %s", err)
	}
	xdpLinker := linker.NewXdpLinker(iface, xdpObjs.ClientIngress)
	tcObjs, err := client.ReadClientTcObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load TC program: %s", err)
	}
	tcLinker := linker.NewTcLinker(iface, tcObjs.ClientEgress, "egress")

	// TODO change this
	err = bpf.Mount()
	if err != nil {
		log.Fatalf("Could not mount BPF filesystem: %s", err)
	}

	clientMap := maps.NewClientMap()
	err = clientMap.CreateClientDataMaps()
	if err != nil {
		log.Fatalf("Could not create client data maps: %s", err)
	}

	return &Client{
		iface:     iface,
		xdpLinker: xdpLinker,
		tcLinker:  tcLinker,
		wg:        &sync.WaitGroup{},
	}, nil
}

func (c *Client) Start() {
	c.wg.Add(2)

	go func() {
		defer c.wg.Done()
		if err := c.xdpLinker.Attach(); err != nil {
			log.WithError(err).Error("couldn't attach XDP program")
		}
	}()

	go func() {
		defer c.wg.Done()
		if err := c.tcLinker.Attach(); err != nil {
			log.WithError(err).Error("couldn't attach TC program")
		}
	}()

	log.Printf("Attached XDP program to iface %q", c.iface.Attrs().Name)
	log.Printf("Press Ctrl-C to exit and remove the program")
}

func (c *Client) Stop() {
	if err := c.xdpLinker.Detach(); err != nil {
		log.WithError(err).Error("couldn't detach XDP program")
	}
	if err := c.tcLinker.Detach(); err != nil {
		log.WithError(err).Error("couldn't detach TC program")
	}
	c.wg.Wait()
}
