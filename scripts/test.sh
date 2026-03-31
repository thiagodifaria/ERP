#!/usr/bin/env bash
set -euo pipefail

# Este script centraliza validacoes tecnicas em modo container-first.

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker-compose.yml"
ENV_FILE="$ROOT_DIR/.env"

if [[ ! -f "$ENV_FILE" ]]; then
  ENV_FILE="$ROOT_DIR/.env.example"
fi

COMPOSE_CMD=(docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE")

run_go_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-golang/edge:/workspace" \
    -w /workspace \
    golang:1.24-alpine \
    go test ./...
}

run_dotnet_build() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-csharp/identity:/workspace" \
    -w /workspace \
    mcr.microsoft.com/dotnet/sdk:8.0 \
    dotnet build src/Identity.Api/Identity.Api.csproj -c Release
}

run_rust_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-rust/webhook-hub:/workspace" \
    -w /workspace \
    rust:1 \
    cargo test
}

run_smoke() {
  "${COMPOSE_CMD[@]}" ps
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
      echo "[test] integration suites will be wired as persistence and messaging adapters land"
      ;;
    contract)
      echo "[test] contract suites will be wired as public endpoints and events stabilize"
      ;;
    smoke)
      run_smoke
      ;;
    all)
      run_go_unit
      run_dotnet_build
      run_rust_unit
      echo "[test] integration suites will be wired as persistence and messaging adapters land"
      echo "[test] contract suites will be wired as public endpoints and events stabilize"
      ;;
    *)
      usage
      exit 1
      ;;
  esac
}

main "$@"
