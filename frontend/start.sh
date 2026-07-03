#!/bin/sh
set -e

PORT="${PORT:-3000}"

echo "Starting frontend on port ${PORT}..."

# serve v14+ exige formato tcp://host:port
exec serve -s dist --listen "tcp://0.0.0.0:${PORT}"
