package hello

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schorlet/gcp/world/api/grpc/interceptor"
	"github.com/schorlet/gcp/world/api/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type Client struct {
	pb.HelloClient
	io.Closer
}

func NewClient(addr string) (*Client, error) {
	opts := []grpc.DialOption{
		grpc.WithReturnConnectionError(), // Implies WithBlock()
		grpc.WithUnaryInterceptor(interceptor.UnaryClientLogger),
		// grpc.WithStatsHandler(h stats.Handler),
	}

	caFile := os.Getenv("HELLO_SERVER_CA_FILE")
	if caFile != "" {
		creds, err := credentials.NewClientTLSFromFile(caFile, "")
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	ctx, timeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer timeout()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("dialing server: %v\n", err)
	}

	return &Client{
		HelloClient: pb.NewHelloClient(conn),
		Closer:      conn,
	}, nil
}

func HandleWithClient(client pb.HelloClient) (http.HandlerFunc, error) {
	grpcmux := runtime.NewServeMux()
	err := pb.RegisterHelloHandlerClient(context.Background(), grpcmux, client)
	if err != nil {
		return nil, err
	}
	return grpcmux.ServeHTTP, err
}
