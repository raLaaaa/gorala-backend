FROM ubuntu:12.04

RUN apt-get install -y python-software-properties

RUN add-apt-repository ppa:duh/golang
RUN apt-get update
RUN apt-get install -y golang

ADD server.go /var/server/server.go

EXPOSE 4334

CMD ["go", "run", "/var/server/server.go"]