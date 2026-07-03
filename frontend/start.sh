#!/bin/sh
set -e

echo "Building frontend..."
npm ci
npm run build

echo "Starting frontend on port ${PORT:-3000}..."
exec npx serve dist -s -l "${PORT:-3000}"
