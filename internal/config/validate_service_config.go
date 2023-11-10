package config

import (
	"net"

	"github.com/go-playground/validator"
)

func ServiceConfigValidation(sl validator.StructLevel) {
	serviceCfg := sl.Current().Interface().(ServiceConfig)

	if len(serviceCfg.DomainName) == 0 && len(serviceCfg.Ipv6Addresses) == 0 {
		sl.ReportError(serviceCfg.DomainName, "domain_name", "", "domain_name or ipv6_addresses is required", "")
		sl.ReportError(serviceCfg.Ipv6Addresses, "ipv6_addresses", "", "domain_name or ipv6_addresses is required", "")
	}

	if serviceCfg.Ipv6Addresses != nil {
		for _, ipv6 := range serviceCfg.Ipv6Addresses {
			ip := net.ParseIP(ipv6)
			if ip == nil || ip.To4() != nil {
				sl.ReportError(ipv6, "ipv6_addresses", "", "ipv6_addresses must be valid ipv6 addresses", "")
			}
		}
	}

	if len(serviceCfg.Applications) == 0 {
		sl.ReportError(serviceCfg.Applications, "applications", "", "at least one application must be specified", "")
	}
}
