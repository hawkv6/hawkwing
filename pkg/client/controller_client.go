package client

import (
	"fmt"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/messaging"
	"github.com/hawkv6/hawkwing/pkg/syncer"
)

type ControllerClient struct {
	mainErrCh         chan error
	adapterChannels   *messaging.AdapterChannels
	messagingChannels *messaging.MessagingChannels
	messenger         *messaging.MessagingClient
	adapter           *messaging.MessagingAdapter
	syncer            *syncer.Syncer
}

func NewControllerClient(mainErrCh chan error, clientMap *maps.ClientMap) *ControllerClient {
	messagingChannels := messaging.NewMessagingChannels()
	adapterChannels := messaging.NewAdapterChannels()
	messenger := messaging.NewMessagingClient(messagingChannels)
	adapter := messaging.NewMessagingAdapter(messagingChannels, adapterChannels)

	realBpf := &bpf.RealBpf{}
	syncer := syncer.NewSyncer(realBpf, adapterChannels, clientMap)

	return &ControllerClient{
		mainErrCh:         mainErrCh,
		adapterChannels:   adapterChannels,
		messagingChannels: messagingChannels,
		messenger:         messenger,
		adapter:           adapter,
		syncer:            syncer,
	}
}

func (cc *ControllerClient) Start() {
	go func() {
		cc.messenger.Start()
		cc.adapter.Start()
		cc.syncer.Start()
		cc.syncer.FetchAll()
		select {
		case err := <-cc.messenger.ErrCh:
			cc.mainErrCh <- fmt.Errorf("messenger error: %s", err)
		case err := <-cc.syncer.ErrCh:
			cc.mainErrCh <- fmt.Errorf("syncer error: %s", err)
		}
	}()
}
