FROM golang:alpine as build

LABEL maintainer "github.com/jxsl13"

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /build
COPY . ./
COPY go.* ./

ENV CGO_ENABLED=0
ENV GOOS=linux 

RUN go get -d && go build -a -ldflags '-w -extldflags "-static"' -o monitor-zcatch .


FROM alpine:latest as minimal

ENV MONITOR_BROKER_ADDRESS=tcp://mosquitto:1883

WORKDIR /app
COPY --from=build /build/monitor-zcatch .
VOLUME ["/data", "/app/.env"]
ENTRYPOINT ["/app/monitor-zcatch"]
