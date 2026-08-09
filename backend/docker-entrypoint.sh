#!/bin/sh
set -e

# Run database migrations before starting the server unless explicitly disabled.
if [ "${AUTO_MIGRATE:-true}" = "true" ]; then
  echo "[entrypoint] running database migrations..."
  /migrate up
else
  echo "[entrypoint] AUTO_MIGRATE=false, skipping migrations"
fi

echo "[entrypoint] starting server..."
exec /server
