package mux

import (
	"net"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// InsecureOptions holds the address where to listen for TCP packets
type InsecureOptions struct {
	Address string
}

// NewMuxWithInsecure returns a clear-text *GrpcWebMux. Use only for development.
func NewMuxWithInsecure(grpcServer *grpc.Server, opts InsecureOptions) (*GrpcWebMux, error) {

	lis, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return nil, err
	}

	mux := cmux.New(lis)
	return &GrpcWebMux{mux: mux, Listener: lis, GrpcServer: grpcServer}, nil
}
