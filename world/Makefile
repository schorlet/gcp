
.DEFAULT_GOAL := all

.PHONY: all
## all: build all (default target)
all:
	$(MAKE) -C ./api install
	$(MAKE) -C ./hello build
	$(MAKE) -C ./webapp build

.PHONY: clean
## clean: clean all
clean:
	$(MAKE) -C ./api clean
	$(MAKE) -C ./hello clean
	$(MAKE) -C ./webapp clean

.PHONY: test
## test: run tests
test:
	$(MAKE) -C ./hello test


.PHONY: run-hello
## run-hello: run the hello server
run-hello:
	$(MAKE) -C ./hello run

.PHONY: run-webapp
## run-hello: run the webapp server
run-webapp:
	$(MAKE) -C ./webapp run



.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##/ /p' $(MAKEFILE_LIST) | column -t -s ':'
