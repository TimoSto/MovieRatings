FROM golang:1.16-alpine

ENV GOOS=windows

ENV GOARCH=amd64

ADD ./SaveTrendsToDB ./saveTrendsToDB

WORKDIR ./SaveTrendsToDB

RUN go mod tidy

RUN go build SaveTrendsOfWeek.go