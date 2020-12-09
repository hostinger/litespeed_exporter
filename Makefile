GO 	 	    := go
GO111MODULE := on
GOHOSTOS    := $(shell $(GO) env GOHOSTOS)
GOHOSTARCH  := $(shell $(GO) env GOHOSTARCH)

BINARY := litespeed_exporter

PLATFORMS := darwin/amd64 dragonfly/amd64 freebsd/amd64 linux/amd64 netbsd/amd64 openbsd/amd64 windows/amd64 linux/386 freebsd/386 netbsd/386 openbsd/386 windows/386

BUILD_VERSION  := $(shell cat VERSION)
BUILD_DATE 	   := $(shell date '+%FT%T')
BUILD_REVISION := $(shell git rev-parse HEAD)

LDFLAGS=-ldflags "-X main.Version=$(BUILD_VERSION) -X main.Date=$(BUILD_DATE) -X main.Revision=$(BUILD_REVISION)"

all: check test clean build

.PHONY: test
test:
	$(GO) test -v -cover -race -coverprofile=coverage.txt ./...

.PHONY: check
check:
	$(GO) vet ./...

.PHONY: clean
clean:
	rm -rf bin/*

build:
	GO111MODULE=$(GO111MODULE) GOOS=$(GOHOSTOS) GOARCH=$(GOHOSTARCH) $(GO) build $(LDFLAGS) -v -o bin/$(BINARY)-$(BUILD_VERSION)

temp = $(subst /, ,$@)
OS = $(word 1, $(temp))
ARCH = $(word 2, $(temp))

.PHONY: build-all $(PLATFORMS)
build-all: $(PLATFORMS)

$(PLATFORMS):
	GO111MODULE=$(GO111MODULE) GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(LDFLAGS) -v -o bin/$(BINARY)-$(BUILD_VERSION).$(OS)-$(ARCH)
