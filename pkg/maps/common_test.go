package maps

import "testing"

func TestFormatDNS(t *testing.T) {
	domain := "www.example.com"
	expected := [256]byte{3, 119, 119, 119, 7, 101, 120, 97, 109, 112, 108, 101, 3, 99, 111, 109, 0}
	result, err := FormatDNSName(domain)
	if err != nil {
		t.Errorf("error formatting domain name: %v", err)
	}
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
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
	sids := []string{"2001:db8::1", "2001:db8::2", "2001:db8::3"}
	expected := [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }{
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
	}
	result := SidToInet6Sid(sids)
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestGenerateSidLookupValue(t *testing.T) {
	type args struct {
		sids []string
	}
	tests := []struct {
		name string
		args args
		want SidListData
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSidLookupValue(tt.args.sids)
			if result != tt.want {
				t.Errorf("expected %v, got %v", tt.want, result)
			}
		})
	}

}
