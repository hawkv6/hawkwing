package messaging

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
	a.HandleIntent()
}

func (a *MessagingAdapter) HandleIntent() {
	go func() {
		for {
			intentRequest := <-a.adapterChannels.ChAdapterIntentRequest
			a.messagingChannels.ChMessageIntentRequest <- intentRequest.Marshal()
		}
	}()
	go func() {
		for {
			intentResponse := <-a.messagingChannels.ChMessageIntentResponse
			a.adapterChannels.ChAdapterIntentResponse <- &IntentResponse{
				SidList: intentResponse.Ipv6Addresses,
			}
		}
	}()
}
