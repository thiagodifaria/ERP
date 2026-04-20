#!/usr/bin/env bash
set -euo pipefail

# Este script centraliza validacoes tecnicas em modo container-first.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker-compose.yml"
ENV_FILE="$ROOT_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  ENV_FILE="$ROOT_DIR/.env.example"
fi

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

load_env_file_preserving_env

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

  if [[ -z "$requested_port" ]]; then
    return
  fi

  if ! is_tcp_port_in_use "$requested_port"; then
    return
  fi

  local fallback_start=$((requested_port + 1000))
  local fallback_port
  fallback_port="$(find_available_port "$fallback_start")"
  export "$variable_name=$fallback_port"
  echo "[ports] remapped $label host port from $requested_port to $fallback_port because it is already in use"
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
  remap_host_port_if_needed "EDGE_HTTP_PORT" "edge"
  remap_host_port_if_needed "IDENTITY_HTTP_PORT" "identity"
  remap_host_port_if_needed "WEBHOOK_HUB_HTTP_PORT" "webhook-hub"
  remap_host_port_if_needed "CRM_HTTP_PORT" "crm"
  remap_host_port_if_needed "WORKFLOW_CONTROL_HTTP_PORT" "workflow-control"
  remap_host_port_if_needed "WORKFLOW_RUNTIME_HTTP_PORT" "workflow-runtime"
  remap_host_port_if_needed "ANALYTICS_HTTP_PORT" "analytics"
  remap_host_port_if_needed "SALES_HTTP_PORT" "sales"
  remap_host_port_if_needed "ENGAGEMENT_HTTP_PORT" "engagement"
}

prepare_runtime_ports
export ERP_HOST_PORTS_LOCKED=1

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

  docker run --rm \
    -v "$ROOT_DIR/service-api/service-golang/sales:/workspace" \
    -w /workspace \
    golang:1.24-alpine \
    go test ./...
}

run_typescript_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-typescript/workflow-control:/workspace" \
    -w /workspace \
    node:22-alpine \
    sh -lc "npm install && npm run test:unit && rm -rf node_modules dist"

  docker run --rm \
    -v "$ROOT_DIR/service-api/service-typescript/engagement:/workspace" \
    -w /workspace \
    node:22-alpine \
    sh -lc "npm install && npm run test:unit && rm -rf node_modules dist"
}

run_typescript_contract() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-typescript/workflow-control:/workspace" \
    -w /workspace \
    node:22-alpine \
    sh -lc "npm install && npm run test:contract && rm -rf node_modules dist"

  docker run --rm \
    -v "$ROOT_DIR/service-api/service-typescript/engagement:/workspace" \
    -w /workspace \
    node:22-alpine \
    sh -lc "npm install && npm run test:contract && rm -rf node_modules dist"
}

run_elixir_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-elixir/workflow-runtime:/workspace" \
    -w /workspace \
    elixir:1.17-alpine \
    sh -lc "apk add --no-cache build-base git >/dev/null && mix local.hex --force >/dev/null && mix local.rebar --force >/dev/null && mix deps.get >/dev/null && mix test"
}

run_python_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-python/analytics:/workspace" \
    -w /workspace \
    python:3.12-slim \
    sh -lc "pip install --no-cache-dir -e .[dev] >/dev/null && pytest && rm -rf .pytest_cache analytics.egg-info"
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

run_go_contract() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-golang/crm:/workspace" \
    -w /workspace \
    golang:1.24-alpine \
    go test -tags=contract ./tests/contract/...

  docker run --rm \
    -v "$ROOT_DIR/service-api/service-golang/sales:/workspace" \
    -w /workspace \
    golang:1.24-alpine \
    go test -tags=contract ./tests/contract/...
}

run_rust_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-rust/webhook-hub:/workspace" \
    -w /workspace \
    rust:1 \
    cargo test
}

run_platform_runtime_smoke() {
  local keycloak_url="http://localhost:${KEYCLOAK_PORT:-8089}"
  local openfga_url="http://localhost:${OPENFGA_HTTP_PORT:-8090}"
  local prometheus_url="http://localhost:${PROMETHEUS_PORT:-9090}"
  local grafana_url="http://localhost:${GRAFANA_PORT:-3000}"
  local realm_response
  local openfga_store_response
  local openfga_list_response
  local prometheus_targets_response
  local prometheus_probe_response
  local grafana_response
  local kafka_list_response

  "${COMPOSE_CMD[@]}" up -d service-postgresql service-redis service-kafka service-keycloak service-openfga service-blackbox-exporter service-prometheus service-grafana
  wait_for_http_ready "$keycloak_url/realms/${KEYCLOAK_REALM:-erp-local}"
  wait_for_http_ready "$openfga_url/healthz"
  wait_for_http_ready "$prometheus_url/-/ready"
  wait_for_http_ready "$grafana_url/api/health"
  wait_for_prometheus_probe_success "$prometheus_url" "blackbox-http" "http://service-keycloak:8080/realms/${KEYCLOAK_REALM:-erp-local}"
  wait_for_prometheus_probe_success "$prometheus_url" "blackbox-http" "http://service-openfga:8080/healthz"
  wait_for_prometheus_probe_success "$prometheus_url" "blackbox-http" "http://service-grafana:3000/api/health"
  wait_for_prometheus_probe_success "$prometheus_url" "blackbox-tcp" "service-postgresql:5432"
  wait_for_prometheus_probe_success "$prometheus_url" "blackbox-tcp" "service-redis:6379"
  wait_for_prometheus_probe_success "$prometheus_url" "blackbox-tcp" "service-kafka:9092"

  realm_response="$(curl -fsS "$keycloak_url/realms/${KEYCLOAK_REALM:-erp-local}")"
  openfga_store_response="$(curl -fsS -X POST "$openfga_url/stores" -H 'content-type: application/json' -d "{\"name\":\"${OPENFGA_STORE_NAME:-erp-local-store}\"}")"
  openfga_list_response="$(curl -fsS "$openfga_url/stores")"
  prometheus_targets_response="$(curl -fsS "$prometheus_url/api/v1/targets")"
  prometheus_probe_response="$(curl -fsS --get --data-urlencode 'query=probe_success' "$prometheus_url/api/v1/query")"
  grafana_response="$(curl -fsS "$grafana_url/api/health")"
  kafka_list_response="$("${COMPOSE_CMD[@]}" exec -T service-kafka kafka-topics --bootstrap-server service-kafka:9092 --list || true)"

  echo "[test] platform keycloak realm => $realm_response"
  echo "[test] platform openfga store create => $openfga_store_response"
  echo "[test] platform openfga stores => $openfga_list_response"
  echo "[test] platform prometheus targets => $prometheus_targets_response"
  echo "[test] platform prometheus probes => $prometheus_probe_response"
  echo "[test] platform grafana health => $grafana_response"
  echo "[test] platform kafka topics => $kafka_list_response"

  if [[ "$realm_response" != *"\"realm\":\"${KEYCLOAK_REALM:-erp-local}\""* || "$openfga_store_response" != *"\"name\":\"${OPENFGA_STORE_NAME:-erp-local-store}\""* || "$openfga_list_response" != *"\"name\":\"${OPENFGA_STORE_NAME:-erp-local-store}\""* || "$prometheus_targets_response" != *"service-keycloak:8080/realms/${KEYCLOAK_REALM:-erp-local}"* || "$prometheus_targets_response" != *'service-openfga:8080/healthz'* || "$prometheus_targets_response" != *'service-postgresql:5432'* || "$prometheus_targets_response" != *'service-kafka:9092'* || "$prometheus_probe_response" != *'"status":"success"'* || "$grafana_response" != *'"database"'* || "$grafana_response" != *'"ok"'* ]]; then
    echo "[test] platform runtime stack did not expose the expected local foundations"
    exit 1
  fi
}

run_webhook_hub_runtime_smoke() {
  local base_url="http://localhost:${WEBHOOK_HUB_HTTP_PORT:-8082}"
  local health_details_response
  local list_response
  local create_response
  local duplicate_response
  local created_public_id
  local validate_response
  local queue_response
  local process_response
  local forward_response
  local transitions_response
  local detail_response
  local filtered_response
  local summary_response
  local db_summary

  "${COMPOSE_CMD[@]}" up -d --build webhook-hub
  wait_for_http_ready "$base_url/health/ready"

  health_details_response="$(curl -fsS "$base_url/health/details")"
  list_response="$(curl -fsS "$base_url/api/webhook-hub/events")"
  echo "[test] webhook-hub health details => $health_details_response"
  echo "[test] webhook-hub list => $list_response"

  if [[ "$health_details_response" != *'"name":"signature-validation","status":"ready"'* || "$health_details_response" != *'"name":"postgresql","status":"ready"'* || "$list_response" != '[]' ]]; then
    echo "[test] webhook-hub bootstrap runtime state was not clean"
    exit 1
  fi

  create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"provider":"stripe","event_type":"payment.succeeded","external_id":"evt_runtime_001","payload_summary":"Pagamento confirmado pelo smoke do webhook-hub."}' \
    "$base_url/api/webhook-hub/events")"
  duplicate_response="$(curl -sS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"provider":"stripe","event_type":"payment.succeeded","external_id":"evt_runtime_001","payload_summary":"Pagamento duplicado do smoke do webhook-hub."}' \
    -w ' HTTP_STATUS:%{http_code}' \
    "$base_url/api/webhook-hub/events")"
  created_public_id="$(echo "$create_response" | sed -n 's/.*"public_id":"\([^"]*\)".*/\1/p')"
  validate_response="$(curl -fsS -X POST "$base_url/api/webhook-hub/events/$created_public_id/validate")"
  queue_response="$(curl -fsS -X POST "$base_url/api/webhook-hub/events/$created_public_id/queue")"
  process_response="$(curl -fsS -X POST "$base_url/api/webhook-hub/events/$created_public_id/process")"
  forward_response="$(curl -fsS -X POST "$base_url/api/webhook-hub/events/$created_public_id/forward")"
  transitions_response="$(curl -fsS "$base_url/api/webhook-hub/events/$created_public_id/transitions")"
  detail_response="$(curl -fsS "$base_url/api/webhook-hub/events/$created_public_id")"
  filtered_response="$(curl -fsS "$base_url/api/webhook-hub/events?provider=stripe&event_type=payment.succeeded&status=forwarded")"
  summary_response="$(curl -fsS "$base_url/api/webhook-hub/events/summary")"
  list_response="$(curl -fsS "$base_url/api/webhook-hub/events")"
  db_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        count(*) || '|' ||
        count(*) FILTER (WHERE status = 'forwarded') || '|' ||
        count(*) FILTER (WHERE provider = 'stripe') || '|' ||
        (SELECT count(*) FROM webhook_hub.webhook_event_transitions)
      FROM webhook_hub.webhook_events;
    ")"
  echo "[test] webhook-hub create => $create_response"
  echo "[test] webhook-hub duplicate => $duplicate_response"
  echo "[test] webhook-hub validate => $validate_response"
  echo "[test] webhook-hub queue => $queue_response"
  echo "[test] webhook-hub process => $process_response"
  echo "[test] webhook-hub forward => $forward_response"
  echo "[test] webhook-hub transitions => $transitions_response"
  echo "[test] webhook-hub detail => $detail_response"
  echo "[test] webhook-hub filtered list => $filtered_response"
  echo "[test] webhook-hub list after create => $list_response"
  echo "[test] webhook-hub summary => $summary_response"
  echo "[test] webhook-hub db summary => $db_summary"

  if [[ -z "$created_public_id" || "$create_response" != *'"provider":"stripe"'* || "$create_response" != *'"event_type":"payment.succeeded"'* || "$duplicate_response" != *'"code":"webhook_event_conflict"'* || "$duplicate_response" != *'HTTP_STATUS:409'* || "$validate_response" != *'"status":"validated"'* || "$queue_response" != *'"status":"queued"'* || "$process_response" != *'"status":"processing"'* || "$forward_response" != *'"status":"forwarded"'* || "$transitions_response" != *'"status":"received"'* || "$transitions_response" != *'"status":"validated"'* || "$transitions_response" != *'"status":"queued"'* || "$transitions_response" != *'"status":"processing"'* || "$transitions_response" != *'"status":"forwarded"'* || "$detail_response" != *"\"public_id\":\"$created_public_id\""* || "$detail_response" != *'"external_id":"evt_runtime_001"'* || "$detail_response" != *'"status":"forwarded"'* || "$filtered_response" != *'"provider":"stripe"'* || "$filtered_response" != *'"status":"forwarded"'* || "$list_response" != *'"external_id":"evt_runtime_001"'* || "$summary_response" != *'"total":1'* || "$summary_response" != *'"received":0'* || "$summary_response" != *'"pending_delivery":0'* || "$summary_response" != *'"handled":1'* || "$summary_response" != *'"stripe":1'* || "$summary_response" != *'"forwarded":1'* || "$db_summary" != '1|1|1|5' ]]; then
    echo "[test] webhook-hub runtime ingestion did not persist"
    exit 1
  fi
}

run_identity_database_smoke() {
  local smoke_slug="smoke-identity-bootstrap"
  local summary
  local crm_summary
  local sales_summary
  local workflow_control_summary

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

  sales_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        tenant.slug || '|' ||
        (SELECT count(*) FROM sales.opportunities AS opportunity WHERE opportunity.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM sales.proposals AS proposal WHERE proposal.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM sales.opportunities AS opportunity WHERE opportunity.tenant_id = tenant.id AND opportunity.stage = 'won') || '|' ||
        (SELECT count(*) FROM sales.opportunities AS opportunity WHERE opportunity.tenant_id = tenant.id AND opportunity.stage = 'lost') || '|' ||
        (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id AND sale.status = 'active') || '|' ||
        (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id AND sale.status = 'invoiced') || '|' ||
        (SELECT count(*) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id AND invoice.status = 'paid') || '|' ||
        (SELECT count(*) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id AND invoice.status NOT IN ('paid', 'cancelled')) || '|' ||
        (SELECT COALESCE(sum(sale.amount_cents), 0) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id AND sale.status <> 'cancelled') || '|' ||
        (SELECT COALESCE(sum(invoice.amount_cents), 0) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id AND invoice.status = 'paid')
      FROM identity.tenants AS tenant
      WHERE tenant.slug = '$smoke_slug';
    ")"

  echo "[test] sales db smoke => $sales_summary"

  if [[ "$sales_summary" != "$smoke_slug|1|1|1|1|1|0|1|0|0|1|125000|0" ]]; then
    echo "[test] unexpected sales db smoke summary"
    exit 1
  fi

  workflow_control_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        tenant.slug || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definition_versions AS version WHERE version.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_run_events AS workflow_run_event WHERE workflow_run_event.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'draft') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'active') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'archived') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'pending') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'running') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'completed') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'failed') || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'cancelled') || '|' ||
        (SELECT COALESCE(MAX(version.version_number), 0) FROM workflow_control.workflow_definition_versions AS version WHERE version.tenant_id = tenant.id)
      FROM identity.tenants AS tenant
      WHERE tenant.slug = '$smoke_slug';
    ")" 

  echo "[test] workflow-control db smoke => $workflow_control_summary"

  if [[ "$workflow_control_summary" != "$smoke_slug|1|1|1|1|0|1|0|0|1|0|0|0|1" ]]; then
    echo "[test] unexpected workflow-control db smoke summary"
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

