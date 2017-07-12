VERSION := 0.3.0
RELEASEBIN := $(CURDIR)/goat
INSTALLBIN := /usr/bin/goat

build:
	@go build -ldflags "-X main.VERSION=$(VERSION)" .

install:
	@go install .

deps:
	@go get -u github.com/golang/dep
	@dep ensure

test: build
	@go fmt . 
	@go vet .
	@go test -v

download:
	@curl -L https://github.com/sevagh/goat/releases/download/$(VERSION)/goat --output $(RELEASEBIN)
	@chmod +x $(RELEASEBIN)

install_release:
	@cp $(RELEASEBIN) $(INSTALLBIN)

rpm:
	@rpmlint specfile.spec
	@rpmbuild -ba specfile.spec --define "_sourcedir $$PWD" --define "_version $(VERSION)"
