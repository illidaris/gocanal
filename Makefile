GO ?= go
GOOS ?= linux
GOARCH ?= amd64
BINARY := ./bin/gocanal

.PHONY: build tar

# Disable cgo so release binaries do not depend on the build host's glibc.
build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -trimpath -o $(BINARY) .

tar: build
	tar -czvf gocanal.tar.gz ./bin