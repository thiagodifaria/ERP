"""Relatorio consolidado do tenant cruzando identidade, comercial e automacao."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import append_filter, tenant_filter


def build_tenant_360(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_tenant_360(tenant_slug)

    return build_static_tenant_360(tenant_slug)


def build_static_tenant_360(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "identity": {
            "companies": 3,
            "users": 18,
            "teams": 6,
            "roles": 5,
            "teamMemberships": 14,
        },
        "commercial": {
            "leads": 128,
            "qualifiedLeads": 37,
            "assignedLeads": 96,
            "leadNotes": 52,
        },
        "automation": {
            "activeDefinitions": 6,
            "workflowRuns": 41,
            "workflowRunEvents": 84,
            "runtimeExecutions": 44,
            "runtimeCompleted": 28,
            "runtimeFailed": 8,
            "runtimeCancelled": 4,
        },
    }


def build_postgres_tenant_360(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        identity_metrics = fetch_identity_metrics(connection, tenant_slug)
        commercial_metrics = fetch_commercial_metrics(connection, tenant_slug)
        automation_metrics = fetch_automation_metrics(connection, tenant_slug)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "identity": identity_metrics,
        "commercial": commercial_metrics,
        "automation": automation_metrics,
    }


def fetch_identity_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM identity.companies
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS companies,
            (
                SELECT count(*)
                FROM identity.users
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS users,
            (
                SELECT count(*)
                FROM identity.teams
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS teams,
            (
                SELECT count(*)
                FROM identity.roles
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS roles,
            (
                SELECT count(*)
                FROM identity.team_memberships
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS team_memberships
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 5)
        row = cursor.fetchone() or {}

    return {
        "companies": int(row.get("companies", 0) or 0),
        "users": int(row.get("users", 0) or 0),
        "teams": int(row.get("teams", 0) or 0),
        "roles": int(row.get("roles", 0) or 0),
        "teamMemberships": int(row.get("team_memberships", 0) or 0),
    }


def fetch_commercial_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    leads_query = f"""
        SELECT
            count(*) AS leads,
            count(*) FILTER (WHERE lead.status = 'qualified') AS qualified_leads,
            count(*) FILTER (WHERE lead.owner_user_public_id IS NOT NULL) AS assigned_leads
        FROM crm.leads AS lead
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    notes_query = f"""
        SELECT count(*) AS lead_notes
        FROM crm.lead_notes AS note
        JOIN crm.leads AS lead ON lead.id = note.lead_id
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(leads_query, params)
        leads_row = cursor.fetchone() or {}
        cursor.execute(notes_query, params)
        notes_row = cursor.fetchone() or {}

    return {
        "leads": int(leads_row.get("leads", 0) or 0),
        "qualifiedLeads": int(leads_row.get("qualified_leads", 0) or 0),
        "assignedLeads": int(leads_row.get("assigned_leads", 0) or 0),
        "leadNotes": int(notes_row.get("lead_notes", 0) or 0),
    }


def fetch_automation_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    definitions_query = f"""
        SELECT count(*) AS active_definitions
        FROM workflow_control.workflow_definitions AS definition
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        WHERE definition.status = 'active'
        {append_filter(filter_sql)}
    """

    runs_query = f"""
        SELECT
            count(*) AS workflow_runs,
            count(*) FILTER (WHERE run.status = 'completed') AS runs_completed,
            count(*) FILTER (WHERE run.status = 'failed') AS runs_failed,
            count(*) FILTER (WHERE run.status = 'cancelled') AS runs_cancelled
        FROM workflow_control.workflow_runs AS run
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    events_query = f"""
        SELECT count(*) AS workflow_run_events
        FROM workflow_control.workflow_run_events AS event
        JOIN workflow_control.workflow_runs AS run ON run.id = event.workflow_run_id
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    runtime_query = f"""
        SELECT
            count(*) AS runtime_executions,
            count(*) FILTER (WHERE execution.status = 'completed') AS runtime_completed,
            count(*) FILTER (WHERE execution.status = 'failed') AS runtime_failed,
            count(*) FILTER (WHERE execution.status = 'cancelled') AS runtime_cancelled
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(definitions_query, params)
        definitions_row = cursor.fetchone() or {}
        cursor.execute(runs_query, params)
        runs_row = cursor.fetchone() or {}
        cursor.execute(events_query, params)
        events_row = cursor.fetchone() or {}
        cursor.execute(runtime_query, params)
        runtime_row = cursor.fetchone() or {}

    return {
        "activeDefinitions": int(definitions_row.get("active_definitions", 0) or 0),
        "workflowRuns": int(runs_row.get("workflow_runs", 0) or 0),
        "workflowRunEvents": int(events_row.get("workflow_run_events", 0) or 0),
        "runtimeExecutions": int(runtime_row.get("runtime_executions", 0) or 0),
        "runtimeCompleted": int(runtime_row.get("runtime_completed", 0) or 0),
        "runtimeFailed": int(runtime_row.get("runtime_failed", 0) or 0),
        "runtimeCancelled": int(runtime_row.get("runtime_cancelled", 0) or 0),
    }
