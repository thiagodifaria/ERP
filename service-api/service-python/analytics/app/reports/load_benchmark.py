"""Relatorio consolidado do benchmark por carga."""

from datetime import datetime, timezone
import json

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_load_benchmark(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_load_benchmark(tenant_slug)

    return build_static_load_benchmark(tenant_slug)


def build_static_load_benchmark(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "summary": {"totalRuns": 2, "stable": 1, "attention": 1, "critical": 0},
        "latest": {
            "publicId": "00000000-0000-0000-0000-00000000lb01",
            "status": "attention",
            "avgLatencyMs": 182,
            "p95LatencyMs": 332,
            "throughputRps": 46.8,
            "cpuLoadPercent": 74,
            "memoryMb": 612,
        },
        "recent": [],
    }


def build_postgres_load_benchmark(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        params: list[str] = []
        where_clause = ""

        if tenant_slug:
            where_clause = "WHERE tenant.slug = %s"
            params.append(tenant_slug)

        with connection.cursor() as cursor:
            cursor.execute(
                f"""
                    SELECT
                        benchmark.public_id::text AS public_id,
                        benchmark.output_payload,
                        to_char(benchmark.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SSOF') AS created_at
                    FROM simulation.load_benchmark_runs AS benchmark
                    LEFT JOIN identity.tenants AS tenant ON tenant.id = benchmark.tenant_id
                    {where_clause}
                    ORDER BY benchmark.created_at DESC
                    LIMIT 5
                """,
                params,
            )
            rows = cursor.fetchall()

    recent = []
    stable = 0
    attention = 0
    critical = 0

    for row in rows:
        results = normalize_payload(row.get("output_payload")).get("results", {})
        status = results.get("status", "stable")

        if status == "critical":
            critical += 1
        elif status == "attention":
            attention += 1
        else:
            stable += 1

        recent.append(
            {
                "publicId": row["public_id"],
                "status": status,
                "avgLatencyMs": int(results.get("avgLatencyMs", 0) or 0),
                "p95LatencyMs": int(results.get("p95LatencyMs", 0) or 0),
                "throughputRps": float(results.get("throughputRps", 0.0) or 0.0),
                "cpuLoadPercent": int(results.get("cpuLoadPercent", 0) or 0),
                "memoryMb": int(results.get("memoryMb", 0) or 0),
                "createdAt": row["created_at"],
            }
        )

    latest = recent[0] if recent else None

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "summary": {
            "totalRuns": len(recent),
            "stable": stable,
            "attention": attention,
            "critical": critical,
        },
        "latest": latest,
        "recent": recent,
    }


def normalize_payload(payload: object) -> dict:
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
