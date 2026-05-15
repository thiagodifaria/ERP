#!/usr/bin/env bash
set -euo pipefail

target_file="${1:-.env}"

if [[ -e "$target_file" ]]; then
  echo "Refusing to overwrite existing $target_file" >&2
  exit 1
fi

if [[ ! -f .env.example ]]; then
  echo "Run this script from the repository root." >&2
  exit 1
fi

generate_secret() {
  openssl rand -base64 48 | tr -d '\n'
}

cp .env.example "$target_file"

replace_value() {
  local key="$1"
  local value="$2"

  sed -i "s|^${key}=.*|${key}=${value}|" "$target_file"
}

replace_value "ERP_JWT_HS256_SECRET" "$(generate_secret)"
replace_value "ERP_INTERNAL_SERVICE_TOKEN" "$(generate_secret)"
replace_value "ERP_POSTGRES_PASSWORD" "$(generate_secret)"
replace_value "KEYCLOAK_ADMIN_PASSWORD" "$(generate_secret)"
replace_value "GRAFANA_ADMIN_PASSWORD" "$(generate_secret)"
replace_value "IDENTITY_BOOTSTRAP_PASSWORD" "$(generate_secret)"
replace_value "DOCUMENTS_ACCESS_TOKEN_SECRET" "$(generate_secret)"

echo "Generated $target_file with local high-entropy secrets."
