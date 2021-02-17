// Package mux helps multiplexing gRPC and gRPC Web on the same port, switching on HTTP Content-Type Header.
// It features insecure clear-text, TLS termination and Onion service
package mux

import (
	"net"
	"net/http"
	"strings"

	"github.com/cretz/bine/tor"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"

	"google.golang.org/grpc"
)

// Mux holds a net.Listener, a *grpc.Server and if onion service a *tor.Tor client
type Mux struct {
	mux        cmux.CMux
	Listener   net.Listener
	GrpcServer *grpc.Server

	torClient *tor.Tor

	extraServers []*HTTPServer
}

// HTTPServer holds a net.Listener and an http.Handler
type HTTPServer struct {
	Listener net.Listener
	Handler  http.Handler
}

// Serve starts multiplexing gRPC and gRPC Web on the same port. Serve blocks and perhaps should be invoked concurrently within a go routine.
func (m *Mux) Serve() error {
	grpcL := m.mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := m.mux.Match(cmux.HTTP1Fast())

	grpcWebServer := grpcweb.WrapServer(
		m.GrpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	go m.GrpcServer.Serve(grpcL)
	go http.Serve(httpL, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if isValidRequest(req) {
			grpcWebServer.ServeHTTP(resp, req)
		} else {
			http.NotFound(resp, req)
		}
	}))

	if len(m.extraServers) > 0 {
		for _, s := range m.extraServers {
			go http.Serve(s.Listener, s.Handler)
		}
	}

	return m.mux.Serve()
}

// WithExtraHTTP1 adds to the Mux an additional given HTTP1 handler and optional content-type.
// if contentTypes is nil or empty array, any HTTP1 request will be matched
// any of the contentTypes sourced MUST be different than `application/grpc-web-text` and `application/grpc-web`
// because it will have matched before in the grpc-web server
func (m *Mux) WithExtraHTTP1(handler http.Handler, contentTypes []string) {
	if contentTypes == nil || len(contentTypes) == 0 {
		lis := m.mux.Match(cmux.HTTP1())

		m.extraServers = append(m.extraServers, &HTTPServer{
			Listener: lis,
			Handler:  handler,
		})
		return
	}

	for _, ct := range contentTypes {
		lis := m.mux.Match(cmux.HTTP1HeaderFieldPrefix("content-type", ct))

		m.extraServers = append(m.extraServers, &HTTPServer{
			Listener: lis,
			Handler:  handler,
		})
	}

}

// Close closes the TCP listener connection and in case of onion service it will also halt the tor client.
func (m *Mux) Close() {
	if m.torClient != nil {
		m.torClient.Close()
	}
	m.Listener.Close()
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
