package main

import (
	"bytes"
	"log"
	"net/netip"

	"github.com/phuslu/fastdns"
)

type DNSHandler struct {
	Debug bool
}

func (h *DNSHandler) ServeDNS(rw fastdns.ResponseWriter, req *fastdns.Message) {
	if h.Debug {
		log.Printf("%s: CLASS %s TYPE %s\n", req.Domain, req.Question.Class, req.Question.Type)
	}

	hostB := []byte("wb.hawk.net")
	hostC := []byte("wc.hawk.net")
	if bytes.Equal(req.Domain, hostB) {
		switch req.Question.Type {
		// case fastdns.TypeA:
		// 	fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr("fcbb:cc00:2::a")})
		case fastdns.TypeAAAA:
			fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr("fcbb:cc00:2::a")})
		default:
			fastdns.Error(rw, req, fastdns.RcodeNXDomain)
		}
	} else if bytes.Equal(req.Domain, hostC) {
		switch req.Question.Type {
		// case fastdns.TypeA:
		// 	fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr("fcbb:cc00:3::a")})
		case fastdns.TypeAAAA:
			fastdns.HOST(rw, req, 60, []netip.Addr{netip.MustParseAddr("fcbb:cc00:3::a")})
		default:
			fastdns.Error(rw, req, fastdns.RcodeNXDomain)
		}
	}
}

func main() {
	addr := ":53"

	server := &fastdns.ForkServer{
		Handler: &DNSHandler{
			// Debug: os.Getenv("DEBUG") != "",
			Debug: true,
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
