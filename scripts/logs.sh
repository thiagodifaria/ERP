#!/usr/bin/env bash
set -euo pipefail

# Este script acompanha os logs do ecossistema local do ERP.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker-compose.yml"
ENV_FILE="$ROOT_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  ENV_FILE="$ROOT_DIR/.env.example"
fi

docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" logs -f "$@"
