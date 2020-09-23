package grpc

import (
	"net"
	"os"

	"github.com/schorlet/gcp/world/api/grpc/interceptor"
	"github.com/schorlet/gcp/world/api/pb"
	"github.com/schorlet/gcp/world/hello/domain"
	"github.com/schorlet/gcp/world/hello/grpc/hello"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	addr       string
	grpcServer *grpc.Server
}

func NewServer(addr string, helloService domain.HelloService) (*Server, error) {
	opts := []grpc.ServerOption{
		// grpc.StatsHandler(h stats.Handler),
		grpc.UnaryInterceptor(interceptor.UnaryServerLogger),
	}

	certFile := os.Getenv("HELLO_SERVER_CRT_FILE")
	keyFile := os.Getenv("HELLO_SERVER_KEY_FILE")
	if certFile != "" && keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterHelloServer(grpcServer, hello.NewServer(helloService))
	reflection.Register(grpcServer)

	return &Server{
		addr:       addr,
		grpcServer: grpcServer,
	}, nil
}

func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) ListenAndServeTLS() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(lis)
}
