package logging

import "testing"

func TestIntializeDefaultLogger(t *testing.T) {
	logger := InitializeDefaultLogger()
	if logger == nil {
		t.Errorf("logger is nil")
	}
}
