package messaging

import (
	"context"
	"strconv"
	"time"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const Subsystem = "go-messaging"

type MessagingClient struct {
	api.IntentControllerClient
	messagingChannels *MessagingChannels
	conn              *grpc.ClientConn
	streamErrors      chan error
}

func NewMessagingClient(messagingChannels *MessagingChannels) *MessagingClient {
	return &MessagingClient{
		messagingChannels: messagingChannels,
		streamErrors:      make(chan error),
	}
}

func (c *MessagingClient) Start() {
	c.connect()
	go c.manageStreams()

}

func (c *MessagingClient) connect() {
	connectionAddress := "[" + config.Params.HawkEye.Hostname + "]" + ":" + strconv.Itoa(config.Params.HawkEye.Port)
	for {
		conn, err := grpc.Dial(connectionAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("failed to dial: %v, retrying...", err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.conn = conn
		c.IntentControllerClient = api.NewIntentControllerClient(conn)
		break
	}
	log.Printf("connected to %s", connectionAddress)
}

func (c *MessagingClient) manageStreams() {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := c.IntentControllerClient.GetIntentPath(ctx)
		if err != nil {
			log.Printf("failed to get intent path: %v, retrying...", err)
			c.streamErrors <- err
			continue
		}

		go c.handleGetIntentPathRequests(ctx, stream)
		go c.handleGetIntentPathResults(ctx, stream)

		err = <-c.streamErrors
		cancel()

		log.Fatalf("stream error: %v", err)
		// log.WithError(err).Errorf("stream error")
	}
}

func (c *MessagingClient) handleGetIntentPathRequests(ctx context.Context, stream api.IntentController_GetIntentPathClient) {
	for {
		select {
		case request := <-c.messagingChannels.ChMessageIntentRequest:
			if err := stream.Send(request); err != nil {
				log.Printf("failed to send message: %v", err)
				c.streamErrors <- err
				return
			}
			log.Printf("sent intent request for %s", request.Ipv6DestinationAddress)
		case <-ctx.Done():
			return
		}
	}
}

func (c *MessagingClient) handleGetIntentPathResults(ctx context.Context, stream api.IntentController_GetIntentPathClient) {
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("failed to receive message: %v", err)
			c.streamErrors <- err
			return
		}
		log.Printf("received intent result for %s", in.Ipv6DestinationAddress)
		select {
		case c.messagingChannels.ChMessageIntentResponse <- in:
		case <-ctx.Done():
			return
		}
	}
}
