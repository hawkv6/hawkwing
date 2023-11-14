package messaging

import (
	"testing"
)

func TestNewMessagingChannels(t *testing.T) {
	mc := NewMessagingChannels()
	if mc == nil {
		t.Errorf("NewMessagingChannels() = %v, want %v", mc, "not nil")
	}
	if mc.ChMessageIntentRequest == nil {
		t.Errorf("NewMessagingChannels() = %v, want %v", mc.ChMessageIntentRequest, "not nil")
	}
	if mc.ChMessageIntentResponse == nil {
		t.Errorf("NewMessagingChannels() = %v, want %v", mc.ChMessageIntentResponse, "not nil")
	}
}

func TestNewAdapterChannels(t *testing.T) {
	ac := NewAdapterChannels()
	if ac == nil {
		t.Errorf("NewAdapterChannels() = %v, want %v", ac, "not nil")
	}
	if ac.ChAdapterIntentRequest == nil {
		t.Errorf("NewAdapterChannels() = %v, want %v", ac.ChAdapterIntentRequest, "not nil")
	}
	if ac.ChAdapterIntentResponse == nil {
		t.Errorf("NewAdapterChannels() = %v, want %v", ac.ChAdapterIntentResponse, "not nil")
	}
}
