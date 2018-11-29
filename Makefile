VERSION = 1.0.0
CUR_DIR = $(shell pwd)
NAME = awsudo
WORKSPACE = /go/src/github.com/hanks/awsudo-go
DEV_IMAGE = hanks/awsudo-go-dev:1.1.0
OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')
CMD ?= help
AWSUDO_DEBUG ?= false

.PHONY: build clean dev debug install push run test uninstall

default: test

build: test clean
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) -e "CGO_ENABLED=0" -e "GOARCH=amd64" -e "GOOS=linux" $(DEV_IMAGE) go build -o ./dist/bin/$(NAME)_linux_amd64_$(VERSION) main.go
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) -e "CGO_ENABLED=0" -e "GOARCH=amd64" -e "GOOS=darwin" $(DEV_IMAGE) go build -o ./dist/bin/$(NAME)_darwin_amd64_$(VERSION) main.go

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
	docker run -it --rm -e AWSUDO_DEBUG=$(AWSUDO_DEBUG) -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) go run main.go $(CMD)

simple-test:
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) sh -c 'go test -v -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v /configs | grep -v /version)'
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) go tool cover -html=coverage.out -o coverage.html

test:
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) sh -c 'go vet $$(go list ./...)'
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) sh -c 'golint -set_exit_status $$(go list ./...)'
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) sh -c 'go test -v -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v /configs | grep -v /version)'
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) go tool cover -html=coverage.out -o coverage.html

coveralls:
	docker run -it --rm -v $(CUR_DIR):$(WORKSPACE) $(DEV_IMAGE) goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $(COVERALLS_TOKEN)

uninstall:
	rm /usr/local/bin/$(NAME)
