FROM golang:1.20.7-alpine3.18 as builder
MAINTAINER Alexandre Ferland <me@alexferl.com>

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/tinysyslogd

FROM scratch
COPY --from=builder /build/tinysyslogd /tinysyslogd

EXPOSE 5140 5140/udp

ENTRYPOINT ["/tinysyslogd", "--bind-addr", "0.0.0.0:5140"]
