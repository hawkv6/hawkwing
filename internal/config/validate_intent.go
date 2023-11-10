package config

import (
	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkwing/pkg/types"
)

func IntentValidation(sl validator.StructLevel) {
	intent := sl.Current().Interface().(Intent)

	intentType, err := types.ParseIntentType(intent.Intent)
	if err != nil {
		sl.ReportError(intent.Intent, "intent", "", "invalid intent", "")
	}

	if intentType == types.IntentTypeFlexAlgo || intentType == types.IntentTypeSfc {
		if intent.MinValue > 0 {
			sl.ReportError(intent.MinValue, "min_value", "", "min_value and max_value are not allowed for flex-algo and sfc intents", "")
		}
		if intent.MaxValue > 0 {
			sl.ReportError(intent.MaxValue, "max_value", "", "min_value and max_value are not allowed for flex-algo and sfc intents", "")
		}
	}

	if intentType == types.IntentTypeFlexAlgo {
		if intent.FlexAlgoNr == 0 {
			sl.ReportError(intent.FlexAlgoNr, "flex_algo_number", "", "flex_algo_number is required when using an flex_algo intent", "")
		}
	}

	if intentType == types.IntentTypeSfc {
		if intent.Functions == nil {
			sl.ReportError(intent.Functions, "functions", "", "functions is required when using an sfc intent", "")
		}
	}

}
