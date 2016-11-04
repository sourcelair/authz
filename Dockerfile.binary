FROM golang:1.7-alpine
MAINTAINER Antonis Kalipetis <akalipetis@sourcelair.com>

WORKDIR /go/src/github.com/twistlock/authz
CMD sh

RUN apk update && apk add git make && \
    go get github.com/tools/godep && \
    go get golang.org/x/sys/unix && \
    go get -v github.com/golang/lint/golint

COPY Godeps /go/src/github.com/twistlock/authz/Godeps
RUN godep restore -v

COPY . /go/src/github.com/twistlock/authz
RUN make test
