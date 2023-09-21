package server

import (
	"fmt"
	"log"
	"sync"

	"github.com/hawkv6/hawkwing/pkg/bpf"
	"github.com/hawkv6/hawkwing/pkg/bpf/server"
	"github.com/hawkv6/hawkwing/pkg/linker"
	"github.com/vishvananda/netlink"
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
	xdpObjs, err := server.ReadServerXdpObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load TC program: %s", err)
	}
	xdpLinker := linker.NewXdpLinker(iface, xdpObjs.FilterIngress)

	tcObjs, err := server.ReadServerTcObjects()
	if err != nil {
		return nil, fmt.Errorf("could not load TC program: %s", err)
	}
	tcLinker := linker.NewTcLinker(iface, tcObjs.FilterIngress, "ingress")

	err = bpf.Mount()
	if err != nil {
		log.Fatalf("Could not mount BPF filesystem: %s", err)
	}

	return &Server{
		iface:     iface,
		xdpLinker: xdpLinker,
		tcLinker:  tcLinker,
		wg:        &sync.WaitGroup{},
	}, nil
}

func (s *Server) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := s.xdpLinker.Attach()
		// err := s.tcLinker.Attach()
		if err != nil {
			log.Fatalf("Could not attach TC program: %s", err)
		}
	}()

	log.Printf("Server started on interface %s", s.iface.Attrs().Name)
	log.Printf("Press Ctrl+C to exit and remove the program")
}

func (s *Server) Stop() {
	// if err := s.tcLinker.Detach(); err != nil {
	if err := s.xdpLinker.Detach(); err != nil {
		log.Fatalf("Could not detach TC program: %s", err)
	}
	s.wg.Wait()
}
