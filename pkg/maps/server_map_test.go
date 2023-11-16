package maps

import (
	"testing"

	"github.com/hawkv6/hawkwing/pkg/bpf"
)

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
	mockBpf := &bpf.MockBpf{}
	sm, err := NewServerMap(mockBpf)
	if err != nil {
		t.Errorf("NewServerMap() should not return an error: %s", err)
	}
	err = sm.Create()
	if err != nil {
		t.Errorf("Create() should not return an error: %s", err)
	}
}
