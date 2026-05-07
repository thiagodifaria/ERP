"""Painel consolidado dos novos contextos de core expansion."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_core_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": slug,
            "generatedAt": datetime.now(timezone.utc).isoformat(),
            "summary": {"catalogItems": 12, "suppliers": 6, "supportCases": 8, "notifications": 15},
            "catalog": {"categories": 4, "activeItems": 10},
            "supplier": {"active": 5, "watchlist": 1},
            "support": {"open": 3, "overdue": 1},
            "notification": {"unread": 5, "critical": 2},
            "readiness": {"status": "stable", "catalogReady": True, "supplierReady": True, "supportReady": True, "notificationReady": True},
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  (SELECT count(*) FROM catalog.items AS item JOIN identity.tenants AS tenant_item ON tenant_item.id = item.tenant_id WHERE tenant_item.slug = %s) AS catalog_items,
                  (SELECT count(*) FROM catalog.items AS item JOIN identity.tenants AS tenant_item ON tenant_item.id = item.tenant_id WHERE tenant_item.slug = %s AND item.active = TRUE) AS active_items,
                  (SELECT count(*) FROM catalog.categories AS category JOIN identity.tenants AS tenant_category ON tenant_category.id = category.tenant_id WHERE tenant_category.slug = %s) AS categories_total,
                  (SELECT count(*) FROM supplier.suppliers AS supplier JOIN identity.tenants AS tenant_supplier ON tenant_supplier.id = supplier.tenant_id WHERE tenant_supplier.slug = %s) AS suppliers_total,
                  (SELECT count(*) FROM supplier.suppliers AS supplier JOIN identity.tenants AS tenant_supplier ON tenant_supplier.id = supplier.tenant_id WHERE tenant_supplier.slug = %s AND supplier.status = 'active') AS suppliers_active,
                  (SELECT count(*) FROM supplier.suppliers AS supplier JOIN identity.tenants AS tenant_supplier ON tenant_supplier.id = supplier.tenant_id WHERE tenant_supplier.slug = %s AND supplier.status = 'watchlist') AS suppliers_watchlist,
                  (SELECT count(*) FROM support.cases AS support_case JOIN identity.tenants AS tenant_support ON tenant_support.id = support_case.tenant_id WHERE tenant_support.slug = %s) AS support_cases_total,
                  (SELECT count(*) FROM support.cases AS support_case JOIN identity.tenants AS tenant_support ON tenant_support.id = support_case.tenant_id WHERE tenant_support.slug = %s AND support_case.status IN ('open', 'in_progress', 'waiting_customer')) AS support_cases_open,
                  (SELECT count(*) FROM support.cases AS support_case JOIN identity.tenants AS tenant_support ON tenant_support.id = support_case.tenant_id WHERE tenant_support.slug = %s AND support_case.status NOT IN ('resolved', 'closed') AND support_case.sla_due_at < NOW()) AS support_cases_overdue,
                  (SELECT count(*) FROM notification.notifications AS notification JOIN identity.tenants AS tenant_notification ON tenant_notification.id = notification.tenant_id WHERE tenant_notification.slug = %s) AS notifications_total,
                  (SELECT count(*) FROM notification.notifications AS notification JOIN identity.tenants AS tenant_notification ON tenant_notification.id = notification.tenant_id WHERE tenant_notification.slug = %s AND notification.status = 'unread') AS notifications_unread,
                  (SELECT count(*) FROM notification.notifications AS notification JOIN identity.tenants AS tenant_notification ON tenant_notification.id = notification.tenant_id WHERE tenant_notification.slug = %s AND notification.severity = 'critical') AS notifications_critical
                """,
                (slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug),
            )
            row = cursor.fetchone() or {}

    readiness_status = "stable"
    if int(row.get("support_cases_overdue", 0) or 0) > 0 or int(row.get("notifications_critical", 0) or 0) > 0:
        readiness_status = "attention"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "summary": {
            "catalogItems": int(row.get("catalog_items", 0) or 0),
            "suppliers": int(row.get("suppliers_total", 0) or 0),
            "supportCases": int(row.get("support_cases_total", 0) or 0),
            "notifications": int(row.get("notifications_total", 0) or 0),
        },
        "catalog": {
            "categories": int(row.get("categories_total", 0) or 0),
            "activeItems": int(row.get("active_items", 0) or 0),
        },
        "supplier": {
            "active": int(row.get("suppliers_active", 0) or 0),
            "watchlist": int(row.get("suppliers_watchlist", 0) or 0),
        },
        "support": {
            "open": int(row.get("support_cases_open", 0) or 0),
            "overdue": int(row.get("support_cases_overdue", 0) or 0),
        },
        "notification": {
            "unread": int(row.get("notifications_unread", 0) or 0),
            "critical": int(row.get("notifications_critical", 0) or 0),
        },
        "readiness": {
            "status": readiness_status,
            "catalogReady": int(row.get("catalog_items", 0) or 0) > 0,
            "supplierReady": int(row.get("suppliers_total", 0) or 0) > 0,
            "supportReady": int(row.get("support_cases_total", 0) or 0) > 0,
            "notificationReady": int(row.get("notifications_total", 0) or 0) > 0,
        },
    }
