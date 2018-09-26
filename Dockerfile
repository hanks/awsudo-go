FROM golang:1.10.3-alpine3.8

WORKDIR /go/src/github.com/hanks/awsudo-go

RUN apk --no-cache update && \
        apk add --no-cache python py-pip py-setuptools ca-certificates git gcc && \
        pip --no-cache-dir install awscli && \
        go get -u github.com/derekparker/delve/cmd/dlv && \
        go get github.com/golang/lint/golint && \
        mkdir -p ./dist/bin && \
        rm -rf /var/cache/apk/*

CMD ["sh"]
