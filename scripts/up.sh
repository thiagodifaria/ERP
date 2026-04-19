#!/usr/bin/env bash
set -euo pipefail

# Este script sobe o ecossistema local minimo do ERP.

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
  echo "[up] remapped $label host port from $requested_port to $fallback_port because it is already in use"
}

remap_host_port_if_needed "POSTGRES_PORT" "postgresql"
remap_host_port_if_needed "REDIS_PORT" "redis"
remap_host_port_if_needed "EDGE_HTTP_PORT" "edge"
remap_host_port_if_needed "IDENTITY_HTTP_PORT" "identity"
remap_host_port_if_needed "WEBHOOK_HUB_HTTP_PORT" "webhook-hub"
remap_host_port_if_needed "CRM_HTTP_PORT" "crm"
remap_host_port_if_needed "WORKFLOW_CONTROL_HTTP_PORT" "workflow-control"
remap_host_port_if_needed "WORKFLOW_RUNTIME_HTTP_PORT" "workflow-runtime"
remap_host_port_if_needed "ANALYTICS_HTTP_PORT" "analytics"
remap_host_port_if_needed "SALES_HTTP_PORT" "sales"
export ERP_HOST_PORTS_LOCKED=1

docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build "$@"
