package syncer

import (
	"fmt"
	"slices"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/entities"
	"github.com/hawkv6/hawkwing/pkg/maps"
)

// storeSidList stores the sid list received from the adapter in the client map.
// The sid list is stored in the inner map corresponding to the domain name and
// intent name of the intent response.
//
// Parameters:
//   - intentResponse: The intent response containing the sid list to store.
//
// Returns:
//   - An error if the sid list could not be stored.
func (s *Syncer) storeSidList(intentResponse *entities.PathResult) error {
	var err error

	var innerMapId uint32
	err = s.cm.Reverse.Map.Lookup(maps.Ipv6ToInet6(intentResponse.Ipv6DestinationAddress), &innerMapId)
	if err != nil {
		return fmt.Errorf("could not find a value in the reverse map for %s: %s", intentResponse.Ipv6DestinationAddress, err)
	}

	// ports := s.getPortsToQuery(intentResponse.Ipv6DestinationAddress, intentResponse.Intent)
	portToUpdate := s.getApplicationPortToUpdate(intentResponse)
	sidListData := maps.GenerateSidLookupValue(intentResponse.Ipv6SidAddresses)

	err = s.cm.Outer.UpdateInner(innerMapId, uint16(portToUpdate), sidListData)
	if err != nil {
		return fmt.Errorf("could not update inner map: %s", err)
	}
	return nil
	// for _, port := range ports {
	// 	err = s.cm.Outer.UpdateInner(innerMapId, uint16(port), sidListData)
	// 	if err != nil {
	// 		return fmt.Errorf("could not update inner map: %s", err)
	// 	}
	// }
	// return nil
}

func (s *Syncer) getApplicationPortToUpdate(intentResult *entities.PathResult) int {
	configIntents := s.getApplicationConfigIntents(intentResult)
	resultIntents := s.getApplicationResultIntents(intentResult)

	for _, services := range config.Params.Services {
		if slices.Contains(services.Ipv6Addresses, intentResult.Ipv6DestinationAddress) {
			for _, application := range services.Applications {
				for _, ri := range resultIntents {
					if slices.Contains(configIntents, ri) {
						return application.Port
					}
				}
			}
		}
	}

	return 0
}

func (s *Syncer) getApplicationResultIntents(intentResult *entities.PathResult) []string {
	var intents []string
	for _, ir := range intentResult.Intents {
		intents = append(intents, ir.IntentType.String())
	}
	return intents
}

func (s *Syncer) getApplicationConfigIntents(intentResult *entities.PathResult) []string {
	var intents []string
	for _, service := range config.Params.Services {
		if slices.Contains(service.Ipv6Addresses, intentResult.Ipv6DestinationAddress) {
			for _, application := range service.Applications {
				for _, intent := range application.Intents {
					intents = append(intents, intent.Intent)
				}
			}
		}
	}
	return intents
}

// getPortsToQuery returns a list of ports to query for the given ipv6 address
// and intent type.
//
// Parameters:
//   - ipv6Addr: The ipv6 address to query.
//   - intentType: The intent type to query.
//
// Returns:
//   - A list of ports to query.
// func (s *Syncer) getPortsToQuery(ipv6Addr string, intentType entities.IntentType) []int {
// 	var ports []int
// 	for _, service := range config.Params.Services {
// 		for _, ip := range service.Ipv6Addresses {
// 			if ip == ipv6Addr {
// 				for _, intent := range service.Intents {
// 					if intentType.String() == intent.Intent {
// 						ports = append(ports, intent.Port)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return ports
// }
