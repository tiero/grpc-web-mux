package main

import (
	"context"
	"log"

	"github.com/tiero/grpc-web-mux/pkg/mux"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
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
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})

	myGrpcServer := grpc.NewServer()

	insecureMux, err := mux.NewMuxWithInsecure(
		myGrpcServer,
		mux.InsecureOptions{Address: ":8080"},
	)
	if err != nil {
		log.Panic(err)
	}

	insecureMux.Serve()
	log.Println(insecureMux.Listener.Addr().String())
}
