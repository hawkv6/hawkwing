package syncer

import (
	"context"
	"net"
	"time"

	"github.com/hawkv6/hawkwing/internal/config"
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

func (rs *ResolverService) Start() {
	ticker := time.NewTicker(10 * time.Second)
	if err := rs.ProcessConfig(context.Background()); err != nil {
		log.WithError(err).Error("could not process config")
	}
	go func() {
		for {
			<-ticker.C
			ctx := context.Background()
			if err := rs.ProcessConfig(ctx); err != nil {
				log.WithError(err).Error("could not process config")
			}
		}
	}()
}

func (rs *ResolverService) ProcessConfig(ctx context.Context) error {
	for key, serviceCfg := range config.Params.Services {
		if serviceCfg.Ipv6Addresses == nil {
			ipv6Addresses, err := rs.ResolveIPv6Addresses(ctx, serviceCfg.DomainName)
			if err != nil {
				return err
			}
			var ipv6AddrStrs []string
			for _, ipAddr := range ipv6Addresses {
				ipv6AddrStrs = append(ipv6AddrStrs, ipAddr.String())
			}

			serviceCfg.Ipv6Addresses = ipv6AddrStrs
			config.Params.Services[key] = serviceCfg

			pathRequests := entities.CreatePathRequestsForService(key)
			for _, pathRequest := range pathRequests {
				pr := pathRequest
				rs.reqChan <- &pr
			}
		}
	}
	return nil
}

func (rs *ResolverService) ResolveIPv6Addresses(ctx context.Context, domainName string) ([]net.IPAddr, error) {
	log.Info("resolving IPv6 addresses for ", domainName)
	if ips, found := rs.cache[domainName]; found {
		if time.Now().Before(rs.cacheTTL[domainName]) {
			return ips, nil
		}
	}

	ips, err := rs.dNSResolver.LookupIPAddr(ctx, domainName)
	if err != nil {
		return nil, err
	}
	log.Info("resolved IPv6 addresses for ", domainName, ": ", ips)

	var ipv6Addresses []net.IPAddr
	for _, ip := range ips {
		if ip.IP.To4() == nil {
			ipv6Addresses = append(ipv6Addresses, ip)
		}
	}

	ttl := time.Now().Add(1 * time.Minute)
	rs.cache[domainName] = ipv6Addresses
	rs.cacheTTL[domainName] = ttl

	return ipv6Addresses, nil
}
