# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

RUN mkdir -p /app

WORKDIR /app

ADD . /app

RUN go build ./main.go

CMD ["./main"]