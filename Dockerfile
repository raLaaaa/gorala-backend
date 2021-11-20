# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go install

COPY *.go ./

RUN go build -o /gorala

EXPOSE 4334

CMD [ "/gorala" ]