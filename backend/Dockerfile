# Build
FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/server

# Run (Debian slim — mais estável que Alpine no Railway)
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /app/uploads

ENV GIN_MODE=release
ENV UPLOAD_DIR=/app/uploads

EXPOSE 8080

CMD ["./server"]
