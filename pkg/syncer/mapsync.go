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

	portsToUpdate := s.getApplicationPortsToUpdate(intentResponse)
	sidListData, err := maps.GenerateSidLookupValue(intentResponse.Ipv6SidAddresses)
	if err != nil {
		return fmt.Errorf("could not generate sid lookup value: %s", err)
	}

	for _, portToUpdate := range portsToUpdate {
		err = s.cm.Outer.UpdateInner(innerMapId, uint16(portToUpdate), sidListData)
		if err != nil {
			return fmt.Errorf("could not update inner map: %s", err)
		}
		log.Infof("stored result for %s and port %d", intentResponse.Ipv6DestinationAddress, portToUpdate)
	}

	return nil
}

// getApplicationPortsToUpdate returns the ports of the application that needs to
// be updated based on the intent result.
//
// Parameters:
//   - intentResult: The intent result containing the intents that were
//     satisfied.
//
// Returns:
//   - The ports of the application that needs to be updated.
func (s *Syncer) getApplicationPortsToUpdate(intentResult *entities.PathResult) []int {
	configIntents := s.getApplicationConfigIntents(intentResult)
	resultIntents := s.getApplicationResultIntents(intentResult)
	var ports []int
	for _, services := range config.Params.Services {
		if slices.Contains(services.Ipv6Addresses, intentResult.Ipv6DestinationAddress) {
			for _, application := range services.Applications {
				for _, ri := range resultIntents {
					if slices.Contains(configIntents, ri) {
						ports = append(ports, application.Port)
					}
				}
			}
		}
	}

	return ports
}

// getApplicationResultIntents returns the intents that were satisfied in the
// intent result.
//
// Parameters:
//   - intentResult: The intent result containing the intents that were
//     satisfied.
//
// Returns:
//   - The intents that were satisfied in the intent result.
func (s *Syncer) getApplicationResultIntents(intentResult *entities.PathResult) []string {
	var intents []string
	for _, ir := range intentResult.Intents {
		intents = append(intents, ir.IntentType.String())
	}
	return intents
}

// getApplicationConfigIntents returns the intents that were configured for the
// application.
//
// Parameters:
//   - intentResult: The intent result containing the intents that were
//     satisfied.
//
// Returns:
//   - The intents that were configured for the application.
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
