module github.com/schorlet/gcp/world/hello

go 1.14

require (
	github.com/schorlet/gcp/world/api v0.0.0-20200923193627-66e87d3a957d
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sync v0.1.0
	google.golang.org/grpc v1.53.0
)

// replace github.com/schorlet/gcp/world/api => ../api
