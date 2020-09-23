package grpc

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	hello_client "github.com/schorlet/gcp/world/api/grpc/hello"
	"github.com/schorlet/gcp/world/api/pb"
	"github.com/schorlet/gcp/world/hello/app"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	grpcAddr = "localhost:7999"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")
}

func startServer(t *testing.T) {
	helloService := app.NewHelloService("testing")
	grpcServer, err := NewServer(grpcAddr, helloService)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := grpcServer.ListenAndServeTLS()
		if err != nil {
			panic(err)
		}
	}()

	t.Cleanup(func() {
		grpcServer.GracefulStop()
	})
}

func helloClient(t *testing.T) pb.HelloClient {
	client, err := hello_client.NewClient(grpcAddr)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		client.Close()
	})
	return client
}

func TestHello(t *testing.T) {
	startServer(t)
	client := helloClient(t)

	md := metadata.MD{}
	md["foo"] = []string{"bar"}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, md)

	mdh := metadata.MD{}
	res, err := client.Hello(ctx, &pb.HelloRequest{Name: "world"},
		grpc.WaitForReady(true), grpc.Header(&mdh),
	)
	if err != nil {
		t.Fatalf("Error: hello: %v\n", err)
	}

	log.Printf("Metadata received: %v\n", mdh)

	if res.Message != "Hello world" {
		t.Errorf("Bad response: %s, want: %s\n", res.Message, "hello world")
	}
	if res.Version != "testing" {
		t.Errorf("Bad response: %s, want: %s\n", res.Version, "testing")
	}
}
