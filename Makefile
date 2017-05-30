build:
	@go build .

install:
	@go install .

deps:
	@go get -u github.com/golang/dep
	@dep ensure

test: build
	@golint
	@go fmt . 
	@go vet .
	@go test -v

