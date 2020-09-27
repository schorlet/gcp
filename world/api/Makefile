
GOOS      ?= $(shell go env GOOS)
GOARCH    ?= $(shell go env GOARCH)
SHORT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)

GO_LDFLAGS  := -X 'api.Version=$(SHORT_SHA)'
ifeq ($(STATIC),y)
	# -w Omit DWARF symbol table
	# -s Omit symbol table and debug info
	# -extldflags -static and CGO_ENABLED=0 to make pure static
	GO_LDFLAGS += -w -s -extldflags "-static"
	export CGO_ENABLED=0
endif
export GO111MODULE=on

.DEFAULT_GOAL := api

GENERATED_FILES := pb/hello.pb.go pb/hello.pb.gw.go pb/hello.swagger.json
.PHONY: api
## api: generate protobuf and grpc files (default target)
api: $(GENERATED_FILES)

pb/hello.pb.go: pb/hello.proto
# 	protoc -I=./pb -I=/opt/protoc-3.12.4/include \
# 		--go_out=./pb \
# 		--go_opt=paths=source_relative \
# 		./pb/hello.proto
	protoc -I=./pb -I=/opt/protoc-3.12.4/include \
		--go_out=plugins=grpc:./pb \
		--go_opt=paths=source_relative \
		./pb/hello.proto

# pb/hello_grpc.pb.go: pb/hello.proto
# 	protoc -I=./pb -I=/opt/protoc-3.12.4/include \
# 		--go-grpc_out=./pb \
# 		--go-grpc_opt=paths=source_relative \
# 		./pb/hello.proto

pb/hello.pb.gw.go: pb/hello.proto pb/hello.yaml
	protoc -I=./pb -I=/opt/protoc-3.12.4/include \
		--grpc-gateway_out=logtostderr=true,grpc_api_configuration=./pb/hello.yaml:./pb \
		 --grpc-gateway_opt=paths=source_relative \
		./pb/hello.proto

pb/hello.swagger.json: pb/hello.proto pb/hello.yaml
# 	protoc -I=./pb -I=/opt/protoc-3.12.4/include \
# 		--swagger_out=logtostderr=true,grpc_api_configuration=./pb/hello.yaml:./pb \
# 		./pb/hello.proto
	protoc -I=./pb -I=/opt/protoc-3.12.4/include \
		--openapiv2_out=logtostderr=true,grpc_api_configuration=./pb/hello.yaml:./pb \
		./pb/hello.proto


.PHONY: install
## install: install the application
install: api
	goimports -l -w .
	go install \
		-v -trimpath -mod=readonly \
		-ldflags "$(GO_LDFLAGS)" \
		./...


.PHONY: clean
## clean: clean installed files
clean:
	go clean -i -x
	touch pb/*.proto
	#rm -fv $(GENERATED_FILES)


.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##/ /p' $(MAKEFILE_LIST) | column -t -s ':'
