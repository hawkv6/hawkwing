package messaging

import (
	"fmt"
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
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "GetIntentPath returns no stream",
			testFunc: func(t *testing.T) {
				mc := NewMessagingChannels()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockClient := api.NewMockIntentControllerClient(ctrl)
				messagingClient := &MessagingClient{
					IntentControllerClient: mockClient,
					messagingChannels:      mc,
					streamErrors:           make(chan error),
					ErrCh:                  make(chan error),
				}

				mockClient.EXPECT().GetIntentPath(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

				go messagingClient.manageStreams()

				select {
				case <-messagingClient.ErrCh:
				case <-time.After(time.Second * 15):
					t.Errorf("manageStreams() = want %v, got %v", "message in messagingClient.ErrCh", "no message in messagingClient.ErrCh")
				}
			},
		},
		{
			name: "send returns error",
			testFunc: func(t *testing.T) {
				mc := NewMessagingChannels()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockClient := api.NewMockIntentControllerClient(ctrl)
				messagingClient := &MessagingClient{
					IntentControllerClient: mockClient,
					messagingChannels:      mc,
					streamErrors:           make(chan error),
					ErrCh:                  make(chan error),
				}
				mockStream := api.NewMockIntentController_GetIntentPathClient(ctrl)
				mockGrpcStreamCh := make(chan *api.PathRequest)

				mockClient.EXPECT().GetIntentPath(gomock.Any()).Return(mockStream, nil).AnyTimes()

				mockStream.EXPECT().
					Send(gomock.Any()).
					DoAndReturn(func(msg *api.PathRequest) error {
						mockGrpcStreamCh <- msg
						return fmt.Errorf("error")
					}).
					AnyTimes()

				mockStream.EXPECT().
					Recv().
					DoAndReturn(func() (*api.PathResult, error) {
						select {
						case <-mockGrpcStreamCh:
							return &api.PathResult{}, nil
						case <-time.After(time.Second * 5):
							return nil, io.EOF
						}
					}).
					AnyTimes()

				go messagingClient.manageStreams()

				// Give goroutines time to start
				time.Sleep(100 * time.Millisecond)

				mc.ChMessageIntentRequest <- &api.PathRequest{}
				select {
				case <-messagingClient.ErrCh:
				case <-time.After(time.Second * 15):
					t.Errorf("manageStreams() = want %v, got %v", "message in messagingClient.ErrCh", "no message in messagingClient.ErrCh")
				}
			},
		},
		{
			name: "receive returns error",
			testFunc: func(t *testing.T) {
				mc := NewMessagingChannels()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockClient := api.NewMockIntentControllerClient(ctrl)
				messagingClient := &MessagingClient{
					IntentControllerClient: mockClient,
					messagingChannels:      mc,
					streamErrors:           make(chan error),
					ErrCh:                  make(chan error),
				}
				mockStream := api.NewMockIntentController_GetIntentPathClient(ctrl)
				mockGrpcStreamCh := make(chan *api.PathRequest)

				mockClient.EXPECT().GetIntentPath(gomock.Any()).Return(mockStream, nil).AnyTimes()

				mockStream.EXPECT().
					Send(gomock.Any()).
					DoAndReturn(func(msg *api.PathRequest) error {
						mockGrpcStreamCh <- msg
						return nil
					}).
					AnyTimes()

				mockStream.EXPECT().
					Recv().
					DoAndReturn(func() (*api.PathResult, error) {
						select {
						case <-mockGrpcStreamCh:
							return &api.PathResult{}, fmt.Errorf("error")
						case <-time.After(time.Second * 5):
							return nil, io.EOF
						}
					}).
					AnyTimes()

				go messagingClient.manageStreams()

				// Give goroutines time to start
				time.Sleep(100 * time.Millisecond)

				mc.ChMessageIntentRequest <- &api.PathRequest{}
			},
		},
		{
			name: "no error",
			testFunc: func(t *testing.T) {
				mc := NewMessagingChannels()
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				mockClient := api.NewMockIntentControllerClient(ctrl)
				messagingClient := &MessagingClient{
					IntentControllerClient: mockClient,
					messagingChannels:      mc,
					streamErrors:           make(chan error),
					ErrCh:                  make(chan error),
				}
				mockStream := api.NewMockIntentController_GetIntentPathClient(ctrl)
				mockGrpcStreamCh := make(chan *api.PathRequest, 1)

				mockClient.EXPECT().GetIntentPath(gomock.Any()).Return(mockStream, nil).AnyTimes()

				sendReceived := make(chan struct{})
				mockStream.EXPECT().
					Send(gomock.Any()).
					DoAndReturn(func(msg *api.PathRequest) error {
						mockGrpcStreamCh <- msg
						close(sendReceived)
						return nil
					}).
					Times(1)

				mockStream.EXPECT().
					Recv().
					DoAndReturn(func() (*api.PathResult, error) {
						// Wait for Send to complete first to avoid race
						<-sendReceived
						select {
						case <-mockGrpcStreamCh:
							return &api.PathResult{}, nil
						case <-time.After(time.Second * 5):
							return nil, io.EOF
						}
					}).
					AnyTimes()

				go messagingClient.manageStreams()

				// Give goroutines time to start
				time.Sleep(100 * time.Millisecond)

				mc.ChMessageIntentRequest <- &api.PathRequest{}

				select {
				case <-mc.ChMessageIntentResponse:
				case <-time.After(time.Second * 15):
					t.Errorf("manageStreams() = want %v, got %v", "message in mc.ChMessageIntentResponse", "no message in mc.ChMessageIntentResponse")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}
