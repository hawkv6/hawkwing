package client

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/messaging"
	"github.com/hawkv6/hawkwing/pkg/syncer"
	"github.com/vishvananda/netlink"
)

type Client struct {
	mainErrCh         chan error
	iface             netlink.Link
	xdpLinker         *linker.XdpLinker
	tcLinker          *linker.TcLinker
	wg                *sync.WaitGroup
	adapterChannels   *messaging.AdapterChannels
	messagingChannels *messaging.MessagingChannels
	messenger         *messaging.MessagingClient
	adapter           *messaging.MessagingAdapter
	syncer            *syncer.Syncer
}

func NewClient(mainErrCh chan error, interfaceName string) (*Client, error) {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup network iface %q: %s", interfaceName, err)
	}

	realClientBpfReader := &client.RealClientBpfReader{}
	clientObjs, err := realClientBpfReader.ReadClientBpfObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load client BPF objects: %s", err)
	}

	xdpLinker := linker.NewXdpLinker(iface, clientObjs.ClientIngress)
	tcLinker := linker.NewTcLinker(iface, clientObjs.ClientEgress, "egress")

	err = bpf.Mount()
	if err != nil {
		return nil, fmt.Errorf("could not mount BPF filesystem: %s", err)
	}

	realBpf := &bpf.RealBpf{}
	clientMap, err := maps.NewClientMap(realBpf, realClientBpfReader)
	if err != nil {
		return nil, fmt.Errorf("could not create client map: %s", err)
	}
	err = clientMap.Create()
	if err != nil {
		return nil, fmt.Errorf("could not create client lookup map: %s", err)
	}

	messagingChannels := messaging.NewMessagingChannels()
	adapterChannels := messaging.NewAdapterChannels()
	messenger := messaging.NewMessagingClient(messagingChannels)
	adapter := messaging.NewMessagingAdapter(messagingChannels, adapterChannels)
	syncer := syncer.NewSyncer(realBpf, adapterChannels, clientMap)

	return &Client{
		mainErrCh:         mainErrCh,
		iface:             iface,
		xdpLinker:         xdpLinker,
		tcLinker:          tcLinker,
		wg:                &sync.WaitGroup{},
		adapterChannels:   adapterChannels,
		messagingChannels: messagingChannels,
		messenger:         messenger,
		adapter:           adapter,
		syncer:            syncer,
	}, nil
}

func (c *Client) Start() {
	c.wg.Add(2)

	go func() {
		c.messenger.Start()
		c.adapter.Start()
		c.syncer.Start()
		c.syncer.FetchAll()
		select {
		case err := <-c.messenger.ErrCh:
			c.mainErrCh <- fmt.Errorf("messenger error: %s", err)
		case err := <-c.syncer.ErrCh:
			c.mainErrCh <- fmt.Errorf("syncer error: %s", err)
		}
	}()

	go func() {
		defer c.wg.Done()
		if err := c.xdpLinker.Attach(); err != nil {
			c.mainErrCh <- fmt.Errorf("could not attach client XDP program: %s", err)
		}
	}()

	go func() {
		defer c.wg.Done()
		if err := c.tcLinker.Attach(); err != nil {
			c.mainErrCh <- fmt.Errorf("could not attach client TC program: %s", err)
		}
	}()

	fmt.Printf("Client started on interface %q\n", c.iface.Attrs().Name)
	fmt.Println("Press Ctrl-C to exit and remove the program")
}

func (c *Client) Stop() {
	if err := c.xdpLinker.Detach(); err != nil {
		c.mainErrCh <- fmt.Errorf("could not detach client XDP program: %s", err)
	}
	if err := c.tcLinker.Detach(); err != nil {
		c.mainErrCh <- fmt.Errorf("could not detach client TC program: %s", err)
	}
	c.wg.Wait()
}