wait_for_prometheus_probe_success() {
  local prometheus_url="$1"
  local job="$2"
  local instance="$3"
  local attempts=0
  local response

  until response="$(curl -fsS --get --data-urlencode "query=probe_success{job=\"$job\",instance=\"$instance\"}" "$prometheus_url/api/v1/query")" && [[ "$response" == *'"value":'*'"1"'* ]]; do
    attempts=$((attempts + 1))

    if [[ "$attempts" -ge 45 ]]; then
      echo "[test] prometheus probe did not become successful for $job => $instance"
      exit 1
    fi

    sleep 1
  done
}

compute_totp() {
  local secret="$1"

  if command -v python3 >/dev/null 2>&1; then
    python3 - "$secret" <<'PY'
import base64
import hashlib
import hmac
import sys
import time

secret = sys.argv[1].strip().replace("=", "").upper()
padding = "=" * ((8 - len(secret) % 8) % 8)
key = base64.b32decode(secret + padding, casefold=True)
counter = int(time.time()) // 30
message = counter.to_bytes(8, "big")
digest = hmac.new(key, message, hashlib.sha1).digest()
offset = digest[-1] & 0x0F
code = (
  ((digest[offset] & 0x7F) << 24)
  | (digest[offset + 1] << 16)
  | (digest[offset + 2] << 8)
  | digest[offset + 3]
) % 1000000
print(f"{code:06d}")
PY
    return
  fi

  docker run --rm python:3.12-alpine python - "$secret" <<'PY'
import base64
import hashlib
import hmac
import sys
import time

secret = sys.argv[1].strip().replace("=", "").upper()
padding = "=" * ((8 - len(secret) % 8) % 8)
key = base64.b32decode(secret + padding, casefold=True)
counter = int(time.time()) // 30
message = counter.to_bytes(8, "big")
digest = hmac.new(key, message, hashlib.sha1).digest()
offset = digest[-1] & 0x0F
code = (
  ((digest[offset] & 0x7F) << 24)
  | (digest[offset + 1] << 16)
  | (digest[offset + 2] << 8)
  | digest[offset + 3]
) % 1000000
print(f"{code:06d}")
PY
}

run_crm_runtime_smoke() {
  local base_url="http://localhost:${CRM_HTTP_PORT:-8083}"
  local bootstrap_lead_public_id
  local owner_public_id
  local lead_list
  local bootstrap_notes_response
  local create_response
  local created_public_id
  local created_note_response
  local notes_response
  local profile_response
  local owner_response
  local status_response
  local summary_response

  "${COMPOSE_CMD[@]}" up -d --build crm
  wait_for_http_ready "$base_url/health/ready"

  local details_response
  bootstrap_lead_public_id="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT lead.public_id::text
      FROM crm.leads AS lead
      INNER JOIN identity.tenants AS tenant
        ON tenant.id = lead.tenant_id
      WHERE tenant.slug = 'bootstrap-ops'
        AND lead.email = 'lead@bootstrap-ops.local'
      LIMIT 1;
    ")"
  lead_list="$(curl -fsS "$base_url/api/crm/leads")"
  details_response="$(curl -fsS "$base_url/health/details")"
  bootstrap_notes_response="$(curl -fsS "$base_url/api/crm/leads/$bootstrap_lead_public_id/notes")"
  echo "[test] crm api list => $lead_list"
  echo "[test] crm health details => $details_response"
  echo "[test] crm bootstrap notes => $bootstrap_notes_response"

  if [[ "$lead_list" != *'"email":"lead@bootstrap-ops.local"'* ]]; then
    echo "[test] bootstrap CRM lead was not returned by the live API"
    exit 1
  fi

  if [[ "$details_response" != *'"name":"postgresql","status":"ready"'* ]]; then
    echo "[test] crm health details did not report postgresql ready"
    exit 1
  fi

  if [[ "$bootstrap_notes_response" != *'"category":"qualification"'* ]]; then
    echo "[test] bootstrap CRM note was not returned by the live API"
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

  created_note_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"body":"Cliente pediu comparativo com onboarding premium.","category":"follow-up"}' \
    "$base_url/api/crm/leads/$created_public_id/notes")"
  echo "[test] crm api create note => $created_note_response"

  if [[ "$created_note_response" != *'"category":"follow-up"'* ]]; then
    echo "[test] runtime lead note create did not persist"
    exit 1
  fi

  notes_response="$(curl -fsS "$base_url/api/crm/leads/$created_public_id/notes")"
  echo "[test] crm api list notes => $notes_response"

  if [[ "$notes_response" != *'"body":"Cliente pediu comparativo com onboarding premium."'* ]]; then
    echo "[test] runtime lead notes read did not return the created note"
    exit 1
  fi

  profile_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"name":"Runtime Lead Prime","email":"runtime.lead.prime@example.com","source":"instagram"}' \
    "$base_url/api/crm/leads/$created_public_id")"
  echo "[test] crm api profile => $profile_response"

  if [[ "$profile_response" != *'"name":"Runtime Lead Prime"'* || "$profile_response" != *'"email":"runtime.lead.prime@example.com"'* || "$profile_response" != *'"source":"instagram"'* ]]; then
    echo "[test] runtime lead profile update did not persist"
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

  if [[ "$summary_response" != *'"total":2'* || "$summary_response" != *'"assigned":2'* || "$summary_response" != *'"contacted":1'* || "$summary_response" != *'"instagram":1'* ]]; then
    echo "[test] runtime CRM summary did not reflect live updates"
    exit 1
  fi

  CRM_RUNTIME_LEAD_PUBLIC_ID="$created_public_id"
  CRM_RUNTIME_OWNER_PUBLIC_ID="$owner_public_id"
}

