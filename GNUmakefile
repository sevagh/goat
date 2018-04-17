VERSION:=0.6.0
GOAT_FILES?=$$(find . -name '*.go' | grep -v vendor)

all: build

build: deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static" -X main.VERSION=$(VERSION)' -o  bin/goat

test:
	@go vet ./...
	@go test -v ./...

deps:
	@command -v dep 2>&1 >/dev/null || go get -u github.com/golang/dep/cmd/dep
	@dep ensure

lint:
	@gofmt -s -w $(GOAT_FILES)
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

.PHONY: clean test
