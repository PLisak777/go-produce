# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

CMD ["./go-produce"]