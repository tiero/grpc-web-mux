package main

import (
	"context"
	"log"
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	log.Println("shutting down mux")

}
