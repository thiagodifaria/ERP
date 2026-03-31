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

run_psql_query() {
  local query="$1"
  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -c "$query"
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
  ./scripts/db.sh seed identity
  ./scripts/db.sh seed all
  ./scripts/db.sh summary identity [tenant-slug]
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
    seed)
      ensure_postgres
      case "$scope" in
        identity)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/seeds"
          ;;
        all)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/seeds"
          ;;
        *)
          usage
          exit 1
          ;;
      esac
      ;;
    summary)
      ensure_postgres
      case "$scope" in
        identity)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM identity.companies AS company WHERE company.tenant_id = tenant.id) AS companies,
              (SELECT count(*) FROM identity.users AS \"user\" WHERE \"user\".tenant_id = tenant.id) AS users,
              (SELECT count(*) FROM identity.teams AS team WHERE team.tenant_id = tenant.id) AS teams,
              (SELECT count(*) FROM identity.roles AS role WHERE role.tenant_id = tenant.id) AS roles,
              (SELECT count(*) FROM identity.team_memberships AS membership WHERE membership.tenant_id = tenant.id) AS team_memberships,
              (SELECT count(*) FROM identity.user_roles AS user_role WHERE user_role.tenant_id = tenant.id) AS user_roles
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
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
