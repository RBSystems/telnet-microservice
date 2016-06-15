FROM golang:1.6.2

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/telnet-microservice

WORKDIR /go/src/github.com/byuoitav/telnet-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/telnet-microservice"]

EXPOSE 8001
