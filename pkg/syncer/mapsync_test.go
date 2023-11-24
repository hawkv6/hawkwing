package syncer

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
	"github.com/hawkv6/hawkwing/pkg/entities"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/hawkv6/hawkwing/pkg/types"
	"go.uber.org/mock/gomock"
)

func TestSyncer_getApplicationConfigIntents(t *testing.T) {
	test.SetupTestConfig(t)

	type args struct {
		intentResult *entities.PathResult
	}
	tests := []struct {
		name string
		args args
		want map[int][]string
	}{
		{
			name: "single intent",
			args: args{
				intentResult: &entities.PathResult{
					Ipv6DestinationAddress: "fcbb:cc00:4::a",
					Intents: []entities.Intent{
						{
							IntentType: types.IntentTypeFlexAlgo,
						},
						{
							IntentType: types.IntentTypeLowBandwidth,
						},
						{
							IntentType: types.IntentTypeLowLatency,
						},
					},
				},
			},
			want: map[int][]string{
				443:  {"flex-algo", "low-bandwidth", "low-latency"},
				8080: {"sfc"},
				18:   {"low-bandwidth"},
				19:   {"low-bandwidth"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Syncer{}
			got := s.getApplicationConfigIntents(tt.args.intentResult)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getApplicationConfigIntents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncer_getApplicationResultIntents(t *testing.T) {
	type args struct {
		intentResult *entities.PathResult
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single intent",
			args: args{
				intentResult: &entities.PathResult{
					Ipv6DestinationAddress: "fcbb:cc00:4::a",
					Intents: []entities.Intent{
						{
							IntentType: types.IntentTypeFlexAlgo,
						},
						{
							IntentType: types.IntentTypeLowBandwidth,
						},
						{
							IntentType: types.IntentTypeLowLatency,
						},
					},
				},
			},
			want: []string{"flex-algo", "low-bandwidth", "low-latency"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Syncer{}
			got := s.getResultIntents(tt.args.intentResult)
			if len(got) != len(tt.want) {
				t.Errorf("getApplicationResultIntents() got = %v, want %v", got, tt.want)
			}
			for _, w := range tt.want {
				if !slices.Contains(got, w) {
					t.Errorf("getApplicationResultIntents() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSyncer_getApplicationPortsToUpdate(t *testing.T) {
	test.SetupTestConfig(t)

	type args struct {
		intentResult *entities.PathResult
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "service3 - 443",
			args: args{
				intentResult: &entities.PathResult{
					Ipv6DestinationAddress: "fcbb:cc00:4::a",
					Intents: []entities.Intent{
						{
							IntentType: types.IntentTypeFlexAlgo,
						},
						{
							IntentType: types.IntentTypeLowBandwidth,
						},
						{
							IntentType: types.IntentTypeLowLatency,
						},
					},
				},
			},
			want: []int{443},
		},
		{
			name: "service3 - 8080",
			args: args{
				intentResult: &entities.PathResult{
					Ipv6DestinationAddress: "fcbb:cc00:4::b",
					Intents: []entities.Intent{
						{
							IntentType: types.IntentTypeSfc,
						},
					},
				},
			},
			want: []int{8080},
		},
		{
			name: "service3 - 18, 19",
			args: args{
				intentResult: &entities.PathResult{
					Ipv6DestinationAddress: "fcbb:cc00:4::c",
					Intents: []entities.Intent{
						{
							IntentType: types.IntentTypeLowBandwidth,
						},
					},
				},
			},
			want: []int{18, 19},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Syncer{}
			got := s.getApplicationPortsToUpdate(tt.args.intentResult)
			if len(got) != len(tt.want) {
				t.Errorf("getApplicationPortsToUpdate() got = %v, want %v", got, tt.want)
			}
			for _, w := range tt.want {
				if !slices.Contains(got, w) {
					t.Errorf("getApplicationPortsToUpdate() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSyncer_storeSidList(t *testing.T) {
	test.SetupTestConfig(t)
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)
	mockBpfReader := client.NewMockClientBpfReader(ctrl)
	mockBpfReader.EXPECT().ReadClientBpfSpecs().Return(&test.MockClientCollectionSpec, nil)
	clientMap, err := maps.NewClientMap(mockBpf, mockBpfReader)
	if err != nil {
		t.Fatalf("could not create client map: %v", err)
	}
	pathResult := &entities.PathResult{
		Ipv6DestinationAddress: "fcbb:cc00:4::a",
		Intents: []entities.Intent{
			{
				IntentType: types.IntentTypeFlexAlgo,
				IntentValues: []entities.IntentValue{
					{
						IntentValueType: types.IntentValueTypeFlexAlgoNr,
						NumberValue:     1,
					},
				},
			},
			{
				IntentType: types.IntentTypeLowBandwidth,
				IntentValues: []entities.IntentValue{
					{
						IntentValueType: types.IntentValueTypeMinValue,
						NumberValue:     1,
					},
				},
			},
			{
				IntentType: types.IntentTypeLowLatency,
				IntentValues: []entities.IntentValue{
					{
						IntentValueType: types.IntentValueTypeMaxValue,
						NumberValue:     2,
					},
				},
			},
		},
	}

	type args struct {
		intentResponse *entities.PathResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name: "Lookup returns error",
			args: args{
				intentResponse: pathResult,
			},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap(gomock.Any()).Return(&ebpf.Map{}, nil).AnyTimes()
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).AnyTimes()
			},
		},
		{
			name: "Generate sidlookup value returns error",
			args: args{
				intentResponse: pathResult,
			},
			wantErr: true,
			mockBpf: func() {
				pathResult.Ipv6SidAddresses = []string{"2001:db8::1", "2001:db8::2", "2001:db8::3", "2001:db8::4", "2001:db8::5", "2001:db8::6", "2001:db8::7", "2001:db8::8", "2001:db8::9", "2001:db8::a", "2001:db8::b"}
				mockBpf.EXPECT().LoadPinnedMap(gomock.Any()).Return(&ebpf.Map{}, nil).AnyTimes()
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			},
		},
		{
			name: "UpdateInner returns error",
			args: args{
				intentResponse: pathResult,
			},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap(gomock.Any()).Return(&ebpf.Map{}, nil).AnyTimes()
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockBpf.EXPECT().LoadMapFromId(gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()
			},
		},
		{
			name: "No error",
			args: args{
				intentResponse: pathResult,
			},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap(gomock.Any()).Return(&ebpf.Map{}, nil).AnyTimes()
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				mockBpf.EXPECT().LoadMapFromId(gomock.Any()).Return(nil, nil).AnyTimes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBpf()
			s := &Syncer{
				bpf: mockBpf,
				cm:  clientMap,
			}
			if err := s.storeSidList(tt.args.intentResponse); (err != nil) != tt.wantErr {
				t.Errorf("storeSidList() error = %v, wantErr %v", err, tt.wantErr)
			}
			ctrl.Finish()
		})
	}
}
