#!/usr/bin/env bash
set -euo pipefail

# ERP - comando central de build, runtime e operacao local.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker-compose.yml"
ENV_FILE="$ROOT_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  ENV_FILE="$ROOT_DIR/.env.example"
fi

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

log() { echo -e "${BLUE}[erp]${NC} $*"; }
ok() { echo -e "${GREEN}[ok]${NC} $*"; }
warn() { echo -e "${YELLOW}[warn]${NC} $*"; }
fail() { echo -e "${RED}[erro]${NC} $*" >&2; exit 1; }
section() {
  echo ""
  echo -e "${CYAN}== $* ==${NC}"
}

COMPOSE_CMD=(docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE")
DB_NAME="${ERP_POSTGRES_DB:-erp}"
DB_USER="${ERP_POSTGRES_USER:-erp}"
MIGRATION_DOMAINS=(
  common
  identity
  crm
  sales
  rentals
  finance
  billing
  documents
  analytics
  simulation
  catalog
  platform-control
  support
  supplier
  notification
  fiscal
  engagement
  webhook-hub
  workflow-control
  workflow-runtime
)
SEED_DOMAINS=(
  identity
  crm
  sales
  engagement
  workflow-control
)

load_env_file_preserving_env() {
  local line
  local key

  while IFS= read -r line || [[ -n "$line" ]]; do
    [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]] && continue
    key="${line%%=*}"

    if [[ -n "${!key+x}" ]]; then
      continue
    fi

    export "$line"
  done < "$ENV_FILE"
}

is_tcp_port_in_use() {
  local port="$1"

  if command -v ss >/dev/null 2>&1; then
    ss -H -ltn "( sport = :$port )" 2>/dev/null | grep -q .
    return
  fi

  if command -v netstat >/dev/null 2>&1; then
    netstat -ltn 2>/dev/null | awk '{print $4}' | grep -Eq "(^|:)$port$"
    return
  fi

  (echo >"/dev/tcp/127.0.0.1/$port") >/dev/null 2>&1
}

find_available_port() {
  local port="$1"

  while is_tcp_port_in_use "$port"; do
    port=$((port + 1))
  done

  echo "$port"
}

remap_host_port_if_needed() {
  local variable_name="$1"
  local label="$2"
  local requested_port="${!variable_name:-}"

  if [[ -z "$requested_port" || -n "${ERP_HOST_PORTS_LOCKED:-}" ]]; then
    return
  fi

  if ! is_tcp_port_in_use "$requested_port"; then
    return
  fi

  local fallback_start=$((requested_port + 1000))
  local fallback_port
  fallback_port="$(find_available_port "$fallback_start")"
  export "$variable_name=$fallback_port"
  warn "porta $label remapeada de $requested_port para $fallback_port"
}

prepare_runtime_ports() {
  remap_host_port_if_needed "POSTGRES_PORT" "postgresql"
  remap_host_port_if_needed "REDIS_PORT" "redis"
  remap_host_port_if_needed "KAFKA_PORT" "kafka"
  remap_host_port_if_needed "PROMETHEUS_PORT" "prometheus"
  remap_host_port_if_needed "GRAFANA_PORT" "grafana"
  remap_host_port_if_needed "KEYCLOAK_PORT" "keycloak"
  remap_host_port_if_needed "OPENFGA_HTTP_PORT" "openfga-http"
  remap_host_port_if_needed "OPENFGA_GRPC_PORT" "openfga-grpc"
  remap_host_port_if_needed "OPENFGA_PLAYGROUND_PORT" "openfga-playground"
  remap_host_port_if_needed "GATEWAY_HTTP_PORT" "gateway"
  remap_host_port_if_needed "EDGE_HTTP_PORT" "edge"
  remap_host_port_if_needed "IDENTITY_HTTP_PORT" "identity"
  remap_host_port_if_needed "WEBHOOK_HUB_HTTP_PORT" "webhook-hub"
  remap_host_port_if_needed "CRM_HTTP_PORT" "crm"
  remap_host_port_if_needed "WORKFLOW_CONTROL_HTTP_PORT" "workflow-control"
  remap_host_port_if_needed "WORKFLOW_RUNTIME_HTTP_PORT" "workflow-runtime"
  remap_host_port_if_needed "ANALYTICS_HTTP_PORT" "analytics"
  remap_host_port_if_needed "SIMULATION_HTTP_PORT" "simulation"
  remap_host_port_if_needed "SALES_HTTP_PORT" "sales"
  remap_host_port_if_needed "ENGAGEMENT_HTTP_PORT" "engagement"
  remap_host_port_if_needed "FINANCE_HTTP_PORT" "finance"
  remap_host_port_if_needed "DOCUMENTS_HTTP_PORT" "documents"
  remap_host_port_if_needed "RENTALS_HTTP_PORT" "rentals"
  remap_host_port_if_needed "CATALOG_HTTP_PORT" "catalog"
  remap_host_port_if_needed "PLATFORM_CONTROL_HTTP_PORT" "platform-control"
  remap_host_port_if_needed "SUPPORT_HTTP_PORT" "support"
  remap_host_port_if_needed "SUPPLIER_HTTP_PORT" "supplier"
  remap_host_port_if_needed "NOTIFICATION_HTTP_PORT" "notification"
  remap_host_port_if_needed "FISCAL_HTTP_PORT" "fiscal"
  export ERP_HOST_PORTS_LOCKED=1
}

