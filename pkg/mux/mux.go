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

// GrpcWebMux holds a net.Listener, a *grpc.Server and if onion service a *tor.Tor client
type GrpcWebMux struct {
	mux        cmux.CMux
	Listener   net.Listener
	GrpcServer *grpc.Server

	torClient *tor.Tor

	handlerWithContentTypes *HandlerWithContentTypes
}

// HandlerWithContentTypes ...
type HandlerWithContentTypes struct {
	Handler      http.Handler
	ContentTypes []string
}

// Serve starts multiplexing gRPC and gRPC Web on the same port. Serve blocks and perhaps should be invoked concurrently within a go routine.
func (m *GrpcWebMux) Serve() error {
	grpcL := m.mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := m.mux.Match(cmux.HTTP1())

	grpcWebServer := grpcweb.WrapServer(
		m.GrpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	go m.GrpcServer.Serve(grpcL)
	go http.Serve(httpL, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if isValidOptionOrGrpcWebRequest(req) {
			grpcWebServer.ServeHTTP(resp, req)
		} else if m.handlerWithContentTypes.hasValidContentType(req) {
			m.handlerWithContentTypes.Handler.ServeHTTP(resp, req)
		} else {
			http.NotFound(resp, req)
		}
	}))

	return m.mux.Serve()
}

// WithExtraHandler adds to the GrpcWebMux an additional given HTTP1 handler and optional array of content-type to match.
// if contentTypes is nil or empty array, any HTTP1 request will be matched any of the contentTypes sourced
// MUST be different than `application/grpc-web-text` and `application/grpc-web` because it will have matched before
//in the grpc-web server listener
func (m *GrpcWebMux) WithExtraHandler(handler http.Handler, contentTypes []string) {
	m.handlerWithContentTypes = &HandlerWithContentTypes{
		Handler:      handler,
		ContentTypes: contentTypes,
	}
}

// Close closes the TCP listener connection and in case of onion service it will also halt the tor client.
func (m *GrpcWebMux) Close() {
	if m.torClient != nil {
		m.torClient.Close()
	}
	m.Listener.Close()
}

func (h HandlerWithContentTypes) hasValidContentType(req *http.Request) (found bool) {
	if len(h.ContentTypes) == 0 || h.ContentTypes == nil {
		return true
	}

	for _, c := range h.ContentTypes {
		if strings.HasPrefix(req.Header.Get("content-type"), c) {
			found = true
			break
		}
	}

	return
}

func isValidOptionOrGrpcWebRequest(req *http.Request) bool {
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
