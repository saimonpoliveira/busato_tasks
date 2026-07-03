#!/bin/sh
set -e

PORT="${PORT:-80}"

echo "Starting frontend on port ${PORT}..."

if [ -f /etc/nginx/conf.d/default.conf ]; then
  sed -i "s/listen 80;/listen ${PORT};/" /etc/nginx/conf.d/default.conf
  exec nginx -g 'daemon off;'
fi

if [ -d dist ]; then
  exec npx --yes serve dist -s -l "${PORT}"
fi

echo "No dist/ or nginx config found"
exit 1
