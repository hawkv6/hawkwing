package syncer

import (
	"net"
	"time"

	"github.com/hawkv6/hawkwing/pkg/entities"
)

type ResolverService struct {
	dNSResolver *net.Resolver
	cache       map[string][]net.IPAddr
	cacheTTL    map[string]time.Time
	reqChan     chan *entities.PathRequest
}

func NewResolverService(reqChan chan *entities.PathRequest) *ResolverService {
	return &ResolverService{
		dNSResolver: net.DefaultResolver,
		cache:       make(map[string][]net.IPAddr),
		cacheTTL:    make(map[string]time.Time),
		reqChan:     reqChan,
	}
}

// func (rs *ResolverService) Start() error {
// 	return rs.ProcessConfig(context.Background())
// }

// func (rs *ResolverService) ProcessConfig(ctx context.Context) error {
// 	for key, svcConfig := range config.Params.Services {
// 		if len(svcConfig.Ipv6Addresses) == 0 {
// 			ipv6Addresses, err := rs.ResolveIPv6Addresses(ctx, svcConfig.DomainName)
// 			if err != nil {
// 				return err
// 			}

// 			var ipv6AddrStrs []string
// 			for _, ipAddr := range ipv6Addresses {
// 				ipv6AddrStrs = append(ipv6AddrStrs, ipAddr.String())
// 			}
// 			svcConfig.Ipv6Addresses = ipv6AddrStrs
// 		}
// 	}
// 	return nil
// }

// func (rs *ResolverService) ResolveIPv6Addresses(ctx context.Context, domainName string) ([]net.IPAddr, error) {
// 	if ips, found := rs.cache[domainName]; found {
// 		if time.Now().Before(rs.cacheTTL[domainName]) {
// 			return ips, nil
// 		}
// 	}

// 	ips, err := rs.dNSResolver.LookupIPAddr(ctx, domainName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var ipv6Addresses []net.IPAddr
// 	for _, ip := range ips {
// 		if ip.IP.To4() == nil {
// 			ipv6Addresses = append(ipv6Addresses, ip)
// 		}
// 	}

// 	ttl := time.Now().Add(1 * time.Hour)
// 	rs.cache[domainName] = ipv6Addresses
// 	rs.cacheTTL[domainName] = ttl

// 	return ipv6Addresses, nil
// }
