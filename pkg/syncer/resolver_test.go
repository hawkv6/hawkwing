package syncer

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/hawkv6/hawkwing/internal/config"
	"github.com/hawkv6/hawkwing/internal/test"
	"github.com/hawkv6/hawkwing/pkg/entities"
	"go.uber.org/mock/gomock"
)

func TestResolver_NewResolverService(t *testing.T) {
	reqChan := make(chan *entities.PathRequest, 1)
	rs := NewResolverService(reqChan)

	if rs.dNSResolver != net.DefaultResolver {
		t.Errorf("NewResolverService() dNSResolver = %v, want %v", rs.dNSResolver, net.DefaultResolver)
	}
	if rs.cache == nil {
		t.Errorf("NewResolverService() cache = %v, want %v", rs.cache, make(map[string][]net.IPAddr))
	}
	if rs.cacheTTL == nil {
		t.Errorf("NewResolverService() cacheTTL = %v, want %v", rs.cacheTTL, make(map[string]time.Time))
	}
	if rs.reqChan != reqChan {
		t.Errorf("NewResolverService() reqChan = %v, want %v", rs.reqChan, reqChan)
	}
}

func TestResolver_ResolveIPv6Addresses(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDnsResolver := NewMockDNSResolver(ctrl)

	tests := []struct {
		name           string
		domainName     string
		mockResolverFn func()
		want           []net.IPAddr
		wantErr        bool
	}{
		{
			name:       "success",
			domainName: "test.com",
			mockResolverFn: func() {
				mockDnsResolver.EXPECT().LookupIPAddr(gomock.Any(), "test.com").Return([]net.IPAddr{
					{
						IP:   net.ParseIP("2001:db8::68"),
						Zone: "",
					},
				}, nil)
			},
			want: []net.IPAddr{
				{
					IP:   net.ParseIP("2001:db8::68"),
					Zone: "",
				},
			},
			wantErr: false,
		},
		{
			name:       "error",
			domainName: "test.com",
			mockResolverFn: func() {
				mockDnsResolver.EXPECT().LookupIPAddr(gomock.Any(), "test.com").Return(nil, net.UnknownNetworkError("unknown network"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockResolverFn()

			rs := &ResolverService{
				dNSResolver: mockDnsResolver,
				cache:       make(map[string][]net.IPAddr),
				cacheTTL:    make(map[string]time.Time),
				reqChan:     make(chan *entities.PathRequest, 1),
			}

			got, err := rs.ResolveIPv6Addresses(context.Background(), tt.domainName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveIPv6Addresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveIPv6Addresses() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_ResolveIPCacheHit(t *testing.T) {
	resolver := &ResolverService{
		cache: map[string][]net.IPAddr{
			"test.com": {
				{
					IP:   net.ParseIP("2001:db8::68"),
					Zone: "",
				},
			},
		},
		cacheTTL: map[string]time.Time{
			"test.com": time.Now().Add(10 * time.Second),
		},
	}

	got, err := resolver.ResolveIPv6Addresses(context.Background(), "test.com")
	if err != nil {
		t.Errorf("ResolveIPv6Addresses() error = %v", err)
		return
	}
	want := []net.IPAddr{
		{
			IP:   net.ParseIP("2001:db8::68"),
			Zone: "",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveIPv6Addresses() got = %v, want %v", got, want)
	}

}

func TestResolver_ProcessConfig(t *testing.T) {
	tests := []struct {
		name           string
		mockResolverFn func(*ResolverService, *MockDNSResolver)
	}{
		{
			name: "no error",
			mockResolverFn: func(resolverService *ResolverService, mockDnsResolver *MockDNSResolver) {
				mockDnsResolver.EXPECT().LookupIPAddr(gomock.Any(), gomock.Any()).Return([]net.IPAddr{
					{
						IP:   net.ParseIP("2001:db8::68"),
						Zone: "",
					},
				}, nil).AnyTimes()

				err := resolverService.ProcessConfig(context.Background())
				if err != nil {
					t.Errorf("ProcessConfig() error = %v", err)
				}
			},
		},
		{
			name: "create path request returns error",
			mockResolverFn: func(resolverService *ResolverService, mockDnsResolver *MockDNSResolver) {
				mockDnsResolver.EXPECT().LookupIPAddr(gomock.Any(), gomock.Any()).Return([]net.IPAddr{
					{
						IP:   net.ParseIP("2001:db8::68"),
						Zone: "",
					},
				}, nil).AnyTimes()
				config.Params.Services["service1"].Applications[0].Intents[0].Intent = "invalid"
				err := resolverService.ProcessConfig(context.Background())
				if err == nil {
					t.Errorf("ProcessConfig() error = %v", err)
				}
			},
		},
		{
			name: "resolve ipv6 address error",
			mockResolverFn: func(resolverService *ResolverService, mockDnsResolver *MockDNSResolver) {
				mockDnsResolver.EXPECT().LookupIPAddr(gomock.Any(), gomock.Any()).Return(nil, net.UnknownNetworkError("unknown network")).AnyTimes()
				err := resolverService.ProcessConfig(context.Background())
				if err == nil {
					t.Errorf("ProcessConfig() error = %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.SetupTestConfig(t)
			ctrl := gomock.NewController(t)
			mockDnsResolver := NewMockDNSResolver(ctrl)
			rs := &ResolverService{
				dNSResolver: mockDnsResolver,
				cache:       make(map[string][]net.IPAddr),
				cacheTTL:    make(map[string]time.Time),
				reqChan:     make(chan *entities.PathRequest, 10),
			}
			tt.mockResolverFn(rs, mockDnsResolver)
		})
	}
}
