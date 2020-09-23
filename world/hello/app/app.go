package app

import (
	"os"

	"github.com/schorlet/gcp/world/hello/domain"
)

type HelloService struct {
	version string
}

func NewHelloService(version string) *HelloService {
	return &HelloService{version}
}

func (service *HelloService) Hello(name string) (*domain.Greeting, error) {
	if name == "panic" {
		panic("fake panic")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &domain.Greeting{
		Message:  "Hello " + name,
		Version:  service.version,
		Hostname: hostname,
	}, nil
}
