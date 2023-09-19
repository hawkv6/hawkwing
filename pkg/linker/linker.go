package linker

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/hawkv6/hawkwing/pkg/logging"
	"github.com/vishvananda/netlink"
)

var log = logging.DefaultLogger.WithField("subsystem", Subsystem)

const (
	Subsystem = "ebpf-linker"
)

type Linker interface {
	Attach() error
	Detach() error
}

type XdpLinker struct {
	iface   netlink.Link
	program *ebpf.Program
	link    link.Link
}

type TcLinker struct {
	iface       netlink.Link
	program     *ebpf.Program
	programName string
	direction   uint32
}

func NewXdpLinker(iface netlink.Link, program *ebpf.Program) *XdpLinker {
	return &XdpLinker{
		iface:   iface,
		program: program,
	}
}

func (l *XdpLinker) Attach() error {
	link, err := l.attachXdp()
	if err != nil {
		log.WithError(err).Errorf("couldn't attach XDP program to interface %s", l.iface.Attrs().Name)
	}
	l.link = link
	return nil
}

func (l *XdpLinker) Detach() error {
	err := l.detachXdp()
	if err != nil {
		log.WithError(err).Errorf("couldn't dettach XDP program from interface %s", l.iface.Attrs().Name)
	}
	return nil
}

func NewTcLinker(iface netlink.Link, program *ebpf.Program, direction string) *TcLinker {
	return &TcLinker{
		iface:     iface,
		program:   program,
		direction: directionToParentDisc(direction),
	}
}

func (l *TcLinker) Attach() error {
	err := l.attachTCProgram()
	if err != nil {
		log.WithError(err).Errorf("couldn't attach TC program to interface %s", l.iface.Attrs().Name)
	}
	return nil
}

func (l *TcLinker) Detach() error {
	err := l.removeTCFilters()
	if err != nil {
		log.WithError(err).Errorf("couldn't remove TC filters from interface %s", l.iface.Attrs().Name)
	}
	return nil
}
