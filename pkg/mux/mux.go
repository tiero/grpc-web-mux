// Package mux helps multiplexing gRPC and gRPC Web on the same port, switching on HTTP Content-Type Header.
// It features insecure clear-text, TLS termination and Onion service
package mux

import (
	"net"
	"net/http"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// Mux holds a net.Listener and a *grpc.Server
type Mux struct {
	Listener   net.Listener
	GrpcServer *grpc.Server
}

// Serve starts multiplexing gRPC and gRPC Web on the same port. Serve blocks and perhaps should be invoked concurrently within a go routine.
func (m *Mux) Serve() error {
	mux := cmux.New(m.Listener)
	grpcL := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := mux.Match(cmux.HTTP1Fast())

	grpcWebServer := grpcweb.WrapServer(
		m.GrpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	go m.GrpcServer.Serve(grpcL)
	go http.Serve(httpL, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if isValidRequest(req) {
			grpcWebServer.ServeHTTP(resp, req)
		}
	}))

	return mux.Serve()
}

func isValidRequest(req *http.Request) bool {
	return isValidGrpcWebOptionRequest(req) || isValidGrpcWebRequest(req)
}

func isValidGrpcWebRequest(req *http.Request) bool {
	return req.Method == http.MethodPost && isValidGrpcContentTypeHeader(req.Header.Get("content-type"))
}

func isValidGrpcContentTypeHeader(contentType string) bool {
	return strings.HasPrefix(contentType, "application/grpc-web-text") ||
		strings.HasPrefix(contentType, "application/grpc-web")
}

func isValidGrpcWebOptionRequest(req *http.Request) bool {
	accessControlHeader := req.Header.Get("Access-Control-Request-Headers")
	return req.Method == http.MethodOptions &&
		strings.Contains(accessControlHeader, "x-grpc-web") &&
		strings.Contains(accessControlHeader, "content-type")
}
