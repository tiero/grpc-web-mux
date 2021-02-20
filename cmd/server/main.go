package main

import (
	"flag"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"github.com/tiero/grpc-web-mux/pkg/mux"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
)

var useInsecure = flag.Bool("insecure", false, "if running insecure mux")

func main() {
	flag.Parse()

	myGrpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(myGrpcServer, newServer())

	var err error
	var serverMux *mux.GrpcWebMux
	if *useInsecure {
		serverMux, err = mux.NewMuxWithInsecure(myGrpcServer, mux.InsecureOptions{Address: ":8080"})
	} else {
		options := mux.OnionOptions{Port: 80}

		// If an argument is given means we got a private key
		if len(os.Args) > 1 {
			options = mux.OnionOptions{
				Port:       80,
				PrivateKey: os.Args[1],
			}
		}

		serverMux, err = mux.NewMuxWithOnion(
			myGrpcServer,
			options,
		)
	}

	if err != nil {
		log.Panic(err)
	}

	serverMux.WithExtraHandler(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hello my friend!"))
		}),
		[]string{"application/json"},
	)

	log.Printf("Serving mux at %s\n", serverMux.Listener.Addr().String())

	defer serverMux.Close()
	serverMux.Serve()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Println("shutting down mux")

}
