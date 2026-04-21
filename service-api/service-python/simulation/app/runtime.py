"""Motor de simulacao operacional e benchmark por carga."""

from __future__ import annotations

from datetime import datetime, timezone
import json
import math
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


_MEMORY_SCENARIO_RUNS: list[dict] = []
_MEMORY_BENCHMARK_RUNS: list[dict] = []


def build_catalog() -> dict:
    return {
        "service": settings.service_name,
        "scenarioTypes": [
            {
                "key": "operational-load",
                "name": "Operational load",
                "description": "Projects leads, workflows, webhooks and storage pressure from current tenant activity.",
                "inputs": [
                    "tenantSlug",
                    "leadMultiplier",
                    "automationMultiplier",
                    "webhookMultiplier",
                    "teamSize",
                    "planningHorizonDays",
                    "storageGrowthGb",
                ],
            },
            {
                "key": "load-benchmark",
                "name": "Load benchmark",
                "description": "Estimates latency, throughput and resource pressure for the projected workload.",
                "inputs": [
                    "tenantSlug",
                    "leadMultiplier",
                    "automationMultiplier",
                    "webhookMultiplier",
                    "sampleSize",
                ],
            },
        ],
    }


def create_operational_load_scenario(payload: dict) -> dict:
    tenant_slug = str(payload.get("tenantSlug") or "global").strip() or "global"
    lead_multiplier = normalize_multiplier(payload.get("leadMultiplier"), default=1.0)
    automation_multiplier = normalize_multiplier(payload.get("automationMultiplier"), default=1.0)
    webhook_multiplier = normalize_multiplier(payload.get("webhookMultiplier"), default=1.0)
    planning_horizon_days = max(int(payload.get("planningHorizonDays") or 30), 1)
    team_size = max(int(payload.get("teamSize") or 1), 1)
    storage_growth_gb = max(int(payload.get("storageGrowthGb") or 0), 0)

    base_metrics = fetch_base_metrics(tenant_slug)
    projected = project_operational_load(
        base_metrics,
        lead_multiplier,
        automation_multiplier,
        webhook_multiplier,
        planning_horizon_days,
        team_size,
        storage_growth_gb,
    )

    record = {
        "publicId": str(uuid.uuid4()),
        "scenarioKey": "operational-load",
        "tenantSlug": tenant_slug,
        "createdAt": datetime.now(timezone.utc).isoformat(),
        "inputs": {
            "leadMultiplier": lead_multiplier,
            "automationMultiplier": automation_multiplier,
            "webhookMultiplier": webhook_multiplier,
            "planningHorizonDays": planning_horizon_days,
            "teamSize": team_size,
            "storageGrowthGb": storage_growth_gb,
        },
        "baseline": base_metrics,
        "projection": projected,
    }

    persist_scenario_record(record)
    return record


def list_operational_load_scenarios(tenant_slug: str | None = None) -> list[dict]:
    if settings.repository_driver == "postgres":
        with connect() as connection:
            where_clause = ""
            params: list[str] = []

            if tenant_slug:
                where_clause = "WHERE tenant.slug = %s"
                params.append(tenant_slug)

            with connection.cursor() as cursor:
                cursor.execute(
                    f"""
                        SELECT
                            scenario.public_id::text AS public_id,
                            tenant.slug AS tenant_slug,
                            scenario.scenario_key,
                            scenario.input_payload,
                            scenario.output_payload,
                            to_char(scenario.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SSOF') AS created_at
                        FROM simulation.scenario_runs AS scenario
                        LEFT JOIN identity.tenants AS tenant ON tenant.id = scenario.tenant_id
                        {where_clause}
                        ORDER BY scenario.created_at DESC
                    """,
                    params,
                )
                rows = cursor.fetchall()

        return [restore_scenario_record(row) for row in rows]

    records = _MEMORY_SCENARIO_RUNS
    if tenant_slug:
        records = [record for record in records if record["tenantSlug"] == tenant_slug]
    return list(reversed(records))


