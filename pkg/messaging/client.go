package messaging

import (
	"context"
	"log"
	"time"

	"github.com/hawkv6/hawkwing/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MessagingClient struct {
	api.IntentServiceClient
	messagingChannels *MessagingChannels
	conn              *grpc.ClientConn
}

func NewMessagingClient() *MessagingClient {
	return &MessagingClient{}
}

func (c *MessagingClient) Start() {
	c.connect()
	c.handleIntentRequest()

}

func (c *MessagingClient) connect() {
	connectionAddress := "localhost:50051"
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
		}
	}()
}
