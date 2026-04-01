"""Relatorio inicial de pipeline comercial para o plano analitico."""

from datetime import datetime, timezone

import psycopg

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_pipeline_summary(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_pipeline_summary(tenant_slug)

    return build_static_pipeline_summary(tenant_slug)


def build_static_pipeline_summary(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "metrics": {
            "leadsCaptured": 128,
            "leadsQualified": 37,
            "conversions": 12,
            "conversionRate": 0.3243,
        },
        "bySource": {
            "whatsapp": 46,
            "meta_ads": 39,
            "referral": 23,
            "landing_page": 20,
        },
        "backlog": {
            "pendingContact": 18,
            "runningAutomations": 7,
            "awaitingFinancialReview": 3,
        },
    }


def build_postgres_pipeline_summary(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        lead_metrics = fetch_lead_metrics(connection, tenant_slug)
        source_rows = fetch_source_rows(connection, tenant_slug)
        running_automations = fetch_running_automations(connection, tenant_slug)
        conversions = fetch_completed_runtime_executions(connection, tenant_slug)
        awaiting_financial_review = fetch_webhook_pending_review(connection)

    leads_captured = lead_metrics["total"]
    conversion_rate = round(conversions / leads_captured, 4) if leads_captured > 0 else 0.0

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "metrics": {
            "leadsCaptured": leads_captured,
            "leadsQualified": lead_metrics["qualified"],
            "conversions": conversions,
            "conversionRate": conversion_rate,
        },
        "bySource": {row["source"]: row["total"] for row in source_rows},
        "backlog": {
            "pendingContact": lead_metrics["captured"],
            "runningAutomations": running_automations,
            "awaitingFinancialReview": awaiting_financial_review,
        },
    }


def fetch_lead_metrics(connection: psycopg.Connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE lead.status = 'qualified') AS qualified,
            count(*) FILTER (WHERE lead.status = 'captured') AS captured
        FROM crm.leads AS lead
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {"total": 0, "qualified": 0, "captured": 0}

    return {
        "total": int(row["total"] or 0),
        "qualified": int(row["qualified"] or 0),
        "captured": int(row["captured"] or 0),
    }


def fetch_source_rows(connection: psycopg.Connection, tenant_slug: str | None) -> list[dict]:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            lead.source,
            count(*) AS total
        FROM crm.leads AS lead
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
        GROUP BY lead.source
        ORDER BY total DESC, lead.source ASC
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall()

    return [{"source": row["source"], "total": int(row["total"] or 0)} for row in rows]


def fetch_running_automations(connection: psycopg.Connection, tenant_slug: str | None) -> int:
    workflow_control_filter, workflow_control_params = tenant_filter("tenant.slug = %s", tenant_slug)
    workflow_runtime_filter, workflow_runtime_params = tenant_filter("tenant.slug = %s", tenant_slug)

    with connection.cursor() as cursor:
        cursor.execute(
            f"""
                SELECT count(*) AS total
                FROM workflow_control.workflow_runs AS run
                JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
                WHERE run.status = 'running'
                {append_filter(workflow_control_filter)}
            """,
            workflow_control_params,
        )
        workflow_control_row = cursor.fetchone() or {"total": 0}

        cursor.execute(
            f"""
                SELECT count(*) AS total
                FROM workflow_runtime.executions AS execution
                JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
                WHERE execution.status = 'running'
                {append_filter(workflow_runtime_filter)}
            """,
            workflow_runtime_params,
        )
        workflow_runtime_row = cursor.fetchone() or {"total": 0}

    return int(workflow_control_row["total"] or 0) + int(workflow_runtime_row["total"] or 0)


def fetch_completed_runtime_executions(connection: psycopg.Connection, tenant_slug: str | None) -> int:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT count(*) AS total
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        WHERE execution.status = 'completed'
        {append_filter(filter_sql)}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {"total": 0}

    return int(row["total"] or 0)


def fetch_webhook_pending_review(connection: psycopg.Connection) -> int:
    with connection.cursor() as cursor:
        cursor.execute(
            """
                SELECT count(*) AS total
                FROM webhook_hub.webhook_events
                WHERE status IN ('validated', 'queued', 'processing')
            """
        )
        row = cursor.fetchone() or {"total": 0}

    return int(row["total"] or 0)


def tenant_filter(condition: str, tenant_slug: str | None) -> tuple[str, list[str]]:
    if tenant_slug:
        return f"WHERE {condition}", [tenant_slug]
    return "", []


def append_filter(filter_sql: str) -> str:
    if not filter_sql:
        return ""
    return f" AND {filter_sql.removeprefix('WHERE ')}"