def get_operational_load_scenario(public_id: str) -> dict | None:
    if settings.repository_driver == "postgres":
        with connect() as connection:
            with connection.cursor() as cursor:
                cursor.execute(
                    """
                        SELECT
                            scenario.public_id::text AS public_id,
                            tenant.slug AS tenant_slug,
                            scenario.scenario_key,
                            scenario.input_payload,
                            scenario.output_payload,
                            to_char(scenario.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SSOF') AS created_at
                        FROM simulation.scenario_runs AS scenario
                        LEFT JOIN identity.tenants AS tenant ON tenant.id = scenario.tenant_id
                        WHERE scenario.public_id = %s::uuid
                    """,
                    [public_id],
                )
                row = cursor.fetchone()

        if not row:
            return None

        return restore_scenario_record(row)

    for record in reversed(_MEMORY_SCENARIO_RUNS):
        if record["publicId"] == public_id:
            return record
    return None


def create_load_benchmark(payload: dict) -> dict:
    tenant_slug = str(payload.get("tenantSlug") or "global").strip() or "global"
    lead_multiplier = normalize_multiplier(payload.get("leadMultiplier"), default=1.0)
    automation_multiplier = normalize_multiplier(payload.get("automationMultiplier"), default=1.0)
    webhook_multiplier = normalize_multiplier(payload.get("webhookMultiplier"), default=1.0)
    sample_size = max(int(payload.get("sampleSize") or 120), 30)

    base_metrics = fetch_base_metrics(tenant_slug)
    projection = project_operational_load(
        base_metrics,
        lead_multiplier,
        automation_multiplier,
        webhook_multiplier,
        30,
        max(int(payload.get("teamSize") or 1), 1),
        max(int(payload.get("storageGrowthGb") or 0), 0),
    )

    avg_latency_ms = int(
        55
        + projection["workflowRunsProjected"] / 18
        + projection["webhookEventsProjected"] / 20
        + max(projection["teamCapacityGap"], 0) * 14
    )
    p95_latency_ms = int(avg_latency_ms * 1.85)
    max_latency_ms = int(p95_latency_ms * 1.35)
    throughput_rps = round(
        max(
            8.0,
            280.0
            / max(
                1.0,
                lead_multiplier + automation_multiplier + (webhook_multiplier / 2),
            ),
        ),
        2,
    )
    cpu_load_percent = min(
        96,
        int(
            18
            + projection["workflowRunsProjected"] / 12
            + projection["webhookEventsProjected"] / 18
            + max(projection["teamCapacityGap"], 0) * 7
        ),
    )
    memory_mb = 256 + int(projection["storageProjectedMb"] / 4) + int(projection["workflowRunsProjected"] / 3)
    status = classify_benchmark_status(p95_latency_ms, throughput_rps, cpu_load_percent)

    record = {
        "publicId": str(uuid.uuid4()),
        "benchmarkKey": "load-benchmark",
        "tenantSlug": tenant_slug,
        "createdAt": datetime.now(timezone.utc).isoformat(),
        "inputs": {
            "leadMultiplier": lead_multiplier,
            "automationMultiplier": automation_multiplier,
            "webhookMultiplier": webhook_multiplier,
            "sampleSize": sample_size,
        },
        "baseline": base_metrics,
        "results": {
            "sampleSize": sample_size,
            "avgLatencyMs": avg_latency_ms,
            "p95LatencyMs": p95_latency_ms,
            "maxLatencyMs": max_latency_ms,
            "throughputRps": throughput_rps,
            "cpuLoadPercent": cpu_load_percent,
            "memoryMb": memory_mb,
            "status": status,
        },
    }

    persist_benchmark_record(record)
    return record


def list_load_benchmarks(tenant_slug: str | None = None) -> list[dict]:
    if settings.repository_driver == "postgres":
        with connect() as connection:
            where_clause = ""
            params: list[str] = []

            if tenant_slug:
                where_clause = "WHERE tenant.slug = %s"
                params.append(tenant_slug)

            with connection.cursor() as cursor:
                cursor.execute(
                    f"""
                        SELECT
                            benchmark.public_id::text AS public_id,
                            tenant.slug AS tenant_slug,
                            benchmark.benchmark_key,
                            benchmark.input_payload,
                            benchmark.output_payload,
                            to_char(benchmark.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SSOF') AS created_at
                        FROM simulation.load_benchmark_runs AS benchmark
                        LEFT JOIN identity.tenants AS tenant ON tenant.id = benchmark.tenant_id
                        {where_clause}
                        ORDER BY benchmark.created_at DESC
                    """,
                    params,
                )
                rows = cursor.fetchall()

        return [restore_benchmark_record(row) for row in rows]

    records = _MEMORY_BENCHMARK_RUNS
    if tenant_slug:
        records = [record for record in records if record["tenantSlug"] == tenant_slug]
    return list(reversed(records))


