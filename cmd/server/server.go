package main

import (
	"encoding/json"
	"log"
	"sync"

	pb "google.golang.org/grpc/examples/route_guide/routeguide"
)

type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	routeNotes map[string][]*pb.RouteNote
}

func newServer() *routeGuideServer {
	s := &routeGuideServer{routeNotes: make(map[string][]*pb.RouteNote)}
	if err := json.Unmarshal(exampleData, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
	return s
}

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		if err := stream.Send(feature); err != nil {
			return err
		}
	}
	return nil
}

var exampleData = []byte(`[{
	"location": {
			"latitude": 407838351,
			"longitude": -746143763
	},
	"name": "Patriots Path, Mendham, NJ 07945, USA"
}, {
	"location": {
			"latitude": 408122808,
			"longitude": -743999179
	},
	"name": "101 New Jersey 10, Whippany, NJ 07981, USA"
}, {
	"location": {
			"latitude": 413628156,
			"longitude": -749015468
	},
	"name": "U.S. 6, Shohola, PA 18458, USA"
}, {
	"location": {
			"latitude": 419999544,
			"longitude": -740371136
	},
	"name": "5 Conners Road, Kingston, NY 12401, USA"
}, {
	"location": {
			"latitude": 414008389,
			"longitude": -743951297
	},
	"name": "Mid Hudson Psychiatric Center, New Hampton, NY 10958, USA"
}, {
	"location": {
			"latitude": 419611318,
			"longitude": -746524769
	},
	"name": "287 Flugertown Road, Livingston Manor, NY 12758, USA"
}]`)
