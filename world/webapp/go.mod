module github.com/schorlet/gcp/world/webapp

go 1.14

require (
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	github.com/schorlet/gcp/world/api v0.0.0-00010101000000-000000000000
)

replace github.com/schorlet/gcp/world/api => ../api
