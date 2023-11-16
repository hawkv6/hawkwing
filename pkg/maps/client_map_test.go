package maps

import (
	"testing"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/client"
)

func TestNewClientMap(t *testing.T) {
	mockClientBpfReader := &client.MockClientBpfReader{}
	mockBpf := &bpf.MockBpf{}
	_, err := NewClientMap(mockBpf, mockClientBpfReader)
	if err != nil {
		t.Errorf("NewClientMap() error = %v", err)
	}
}
