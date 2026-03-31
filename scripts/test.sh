#!/usr/bin/env bash
set -euo pipefail

# Este script centraliza validacoes tecnicas em modo container-first.

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

run_go_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-golang/edge:/workspace" \
    -w /workspace \
    golang:1.24-alpine \
    go test ./...

  docker run --rm \
    -v "$ROOT_DIR/service-api/service-golang/crm:/workspace" \
    -w /workspace \
    golang:1.24-alpine \
    go test ./...
}

run_dotnet_build() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-csharp/identity:/workspace" \
    -w /workspace \
    mcr.microsoft.com/dotnet/sdk:8.0 \
    dotnet test tests/Identity.UnitTests/Identity.UnitTests.csproj -c Release
}

run_dotnet_integration() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-csharp/identity:/workspace" \
    -w /workspace \
    mcr.microsoft.com/dotnet/sdk:8.0 \
    dotnet test tests/Identity.IntegrationTests/Identity.IntegrationTests.csproj -c Release
}

run_dotnet_contract() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-csharp/identity:/workspace" \
    -w /workspace \
    mcr.microsoft.com/dotnet/sdk:8.0 \
    dotnet test tests/Identity.ContractTests/Identity.ContractTests.csproj -c Release
}

run_rust_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-rust/webhook-hub:/workspace" \
    -w /workspace \
    rust:1 \
    cargo test
}

run_identity_database_smoke() {
  local smoke_slug="smoke-identity-bootstrap"
  local summary
  local crm_summary

  bash "$ROOT_DIR/scripts/db.sh" up
  bash "$ROOT_DIR/scripts/db.sh" migrate all

  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" -c "
      INSERT INTO identity.tenants (public_id, slug, display_name, status)
      SELECT gen_random_uuid(), '$smoke_slug', 'Smoke Identity Bootstrap', 'active'
      WHERE NOT EXISTS (
        SELECT 1
        FROM identity.tenants
        WHERE slug = '$smoke_slug'
      );
    "

  bash "$ROOT_DIR/scripts/db.sh" seed all
  bash "$ROOT_DIR/scripts/db.sh" seed all

  summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        tenant.slug || '|' ||
        (SELECT count(*) FROM identity.companies AS company WHERE company.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM identity.users AS \"user\" WHERE \"user\".tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM identity.teams AS team WHERE team.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM identity.roles AS role WHERE role.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM identity.team_memberships AS membership WHERE membership.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM identity.user_roles AS user_role WHERE user_role.tenant_id = tenant.id)
      FROM identity.tenants AS tenant
      WHERE tenant.slug = '$smoke_slug';
    ")"

  echo "[test] identity db smoke => $summary"

  if [[ "$summary" != "$smoke_slug|1|1|1|5|1|1" ]]; then
    echo "[test] unexpected identity db smoke summary"
    exit 1
  fi

  crm_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        tenant.slug || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'captured') || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'contacted') || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'qualified') || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.status = 'disqualified') || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.owner_user_public_id IS NOT NULL) || '|' ||
        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id AND lead.owner_user_public_id IS NULL)
      FROM identity.tenants AS tenant
      WHERE tenant.slug = '$smoke_slug';
    ")"

  echo "[test] crm db smoke => $crm_summary"

  if [[ "$crm_summary" != "$smoke_slug|1|1|0|0|0|1|0" ]]; then
    echo "[test] unexpected crm db smoke summary"
    exit 1
  fi
}

wait_for_http_ready() {
  local url="$1"
  local attempts=0

  until curl -fsS "$url" >/dev/null 2>&1; do
    attempts=$((attempts + 1))

    if [[ "$attempts" -ge 30 ]]; then
      echo "[test] http endpoint did not become ready: $url"
      exit 1
    fi

    sleep 1
  done
}

run_crm_runtime_smoke() {
  local base_url="http://localhost:${CRM_HTTP_PORT:-8083}"
  local owner_public_id
  local lead_list
  local create_response
  local created_public_id
  local owner_response
  local status_response
  local summary_response

  "${COMPOSE_CMD[@]}" up -d --build crm
  wait_for_http_ready "$base_url/health/ready"

  lead_list="$(curl -fsS "$base_url/api/crm/leads")"
  echo "[test] crm api list => $lead_list"

  if [[ "$lead_list" != *'"email":"lead@bootstrap-ops.local"'* ]]; then
    echo "[test] bootstrap CRM lead was not returned by the live API"
    exit 1
  fi

  owner_public_id="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT \"user\".public_id::text
      FROM identity.users AS \"user\"
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = \"user\".tenant_id
      WHERE tenant.slug = 'bootstrap-ops'
        AND \"user\".email = 'owner@bootstrap-ops.local';
    ")"

  if [[ -z "$owner_public_id" ]]; then
    echo "[test] bootstrap CRM owner public id was not found"
    exit 1
  fi

  create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"name":"Runtime Lead","email":"runtime.lead@example.com","source":"whatsapp"}' \
    "$base_url/api/crm/leads")"
  echo "[test] crm api create => $create_response"

  created_public_id="$(echo "$create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  if [[ -z "$created_public_id" ]]; then
    echo "[test] runtime lead public id was not returned by create response"
    exit 1
  fi

  owner_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d "{\"ownerUserId\":\"$owner_public_id\"}" \
    "$base_url/api/crm/leads/$created_public_id/owner")"
  echo "[test] crm api owner => $owner_response"

  if [[ "$owner_response" != *"\"ownerUserId\":\"$owner_public_id\""* ]]; then
    echo "[test] runtime lead owner update did not persist"
    exit 1
  fi

  status_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"contacted"}' \
    "$base_url/api/crm/leads/$created_public_id/status")"
  echo "[test] crm api status => $status_response"

  if [[ "$status_response" != *'"status":"contacted"'* ]]; then
    echo "[test] runtime lead status update did not persist"
    exit 1
  fi

  summary_response="$(curl -fsS "$base_url/api/crm/leads/summary")"
  echo "[test] crm api summary => $summary_response"

  if [[ "$summary_response" != *'"total":2'* || "$summary_response" != *'"assigned":2'* || "$summary_response" != *'"contacted":1'* ]]; then
    echo "[test] runtime CRM summary did not reflect live updates"
    exit 1
  fi
}

run_smoke() {
  trap 'bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true' RETURN
  bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true
  "${COMPOSE_CMD[@]}" ps
  run_identity_database_smoke
  run_crm_runtime_smoke
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/test.sh unit
  ./scripts/test.sh integration
  ./scripts/test.sh contract
  ./scripts/test.sh smoke
  ./scripts/test.sh all
EOF
}

main() {
  local command="${1:-}"

  case "$command" in
    unit)
      run_go_unit
      run_dotnet_build
      run_rust_unit
      ;;
    integration)
      run_dotnet_integration
      ;;
    contract)
      run_dotnet_contract
      ;;
    smoke)
      run_smoke
      ;;
    all)
      run_go_unit
      run_dotnet_build
      run_dotnet_integration
      run_dotnet_contract
      run_rust_unit
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
