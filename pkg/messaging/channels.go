package messaging

import "github.com/hawkv6/hawkwing/pkg/api"

type MessagingChannels struct {
	ChMessageIntentRequest  chan *api.Intent
	ChMessageIntentResponse chan *api.Response
}

func NewMessagingChannels() *MessagingChannels {
	return &MessagingChannels{
		ChMessageIntentRequest:  make(chan *api.Intent),
		ChMessageIntentResponse: make(chan *api.Response),
	}
}

type AdapterChannels struct {
	ChAdapterIntentRequest  chan *IntentRequest
	ChAdapterIntentResponse chan *IntentResponse
}

func NewAdapterChannels() *AdapterChannels {
	return &AdapterChannels{
		ChAdapterIntentRequest:  make(chan *IntentRequest),
		ChAdapterIntentResponse: make(chan *IntentResponse),
	}
}
