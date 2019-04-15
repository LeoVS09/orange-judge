FROM golang:1.12.4-stretch

WORKDIR /go/src/orange-judge

RUN go get github.com/fatih/color
COPY . .

RUN go get -v ./...

RUN go build -o server -v main.go && \
    chmod +x server

EXPOSE 3010

#ENTRYPOINT ["./docker-entrypoint.sh"]

CMD ["/bin/bash", "./server",  "-d"]
