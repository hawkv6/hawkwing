// Package client provides the client specific functions.
//
// This package is responsible for handling the client side of the BPF
// implementation. It attaches the client eBPF programs to the interfaces
// on the host. It also creates the client side maps.
//
// This package is also responsible to start the gRPC client and connect
// to the controller. It manages the different goroutines and channels.
package client
