NAME:=goat
VERSION:=$(shell git describe --tags)
OS:=linux
ARCH:=amd64
GOAT_FILES:=$$(find . -name '*.go' | grep -v vendor)
BINPATH:=$(PWD)/bin
PKGDIR?=$(PWD)/pkg
DEBIANDIR:=debian/$(NAME)-$(VERSION)

all: build

docker-build:
	mkdir -p $(PKGDIR)
	docker build -q -t "goat-builder" -f Dockerfile.build .
	docker run -v $(PKGDIR):/goat-pkg goat-builder

build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -mod=vendor -a -tags netgo -ldflags '-w -extldflags "-static" -X main.VERSION=$(VERSION)' -o $(BINPATH)/$(NAME)
	strip $(BINPATH)/$(NAME)
	echo $(VERSION) > $(BINPATH)/version-file

pkgclean:
	rm -rf $(PKGDIR)

tarball:
	@tar -czf $(PKGDIR)/$(NAME)_$(VERSION)_$(OS)_$(ARCH).tar.gz -C $(BINPATH) $(NAME)

zip:
	@zip $(PKGDIR)/$(NAME)_$(VERSION)_$(OS)_$(ARCH).zip -j $(BINPATH)/$(NAME)

test:
	@go test -v ./...

fmt:
	@gofmt -s -w $(GOAT_FILES)

dev-env: ## Build a local development environment using Docker
	@docker run -it --rm \
		-v $(shell pwd):/$(NAME) \
		-w /$(NAME) \
		golang:1.12 \
		/bin/bash -c 'make install; bash'

install: ## Build and install locally the binary (dev purpose)
	go install .

.PHONY: docker-build dev-env clean install test
