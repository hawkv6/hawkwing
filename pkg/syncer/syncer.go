package syncer

import (
	"fmt"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/messaging"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const Subsystem = "go-syncer"

type Syncer struct {
	adapterChannels *messaging.AdapterChannels
	cm              *maps.ClientMap
	reqChan         chan *messaging.IntentRequest
}

func NewSyncer(adapterChannels *messaging.AdapterChannels, cm *maps.ClientMap) *Syncer {
	return &Syncer{
		adapterChannels: adapterChannels,
		cm:              cm,
		reqChan:         make(chan *messaging.IntentRequest),
	}
}

func (s *Syncer) Start() {
	log.Printf("start synchronization process")
	s.handleIntentMessages()
}

func (s *Syncer) FetchAll() {
	log.Printf("fetching all needed intent details")
	for key, services := range config.Params.Services {
		for _, service := range services {
			if service.Sid == nil {
				s.reqChan <- &messaging.IntentRequest{
					DomainName: key,
					IntentName: service.Intent,
				}
			}
		}
	}
}

func (s *Syncer) handleIntentMessages() {
	go func() {
		for {
			intentRequest := <-s.reqChan
			s.adapterChannels.ChAdapterIntentRequest <- intentRequest
			log.Printf("sent intent request for [domain | intent]: [%s | %s]", intentRequest.DomainName, intentRequest.IntentName)
		}
	}()
	go func() {
		for {
			intentResponse := <-s.adapterChannels.ChAdapterIntentResponse
			if err := s.storeSidList(intentResponse); err != nil {
				fmt.Printf("could not store sid list: %s", err)
			}
			log.Printf("stored sid list for [domain | intent]: [%s | %s]", intentResponse.DomainName, intentResponse.IntentName)
		}
	}()
}