run_sales_runtime_smoke() {
  local base_url="http://localhost:${SALES_HTTP_PORT:-8087}"
  local workflow_control_base_url="http://localhost:${WORKFLOW_CONTROL_HTTP_PORT:-8084}"
  local workflow_runtime_base_url="http://localhost:${WORKFLOW_RUNTIME_HTTP_PORT:-8085}"
  local health_details_response
  local list_response
  local create_opportunity_response
  local created_opportunity_public_id
  local opportunity_detail_response
  local create_proposal_response
  local created_proposal_public_id
  local proposals_response
  local proposal_status_response
  local opportunity_stage_response
  local proposal_detail_response
  local convert_response
  local created_sale_public_id
  local create_invoice_response
  local created_invoice_public_id
  local invoice_status_response
  local sale_detail_response
  local invoice_detail_response
  local invoice_list_response
  local opportunity_summary_response
  local sales_summary_response
  local invoice_summary_response
  local workflow_run_create_response
  local workflow_run_public_id
  local workflow_run_start_response
  local workflow_run_complete_response
  local workflow_execution_create_response
  local workflow_execution_public_id
  local workflow_execution_start_response
  local workflow_execution_complete_response
  local db_summary

  if [[ -z "${CRM_RUNTIME_LEAD_PUBLIC_ID:-}" || -z "${CRM_RUNTIME_OWNER_PUBLIC_ID:-}" ]]; then
    echo "[test] CRM runtime identifiers were not captured before sales smoke"
    exit 1
  fi

  "${COMPOSE_CMD[@]}" up -d --build sales
  wait_for_http_ready "$base_url/health/ready"
  wait_for_http_ready "$workflow_control_base_url/health/ready"
  wait_for_http_ready "$workflow_runtime_base_url/health/ready"

  health_details_response="$(curl -fsS "$base_url/health/details")"
  list_response="$(curl -fsS "$base_url/api/sales/opportunities")"
  echo "[test] sales health details => $health_details_response"
  echo "[test] sales opportunity list => $list_response"

  if [[ "$health_details_response" != *'"name":"postgresql","status":"ready"'* || "$list_response" != *'"title":"Bootstrap Ops Opportunity"'* ]]; then
    echo "[test] sales bootstrap runtime state was not ready"
    exit 1
  fi

  create_opportunity_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"leadPublicId\":\"$CRM_RUNTIME_LEAD_PUBLIC_ID\",\"title\":\"Runtime Expansion Opportunity\",\"ownerUserId\":\"$CRM_RUNTIME_OWNER_PUBLIC_ID\",\"amountCents\":99000}" \
    "$base_url/api/sales/opportunities")"
  created_opportunity_public_id="$(echo "$create_opportunity_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"

  create_proposal_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"title":"Runtime Expansion Proposal","amountCents":99000}' \
    "$base_url/api/sales/opportunities/$created_opportunity_public_id/proposals")"
  created_proposal_public_id="$(echo "$create_proposal_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"

  proposals_response="$(curl -fsS "$base_url/api/sales/opportunities/$created_opportunity_public_id/proposals")"
  proposal_status_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"sent"}' \
    "$base_url/api/sales/proposals/$created_proposal_public_id/status")"
  opportunity_stage_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"stage":"negotiation"}' \
    "$base_url/api/sales/opportunities/$created_opportunity_public_id/stage")"
  convert_response="$(curl -fsS -X POST "$base_url/api/sales/proposals/$created_proposal_public_id/convert")"
  created_sale_public_id="$(echo "$convert_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  create_invoice_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"number":"RUNTIME-INV-0001","dueDate":"2026-05-25"}' \
    "$base_url/api/sales/sales/$created_sale_public_id/invoice")"
  created_invoice_public_id="$(echo "$create_invoice_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  invoice_status_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"paid"}' \
    "$base_url/api/sales/invoices/$created_invoice_public_id/status")"
  opportunity_detail_response="$(curl -fsS "$base_url/api/sales/opportunities/$created_opportunity_public_id")"
  proposal_detail_response="$(curl -fsS "$base_url/api/sales/proposals/$created_proposal_public_id")"
  sale_detail_response="$(curl -fsS "$base_url/api/sales/sales/$created_sale_public_id")"
  invoice_detail_response="$(curl -fsS "$base_url/api/sales/invoices/$created_invoice_public_id")"
  invoice_list_response="$(curl -fsS "$base_url/api/sales/invoices")"
  opportunity_summary_response="$(curl -fsS "$base_url/api/sales/opportunities/summary")"
  sales_summary_response="$(curl -fsS "$base_url/api/sales/sales/summary")"
  invoice_summary_response="$(curl -fsS "$base_url/api/sales/invoices/summary")"

  workflow_run_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"workflowDefinitionKey\":\"lead-follow-up\",\"subjectType\":\"crm.lead\",\"subjectPublicId\":\"$CRM_RUNTIME_LEAD_PUBLIC_ID\",\"initiatedBy\":\"sales-mvp\"}" \
    "$workflow_control_base_url/api/workflow-control/runs")"
  workflow_run_public_id="$(echo "$workflow_run_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  workflow_run_start_response="$(curl -fsS -X POST "$workflow_control_base_url/api/workflow-control/runs/$workflow_run_public_id/start")"
  workflow_run_complete_response="$(curl -fsS -X POST "$workflow_control_base_url/api/workflow-control/runs/$workflow_run_public_id/complete")"

  workflow_execution_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"tenantSlug\":\"bootstrap-ops\",\"workflowDefinitionKey\":\"lead-follow-up\",\"subjectType\":\"crm.lead\",\"subjectPublicId\":\"$CRM_RUNTIME_LEAD_PUBLIC_ID\",\"initiatedBy\":\"sales-mvp\"}" \
    "$workflow_runtime_base_url/api/workflow-runtime/executions")"
  workflow_execution_public_id="$(echo "$workflow_execution_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  workflow_execution_start_response="$(curl -fsS -X POST "$workflow_runtime_base_url/api/workflow-runtime/executions/$workflow_execution_public_id/start")"
  workflow_execution_complete_response="$(curl -fsS -X POST "$workflow_runtime_base_url/api/workflow-runtime/executions/$workflow_execution_public_id/complete")"

  db_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        (SELECT count(*) FROM sales.opportunities AS opportunity INNER JOIN identity.tenants AS tenant ON tenant.id = opportunity.tenant_id WHERE tenant.slug = 'bootstrap-ops') || '|' ||
        (SELECT count(*) FROM sales.proposals AS proposal INNER JOIN identity.tenants AS tenant ON tenant.id = proposal.tenant_id WHERE tenant.slug = 'bootstrap-ops') || '|' ||
        (SELECT count(*) FROM sales.sales AS sale INNER JOIN identity.tenants AS tenant ON tenant.id = sale.tenant_id WHERE tenant.slug = 'bootstrap-ops') || '|' ||
        (SELECT count(*) FROM sales.invoices AS invoice INNER JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id WHERE tenant.slug = 'bootstrap-ops') || '|' ||
        (SELECT count(*) FROM sales.opportunities AS opportunity INNER JOIN identity.tenants AS tenant ON tenant.id = opportunity.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND opportunity.stage = 'won') || '|' ||
        (SELECT count(*) FROM sales.sales AS sale INNER JOIN identity.tenants AS tenant ON tenant.id = sale.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND sale.status = 'active') || '|' ||
        (SELECT count(*) FROM sales.sales AS sale INNER JOIN identity.tenants AS tenant ON tenant.id = sale.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND sale.status = 'invoiced') || '|' ||
        (SELECT count(*) FROM sales.invoices AS invoice INNER JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND invoice.status = 'paid') || '|' ||
        (SELECT count(*) FROM sales.invoices AS invoice INNER JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND invoice.status NOT IN ('paid', 'cancelled')) || '|' ||
        (SELECT COALESCE(sum(sale.amount_cents), 0) FROM sales.sales AS sale INNER JOIN identity.tenants AS tenant ON tenant.id = sale.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND sale.status <> 'cancelled') || '|' ||
        (SELECT COALESCE(sum(invoice.amount_cents), 0) FROM sales.invoices AS invoice INNER JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND invoice.status = 'paid');
    ")"

  echo "[test] sales create opportunity => $create_opportunity_response"
  echo "[test] sales create proposal => $create_proposal_response"
  echo "[test] sales proposal list => $proposals_response"
  echo "[test] sales proposal status => $proposal_status_response"
  echo "[test] sales opportunity stage => $opportunity_stage_response"
  echo "[test] sales convert proposal => $convert_response"
  echo "[test] sales create invoice => $create_invoice_response"
  echo "[test] sales invoice status => $invoice_status_response"
  echo "[test] sales opportunity detail => $opportunity_detail_response"
  echo "[test] sales proposal detail => $proposal_detail_response"
  echo "[test] sales sale detail => $sale_detail_response"
  echo "[test] sales invoice detail => $invoice_detail_response"
  echo "[test] sales invoice list => $invoice_list_response"
  echo "[test] sales opportunity summary => $opportunity_summary_response"
  echo "[test] sales sales summary => $sales_summary_response"
  echo "[test] sales invoice summary => $invoice_summary_response"
  echo "[test] sales workflow-control create => $workflow_run_create_response"
  echo "[test] sales workflow-control start => $workflow_run_start_response"
  echo "[test] sales workflow-control complete => $workflow_run_complete_response"
  echo "[test] sales workflow-runtime create => $workflow_execution_create_response"
  echo "[test] sales workflow-runtime start => $workflow_execution_start_response"
  echo "[test] sales workflow-runtime complete => $workflow_execution_complete_response"
  echo "[test] sales db summary => $db_summary"

  if [[ -z "$created_opportunity_public_id" || -z "$created_proposal_public_id" || -z "$created_sale_public_id" || -z "$created_invoice_public_id" || -z "$workflow_run_public_id" || -z "$workflow_execution_public_id" || "$create_opportunity_response" != *'"stage":"qualified"'* || "$create_opportunity_response" != *"\"leadPublicId\":\"$CRM_RUNTIME_LEAD_PUBLIC_ID\""* || "$create_proposal_response" != *'"status":"draft"'* || "$proposals_response" != *'"title":"Runtime Expansion Proposal"'* || "$proposal_status_response" != *'"status":"sent"'* || "$opportunity_stage_response" != *'"stage":"negotiation"'* || "$convert_response" != *'"status":"active"'* || "$create_invoice_response" != *'"number":"RUNTIME-INV-0001"'* || "$create_invoice_response" != *'"status":"draft"'* || "$invoice_status_response" != *'"status":"paid"'* || "$invoice_status_response" != *'"paidAt":"'*
    || "$opportunity_detail_response" != *'"stage":"won"'* || "$proposal_detail_response" != *'"status":"accepted"'* || "$sale_detail_response" != *'"status":"invoiced"'* || "$invoice_detail_response" != *'"number":"RUNTIME-INV-0001"'* || "$invoice_detail_response" != *'"status":"paid"'* || "$invoice_list_response" != *'"number":"RUNTIME-INV-0001"'* || "$opportunity_summary_response" != *'"total":2'* || "$opportunity_summary_response" != *'"totalAmountCents":224000'* || "$opportunity_summary_response" != *'"won":2'* || "$sales_summary_response" != *'"total":2'* || "$sales_summary_response" != *'"bookedRevenueCents":224000'* || "$sales_summary_response" != *'"active":1'* || "$sales_summary_response" != *'"invoiced":1'* || "$invoice_summary_response" != *'"total":2'* || "$invoice_summary_response" != *'"openAmountCents":125000'* || "$invoice_summary_response" != *'"paidAmountCents":99000'* || "$invoice_summary_response" != *'"overdueCount":0'* || "$invoice_summary_response" != *'"sent":1'* || "$invoice_summary_response" != *'"paid":1'*
    || "$workflow_run_create_response" != *'"status":"pending"'* || "$workflow_run_start_response" != *'"status":"running"'* || "$workflow_run_complete_response" != *'"status":"completed"'* || "$workflow_execution_create_response" != *'"status":"pending"'* || "$workflow_execution_start_response" != *'"status":"running"'* || "$workflow_execution_complete_response" != *'"status":"completed"'* || "$db_summary" != '2|2|2|2|2|1|1|1|1|224000|99000' ]]; then
    echo "[test] sales runtime pipeline did not persist the expected vertical slice"
    exit 1
  fi
}

