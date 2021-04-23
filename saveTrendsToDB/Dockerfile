FROM golang:1.16-alpine

ENV GOOS=windows

ENV GOARCH=amd64

ADD ./saveTrendsToDB ./saveTrendsToDB

WORKDIR ./saveTrendsToDB

RUN go mod tidy

RUN go build SaveTrendsOfWeek.go

RUN mv /go/saveTrendsToDB/SaveTrendsOfWeek.exe /tmp