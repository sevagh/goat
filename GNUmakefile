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

build: builddir deps
	@cd cmd/goat && go build $(DYNAMIC_FLAGS) -o ../../$(BIN_DIR)/$(GOAT_NAME)

build_static: builddir deps
	@cd cmd/goat && go build $(STATIC_FLAGS) -o ../../$(BIN_DIR)/$(GOAT_NAME)

release: builddir deps
	@cd cmd/goat && go build $(RELEASE_FLAGS) -o ../../$(BIN_DIR)/$(GOAT_NAME)

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

test:
	@go vet ./pkg/...
	@go vet ./cmd/goat/...
	@go test -v ./pkg/...
	@go test -v ./cmd/goat/...

lint: lintsetup
	@gofmt -s -w $(GOAT_FILES)
	@gometalinter.v2 --enable-all $(GOAT_FILES) --exclude=_test.go

lintsetup:
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install 2>&1 >/dev/null

clean:
	-rm -rf build

package:
	@GOAT_VERSION=$(VERSION) $(MAKE) -C ./centos-package/

.PHONY: clean test
