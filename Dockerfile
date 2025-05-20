FROM golang:bookworm AS builder

WORKDIR /build
COPY . .
RUN go build -o ./bin/statesman ./cmd/...

FROM debian:bookworm
COPY --from=builder /build/bin/statesman /usr/bin/statesman