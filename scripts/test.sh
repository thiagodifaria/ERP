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

run_typescript_unit() {
  docker run --rm \
    -v "$ROOT_DIR/service-api/service-typescript/workflow-control:/workspace" \
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

  workflow_control_summary="$("${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -U "$DB_USER" -d "$DB_NAME" -At -c "
      SELECT
        tenant.slug || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_definition_versions AS version WHERE version.tenant_id = tenant.id) || '|' ||
        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id) || '|' ||
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

  if [[ "$workflow_control_summary" != "$smoke_slug|1|1|1|0|1|0|0|1|0|0|0|1" ]]; then
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
}

run_workflow_control_runtime_smoke() {
  local base_url="http://localhost:${WORKFLOW_CONTROL_HTTP_PORT:-8084}"
  local list_response
  local create_response
  local created_key="runtime-flow"
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
  echo "[test] workflow-control list => $list_response"

  if [[ "$list_response" != *'"key":"lead-follow-up"'* || "$list_response" != *'"status":"active"'* ]]; then
    echo "[test] workflow-control bootstrap catalog was not returned by the live API"
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

  if [[ "$details_response" != *'"name":"postgresql","status":"ready"'* ]]; then
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
}

run_smoke() {
  trap 'bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true' RETURN
  export CRM_HTTP_PORT="${CRM_HTTP_PORT_SMOKE:-18083}"
  export IDENTITY_HTTP_PORT="${IDENTITY_HTTP_PORT_SMOKE:-18081}"
  export WORKFLOW_CONTROL_HTTP_PORT="${WORKFLOW_CONTROL_HTTP_PORT_SMOKE:-18084}"
  bash "$ROOT_DIR/scripts/down.sh" -v >/dev/null 2>&1 || true
  "${COMPOSE_CMD[@]}" ps
  run_identity_database_smoke
  run_workflow_control_runtime_smoke
  run_crm_runtime_smoke
  run_identity_runtime_smoke
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
      run_typescript_unit
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
    smoke)
      run_smoke
      ;;
    all)
      run_go_unit
      run_typescript_unit
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
