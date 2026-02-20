#!/usr/bin/env bash
set -euo pipefail

# Load .env if present
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi

echo "Starting goauthly-fiber-example in dev mode..."
go run ./cmd/server/
