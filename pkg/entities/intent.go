package entities

import (
	"log"

	"github.com/hawkv6/hawkwing/internal/config"
)

type IntentValue struct {
	IntentValueType IntentValueType
	NumberValue     int32
	StringValue     string
}

func CreateIntentValueForIntent(intent config.Intent) []IntentValue {
	intentValues := make([]IntentValue, 0)
	if intent.MinValue != 0 {
		intentValues = append(intentValues, IntentValue{
			IntentValueType: IntentValueTypeMinValue,
			NumberValue:     int32(intent.MinValue),
		})
	}
	if intent.MaxValue != 0 {
		intentValues = append(intentValues, IntentValue{
			IntentValueType: IntentValueTypeMaxValue,
			NumberValue:     int32(intent.MaxValue),
		})
	}
	if intent.Sfc != nil {
		for _, sfc := range intent.Sfc {
			intentValues = append(intentValues, IntentValue{
				IntentValueType: IntentValueTypeSFC,
				StringValue:     sfc,
			})
		}
	}
	if intent.FlexAlgo != 0 {
		intentValues = append(intentValues, IntentValue{
			IntentValueType: IntentValueTypeFlexAlgoNr,
			NumberValue:     int32(intent.FlexAlgo),
		})
	}
	return intentValues
}

type Intent struct {
	IntentType   IntentType
	IntentValues []IntentValue
}

func CreateIntentsForServiceApplication(serviceKey string, applicationPort int) []Intent {
	serviceCfg := config.Params.Services[serviceKey]
	intents := make([]Intent, 0)
	for _, application := range serviceCfg.Applications {
		if application.Port == applicationPort {
			for _, intent := range application.Intents {
				intentType, err := ParseIntentType(intent.Intent)
				if err != nil {
					log.Fatalf("failed to parse intent type %s: %v", intent.Intent, err)
				}
				intents = append(intents, Intent{
					IntentType:   intentType,
					IntentValues: CreateIntentValueForIntent(intent),
				})
			}
		}
	}
	return intents
}
