VERSION:=0.4.1
GOAT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOAT_NAME=$(notdir $(shell pwd))

STATIC_ENV:=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
STATIC_FLAGS:=-a -tags netgo -ldflags '-extldflags "-static" -X main.VERSION=$(VERSION)'
RELEASE_FLAGS:=-a -tags netgo -ldflags '-w -extldflags "-static" -X main.VERSION=$(VERSION)'

DIST_DIR=dist/
BIN_DIR=dist/bin

all: build_static

builddir:
	@mkdir -p $(DIST_DIR) $(BIN_DIR)

build:
	@cd cmd/goat && go build $(DYNAMIC_FLAGS) -o ../../$(BIN_DIR)/$(GOAT_NAME)

build_static: builddir
	@cd cmd/goat && go build $(STATIC_FLAGS) -o ../../$(BIN_DIR)/$(GOAT_NAME)

release: builddir
	@cd cmd/goat && go build $(RELEASE_FLAGS) -o ../../$(BIN_DIR)/$(GOAT_NAME)

deps:
	@go get -u github.com/golang/dep
	@dep ensure

test: build lint
	@go vet .
	@go test -v ./...

lint:
	@gofmt -s -w $(GOAT_FILES)

clean:
	-rm -rf build

package:
	@GOAT_VERSION=$(VERSION) $(MAKE) -C ./centos-package/

.PHONY: clean test
