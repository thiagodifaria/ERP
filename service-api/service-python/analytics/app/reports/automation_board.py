"""Painel operacional de automacao cruzando catalogo, ledger, runtime e entrega."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_automation_board(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_automation_board(tenant_slug)

    return build_static_automation_board(tenant_slug)


def build_static_automation_board(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "catalog": {
            "definitionsTotal": 8,
            "definitionsActive": 6,
            "definitionsDraft": 2,
            "publishedVersions": 14,
        },
        "control": {
            "runsTotal": 41,
            "pendingRuns": 3,
            "runningRuns": 7,
            "completedRuns": 28,
            "failedRuns": 2,
            "cancelledRuns": 1,
            "recordedEvents": 84,
            "byWorkflow": [
                {
                    "workflowDefinitionKey": "lead-follow-up",
                    "total": 21,
                    "running": 4,
                    "completed": 15,
                    "failed": 1,
                    "cancelled": 1,
                },
                {
                    "workflowDefinitionKey": "proposal-reminder",
                    "total": 20,
                    "running": 3,
                    "completed": 13,
                    "failed": 1,
                    "cancelled": 0,
                },
            ],
        },
        "runtime": {
            "executionsTotal": 44,
            "pendingExecutions": 4,
            "runningExecutions": 4,
            "completedExecutions": 28,
            "failedExecutions": 8,
            "cancelledExecutions": 4,
            "recordedTransitions": 116,
            "byWorkflow": [
                {
                    "workflowDefinitionKey": "lead-follow-up",
                    "total": 24,
                    "pending": 2,
                    "running": 2,
                    "completed": 16,
                    "failed": 3,
                    "cancelled": 1,
                    "retriesTotal": 3,
                },
                {
                    "workflowDefinitionKey": "proposal-reminder",
                    "total": 20,
                    "pending": 2,
                    "running": 2,
                    "completed": 12,
                    "failed": 5,
                    "cancelled": 3,
                    "retriesTotal": 5,
                },
            ],
        },
        "delivery": {
            "eventsTotal": 93,
            "validated": 1,
            "queued": 2,
            "processing": 1,
            "forwarded": 87,
            "failed": 3,
        },
    }


def build_postgres_automation_board(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        catalog_metrics = fetch_catalog_metrics(connection, tenant_slug)
        control_metrics = fetch_control_metrics(connection, tenant_slug)
        runtime_metrics = fetch_runtime_metrics(connection, tenant_slug)
        delivery_metrics = fetch_delivery_metrics(connection)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "catalog": catalog_metrics,
        "control": control_metrics,
        "runtime": runtime_metrics,
        "delivery": delivery_metrics,
    }


def fetch_catalog_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    definitions_query = f"""
        SELECT
            count(*) AS definitions_total,
            count(*) FILTER (WHERE definition.status = 'active') AS definitions_active,
            count(*) FILTER (WHERE definition.status = 'draft') AS definitions_draft
        FROM workflow_control.workflow_definitions AS definition
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        {filter_sql}
    """

    versions_query = f"""
        SELECT count(*) AS published_versions
        FROM workflow_control.workflow_definition_versions AS version
        JOIN workflow_control.workflow_definitions AS definition ON definition.id = version.workflow_definition_id
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(definitions_query, params)
        definitions_row = cursor.fetchone() or {}
        cursor.execute(versions_query, params)
        versions_row = cursor.fetchone() or {}

    return {
        "definitionsTotal": int(definitions_row.get("definitions_total", 0) or 0),
        "definitionsActive": int(definitions_row.get("definitions_active", 0) or 0),
        "definitionsDraft": int(definitions_row.get("definitions_draft", 0) or 0),
        "publishedVersions": int(versions_row.get("published_versions", 0) or 0),
    }


def fetch_control_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    runs_query = f"""
        SELECT
            count(*) AS runs_total,
            count(*) FILTER (WHERE run.status = 'pending') AS pending_runs,
            count(*) FILTER (WHERE run.status = 'running') AS running_runs,
            count(*) FILTER (WHERE run.status = 'completed') AS completed_runs,
            count(*) FILTER (WHERE run.status = 'failed') AS failed_runs,
            count(*) FILTER (WHERE run.status = 'cancelled') AS cancelled_runs
        FROM workflow_control.workflow_runs AS run
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    events_query = f"""
        SELECT count(*) AS recorded_events
        FROM workflow_control.workflow_run_events AS event
        JOIN workflow_control.workflow_runs AS run ON run.id = event.workflow_run_id
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(runs_query, params)
        runs_row = cursor.fetchone() or {}
        cursor.execute(events_query, params)
        events_row = cursor.fetchone() or {}

    return {
        "runsTotal": int(runs_row.get("runs_total", 0) or 0),
        "pendingRuns": int(runs_row.get("pending_runs", 0) or 0),
        "runningRuns": int(runs_row.get("running_runs", 0) or 0),
        "completedRuns": int(runs_row.get("completed_runs", 0) or 0),
        "failedRuns": int(runs_row.get("failed_runs", 0) or 0),
        "cancelledRuns": int(runs_row.get("cancelled_runs", 0) or 0),
        "recordedEvents": int(events_row.get("recorded_events", 0) or 0),
        "byWorkflow": fetch_control_by_workflow(connection, tenant_slug),
    }


