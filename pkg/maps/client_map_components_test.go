package maps

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
	"go.uber.org/mock/gomock"
)

func TestNewInnerMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		bpf  bpf.Bpf
		spec *ebpf.MapSpec
	}
	tests := []struct {
		name string
		args args
		want *InnerMap
	}{
		{
			name: "creates new inner map",
			args: args{
				bpf:  mockBpf,
				spec: &mockClientInnerMapSpec,
			},
			want: &InnerMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientInnerMapSpec,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewInnerMap(tt.args.bpf, tt.args.spec)
			if !reflect.DeepEqual(im, tt.want) {
				t.Errorf("NewInnerMap() = %v, want %v", im, tt.want)
			}
		})
	}
}

func TestInnerMap_Build(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	tests := []struct {
		name    string
		wantErr bool
		mockBpf func()
	}{
		{
			name: "ebpf returns an error",
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "ebpf returns a map",
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &InnerMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientInnerMapSpec,
				},
			}
			tt.mockBpf()
			if err := im.Build(); (err != nil) != tt.wantErr {
				t.Errorf("InnerMap.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestNewOuterMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		bpf  bpf.Bpf
		spec *ebpf.MapSpec
	}
	tests := []struct {
		name string
		args args
		want *OuterMap
	}{
		{
			name: "creates new outer map",
			args: args{
				bpf:  mockBpf,
				spec: &mockClientOuterMapSpec,
			},
			want: &OuterMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientOuterMapSpec,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewOuterMap(tt.args.bpf, tt.args.spec)
			if !reflect.DeepEqual(im, tt.want) {
				t.Errorf("NewOuterMap() = %v, want %v", im, tt.want)
			}
		})
	}
}

func TestOuterMap_BuildWith(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		inners map[string]*InnerMap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name: "inner map creation returns an error",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "outer map creation returns an error",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "ebpf returns a map",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &OuterMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientOuterMapSpec,
				},
			}
			tt.mockBpf()
			if err := im.BuildWith(tt.args.inners); (err != nil) != tt.wantErr {
				t.Errorf("OuterMap.BuildWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewLookupMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		bpf  bpf.Bpf
		spec *ebpf.MapSpec
	}
	tests := []struct {
		name string
		args args
		want *LookupMap
	}{
		{
			name: "creates new lookup map",
			args: args{
				bpf:  mockBpf,
				spec: &mockClientLookupMapSpec,
			},
			want: &LookupMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientLookupMapSpec,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewLookupMap(tt.args.bpf, tt.args.spec)
			if !reflect.DeepEqual(im, tt.want) {
				t.Errorf("NewLookupMap() = %v, want %v", im, tt.want)
			}
		})
	}
}

func TestLookupMap_BuildWith(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		inners map[string]*InnerMap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name: "lookup map creation returns an error",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "ebpf returns a map",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &LookupMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientLookupMapSpec,
				},
			}
			tt.mockBpf()
			if err := im.BuildWith(tt.args.inners); (err != nil) != tt.wantErr {
				t.Errorf("LookupMap.BuildWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewReverseMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		bpf  bpf.Bpf
		spec *ebpf.MapSpec
	}
	tests := []struct {
		name string
		args args
		want *ReverseMap
	}{
		{
			name: "creates new reverse map",
			args: args{
				bpf:  mockBpf,
				spec: &mockClientReverseMapSpec,
			},
			want: &ReverseMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientReverseMapSpec,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := NewReverseMap(tt.args.bpf, tt.args.spec)
			if !reflect.DeepEqual(im, tt.want) {
				t.Errorf("NewReverseMap() = %v, want %v", im, tt.want)
			}
		})
	}
}

func TestReverseMap_BuildWith(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		inners map[string]*InnerMap
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name: "reverse map creation returns an error",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "ebpf returns a map",
			args: args{
				inners: map[string]*InnerMap{
					"test": {
						Map: Map{
							bpf:  mockBpf,
							spec: &mockClientInnerMapSpec,
						},
					},
				},
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			im := &ReverseMap{
				Map: Map{
					bpf:  mockBpf,
					spec: &mockClientReverseMapSpec,
				},
			}
			tt.mockBpf()
			if err := im.BuildWith(tt.args.inners); (err != nil) != tt.wantErr {
				t.Errorf("ReverseMap.BuildWith() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
