FROM alpine AS base
RUN apk add --no-cache curl wget

FROM golang:1.17 AS go-builder
WORKDIR /go/app
COPY . /go/app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/app/ /go/app/server.go

FROM base
COPY --from=go-builder /go/app/ /main
RUN ls
RUN cd main
RUN ls
CMD ["/main"]