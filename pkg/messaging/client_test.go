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
	mc := NewMessagingChannels()
	ctrl := gomock.NewController(t)
	mockClient := api.NewMockIntentControllerClient(ctrl)
	messagingClient := &MessagingClient{
		IntentControllerClient: mockClient,
		messagingChannels:      mc,
	}
	mockStream := api.NewMockIntentController_GetIntentPathClient(ctrl)
	mockGrpcStreamCh := make(chan *api.PathRequest)

	tests := []struct {
		name     string
		testFunc func()
	}{
		// {
		// 	name: "GetIntentPath returns no stream",
		// 	testFunc: func() {
		// 		mockClient.EXPECT().GetIntentPath(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
		// 		messagingClient.manageStreams()
		// 	},
		// },
		{
			name: "send returns error",
			testFunc: func() {
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
				mc.ChMessageIntentRequest <- &api.PathRequest{}
			},
		},
		{
			name: "receive returns error",
			testFunc: func() {
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
				mc.ChMessageIntentRequest <- &api.PathRequest{}
			},
		},
		{
			name: "no error",
			testFunc: func() {
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc()
			ctrl.Finish()
		})
	}
}
