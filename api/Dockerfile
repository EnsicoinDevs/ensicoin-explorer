FROM golang:latest

RUN go get -u github.com/golang/dep/...

RUN mkdir -p /go/src/github.com/EnsicoinDevs/ensicoin-explorer/api/db
WORKDIR /go/src/github.com/EnsicoinDevs/ensicoin-explorer/api

COPY . .

RUN go mod download

RUN go install

ENTRYPOINT /go/bin/api

EXPOSE 8080
