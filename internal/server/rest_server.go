package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type RESTServer struct {
	server         *http.Server
	listener       net.Listener
	GRPCGatewayMux *runtime.ServeMux
	GRPCServer     *grpc.Server
}

func NewRESTServer(ctx context.Context, port int) (*RESTServer, error) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	httpGatewayMux := http.NewServeMux()
	grpcGatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)

	httpGatewayMux.Handle("/", grpcGatewayMux)
	server := &http.Server{Handler: grpcDispatcher(ctx, grpcServer, httpGatewayMux)}

	log.Printf("REST server is listening on %d \n", port)
	return &RESTServer{listener: listener, server: server, GRPCGatewayMux: grpcGatewayMux, GRPCServer: grpcServer}, nil
}

//GracefulStop tries to gracefully stop the REST server.
func (s RESTServer) GracefulStop(ctx context.Context) {
	log.Println("Gracefully shutting down REST server")
	_ = s.server.Shutdown(ctx)
	log.Println("Gracefully shutting down gRPC server")
	s.GRPCServer.GracefulStop()
}

// Serve starts the REST server
func (s RESTServer) Serve() error {
	return s.server.Serve(s.listener)
}