def persist_scenario_record(record: dict) -> None:
    if settings.repository_driver == "postgres":
        with connect() as connection:
            tenant_id = resolve_tenant_id(connection, record["tenantSlug"])
            with connection.cursor() as cursor:
                cursor.execute(
                    """
                        INSERT INTO simulation.scenario_runs (tenant_id, public_id, scenario_key, input_payload, output_payload)
                        VALUES (%s, %s::uuid, %s, %s::jsonb, %s::jsonb)
                    """,
                    [
                        tenant_id,
                        record["publicId"],
                        record["scenarioKey"],
                        json.dumps(record["inputs"]),
                        json.dumps(
                            {
                                "baseline": record["baseline"],
                                "projection": record["projection"],
                            }
                        ),
                    ],
                )
            connection.commit()
        return

    _MEMORY_SCENARIO_RUNS.append(record)


def persist_benchmark_record(record: dict) -> None:
    if settings.repository_driver == "postgres":
        with connect() as connection:
            tenant_id = resolve_tenant_id(connection, record["tenantSlug"])
            with connection.cursor() as cursor:
                cursor.execute(
                    """
                        INSERT INTO simulation.load_benchmark_runs (tenant_id, public_id, benchmark_key, input_payload, output_payload)
                        VALUES (%s, %s::uuid, %s, %s::jsonb, %s::jsonb)
                    """,
                    [
                        tenant_id,
                        record["publicId"],
                        record["benchmarkKey"],
                        json.dumps(record["inputs"]),
                        json.dumps(
                            {
                                "baseline": record["baseline"],
                                "results": record["results"],
                            }
                        ),
                    ],
                )
            connection.commit()
        return

    _MEMORY_BENCHMARK_RUNS.append(record)


