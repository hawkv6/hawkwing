package syncer

import (
	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/entities"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/messaging"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const Subsystem = "go-syncer"

type Syncer struct {
	adapterChannels *messaging.AdapterChannels
	cm              *maps.ClientMap
	reqChan         chan *entities.PathRequest
	resolver        *ResolverService
}

func NewSyncer(adapterChannels *messaging.AdapterChannels, cm *maps.ClientMap) *Syncer {
	reqChan := make(chan *entities.PathRequest)
	return &Syncer{
		adapterChannels: adapterChannels,
		cm:              cm,
		reqChan:         reqChan,
		resolver:        NewResolverService(reqChan),
	}
}

func (s *Syncer) Start() {
	log.Printf("start synchronization process")
	s.handleIntentMessages()
}

func (s *Syncer) FetchAll() {
	log.Printf("fetching all needed intent details")
	for key, _ := range config.Params.Services {
		pathRequests := entities.CreatePathRequestsForService(key)
		for _, pathRequest := range pathRequests {
			s.reqChan <- &pathRequest
		}

		// for _, service := range services {
		// 	if service.Sid == nil {
		// 		s.reqChan <- &messaging.PathRequest{
		// 			DomainName: key,
		// 			IntentName: service.Intent,
		// 		}
		// 	}
		// }
	}
}

func (s *Syncer) handleIntentMessages() {
	go func() {
		for {
			intentRequest := <-s.reqChan
			s.adapterChannels.ChAdapterIntentRequest <- intentRequest
			log.Printf("requested intent details for %s", intentRequest.Ipv6DestinationAddress)
		}
	}()
	go func() {
		for {
			intentResponse := <-s.adapterChannels.ChAdapterIntentResponse
			if err := s.storeSidList(intentResponse); err != nil {
				log.WithError(err).Error("could not store sid list")
			}
			log.Printf("stored intent details for %s", intentResponse.Ipv6DestinationAddress)
		}
	}()
}
