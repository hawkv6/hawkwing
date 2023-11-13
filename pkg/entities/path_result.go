package entities

import (
	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/types"
)

type PathResult struct {
	Ipv6DestinationAddress string
	Intents                []Intent
	Ipv6SidAddresses       []string
}

func NewPathResult(ipv6daddr string, intents []Intent, ipv6SidAddresses []string) *PathResult {
	return &PathResult{
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
			intentValues = append(intentValues, IntentValue{
				IntentValueType: types.IntentValueType(val.Type),
				NumberValue:     *val.NumberValue,
				StringValue:     *val.StringValue,
			})
		}
		intents = append(intents, Intent{
			IntentType:   types.IntentType(intent.Type),
			IntentValues: intentValues,
		})
	}
	return &PathResult{
		Ipv6DestinationAddress: pr.Ipv6DestinationAddress,
		Intents:                intents,
		Ipv6SidAddresses:       pr.Ipv6SidAddresses,
	}
}
