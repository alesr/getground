FROM golang:1.17-alpine3.14

ENV WAIT_VERSION 2.9.0

WORKDIR /go/src/app

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait

RUN chmod +x /wait

COPY . .

RUN go fmt -x -mod=vendor ./...

RUN go build -v -o party cmd/party/main.go

EXPOSE 3000