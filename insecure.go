package mux

import (
	"net"

	"google.golang.org/grpc"
)

// InsecureOptions holds the address where to listen for TCP packets
type InsecureOptions struct {
	Address string
}

// NewMuxWithInsecure returns a clear-text *Mux. Use only for development.
func NewMuxWithInsecure(grpcServer *grpc.Server, opts InsecureOptions) (*Mux, error) {

	lis, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return nil, err
	}

	return &Mux{Listener: lis, GrpcServer: grpcServer}, nil
}
