package messaging

import (
	"io"
	"testing"
	"time"

	"github.com/hawkv6/hawkwing/pkg/api"
	"go.uber.org/mock/gomock"
)

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

func TestMessagingClient_manageStreams(t *testing.T) {
	mc := NewMessagingChannels()
	ctrl := gomock.NewController(t)

	mockClient := api.NewMockIntentControllerClient(ctrl)
	messagingClient := &MessagingClient{
		IntentControllerClient: mockClient,
		messagingChannels:      mc,
	}

	mockStream := api.NewMockIntentController_GetIntentPathClient(ctrl)

	mockClient.EXPECT().GetIntentPath(gomock.Any()).Return(mockStream, nil).AnyTimes()

	mockStream.EXPECT().
		Send(gomock.Any()).
		DoAndReturn(func(msg *api.PathRequest) error {
			mc.ChMessageIntentRequest <- msg
			return nil
		}).
		AnyTimes()

	mockStream.EXPECT().
		Recv().
		DoAndReturn(func() (*api.PathResult, error) {
			select {
			case <-mc.ChMessageIntentRequest:
				return &api.PathResult{}, nil
			case <-time.After(time.Second * 5):
				return nil, io.EOF
			}
		}).
		AnyTimes()

	go messagingClient.manageStreams()
	mc.ChMessageIntentRequest <- &api.PathRequest{}

	select {
	case <-mc.ChMessageIntentResponse:
	case <-time.After(time.Second * 2):
		t.Errorf("manageStreams() = want %v, got %v", "message in mc.ChMessageIntentResponse", "no message in mc.ChMessageIntentResponse")
	}
}
