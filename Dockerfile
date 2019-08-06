FROM golang:1.8

WORKDIR /go/src/app
COPY src/ .

CMD ["go", "run", "app/server.go"]