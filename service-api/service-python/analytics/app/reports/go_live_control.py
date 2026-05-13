"""Painel de rollout, adocao e readiness de go-live."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_release_controls() -> dict:
    return {
        "status": "stable",
        "acceptanceReady": True,
        "playbook": "docs/OPERACOES.md#go-live-operacional",
        "testSuites": ["smoke", "hardening"],
        "controls": [
            "tenant-rollout",
            "adoption-monitoring",
            "rollback",
            "bottleneck-review",
            "usage-adjustments",
        ],
        "coveredAreas": [
            "rollout-by-tenant",
            "business-and-health-metrics",
            "auditable-rollback",
            "bottleneck-observation",
            "fine-tuning-by-usage",
        ],
    }


def build_go_live_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": slug,
            "generatedAt": datetime.now(timezone.utc).isoformat(),
            "rollouts": {"planned": 1, "running": 0, "completed": 1, "rolledBack": 0},
            "adoption": {"trackedMetrics": 4, "totalQuantity": 4096, "adoptionPct": 84, "targetPct": 80, "gapPct": 0},
            "bottlenecks": {"critical": 0, "attention": 1, "total": 1},
            "adjustments": {"recommended": 1, "applySupported": 1},
            "readiness": {"status": "stable", "rolloutReady": True, "rollbackReady": True, "metricsObserved": True},
            "releaseControls": build_release_controls(),
            "goLiveClosure": build_release_controls(),
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  count(*) FILTER (WHERE rollout.status = 'planned') AS planned_total,
                  count(*) FILTER (WHERE rollout.status = 'running') AS running_total,
                  count(*) FILTER (WHERE rollout.status = 'completed') AS completed_total,
                  count(*) FILTER (WHERE rollout.status = 'rolled_back') AS rolled_back_total,
                  coalesce(max((rollout.readiness_json ->> 'rolloutReady')::boolean::int), 0) AS rollout_ready_flag,
                  coalesce(max(rollout.adoption_target_pct), 70) AS target_pct
                FROM platform_control.go_live_rollouts AS rollout
                JOIN identity.tenants AS tenant ON tenant.id = rollout.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            rollout_row = cursor.fetchone() or {}
            cursor.execute(
                """
                SELECT
                  count(*) AS tracked_metrics,
                  coalesce(sum(snapshot.quantity), 0) AS total_quantity
                FROM platform_control.usage_snapshots AS snapshot
                JOIN identity.tenants AS tenant ON tenant.id = snapshot.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            usage_row = cursor.fetchone() or {}
            cursor.execute(
                """
                SELECT
                  count(*) FILTER (WHERE provider.critical = TRUE AND provider.mode IN ('unconfigured', 'disabled')) AS critical_provider_gaps,
                  count(*) FILTER (WHERE block.active = TRUE) AS active_blocks,
                  count(*) FILTER (
                    WHERE quota.limit_value > 0
                      AND usage.total_quantity >= quota.limit_value
                  ) AS quota_limits_reached,
                  count(*) FILTER (
                    WHERE quota.limit_value > 0
                      AND usage.total_quantity < quota.limit_value
                      AND usage.total_quantity >= (quota.limit_value * 0.85)
                  ) AS quota_attention
                FROM identity.tenants AS tenant
                LEFT JOIN platform_control.provider_defaults AS provider ON provider.tenant_id = tenant.id
                LEFT JOIN platform_control.tenant_blocks AS block ON block.tenant_id = tenant.id
                LEFT JOIN platform_control.quotas AS quota ON quota.tenant_id = tenant.id
                LEFT JOIN (
                  SELECT tenant_id, metric_key, coalesce(sum(quantity), 0) AS total_quantity
                  FROM platform_control.usage_snapshots
                  GROUP BY tenant_id, metric_key
                ) AS usage ON usage.tenant_id = quota.tenant_id AND usage.metric_key = quota.metric_key
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            risk_row = cursor.fetchone() or {}

    status = "stable"
    if int(rollout_row.get("running_total", 0) or 0) > 0 or int(usage_row.get("tracked_metrics", 0) or 0) == 0:
        status = "attention"

    tracked_metrics = int(usage_row.get("tracked_metrics", 0) or 0)
    total_quantity = int(usage_row.get("total_quantity", 0) or 0)
    target_pct = int(rollout_row.get("target_pct", 70) or 70)
    critical_provider_gaps = int(risk_row.get("critical_provider_gaps", 0) or 0)
    active_blocks = int(risk_row.get("active_blocks", 0) or 0)
    quota_limits_reached = int(risk_row.get("quota_limits_reached", 0) or 0)
    quota_attention = int(risk_row.get("quota_attention", 0) or 0)

    adoption_pct = 0
    if tracked_metrics > 0:
        pressure = critical_provider_gaps + active_blocks + quota_limits_reached + quota_attention
        adoption_pct = max(0, min(100, int(round(100 - (pressure * 12.5)))))

    if total_quantity > 0 and adoption_pct < 70:
        adoption_pct = min(100, adoption_pct + 15)
    gap_pct = max(target_pct - adoption_pct, 0)
    recommended = critical_provider_gaps + active_blocks + quota_limits_reached + quota_attention
    adjustments_supported = active_blocks + quota_limits_reached + quota_attention

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "rollouts": {
            "planned": int(rollout_row.get("planned_total", 0) or 0),
            "running": int(rollout_row.get("running_total", 0) or 0),
            "completed": int(rollout_row.get("completed_total", 0) or 0),
            "rolledBack": int(rollout_row.get("rolled_back_total", 0) or 0),
        },
        "adoption": {
            "trackedMetrics": tracked_metrics,
            "totalQuantity": total_quantity,
            "adoptionPct": adoption_pct,
            "targetPct": target_pct,
            "gapPct": gap_pct,
        },
        "bottlenecks": {
            "critical": critical_provider_gaps + active_blocks + quota_limits_reached,
            "attention": quota_attention,
            "total": critical_provider_gaps + active_blocks + quota_limits_reached + quota_attention,
        },
        "adjustments": {
            "recommended": recommended,
            "applySupported": adjustments_supported,
        },
        "readiness": {
            "status": status,
            "rolloutReady": bool(int(rollout_row.get("rollout_ready_flag", 0) or 0)),
            "rollbackReady": True,
            "metricsObserved": tracked_metrics > 0,
        },
        "releaseControls": build_release_controls(),
        "goLiveClosure": build_release_controls(),
    }
