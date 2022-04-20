package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	contentType     = "content-Type"
	grpcContentType = "application/grpc"
)

// Dispatcher for routing to REST or gRPC handler depending upon HTTP1.1 or HTTP2
func grpcDispatcher(ctx context.Context, grpcHander http.Handler, httpHandler http.Handler) http.Handler {
	hf := func(w http.ResponseWriter, r *http.Request) {
		req := r.WithContext(ctx)

		contentTypeHeader := r.Header.Get(contentType)

		if r.ProtoMajor == 2 && strings.HasPrefix(contentTypeHeader, grpcContentType) {
			log.Printf("%s \"%s\", routing to gRPC server \n", contentType, contentTypeHeader)
			grpcHander.ServeHTTP(w, req)
		} else {
			log.Printf("%s \"%s\", routing to HTTP server\n", contentType, contentTypeHeader)
			httpHandler.ServeHTTP(w, req)
		}

	}
	return h2c.NewHandler(http.HandlerFunc(hf), &http2.Server{})
}
