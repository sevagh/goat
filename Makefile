VERSION := 0.2.0
RELEASEBIN := $(CURDIR)/goat
INSTALLBIN := /usr/bin/goat

build:
	@go build .

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
	@rpmbuild -ba specfile.spec --define "_sourcedir $$PWD"
