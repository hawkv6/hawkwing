package entities

import (
	"fmt"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/pkg/types"
)

type IntentValue struct {
	IntentValueType types.IntentValueType
	NumberValue     int32
	StringValue     string
}

func CreateIntentValueForIntent(intent config.Intent) []IntentValue {
	intentValues := make([]IntentValue, 0)
	if intent.Intent == types.IntentTypeSfc.String() {
		for _, function := range intent.Functions {
			intentValues = append(intentValues, IntentValue{
				IntentValueType: types.IntentValueTypeSFC,
				StringValue:     function,
			})
		}
		return intentValues
	}
	if intent.Intent == types.IntentTypeFlexAlgo.String() {
		intentValues = append(intentValues, IntentValue{
			IntentValueType: types.IntentValueTypeFlexAlgoNr,
			NumberValue:     int32(intent.FlexAlgoNr),
		})
		return intentValues
	}
	if intent.MinValue != 0 {
		intentValues = append(intentValues, IntentValue{
			IntentValueType: types.IntentValueTypeMinValue,
			NumberValue:     int32(intent.MinValue),
		})
	}
	if intent.MaxValue != 0 {
		intentValues = append(intentValues, IntentValue{
			IntentValueType: types.IntentValueTypeMaxValue,
			NumberValue:     int32(intent.MaxValue),
		})
	}
	return intentValues
}

type Intent struct {
	IntentType   types.IntentType
	IntentValues []IntentValue
}

func CreateIntentsForServiceApplication(serviceKey string, applicationPort int) ([]Intent, error) {
	serviceCfg := config.Params.Services[serviceKey]
	intents := make([]Intent, 0)
	for _, application := range serviceCfg.Applications {
		if application.Port == applicationPort {
			for _, intent := range application.Intents {
				intentType, err := types.ParseIntentType(intent.Intent)
				if err != nil {
					return nil, fmt.Errorf("failed to parse intent type %s: %v", intent.Intent, err)
				}
				intents = append(intents, Intent{
					IntentType:   intentType,
					IntentValues: CreateIntentValueForIntent(intent),
				})
			}
		}
	}
	return intents, nil
}
