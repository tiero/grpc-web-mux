package mux

import (
	"context"
	"crypto"
	"fmt"
	"io"
	"time"

	"github.com/cretz/bine/control"
	"github.com/cretz/bine/tor"
	"github.com/ipsn/go-libtor"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// OnionOptions defines options of onion service.
// PrivateKey is a base64 blob, if not present, a ED25519 key is generated for OnionV3.
// DataDir is the directory used by Tor. If it is empty, a temporary
// directory is created in TempDataDirBase.
// TempDataDirBase is the parent directory that a temporary data directory
// will be created under for use by Tor. This is ignored if DataDir is not
// empty. If empty it is assumed to be the current working directory
type OnionOptions struct {
	// The port wich the onion service will listen on
	Port int
	//PrivateKey is a base64 blob. If not present, a key is generated based
	PrivateKey string
	// DataDir defines where to store temporary files for the Onion service. If empty will be created one in the path of the process
	DataDir string
	// DebugWriter holds the debug output. If nil, no ouputs
	DebugWriter io.Writer
}

// NewMuxWithOnion returns a *GrpcWebMux publishing an V3 Onion Service with the given private key
// and creating an empty datadir in the current working directory.
// This will take couple of minutes to spin up, so be patient.
func NewMuxWithOnion(grpcServer *grpc.Server, opts OnionOptions) (*GrpcWebMux, error) {
	// Starting tor please wait a bit...
	torClient, err := tor.Start(nil, &tor.StartConf{
		ProcessCreator: libtor.Creator,
		DataDir:        opts.DataDir,
		DebugWriter:    opts.DebugWriter,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to start tor: %v", err)
	}

	// Wait at most a few minutes to publish the service
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// PrivKey
	var privKey crypto.PrivateKey = nil
	if len(opts.PrivateKey) > 0 {
		privKeyInfo, err := control.ED25519KeyFromBlob(opts.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Failed to deserialize private key: %v", err)
		}
		privKey = privKeyInfo.PrivateKey()
	}
	// Create an onion service to listen on any port but show as 80
	onion, err := torClient.Listen(ctx, &tor.ListenConf{
		RemotePorts: []int{opts.Port},
		Version3:    true,
		Key:         privKey,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create onion service: %v", err)
	}

	mux := cmux.New(onion)
	return &GrpcWebMux{
		mux:        mux,
		Listener:   onion,
		GrpcServer: grpcServer,
		torClient:  torClient,
	}, nil
}
