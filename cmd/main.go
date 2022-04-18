package main

import (
	orders "cmd/main.go/orders-api/proto"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server")

	grpcServer := grpc.NewServer()

	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server", err)
	}
	router := runtime.NewServeMux()

	if err = orders.RegisterOrdersAPIHandler(context.Background(), router, conn); err != nil {
		log.Fatalln("Failed to regidter gateway:", err)
	}

	http.ListenAndServe(":8080", httpGrpcRouter(grpcServer, router))
}

func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}
