package config

import (
	"net"

	"github.com/go-playground/validator"
	"github.com/hawkv6/hawkwing/pkg/types"
)

func ApplicationValidation(sl validator.StructLevel) {
	application := sl.Current().Interface().(Application)

	if application.Port > 65535 || application.Port < 0 {
		sl.ReportError(application.Port, "port", "", "port must be between 0 and 65535", "")
	}

	if application.Sid != nil {
		for _, sid := range application.Sid {
			ip := net.ParseIP(sid)
			if ip == nil || ip.To4() != nil {
				sl.ReportError(sid, "sid", "", "sid must be a valid ipv6 address", "")
			}
		}
	}

	if application.Sid == nil && application.Intents == nil {
		sl.ReportError(application.Sid, "sid", "", "sid or intents must be specified", "")
	}

	for k, intent := range application.Intents {
		intentType, err := types.ParseIntentType(intent.Intent)
		if err != nil {
			sl.ReportError(intent.Intent, "intent", "", "invalid intent", "")
		}

		// sfc intent must be the first intent in the list
		if intentType == types.IntentTypeSfc {
			if k != 0 {
				sl.ReportError(intent.Intent, "intent", "", "sfc intent must be the first intent in the list", "")
			}

			if application.Sid == nil {
				sl.ReportError(intent.Intent, "sid", "", "sid is required as backup path when using an sfc intent", "")
			}
		}

		// flex-algo intent must be the first intent in the list if there is no sfc intent
		if intentType == types.IntentTypeFlexAlgo {
			if k >= 1 && application.Intents[0].Intent != types.IntentTypeSfc.String() {
				sl.ReportError(intent.Intent, "intent", "", "flex-algo intent must be the first intent in the list", "")
			}

			if application.Sid == nil {
				sl.ReportError(intent.Intent, "sid", "", "sid is required as backup path when using an flex-algo intent", "")
			}
		}
	}
}
