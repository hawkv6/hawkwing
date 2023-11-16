package server

import (
	"fmt"
	"sync"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/server"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"github.com/hawkv6/hawkwing/pkg/maps"
	"github.com/vishvananda/netlink"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const (
	Subsystem = "go-server"
)

type Server struct {
	iface     netlink.Link
	xdpLinker *linker.XdpLinker
	tcLinker  *linker.TcLinker
	wg        *sync.WaitGroup
}

func NewServer(interfaceName string) (*Server, error) {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("could not lookup network iface %q: %s", interfaceName, err)
	}

	serverObjs, err := server.ReadServerBpfObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load server bpf objects: %s", err)
	}

	xdpLinker := linker.NewXdpLinker(iface, serverObjs.ServerIngress)
	tcLinker := linker.NewTcLinker(iface, serverObjs.ServerEgress, "egress")

	err = bpf.Mount()
	if err != nil {
		return nil, fmt.Errorf("could not mount BPF filesystem: %s", err)
	}

	realBpf := &bpf.RealBpf{}
	serverMap, err := maps.NewServerMap(realBpf)
	if err != nil {
		return nil, fmt.Errorf("could not create server map: %s", err)
	}

	err = serverMap.Create()
	if err != nil {
		return nil, fmt.Errorf("could not create server map: %s", err)
	}

	return &Server{
		iface:     iface,
		xdpLinker: xdpLinker,
		tcLinker:  tcLinker,
		wg:        &sync.WaitGroup{},
	}, nil
}

func (s *Server) Start() {
	s.wg.Add(2)
	go func() {
		defer s.wg.Done()
		err := s.xdpLinker.Attach()
		if err != nil {
			log.Fatalf("could not attach server XDP program: %s", err)
		}
	}()

	go func() {
		defer s.wg.Done()
		err := s.tcLinker.Attach()
		if err != nil {
			log.Fatalf("could not attach server TC program: %s", err)
		}
	}()

	fmt.Printf("Server started on interface %s\n", s.iface.Attrs().Name)
	fmt.Println("Press Ctrl+C to exit and remove the program")
}

func (s *Server) Stop() {
	if err := s.xdpLinker.Detach(); err != nil {
		log.Fatalf("could not detach server XDP program: %s", err)
	}
	if err := s.tcLinker.Detach(); err != nil {
		log.Fatalf("could not detach server TC program: %s", err)
	}
	s.wg.Wait()
}