run_workflow_control_runtime_smoke() {
  local base_url="http://localhost:${WORKFLOW_CONTROL_HTTP_PORT:-8084}"
  local list_response
  local runs_response
  local run_events_response
  local run_events_status_response
  local run_events_creator_response
  local run_events_summary_response
  local run_note_response
  local runs_summary_response
  local create_response
  local created_key="runtime-flow"
  local run_create_response
  local created_run_public_id
  local filtered_runs_response
  local run_start_response
  local run_complete_create_response
  local completed_run_public_id
  local run_complete_response
  local run_fail_create_response
  local failed_run_public_id
  local run_fail_response
  local run_cancel_create_response
  local cancelled_run_public_id
  local run_cancel_response
  local versions_response
  local current_version_response
  local publish_response
  local publish_response_v2
  local version_detail_response
  local restore_response
  local profile_response
  local detail_response
  local status_response
  local health_details_response

  "${COMPOSE_CMD[@]}" up -d --build workflow-control
  wait_for_http_ready "$base_url/health/ready"

  health_details_response="$(curl -fsS "$base_url/health/details")"
  echo "[test] workflow-control health details => $health_details_response"

  if [[ "$health_details_response" != *'"name":"postgresql","status":"ready"'* || "$health_details_response" != *'"name":"definitions-catalog","status":"ready"'* ]]; then
    echo "[test] workflow-control health details did not expose the expected postgresql-backed catalog dependencies"
    exit 1
  fi

  list_response="$(curl -fsS "$base_url/api/workflow-control/definitions")"
  runs_response="$(curl -fsS "$base_url/api/workflow-control/runs")"
  run_events_response="$(curl -fsS "$base_url/api/workflow-control/runs/00000000-0000-0000-0000-000000000301/events")"
  runs_summary_response="$(curl -fsS "$base_url/api/workflow-control/runs/summary")"
  echo "[test] workflow-control list => $list_response"
  echo "[test] workflow-control runs => $runs_response"
  echo "[test] workflow-control run events => $run_events_response"
  echo "[test] workflow-control runs summary => $runs_summary_response"

  if [[ "$list_response" != *'"key":"lead-follow-up"'* || "$list_response" != *'"status":"active"'* ]]; then
    echo "[test] workflow-control bootstrap catalog was not returned by the live API"
    exit 1
  fi

  if [[ "$runs_response" != *'"status":"running"'* || "$runs_response" != *'"triggerEvent":"lead.created"'* || "$run_events_response" != *'"category":"note"'* || "$run_events_response" != *'"createdBy":"bootstrap-seed"'* || "$runs_summary_response" != *'"total":1'* || "$runs_summary_response" != *'"running":1'* ]]; then
    echo "[test] workflow-control bootstrap run ledger was not returned by the live API"
    exit 1
  fi

  versions_response="$(curl -fsS "$base_url/api/workflow-control/definitions/lead-follow-up/versions")"
  current_version_response="$(curl -fsS "$base_url/api/workflow-control/definitions/lead-follow-up/versions/current")"
  echo "[test] workflow-control versions => $versions_response"
  echo "[test] workflow-control current version => $current_version_response"

  if [[ "$versions_response" != *'"versionNumber":1'* || "$current_version_response" != *'"versionNumber":1'* ]]; then
    echo "[test] workflow-control bootstrap version history was not returned by the live API"
    exit 1
  fi

  create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"key":"runtime-flow","name":"Runtime Flow","description":"Fluxo criado no smoke HTTP do workflow-control.","trigger":"lead.created"}' \
    "$base_url/api/workflow-control/definitions")"
  echo "[test] workflow-control create => $create_response"

  if [[ "$create_response" != *"\"key\":\"$created_key\""* || "$create_response" != *'"status":"draft"'* ]]; then
    echo "[test] workflow-control create did not return the expected resource"
    exit 1
  fi

  run_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"workflowDefinitionKey":"lead-follow-up","subjectType":"crm.lead","subjectPublicId":"00000000-0000-0000-0000-000000009998","initiatedBy":"smoke-ops"}' \
    "$base_url/api/workflow-control/runs")"
  echo "[test] workflow-control create run => $run_create_response"

  created_run_public_id="$(echo "$run_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  if [[ -z "$created_run_public_id" || "$run_create_response" != *'"status":"pending"'* || "$run_create_response" != *'"initiatedBy":"smoke-ops"'* ]]; then
    echo "[test] workflow-control run create did not persist"
    exit 1
  fi

  filtered_runs_response="$(curl -fsS "$base_url/api/workflow-control/runs?status=pending&workflowDefinitionKey=lead-follow-up&subjectType=crm.lead&initiatedBy=smoke-ops")"
  runs_summary_response="$(curl -fsS "$base_url/api/workflow-control/runs/summary")"
  echo "[test] workflow-control filtered runs => $filtered_runs_response"
  echo "[test] workflow-control runs summary after create => $runs_summary_response"

  if [[ "$filtered_runs_response" != *"\"publicId\":\"$created_run_public_id\""* || "$runs_summary_response" != *'"total":2'* || "$runs_summary_response" != *'"pending":1'* || "$runs_summary_response" != *'"running":1'* ]]; then
    echo "[test] workflow-control run filters or summary did not reflect live writes"
    exit 1
  fi

  run_start_response="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/runs/$created_run_public_id/start")"
  echo "[test] workflow-control start run => $run_start_response"

  if [[ "$run_start_response" != *"\"publicId\":\"$created_run_public_id\""* || "$run_start_response" != *'"status":"running"'* || "$run_start_response" != *'"startedAt":"'*
  ]]; then
    echo "[test] workflow-control run start did not persist"
    exit 1
  fi

  run_note_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"body":"Anotacao operacional criada no smoke do workflow-control.","createdBy":"smoke-ops"}' \
    "$base_url/api/workflow-control/runs/$created_run_public_id/events")"
  run_events_response="$(curl -fsS "$base_url/api/workflow-control/runs/$created_run_public_id/events")"
  run_events_status_response="$(curl -fsS "$base_url/api/workflow-control/runs/$created_run_public_id/events?category=status")"
  run_events_creator_response="$(curl -fsS "$base_url/api/workflow-control/runs/$created_run_public_id/events?createdBy=smoke-ops")"
  run_events_summary_response="$(curl -fsS "$base_url/api/workflow-control/runs/$created_run_public_id/events/summary")"
  echo "[test] workflow-control create run note => $run_note_response"
  echo "[test] workflow-control runtime run events => $run_events_response"
  echo "[test] workflow-control runtime run events status filter => $run_events_status_response"
  echo "[test] workflow-control runtime run events creator filter => $run_events_creator_response"
  echo "[test] workflow-control runtime run events summary => $run_events_summary_response"

  if [[ "$run_note_response" != *"\"workflowRunPublicId\":\"$created_run_public_id\""* || "$run_note_response" != *'"category":"note"'* || "$run_events_response" != *'"createdBy":"smoke-ops"'* || "$run_events_response" != *'"body":"Anotacao operacional criada no smoke do workflow-control."'* || "$run_events_response" != *'"body":"Workflow run moved to running."'* || "$run_events_response" != *'"createdBy":"workflow-control"'* || "$run_events_status_response" != *'"category":"status"'* || "$run_events_status_response" != *'"createdBy":"workflow-control"'* || "$run_events_creator_response" != *'"category":"note"'* || "$run_events_creator_response" != *'"createdBy":"smoke-ops"'* || "$run_events_summary_response" != *"\"workflowRunPublicId\":\"$created_run_public_id\""* || "$run_events_summary_response" != *'"total":2'* || "$run_events_summary_response" != *'"status":1'* || "$run_events_summary_response" != *'"note":1'* || "$run_events_summary_response" != *'"latestCategory":"note"'* ]]; then
    echo "[test] workflow-control run note create did not persist"
    exit 1
  fi

  run_complete_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"workflowDefinitionKey":"lead-follow-up","subjectType":"crm.lead","subjectPublicId":"00000000-0000-0000-0000-000000009997","initiatedBy":"smoke-complete"}' \
    "$base_url/api/workflow-control/runs")"
  completed_run_public_id="$(echo "$run_complete_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"

  if [[ -z "$completed_run_public_id" ]]; then
    echo "[test] workflow-control complete-path run create did not return public id"
    exit 1
  fi

  curl -fsS -X POST "$base_url/api/workflow-control/runs/$completed_run_public_id/start" >/dev/null
  run_complete_response="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/runs/$completed_run_public_id/complete")"
  echo "[test] workflow-control complete run => $run_complete_response"

  if [[ "$run_complete_response" != *"\"publicId\":\"$completed_run_public_id\""* || "$run_complete_response" != *'"status":"completed"'* || "$run_complete_response" != *'"completedAt":"'*
  ]]; then
    echo "[test] workflow-control run complete did not persist"
    exit 1
  fi

  run_fail_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"workflowDefinitionKey":"lead-follow-up","subjectType":"crm.lead","subjectPublicId":"00000000-0000-0000-0000-000000009996","initiatedBy":"smoke-fail"}' \
    "$base_url/api/workflow-control/runs")"
  failed_run_public_id="$(echo "$run_fail_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"

  if [[ -z "$failed_run_public_id" ]]; then
    echo "[test] workflow-control fail-path run create did not return public id"
    exit 1
  fi

  curl -fsS -X POST "$base_url/api/workflow-control/runs/$failed_run_public_id/start" >/dev/null
  run_fail_response="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/runs/$failed_run_public_id/fail")"
  echo "[test] workflow-control fail run => $run_fail_response"

  if [[ "$run_fail_response" != *"\"publicId\":\"$failed_run_public_id\""* || "$run_fail_response" != *'"status":"failed"'* || "$run_fail_response" != *'"failedAt":"'*
  ]]; then
    echo "[test] workflow-control run fail did not persist"
    exit 1
  fi

  run_cancel_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"workflowDefinitionKey":"lead-follow-up","subjectType":"crm.lead","subjectPublicId":"00000000-0000-0000-0000-000000009995","initiatedBy":"smoke-cancel"}' \
    "$base_url/api/workflow-control/runs")"
  cancelled_run_public_id="$(echo "$run_cancel_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"

  if [[ -z "$cancelled_run_public_id" ]]; then
    echo "[test] workflow-control cancel-path run create did not return public id"
    exit 1
  fi

  run_cancel_response="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/runs/$cancelled_run_public_id/cancel")"
  echo "[test] workflow-control cancel run => $run_cancel_response"

  if [[ "$run_cancel_response" != *"\"publicId\":\"$cancelled_run_public_id\""* || "$run_cancel_response" != *'"status":"cancelled"'* || "$run_cancel_response" != *'"cancelledAt":"'*
  ]]; then
    echo "[test] workflow-control run cancel did not persist"
    exit 1
  fi

  runs_summary_response="$(curl -fsS "$base_url/api/workflow-control/runs/summary")"
  echo "[test] workflow-control runs summary after transitions => $runs_summary_response"

  if [[ "$runs_summary_response" != *'"total":5'* || "$runs_summary_response" != *'"running":2'* || "$runs_summary_response" != *'"completed":1'* || "$runs_summary_response" != *'"failed":1'* || "$runs_summary_response" != *'"cancelled":1'* ]]; then
    echo "[test] workflow-control run transitions did not reflect on live summary"
    exit 1
  fi

  publish_response="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/definitions/$created_key/versions")"
  echo "[test] workflow-control publish => $publish_response"

  if [[ "$publish_response" != *'"versionNumber":1'* || "$publish_response" != *'"snapshotName":"Runtime Flow"'* ]]; then
    echo "[test] workflow-control version publish did not persist"
    exit 1
  fi

  profile_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"name":"Runtime Flow Prime","description":"Fluxo atualizado no smoke do workflow-control.","trigger":"lead.qualified"}' \
    "$base_url/api/workflow-control/definitions/$created_key")"
  echo "[test] workflow-control profile => $profile_response"

  if [[ "$profile_response" != *"\"key\":\"$created_key\""* || "$profile_response" != *'"name":"Runtime Flow Prime"'* || "$profile_response" != *'"trigger":"lead.qualified"'* ]]; then
    echo "[test] workflow-control metadata update did not persist"
    exit 1
  fi

  detail_response="$(curl -fsS "$base_url/api/workflow-control/definitions/$created_key")"
  echo "[test] workflow-control detail => $detail_response"

  if [[ "$detail_response" != *"\"key\":\"$created_key\""* || "$detail_response" != *'"name":"Runtime Flow Prime"'* || "$detail_response" != *'"trigger":"lead.qualified"'* ]]; then
    echo "[test] workflow-control detail did not return the created resource"
    exit 1
  fi

  status_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"active"}' \
    "$base_url/api/workflow-control/definitions/$created_key/status")"
  echo "[test] workflow-control status => $status_response"

  if [[ "$status_response" != *"\"key\":\"$created_key\""* || "$status_response" != *'"status":"active"'* ]]; then
    echo "[test] workflow-control status update did not persist"
    exit 1
  fi

  publish_response_v2="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/definitions/$created_key/versions")"
  current_version_response="$(curl -fsS "$base_url/api/workflow-control/definitions/$created_key/versions/current")"
  version_detail_response="$(curl -fsS "$base_url/api/workflow-control/definitions/$created_key/versions/1")"
  echo "[test] workflow-control publish v2 => $publish_response_v2"
  echo "[test] workflow-control current version runtime-flow => $current_version_response"
  echo "[test] workflow-control version detail runtime-flow => $version_detail_response"

  if [[ "$publish_response_v2" != *'"versionNumber":2'* || "$current_version_response" != *'"versionNumber":2'* || "$version_detail_response" != *'"versionNumber":1'* || "$version_detail_response" != *'"snapshotTrigger":"lead.created"'* ]]; then
    echo "[test] workflow-control version history did not reflect live publish operations"
    exit 1
  fi

  restore_response="$(curl -fsS \
    -X POST \
    "$base_url/api/workflow-control/definitions/$created_key/versions/1/restore")"
  echo "[test] workflow-control restore => $restore_response"

  if [[ "$restore_response" != *'"name":"Runtime Flow"'* || "$restore_response" != *'"status":"draft"'* || "$restore_response" != *'"trigger":"lead.created"'* ]]; then
    echo "[test] workflow-control restore did not bring the definition back to the published snapshot"
    exit 1
  fi
}