check_dependencies() {
  section "Dependencias"

  command -v docker >/dev/null 2>&1 || fail "Docker nao encontrado"
  docker compose version >/dev/null 2>&1 || fail "Docker Compose nao encontrado"

  ok "Docker e Docker Compose disponiveis"
}

compose_build() {
  section "Build"
  "${COMPOSE_CMD[@]}" build "$@"
}

compose_up() {
  section "Subindo stack"
  prepare_runtime_ports
  "${COMPOSE_CMD[@]}" up -d --build "$@"
}

compose_down() {
  section "Derrubando stack"
  "${COMPOSE_CMD[@]}" down --remove-orphans "$@"
}

compose_logs() {
  "${COMPOSE_CMD[@]}" logs -f "$@"
}

compose_ps() {
  "${COMPOSE_CMD[@]}" ps "$@"
}

ensure_postgres() {
  prepare_runtime_ports
  "${COMPOSE_CMD[@]}" up -d service-postgresql

  local attempts=0
  until "${COMPOSE_CMD[@]}" exec -T service-postgresql pg_isready -U "$DB_USER" -d "$DB_NAME" >/dev/null 2>&1; do
    attempts=$((attempts + 1))

    if [[ "$attempts" -ge 30 ]]; then
      fail "PostgreSQL nao ficou pronto a tempo"
    fi

    sleep 1
  done
}

run_psql_file() {
  local file_path="$1"
  log "aplicando $(basename "$file_path")"
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
  local required="${2:-required}"

  if [[ ! -d "$directory_path" ]]; then
    if [[ "$required" == "optional" ]]; then
      return
    fi
    fail "diretorio nao encontrado: $directory_path"
  fi

  while IFS= read -r file_path; do
    run_psql_file "$file_path"
  done < <(find "$directory_path" -maxdepth 1 -type f -name '*.sql' | sort)
}

database_up() {
  section "PostgreSQL"
  ensure_postgres
  ok "PostgreSQL pronto"
}

migrate_domain() {
  local domain="$1"
  apply_directory "$ROOT_DIR/service-api/service-postgresql/$domain/migrations"
}

seed_domain() {
  local domain="$1"
  apply_directory "$ROOT_DIR/service-api/service-postgresql/$domain/seeds"
}

database_migrate() {
  local scope="${1:-}"
  [[ -n "$scope" ]] || fail "uso: ./scripts/build.sh migrate <contexto|all>"

  section "Migrations: $scope"
  ensure_postgres

  if [[ "$scope" == "all" ]]; then
    local domain
    for domain in "${MIGRATION_DOMAINS[@]}"; do
      migrate_domain "$domain"
    done
    ok "migrations aplicadas"
    return
  fi

  migrate_domain "$scope"
  ok "migration aplicada para $scope"
}

database_seed() {
  local scope="${1:-}"
  [[ -n "$scope" ]] || fail "uso: ./scripts/build.sh seed <contexto|all>"

  section "Seeds: $scope"
  ensure_postgres

  if [[ "$scope" == "all" ]]; then
    local domain
    for domain in "${SEED_DOMAINS[@]}"; do
      apply_directory "$ROOT_DIR/service-api/service-postgresql/$domain/seeds" optional
    done
    ok "seeds aplicados"
    return
  fi

  seed_domain "$scope"
  ok "seed aplicado para $scope"
}

database_backup() {
  local output_path="${1:-/tmp/erp-local-backup.sql}"
  local output_directory

  section "Backup"
  ensure_postgres
  output_directory="$(dirname "$output_path")"
  mkdir -p "$output_directory"

  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    pg_dump --clean --if-exists --no-owner --no-privileges -U "$DB_USER" -d "$DB_NAME" > "$output_path"
  ok "backup salvo em $output_path"
}

database_backup_encrypted() {
  local output_path="${1:-/tmp/erp-local-backup.sql.enc}"
  local temporary_backup

  [[ -n "${ERP_BACKUP_ENCRYPTION_KEY:-}" ]] || fail "ERP_BACKUP_ENCRYPTION_KEY obrigatorio para backup criptografado"
  command -v openssl >/dev/null 2>&1 || fail "openssl nao encontrado para backup criptografado"

  temporary_backup="$(mktemp "${TMPDIR:-/tmp}/erp-backup-plain-XXXXXX.sql")"
  database_backup "$temporary_backup"
  mkdir -p "$(dirname "$output_path")"
  openssl enc -aes-256-cbc -pbkdf2 -salt -pass env:ERP_BACKUP_ENCRYPTION_KEY -in "$temporary_backup" -out "$output_path"
  rm -f "$temporary_backup"
  ok "backup criptografado salvo em $output_path"
}

database_restore() {
  local input_path="${1:-}"
  [[ -n "$input_path" ]] || fail "uso: ./scripts/build.sh restore <arquivo.sql>"
  [[ -f "$input_path" ]] || fail "backup nao encontrado: $input_path"

  section "Restore"
  ensure_postgres
  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" < "$input_path"
  ok "restore aplicado"
}

