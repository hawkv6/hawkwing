package maps

import (
	"testing"

	"github.com/cilium/ebpf"
)

func TestNewMap(t *testing.T) {
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
	}
	m := NewMap(spec)
	if m == nil {
		t.Error("NewMap() should not return nil")
	}
}

func TestSetPath(t *testing.T) {
	spec := &ebpf.MapSpec{
		Name:      "test",
		Type:      ebpf.Array,
		KeySize:   4,
		ValueSize: 4,
	}
	m := NewMap(spec)
	m.setPath()
	if m.path != "/sys/fs/bpf" {
		t.Error("setPath() should set path to /sys/fs/bpf/test")
	}
}
