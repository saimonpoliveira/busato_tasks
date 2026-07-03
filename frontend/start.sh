#!/bin/sh
set -e

PORT="${PORT:-80}"

echo "Starting frontend on port ${PORT}..."

# Debian nginx config path
if [ -f /etc/nginx/sites-available/default ]; then
  sed -i "s/listen 80;/listen ${PORT};/" /etc/nginx/sites-available/default
  exec nginx -g 'daemon off;'
fi

if [ -d /usr/share/nginx/html ]; then
  exec npx --yes serve /usr/share/nginx/html -s -l "${PORT}"
fi

echo "No nginx config or static files found"
exit 1
