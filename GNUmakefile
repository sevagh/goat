VERSION:=0.5.0
GOAT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOAT_CMDS=$(shell find cmd/ -maxdepth 1 -mindepth 1 -type d)

STATIC_ENV:=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
STATIC_FLAGS:=-a -tags netgo -ldflags '-extldflags "-static" -X main.VERSION=$(VERSION)'
RELEASE_FLAGS:=-a -tags netgo -ldflags '-w -extldflags "-static" -X main.VERSION=$(VERSION)'

all: build_static

build: deps
	@$(foreach cmd,$(GOAT_CMDS),\
		cd $(cmd) &&\
		$(STATIC_ENV) go build $(STATIC_FLAGS) \
			-o ../../bin/$(notdir $(cmd)) &&\
		cd - 2>&1 >/dev/null;)

release: deps
	@$(foreach cmd,$(GOAT_CMDS),\
		cd $(cmd) &&\
		$(STATIC_ENV) go build $(RELEASE_FLAGS) \
			-o ../../bin/$(notdir $(cmd)) &&\
		cd - 2>&1 >/dev/null;)

test:
	@$(foreach cmd,$(GOAT_CMDS),\
		go vet ./$(cmd) &&\
		go test -v ./$(cmd);)
	@go vet ./pkg/...
	@go test -v ./pkg/...

deps:
	@command -v dep 2>&1 >/dev/null || go get -u github.com/golang/dep/cmd/dep
	@dep ensure

lint:
	@gofmt -s -w $(GOAT_FILES)
	-gometalinter.v2 --enable-all $(GOAT_FILES) --exclude=_test.go

lintsetup:
	@go get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter.v2 --install 2>&1 >/dev/null

clean:
	-rm -rf bin

package: release
	@GOAT_VERSION=$(VERSION) $(MAKE) -C ./rpm-package/

.PHONY: clean test
