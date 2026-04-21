"""Relatorio de custo estimado por cenario de simulacao."""

from datetime import datetime, timezone
import json

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_cost_estimator(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_cost_estimator(tenant_slug)

    return build_static_cost_estimator(tenant_slug)


def build_static_cost_estimator(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "latestScenarioRunPublicId": "00000000-0000-0000-0000-00000000ce01",
        "costs": {
            "infraCostCents": 22400,
            "storageCostCents": 9420,
            "supportCostCents": 27000,
            "estimatedMonthlyCostCents": 58820,
        },
        "projection": {
            "monthlyOperationsProjected": 361,
            "requiredTeamCapacity": 6,
            "teamCapacityGap": 1,
            "storageProjectedMb": 3140,
            "risk": "attention",
        },
        "recommendations": {
            "plan": "Scale analytics and workflow throughput before expanding outbound volume.",
            "needsTeamExpansion": True,
            "suggestedStorageTier": "warm-object-storage",
        },
    }


def build_postgres_cost_estimator(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        with connection.cursor() as cursor:
            params: list[str] = []
            where_clause = ""

            if tenant_slug:
                where_clause = "WHERE tenant.slug = %s"
                params.append(tenant_slug)

            cursor.execute(
                f"""
                    SELECT
                        scenario.public_id::text AS public_id,
                        scenario.output_payload
                    FROM simulation.scenario_runs AS scenario
                    LEFT JOIN identity.tenants AS tenant ON tenant.id = scenario.tenant_id
                    {where_clause}
                    ORDER BY scenario.created_at DESC
                    LIMIT 1
                """,
                params,
            )
            row = cursor.fetchone()

    payload = normalize_payload(row.get("output_payload")) if row else {}
    projection = payload.get("projection", {})
    costs = projection.get("costBreakdown", {})

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "latestScenarioRunPublicId": row.get("public_id") if row else None,
        "costs": {
            "infraCostCents": int(costs.get("infraCostCents", 0) or 0),
            "storageCostCents": int(costs.get("storageCostCents", 0) or 0),
            "supportCostCents": int(costs.get("supportCostCents", 0) or 0),
            "estimatedMonthlyCostCents": int(costs.get("estimatedMonthlyCostCents", 0) or 0),
        },
        "projection": {
            "monthlyOperationsProjected": int(projection.get("monthlyOperationsProjected", 0) or 0),
            "requiredTeamCapacity": int(projection.get("requiredTeamCapacity", 0) or 0),
            "teamCapacityGap": int(projection.get("teamCapacityGap", 0) or 0),
            "storageProjectedMb": int(projection.get("storageProjectedMb", 0) or 0),
            "risk": projection.get("risk", "stable"),
        },
        "recommendations": {
            "plan": build_plan_recommendation(projection),
            "needsTeamExpansion": int(projection.get("teamCapacityGap", 0) or 0) > 0,
            "suggestedStorageTier": build_storage_tier(int(projection.get("storageProjectedMb", 0) or 0)),
        },
    }


def build_plan_recommendation(projection: dict) -> str:
    risk = projection.get("risk", "stable")
    team_gap = int(projection.get("teamCapacityGap", 0) or 0)

    if risk == "critical":
        return "Prioritize automation scaling and expand the operational team before taking the next demand wave."
    if team_gap > 0:
        return "Expand the operational team and storage envelope before unlocking the projected load."
    return "Current platform footprint supports the projected scenario without immediate expansion."


def build_storage_tier(storage_projected_mb: int) -> str:
    if storage_projected_mb >= 4096:
        return "cold-object-storage"
    if storage_projected_mb >= 2048:
        return "warm-object-storage"
    return "standard-object-storage"


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