run_workflow_runtime_smoke() {
  local base_url="http://localhost:${WORKFLOW_RUNTIME_HTTP_PORT:-8085}"
  local health_details_response
  local list_response
  local summary_response
  local create_response
  local created_public_id
  local start_response
  local complete_response
  local transitions_response
  local filtered_response
  local filtered_summary_response
  local cancel_create_response
  local cancel_public_id
  local cancel_response
  local fail_create_response
  local fail_public_id
  local fail_response
  local retry_response
  local retry_start_response
  local retry_complete_response
  local retry_transitions_response
  local grouped_summary_response
  local db_summary

  "${COMPOSE_CMD[@]}" up -d --build workflow-runtime
  wait_for_http_ready "$base_url/health/ready"

  health_details_response="$(curl -fsS "$base_url/health/details")"
  list_response="$(curl -fsS "$base_url/api/workflow-runtime/executions")"
  summary_response="$(curl -fsS "$base_url/api/workflow-runtime/executions/summary")"
  echo "[test] workflow-runtime health details => $health_details_response"
  echo "[test] workflow-runtime list => $list_response"
  echo "[test] workflow-runtime summary => $summary_response"

  if [[ "$health_details_response" != *'"name":"postgresql","status":"ready"'* || "$health_details_response" != *'"name":"timer-wheel","status":"ready"'* || "$health_details_response" != *'"name":"workflow-catalog","status":"ready"'* || "$list_response" != '[]' || "$summary_response" != *'"total":0'* ]]; then
    echo "[test] workflow-runtime bootstrap runtime state was not clean"
    exit 1
  fi

  create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"tenantSlug":"bootstrap-ops","workflowDefinitionKey":"lead-follow-up","subjectType":"crm.lead","subjectPublicId":"00000000-0000-0000-0000-000000008851","initiatedBy":"runtime-smoke"}' \
    "$base_url/api/workflow-runtime/executions")"
  created_public_id="$(echo "$create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"

  start_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$created_public_id/start")"
  complete_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$created_public_id/complete")"
  transitions_response="$(curl -fsS "$base_url/api/workflow-runtime/executions/$created_public_id/transitions")"
  filtered_response="$(curl -fsS "$base_url/api/workflow-runtime/executions?tenantSlug=bootstrap-ops&status=completed")"
  filtered_summary_response="$(curl -fsS "$base_url/api/workflow-runtime/executions/summary?tenantSlug=bootstrap-ops&workflowDefinitionKey=lead-follow-up")"

  cancel_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"tenantSlug":"northwind-group","workflowDefinitionKey":"lead-follow-up","subjectType":"crm.deal","subjectPublicId":"00000000-0000-0000-0000-000000008852","initiatedBy":"runtime-cancel"}' \
    "$base_url/api/workflow-runtime/executions")"
  cancel_public_id="$(echo "$cancel_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  cancel_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$cancel_public_id/cancel")"

  fail_create_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"workflowDefinitionKey":"lead-follow-up","subjectType":"sales.quote","subjectPublicId":"00000000-0000-0000-0000-000000008853","initiatedBy":"runtime-fail"}' \
    "$base_url/api/workflow-runtime/executions")"
  fail_public_id="$(echo "$fail_create_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  fail_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$fail_public_id/fail")"
  retry_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$fail_public_id/retry")"
  retry_start_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$fail_public_id/start")"
  retry_complete_response="$(curl -fsS -X POST "$base_url/api/workflow-runtime/executions/$fail_public_id/complete")"
  retry_transitions_response="$(curl -fsS "$base_url/api/workflow-runtime/executions/$fail_public_id/transitions")"
  grouped_summary_response="$(curl -fsS "$base_url/api/workflow-runtime/executions/summary/by-workflow?tenantSlug=bootstrap-ops")"
  summary_response="$(curl -fsS "$base_url/api/workflow-runtime/executions/summary")"
  list_response="$(curl -fsS "$base_url/api/workflow-runtime/executions")"
  db_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        count(*) || '|' ||
        count(*) FILTER (WHERE status = 'completed') || '|' ||
        count(*) FILTER (WHERE status = 'failed') || '|' ||
        count(*) FILTER (WHERE status = 'cancelled') || '|' ||
        (SELECT count(*) FROM workflow_runtime.execution_transitions) || '|' ||
        coalesce(sum(retry_count), 0)
      FROM workflow_runtime.executions;
    ")"

  echo "[test] workflow-runtime create => $create_response"
  echo "[test] workflow-runtime start => $start_response"
  echo "[test] workflow-runtime complete => $complete_response"
  echo "[test] workflow-runtime transitions => $transitions_response"
  echo "[test] workflow-runtime filtered list => $filtered_response"
  echo "[test] workflow-runtime filtered summary => $filtered_summary_response"
  echo "[test] workflow-runtime cancel => $cancel_response"
  echo "[test] workflow-runtime fail => $fail_response"
  echo "[test] workflow-runtime retry => $retry_response"
  echo "[test] workflow-runtime retry start => $retry_start_response"
  echo "[test] workflow-runtime retry complete => $retry_complete_response"
  echo "[test] workflow-runtime retry transitions => $retry_transitions_response"
  echo "[test] workflow-runtime grouped summary => $grouped_summary_response"
  echo "[test] workflow-runtime summary after transitions => $summary_response"
  echo "[test] workflow-runtime list after transitions => $list_response"
  echo "[test] workflow-runtime db summary => $db_summary"

  if [[ -z "$created_public_id" || -z "$cancel_public_id" || -z "$fail_public_id" || "$create_response" != *'"tenantSlug":"bootstrap-ops"'* || "$start_response" != *'"status":"running"'* || "$complete_response" != *'"status":"completed"'* || "$complete_response" != *'"completedAt":"'*
    || "$transitions_response" != *'"status":"pending"'* || "$transitions_response" != *'"status":"running"'* || "$transitions_response" != *'"status":"completed"'* || "$filtered_response" != *"\"publicId\":\"$created_public_id\""* || "$filtered_response" != *'"tenantSlug":"bootstrap-ops"'* || "$filtered_summary_response" != *'"total":1'* || "$filtered_summary_response" != *'"completed":1'* || "$filtered_summary_response" != *'"failed":0'* || "$cancel_response" != *'"status":"cancelled"'* || "$cancel_response" != *'"tenantSlug":"northwind-group"'* || "$cancel_response" != *'"cancelledAt":"'*
    || "$fail_response" != *'"status":"failed"'* || "$fail_response" != *'"failedAt":"'* || "$retry_response" != *'"status":"pending"'* || "$retry_response" != *'"retryCount":1'* || "$retry_start_response" != *'"status":"running"'* || "$retry_complete_response" != *'"status":"completed"'* || "$retry_complete_response" != *'"retryCount":1'* || "$retry_transitions_response" != *'"status":"failed"'* || "$retry_transitions_response" != *'"status":"completed"'* || "$grouped_summary_response" != *'"workflowDefinitionKey":"lead-follow-up"'* || "$grouped_summary_response" != *'"retriesTotal":1'*
    || "$summary_response" != *'"total":3'* || "$summary_response" != *'"completed":2'* || "$summary_response" != *'"failed":0'* || "$summary_response" != *'"cancelled":1'* || "$list_response" != *"\"publicId\":\"$cancel_public_id\""* || "$list_response" != *"\"publicId\":\"$fail_public_id\""* || "$db_summary" != '3|2|0|1|10|1' ]]; then
    echo "[test] workflow-runtime runtime lifecycle did not persist in postgresql as expected"
    exit 1
  fi
}

run_engagement_runtime_smoke() {
  local base_url="http://localhost:${ENGAGEMENT_HTTP_PORT:-8088}"
  local health_details_response
  local campaigns_response
  local bootstrap_summary_response
  local create_campaign_response
  local created_campaign_public_id
  local update_campaign_status_response
  local create_touchpoint_response
  local created_touchpoint_public_id
  local update_touchpoint_status_response
  local touchpoint_detail_response
  local filtered_touchpoints_response
  local summary_response
  local db_summary

  if [[ -z "${CRM_RUNTIME_LEAD_PUBLIC_ID:-}" ]]; then
    echo "[test] CRM runtime identifiers were not captured before engagement smoke"
    exit 1
  fi

  "${COMPOSE_CMD[@]}" up -d --build engagement
  wait_for_http_ready "$base_url/health/ready"

  health_details_response="$(curl -fsS "$base_url/health/details")"
  campaigns_response="$(curl -fsS "$base_url/api/engagement/campaigns?tenantSlug=bootstrap-ops")"
  bootstrap_summary_response="$(curl -fsS "$base_url/api/engagement/touchpoints/summary?tenantSlug=bootstrap-ops")"

  echo "[test] engagement health details => $health_details_response"
  echo "[test] engagement campaigns => $campaigns_response"
  echo "[test] engagement bootstrap summary => $bootstrap_summary_response"

  if [[ "$health_details_response" != *'"name":"postgresql","status":"ready"'* || "$campaigns_response" != *'"key":"lead-follow-up-campaign"'* || "$bootstrap_summary_response" != *'"workflowDispatched":1'* || "$bootstrap_summary_response" != *'"responded":1'* ]]; then
    echo "[test] engagement bootstrap runtime state was not ready"
    exit 1
  fi

  create_campaign_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"tenantSlug":"bootstrap-ops","key":"runtime-reactivation","name":"Runtime Reactivation","description":"Campanha criada no smoke do engagement.","channel":"email","touchpointGoal":"revive-lead","workflowDefinitionKey":"lead-follow-up","budgetCents":48000}' \
    "$base_url/api/engagement/campaigns")"
  created_campaign_public_id="$(echo "$create_campaign_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  update_campaign_status_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"active"}' \
    "$base_url/api/engagement/campaigns/$created_campaign_public_id/status")"
  create_touchpoint_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"tenantSlug\":\"bootstrap-ops\",\"campaignPublicId\":\"$created_campaign_public_id\",\"leadPublicId\":\"$CRM_RUNTIME_LEAD_PUBLIC_ID\",\"contactValue\":\"runtime.lead.prime@example.com\",\"source\":\"sales\",\"createdBy\":\"engagement-smoke\",\"notes\":\"Touchpoint criado no smoke do engagement.\"}" \
    "$base_url/api/engagement/touchpoints")"
  created_touchpoint_public_id="$(echo "$create_touchpoint_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  update_touchpoint_status_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"converted","lastWorkflowRunPublicId":"00000000-0000-0000-0000-000000000390"}' \
    "$base_url/api/engagement/touchpoints/$created_touchpoint_public_id/status")"
  touchpoint_detail_response="$(curl -fsS "$base_url/api/engagement/touchpoints/$created_touchpoint_public_id")"
  filtered_touchpoints_response="$(curl -fsS "$base_url/api/engagement/touchpoints?tenantSlug=bootstrap-ops&status=converted&channel=email")"
  summary_response="$(curl -fsS "$base_url/api/engagement/touchpoints/summary?tenantSlug=bootstrap-ops")"
  db_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        (SELECT count(*) FROM engagement.campaigns AS campaign INNER JOIN identity.tenants AS tenant ON tenant.id = campaign.tenant_id WHERE tenant.slug = 'bootstrap-ops') || '|' ||
        (SELECT count(*) FROM engagement.touchpoints AS touchpoint INNER JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id WHERE tenant.slug = 'bootstrap-ops') || '|' ||
        (SELECT count(*) FROM engagement.campaigns AS campaign INNER JOIN identity.tenants AS tenant ON tenant.id = campaign.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND campaign.status = 'active') || '|' ||
        (SELECT count(*) FROM engagement.touchpoints AS touchpoint INNER JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND touchpoint.status = 'responded') || '|' ||
        (SELECT count(*) FROM engagement.touchpoints AS touchpoint INNER JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND touchpoint.status = 'converted') || '|' ||
        (SELECT count(*) FROM engagement.touchpoints AS touchpoint INNER JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id WHERE tenant.slug = 'bootstrap-ops' AND touchpoint.last_workflow_run_public_id IS NOT NULL);
    ")"

  echo "[test] engagement create campaign => $create_campaign_response"
  echo "[test] engagement update campaign status => $update_campaign_status_response"
  echo "[test] engagement create touchpoint => $create_touchpoint_response"
  echo "[test] engagement update touchpoint status => $update_touchpoint_status_response"
  echo "[test] engagement touchpoint detail => $touchpoint_detail_response"
  echo "[test] engagement filtered touchpoints => $filtered_touchpoints_response"
  echo "[test] engagement summary => $summary_response"
  echo "[test] engagement db summary => $db_summary"

  if [[ -z "$created_campaign_public_id" || -z "$created_touchpoint_public_id" || "$update_campaign_status_response" != *'"status":"active"'* || "$create_touchpoint_response" != *"\"leadPublicId\":\"$CRM_RUNTIME_LEAD_PUBLIC_ID\""* || "$create_touchpoint_response" != *'"status":"queued"'* || "$update_touchpoint_status_response" != *'"status":"converted"'* || "$update_touchpoint_status_response" != *'"lastWorkflowRunPublicId":"00000000-0000-0000-0000-000000000390"'* || "$touchpoint_detail_response" != *'"channel":"email"'* || "$filtered_touchpoints_response" != *"\"publicId\":\"$created_touchpoint_public_id\""* || "$summary_response" != *'"touchpoints":2'* || "$summary_response" != *'"workflowDispatched":2'* || "$summary_response" != *'"converted":1'* || "$summary_response" != *'"email":1'* || "$db_summary" != "3|2|2|1|1|2" ]]; then
    echo "[test] engagement runtime flow did not persist the expected omnichannel slice"
    exit 1
  fi
}

