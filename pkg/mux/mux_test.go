package mux_test

import (
	"log"
	"net/http"

	"github.com/tiero/grpc-web-mux/pkg/mux"
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

	log.Printf("Serving mux at %s\n", insecureMux.Listener.Addr().String())

	defer insecureMux.Close()
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

	log.Printf("Serving mux at %s\n", tlsMux.Listener.Addr().String())

	defer tlsMux.Close()
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

	log.Printf("Serving mux at %s\n", onionMux.Listener.Addr().String())

	defer onionMux.Close()
	onionMux.Serve()
}

func ExampleMux_WithHTTP1Handler() {

	myGrpcServer := grpc.NewServer()

	insecureMux, err := mux.NewMuxWithInsecure(
		myGrpcServer,
		mux.InsecureOptions{Address: ":8080"},
	)
	if err != nil {
		log.Panic(err)
	}

	insecureMux.WithHTTP1Handler(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hello Insecure!"))
		}),
		nil,
	)

	log.Printf("Serving mux at %s\n", insecureMux.Listener.Addr().String())

	defer insecureMux.Close()
	insecureMux.Serve()

}
