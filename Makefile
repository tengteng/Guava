GOPATH := ${PWD}:${PWD}/_vendor:${GOPATH}
export GOPATH

default: install

env: env
	go env

install: fmt
	go install ./...

build: fmt vet
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

vendor_clean:
	rm -rf ./_vendor/src

clean:
	rm -rf ./pkg ./bin
