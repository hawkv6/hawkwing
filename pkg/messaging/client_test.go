package messaging

import "testing"

func TestNewMessagingClient(t *testing.T) {
	mc := NewMessagingChannels()
	m := NewMessagingClient(mc)
	if m == nil {
		t.Errorf("NewMessagingClient() = %v, want %v", m, "not nil")
	}
	if m.messagingChannels != mc {
		t.Errorf("NewMessagingClient() = %v, want %v", m.messagingChannels, mc)
	}
}
