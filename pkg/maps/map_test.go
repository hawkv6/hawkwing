package maps

import (
	"testing"

	"github.com/cilium/ebpf"
	"github.com/hawkv6/hawkwing/pkg/bpf"
)

func TestNewMap(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
	}
	m := NewMap(mockBpf, spec)
	if m == nil {
		t.Error("NewMap() should not return nil")
	}
}

func TestSetPath(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
	}
	m := NewMap(mockBpf, spec)
	m.setPath()
	if m.path != "/sys/fs/bpf" {
		t.Error("setPath() should set path to /sys/fs/bpf/test")
	}
}

func TestOpenOrCreate(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
		Pinning:   ebpf.PinByName,
	}
	m := NewMap(mockBpf, spec)
	err := m.OpenOrCreate()
	if err != nil {
		t.Errorf("OpenOrCreate() should not return an error: %s", err)
	}
	m = NewMap(mockBpf, nil)
	err = m.OpenOrCreate()
	if err == nil {
		t.Error("OpenOrCreate() should return an error when spec is nil")
	}
}

func TestOpen(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
		Pinning:   ebpf.PinByName,
	}
	m := NewMap(mockBpf, spec)
	err := m.Open()
	if err == nil {
		t.Error("Open() should return an error when path is empty")
	}
	m.setPath()
	err = m.Open()
	if err != nil {
		t.Errorf("Open() should not return an error: %s", err)
	}
}

func TestLookup(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
		Pinning:   ebpf.PinByName,
	}
	m := NewMap(mockBpf, spec)
	err := m.Lookup(nil, nil)
	if err == nil {
		t.Error("Lookup() should return an error when map is nil")
	}
	err = m.Lookup(&ebpf.Map{}, nil)
	if err == nil {
		t.Error("Lookup() should return an error when key is nil")
	}
	err = m.Lookup(&ebpf.Map{}, &ebpf.Map{})
	if err == nil {
		t.Error("Lookup() should return an error when path is empty")
	}
	m.setPath()
	err = m.Lookup(&ebpf.Map{}, &ebpf.Map{})
	if err != nil {
		t.Errorf("Lookup() should not return an error: %s", err)
	}
}

func TestUpdateInner(t *testing.T) {
	mockBpf := &bpf.MockBpf{}
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
		Pinning:   ebpf.PinByName,
	}
	m := NewMap(mockBpf, spec)

	type args struct {
		outerKey   interface{}
		innerKey   interface{}
		innerValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "outerKey is nil",
			args:    args{outerKey: nil, innerKey: nil, innerValue: nil},
			wantErr: true,
		},
		{
			name:    "innerKey is nil",
			args:    args{outerKey: &ebpf.Map{}, innerKey: nil, innerValue: nil},
			wantErr: true,
		},
		{
			name:    "innerValue is nil",
			args:    args{outerKey: &ebpf.Map{}, innerKey: &ebpf.Map{}, innerValue: nil},
			wantErr: true,
		},
		{
			name:    "path is empty",
			args:    args{outerKey: &ebpf.Map{}, innerKey: &ebpf.Map{}, innerValue: &ebpf.Map{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.UpdateInner(tt.args.outerKey, tt.args.innerKey, tt.args.innerValue); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
