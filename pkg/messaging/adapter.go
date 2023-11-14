package messaging

import "github.com/hawkv6/hawkwing/pkg/entities"

type MessagingAdapter struct {
	messagingChannels *MessagingChannels
	adapterChannels   *AdapterChannels
}

func NewMessagingAdapter(messagingChannels *MessagingChannels, adapterChannels *AdapterChannels) *MessagingAdapter {
	return &MessagingAdapter{
		messagingChannels: messagingChannels,
		adapterChannels:   adapterChannels,
	}
}

func (a *MessagingAdapter) Start() {
	go func() {
		for {
			intentRequest := <-a.adapterChannels.ChAdapterIntentRequest
			a.messagingChannels.ChMessageIntentRequest <- intentRequest.Marshal()
		}
	}()
	go func() {
		for {
			intentResponse := <-a.messagingChannels.ChMessageIntentResponse
			pathResult := entities.UnmarshalPathResult(intentResponse)
			a.adapterChannels.ChAdapterIntentResponse <- pathResult
		}
	}()
}
