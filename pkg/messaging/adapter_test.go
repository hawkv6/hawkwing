package messaging

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/entities"
)

func TestNewMessagingAdapter(t *testing.T) {
	mc := NewMessagingChannels()
	ac := NewAdapterChannels()
	ma := NewMessagingAdapter(mc, ac)
	if ma == nil {
		t.Errorf("NewMessagingAdapter() = %v, want %v", ma, "not nil")
	}
	if ma.messagingChannels != mc {
		t.Errorf("NewMessagingAdapter() = %v, want %v", ma.messagingChannels, mc)
	}
	if ma.adapterChannels != ac {
		t.Errorf("NewMessagingAdapter() = %v, want %v", ma.adapterChannels, ac)
	}
}

func TestMessagingAdapter_Start(t *testing.T) {
	test.SetupTestConfig(t)
	mc := NewMessagingChannels()
	ac := NewAdapterChannels()
	ma := NewMessagingAdapter(mc, ac)

	mockPathRequests, _ := entities.CreatePathRequestsForService("service3")
	fmt.Println(mockPathRequests)

	mockPathResult := api.PathResult{
		Ipv6DestinationAddress: "fcbb:cc00:4::a",
		Ipv6SidAddresses: []string{
			"2001:db8:0:1::1",
		},
	}

	ma.Start()

	ac.ChAdapterIntentRequest <- &mockPathRequests[0]
	select {
	case receivedRequest := <-mc.ChMessageIntentRequest:
		if reflect.DeepEqual(receivedRequest, mockPathRequests[0]) {
			t.Errorf("receivedRequest = %v, want %v", receivedRequest, mockPathRequests[0])
		}
	case <-time.After(time.Second * 1):
		t.Error("Timeout: did not receive message on ChMessageIntentRequest")
	}

	mc.ChMessageIntentResponse <- &mockPathResult
	select {
	case receivedResponse := <-ac.ChAdapterIntentResponse:
		if receivedResponse == nil {
			t.Error("receivedResponse = nil, want not nil")
		}
	case <-time.After(time.Second * 1):
		t.Error("Timeout: did not receive message on ChAdapterIntentResponse")
	}
}
