package entities

import (
	"reflect"
	"testing"

	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/types"
)

func TestNewPathRequest(t *testing.T) {
	type args struct {
		ipv6daddr string
		intents   []Intent
	}
	tests := []struct {
		name string
		args args
		want PathRequest
	}{
		{
			name: "TestNewPathRequest",
			args: args{
				ipv6daddr: "2001:db8::1",
				intents: []Intent{
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
			want: PathRequest{
				Ipv6DestinationAddress: "2001:db8::1",
				Intents: []Intent{
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
		},
	}
	for _, tt := range tests {
		if got := NewPathRequest(tt.args.ipv6daddr, tt.args.intents); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("NewPathRequest() = %v, want %v", got, tt.want)
		}
	}

}

func TestPathRequest_Marshal(t *testing.T) {
	pr := NewPathRequest("2001:db8::1", []Intent{
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
	})
	if pr.Marshal() == nil {
		t.Errorf("PathRequest.Marshal() = nil")
	}
	if pr.Marshal().Ipv6DestinationAddress != "2001:db8::1" {
		t.Errorf("PathRequest.Marshal() = %v", pr.Marshal().Ipv6DestinationAddress)
	}
	if len(pr.Marshal().Intents) != 1 {
		t.Errorf("PathRequest.Marshal() = %v", len(pr.Marshal().Intents))
	}
	if pr.Marshal().Intents[0].Type != api.IntentType(types.IntentTypeSfc) {
		t.Errorf("PathRequest.Marshal() = %v", pr.Marshal().Intents[0].Type)
	}
	if len(pr.Marshal().Intents[0].Values) != 2 {
		t.Errorf("PathRequest.Marshal() = %v", len(pr.Marshal().Intents[0].Values))
	}
	if pr.Marshal().Intents[0].Values[0].Type != api.ValueType(types.IntentValueTypeSFC) {
		t.Errorf("PathRequest.Marshal() = %v", pr.Marshal().Intents[0].Values[0].Type)
	}
	if pr.Marshal().Intents[0].Values[1].Type != api.ValueType(types.IntentValueTypeSFC) {
		t.Errorf("PathRequest.Marshal() = %v", pr.Marshal().Intents[0].Values[1].Type)
	}

}

func TestCreatePathRequestsForService(t *testing.T) {
	test.SetupTestConfig(t)

	type args struct {
		serviceKey string
	}
	tests := []struct {
		name string
		args args
		want []PathRequest
	}{
		{
			name: "Mutliple intents",
			args: args{
				serviceKey: "service3",
			},
			want: []PathRequest{
				{
					Ipv6DestinationAddress: "fcbb:cc00:4::a",
					Intents: []Intent{
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
				{
					Ipv6DestinationAddress: "fcbb:cc00:4::b",
					Intents: []Intent{
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
				{
					Ipv6DestinationAddress: "fcbb:cc00:4::c",
					Intents: []Intent{
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
			},
		},
	}
	for _, tt := range tests {
		if got, err := CreatePathRequestsForService(tt.args.serviceKey); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CreatePathRequestsForService() = %v, want %v", got, tt.want)
		} else if err != nil {
			t.Errorf("CreatePathRequestsForService() = %v", err)
		}
		// if got, err := CreatePathRequestsForService(tt.args.serviceKey); !reflect.DeepEqual(got, tt.want) {
		// 	t.Errorf("CreatePathRequestsForService() = %v, want %v", got, tt.want)
		// }
	}
}
