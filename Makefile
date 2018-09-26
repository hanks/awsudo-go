VERSION = 1.0.0
CUR_DIR = $(shell pwd)
NAME = awsudo
WORKSPACE = /go/src/github.com/hanks/awsudo-go
DEV_IMAGE = hanks/awsudo-go-dev:1.0.0
OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')

.PHONY: build clean dev debug install push run test uninstall

default: test

build: test clean
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./dist/bin/$(NAME)_linux_amd64_$(VERSION) main.go
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o ./dist/bin/$(NAME)_darwin_amd64_$(VERSION) main.go

clean:
	rm -rf ./dist

dev:
	docker build -t $(DEV_IMAGE) .

debug:
	docker run -it --rm --security-opt=seccomp:unconfined -e CGO_ENABLED=0 -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) dlv debug main.go

install:
	cp ./dist/bin/$(NAME)_$(OS)_amd64_$(VERSION) /usr/local/bin/$(NAME)

push:
	docker push $(DEV_IMAGE)

run:
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) go run main.go help

test:
	echo "unit test..."
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) go vet ./...
	#docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) golint -set_exit_status $(go list ./...)

uninstall:
	rm /usr/local/bin/$(NAME)