database_restore_encrypted() {
  local input_path="${1:-}"
  local temporary_backup

  [[ -n "$input_path" ]] || fail "uso: ./scripts/build.sh restore-encrypted <arquivo.sql.enc>"
  [[ -f "$input_path" ]] || fail "backup criptografado nao encontrado: $input_path"
  [[ -n "${ERP_BACKUP_ENCRYPTION_KEY:-}" ]] || fail "ERP_BACKUP_ENCRYPTION_KEY obrigatorio para restore criptografado"
  command -v openssl >/dev/null 2>&1 || fail "openssl nao encontrado para restore criptografado"

  temporary_backup="$(mktemp "${TMPDIR:-/tmp}/erp-backup-decrypted-XXXXXX.sql")"
  openssl enc -d -aes-256-cbc -pbkdf2 -pass env:ERP_BACKUP_ENCRYPTION_KEY -in "$input_path" -out "$temporary_backup"
  database_restore "$temporary_backup"
  rm -f "$temporary_backup"
  ok "restore criptografado aplicado"
}

database_summary() {
  local schema="${1:-}"
  [[ -n "$schema" ]] || fail "uso: ./scripts/build.sh summary <schema>"

  section "Resumo: $schema"
  ensure_postgres
  run_psql_query "
    SELECT
      schemaname,
      relname AS table_name,
      n_live_tup AS estimated_rows
    FROM pg_stat_user_tables
    WHERE schemaname = replace('$schema', '-', '_')
    ORDER BY relname;
  "
}

database_psql() {
  ensure_postgres
  "${COMPOSE_CMD[@]}" exec service-postgresql psql -U "$DB_USER" -d "$DB_NAME"
}

full_build() {
  check_dependencies
  compose_down >/dev/null 2>&1 || true
  compose_build
  compose_up
  database_migrate all
  database_seed all
  compose_ps
}

usage() {
  cat <<'EOF'
Uso:
  ./scripts/build.sh                         build completo, sobe stack, migra e semeia
  ./scripts/build.sh build [servico...]      constroi imagens
  ./scripts/build.sh up [servico...]         sobe stack local com build
  ./scripts/build.sh down [-v]               derruba stack local
  ./scripts/build.sh restart                 reinicia stack local
  ./scripts/build.sh logs [servico...]       acompanha logs
  ./scripts/build.sh ps                      lista containers
  ./scripts/build.sh migrate <contexto|all>  aplica migrations
  ./scripts/build.sh seed <contexto|all>     aplica seeds
  ./scripts/build.sh backup [arquivo.sql]    gera dump do PostgreSQL
  ./scripts/build.sh backup-encrypted [arquivo.sql.enc] gera dump criptografado
  ./scripts/build.sh restore <arquivo.sql>   restaura dump do PostgreSQL
  ./scripts/build.sh restore-encrypted <arquivo.sql.enc> restaura dump criptografado
  ./scripts/build.sh summary <schema>        resumo relacional simples
  ./scripts/build.sh psql                    abre psql no PostgreSQL local
  ./scripts/build.sh db <comando>            alias para comandos de banco

Testes continuam centralizados em:
  ./scripts/test.sh unit|integration|contract|platform|smoke|performance|backup-restore|hardening|all
EOF
}

dispatch_database() {
  local command="${1:-}"
  shift || true

  case "$command" in
    up) database_up "$@" ;;
    migrate) database_migrate "$@" ;;
    seed) database_seed "$@" ;;
    backup) database_backup "$@" ;;
    backup-encrypted) database_backup_encrypted "$@" ;;
    restore) database_restore "$@" ;;
    restore-encrypted) database_restore_encrypted "$@" ;;
    summary) database_summary "$@" ;;
    psql) database_psql "$@" ;;
    *) usage; exit 1 ;;
  esac
}

main() {
  cd "$ROOT_DIR"
  load_env_file_preserving_env
  DB_NAME="${ERP_POSTGRES_DB:-erp}"
  DB_USER="${ERP_POSTGRES_USER:-erp}"

  local command="${1:-}"
  if [[ -z "$command" ]]; then
    full_build
    return
  fi
  shift || true

  case "$command" in
    build | build-only) check_dependencies; compose_build "$@" ;;
    up | start) check_dependencies; compose_up "$@" ;;
    down | stop) compose_down "$@" ;;
    restart) compose_down; compose_up ;;
    logs) compose_logs "$@" ;;
    ps | status) compose_ps "$@" ;;
    migrate) database_migrate "$@" ;;
    seed) database_seed "$@" ;;
    backup) database_backup "$@" ;;
    backup-encrypted) database_backup_encrypted "$@" ;;
    restore) database_restore "$@" ;;
    restore-encrypted) database_restore_encrypted "$@" ;;
    summary) database_summary "$@" ;;
    psql) database_psql "$@" ;;
    db) dispatch_database "$@" ;;
    help | --help | -h) usage ;;
    *) usage; exit 1 ;;
  esac
}

main "$@"
