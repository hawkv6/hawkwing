package types

import (
	"fmt"

	"github.com/hawkv6/hawkwing/pkg/api"
)

type IntentValueType api.ValueType

const (
	IntentValueTypeUnspecified = IntentValueType(api.ValueType_VALUE_TYPE_UNSPECIFIED)
	IntentValueTypeMinValue    = IntentValueType(api.ValueType_VALUE_TYPE_MIN_VALUE)
	IntentValueTypeMaxValue    = IntentValueType(api.ValueType_VALUE_TYPE_MAX_VALUE)
	IntentValueTypeSFC         = IntentValueType(api.ValueType_VALUE_TYPE_SFC)
	IntentValueTypeFlexAlgoNr  = IntentValueType(api.ValueType_VALUE_TYPE_FLEX_ALGO_NR)
)

func (ivt IntentValueType) String() string {
	switch ivt {
	case IntentValueTypeMinValue:
		return "min-value"
	case IntentValueTypeMaxValue:
		return "max-value"
	case IntentValueTypeSFC:
		return "sfc"
	case IntentValueTypeFlexAlgoNr:
		return "flex-algo-nr"
	default:
		return "unspecified"
	}
}

func ParseIntentValueType(s string) (IntentValueType, error) {
	switch s {
	case "min-value":
		return IntentValueTypeMinValue, nil
	case "max-value":
		return IntentValueTypeMaxValue, nil
	case "sfc":
		return IntentValueTypeSFC, nil
	case "flex-algo-nr":
		return IntentValueTypeFlexAlgoNr, nil
	case "unspecified":
		return IntentValueTypeUnspecified, nil
	default:
		return 0, fmt.Errorf("invalid IntentValueType: %s", s)
	}
}

type IntentType api.IntentType

const (
	IntentTypeUnspecified    = IntentType(api.IntentType_INTENT_TYPE_UNSPECIFIED)
	IntentTypeHighBandwidth  = IntentType(api.IntentType_INTENT_TYPE_HIGH_BANDWIDTH)
	IntentTypeLowBandwidth   = IntentType(api.IntentType_INTENT_TYPE_LOW_BANDWIDTH)
	IntentTypeLowLatency     = IntentType(api.IntentType_INTENT_TYPE_LOW_LATENCY)
	IntentTypeLowPacketLoss  = IntentType(api.IntentType_INTENT_TYPE_LOW_PACKET_LOSS)
	IntentTypeLowJitter      = IntentType(api.IntentType_INTENT_TYPE_LOW_JITTER)
	IntentTypeFlexAlgo       = IntentType(api.IntentType_INTENT_TYPE_FLEX_ALGO)
	IntentTypeSfc            = IntentType(api.IntentType_INTENT_TYPE_SFC)
	IntentTypeLowUtilization = IntentType(api.IntentType_INTENT_TYPE_LOW_UTILIZATION)
)

func (it IntentType) String() string {
	switch it {
	case IntentTypeHighBandwidth:
		return "high-bandwidth"
	case IntentTypeLowBandwidth:
		return "low-bandwidth"
	case IntentTypeLowLatency:
		return "low-latency"
	case IntentTypeLowPacketLoss:
		return "low-packet-loss"
	case IntentTypeLowJitter:
		return "low-jitter"
	case IntentTypeFlexAlgo:
		return "flex-algo"
	case IntentTypeSfc:
		return "sfc"
	case IntentTypeLowUtilization:
		return "low-utilization"
	default:
		return "unspecified"
	}
}

func ParseIntentType(s string) (IntentType, error) {
	switch s {
	case "high-bandwidth":
		return IntentTypeHighBandwidth, nil
	case "low-bandwidth":
		return IntentTypeLowBandwidth, nil
	case "low-latency":
		return IntentTypeLowLatency, nil
	case "low-packet-loss":
		return IntentTypeLowPacketLoss, nil
	case "low-jitter":
		return IntentTypeLowJitter, nil
	case "unspecified":
		return IntentTypeUnspecified, nil
	case "flex-algo":
		return IntentTypeFlexAlgo, nil
	case "sfc":
		return IntentTypeSfc, nil
	case "low-utilization":
		return IntentTypeLowUtilization, nil
	default:
		return 0, fmt.Errorf("invalid IntentType: %s", s)
	}
}
