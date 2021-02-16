package mux_test

import (
	"log"

	mux "github.com/tiero/grpc-web-mux"
	"google.golang.org/grpc"
)

func ExampleNewMuxWithInsecure() {

	myGrpcServer := grpc.NewServer()

	insecureMux, err := mux.NewMuxWithInsecure(
		myGrpcServer,
		mux.InsecureOptions{Address: ":8080"},
	)
	if err != nil {
		log.Panic(err)
	}

	insecureMux.Serve()
}

func ExampleNewMuxWithTLS() {

	myGrpcServer := grpc.NewServer()

	tlsMux, err := mux.NewMuxWithTLS(
		myGrpcServer,
		mux.TLSOptions{
			Address: ":9945",
			Domain:  "mydomain.com",
		},
	)
	if err != nil {
		log.Panic(err)
	}

	tlsMux.Serve()
}

func ExampleNewMuxWithOnion() {

	myGrpcServer := grpc.NewServer()

	onionMux, err := mux.NewMuxWithOnion(
		myGrpcServer,
		mux.OnionOptions{
			Port:       80,
			PrivateKey: "myBase64SerializedPrivateKey",
		},
	)
	if err != nil {
		log.Panic(err)
	}

	onionMux.Serve()
}
