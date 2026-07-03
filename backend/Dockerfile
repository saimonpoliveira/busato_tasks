# Single-stage build (mais estável no Railway que multi-stage)
FROM golang:1.22-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o server ./cmd/server

RUN mkdir -p /app/uploads

ENV GIN_MODE=release
ENV UPLOAD_DIR=/app/uploads

EXPOSE 8080

CMD ["./server"]
