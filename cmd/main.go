package main

import (
	"context"
	"log"

	"github.com/abhishekmandhare/go-grpc-rest/cmd/startup"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx := context.Background()
	errGrp, gCtx := errgroup.WithContext(ctx)
	errGrp.Go(startup.RunAppServer(gCtx))
	errGrp.Go(startup.RunSignalListener(ctx))
	if err := errGrp.Wait(); err != nil {
		log.Fatalf("Terminating with error %s\n", err)
	}
}
