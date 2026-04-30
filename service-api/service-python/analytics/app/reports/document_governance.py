"""Relatorio de governanca documental e operacao de uploads."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_document_governance(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_document_governance(tenant_slug)

    return build_static_document_governance(tenant_slug)


def build_static_document_governance(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "inventory": {
            "attachmentsTotal": 24,
            "activeAttachments": 19,
            "archivedAttachments": 5,
            "retentionLongTerm": 11,
        },
        "visibility": {
            "internal": 8,
            "restricted": 10,
            "public": 6,
        },
        "storage": {
            "manual": 7,
            "external": 17,
            "drivers": {"local": 5, "manual": 7, "r2": 9, "s3": 3},
        },
        "uploads": {
            "sessionsTotal": 9,
            "pending": 2,
            "completed": 6,
            "expired": 1,
        },
        "ownership": {
            "crm.lead": 7,
            "crm.customer": 9,
            "rentals.contract": 5,
            "sales.sale": 3,
        },
    }


def build_postgres_document_governance(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        inventory = fetch_inventory(connection, tenant_slug)
        visibility = fetch_visibility(connection, tenant_slug)
        storage = fetch_storage(connection, tenant_slug)
        uploads = fetch_uploads(connection, tenant_slug)
        ownership = fetch_ownership(connection, tenant_slug)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "inventory": inventory,
        "visibility": visibility,
        "storage": storage,
        "uploads": uploads,
        "ownership": ownership,
    }


def fetch_inventory(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM documents.attachments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS attachments_total,
            (
                SELECT count(*)
                FROM documents.attachments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND archived_at IS NULL
            ) AS active_attachments,
            (
                SELECT count(*)
                FROM documents.attachments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND archived_at IS NOT NULL
            ) AS archived_attachments,
            (
                SELECT count(*)
                FROM documents.attachments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND retention_days >= 180
            ) AS retention_long_term
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 4)
        row = cursor.fetchone() or {}

    return {
        "attachmentsTotal": int(row.get("attachments_total", 0) or 0),
        "activeAttachments": int(row.get("active_attachments", 0) or 0),
        "archivedAttachments": int(row.get("archived_attachments", 0) or 0),
        "retentionLongTerm": int(row.get("retention_long_term", 0) or 0),
    }


def fetch_visibility(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) FILTER (WHERE visibility = 'internal') AS internal_count,
            count(*) FILTER (WHERE visibility = 'restricted') AS restricted_count,
            count(*) FILTER (WHERE visibility = 'public') AS public_count
        FROM documents.attachments
        WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "internal": int(row.get("internal_count", 0) or 0),
        "restricted": int(row.get("restricted_count", 0) or 0),
        "public": int(row.get("public_count", 0) or 0),
    }


def fetch_storage(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT storage_driver, count(*) AS total
        FROM documents.attachments
        WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
        GROUP BY storage_driver
        ORDER BY storage_driver
    """

    drivers: dict[str, int] = {}
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall() or []

    for row in rows:
        drivers[str(row.get("storage_driver", "") or "")] = int(row.get("total", 0) or 0)

    manual_total = int(drivers.get("manual", 0) or 0)
    external_total = sum(value for key, value in drivers.items() if key != "manual")

    return {
        "manual": manual_total,
        "external": external_total,
        "drivers": drivers,
    }


def fetch_uploads(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS sessions_total,
            count(*) FILTER (WHERE status = 'pending_upload') AS pending_count,
            count(*) FILTER (WHERE status = 'completed') AS completed_count,
            count(*) FILTER (WHERE status = 'expired') AS expired_count
        FROM documents.upload_sessions
        WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "sessionsTotal": int(row.get("sessions_total", 0) or 0),
        "pending": int(row.get("pending_count", 0) or 0),
        "completed": int(row.get("completed_count", 0) or 0),
        "expired": int(row.get("expired_count", 0) or 0),
    }


def fetch_ownership(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT owner_type, count(*) AS total
        FROM documents.attachments
        WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
        GROUP BY owner_type
        ORDER BY owner_type
    """

    ownership: dict[str, int] = {}
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall() or []

    for row in rows:
        ownership[str(row.get("owner_type", "") or "")] = int(row.get("total", 0) or 0)

    return ownership
