#!/usr/bin/env bash
set -euo pipefail

# Este script centraliza a aplicacao das migrations do PostgreSQL em ambiente local.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker-compose.yml"
ENV_FILE="$ROOT_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  ENV_FILE="$ROOT_DIR/.env.example"
fi

set -a
source "$ENV_FILE"
set +a

COMPOSE_CMD=(docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE")
DB_NAME="${ERP_POSTGRES_DB:-erp}"
DB_USER="${ERP_POSTGRES_USER:-erp}"

run_psql_file() {
  local file_path="$1"
  echo "[db] applying $(basename "$file_path")"
  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" < "$file_path"
}

apply_directory() {
  local directory_path="$1"

  if [[ ! -d "$directory_path" ]]; then
    echo "[db] directory not found: $directory_path"
    exit 1
  fi

  while IFS= read -r file_path; do
    run_psql_file "$file_path"
  done < <(find "$directory_path" -maxdepth 1 -type f -name '*.sql' | sort)
}

ensure_postgres() {
  "${COMPOSE_CMD[@]}" up -d service-postgresql
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/db.sh up
  ./scripts/db.sh migrate common
  ./scripts/db.sh migrate identity
  ./scripts/db.sh migrate all
  ./scripts/db.sh psql
EOF
}

main() {
  local command="${1:-}"
  local scope="${2:-}"

  case "$command" in
    up)
      ensure_postgres
      ;;
    migrate)
      ensure_postgres
      case "$scope" in
        common)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/common/migrations"
          ;;
        identity)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/migrations"
          ;;
        all)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/common/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/migrations"
          ;;
        *)
          usage
          exit 1
          ;;
      esac
      ;;
    psql)
      ensure_postgres
      "${COMPOSE_CMD[@]}" exec service-postgresql psql -U "$DB_USER" -d "$DB_NAME"
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
