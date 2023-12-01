package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/netip"

	"github.com/phuslu/fastdns"
	"github.com/spf13/viper"
)

var (
	viperInstance = viper.NewWithOptions(viper.KeyDelimiter("\\"))
	Params        Config
	CfgFile       string
)

type Host struct {
	DomainName    string   `mapstructure:"domain_name"`
	Ipv6Addresses []string `mapstructure:"ipv6_addresses"`
}

type Config struct {
	Debugging bool   `mapstructure:"debugging"`
	Hosts     []Host `mapstructure:"hosts"`
}

func ParseConfig() error {
	if len(viperInstance.ConfigFileUsed()) != 0 {
		if err := viperInstance.ReadInConfig(); err != nil {
			return fmt.Errorf("failed to load config file %s: %v", viperInstance.ConfigFileUsed(), err)
		}
	}

	if err := viperInstance.Unmarshal(&Params); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}

func GetConfigInstance() *viper.Viper {
	return viperInstance
}

type DNSHandler struct {
	Debug bool
}

func (h *DNSHandler) ServeDNS(rw fastdns.ResponseWriter, req *fastdns.Message) {
	if h.Debug {
		log.Printf("%s: CLASS %s TYPE %s\n", req.Domain, req.Question.Class, req.Question.Type)
	}

	for _, host := range Params.Hosts {
		if bytes.Equal(req.Domain, []byte(host.DomainName)) {
			switch req.Question.Type {
			case fastdns.TypeA:
				for _, addr := range host.Ipv6Addresses {
					fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr(addr)})
				}
			case fastdns.TypeAAAA:
				for _, addr := range host.Ipv6Addresses {
					fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr(addr)})
				}
			default:
				fastdns.Error(rw, req, fastdns.RcodeNXDomain)
			}
		}
	}
}

func Start() {
	addr := ":53"

	server := &fastdns.ForkServer{
		Handler: &DNSHandler{
			Debug: Params.Debugging,
		},
		Stats: &fastdns.CoreStats{
			Prefix: "coredns_",
			Family: "1",
			Proto:  "udp",
			Server: "dns://" + addr,
			Zone:   ".",
		},
		ErrorLog: log.Default(),
	}

	err := server.ListenAndServe(addr)
	if err != nil {
		log.Fatalf("dnsserver error: %+v", err)
	}
}

func init() {
	flag.StringVar(&CfgFile, "config", "./config.yaml", "config file (default is ./config.yaml)")
}

func main() {
	flag.Parse()
	viperInstance.SetConfigFile(CfgFile)
	err := ParseConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	Start()
}
