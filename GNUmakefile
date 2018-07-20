NAME:=goat
VERSION:=0.7.0
GOAT_FILES?=$$(find . -name '*.go' | grep -v vendor)

all: build

build: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static" -X main.VERSION=$(VERSION)' -o bin/$(NAME)
	strip bin/$(NAME)

test:
	@go vet ./...
	@go test -v ./...

deps:
	@command -v dep 2>&1 >/dev/null || go get -u github.com/golang/dep/cmd/dep
	@dep ensure -v

fmt:
	@gofmt -s -w $(GOAT_FILES)

lint:
	-gometalinter.v2 --enable-all $(GOAT_FILES) --exclude=_test.go

lintsetup:
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install 2>&1 >/dev/null
	@go install ./...

clean:
	-rm -rf bin

rpm: build
	@cp bin/goat rpm-package/
	GOAT_VERSION=$(VERSION) $(MAKE) -C ./rpm-package/

dev-env: ## Build a local development environment using Docker
	@docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/sevagh/$(NAME) \
		-w /go/src/github.com/sevagh/$(NAME) \
		golang:1.10 \
		/bin/bash -c 'make deps; make install; bash'

install: ## Build and install locally the binary (dev purpose)
	go install .

.PHONY: dev-env clean install test
