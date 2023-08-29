package bpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -no-global-types -cc $BPF_CLANG -cflags $BPF_CFLAGS xdp ../../src/main.c -- -I../../src

type IntentServiceData struct {
	DomainName [256]byte
	SrcAddrV6  struct{ In6U struct{ U6Addr8 [16]uint8 } }
	DstAddrV6  struct{ In6U struct{ U6Addr8 [16]uint8 } }
	DstPort    uint16
	Segments   [10]struct{ In6U struct{ U6Addr8 [16]uint8 } }
}
