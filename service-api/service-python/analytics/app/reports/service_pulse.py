"""Pulso operacional transversal entre os principais serviços do ERP."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import append_filter, tenant_filter


def build_service_pulse(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_service_pulse(tenant_slug)

    return build_static_service_pulse(tenant_slug)


def build_static_service_pulse(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "services": {
            "crm": {"totalLeads": 128, "captured": 18, "contacted": 42, "qualified": 37, "disqualified": 31},
            "workflowControl": {"activeDefinitions": 6, "runsRunning": 7, "runsCompleted": 31, "runsFailed": 2, "runsCancelled": 1},
            "workflowRuntime": {"totalExecutions": 44, "running": 4, "completed": 28, "failed": 8, "cancelled": 4},
            "webhookHub": {"totalEvents": 93, "forwarded": 87, "queued": 2, "processing": 1, "failed": 3},
        },
    }


def build_postgres_service_pulse(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        crm_metrics = fetch_crm_metrics(connection, tenant_slug)
        workflow_control_metrics = fetch_workflow_control_metrics(connection, tenant_slug)
        workflow_runtime_metrics = fetch_workflow_runtime_metrics(connection, tenant_slug)
        webhook_hub_metrics = fetch_webhook_hub_metrics(connection)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "services": {
            "crm": crm_metrics,
            "workflowControl": workflow_control_metrics,
            "workflowRuntime": workflow_runtime_metrics,
            "webhookHub": webhook_hub_metrics,
        },
    }


def fetch_crm_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total_leads,
            count(*) FILTER (WHERE lead.status = 'captured') AS captured,
            count(*) FILTER (WHERE lead.status = 'contacted') AS contacted,
            count(*) FILTER (WHERE lead.status = 'qualified') AS qualified,
            count(*) FILTER (WHERE lead.status = 'disqualified') AS disqualified
        FROM crm.leads AS lead
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "totalLeads": int(row.get("total_leads", 0) or 0),
        "captured": int(row.get("captured", 0) or 0),
        "contacted": int(row.get("contacted", 0) or 0),
        "qualified": int(row.get("qualified", 0) or 0),
        "disqualified": int(row.get("disqualified", 0) or 0),
    }


def fetch_workflow_control_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    definitions_query = f"""
        SELECT count(*) AS total
        FROM workflow_control.workflow_definitions AS definition
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        WHERE definition.status = 'active'
        {append_filter(filter_sql)}
    """

    runs_query = f"""
        SELECT
            count(*) FILTER (WHERE run.status = 'running') AS running,
            count(*) FILTER (WHERE run.status = 'completed') AS completed,
            count(*) FILTER (WHERE run.status = 'failed') AS failed,
            count(*) FILTER (WHERE run.status = 'cancelled') AS cancelled
        FROM workflow_control.workflow_runs AS run
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(definitions_query, params)
        definitions_row = cursor.fetchone() or {"total": 0}
        cursor.execute(runs_query, params)
        runs_row = cursor.fetchone() or {}

    return {
        "activeDefinitions": int(definitions_row.get("total", 0) or 0),
        "runsRunning": int(runs_row.get("running", 0) or 0),
        "runsCompleted": int(runs_row.get("completed", 0) or 0),
        "runsFailed": int(runs_row.get("failed", 0) or 0),
        "runsCancelled": int(runs_row.get("cancelled", 0) or 0),
    }


def fetch_workflow_runtime_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE execution.status = 'running') AS running,
            count(*) FILTER (WHERE execution.status = 'completed') AS completed,
            count(*) FILTER (WHERE execution.status = 'failed') AS failed,
            count(*) FILTER (WHERE execution.status = 'cancelled') AS cancelled
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "totalExecutions": int(row.get("total", 0) or 0),
        "running": int(row.get("running", 0) or 0),
        "completed": int(row.get("completed", 0) or 0),
        "failed": int(row.get("failed", 0) or 0),
        "cancelled": int(row.get("cancelled", 0) or 0),
    }


def fetch_webhook_hub_metrics(connection) -> dict:
    with connection.cursor() as cursor:
        cursor.execute(
            """
                SELECT
                    count(*) AS total,
                    count(*) FILTER (WHERE status = 'forwarded') AS forwarded,
                    count(*) FILTER (WHERE status = 'queued') AS queued,
                    count(*) FILTER (WHERE status = 'processing') AS processing,
                    count(*) FILTER (WHERE status = 'failed') AS failed
                FROM webhook_hub.webhook_events
            """
        )
        row = cursor.fetchone() or {}

    return {
        "totalEvents": int(row.get("total", 0) or 0),
        "forwarded": int(row.get("forwarded", 0) or 0),
        "queued": int(row.get("queued", 0) or 0),
        "processing": int(row.get("processing", 0) or 0),
        "failed": int(row.get("failed", 0) or 0),
    }
