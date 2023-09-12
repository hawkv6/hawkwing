package client

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/hawkv6/hawkwing/pkg/logging"
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

const (
	clientReverseMapPath = "/sys/fs/bpf/client_reverse_map"
	clientInnerMapPath   = "/sys/fs/bpf/client_inner_map"
	clientOuterMapPath   = "/sys/fs/bpf/client_outer_map"
	clientMapPath        = "/sys/fs/bpf/client_map"
)

func NewClient(interfaceName string) (*Client, error) {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup network iface %q: %s", interfaceName, err)
	}
	xdpObjs, err := client.ReadClientXdpObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load XDP program: %s", err)
	}
	xdpLinker := linker.NewXdpLinker(iface, xdpObjs.InterceptDns)
	tcObjs, err := client.ReadClientTcObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load TC program: %s", err)
	}
	tcLinker := linker.NewTcLinker(iface, tcObjs.EncapEgress)

	// TODO change this
	err = bpf.Mount()
	if err != nil {
		log.Fatalf("Could not mount BPF filesystem: %s", err)
	}

	// err = client.InitializeBpfMap(xdpObjs.ClientMap)
	// if err != nil {
	// 	log.Fatalf("Could not initialize BPF map: %s", err)
	// }
	// err = client.InitializeInnerMap(xdpObjs.ClientInnerMap, xdpObjs.ClientInnerMap)
	// if err != nil {
	// 	log.Fatalf("Could not initialize BPF map: %s", err)
	// }
	// innerMap1, innerMap2, err := client.CreateInnerMaps()
	// if err != nil {
	// 	log.Fatalf("Could not create inner maps: %s", err)
	// }
	// err = client.InitializeOuterMap(xdpObjs.ClientOuterMap, innerMap1, innerMap2)
	// if err != nil {
	// 	log.Fatalf("Could not initialize BPF map: %s", err)
	// }
	_ = client.Hope()

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
