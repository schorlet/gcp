package hello

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/schorlet/gcp/world/api/pb"
	"github.com/schorlet/gcp/world/hello/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var _ pb.HelloServer = (*Server)(nil)

type Server struct {
	pb.UnimplementedHelloServer
	service domain.HelloService
}

func NewServer(service domain.HelloService) *Server {
	return &Server{service: service}
}

func (server *Server) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// simulate latency
	random := rand.New(rand.NewSource(time.Now().Unix()))
	time.Sleep(time.Duration(random.Intn(1000)) * time.Millisecond)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("Metadata received: %v\n", md)
	}

	if p, ok := peer.FromContext(ctx); ok {
		if err := grpc.SendHeader(ctx, metadata.MD{
			"peer-address": []string{p.Addr.String()},
		}); err != nil {
			return nil, err
		}
	}

	greeting, err := server.service.Hello(req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.HelloResponse{
		Message:  greeting.Message,
		Version:  greeting.Version,
		Hostname: greeting.Hostname,
	}, nil
}
