# syntax=docker/dockerfile:1.7
#
# Command server — raw TCP (:8080) + UDP (:8081) real-time relay.
# Pure Go, no cgo, so it runs from a scratch-like alpine image.
#
# Must track apps/command/go.mod's `go` directive.
FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod ./
COPY go.su[m] ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/command .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tini && adduser -D -u 10001 app
WORKDIR /app
COPY --from=builder /out/command /app/command
USER app

EXPOSE 8080/tcp 8081/udp
ENTRYPOINT ["/sbin/tini", "--"]
# Servers only; the simulated test clients are for `go run . -mode=all` locally.
CMD ["/app/command", "-mode=server"]
