package linker

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"

	"github.com/cilium/ebpf/link"
	"github.com/vishvananda/netlink"
)

func directionToParentDisc(direction string) uint32 {
	switch direction {
	case "ingress":
		return netlink.HANDLE_MIN_INGRESS
	case "egress":
		return netlink.HANDLE_MIN_EGRESS
	}
	return 0
}

func (l *XdpLinker) detachXdp() error {
	return l.link.Close()
}

func (l *XdpLinker) attachXdp() (link.Link, error) {
	if l.program == nil {
		return nil, errors.New("cannot attach a nil program")
	}
	link, err := link.AttachXDP(link.XDPOptions{
		Program:   l.program,
		Interface: l.iface.Attrs().Index,
		Flags:     link.XDPGenericMode,
	})
	if err != nil {
		return nil, fmt.Errorf("attaching XDP program to interface %s: %v", l.iface.Attrs().Name, err)
	}
	return link, nil
}

func (l *TcLinker) replaceQdisc() error {
	attrs := netlink.QdiscAttrs{
		LinkIndex: l.iface.Attrs().Index,
		Handle:    netlink.MakeHandle(0xffff, 0),
		Parent:    netlink.HANDLE_CLSACT,
	}

	qdisc := &netlink.GenericQdisc{
		QdiscAttrs: attrs,
		QdiscType:  "clsact",
	}

	return netlink.QdiscReplace(qdisc)
}

func (l *TcLinker) removeTCFilters() error {
	filters, err := netlink.FilterList(l.iface, l.direction)
	if err != nil {
		return err
	}

	for _, f := range filters {
		if err := netlink.FilterDel(f); err != nil {
			return err
		}
	}

	return nil
}

func (l *TcLinker) attachTCProgram() error {
	if l.program == nil {
		return errors.New("cannot attach a nil program")
	}

	if err := l.replaceQdisc(); err != nil {
		return fmt.Errorf("replacing qdisc on interface %s: %v", l.iface.Attrs().Name, err)
	}

	filter := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: l.iface.Attrs().Index,
			Parent:    l.direction,
			Handle:    1,
			Protocol:  unix.ETH_P_ALL,
			Priority:  1,
		},
		Fd:           l.program.FD(),
		Name:         l.programName,
		DirectAction: true,
	}
	if err := netlink.FilterReplace(filter); err != nil {
		return fmt.Errorf("attaching TC program to interface %s: %v", l.iface.Attrs().Name, err)
	}
	return nil
}
