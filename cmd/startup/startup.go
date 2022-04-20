package startup

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhishekmandhare/go-grpc-rest/internal/api"
	"github.com/abhishekmandhare/go-grpc-rest/internal/server"
	orders "github.com/abhishekmandhare/go-grpc-rest/orders-api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunAppServer(ctx context.Context) func() error {

	return func() error {
		port := 8080
		restServer, err := server.NewRESTServer(ctx, port)
		if err != nil {
			log.Printf("Error creating listener for REST server : %v", err)
			return err
		}

		orders.RegisterOrdersAPIServer(restServer.GRPCServer, api.NewOrdersAPI())

		if err := orders.RegisterOrdersAPIHandlerFromEndpoint(ctx, restServer.GRPCGatewayMux, fmt.Sprintf(":%d", port), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
			return err
		}

		errChan := make(chan error)

		go func(s *server.RESTServer, errChan chan error) {
			if err := s.Serve(); err != nil {
				errChan <- err
			}
		}(restServer, errChan)

		select {
		case <-ctx.Done():
			log.Println("REST server terminated by upstream")
			restServer.GracefulStop(ctx)
		case err := <-errChan:
			log.Printf("REST server terminated by unexpected error : %s \n", err)
			return err
		}

		return nil
	}
}

// RunSignalListener returns a function that starts a listener for system signals.
func RunSignalListener(ctx context.Context) func() error {
	return func() error {
		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)

		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-sigChan:
			return fmt.Errorf("Terminated by SIGTERM")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
