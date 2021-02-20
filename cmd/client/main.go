package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"

	"golang.org/x/net/proxy"
)

const (
	defaultAddress = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	var dialOpts []grpc.DialOption = []grpc.DialOption{grpc.WithInsecure()}
	var address string = defaultAddress

	if len(os.Args) > 1 {
		address = os.Args[1]
		dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9150", nil, nil)
		if err != nil {
			log.Fatal(err)
		}

		dialOpts = append(
			dialOpts,
			grpc.WithTimeout(15*time.Second),
			grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) { return dialer.Dial("tcp", addr) }),
		)
	}

	conn, err := grpc.Dial(address, dialOpts...)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRouteGuideClient(conn)

	rect := &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	}

	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", c, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", c, err)
		}
		log.Printf("Feature: name: %q, point:(%v, %v)", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}
