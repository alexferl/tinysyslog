FROM golang:1.11.2-alpine as builder
MAINTAINER Alexandre Ferland <aferlandqc@gmail.com>

ENV GO111MODULE=on

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /build/tinysyslog /tinysyslog
EXPOSE 5140 5140/udp
ENTRYPOINT ["/tinysyslog", "--address", "0.0.0.0:5140"]
