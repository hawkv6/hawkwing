package messaging

import (
	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/entities"
)

type MessagingChannels struct {
	ChMessageIntentRequest  chan *api.PathRequest
	ChMessageIntentResponse chan *api.PathResult
}

func NewMessagingChannels() *MessagingChannels {
	return &MessagingChannels{
		ChMessageIntentRequest:  make(chan *api.PathRequest),
		ChMessageIntentResponse: make(chan *api.PathResult),
	}
}

type AdapterChannels struct {
	ChAdapterIntentRequest  chan *entities.PathRequest
	ChAdapterIntentResponse chan *entities.PathResult
}

func NewAdapterChannels() *AdapterChannels {
	return &AdapterChannels{
		ChAdapterIntentRequest:  make(chan *entities.PathRequest),
		ChAdapterIntentResponse: make(chan *entities.PathResult),
	}
}
