"""Painel de rollout, adocao e readiness de go-live."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_go_live_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": slug,
            "generatedAt": datetime.now(timezone.utc).isoformat(),
            "rollouts": {"planned": 1, "running": 0, "completed": 1, "rolledBack": 0},
            "adoption": {"trackedMetrics": 4, "totalQuantity": 4096},
            "readiness": {"status": "stable", "rolloutReady": True, "rollbackReady": True, "metricsObserved": True},
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
                  coalesce(max((rollout.readiness_json ->> 'rolloutReady')::boolean::int), 0) AS rollout_ready_flag
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

    status = "stable"
    if int(rollout_row.get("running_total", 0) or 0) > 0 or int(usage_row.get("tracked_metrics", 0) or 0) == 0:
        status = "attention"

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
            "trackedMetrics": int(usage_row.get("tracked_metrics", 0) or 0),
            "totalQuantity": int(usage_row.get("total_quantity", 0) or 0),
        },
        "readiness": {
            "status": status,
            "rolloutReady": bool(int(rollout_row.get("rollout_ready_flag", 0) or 0)),
            "rollbackReady": True,
            "metricsObserved": int(usage_row.get("tracked_metrics", 0) or 0) > 0,
        },
    }
