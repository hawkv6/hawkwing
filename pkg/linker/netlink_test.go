package linker

import (
	"testing"

	"github.com/vishvananda/netlink"
)

func TestDirectionToParentDisc(t *testing.T) {
	type args struct {
		direction string
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "ingress",
			args: args{
				direction: "ingress",
			},
			want: netlink.HANDLE_MIN_INGRESS,
		},
		{
			name: "egress",
			args: args{
				direction: "egress",
			},
			want: netlink.HANDLE_MIN_EGRESS,
		},
		{
			name: "false",
			args: args{
				direction: "test",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		if got := directionToParentDisc(tt.args.direction); got != tt.want {
			t.Errorf("%q. directionToParentDisc() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
