package maps

import (
	"fmt"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
	"go.uber.org/mock/gomock"
)

var testMapSpec = &ebpf.MapSpec{
	Name:      "test",
	Type:      ebpf.Array,
	KeySize:   4,
	ValueSize: 4,
	Pinning:   ebpf.PinByName,
}

func TestNewMap(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	m := NewMap(mockBpf, testMapSpec)
	if m == nil {
		t.Error("NewMap() should not return nil")
	}

}

func TestSetPath(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	m := NewMap(mockBpf, testMapSpec)
	m.setPath()
	if m.path != "/sys/fs/bpf" {
		t.Error("setPath() should set path to /sys/fs/bpf/test")
	}
}

func TestOpenOrCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		pin bool
	}
	tests := []struct {
		name        string
		args        args
		mapCreation func() *Map
		wantErr     bool
		mockBpf     func()
	}{
		{
			name: "spec is nil",
			args: args{pin: true},
			mapCreation: func() *Map {
				return NewMap(mockBpf, nil)
			},
			wantErr: true,
			mockBpf: func() {},
		},
		{
			name: "CreateMap returns error",
			args: args{pin: true},
			mapCreation: func() *Map {
				return NewMap(mockBpf, testMapSpec)
			},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(testMapSpec, "/sys/fs/bpf").Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name: "no error",
			args: args{pin: true},
			mapCreation: func() *Map {
				return NewMap(mockBpf, testMapSpec)
			},
			wantErr: false,
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(testMapSpec, "/sys/fs/bpf").Return(nil, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBpf()
			m := tt.mapCreation()
			err := m.OpenOrCreate()
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenOrCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name:    "path is empty",
			args:    args{path: ""},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap("").Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name:    "open returns error",
			args:    args{path: "/sys/fs/bpf"},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap("/sys/fs/bpf").Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBpf()
			m := NewMap(mockBpf, testMapSpec)
			m.path = tt.args.path
			err := m.Open()
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	type args struct {
		key   interface{}
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name:    "key is nil",
			args:    args{key: nil, value: nil},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap("").Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name:    "value is nil",
			args:    args{key: &ebpf.Map{}, value: nil},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap("").Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name:    "path is empty",
			args:    args{key: &ebpf.Map{}, value: &ebpf.Map{}},
			wantErr: true,
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap("").Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name: "open returns error",
			args: args{key: &ebpf.Map{}, value: &ebpf.Map{}},
			mockBpf: func() {
				mockBpf.EXPECT().LoadPinnedMap("").Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBpf()
			m := NewMap(mockBpf, testMapSpec)
			err := m.Lookup(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateInner(t *testing.T) {
	var mockBpf *bpf.MockBpf
	var m *Map
	ctrl := gomock.NewController(t)

	type args struct {
		outerKey   interface{}
		innerKey   interface{}
		innerValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockBpf func()
	}{
		{
			name: "lookup fails",
			args: args{outerKey: nil, innerKey: nil, innerValue: nil},
			mockBpf: func() {
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "inner map loading fails",
			args: args{outerKey: &ebpf.Map{}, innerKey: &ebpf.Map{}, innerValue: &ebpf.Map{}},
			mockBpf: func() {
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockBpf.EXPECT().LoadMapFromId(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "update fails",
			args: args{outerKey: &ebpf.Map{}, innerKey: &ebpf.Map{}, innerValue: &ebpf.Map{}},
			mockBpf: func() {
				mockBpf.EXPECT().LookupMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockBpf.EXPECT().LoadMapFromId(gomock.Any()).Return(nil, nil)
				mockBpf.EXPECT().PutMap(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBpf = bpf.NewMockBpf(ctrl)
			m = NewMap(mockBpf, testMapSpec)
			m.m = &ebpf.Map{}
			tt.mockBpf()
			err := m.UpdateInner(tt.args.outerKey, tt.args.innerKey, tt.args.innerValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateInner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
