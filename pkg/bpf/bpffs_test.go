package bpf

import (
	"fmt"
	"testing"
)

func TestMapToPath(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic", args{"foo"}, "/sys/fs/bpf/foo"},
		{"slash", args{"/"}, "/sys/fs/bpf"},
		{"empty", args{""}, "/sys/fs/bpf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToPath(tt.args.name); got != tt.want {
				fmt.Println(got)
				t.Errorf("MapToPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
