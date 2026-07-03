#!/bin/sh
set -e

echo "Building Busato Tasks API..."
go build -ldflags="-s -w" -o server ./cmd/server

echo "Starting Busato Tasks API..."
exec ./server
