#!/bin/sh
set -e

echo "Building frontend..."
npm ci
npm run build

echo "Starting frontend..."
exec ./start.sh
