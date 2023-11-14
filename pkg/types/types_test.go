package types

import "testing"

func TestIntentValueTypeToString(t *testing.T) {
	type args struct {
		ivt IntentValueType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "IntentValueTypeMinValue",
			args: args{
				ivt: IntentValueTypeMinValue,
			},
			want: "min-value",
		},
		{
			name: "IntentValueTypeMaxValue",
			args: args{
				ivt: IntentValueTypeMaxValue,
			},
			want: "max-value",
		},
		{
			name: "IntentValueTypeSFC",
			args: args{
				ivt: IntentValueTypeSFC,
			},
			want: "sfc",
		},
		{
			name: "IntentValueTypeFlexAlgoNr",
			args: args{
				ivt: IntentValueTypeFlexAlgoNr,
			},
			want: "flex-algo-nr",
		},
		{
			name: "IntentValueTypeUnspecified",
			args: args{
				ivt: IntentValueTypeUnspecified,
			},
			want: "unspecified",
		},
		{
			name: "IntentValueTypeInvalid",
			args: args{
				ivt: IntentValueType(100),
			},
			want: "unspecified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.ivt.String(); got != tt.want {
				t.Errorf("IntentValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIntentValueType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    IntentValueType
		wantErr bool
	}{
		{
			name: "IntentValueTypeMinValue",
			args: args{
				s: "min-value",
			},
			want:    IntentValueTypeMinValue,
			wantErr: false,
		},
		{
			name: "IntentValueTypeMaxValue",
			args: args{
				s: "max-value",
			},
			want:    IntentValueTypeMaxValue,
			wantErr: false,
		},
		{
			name: "IntentValueTypeSFC",
			args: args{
				s: "sfc",
			},
			want:    IntentValueTypeSFC,
			wantErr: false,
		},
		{
			name: "IntentValueTypeFlexAlgoNr",
			args: args{
				s: "flex-algo-nr",
			},
			want:    IntentValueTypeFlexAlgoNr,
			wantErr: false,
		},
		{
			name: "IntentValueTypeUnspecified",
			args: args{
				s: "unspecified",
			},
			want:    IntentValueTypeUnspecified,
			wantErr: false,
		},
		{
			name: "IntentValueTypeInvalid",
			args: args{
				s: "invalid",
			},
			want:    IntentValueTypeUnspecified,
			wantErr: true,
		},
		{
			name: "IntentValueTypeEmpty",
			args: args{
				s: "",
			},
			want:    IntentValueTypeUnspecified,
			wantErr: true,
		},
		{
			name: "IntentValueTypeSpace",
			args: args{
				s: " ",
			},
			want:    IntentValueTypeUnspecified,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntentValueType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntentValueType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want && err == nil {
				t.Errorf("ParseIntentValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntentTypeToString(t *testing.T) {
	type args struct {
		it IntentType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "IntentTypeHighBandwidth",
			args: args{
				it: IntentTypeHighBandwidth,
			},
			want: "high-bandwidth",
		},
		{
			name: "IntentTypeLowBandwidth",
			args: args{
				it: IntentTypeLowBandwidth,
			},
			want: "low-bandwidth",
		},
		{
			name: "IntentTypeLowLatency",
			args: args{
				it: IntentTypeLowLatency,
			},
			want: "low-latency",
		},
		{
			name: "IntentTypeLowPacketLoss",
			args: args{
				it: IntentTypeLowPacketLoss,
			},
			want: "low-packet-loss",
		},
		{
			name: "IntentTypeLowJitter",
			args: args{
				it: IntentTypeLowJitter,
			},
			want: "low-jitter",
		},
		{
			name: "IntentTypeFlexAlgo",
			args: args{
				it: IntentTypeFlexAlgo,
			},
			want: "flex-algo",
		},
		{
			name: "IntentTypeSfc",
			args: args{
				it: IntentTypeSfc,
			},
			want: "sfc",
		},
		{
			name: "IntentTypeUnspecified",
			args: args{
				it: IntentTypeUnspecified,
			},
			want: "unspecified",
		},
		{
			name: "IntentTypeInvalid",
			args: args{
				it: IntentType(100),
			},
			want: "unspecified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.it.String(); got != tt.want {
				t.Errorf("IntentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseIntentType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    IntentType
		wantErr bool
	}{
		{
			name: "IntentTypeHighBandwidth",
			args: args{
				s: "high-bandwidth",
			},
			want:    IntentTypeHighBandwidth,
			wantErr: false,
		},
		{
			name: "IntentTypeLowBandwidth",
			args: args{
				s: "low-bandwidth",
			},
			want:    IntentTypeLowBandwidth,
			wantErr: false,
		},
		{
			name: "IntentTypeLowLatency",
			args: args{
				s: "low-latency",
			},
			want:    IntentTypeLowLatency,
			wantErr: false,
		},
		{
			name: "IntentTypeLowPacketLoss",
			args: args{
				s: "low-packet-loss",
			},
			want:    IntentTypeLowPacketLoss,
			wantErr: false,
		},
		{
			name: "IntentTypeLowJitter",
			args: args{
				s: "low-jitter",
			},
			want:    IntentTypeLowJitter,
			wantErr: false,
		},
		{
			name: "IntentTypeFlexAlgo",
			args: args{
				s: "flex-algo",
			},
			want:    IntentTypeFlexAlgo,
			wantErr: false,
		},
		{
			name: "IntentTypeSfc",
			args: args{
				s: "sfc",
			},
			want:    IntentTypeSfc,
			wantErr: false,
		},
		{
			name: "IntentTypeUnspecified",
			args: args{
				s: "unspecified",
			},
			want:    IntentTypeUnspecified,
			wantErr: false,
		},
		{
			name: "IntentTypeInvalid",
			args: args{
				s: "invalid",
			},
			want:    IntentTypeUnspecified,
			wantErr: true,
		},
		{
			name: "IntentTypeEmpty",
			args: args{
				s: "",
			},
			want:    IntentTypeUnspecified,
			wantErr: true,
		},
		{
			name: "IntentTypeSpace",
			args: args{
				s: " ",
			},
			want:    IntentTypeUnspecified,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntentType(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIntentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want && err == nil {
				t.Errorf("ParseIntentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
