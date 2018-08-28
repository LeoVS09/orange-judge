FROM golang:1.11-stretch

WORKDIR /go/src/orange-judge

ADD . .
RUN go get -v ./...

RUN go build -o server -v main.go

EXPOSE 3010

CMD ["/go/src/orange-judge/server",  "-d", "-tc=false"]
