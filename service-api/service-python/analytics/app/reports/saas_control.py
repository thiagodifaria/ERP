"""Painel executivo de consumo, limites e lifecycle de tenants."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_saas_control(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver != "postgres":
        slug = tenant_slug or "bootstrap-ops"
        return {
            "tenantSlug": slug,
            "generatedAt": datetime.now(timezone.utc).isoformat(),
            "entitlements": {"total": 6, "enabled": 5},
            "quotas": {"active": 4, "attention": 1, "hard": 1},
            "metering": {"trackedMetrics": 4, "totalQuantity": 4096},
            "blocks": {"active": 0},
            "lifecycle": {"queued": 1, "running": 0, "completed": 2, "failed": 0},
            "readiness": {"status": "stable", "onboardingReady": True, "offboardingReady": True, "usageBasedBillingReady": True},
        }

    slug = tenant_slug or "bootstrap-ops"
    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  count(*) FILTER (WHERE entitlement.enabled) AS entitlements_enabled,
                  count(*) AS entitlements_total
                FROM platform_control.entitlements AS entitlement
                JOIN identity.tenants AS tenant ON tenant.id = entitlement.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            entitlements = cursor.fetchone() or {}

            cursor.execute(
                """
                SELECT
                  count(*) AS quotas_total,
                  count(*) FILTER (WHERE enforcement_mode = 'hard') AS quotas_hard
                FROM platform_control.quotas AS quota
                JOIN identity.tenants AS tenant ON tenant.id = quota.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            quotas = cursor.fetchone() or {}

            cursor.execute(
                """
                SELECT count(*) AS blocks_active
                FROM platform_control.tenant_blocks AS block
                JOIN identity.tenants AS tenant ON tenant.id = block.tenant_id
                WHERE tenant.slug = %s AND block.active = TRUE
                """,
                (slug,),
            )
            blocks = cursor.fetchone() or {}

            cursor.execute(
                """
                SELECT
                  count(*) AS metrics_total,
                  coalesce(sum(quantity), 0) AS quantity_total
                FROM platform_control.usage_snapshots AS snapshot
                JOIN identity.tenants AS tenant ON tenant.id = snapshot.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            metering = cursor.fetchone() or {}

            cursor.execute(
                """
                SELECT
                  count(*) FILTER (WHERE status = 'queued') AS queued,
                  count(*) FILTER (WHERE status = 'running') AS running,
                  count(*) FILTER (WHERE status = 'completed') AS completed,
                  count(*) FILTER (WHERE status = 'failed') AS failed
                FROM platform_control.lifecycle_jobs AS job
                JOIN identity.tenants AS tenant ON tenant.id = job.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            lifecycle = cursor.fetchone() or {}

    attention = 1 if int(blocks.get("blocks_active", 0) or 0) > 0 else 0
    status = "stable" if attention == 0 else "attention"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "entitlements": {
            "total": int(entitlements.get("entitlements_total", 0) or 0),
            "enabled": int(entitlements.get("entitlements_enabled", 0) or 0),
        },
        "quotas": {
            "active": int(quotas.get("quotas_total", 0) or 0),
            "attention": attention,
            "hard": int(quotas.get("quotas_hard", 0) or 0),
        },
        "metering": {
            "trackedMetrics": int(metering.get("metrics_total", 0) or 0),
            "totalQuantity": int(metering.get("quantity_total", 0) or 0),
        },
        "blocks": {
            "active": int(blocks.get("blocks_active", 0) or 0),
        },
        "lifecycle": {
            "queued": int(lifecycle.get("queued", 0) or 0),
            "running": int(lifecycle.get("running", 0) or 0),
            "completed": int(lifecycle.get("completed", 0) or 0),
            "failed": int(lifecycle.get("failed", 0) or 0),
        },
        "readiness": {
            "status": status,
            "onboardingReady": True,
            "offboardingReady": True,
            "usageBasedBillingReady": int(quotas.get("quotas_total", 0) or 0) > 0,
        },
    }
