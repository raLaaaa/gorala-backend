FROM golang:1.17
ADD server.go /go/src/server.go
EXPOSE 4334
WORKDIR /go/src
ENTRYPOINT [ "go", "run",  "server.go"]