package entities

import (
	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/types"
)

type PathResult struct {
	Ipv6SourceAddress      string
	Ipv6DestinationAddress string
	Intents                []Intent
	Ipv6SidAddresses       []string
}

func NewPathResult(ipv6saddr string, ipv6daddr string, intents []Intent, ipv6SidAddresses []string) *PathResult {
	return &PathResult{
		Ipv6SourceAddress:      ipv6saddr,
		Ipv6DestinationAddress: ipv6daddr,
		Intents:                intents,
		Ipv6SidAddresses:       ipv6SidAddresses,
	}
}

func UnmarshalPathResult(pr *api.PathResult) *PathResult {
	intents := make([]Intent, 0, len(pr.Intents))
	for _, intent := range pr.Intents {
		intentValues := make([]IntentValue, 0, len(intent.Values))
		for _, val := range intent.Values {
			numberValue := int32(0)
			stringValue := ""
			if val.NumberValue != nil {
				numberValue = *val.NumberValue
			}
			if val.StringValue != nil {
				stringValue = *val.StringValue
			}
			intentValues = append(intentValues, IntentValue{
				IntentValueType: types.IntentValueType(val.Type),
				NumberValue:     numberValue,
				StringValue:     stringValue,
			})
		}
		intents = append(intents, Intent{
			IntentType:   types.IntentType(intent.Type),
			IntentValues: intentValues,
		})
	}
	return &PathResult{
		Ipv6SourceAddress:      pr.Ipv6SourceAddress,
		Ipv6DestinationAddress: pr.Ipv6DestinationAddress,
		Intents:                intents,
		Ipv6SidAddresses:       pr.Ipv6SidAddresses,
	}
}
