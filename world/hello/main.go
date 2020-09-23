package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/schorlet/gcp/world/hello/app"
	"github.com/schorlet/gcp/world/hello/grpc"

	"golang.org/x/sync/errgroup"
)

var (
	Version     = "unset"
	DefaultPort = "8021"
)

func main() {
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	// base context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// errgroup context
	g, ctx := errgroup.WithContext(ctx)

	var grpcServer *grpc.Server

	// interruption handler
	g.Go(func() error {
		interrupted := make(chan os.Signal, 1)
		signal.Notify(interrupted,
			// os.Interrupt
			syscall.SIGINT, // kill -SIGINT $(pgrep wiki) or `Ctrl+c`
			syscall.SIGTERM,
		)

		select {
		case <-interrupted:
			cancel()
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// cancellation handler
	g.Go(func() error {
		<-ctx.Done()

		if grpcServer != nil {
			grpcServer.GracefulStop()
		}

		return ctx.Err()
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	helloService := app.NewHelloService(Version)

	// grpc server
	g.Go(func() error {
		var err error
		grpcServer, err = grpc.NewServer(":"+port, helloService)
		if err != nil {
			return err
		}

		return grpcServer.ListenAndServeTLS()
	})

	if err := g.Wait(); err != nil {
		log.Printf("Error: %v", err)
	}
}
