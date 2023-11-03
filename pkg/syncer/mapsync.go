package syncer

import (
	"fmt"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/messaging"
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
func (s *Syncer) storeSidList(intentResponse *messaging.IntentResponse) error {
	ports := s.getPortsToQuery(intentResponse.DomainName, intentResponse.IntentName)
	formattedDomainName, err := maps.FormatDNSName(intentResponse.DomainName)
	if err != nil {
		return fmt.Errorf("could not format domain name: %s", err)
	}

	var lookupValue uint32
	err = s.cm.Lookup.Map.Lookup(formattedDomainName, &lookupValue)
	if err != nil {
		return fmt.Errorf("could not find a value in the lookup map for %s: %s", intentResponse.DomainName, err)
	}

	sidListData := maps.GenerateSidLookupValue(intentResponse.SidList)

	for _, port := range ports {
		err = s.cm.Outer.UpdateInner(lookupValue, uint16(port), sidListData)
		if err != nil {
			return fmt.Errorf("could not update inner map: %s", err)
		}
	}
	return nil
}

// getPortsToQuery returns a list of ports to query for the given domain name
// and intent.
//
// Parameters:
//   - domainName: The domain name to query.
//   - intent: The intent to query.
//
// Returns:
//   - A list of ports to query.
func (s *Syncer) getPortsToQuery(domainName string, intent string) []int {
	var ports []int
	for _, service := range config.Params.Services[domainName] {
		if service.Intent == intent && service.Sid == nil {
			ports = append(ports, service.Port)
		}
	}
	return ports
}
