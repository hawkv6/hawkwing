package entities

import (
	"testing"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/types"
)

func TestCreateIntentValueForIntent(t *testing.T) {
	type args struct {
		intent config.Intent
	}
	tests := []struct {
		name string
		args args
		want []IntentValue
	}{
		{
			name: "SFC",
			args: args{
				intent: config.Intent{
					Intent: types.IntentTypeSfc.String(),
					Functions: []string{
						"function1",
						"function2",
					},
				},
			},
			want: []IntentValue{
				{
					IntentValueType: types.IntentValueTypeSFC,
					StringValue:     "function1",
				},
				{
					IntentValueType: types.IntentValueTypeSFC,
					StringValue:     "function2",
				},
			},
		},
		{
			name: "FlexAlgo",
			args: args{
				intent: config.Intent{
					Intent:     types.IntentTypeFlexAlgo.String(),
					FlexAlgoNr: 1,
				},
			},
			want: []IntentValue{
				{
					IntentValueType: types.IntentValueTypeFlexAlgoNr,
					NumberValue:     1,
				},
			},
		},
		{
			name: "HighBandwidth",
			args: args{
				intent: config.Intent{
					Intent:   types.IntentTypeHighBandwidth.String(),
					MinValue: 1,
					MaxValue: 2,
				},
			},
			want: []IntentValue{
				{
					IntentValueType: types.IntentValueTypeMinValue,
					NumberValue:     1,
				},
				{
					IntentValueType: types.IntentValueTypeMaxValue,
					NumberValue:     2,
				},
			},
		},
		{
			name: "LowLatency",
			args: args{
				intent: config.Intent{
					Intent:   types.IntentTypeLowLatency.String(),
					MinValue: 1,
					MaxValue: 2,
				},
			},
			want: []IntentValue{
				{
					IntentValueType: types.IntentValueTypeMinValue,
					NumberValue:     1,
				},
				{
					IntentValueType: types.IntentValueTypeMaxValue,
					NumberValue:     2,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intentValues := CreateIntentValueForIntent(tt.args.intent)
			if len(intentValues) != len(tt.want) {
				t.Errorf("Expected %d intent values, got %d", len(tt.want), len(intentValues))
			}
			for i, intentValue := range intentValues {
				if intentValue.IntentValueType != tt.want[i].IntentValueType {
					t.Errorf("Expected intent value type %s, got %s", tt.want[i].IntentValueType, intentValue.IntentValueType)
				}
				if intentValue.NumberValue != tt.want[i].NumberValue {
					t.Errorf("Expected number value %d, got %d", tt.want[i].NumberValue, intentValue.NumberValue)
				}
				if intentValue.StringValue != tt.want[i].StringValue {
					t.Errorf("Expected string value %s, got %s", tt.want[i].StringValue, intentValue.StringValue)
				}
			}
		})
	}
}

func TestCreateIntentsForServiceApplication(t *testing.T) {
	test.SetupTestConfig(t)
	type args struct {
		serviceKey      string
		applicationPort int
	}

	tests := []struct {
		name string
		args args
		want []Intent
	}{
		{
			name: "SFC",
			args: args{
				serviceKey:      "service1",
				applicationPort: 80,
			},
			want: []Intent{
				{
					IntentType: types.IntentTypeSfc,
					IntentValues: []IntentValue{
						{
							IntentValueType: types.IntentValueTypeSFC,
							StringValue:     "function1",
						},
						{
							IntentValueType: types.IntentValueTypeSFC,
							StringValue:     "function2",
						},
					},
				},
			},
		},
		{
			name: "FlexAlgo",
			args: args{
				serviceKey:      "service1",
				applicationPort: 8080,
			},
			want: []Intent{
				{
					IntentType: types.IntentTypeFlexAlgo,
					IntentValues: []IntentValue{
						{
							IntentValueType: types.IntentValueTypeFlexAlgoNr,
							NumberValue:     1,
						},
					},
				},
			},
		},
		{
			name: "HighBandwidth",
			args: args{
				serviceKey:      "service2",
				applicationPort: 1433,
			},
			want: []Intent{
				{
					IntentType: types.IntentTypeHighBandwidth,
					IntentValues: []IntentValue{
						{
							IntentValueType: types.IntentValueTypeMinValue,
							NumberValue:     1,
						},
						{
							IntentValueType: types.IntentValueTypeMaxValue,
							NumberValue:     2,
						},
					},
				},
			},
		},
		{
			name: "MultipleIntents",
			args: args{
				serviceKey:      "service3",
				applicationPort: 443,
			},
			want: []Intent{
				{
					IntentType: types.IntentTypeFlexAlgo,
					IntentValues: []IntentValue{
						{
							IntentValueType: types.IntentValueTypeFlexAlgoNr,
							NumberValue:     1,
						},
					},
				},
				{
					IntentType: types.IntentTypeLowBandwidth,
					IntentValues: []IntentValue{
						{
							IntentValueType: types.IntentValueTypeMinValue,
							NumberValue:     1,
						},
						{
							IntentValueType: types.IntentValueTypeMaxValue,
							NumberValue:     2,
						},
					},
				},
				{
					IntentType: types.IntentTypeLowLatency,
					IntentValues: []IntentValue{
						{
							IntentValueType: types.IntentValueTypeMinValue,
							NumberValue:     1,
						},
						{
							IntentValueType: types.IntentValueTypeMaxValue,
							NumberValue:     2,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intents, err := CreateIntentsForServiceApplication(tt.args.serviceKey, tt.args.applicationPort)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(intents) != len(tt.want) {
				t.Errorf("Expected %d intents, got %d", len(tt.want), len(intents))
			}
			for i, intent := range intents {
				if intent.IntentType != tt.want[i].IntentType {
					t.Errorf("Expected intent type %s, got %s", tt.want[i].IntentType, intent.IntentType)
				}
				if len(intent.IntentValues) != len(tt.want[i].IntentValues) {
					t.Errorf("Expected %d intent values, got %d", len(tt.want[i].IntentValues), len(intent.IntentValues))
				}
				for j, intentValue := range intent.IntentValues {
					if intentValue.IntentValueType != tt.want[i].IntentValues[j].IntentValueType {
						t.Errorf("Expected intent value type %s, got %s", tt.want[i].IntentValues[j].IntentValueType, intentValue.IntentValueType)
					}
					if intentValue.NumberValue != tt.want[i].IntentValues[j].NumberValue {
						t.Errorf("Expected number value %d, got %d", tt.want[i].IntentValues[j].NumberValue, intentValue.NumberValue)
					}
					if intentValue.StringValue != tt.want[i].IntentValues[j].StringValue {
						t.Errorf("Expected string value %s, got %s", tt.want[i].IntentValues[j].StringValue, intentValue.StringValue)
					}
				}
			}
		})
	}

	application0Intent0 := config.Params.Services["service1"].Applications[0].Intents[0]
	application0Intent0.Intent = "invalid"
	config.Params.Services["service1"].Applications[0].Intents[0] = application0Intent0

	_, err := CreateIntentsForServiceApplication("service1", 80)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
