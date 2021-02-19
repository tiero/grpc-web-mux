package main

import (
	"context"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"github.com/tiero/grpc-web-mux/pkg/mux"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {

	myGrpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(myGrpcServer, &server{})

	options := mux.OnionOptions{Port: 80}
	if len(os.Args) > 1 {
		options = mux.OnionOptions{
			Port:       80,
			PrivateKey: os.Args[1],
		}
	}
	serverMux, err := mux.NewMuxWithOnion(
		myGrpcServer,
		options,
	)
	if err != nil {
		log.Panic(err)
	}

	serverMux.WithHTTP1Handler(
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