def fetch_runtime_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    executions_query = f"""
        SELECT
            count(*) AS executions_total,
            count(*) FILTER (WHERE execution.status = 'pending') AS pending_executions,
            count(*) FILTER (WHERE execution.status = 'running') AS running_executions,
            count(*) FILTER (WHERE execution.status = 'completed') AS completed_executions,
            count(*) FILTER (WHERE execution.status = 'failed') AS failed_executions,
            count(*) FILTER (WHERE execution.status = 'cancelled') AS cancelled_executions
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
    """

    transitions_query = f"""
        SELECT count(*) AS recorded_transitions
        FROM workflow_runtime.execution_transitions AS transition
        JOIN workflow_runtime.executions AS execution ON execution.id = transition.execution_id
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(executions_query, params)
        executions_row = cursor.fetchone() or {}
        cursor.execute(transitions_query, params)
        transitions_row = cursor.fetchone() or {}

    return {
        "executionsTotal": int(executions_row.get("executions_total", 0) or 0),
        "pendingExecutions": int(executions_row.get("pending_executions", 0) or 0),
        "runningExecutions": int(executions_row.get("running_executions", 0) or 0),
        "completedExecutions": int(executions_row.get("completed_executions", 0) or 0),
        "failedExecutions": int(executions_row.get("failed_executions", 0) or 0),
        "cancelledExecutions": int(executions_row.get("cancelled_executions", 0) or 0),
        "recordedTransitions": int(transitions_row.get("recorded_transitions", 0) or 0),
        "byWorkflow": fetch_runtime_by_workflow(connection, tenant_slug),
    }


def fetch_control_by_workflow(connection, tenant_slug: str | None) -> list[dict]:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            definition.key AS workflow_definition_key,
            count(*) AS total,
            count(*) FILTER (WHERE run.status = 'running') AS running,
            count(*) FILTER (WHERE run.status = 'completed') AS completed,
            count(*) FILTER (WHERE run.status = 'failed') AS failed,
            count(*) FILTER (WHERE run.status = 'cancelled') AS cancelled
        FROM workflow_control.workflow_runs AS run
        JOIN workflow_control.workflow_definitions AS definition ON definition.id = run.workflow_definition_id
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
        GROUP BY definition.key
        ORDER BY definition.key ASC
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall() or []

    return [
        {
            "workflowDefinitionKey": row["workflow_definition_key"],
            "total": int(row.get("total", 0) or 0),
            "running": int(row.get("running", 0) or 0),
            "completed": int(row.get("completed", 0) or 0),
            "failed": int(row.get("failed", 0) or 0),
            "cancelled": int(row.get("cancelled", 0) or 0),
        }
        for row in rows
    ]


def fetch_runtime_by_workflow(connection, tenant_slug: str | None) -> list[dict]:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            execution.workflow_definition_key,
            count(*) AS total,
            count(*) FILTER (WHERE execution.status = 'pending') AS pending,
            count(*) FILTER (WHERE execution.status = 'running') AS running,
            count(*) FILTER (WHERE execution.status = 'completed') AS completed,
            count(*) FILTER (WHERE execution.status = 'failed') AS failed,
            count(*) FILTER (WHERE execution.status = 'cancelled') AS cancelled,
            coalesce(sum(execution.retry_count), 0) AS retries_total
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
        GROUP BY execution.workflow_definition_key
        ORDER BY execution.workflow_definition_key ASC
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall() or []

    return [
        {
            "workflowDefinitionKey": row["workflow_definition_key"],
            "total": int(row.get("total", 0) or 0),
            "pending": int(row.get("pending", 0) or 0),
            "running": int(row.get("running", 0) or 0),
            "completed": int(row.get("completed", 0) or 0),
            "failed": int(row.get("failed", 0) or 0),
            "cancelled": int(row.get("cancelled", 0) or 0),
            "retriesTotal": int(row.get("retries_total", 0) or 0),
        }
        for row in rows
    ]


def fetch_delivery_metrics(connection) -> dict:
    with connection.cursor() as cursor:
        cursor.execute(
            """
                SELECT
                    count(*) AS events_total,
                    count(*) FILTER (WHERE status = 'validated') AS validated,
                    count(*) FILTER (WHERE status = 'queued') AS queued,
                    count(*) FILTER (WHERE status = 'processing') AS processing,
                    count(*) FILTER (WHERE status = 'forwarded') AS forwarded,
                    count(*) FILTER (WHERE status = 'failed') AS failed
                FROM webhook_hub.webhook_events
            """
        )
        row = cursor.fetchone() or {}

    return {
        "eventsTotal": int(row.get("events_total", 0) or 0),
        "validated": int(row.get("validated", 0) or 0),
        "queued": int(row.get("queued", 0) or 0),
        "processing": int(row.get("processing", 0) or 0),
        "forwarded": int(row.get("forwarded", 0) or 0),
        "failed": int(row.get("failed", 0) or 0),
    }
