FROM golang:1.15-alpine

ENV GOOS=windows

ENV GOARCH=amd64

ADD ./databaseConnector ./databaseConnector

WORKDIR ./databaseConnector

RUN go mod tidy

RUN go build SaveTrendsOfWeek.go