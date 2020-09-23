package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	hello_adapter "github.com/schorlet/gcp/world/webapp/adapters/grpc/hello"
	http_ports "github.com/schorlet/gcp/world/webapp/ports/http"
)

var (
	Version = "unset"
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

	var httpServer *http.Server

	// cancellation handler
	g.Go(func() error {
		<-ctx.Done()

		if httpServer != nil {
			ctx, timeout := context.WithTimeout(context.Background(), 1*time.Second)
			defer timeout()
			if err := httpServer.Shutdown(ctx); err != nil {
				return err
			}
		}

		return ctx.Err()
	})

	// http server
	g.Go(func() error {
		helloClient, err := hello_adapter.NewClient()
		if err != nil {
			return err
		}
		defer helloClient.Close()

		helloHandler, err := hello_adapter.HandleWithClient(helloClient)
		if err != nil {
			return err
		}

		router, err := http_ports.NewRouter(Version, helloHandler)
		if err != nil {
			return err
		}

		httpServer, err = http_ports.NewServer(router)
		if err != nil {
			return err
		}

		return httpServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Printf("Error: %v", err)
	}
}
