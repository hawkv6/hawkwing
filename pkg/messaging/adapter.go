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

func (a *MessagingAdapter) Start() error {
	return nil
}
