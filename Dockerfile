FROM golang:1.24.2-alpine AS builder
LABEL maintainer="Alexandre Ferland <me@alexferl.com>"

WORKDIR /build

RUN apk add --no-cache git
RUN adduser -D -u 1337 tinysyslog

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/tinysyslogd

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/tinysyslogd /tinysyslogd

USER tinysyslog

EXPOSE 5140/tcp 5140/udp

ENTRYPOINT ["/tinysyslogd", "--bind-addr", "0.0.0.0:5140"]