run_analytics_runtime_smoke() {
  local base_url="http://localhost:${ANALYTICS_HTTP_PORT:-8086}"
  local health_details_response
  local pipeline_summary_response
  local service_pulse_response
  local sales_journey_response
  local tenant_360_response
  local automation_board_response
  local workflow_definition_health_response
  local delivery_reliability_response
  local revenue_operations_response

  "${COMPOSE_CMD[@]}" up -d --build analytics
  wait_for_http_ready "$base_url/health/ready"

  health_details_response="$(curl -fsS "$base_url/health/details")"
  pipeline_summary_response="$(curl -fsS "$base_url/api/analytics/reports/pipeline-summary?tenant_slug=bootstrap-ops")"
  service_pulse_response="$(curl -fsS "$base_url/api/analytics/reports/service-pulse?tenant_slug=bootstrap-ops")"
  sales_journey_response="$(curl -fsS "$base_url/api/analytics/reports/sales-journey?tenant_slug=bootstrap-ops")"
  tenant_360_response="$(curl -fsS "$base_url/api/analytics/reports/tenant-360?tenant_slug=bootstrap-ops")"
  automation_board_response="$(curl -fsS "$base_url/api/analytics/reports/automation-board?tenant_slug=bootstrap-ops")"
  workflow_definition_health_response="$(curl -fsS "$base_url/api/analytics/reports/workflow-definition-health?tenant_slug=bootstrap-ops")"
  delivery_reliability_response="$(curl -fsS "$base_url/api/analytics/reports/delivery-reliability?provider=stripe")"
  revenue_operations_response="$(curl -fsS "$base_url/api/analytics/reports/revenue-operations?tenant_slug=bootstrap-ops")"
  echo "[test] analytics health details => $health_details_response"
  echo "[test] analytics pipeline summary => $pipeline_summary_response"
  echo "[test] analytics service pulse => $service_pulse_response"
  echo "[test] analytics sales journey => $sales_journey_response"
  echo "[test] analytics tenant 360 => $tenant_360_response"
  echo "[test] analytics automation board => $automation_board_response"
  echo "[test] analytics workflow definition health => $workflow_definition_health_response"
  echo "[test] analytics delivery reliability => $delivery_reliability_response"
  echo "[test] analytics revenue operations => $revenue_operations_response"

  if [[ "$health_details_response" != *'"name":"report-engine","status":"ready"'* || "$health_details_response" != *'"name":"postgresql","status":"ready"'* || "$pipeline_summary_response" != *'"tenantSlug":"bootstrap-ops"'* || "$pipeline_summary_response" != *'"dataSource":"postgresql"'* || "$pipeline_summary_response" != *'"leadsCaptured":2'* || "$pipeline_summary_response" != *'"conversions":2'* || "$pipeline_summary_response" != *'"manual":1'* || "$pipeline_summary_response" != *'"instagram":1'* || "$pipeline_summary_response" != *'"runningAutomations":2'* || "$service_pulse_response" != *'"tenantSlug":"bootstrap-ops"'* || "$service_pulse_response" != *'"dataSource":"postgresql"'* || "$service_pulse_response" != *'"totalLeads":2'* || "$service_pulse_response" != *'"salesTotal":2'* || "$service_pulse_response" != *'"bookedRevenueCents":224000'* || "$service_pulse_response" != *'"activeDefinitions":1'* || "$service_pulse_response" != *'"runsRunning":2'* || "$service_pulse_response" != *'"runsCompleted":2'* || "$service_pulse_response" != *'"runsFailed":1'* || "$service_pulse_response" != *'"runsCancelled":1'* || "$service_pulse_response" != *'"totalExecutions":3'* || "$service_pulse_response" != *'"completed":3'* || "$service_pulse_response" != *'"failed":0'* || "$service_pulse_response" != *'"forwarded":1'* || "$sales_journey_response" != *'"tenantSlug":"bootstrap-ops"'* || "$sales_journey_response" != *'"dataSource":"postgresql"'* || "$sales_journey_response" != *'"leadsCaptured":2'* || "$sales_journey_response" != *'"leadsWithOpportunity":2'* || "$sales_journey_response" != *'"opportunitiesWithProposal":2'* || "$sales_journey_response" != *'"proposalsConverted":2'* || "$sales_journey_response" != *'"salesWon":2'* || "$sales_journey_response" != *'"leadToSaleConversionRate":1.0'* || "$sales_journey_response" != *'"totalAmountCents":224000'* || "$sales_journey_response" != *'"won":2'* || "$sales_journey_response" != *'"accepted":2'* || "$sales_journey_response" != *'"bookedRevenueCents":224000'* || "$sales_journey_response" != *'"active":1'* || "$sales_journey_response" != *'"invoiced":1'* || "$sales_journey_response" != *'"controlRuns":1'* || "$sales_journey_response" != *'"runtimeCompleted":1'* || "$tenant_360_response" != *'"tenantSlug":"bootstrap-ops"'* || "$tenant_360_response" != *'"dataSource":"postgresql"'* || "$tenant_360_response" != *'"companies":1'* || "$tenant_360_response" != *'"users":1'* || "$tenant_360_response" != *'"teams":1'* || "$tenant_360_response" != *'"roles":5'* || "$tenant_360_response" != *'"assignedLeads":2'* || "$tenant_360_response" != *'"leadNotes":2'* || "$tenant_360_response" != *'"opportunities":2'* || "$tenant_360_response" != *'"proposals":2'* || "$tenant_360_response" != *'"sales":2'* || "$tenant_360_response" != *'"bookedRevenueCents":224000'* || "$tenant_360_response" != *'"workflowRuns":6'* || "$tenant_360_response" != *'"workflowRunEvents":10'* || "$tenant_360_response" != *'"runtimeExecutions":3'* || "$tenant_360_response" != *'"runtimeCompleted":3'* || "$tenant_360_response" != *'"runtimeFailed":0'* || "$automation_board_response" != *'"tenantSlug":"bootstrap-ops"'* || "$automation_board_response" != *'"dataSource":"postgresql"'* || "$automation_board_response" != *'"definitionsTotal":2'* || "$automation_board_response" != *'"definitionsActive":1'* || "$automation_board_response" != *'"definitionsDraft":1'* || "$automation_board_response" != *'"publishedVersions":3'* || "$automation_board_response" != *'"runsTotal":6'* || "$automation_board_response" != *'"runningRuns":2'* || "$automation_board_response" != *'"recordedEvents":10'* || "$automation_board_response" != *'"workflowDefinitionKey":"lead-follow-up"'* || "$automation_board_response" != *'"executionsTotal":3'* || "$automation_board_response" != *'"completedExecutions":3'* || "$automation_board_response" != *'"failedExecutions":0'* || "$automation_board_response" != *'"retriesTotal":1'* || "$automation_board_response" != *'"recordedTransitions":11'* || "$automation_board_response" != *'"forwarded":1'* || "$workflow_definition_health_response" != *'"tenantSlug":"bootstrap-ops"'* || "$workflow_definition_health_response" != *'"dataSource":"postgresql"'* || "$workflow_definition_health_response" != *'"definitionsTotal":2'* || "$workflow_definition_health_response" != *'"stable":1'* || "$workflow_definition_health_response" != *'"attention":1'* || "$workflow_definition_health_response" != *'"critical":0'* || "$workflow_definition_health_response" != *'"workflowDefinitionKey":"lead-follow-up"'* || "$workflow_definition_health_response" != *'"health":"stable"'* || "$workflow_definition_health_response" != *'"workflowDefinitionKey":"runtime-flow"'* || "$workflow_definition_health_response" != *'"health":"attention"'* || "$workflow_definition_health_response" != *'"definition-not-active"'* || "$delivery_reliability_response" != *'"provider":"stripe"'* || "$delivery_reliability_response" != *'"dataSource":"postgresql"'* || "$delivery_reliability_response" != *'"totalEvents":1'* || "$delivery_reliability_response" != *'"handledEvents":1'* || "$delivery_reliability_response" != *'"avgTransitionsPerEvent":5.0'* || "$delivery_reliability_response" != *'"received":1'* || "$delivery_reliability_response" != *'"validated":1'* || "$delivery_reliability_response" != *'"queued":1'* || "$delivery_reliability_response" != *'"processing":1'* || "$delivery_reliability_response" != *'"forwarded":1'* || "$revenue_operations_response" != *'"tenantSlug":"bootstrap-ops"'* || "$revenue_operations_response" != *'"dataSource":"postgresql"'* || "$revenue_operations_response" != *'"bookedRevenueCents":224000'* || "$revenue_operations_response" != *'"total":2'* || "$revenue_operations_response" != *'"openAmountCents":125000'* || "$revenue_operations_response" != *'"paidAmountCents":99000'* || "$revenue_operations_response" != *'"overdueAmountCents":0'* || "$revenue_operations_response" != *'"overdueCount":0'* || "$revenue_operations_response" != *'"invoiceCoverageRate":1.0'* || "$revenue_operations_response" != *'"collectionRate":0.442'* || "$revenue_operations_response" != *'"averageTicketCents":112000'* || "$revenue_operations_response" != *'"invoicesDueSoon":0'* ]]; then
    echo "[test] analytics runtime bootstrap report did not return the expected payload"
    exit 1
  fi
}

