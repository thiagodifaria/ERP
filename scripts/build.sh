#!/usr/bin/env bash
set -euo pipefail

# Este script centraliza o build inicial das imagens locais do ERP.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker-compose.yml"
ENV_FILE="$ROOT_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  ENV_FILE="$ROOT_DIR/.env.example"
fi

docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" build "$@"
