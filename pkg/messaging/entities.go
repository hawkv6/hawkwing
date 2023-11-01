package messaging

import "github.com/hawkv6/hawkwing/pkg/api"

type IntentRequest struct {
	DomainName string
	IntentName string
}

type IntentResponse struct {
	SidList []string
}

func (i *IntentRequest) Marshal() *api.Intent {
	return &api.Intent{
		DomainName: i.DomainName,
		Intent:     intentStringToEnum(i.IntentName),
	}
}

func intentStringToEnum(intent string) api.IntentType {
	switch intent {
	case "high-bandwidth":
		return api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH
	case "low-bandwidth":
		return api.IntentType_INTENT_TYPE_LOW_BANDWIDTH
	case "low-latency":
		return api.IntentType_INTENT_TYPE_LOW_LATENCY
	default:
		return api.IntentType_INTENT_TYPE_UNSPECIFIED
	}
}
