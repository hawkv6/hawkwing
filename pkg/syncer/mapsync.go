package syncer

import (
	"github.com/hawkv6/hawkwing/internal/config"
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
	// TODO: implement

	// ports := s.getPortsToQuery(intentResponse.DomainName, intentResponse.IntentName)
	// formattedDomainName, err := maps.FormatDNSName(intentResponse.DomainName)
	// if err != nil {
	// 	return fmt.Errorf("could not format domain name: %s", err)
	// }

	// path := "/sys/fs/bpf/hawkwing/client_outer_map"
	// om, err := ebpf.LoadPinnedMap(path, nil)
	// if err != nil {
	// 	return fmt.Errorf("could not load outer map: %s", err)
	// }
	// lm, err := ebpf.LoadPinnedMap(s.cm.Lookup.Path, nil)
	// if err != nil {
	// 	return fmt.Errorf("could not load lookup map: %s", err)
	// }

	// innerMapId := uint32(0)
	// err = lm.Lookup(formattedDomainName, &innerMapId)
	// if err != nil {
	// 	return fmt.Errorf("could not find inner map for domain %s", intentResponse.DomainName)
	// }

	// var innerMap *ebpf.Map
	// err = om.Lookup(uint32(innerMapId), &innerMap)
	// if err != nil {
	// 	return fmt.Errorf("could not find inner map for domain %s", intentResponse.DomainName)
	// }

	// sidListData := maps.GenerateSidLookupValue(intentResponse.SidList)

	// for _, port := range ports {
	// 	err = innerMap.Put(uint32(port), sidListData)
	// 	if err != nil {
	// 		return fmt.Errorf("could not update inner map: %s", err)
	// 	}
	// }
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
