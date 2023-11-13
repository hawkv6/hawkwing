package entities

import (
	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/api"
)

type PathRequest struct {
	Ipv6DestinationAddress string
	Intents                []Intent
}

func NewPathRequest(ipv6daddr string, intents []Intent) PathRequest {
	return PathRequest{
		Ipv6DestinationAddress: ipv6daddr,
		Intents:                intents,
	}
}

func (pr *PathRequest) Marshal() *api.PathRequest {
	intents := make([]*api.Intent, 0, len(pr.Intents))
	for _, intent := range pr.Intents {
		intentValues := make([]*api.Value, 0, len(intent.IntentValues))
		for _, val := range intent.IntentValues {
			iv := val
			intentValues = append(intentValues, &api.Value{
				Type:        api.ValueType(val.IntentValueType),
				NumberValue: &iv.NumberValue,
				StringValue: &iv.StringValue,
			})
		}
		intents = append(intents, &api.Intent{
			Type:   api.IntentType(intent.IntentType),
			Values: intentValues,
		})
	}
	return &api.PathRequest{
		Ipv6DestinationAddress: pr.Ipv6DestinationAddress,
		Intents:                intents,
	}
}

func CreatePathRequestsForService(serviceKey string) []PathRequest {
	serviceCfg := config.Params.Services[serviceKey]
	pathRequests := make([]PathRequest, 0, len(serviceCfg.Ipv6Addresses))
	for _, application := range serviceCfg.Applications {
		applicationIntents := CreateIntentsForServiceApplication(serviceKey, application.Port)
		for _, ipv6Addr := range serviceCfg.Ipv6Addresses {
			pathRequests = append(pathRequests, NewPathRequest(ipv6Addr, applicationIntents))
		}

	}
	return pathRequests
}
