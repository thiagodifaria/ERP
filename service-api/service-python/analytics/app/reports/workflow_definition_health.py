"""Saude operacional por definicao de workflow."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_workflow_definition_health(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_workflow_definition_health(tenant_slug)

    return build_static_workflow_definition_health(tenant_slug)


def build_static_workflow_definition_health(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    definitions = [
        build_definition_payload(
            workflow_definition_key="lead-follow-up",
            name="Lead Follow-Up",
            trigger="lead.created",
            definition_status="active",
            published_versions=3,
            current_version_number=3,
            current_version_status="active",
            control={
                "runsTotal": 21,
                "pending": 1,
                "running": 4,
                "completed": 15,
                "failed": 1,
                "cancelled": 0,
            },
            runtime={
                "executionsTotal": 24,
                "pending": 2,
                "running": 2,
                "completed": 16,
                "failed": 3,
                "cancelled": 1,
                "retriesTotal": 3,
            },
        ),
        build_definition_payload(
            workflow_definition_key="proposal-reminder",
            name="Proposal Reminder",
            trigger="proposal.created",
            definition_status="draft",
            published_versions=1,
            current_version_number=1,
            current_version_status="active",
            control={
                "runsTotal": 20,
                "pending": 2,
                "running": 3,
                "completed": 13,
                "failed": 1,
                "cancelled": 1,
            },
            runtime={
                "executionsTotal": 20,
                "pending": 2,
                "running": 2,
                "completed": 12,
                "failed": 4,
                "cancelled": 0,
                "retriesTotal": 5,
            },
        ),
    ]

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "summary": build_summary(definitions),
        "definitions": definitions,
    }


def build_postgres_workflow_definition_health(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        definition_rows = fetch_definition_rows(connection, tenant_slug)
        control_by_workflow = fetch_control_by_workflow(connection, tenant_slug)
        runtime_by_workflow = fetch_runtime_by_workflow(connection, tenant_slug)

    control_index = {row["workflowDefinitionKey"]: row for row in control_by_workflow}
    runtime_index = {row["workflowDefinitionKey"]: row for row in runtime_by_workflow}

    definitions = [
        build_definition_payload(
            workflow_definition_key=row["workflow_definition_key"],
            name=row["name"],
            trigger=row["trigger"],
            definition_status=row["definition_status"],
            published_versions=int(row.get("published_versions", 0) or 0),
            current_version_number=row.get("current_version_number"),
            current_version_status=row.get("current_version_status"),
            control=control_index.get(
                row["workflow_definition_key"],
                empty_control_metrics(),
            ),
            runtime=runtime_index.get(
                row["workflow_definition_key"],
                empty_runtime_metrics(),
            ),
        )
        for row in definition_rows
    ]

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "summary": build_summary(definitions),
        "definitions": definitions,
    }


def fetch_definition_rows(connection, tenant_slug: str | None) -> list[dict]:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            definition.key AS workflow_definition_key,
            definition.name,
            definition.trigger,
            definition.status AS definition_status,
            count(version.id) AS published_versions,
            current_version.version_number AS current_version_number,
            current_version.snapshot_status AS current_version_status
        FROM workflow_control.workflow_definitions AS definition
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        LEFT JOIN workflow_control.workflow_definition_versions AS version
            ON version.workflow_definition_id = definition.id
        LEFT JOIN LATERAL (
            SELECT
                version_number,
                snapshot_status
            FROM workflow_control.workflow_definition_versions
            WHERE workflow_definition_id = definition.id
            ORDER BY version_number DESC
            LIMIT 1
        ) AS current_version ON true
        {filter_sql}
        GROUP BY
            definition.id,
            definition.key,
            definition.name,
            definition.trigger,
            definition.status,
            current_version.version_number,
            current_version.snapshot_status
        ORDER BY definition.key ASC
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        return cursor.fetchall() or []


def fetch_control_by_workflow(connection, tenant_slug: str | None) -> list[dict]:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            definition.key AS workflow_definition_key,
            count(*) AS runs_total,
            count(*) FILTER (WHERE run.status = 'pending') AS pending,
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
            "runsTotal": int(row.get("runs_total", 0) or 0),
            "pending": int(row.get("pending", 0) or 0),
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
            count(*) AS executions_total,
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
            "executionsTotal": int(row.get("executions_total", 0) or 0),
            "pending": int(row.get("pending", 0) or 0),
            "running": int(row.get("running", 0) or 0),
            "completed": int(row.get("completed", 0) or 0),
            "failed": int(row.get("failed", 0) or 0),
            "cancelled": int(row.get("cancelled", 0) or 0),
            "retriesTotal": int(row.get("retries_total", 0) or 0),
        }
        for row in rows
    ]


def build_definition_payload(
    workflow_definition_key: str,
    name: str,
    trigger: str,
    definition_status: str,
    published_versions: int,
    current_version_number: int | None,
    current_version_status: str | None,
    control: dict,
    runtime: dict,
) -> dict:
    health, reasons = classify_health(
        definition_status=definition_status,
        published_versions=published_versions,
        current_version_status=current_version_status,
    )

    return {
        "workflowDefinitionKey": workflow_definition_key,
        "name": name,
        "trigger": trigger,
        "definitionStatus": definition_status,
        "publishedVersions": published_versions,
        "currentVersionNumber": current_version_number,
        "currentVersionStatus": current_version_status,
        "health": health,
        "attentionReasons": reasons,
        "control": {
            "runsTotal": int(control.get("runsTotal", 0) or 0),
            "pending": int(control.get("pending", 0) or 0),
            "running": int(control.get("running", 0) or 0),
            "completed": int(control.get("completed", 0) or 0),
            "failed": int(control.get("failed", 0) or 0),
            "cancelled": int(control.get("cancelled", 0) or 0),
        },
        "runtime": {
            "executionsTotal": int(runtime.get("executionsTotal", 0) or 0),
            "pending": int(runtime.get("pending", 0) or 0),
            "running": int(runtime.get("running", 0) or 0),
            "completed": int(runtime.get("completed", 0) or 0),
            "failed": int(runtime.get("failed", 0) or 0),
            "cancelled": int(runtime.get("cancelled", 0) or 0),
            "retriesTotal": int(runtime.get("retriesTotal", 0) or 0),
        },
    }


def classify_health(
    definition_status: str,
    published_versions: int,
    current_version_status: str | None,
) -> tuple[str, list[str]]:
    if published_versions == 0 or current_version_status is None:
        return "critical", ["published-version-missing"]

    reasons: list[str] = []

    if definition_status != "active":
        reasons.append("definition-not-active")

    if current_version_status != "active":
        reasons.append("current-version-not-active")

    if reasons:
        return "attention", reasons

    return "stable", []


def build_summary(definitions: list[dict]) -> dict:
    return {
        "definitionsTotal": len(definitions),
        "stable": sum(1 for definition in definitions if definition["health"] == "stable"),
        "attention": sum(1 for definition in definitions if definition["health"] == "attention"),
        "critical": sum(1 for definition in definitions if definition["health"] == "critical"),
        "withPublishedVersion": sum(1 for definition in definitions if definition["publishedVersions"] > 0),
        "withoutPublishedVersion": sum(1 for definition in definitions if definition["publishedVersions"] == 0),
    }


def empty_control_metrics() -> dict:
    return {
        "runsTotal": 0,
        "pending": 0,
        "running": 0,
        "completed": 0,
        "failed": 0,
        "cancelled": 0,
    }


def empty_runtime_metrics() -> dict:
    return {
        "executionsTotal": 0,
        "pending": 0,
        "running": 0,
        "completed": 0,
        "failed": 0,
        "cancelled": 0,
        "retriesTotal": 0,
    }
