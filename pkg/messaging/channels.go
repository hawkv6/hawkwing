package messaging

import "github.com/hawkv6/hawkwing/pkg/api"

type MessagingChannels struct {
	ChMessageIntentRequest  chan *api.Intent
	ChMessageIntentResponse chan *api.Response
}

type AdapterChannels struct {
	ChAdapterIntentRequest  chan *IntentRequest
	ChAdapterIntentResponse chan *IntentResponse
}
