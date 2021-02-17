package mux

import (
	"crypto/tls"

	"github.com/caddyserver/certmagic"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// TLSOptions holds the address where to listen for TCP packets and
// defines the domain for which we need to obtain and renew a TLS cerficate
type TLSOptions struct {
	Address string
	Domain  string
}

// NewMuxWithTLS returns a *Mux with TLS termination automatically obtaining the TLS certificate with given domain through CertMagic.
// By default, CertMagic stores assets on the local file system in $HOME/.local/share/certmagic (and honors $XDG_DATA_HOME if set).
// CertMagic will create the directory if it does not exist.
//If writes are denied, things will not be happy, so make sure CertMagic can write to it!
func NewMuxWithTLS(grpcServer *grpc.Server, opts TLSOptions) (*Mux, error) {
	tlsConfig, err := certmagic.TLS([]string{opts.Domain})
	if err != nil {
		return nil, err
	}

	const requiredCipher = tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
	tlsConfig.CipherSuites = []uint16{requiredCipher}
	tlsConfig.NextProtos = []string{"http/1.1", "h2", "h2-14"} // h2-14 is just for compatibility. will be eventually removed.

	lis, err := tls.Listen("tcp", opts.Address, tlsConfig)
	if err != nil {
		return nil, err
	}

	mux := cmux.New(lis)
	return &Mux{mux: mux, Listener: lis, GrpcServer: grpcServer}, nil
}
