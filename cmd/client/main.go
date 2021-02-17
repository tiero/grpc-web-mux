package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"golang.org/x/net/proxy"
)

const (
	defaultAddress = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	address := defaultAddress
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(15*time.Second), grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		return dialer.Dial("tcp", addr)
	}))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := "world"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.GetMessage())
}
