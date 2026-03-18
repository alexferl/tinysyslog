FROM golang:1.26-alpine AS builder
LABEL maintainer="Alexandre Ferland <me@alexferl.com>"

WORKDIR /build

RUN apk add --no-cache git

RUN addgroup -g 65532 -S nonroot && \
    adduser  -u 65532 -S -D -G nonroot -s /sbin/nologin nonroot

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG GOOS=linux
ARG GOARCH=amd64

ENV GOOS=$GOOS
ENV GOARCH=$GOARCH

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /build/app ./cmd/tinysyslogd

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /build/app /app

USER nonroot

EXPOSE 5140/tcp 5140/udp

ENTRYPOINT ["/app", "--bind-addr", "0.0.0.0:5140"]
