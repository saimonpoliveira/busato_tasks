#!/bin/sh
set -e

PORT="${PORT:-3000}"
API_URL="${VITE_API_URL:-${API_URL:-/api/v1}}"

echo "Generating runtime config..."
echo "window.__ENV__ = { API_URL: \"${API_URL}\" };" > dist/env-config.js
echo "API URL: ${API_URL}"

echo "Starting frontend on port ${PORT}..."
exec serve -s dist --listen "tcp://0.0.0.0:${PORT}"
