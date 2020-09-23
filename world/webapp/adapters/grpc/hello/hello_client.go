package hello

import (
	"net"
	"net/http"
	"os"

	"github.com/schorlet/gcp/world/api/grpc/hello"
	"github.com/schorlet/gcp/world/api/pb"
)

func NewClient() (*hello.Client, error) {
	helloHost, helloPort := os.Getenv("HELLO_SERVICE_HOST"),
		os.Getenv("HELLO_SERVICE_PORT")

	if helloHost == "" {
		helloHost = "localhost"
	}
	if helloPort == "" {
		helloPort = "8012"
	}

	return hello.NewClient(
		net.JoinHostPort(helloHost, helloPort),
	)
}

func HandleWithClient(client pb.HelloClient) (http.HandlerFunc, error) {
	return hello.HandleWithClient(client)
}
