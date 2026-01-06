# build stage
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o healthcheck

# runtime stage
FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/healthcheck .

RUN apt-get update && apt-get install -y \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash healthuser

RUN mkdir -p /app/logs \
    && chown -R healthuser:healthuser /app/logs

USER healthuser

EXPOSE 8081

CMD ["./healthcheck"]