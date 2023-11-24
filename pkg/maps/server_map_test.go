package maps

import (
	"fmt"
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
	"go.uber.org/mock/gomock"
)

var testServerLookupMapSpec = &ebpf.MapSpec{
	Name:       "server_lookup_map",
	Type:       ebpf.Array,
	KeySize:    4,
	ValueSize:  4,
	MaxEntries: 100,
}

func TestNewServerMap(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	sm, err := NewServerMap(mockBpf)
	if err != nil {
		t.Errorf("NewServerMap() should not return an error: %s", err)
	}
	if sm == nil {
		t.Error("NewServerMap() should not return nil")
	} else if sm.Lookup == nil {
		t.Error("NewServerMap() should not return nil for Lookup")
	}
}

func TestCreateServerMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBpf := bpf.NewMockBpf(ctrl)

	tests := []struct {
		name     string
		mapCreat func() *ServerMap
		mockBpf  func()
		wantErr  bool
	}{
		{
			name: "lookup map without spec",
			mapCreat: func() *ServerMap {
				return &ServerMap{
					bpf:    mockBpf,
					Lookup: NewMap(mockBpf, nil),
				}
			},
			mockBpf: func() {},
			wantErr: true,
		},
		{
			name: "create lookup map without error",
			mapCreat: func() *ServerMap {
				return &ServerMap{
					bpf:    mockBpf,
					Lookup: NewMap(mockBpf, testServerLookupMapSpec),
				}
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "create lookup map with error",
			mapCreat: func() *ServerMap {
				return &ServerMap{
					bpf:    mockBpf,
					Lookup: NewMap(mockBpf, testServerLookupMapSpec),
				}
			},
			mockBpf: func() {
				mockBpf.EXPECT().CreateMap(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBpf()
			sm := tt.mapCreat()
			err := sm.Create()
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
