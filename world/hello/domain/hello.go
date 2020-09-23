package domain

type Greeting struct {
	Message  string `json:"message,omitempty"`
	Version  string `json:"version,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

type HelloService interface {
	Hello(name string) (*Greeting, error)
}
