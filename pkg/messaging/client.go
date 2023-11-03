package messaging

import (
	"context"
	"time"

	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const Subsystem = "go-messaging"

type MessagingClient struct {
	api.IntentServiceClient
	messagingChannels *MessagingChannels
	conn              *grpc.ClientConn
}

func NewMessagingClient(messagingChannels *MessagingChannels) *MessagingClient {
	return &MessagingClient{
		messagingChannels: messagingChannels,
	}
}

func (c *MessagingClient) Start() {
	c.connect()
	c.handleIntentRequest()

}

func (c *MessagingClient) connect() {
	// connectionAddress := config.Params.HawkEye.Hostname + ":" + strconv.Itoa(config.Params.HawkEye.Port)
	connectionAddress := "[fcbb:cc00:5::f]:5001"
	for {
		conn, err := grpc.Dial(connectionAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("failed to dial: %v, retrying...", err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.conn = conn
		c.IntentServiceClient = api.NewIntentServiceClient(conn)
		break
	}
	log.Printf("connected to %s", connectionAddress)
}

func (c *MessagingClient) handleIntentRequest() {
	go func() {
		for {
			intentRequest := <-c.messagingChannels.ChMessageIntentRequest
			intentResponse, err := c.GetIntentDetails(context.Background(), intentRequest)
			if err != nil {
				log.Printf("failed to get intent details: %v, retrying...", err)
				c.connect()
				continue
			}
			c.messagingChannels.ChMessageIntentResponse <- intentResponse
			log.Printf("received intent response for [domain | intent]: [%s | %s]", intentResponse.DomainName, intentEnumToString(intentResponse.Intent))
		}
	}()
}
