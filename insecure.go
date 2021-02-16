package mux

import (
	"net"

	"google.golang.org/grpc"
)

// InsecureOptions ...
type InsecureOptions struct {
	Address string
}

// NexMuxWithInsecure returns a clear-text *Mux. Use only for development.
func NexMuxWithInsecure(grpcServer *grpc.Server, opts InsecureOptions) (*Mux, error) {

	lis, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return nil, err
	}

	return &Mux{listener: lis, grpcServer: grpcServer}, nil
}
