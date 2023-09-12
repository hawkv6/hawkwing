package bpf

type ClientInnerMap struct {
	dstPort  uint16
	segments [10]struct{ in6U struct{ u6Addr8 [16]uint8 } }
}

type ClientOuterMap struct {
	domainName  [256]byte
	innerMapMap map[uint32]ClientInnerMap
}

// func NewClientInnerMap(dstPort uint16, segments []string) *ClientInnerMap {
// 	return &ClientInnerMap{
// 		dstPort:  dstPort,
// 		segments: segments,
// 	}
// }