run_identity_runtime_smoke() {
  local base_url="http://localhost:${IDENTITY_HTTP_PORT:-8081}"
  local tenant_slug="runtime-identity-lab"
  local tenants_response
  local create_tenant_response
  local companies_response
  local create_company_response
  local updated_company_response
  local company_detail_response
  local users_response
  local create_user_response
  local updated_user_response
  local user_detail_response
  local created_user_public_id
  local teams_response
  local create_team_response
  local updated_team_response
  local created_team_public_id
  local members_response
  local removed_member_response
  local assign_role_response
  local revoke_role_response
  local user_roles_response
  local roles_response
  local snapshot_response
  local create_invite_response
  local invites_response
  local invite_token
  local accept_invite_response
  local invited_user_public_id
  local mfa_enroll_response
  local mfa_secret
  local mfa_code
  local mfa_verify_response
  local login_without_otp_response
  local login_with_otp_response
  local session_token
  local refresh_token
  local access_response
  local refresh_session_response
  local blocked_access_response
  local audit_response

  "${COMPOSE_CMD[@]}" up -d --build identity
  wait_for_http_ready "$base_url/health/ready"

  local details_response
  tenants_response="$(curl -fsS "$base_url/api/identity/tenants")"
  details_response="$(curl -fsS "$base_url/health/details")"
  echo "[test] identity api tenants => $tenants_response"
  echo "[test] identity health details => $details_response"

  if [[ "$tenants_response" != *'"slug":"bootstrap-ops"'* ]]; then
    echo "[test] bootstrap identity tenants were not returned by the live API"
    exit 1
  fi

  if [[ "$details_response" != *'"name":"postgresql","status":"ready"'* || "$details_response" != *'"name":"keycloak","status":"ready"'* || "$details_response" != *'"name":"openfga","status":"ready"'* ]]; then
    echo "[test] identity health details did not report postgresql ready"
    exit 1
  fi

  create_tenant_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"slug\":\"$tenant_slug\",\"displayName\":\"Runtime Identity Lab\"}" \
    "$base_url/api/identity/tenants")"
  echo "[test] identity api create tenant => $create_tenant_response"

  if [[ "$create_tenant_response" != *"\"slug\":\"$tenant_slug\""* ]]; then
    echo "[test] runtime identity tenant was not created"
    exit 1
  fi

  companies_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/companies")"
  echo "[test] identity api companies => $companies_response"

  if [[ "$companies_response" != *'"displayName":"Runtime Identity Lab"'* ]]; then
    echo "[test] default runtime identity company was not provisioned"
    exit 1
  fi

  create_company_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"displayName":"Runtime Identity Branch","legalName":"Runtime Identity Branch LTDA","taxId":"12345678901234"}' \
    "$base_url/api/identity/tenants/$tenant_slug/companies")"
  echo "[test] identity api create company => $create_company_response"

  if [[ "$create_company_response" != *'"displayName":"Runtime Identity Branch"'* ]]; then
    echo "[test] runtime identity company create did not persist"
    exit 1
  fi

  local created_company_public_id
  created_company_public_id="$(echo "$create_company_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  if [[ -z "$created_company_public_id" ]]; then
    echo "[test] runtime identity company public id was not returned"
    exit 1
  fi

  updated_company_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"displayName":"Runtime Identity Branch Prime","legalName":"Runtime Identity Branch Prime LTDA","taxId":"55544433322211"}' \
    "$base_url/api/identity/tenants/$tenant_slug/companies/$created_company_public_id")"
  echo "[test] identity api update company => $updated_company_response"

  if [[ "$updated_company_response" != *'"displayName":"Runtime Identity Branch Prime"'* || "$updated_company_response" != *'"taxId":"55544433322211"'* ]]; then
    echo "[test] runtime identity company update did not persist"
    exit 1
  fi

  company_detail_response="$(curl -fsS \
    "$base_url/api/identity/tenants/$tenant_slug/companies/$created_company_public_id")"
  echo "[test] identity api company detail => $company_detail_response"

  if [[ "$company_detail_response" != *'"displayName":"Runtime Identity Branch Prime"'* ]]; then
    echo "[test] runtime identity company detail did not reflect live update"
    exit 1
  fi

  users_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/users")"
  echo "[test] identity api users => $users_response"

  if [[ "$users_response" != *"\"email\":\"owner@$tenant_slug.local\""* ]]; then
    echo "[test] default runtime identity owner was not provisioned"
    exit 1
  fi

  create_user_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"runtime.user@$tenant_slug.local\",\"displayName\":\"Runtime User\",\"givenName\":\"Runtime\",\"familyName\":\"User\"}" \
    "$base_url/api/identity/tenants/$tenant_slug/users")"
  echo "[test] identity api create user => $create_user_response"

  created_user_public_id="$(echo "$create_user_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  if [[ -z "$created_user_public_id" ]]; then
    echo "[test] runtime identity user public id was not returned"
    exit 1
  fi

  updated_user_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"runtime.user.prime@$tenant_slug.local\",\"displayName\":\"Runtime User Prime\",\"givenName\":\"Runtime\",\"familyName\":\"Prime\"}" \
    "$base_url/api/identity/tenants/$tenant_slug/users/$created_user_public_id")"
  echo "[test] identity api update user => $updated_user_response"

  if [[ "$updated_user_response" != *"\"email\":\"runtime.user.prime@$tenant_slug.local\""* || "$updated_user_response" != *'"familyName":"Prime"'* ]]; then
    echo "[test] runtime identity user update did not persist"
    exit 1
  fi

  user_detail_response="$(curl -fsS \
    "$base_url/api/identity/tenants/$tenant_slug/users/$created_user_public_id")"
  echo "[test] identity api user detail => $user_detail_response"

  if [[ "$user_detail_response" != *"\"email\":\"runtime.user.prime@$tenant_slug.local\""* || "$user_detail_response" != *'"displayName":"Runtime User Prime"'* ]]; then
    echo "[test] runtime identity user detail did not reflect live update"
    exit 1
  fi

  teams_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/teams")"
  echo "[test] identity api teams => $teams_response"

  if [[ "$teams_response" != *'"name":"Core"'* ]]; then
    echo "[test] default runtime identity team was not provisioned"
    exit 1
  fi

  create_team_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"name":"Field Ops"}' \
    "$base_url/api/identity/tenants/$tenant_slug/teams")"
  echo "[test] identity api create team => $create_team_response"

  created_team_public_id="$(echo "$create_team_response" | sed -n 's/.*"publicId":"\([^"]*\)".*/\1/p')"
  if [[ -z "$created_team_public_id" ]]; then
    echo "[test] runtime identity team public id was not returned"
    exit 1
  fi

  updated_team_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"name":"Field Ops Prime"}' \
    "$base_url/api/identity/tenants/$tenant_slug/teams/$created_team_public_id")"
  echo "[test] identity api update team => $updated_team_response"

  if [[ "$updated_team_response" != *'"name":"Field Ops Prime"'* ]]; then
    echo "[test] runtime identity team update did not persist"
    exit 1
  fi

  members_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"userPublicId\":\"$created_user_public_id\"}" \
    "$base_url/api/identity/tenants/$tenant_slug/teams/$created_team_public_id/members")"
  echo "[test] identity api add member => $members_response"

  if [[ "$members_response" != *"\"userPublicId\":\"$created_user_public_id\""* ]]; then
    echo "[test] runtime identity team membership did not persist"
    exit 1
  fi

  removed_member_response="$(curl -fsS \
    -X DELETE \
    "$base_url/api/identity/tenants/$tenant_slug/teams/$created_team_public_id/members/$created_user_public_id")"
  echo "[test] identity api remove member => $removed_member_response"

  if [[ "$removed_member_response" != *"\"userPublicId\":\"$created_user_public_id\""* ]]; then
    echo "[test] runtime identity team membership removal did not persist"
    exit 1
  fi

  assign_role_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"roleCode":"viewer"}' \
    "$base_url/api/identity/tenants/$tenant_slug/users/$created_user_public_id/roles")"
  echo "[test] identity api assign role => $assign_role_response"

  if [[ "$assign_role_response" != *'"roleCode":"viewer"'* ]]; then
    echo "[test] runtime identity role assignment did not persist"
    exit 1
  fi

  user_roles_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/users/$created_user_public_id/roles")"
  echo "[test] identity api user roles => $user_roles_response"

  if [[ "$user_roles_response" != *'"roleCode":"viewer"'* ]]; then
    echo "[test] runtime identity user role read did not return viewer"
    exit 1
  fi

  revoke_role_response="$(curl -fsS \
    -X DELETE \
    "$base_url/api/identity/tenants/$tenant_slug/users/$created_user_public_id/roles/viewer")"
  echo "[test] identity api revoke role => $revoke_role_response"

  if [[ "$revoke_role_response" != *'"roleCode":"viewer"'* ]]; then
    echo "[test] runtime identity role revoke did not persist"
    exit 1
  fi

  user_roles_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/users/$created_user_public_id/roles")"
  echo "[test] identity api user roles after revoke => $user_roles_response"

  if [[ "$user_roles_response" == *'"roleCode":"viewer"'* ]]; then
    echo "[test] runtime identity user role read still returned revoked viewer"
    exit 1
  fi

  roles_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/roles")"
  echo "[test] identity api roles => $roles_response"

  if [[ "$roles_response" != *'"code":"owner"'* || "$roles_response" != *'"code":"viewer"'* ]]; then
    echo "[test] runtime identity tenant roles did not return default role catalog"
    exit 1
  fi

  snapshot_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/snapshot")"
  echo "[test] identity api snapshot => $snapshot_response"

  if [[ "$snapshot_response" != *'"companies":2'* || "$snapshot_response" != *'"users":2'* || "$snapshot_response" != *'"teams":2'* || "$snapshot_response" != *'"roles":5'* || "$snapshot_response" != *'"teamMemberships":1'* || "$snapshot_response" != *'"userRoles":1'* || "$snapshot_response" != *'"name":"Field Ops Prime"'* ]]; then
    echo "[test] runtime identity snapshot did not reflect live updates"
    exit 1
  fi

  create_invite_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"invite.flow@$tenant_slug.local\",\"displayName\":\"Invite Flow User\",\"roleCodes\":[\"viewer\"],\"teamPublicIds\":[\"$created_team_public_id\"],\"expiresInDays\":7}" \
    "$base_url/api/identity/tenants/$tenant_slug/invites")"
  echo "[test] identity api create invite => $create_invite_response"

  invite_token="$(echo "$create_invite_response" | sed -n 's/.*"inviteToken":"\([^"]*\)".*/\1/p')"
  if [[ -z "$invite_token" || "$create_invite_response" != *'"status":"pending"'* || "$create_invite_response" != *'"roleCodes":["viewer"]'* ]]; then
    echo "[test] runtime identity invite create did not persist"
    exit 1
  fi

  invites_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/invites")"
  echo "[test] identity api invites => $invites_response"

  if [[ "$invites_response" != *"invite.flow@$tenant_slug.local"* ]]; then
    echo "[test] runtime identity invite list did not return the created invite"
    exit 1
  fi

  accept_invite_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"displayName":"Invite Flow User","givenName":"Invite","familyName":"Flow","password":"PhaseTwo123"}' \
    "$base_url/api/identity/invites/$invite_token/accept")"
  echo "[test] identity api accept invite => $accept_invite_response"

  invited_user_public_id="$(echo "$accept_invite_response" | grep -o '"publicId":"[^"]*"' | tail -n 1 | cut -d'"' -f4)"
  if [[ -z "$invited_user_public_id" || "$accept_invite_response" != *'"inviteStatus":"accepted"'* || "$accept_invite_response" != *"invite.flow@$tenant_slug.local"* ]]; then
    echo "[test] runtime identity invite accept did not activate the invited user"
    exit 1
  fi

  mfa_enroll_response="$(curl -fsS \
    -X POST \
    "$base_url/api/identity/tenants/$tenant_slug/users/$invited_user_public_id/mfa/enroll")"
  echo "[test] identity api mfa enroll => $mfa_enroll_response"

  mfa_secret="$(echo "$mfa_enroll_response" | sed -n 's/.*"secret":"\([^"]*\)".*/\1/p')"
  if [[ -z "$mfa_secret" || "$mfa_enroll_response" != *'"enabled":false'* ]]; then
    echo "[test] runtime identity MFA enrollment did not return a secret"
    exit 1
  fi

  mfa_code="$(compute_totp "$mfa_secret")"
  if [[ -z "$mfa_code" ]]; then
    echo "[test] runtime identity MFA helper did not generate an OTP"
    exit 1
  fi

  mfa_verify_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"otpCode\":\"$mfa_code\"}" \
    "$base_url/api/identity/tenants/$tenant_slug/users/$invited_user_public_id/mfa/verify")"
  echo "[test] identity api mfa verify => $mfa_verify_response"

  if [[ "$mfa_verify_response" != *'"enabled":true'* ]]; then
    echo "[test] runtime identity MFA verification did not enable the factor"
    exit 1
  fi

  login_without_otp_response="$(curl -sS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"tenantSlug\":\"$tenant_slug\",\"email\":\"invite.flow@$tenant_slug.local\",\"password\":\"PhaseTwo123\",\"otpCode\":null}" \
    -w ' HTTP_STATUS:%{http_code}' \
    "$base_url/api/identity/sessions/login")"
  echo "[test] identity api login without otp => $login_without_otp_response"

  if [[ "$login_without_otp_response" != *'"code":"mfa_required"'* || "$login_without_otp_response" != *'HTTP_STATUS:401'* ]]; then
    echo "[test] runtime identity login should require MFA after enablement"
    exit 1
  fi

  login_with_otp_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"tenantSlug\":\"$tenant_slug\",\"email\":\"invite.flow@$tenant_slug.local\",\"password\":\"PhaseTwo123\",\"otpCode\":\"$mfa_code\"}" \
    "$base_url/api/identity/sessions/login")"
  echo "[test] identity api login with otp => $login_with_otp_response"

  session_token="$(echo "$login_with_otp_response" | sed -n 's/.*"sessionToken":"\([^"]*\)".*/\1/p')"
  refresh_token="$(echo "$login_with_otp_response" | sed -n 's/.*"refreshToken":"\([^"]*\)".*/\1/p')"
  if [[ -z "$session_token" || -z "$refresh_token" || "$login_with_otp_response" != *'"mfaEnabled":true'* || "$login_with_otp_response" != *'"roleCodes":["viewer"]'* ]]; then
    echo "[test] runtime identity login with MFA did not create a usable session"
    exit 1
  fi

  access_response="$(curl -fsS \
    -H "Authorization: Bearer $session_token" \
    "$base_url/api/identity/tenants/$tenant_slug/access")"
  echo "[test] identity api access resolve => $access_response"

  if [[ "$access_response" != *'"authorized":true'* || "$access_response" != *'"status":"active"'* || "$access_response" != *'"roleCodes":["viewer"]'* ]]; then
    echo "[test] runtime identity access resolution did not enforce the tenant scope"
    exit 1
  fi

  refresh_session_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"refreshToken\":\"$refresh_token\"}" \
    "$base_url/api/identity/sessions/refresh")"
  echo "[test] identity api refresh session => $refresh_session_response"

  if [[ "$refresh_session_response" != *'"refreshToken":"'* || "$refresh_session_response" == *"\"refreshToken\":\"$refresh_token\""* ]]; then
    echo "[test] runtime identity refresh did not rotate the refresh token"
    exit 1
  fi

  blocked_access_response="$(curl -fsS \
    -X PATCH \
    -H "Content-Type: application/json" \
    -d '{"status":"suspended"}' \
    "$base_url/api/identity/tenants/$tenant_slug/users/$invited_user_public_id/access")"
  echo "[test] identity api block user => $blocked_access_response"

  if [[ "$blocked_access_response" != *'"status":"suspended"'* ]]; then
    echo "[test] runtime identity access block did not persist"
    exit 1
  fi

  login_without_otp_response="$(curl -sS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"tenantSlug\":\"$tenant_slug\",\"email\":\"invite.flow@$tenant_slug.local\",\"password\":\"PhaseTwo123\",\"otpCode\":\"$(compute_totp "$mfa_secret")\"}" \
    -w ' HTTP_STATUS:%{http_code}' \
    "$base_url/api/identity/sessions/login")"
  echo "[test] identity api login after block => $login_without_otp_response"

  if [[ "$login_without_otp_response" != *'"code":"access_blocked"'* || "$login_without_otp_response" != *'HTTP_STATUS:403'* ]]; then
    echo "[test] runtime identity login should be blocked after suspension"
    exit 1
  fi

  access_response="$(curl -sS \
    -H "Authorization: Bearer $session_token" \
    -w ' HTTP_STATUS:%{http_code}' \
    "$base_url/api/identity/tenants/$tenant_slug/access")"
  echo "[test] identity api access after block => $access_response"

  if [[ "$access_response" != *'"code":"invalid_session"'* || "$access_response" != *'HTTP_STATUS:401'* ]]; then
    echo "[test] runtime identity access should reject revoked sessions after suspension"
    exit 1
  fi

  audit_response="$(curl -fsS "$base_url/api/identity/tenants/$tenant_slug/security/audit")"
  echo "[test] identity api security audit => $audit_response"

  if [[ "$audit_response" != *'"eventCode":"invite_created"'* || "$audit_response" != *'"eventCode":"invite_accepted"'* || "$audit_response" != *'"eventCode":"mfa_enabled"'* || "$audit_response" != *'"eventCode":"access_blocked"'* ]]; then
    echo "[test] runtime identity security audit did not capture the expected events"
    exit 1
  fi
}

run_edge_runtime_smoke() {
  local base_url="http://localhost:${EDGE_HTTP_PORT:-8080}"
  local details_response
  local ops_health_response
  local session_response
  local session_token
  local tenant_overview_unauthorized_response
  local tenant_overview_response
  local automation_overview_response
  local sales_overview_response
  local revenue_overview_response

  "${COMPOSE_CMD[@]}" up -d --build edge
  wait_for_http_ready "$base_url/health/ready"

  details_response="$(curl -fsS "$base_url/health/details")"
  ops_health_response="$(curl -fsS "$base_url/api/edge/ops/health")"
  session_response="$(curl -fsS \
    -X POST \
    -H "Content-Type: application/json" \
    -d "{\"tenantSlug\":\"bootstrap-ops\",\"email\":\"owner@bootstrap-ops.local\",\"password\":\"${IDENTITY_BOOTSTRAP_PASSWORD:-Change.Me123!}\",\"otpCode\":null}" \
    "http://localhost:${IDENTITY_HTTP_PORT:-8081}/api/identity/sessions/login")"
  session_token="$(echo "$session_response" | sed -n 's/.*"sessionToken":"\([^"]*\)".*/\1/p')"
  tenant_overview_unauthorized_response="$(curl -sS -w ' HTTP_STATUS:%{http_code}' "$base_url/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops")"
  tenant_overview_response="$(curl -fsS -H "Authorization: Bearer $session_token" "$base_url/api/edge/ops/tenant-overview?tenantSlug=bootstrap-ops")"
  automation_overview_response="$(curl -fsS -H "Authorization: Bearer $session_token" "$base_url/api/edge/ops/automation-overview?tenantSlug=bootstrap-ops")"
  sales_overview_response="$(curl -fsS -H "Authorization: Bearer $session_token" "$base_url/api/edge/ops/sales-overview?tenantSlug=bootstrap-ops")"
  revenue_overview_response="$(curl -fsS -H "Authorization: Bearer $session_token" "$base_url/api/edge/ops/revenue-overview?tenantSlug=bootstrap-ops")"

  echo "[test] edge health details => $details_response"
  echo "[test] edge ops health => $ops_health_response"
  echo "[test] edge bootstrap session => $session_response"
  echo "[test] edge tenant overview unauthorized => $tenant_overview_unauthorized_response"
  echo "[test] edge tenant overview => $tenant_overview_response"
  echo "[test] edge automation overview => $automation_overview_response"
  echo "[test] edge sales overview => $sales_overview_response"
  echo "[test] edge revenue overview => $revenue_overview_response"

  if [[ -z "$session_token" || "$details_response" != *'"service":"edge"'* || "$details_response" != *'"status":"ready"'* || "$details_response" != *'"name":"identity","status":"ready"'* || "$details_response" != *'"name":"analytics","status":"ready"'* || "$details_response" != *'"name":"webhook-hub","status":"ready"'* || "$details_response" != *'"name":"sales","status":"ready"'* || "$ops_health_response" != *'"service":"edge"'* || "$ops_health_response" != *'"status":"ready"'* || "$ops_health_response" != *'"total":7'* || "$ops_health_response" != *'"ready":7'* || "$ops_health_response" != *'"degraded":0'* || "$ops_health_response" != *'"name":"workflow-runtime"'* || "$ops_health_response" != *'"name":"analytics"'* || "$ops_health_response" != *'"name":"sales"'* || "$tenant_overview_unauthorized_response" != *'"code":"session_required"'* || "$tenant_overview_unauthorized_response" != *'HTTP_STATUS:401'* || "$tenant_overview_response" != *'"service":"edge"'* || "$tenant_overview_response" != *'"tenantSlug":"bootstrap-ops"'* || "$tenant_overview_response" != *'"automationBoard"'* || "$tenant_overview_response" != *'"workflowDefinitionKey":"lead-follow-up"'* || "$automation_overview_response" != *'"service":"edge"'* || "$automation_overview_response" != *'"tenantSlug":"bootstrap-ops"'* || "$automation_overview_response" != *'"status":"attention"'* || "$automation_overview_response" != *'"activeDefinitions":1'* || "$automation_overview_response" != *'"stableDefinitions":1'* || "$automation_overview_response" != *'"attentionDefinitions":1'* || "$automation_overview_response" != *'"criticalDefinitions":0'* || "$automation_overview_response" != *'"runningControlRuns":2'* || "$automation_overview_response" != *'"completedRuntimeExecutions":3'* || "$automation_overview_response" != *'"forwardedWebhookEvents":1'* || "$automation_overview_response" != *'"workflowDefinitionHealth"'* || "$sales_overview_response" != *'"service":"edge"'* || "$sales_overview_response" != *'"tenantSlug":"bootstrap-ops"'* || "$sales_overview_response" != *'"status":"stable"'* || "$sales_overview_response" != *'"leadsCaptured":2'* || "$sales_overview_response" != *'"opportunities":2'* || "$sales_overview_response" != *'"proposals":2'* || "$sales_overview_response" != *'"salesWon":2'* || "$sales_overview_response" != *'"bookedRevenueCents":224000'* || "$sales_overview_response" != *'"completedAutomations":1'* || "$sales_overview_response" != *'"salesJourney"'* || "$revenue_overview_response" != *'"service":"edge"'* || "$revenue_overview_response" != *'"tenantSlug":"bootstrap-ops"'* || "$revenue_overview_response" != *'"status":"stable"'* || "$revenue_overview_response" != *'"salesWon":2'* || "$revenue_overview_response" != *'"invoices":2'* || "$revenue_overview_response" != *'"paidInvoices":1'* || "$revenue_overview_response" != *'"openAmountCents":125000'* || "$revenue_overview_response" != *'"paidAmountCents":99000'* || "$revenue_overview_response" != *'"overdueInvoices":0'* || "$revenue_overview_response" != *'"collectionRateBps":4420'* || "$revenue_overview_response" != *'"revenueOperations"'* ]]; then
    echo "[test] edge automation cockpit did not aggregate the expected live payload"
    exit 1
  fi
}

run_smoke() {
  trap 'bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true' RETURN
  export POSTGRES_PORT="${POSTGRES_PORT_SMOKE:-16432}"
  export REDIS_PORT="${REDIS_PORT_SMOKE:-16379}"
  export KAFKA_PORT="${KAFKA_PORT_SMOKE:-19092}"
  export PROMETHEUS_PORT="${PROMETHEUS_PORT_SMOKE:-19090}"
  export GRAFANA_PORT="${GRAFANA_PORT_SMOKE:-13000}"
  export KEYCLOAK_PORT="${KEYCLOAK_PORT_SMOKE:-18089}"
  export OPENFGA_HTTP_PORT="${OPENFGA_HTTP_PORT_SMOKE:-18090}"
  export OPENFGA_GRPC_PORT="${OPENFGA_GRPC_PORT_SMOKE:-18091}"
  export OPENFGA_PLAYGROUND_PORT="${OPENFGA_PLAYGROUND_PORT_SMOKE:-13010}"
  export EDGE_HTTP_PORT="${EDGE_HTTP_PORT_SMOKE:-18080}"
  export CRM_HTTP_PORT="${CRM_HTTP_PORT_SMOKE:-18083}"
  export SALES_HTTP_PORT="${SALES_HTTP_PORT_SMOKE:-18087}"
  export IDENTITY_HTTP_PORT="${IDENTITY_HTTP_PORT_SMOKE:-18081}"
  export WEBHOOK_HUB_HTTP_PORT="${WEBHOOK_HUB_HTTP_PORT_SMOKE:-18082}"
  export WORKFLOW_CONTROL_HTTP_PORT="${WORKFLOW_CONTROL_HTTP_PORT_SMOKE:-18084}"
  export WORKFLOW_RUNTIME_HTTP_PORT="${WORKFLOW_RUNTIME_HTTP_PORT_SMOKE:-18085}"
  export ANALYTICS_HTTP_PORT="${ANALYTICS_HTTP_PORT_SMOKE:-18086}"
  export ENGAGEMENT_HTTP_PORT="${ENGAGEMENT_HTTP_PORT_SMOKE:-18088}"
  bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true
  "${COMPOSE_CMD[@]}" ps
  run_platform_runtime_smoke
  run_identity_database_smoke
  run_webhook_hub_runtime_smoke
  run_workflow_control_runtime_smoke
  run_workflow_runtime_smoke
  run_crm_runtime_smoke
  run_sales_runtime_smoke
  run_engagement_runtime_smoke
  run_analytics_runtime_smoke
  run_identity_runtime_smoke
  run_edge_runtime_smoke
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/test.sh unit
  ./scripts/test.sh integration
  ./scripts/test.sh contract
  ./scripts/test.sh platform
  ./scripts/test.sh smoke
  ./scripts/test.sh all
EOF
}

main() {
  local command="${1:-}"

  case "$command" in
    unit)
      run_go_unit
      run_typescript_unit
      run_elixir_unit
      run_python_unit
      run_dotnet_build
      run_rust_unit
      ;;
    integration)
      run_dotnet_integration
      ;;
    contract)
      run_typescript_contract
      run_go_contract
      run_dotnet_contract
      ;;
    platform)
      trap 'bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true' RETURN
      export POSTGRES_PORT="${POSTGRES_PORT_PLATFORM:-16432}"
      export REDIS_PORT="${REDIS_PORT_PLATFORM:-16379}"
      export KAFKA_PORT="${KAFKA_PORT_PLATFORM:-19092}"
      export PROMETHEUS_PORT="${PROMETHEUS_PORT_PLATFORM:-19090}"
      export GRAFANA_PORT="${GRAFANA_PORT_PLATFORM:-13000}"
      export KEYCLOAK_PORT="${KEYCLOAK_PORT_PLATFORM:-18089}"
      export OPENFGA_HTTP_PORT="${OPENFGA_HTTP_PORT_PLATFORM:-18090}"
      export OPENFGA_GRPC_PORT="${OPENFGA_GRPC_PORT_PLATFORM:-18091}"
      export OPENFGA_PLAYGROUND_PORT="${OPENFGA_PLAYGROUND_PORT_PLATFORM:-13010}"
      bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true
      run_platform_runtime_smoke
      ;;
    smoke)
      run_smoke
      ;;
    all)
      run_go_unit
      run_typescript_unit
      run_elixir_unit
      run_python_unit
      run_dotnet_build
      run_dotnet_integration
      run_typescript_contract
      run_go_contract
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
