NAME:=goat
VERSION:=0.9.0
OS:=linux
ARCH:=amd64
GOAT_FILES?=$$(find . -name '*.go' | grep -v vendor)

all: build

build: deps
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -a -tags netgo -ldflags '-w -extldflags "-static" -X main.VERSION=$(VERSION)' -o usr/bin/$(NAME)
	strip usr/bin/$(NAME)

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

package_all: pkgclean build deb rpm zip

zip:
	@zip pkg/$(NAME)_$(VERSION)_$(OS)_$(ARCH).zip -j usr/bin/$(NAME)

deb:
	@mkdir -p pkg
	fpm -s dir -t deb -n $(NAME) -v $(VERSION) -C . \
		-p pkg/$(NAME)_VERSION_ARCH.deb \
		-d "mdadm" \
		--deb-systemd ./goat@.service \
		usr/bin

rpm:
	@mkdir -p pkg
	fpm -s dir -t rpm -n $(NAME) -v $(VERSION) -C . \
		-p pkg/$(NAME)_VERSION_ARCH.rpm \
		-d "mdadm" \
		--rpm-systemd ./goat@.service \
		usr/bin

pkgclean:
	@rm -rf pkg

lintsetup:
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install 2>&1 >/dev/null
	@go install ./...

dev-env: ## Build a local development environment using Docker
	@docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/sevagh/$(NAME) \
		-w /go/src/github.com/sevagh/$(NAME) \
		golang:1.10 \
		/bin/bash -c 'make deps; make install; bash'

install: ## Build and install locally the binary (dev purpose)
	go install .

.PHONY: dev-env clean install test
