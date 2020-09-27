module github.com/schorlet/gcp/world/hello

go 1.14

require (
	github.com/schorlet/gcp/world/api v0.0.0-20200923193627-66e87d3a957d
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/grpc v1.31.1
)

// replace github.com/schorlet/gcp/world/api => ../api
