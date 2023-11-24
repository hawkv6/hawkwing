package entities

import (
	"reflect"
	"testing"

	"github.com/hawkv6/hawkwing/pkg/api"
	"github.com/hawkv6/hawkwing/pkg/types"
)

func TestNewPathResult(t *testing.T) {
	type args struct {
		ipv6daddr        string
		intents          []Intent
		ipv6SidAddresses []string
	}
	tests := []struct {
		name string
		args args
		want *PathResult
	}{
		{
			name: "TestNewPathResult",
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
				ipv6SidAddresses: []string{
					"2001:db8::2",
					"2001:db8::3",
				},
			},
			want: &PathResult{
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
				Ipv6SidAddresses: []string{
					"2001:db8::2",
					"2001:db8::3",
				},
			},
		},
	}
	for _, tt := range tests {
		if got := NewPathResult(tt.args.ipv6daddr, tt.args.intents, tt.args.ipv6SidAddresses); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("NewPathResult() = %v, want %v", got, tt.want)
		}
	}
}

func TestUnmarshalPathResult(t *testing.T) {
	api_pr := &api.PathResult{
		Ipv6DestinationAddress: "2001:db8::1",
		Intents: []*api.Intent{
			{
				Type: api.IntentType_INTENT_TYPE_SFC,
			},
		},
		Ipv6SidAddresses: []string{
			"2001:db8::2",
			"2001:db8::3",
		},
	}

	wanted_pr := &PathResult{
		Ipv6DestinationAddress: "2001:db8::1",
		Intents: []Intent{
			{
				IntentType:   types.IntentTypeSfc,
				IntentValues: []IntentValue{},
			},
		},
		Ipv6SidAddresses: []string{
			"2001:db8::2",
			"2001:db8::3",
		},
	}

	if got := UnmarshalPathResult(api_pr); !reflect.DeepEqual(got, wanted_pr) {
		t.Errorf("UnmarshalPathResult() = %#v, want %#v", got, wanted_pr)
	}
}
