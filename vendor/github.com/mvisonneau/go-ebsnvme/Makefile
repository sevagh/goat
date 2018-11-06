NAME          := go-ebsnvme
VERSION       := $(shell git describe --tags --abbrev=1)
FILES         := $(shell git ls-files '*.go')
LDFLAGS       := -w -extldflags "-static" -X 'main.version=$(VERSION)'
REGISTRY      := mvisonneau/$(NAME)
.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Format source code
	@command -v goimports 2>&1 >/dev/null || go get -u golang.org/x/tools/cmd/goimports
	goimports -w $(FILES)

.PHONY: lint
lint: ## Run golint and go vet against the codebase
	@command -v golint 2>&1 >/dev/null || go get -u github.com/golang/lint/golint
	golint -set_exit_status .
	go vet ./...

.PHONY: test
test: ## Run the tests against the codebase
	go test -v ./...

.PHONY: install
install: ## Build and install locally the binary (dev purpose)
	go install .

.PHONY: build
build: ## Build the binary
	@command -v gox 2>&1 >/dev/null || go get -u github.com/mitchellh/gox
	mkdir -p dist; rm -rf dist/*
	CGO_ENABLED=0 gox -osarch "linux/386 linux/amd64" -ldflags "$(LDFLAGS)" -output dist/$(NAME)_{{.OS}}_{{.Arch}}
	strip dist/*_linux_*

.PHONY: build-docker
build-docker:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" .
	strip $(NAME)

.PHONY: publish-github
publish-github: ## Send the binaries onto the GitHub release
	@command -v ghr 2>&1 >/dev/null || go get -u github.com/tcnksm/ghr
	ghr -u mvisonneau -replace $(VERSION) dist

.PHONY: deps
deps: ## Fetch all dependencies
	@command -v dep 2>&1 >/dev/null || go get -u github.com/golang/dep/cmd/dep
	@dep ensure -v

.PHONY: imports
imports: ## Fixes the syntax (linting) of the codebase
	goimports -d $(FILES)

.PHONY: clean
clean: ## Remove binary if it exists
	rm -f $(NAME)

.PHONY: coverage
coverage: ## Generates coverage report
	rm -rf *.out
	go test -coverprofile=coverage.out

.PHONY: dev-env
dev-env: ## Build a local development environment using Docker
	@docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/mvisonneau/$(NAME) \
		-w /go/src/github.com/mvisonneau/$(NAME) \
		golang:1.11 \
		/bin/bash -c 'make deps; make install; bash'

.PHONY: all
all: lint imports test coverage build ## Test, builds and ship package for all supported platforms

.PHONY: help
help: ## Displays this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