def fetch_base_metrics(tenant_slug: str) -> dict:
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": tenant_slug,
            "leads": 128,
            "sales": 12,
            "workflowRuns": 41,
            "runtimeExecutions": 28,
            "webhookEvents": 93,
            "attachments": 37,
            "activeUsers": 9,
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                    SELECT
                        tenant.slug,
                        tenant.id,
                        (SELECT count(*) FROM crm.leads AS lead WHERE lead.tenant_id = tenant.id) AS leads,
                        (SELECT count(*) FROM sales.sales AS sale WHERE sale.tenant_id = tenant.id) AS sales,
                        (SELECT count(*) FROM workflow_control.workflow_runs AS workflow_run WHERE workflow_run.tenant_id = tenant.id) AS workflow_runs,
                        (SELECT count(*) FROM workflow_runtime.executions AS execution WHERE execution.tenant_id = tenant.id) AS runtime_executions,
                        (SELECT count(*) FROM documents.attachments AS attachment WHERE attachment.tenant_id = tenant.id) AS attachments,
                        (SELECT count(*) FROM identity.users AS "user" WHERE "user".tenant_id = tenant.id AND "user".status = 'active') AS active_users
                    FROM identity.tenants AS tenant
                    WHERE tenant.slug = %s
                """,
                [tenant_slug],
            )
            row = cursor.fetchone()

            cursor.execute("SELECT count(*) AS total FROM webhook_hub.webhook_events")
            webhook_row = cursor.fetchone() or {"total": 0}

    if not row:
        return {
            "tenantSlug": tenant_slug,
            "leads": 0,
            "sales": 0,
            "workflowRuns": 0,
            "runtimeExecutions": 0,
            "webhookEvents": int(webhook_row.get("total", 0) or 0),
            "attachments": 0,
            "activeUsers": 0,
        }

    return {
        "tenantSlug": row["slug"],
        "leads": int(row["leads"] or 0),
        "sales": int(row["sales"] or 0),
        "workflowRuns": int(row["workflow_runs"] or 0),
        "runtimeExecutions": int(row["runtime_executions"] or 0),
        "webhookEvents": int(webhook_row.get("total", 0) or 0),
        "attachments": int(row["attachments"] or 0),
        "activeUsers": int(row["active_users"] or 0),
    }


def resolve_tenant_id(connection, tenant_slug: str) -> int | None:
    with connection.cursor() as cursor:
        cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", [tenant_slug])
        row = cursor.fetchone()

    if not row:
        return None

    return int(row["id"])


def project_operational_load(
    base_metrics: dict,
    lead_multiplier: float,
    automation_multiplier: float,
    webhook_multiplier: float,
    planning_horizon_days: int,
    team_size: int,
    storage_growth_gb: int,
) -> dict:
    leads_projected = max(1, math.ceil(max(base_metrics["leads"], 1) * lead_multiplier))
    workflow_runs_projected = max(
        1,
        math.ceil(max(base_metrics["workflowRuns"] + base_metrics["runtimeExecutions"], 1) * automation_multiplier),
    )
    webhook_events_projected = max(1, math.ceil(max(base_metrics["webhookEvents"], 1) * webhook_multiplier))
    attachments_projected = max(1, math.ceil(max(base_metrics["attachments"], 1) * max(lead_multiplier, 1.0)))
    storage_projected_mb = attachments_projected * 4 + (storage_growth_gb * 1024)
    monthly_operations = leads_projected + workflow_runs_projected + webhook_events_projected

    operator_demand = math.ceil(
        (leads_projected / 40)
        + (workflow_runs_projected / 180)
        + (webhook_events_projected / 260)
    )
    required_team_capacity = max(team_size, operator_demand)
    team_capacity_gap = max(required_team_capacity - team_size, 0)

    infra_cost_cents = 15000 + (workflow_runs_projected * 12) + (webhook_events_projected * 6)
    storage_cost_cents = storage_projected_mb * 3
    support_cost_cents = required_team_capacity * 4500
    estimated_monthly_cost_cents = infra_cost_cents + storage_cost_cents + support_cost_cents

    return {
        "planningHorizonDays": planning_horizon_days,
        "leadsProjected": leads_projected,
        "workflowRunsProjected": workflow_runs_projected,
        "webhookEventsProjected": webhook_events_projected,
        "attachmentsProjected": attachments_projected,
        "storageProjectedMb": storage_projected_mb,
        "monthlyOperationsProjected": monthly_operations,
        "requiredTeamCapacity": required_team_capacity,
        "teamCapacityGap": team_capacity_gap,
        "costBreakdown": {
            "infraCostCents": infra_cost_cents,
            "storageCostCents": storage_cost_cents,
            "supportCostCents": support_cost_cents,
            "estimatedMonthlyCostCents": estimated_monthly_cost_cents,
        },
        "risk": classify_scenario_risk(monthly_operations, team_capacity_gap, storage_projected_mb),
    }


def classify_scenario_risk(monthly_operations: int, team_capacity_gap: int, storage_projected_mb: int) -> str:
    if monthly_operations >= 900 or team_capacity_gap >= 3 or storage_projected_mb >= 4096:
        return "critical"
    if monthly_operations >= 400 or team_capacity_gap >= 1 or storage_projected_mb >= 2048:
        return "attention"
    return "stable"


def classify_benchmark_status(p95_latency_ms: int, throughput_rps: float, cpu_load_percent: int) -> str:
    if p95_latency_ms >= 420 or throughput_rps <= 12 or cpu_load_percent >= 90:
        return "critical"
    if p95_latency_ms >= 250 or throughput_rps <= 24 or cpu_load_percent >= 75:
        return "attention"
    return "stable"


def normalize_multiplier(raw_value: object, default: float) -> float:
    try:
        value = float(raw_value)
    except (TypeError, ValueError):
        return default

    return max(value, 0.1)


def restore_scenario_record(row: dict) -> dict:
    payload = normalize_json_payload(row.get("output_payload"))
    return {
        "publicId": row["public_id"],
        "scenarioKey": row["scenario_key"],
        "tenantSlug": row.get("tenant_slug") or "global",
        "createdAt": row["created_at"],
        "inputs": row.get("input_payload") or {},
        "baseline": payload.get("baseline", {}),
        "projection": payload.get("projection", {}),
    }


def restore_benchmark_record(row: dict) -> dict:
    payload = normalize_json_payload(row.get("output_payload"))
    return {
        "publicId": row["public_id"],
        "benchmarkKey": row["benchmark_key"],
        "tenantSlug": row.get("tenant_slug") or "global",
        "createdAt": row["created_at"],
        "inputs": row.get("input_payload") or {},
        "baseline": payload.get("baseline", {}),
        "results": payload.get("results", {}),
    }


def normalize_json_payload(payload: object) -> dict:
    if isinstance(payload, dict):
        return payload
    if isinstance(payload, str):
        try:
            loaded = json.loads(payload)
        except json.JSONDecodeError:
            return {}
        if isinstance(loaded, dict):
            return loaded
    return {}
