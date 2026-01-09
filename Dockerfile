# build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY app/go.mod ./
RUN go mod tidy

COPY app/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o healthcheck

# runtime stage
FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/healthcheck .

RUN apt-get update && apt-get install -y \
    sudo \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN update-ca-certificates

RUN useradd -m -s /bin/bash healthuser

RUN mkdir -p /app/logs /app/data \
    && chown -R healthuser:healthuser /app/logs /app/data

USER healthuser

EXPOSE 8081

CMD ["./healthcheck"]