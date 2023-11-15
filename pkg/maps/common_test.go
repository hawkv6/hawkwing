package maps

import "testing"

func TestFormatDNS(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    [256]byte
		wantErr bool
	}{
		{
			name: "root-domain",
			args: args{
				domain: ".",
			},
			want:    [256]byte{},
			wantErr: true,
		},
		{
			name: "valid-domain",
			args: args{
				domain: "www.example.com",
			},
			want:    [256]byte{3, 119, 119, 119, 7, 101, 120, 97, 109, 112, 108, 101, 3, 99, 111, 109, 0},
			wantErr: false,
		},
		{
			name: "empty-domain",
			args: args{
				domain: "",
			},
			want:    [256]byte{},
			wantErr: true,
		},
		{
			name: "invalid-domain",
			args: args{
				domain: "example",
			},
			want:    [256]byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatDNSName(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("error formatting domain name: %v", err)
			}
			if result != tt.want {
				t.Errorf("expected %v, got %v", tt.want, result)
			}
		})
	}
}

func TestIpv6ToInet6(t *testing.T) {
	ipv6Addr := "2001:db8::1"
	expected := struct{ In6U struct{ U6Addr8 [16]uint8 } }{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}}
	result := Ipv6ToInet6(ipv6Addr)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSidToInet6Sid(t *testing.T) {
	type args struct {
		sids []string
	}
	tests := []struct {
		name    string
		args    args
		want    [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
		wantErr bool
	}{
		{
			name: "non-empty-sid",
			args: args{
				sids: []string{"2001:db8::1", "2001:db8::2", "2001:db8::3"},
			},
			want: [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
			},
			wantErr: false,
		},
		{
			name: "too-long-sid",
			args: args{
				sids: []string{"2001:db8::1", "2001:db8::2", "2001:db8::3", "2001:db8::4", "2001:db8::5", "2001:db8::6", "2001:db8::7", "2001:db8::8", "2001:db8::9", "2001:db8::a", "2001:db8::b"},
			},
			want:    [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{},
			wantErr: true,
		},
		{
			name: "empty-sid",
			args: args{
				sids: []string{},
			},
			want: [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SidToInet6Sid(tt.args.sids)
			if (err != nil) != tt.wantErr {
				t.Errorf("error converting SID to inet6: %v", err)
			}
			if result != tt.want {
				t.Errorf("expected %v, got %v", tt.want, result)
			}
		})
	}
}

func TestGenerateSidLookupValue(t *testing.T) {
	type args struct {
		sids []string
	}
	tests := []struct {
		name    string
		args    args
		want    SidListData
		wantErr bool
	}{
		{
			name: "non-empty-sid",
			args: args{
				sids: []string{"2001:db8::1", "2001:db8::2", "2001:db8::3"},
			},
			want: SidListData{
				SidlistSize: 4,
				Sidlist: [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				},
			},
			wantErr: false,
		},
		{
			name: "empty-sid",
			args: args{
				sids: []string{},
			},
			want: SidListData{
				SidlistSize: 0,
				Sidlist: [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
					{In6U: struct{ U6Addr8 [16]uint8 }{U6Addr8: [16]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}},
				},
			},
			wantErr: false,
		},
		{
			name: "too-long-sid",
			args: args{
				sids: []string{"2001:db8::1", "2001:db8::2", "2001:db8::3", "2001:db8::4", "2001:db8::5", "2001:db8::6", "2001:db8::7", "2001:db8::8", "2001:db8::9", "2001:db8::a", "2001:db8::b"},
			},
			want:    SidListData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateSidLookupValue(tt.args.sids)
			if (err != nil) != tt.wantErr {
				t.Errorf("error generating SID lookup value: %v", err)
			}
			if result != tt.want {
				t.Errorf("expected %v, got %v", tt.want, result)
			}
		})
	}

}
