#!/usr/bin/env bash
set -euo pipefail

# Este script centraliza a aplicacao das migrations do PostgreSQL em ambiente local.

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

prepare_database_port() {
  if [[ -n "${ERP_HOST_PORTS_LOCKED:-}" ]]; then
    return
  fi

  local requested_port="${POSTGRES_PORT:-}"

  if [[ -z "$requested_port" ]]; then
    return
  fi

  if ! is_tcp_port_in_use "$requested_port"; then
    return
  fi

  local fallback_start=$((requested_port + 1000))
  local fallback_port
  fallback_port="$(find_available_port "$fallback_start")"
  export POSTGRES_PORT="$fallback_port"
  echo "[db] remapped postgresql host port from $requested_port to $fallback_port because it is already in use"
}

prepare_database_port

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

backup_database() {
  local output_path="$1"
  local output_directory

  output_directory="$(dirname "$output_path")"
  mkdir -p "$output_directory"
  echo "[db] backing up database to $output_path"
  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    pg_dump --clean --if-exists --no-owner --no-privileges -U "$DB_USER" -d "$DB_NAME" > "$output_path"
}

restore_database() {
  local input_path="$1"

  if [[ ! -f "$input_path" ]]; then
    echo "[db] backup file not found: $input_path"
    exit 1
  fi

  echo "[db] restoring database from $input_path"
  "${COMPOSE_CMD[@]}" exec -T service-postgresql \
    psql -v ON_ERROR_STOP=1 -U "$DB_USER" -d "$DB_NAME" < "$input_path"
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
  ./scripts/db.sh migrate sales
  ./scripts/db.sh migrate finance
  ./scripts/db.sh migrate documents
  ./scripts/db.sh migrate analytics
  ./scripts/db.sh migrate simulation
  ./scripts/db.sh migrate engagement
  ./scripts/db.sh migrate webhook-hub
  ./scripts/db.sh migrate workflow-control
  ./scripts/db.sh migrate workflow-runtime
  ./scripts/db.sh migrate all
  ./scripts/db.sh seed identity
  ./scripts/db.sh seed crm
  ./scripts/db.sh seed sales
  ./scripts/db.sh seed engagement
  ./scripts/db.sh seed workflow-control
  ./scripts/db.sh seed all
  ./scripts/db.sh backup /tmp/erp-backup.sql
  ./scripts/db.sh restore /tmp/erp-backup.sql
  ./scripts/db.sh summary identity [tenant-slug]
  ./scripts/db.sh summary crm [tenant-slug]
  ./scripts/db.sh summary sales [tenant-slug]
  ./scripts/db.sh summary finance [tenant-slug]
  ./scripts/db.sh summary documents [tenant-slug]
  ./scripts/db.sh summary analytics [tenant-slug]
  ./scripts/db.sh summary simulation [tenant-slug]
  ./scripts/db.sh summary engagement [tenant-slug]
  ./scripts/db.sh summary webhook-hub
  ./scripts/db.sh summary workflow-control [tenant-slug]
  ./scripts/db.sh summary workflow-runtime [tenant-slug]
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
        sales)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/sales/migrations"
          ;;
        finance)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/finance/migrations"
          ;;
        documents)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/documents/migrations"
          ;;
        analytics)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/analytics/migrations"
          ;;
        simulation)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/simulation/migrations"
          ;;
        engagement)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/engagement/migrations"
          ;;
        webhook-hub)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/webhook-hub/migrations"
          ;;
        workflow-control)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/migrations"
          ;;
        workflow-runtime)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-runtime/migrations"
          ;;
        all)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/common/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/crm/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/sales/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/finance/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/documents/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/analytics/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/simulation/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/engagement/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/webhook-hub/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/migrations"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-runtime/migrations"
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
        sales)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/sales/seeds"
          ;;
        engagement)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/engagement/seeds"
          ;;
        workflow-control)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/seeds"
          ;;
        all)
          apply_directory "$ROOT_DIR/service-api/service-postgresql/identity/seeds"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/crm/seeds"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/sales/seeds"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/engagement/seeds"
          apply_directory "$ROOT_DIR/service-api/service-postgresql/workflow-control/seeds"
          ;;
        *)
          usage
          exit 1
          ;;
      esac
      ;;
    backup)
      ensure_postgres
      local output_path="${2:-$ROOT_DIR/.cache/backups/erp-local-backup.sql}"
      backup_database "$output_path"
      ;;
    restore)
      ensure_postgres
      local input_path="${2:-}"

      if [[ -z "$input_path" ]]; then
        usage
        exit 1
      fi

      restore_database "$input_path"
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
              (SELECT count(*) FROM crm.customers AS customer WHERE customer.tenant_id = tenant.id) AS customers,
              (SELECT count(*) FROM crm.lead_notes AS note WHERE note.tenant_id = tenant.id) AS notes,
              (SELECT count(*) FROM crm.relationship_events AS event WHERE event.tenant_id = tenant.id) AS history_events,
              (SELECT count(*) FROM crm.outbox_events AS event WHERE event.tenant_id = tenant.id AND event.status = 'pending') AS pending_outbox,
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
        sales)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM sales.opportunities AS opportunity WHERE opportunity.tenant_id = tenant.id) AS opportunities,
              (SELECT count(*) FROM sales.proposals AS proposal WHERE proposal.tenant_id = tenant.id) AS proposals,
              (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id) AS sales,
              (SELECT count(*) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id) AS invoices,
              (SELECT count(*) FROM sales.opportunities AS opportunity WHERE opportunity.tenant_id = tenant.id AND opportunity.stage = 'won') AS won,
              (SELECT count(*) FROM sales.opportunities AS opportunity WHERE opportunity.tenant_id = tenant.id AND opportunity.stage = 'lost') AS lost,
              (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id AND sale.status = 'active') AS active_sales,
              (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id AND sale.status = 'invoiced') AS invoiced_sales,
              (SELECT count(*) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id AND invoice.status = 'paid') AS paid_invoices,
              (SELECT count(*) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id AND invoice.status NOT IN ('paid', 'cancelled')) AS open_invoices,
              (SELECT COALESCE(sum(sale.amount_cents), 0) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id AND sale.status <> 'cancelled') AS booked_revenue_cents,
              (SELECT COALESCE(sum(invoice.amount_cents), 0) FROM sales.invoices AS invoice WHERE invoice.tenant_id = tenant.id AND invoice.status = 'paid') AS collected_revenue_cents
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        engagement)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM engagement.campaigns AS campaign WHERE campaign.tenant_id = tenant.id) AS campaigns,
              (SELECT count(*) FROM engagement.touchpoints AS touchpoint WHERE touchpoint.tenant_id = tenant.id) AS touchpoints,
              (SELECT count(*) FROM engagement.campaigns AS campaign WHERE campaign.tenant_id = tenant.id AND campaign.status = 'active') AS active_campaigns,
              (SELECT count(*) FROM engagement.campaigns AS campaign WHERE campaign.tenant_id = tenant.id AND campaign.status = 'paused') AS paused_campaigns,
              (SELECT count(*) FROM engagement.touchpoints AS touchpoint WHERE touchpoint.tenant_id = tenant.id AND touchpoint.status = 'responded') AS responded_touchpoints,
              (SELECT count(*) FROM engagement.touchpoints AS touchpoint WHERE touchpoint.tenant_id = tenant.id AND touchpoint.status = 'converted') AS converted_touchpoints,
              (SELECT count(*) FROM engagement.touchpoints AS touchpoint WHERE touchpoint.tenant_id = tenant.id AND touchpoint.last_workflow_run_public_id IS NOT NULL) AS workflow_dispatched
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        finance)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id) AS projections,
              (SELECT count(*) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id AND projection.status = 'forecast') AS forecast,
              (SELECT count(*) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id AND projection.status = 'open') AS projection_open,
              (SELECT count(*) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id AND projection.status = 'paid') AS projection_paid,
              (SELECT count(*) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id AND projection.status = 'cancelled') AS projection_cancelled,
              (SELECT COALESCE(sum(projection.amount_cents), 0) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id AND projection.status IN ('forecast', 'open')) AS pipeline_amount_cents,
              (SELECT COALESCE(sum(projection.amount_cents), 0) FROM finance.receivable_projections AS projection WHERE projection.tenant_id = tenant.id AND projection.status = 'paid') AS projected_paid_amount_cents,
              (SELECT count(*) FROM finance.receivable_entries AS receivable WHERE receivable.tenant_id = tenant.id) AS receivables,
              (SELECT count(*) FROM finance.receivable_entries AS receivable WHERE receivable.tenant_id = tenant.id AND receivable.status = 'open') AS receivable_open,
              (SELECT count(*) FROM finance.receivable_entries AS receivable WHERE receivable.tenant_id = tenant.id AND receivable.status = 'paid') AS receivable_paid,
              (SELECT count(*) FROM finance.receivable_entries AS receivable WHERE receivable.tenant_id = tenant.id AND receivable.status = 'cancelled') AS receivable_cancelled,
              (SELECT COALESCE(sum(receivable.amount_cents), 0) FROM finance.receivable_entries AS receivable WHERE receivable.tenant_id = tenant.id AND receivable.status = 'open') AS receivable_open_amount_cents,
              (SELECT COALESCE(sum(receivable.amount_cents), 0) FROM finance.receivable_entries AS receivable WHERE receivable.tenant_id = tenant.id AND receivable.status = 'paid') AS receivable_paid_amount_cents,
              (SELECT count(*) FROM finance.commission_entries AS commission WHERE commission.tenant_id = tenant.id) AS commissions,
              (SELECT count(*) FROM finance.commission_entries AS commission WHERE commission.tenant_id = tenant.id AND commission.status = 'pending') AS commission_pending,
              (SELECT count(*) FROM finance.commission_entries AS commission WHERE commission.tenant_id = tenant.id AND commission.status = 'blocked') AS commission_blocked,
              (SELECT count(*) FROM finance.commission_entries AS commission WHERE commission.tenant_id = tenant.id AND commission.status = 'released') AS commission_released,
              (SELECT COALESCE(sum(commission.amount_cents), 0) FROM finance.commission_entries AS commission WHERE commission.tenant_id = tenant.id AND commission.status = 'released') AS commission_released_amount_cents,
              (SELECT count(*) FROM finance.payables AS payable WHERE payable.tenant_id = tenant.id) AS payables,
              (SELECT count(*) FROM finance.payables AS payable WHERE payable.tenant_id = tenant.id AND payable.status = 'open') AS payable_open,
              (SELECT count(*) FROM finance.payables AS payable WHERE payable.tenant_id = tenant.id AND payable.status = 'paid') AS payable_paid,
              (SELECT count(*) FROM finance.payables AS payable WHERE payable.tenant_id = tenant.id AND payable.status = 'cancelled') AS payable_cancelled,
              (SELECT COALESCE(sum(payable.amount_cents), 0) FROM finance.payables AS payable WHERE payable.tenant_id = tenant.id AND payable.status = 'paid') AS payable_paid_amount_cents,
              (SELECT count(*) FROM finance.cost_entries AS cost WHERE cost.tenant_id = tenant.id) AS costs,
              (SELECT COALESCE(sum(cost.amount_cents), 0) FROM finance.cost_entries AS cost WHERE cost.tenant_id = tenant.id) AS cost_amount_cents,
              (SELECT count(*) FROM finance.period_closures AS closure WHERE closure.tenant_id = tenant.id) AS closures
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        documents)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM documents.attachments AS attachment WHERE attachment.tenant_id = tenant.id) AS attachments,
              (SELECT count(*) FROM documents.attachments AS attachment WHERE attachment.tenant_id = tenant.id AND attachment.owner_type = 'crm.lead') AS lead_attachments,
              (SELECT count(*) FROM documents.attachments AS attachment WHERE attachment.tenant_id = tenant.id AND attachment.owner_type = 'crm.customer') AS customer_attachments,
              (SELECT count(*) FROM documents.attachments AS attachment WHERE attachment.tenant_id = tenant.id AND attachment.storage_driver = 'manual') AS manual_storage,
              (SELECT count(*) FROM documents.attachments AS attachment WHERE attachment.tenant_id = tenant.id AND attachment.storage_driver <> 'manual') AS external_storage
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        analytics)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM simulation.load_benchmark_runs AS benchmark WHERE benchmark.tenant_id = tenant.id) AS benchmarks,
              (SELECT count(*) FROM simulation.scenario_runs AS scenario WHERE scenario.tenant_id = tenant.id) AS scenarios
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        simulation)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM simulation.scenario_runs AS scenario WHERE scenario.tenant_id = tenant.id) AS scenarios,
              (SELECT count(*) FROM simulation.load_benchmark_runs AS benchmark WHERE benchmark.tenant_id = tenant.id) AS load_benchmarks
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        webhook-hub)
          run_psql_query "
            SELECT
              count(*) AS total,
              count(*) FILTER (WHERE status = 'received') AS received,
              count(*) FILTER (WHERE status = 'validated') AS validated,
              count(*) FILTER (WHERE status = 'queued') AS queued,
              count(*) FILTER (WHERE status = 'processing') AS processing,
              count(*) FILTER (WHERE status = 'forwarded') AS forwarded,
              count(*) FILTER (WHERE status = 'failed') AS failed,
              count(*) FILTER (WHERE status = 'rejected') AS rejected,
              (SELECT count(*) FROM webhook_hub.webhook_event_transitions) AS transitions
            FROM webhook_hub.webhook_events;
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
              (SELECT count(*) FROM workflow_control.workflow_definition_versions AS version WHERE version.tenant_id = tenant.id) AS versions,
              (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id) AS runs,
              (SELECT count(*) FROM workflow_control.workflow_run_events AS workflow_run_event WHERE workflow_run_event.tenant_id = tenant.id) AS events,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'draft') AS draft,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'active') AS active,
              (SELECT count(*) FROM workflow_control.workflow_definitions AS definition WHERE definition.tenant_id = tenant.id AND definition.status = 'archived') AS archived,
              (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'pending') AS pending_runs,
              (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'running') AS running_runs,
              (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'completed') AS completed_runs,
              (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'failed') AS failed_runs,
              (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id AND workflow_run.status = 'cancelled') AS cancelled_runs,
              (SELECT COALESCE(MAX(version.version_number), 0) FROM workflow_control.workflow_definition_versions AS version WHERE version.tenant_id = tenant.id) AS latest_version
            FROM identity.tenants AS tenant
            $where_clause
            ORDER BY tenant.slug;
          "
          ;;
        workflow-runtime)
          local tenant_slug="${3:-}"
          local where_clause=""

          if [[ -n "$tenant_slug" ]]; then
            where_clause="WHERE tenant.slug = '$tenant_slug'"
          fi

          run_psql_query "
            SELECT
              tenant.slug,
              (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id) AS executions,
              (SELECT count(*) FROM workflow_runtime.execution_transitions AS transition
                JOIN workflow_runtime.executions AS execution ON execution.id = transition.execution_id
               WHERE execution.tenant_id = tenant.id) AS transitions,
              (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id AND execution.status = 'pending') AS pending,
              (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id AND execution.status = 'running') AS running,
              (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id AND execution.status = 'completed') AS completed,
              (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id AND execution.status = 'failed') AS failed,
              (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id AND execution.status = 'cancelled') AS cancelled
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
