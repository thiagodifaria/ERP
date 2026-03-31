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

  local attempts=0
  until "${COMPOSE_CMD[@]}" exec -T service-postgresql pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1; do
    attempts=$((attempts + 1))

    if [[ "$attempts" -ge 30 ]]; then
      echo "[db] postgresql did not become ready in time"
      exit 1
    fi

    sleep 1
  done
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/db.sh up
  ./scripts/db.sh migrate common
  ./scripts/db.sh migrate identity
  ./scripts/db.sh migrate crm
  ./scripts/db.sh migrate workflow-control
  ./scripts/db.sh migrate all
  ./scripts/db.sh seed identity
  ./scripts/db.sh seed crm
  ./scripts/db.sh seed workflow-control
  ./scripts/db.sh seed all
  ./scripts/db.sh summary identity [tenant-slug]
  ./scripts/db.sh summary crm [tenant-slug]
  ./scripts/db.sh summary workflow-control [tenant-slug]
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
        crm)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/crm/migrations"
          ;;
        workflow-control)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/migrations"
          ;;
        all)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/common/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/crm/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/migrations"
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
        crm)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/crm/seeds"
          ;;
        workflow-control)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/seeds"
          ;;
        all)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/seeds"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/crm/seeds"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/seeds"
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
        crm)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id) AS leads,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'captured') AS captured,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'contacted') AS contacted,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'qualified') AS qualified,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'disqualified') AS disqualified,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.owner_user_public_id IS NOT NULL) AS assigned,
              (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.owner_user_public_id IS NULL) AS unassigned
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        workflow-control)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id) AS definitions,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'draft') AS draft,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'active') AS active,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'archived') AS archived
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
