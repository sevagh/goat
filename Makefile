VERSION:=0.4.0
GOAT_FILES?=$$(find . -name '*.go' | grep -v vendor)

build:
	@go build -ldflags "-X main.VERSION=$(VERSION)" .

install:
	@go install .

deps:
	@go get -u github.com/golang/dep
	@dep ensure

test: build lint
	@go vet .
	@go test -v ./...

lint:
	@gofmt -s -w $(GOAT_FILES)

package:
	@GOAT_VERSION=$(VERSION) $(MAKE) -C ./.pkg/
