APP ?= hello
PORT ?= 8021

GOOS      ?= $(shell go env GOOS)
GOARCH    ?= $(shell go env GOARCH)
TARGET    := $(APP)-$(GOOS)-$(GOARCH)
SHORT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)

GOFILES     := $(shell find . -name '*.go' -a ! -name '*_test.go')
GOTESTFILES := $(shell find . -name '*_test.go')

GO_LDFLAGS  := -X 'main.Version=$(SHORT_SHA)'
ifeq ($(STATIC),y)
	# -w Omit DWARF symbol table
	# -s Omit symbol table and debug info
	# -extldflags -static and CGO_ENABLED=0 to make pure static
	GO_LDFLAGS += -w -s -extldflags "-static"
	export CGO_ENABLED=0
endif
export GO111MODULE=on

.DEFAULT_GOAL := build

.PHONY: build
## build: build the application (default target)
build: $(TARGET)

$(TARGET): $(GOFILES)
	goimports -l -w .
	go build \
		-v -trimpath -mod=readonly \
		-ldflags "$(GO_LDFLAGS)" \
		-o $(TARGET)

.PHONY: build-debug
## build-debug: build without optimisation
build-debug:
	goimports -l -w .
	go build \
		-v -mod=readonly \
		-gcflags "-N -l" \
		-o $(TARGET)

.PHONY: install
## install: install the application
install:
	goimports -l -w .
	go install \
		-v -trimpath -mod=readonly \
		-ldflags "$(GO_LDFLAGS)"


.PHONY: lint
## lint: lint source files
lint:
	gofmt -e -d -s $(GOFILES) $(GOTESTFILES)
	go vet
	golangci-lint run ./...

.PHONY: clean
## clean: clean generated files
clean:
	go clean -i -x
	rm -f $(APP)-$(GOOS)-* cover.out


.PHONY: run
## run: run the application
run: $(TARGET)
	# pkill $(TARGET)
	PORT=$(PORT) ./$(TARGET)

.PHONY: attach
## attach: debug the application
attach: build-debug
# 	dlv debug github.com/schorlet/gcp/world/hello
# 	dlv exec ./$(TARGET)
	dlv attach $(shell pgrep $(TARGET))


.PHONY: test
## test: run test
test: $(TARGET)
	goimports -l -w .
	go test -v ./...

.PHONY: cover
## cover: test coverage
cover: cover.out

cover.out: $(TARGET)
	goimports -l -w .
	go test -v -coverprofile=cover.out ./...
	go tool cover -func=cover.out


.PHONY: image
## image: build a docker image
image: Dockerfile
	docker image build \
		--tag $(APP):v1 \
		--build-arg SHORT_SHA=$(SHORT_SHA) \
		-f Dockerfile .

.PHONY: run-container
## run-container: run a docker container
run-container:
	docker container run --rm -d -p $(PORT):$(PORT) $(APP):v1

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##/ /p' $(MAKEFILE_LIST) | column -t -s ':'
