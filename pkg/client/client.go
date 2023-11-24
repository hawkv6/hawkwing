package client

type Client struct {
	mainErrCh        chan error
	ebpfClient       *EbpfClient
	controllerClient *ControllerClient
}

func NewClient(mainErrCh chan error, interfaceName string) (*Client, error) {
	ebpfClient, clientMap, err := NewEbpfClient(interfaceName, mainErrCh)
	if err != nil {
		return nil, err
	}
	controllerClient := NewControllerClient(mainErrCh, clientMap)
	return &Client{
		mainErrCh:        mainErrCh,
		ebpfClient:       ebpfClient,
		controllerClient: controllerClient,
	}, nil
}

func (c *Client) Start() {
	c.ebpfClient.Start()
	c.controllerClient.Start()
}

func (c *Client) Stop() {
	c.ebpfClient.Stop()
}
