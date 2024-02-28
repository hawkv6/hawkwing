package client

import (
	"fmt"
	"net"
	"strings"

	"github.com/hawkv6/hawkwing/internal/config"
)

type Client struct {
	mainErrCh        chan error
	ebpfClient       *EbpfClient
	controllerClient *ControllerClient
}

func NewClient(mainErrCh chan error, interfaceName string) (*Client, error) {
	ief, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, err
	}

	if config.Params.ClientIpv6Address == "" {
		addrs, err := ief.Addrs()
		if err != nil {
			return nil, err
		}
		addr := strings.Split(addrs[0].String(), "/")
		config.Params.ClientIpv6Address = addr[0]
	}

	ebpfClient, clientMap, err := NewEbpfClient(interfaceName, mainErrCh)
	if err != nil {
		return nil, err
	}

	if config.Params.HawkEye.Enabled {
		controllerClient := NewControllerClient(mainErrCh, clientMap)
		return &Client{
			mainErrCh:        mainErrCh,
			ebpfClient:       ebpfClient,
			controllerClient: controllerClient,
		}, nil
	}

	return &Client{
		mainErrCh:  mainErrCh,
		ebpfClient: ebpfClient,
	}, nil
}

func (c *Client) Start() {
	if !config.Params.HawkEye.Enabled {
		fmt.Println("==================INFO==================")
		fmt.Println("Hawkwing is running in standalone mode.")
		fmt.Println("========================================")
		c.ebpfClient.Start()
	} else {
		fmt.Println("==================INFO==================")
		fmt.Println("Hawkwing is running in controller mode.")
		fmt.Println("========================================")
		c.ebpfClient.Start()
		c.controllerClient.Start()
	}
}

func (c *Client) Stop() {
	c.ebpfClient.Stop()
}
