FROM golang:1.7-alpine

RUN apk update && apk upgrade && apk add git

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/telnet-microservice

WORKDIR /go/src/github.com/byuoitav/telnet-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/telnet-microservice"]

EXPOSE 8001
