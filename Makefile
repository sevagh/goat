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

