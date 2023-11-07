package entities

import "github.com/hawkv6/hawkwing/pkg/api"

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

func (pr *PathResult) Marshal() *api.PathResult {
	intents := make([]*api.Intent, 0, len(pr.Intents))
	for _, intent := range pr.Intents {
		intentValues := make([]*api.Value, 0, len(intent.IntentValues))
		for _, val := range intent.IntentValues {
			intentValues = append(intentValues, &api.Value{
				Type:        api.ValueType(val.IntentValueType),
				NumberValue: &val.NumberValue,
				StringValue: &val.StringValue,
			})
		}
		intents = append(intents, &api.Intent{
			Type:   api.IntentType(intent.IntentType),
			Values: intentValues,
		})
	}
	return &api.PathResult{
		Ipv6DestinationAddress: pr.Ipv6DestinationAddress,
		Intents:                intents,
		Ipv6SidAddresses:       pr.Ipv6SidAddresses,
	}
}

func UnmarshalPathResult(pr *api.PathResult) *PathResult {
	intents := make([]Intent, 0, len(pr.Intents))
	for _, intent := range pr.Intents {
		intentValues := make([]IntentValue, 0, len(intent.Values))
		for _, val := range intent.Values {
			intentValues = append(intentValues, IntentValue{
				IntentValueType: IntentValueType(val.Type),
				NumberValue:     *val.NumberValue,
				StringValue:     *val.StringValue,
			})
		}
		intents = append(intents, Intent{
			IntentType:   IntentType(intent.Type),
			IntentValues: intentValues,
		})
	}
	return &PathResult{
		Ipv6DestinationAddress: pr.Ipv6DestinationAddress,
		Intents:                intents,
		Ipv6SidAddresses:       pr.Ipv6SidAddresses,
	}
}
