FROM golang:1.21.11

RUN go install github.com/cosmtrek/air@v1.48.0

WORKDIR /go/src/simple-app
