from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "preferences": [],
    "notifications": [],
}


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def _find_tenant_id(cursor, tenant_slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (tenant_slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def _normalize_limit(limit: int | None, fallback: int = 50) -> int:
    if limit is None or limit <= 0:
        return fallback
    return min(limit, 100)


def _paginate(records: list[dict], cursor: str | None, limit: int) -> dict:
    page_limit = _normalize_limit(limit)
    start_index = 0
    if cursor:
        for index, record in enumerate(records):
            if record["publicId"] == cursor:
                start_index = index + 1
                break
    items = records[start_index : start_index + page_limit]
    next_cursor = None
    if start_index + page_limit < len(records) and items:
        next_cursor = items[-1]["publicId"]
    return {
        "items": items,
        "pageInfo": {
            "cursor": cursor,
            "limit": page_limit,
            "returned": len(items),
            "nextCursor": next_cursor,
            "hasMore": next_cursor is not None,
        },
    }


def upsert_preference(user_public_id: str, payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    normalized_user = user_public_id.strip()
    if normalized_user == "":
        raise ValueError("notification_user_public_id_required")
    quiet_hours = payload.get("quietHours") or {"from": "22:00", "to": "07:00"}
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "userPublicId": normalized_user,
        "inAppEnabled": bool(payload.get("inAppEnabled", True)),
        "emailEnabled": bool(payload.get("emailEnabled", False)),
        "quietHours": quiet_hours,
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        for index, item in enumerate(IN_MEMORY_STATE["preferences"]):
            if item["tenantSlug"] == slug and item["userPublicId"] == normalized_user:
                record["publicId"] = item["publicId"]
                IN_MEMORY_STATE["preferences"][index] = record
                return record
        IN_MEMORY_STATE["preferences"].append(record)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO notification.preferences (
                  tenant_id, public_id, user_public_id, in_app_enabled, email_enabled, quiet_hours_json
                )
                VALUES (%s, %s, %s, %s, %s, %s::jsonb)
                ON CONFLICT (tenant_id, user_public_id)
                DO UPDATE SET
                  in_app_enabled = EXCLUDED.in_app_enabled,
                  email_enabled = EXCLUDED.email_enabled,
                  quiet_hours_json = EXCLUDED.quiet_hours_json
                RETURNING public_id, updated_at
                """,
                (tenant_id, record["publicId"], normalized_user, record["inAppEnabled"], record["emailEnabled"], json.dumps(quiet_hours)),
            )
            row = cursor.fetchone()
            connection.commit()
            record["publicId"] = row["public_id"]
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def get_preference(user_public_id: str, tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["preferences"]:
            if item["tenantSlug"] == slug and item["userPublicId"] == user_public_id:
                return item
        return {
            "publicId": str(uuid.uuid4()),
            "tenantSlug": slug,
            "userPublicId": user_public_id,
            "inAppEnabled": True,
            "emailEnabled": False,
            "quietHours": {"from": "22:00", "to": "07:00"},
            "updatedAt": utc_now(),
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT public_id, in_app_enabled, email_enabled, quiet_hours_json, updated_at
                FROM notification.preferences AS preference
                JOIN identity.tenants AS tenant ON tenant.id = preference.tenant_id
                WHERE tenant.slug = %s AND preference.user_public_id = %s
                """,
                (slug, user_public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return {
                    "publicId": str(uuid.uuid4()),
                    "tenantSlug": slug,
                    "userPublicId": user_public_id,
                    "inAppEnabled": True,
                    "emailEnabled": False,
                    "quietHours": {"from": "22:00", "to": "07:00"},
                    "updatedAt": utc_now(),
                }
            return {
                "publicId": row["public_id"],
                "tenantSlug": slug,
                "userPublicId": user_public_id,
                "inAppEnabled": row["in_app_enabled"],
                "emailEnabled": row["email_enabled"],
                "quietHours": row["quiet_hours_json"] or {"from": "22:00", "to": "07:00"},
                "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            }


def create_notification(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    title = (payload.get("title") or "").strip()
    if title == "":
        raise ValueError("notification_title_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "userPublicId": (payload.get("userPublicId") or "").strip(),
        "title": title,
        "body": (payload.get("body") or "").strip(),
        "severity": (payload.get("severity") or "info").strip().lower(),
        "channel": (payload.get("channel") or "in_app").strip().lower(),
        "status": "unread",
        "sourceModule": (payload.get("sourceModule") or "manual").strip().lower(),
        "entityKind": (payload.get("entityKind") or "").strip().lower() or None,
        "entityPublicId": (payload.get("entityPublicId") or "").strip() or None,
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }
    if record["severity"] not in {"info", "warning", "critical", "success"}:
        raise ValueError("notification_severity_invalid")
    if record["channel"] not in {"in_app", "email"}:
        raise ValueError("notification_channel_invalid")

    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["notifications"].append(record)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO notification.notifications (
                  tenant_id, public_id, user_public_id, title, body, severity, channel, status, source_module,
                  entity_kind, entity_public_id
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, 'unread', %s, %s, %s)
                RETURNING created_at, updated_at
                """,
                (
                    tenant_id,
                    record["publicId"],
                    record["userPublicId"] or None,
                    title,
                    record["body"],
                    record["severity"],
                    record["channel"],
                    record["sourceModule"],
                    record["entityKind"],
                    record["entityPublicId"],
                ),
            )
            row = cursor.fetchone()
            connection.commit()
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def bulk_create_notifications(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    items = payload.get("items") or []
    if not isinstance(items, list):
        raise ValueError("notification_bulk_items_invalid")

    results: list[dict] = []
    succeeded = 0
    failed = 0
    for index, item in enumerate(items):
        candidate = {"tenantSlug": slug, **item}
        try:
            notification = create_notification(candidate)
            succeeded += 1
            results.append({"index": index, "status": "created", "notification": notification})
        except ValueError as error:
            failed += 1
            results.append(
                {
                    "index": index,
                    "status": "failed",
                    "errorCode": str(error),
                    "message": "Notification payload is invalid.",
                }
            )

    return {
        "tenantSlug": slug,
        "summary": {
            "requested": len(items),
            "succeeded": succeeded,
            "failed": failed,
            "partialSuccess": failed > 0,
        },
        "items": results,
    }


def list_notifications(tenant_slug: str | None = None, status: str | None = None, cursor: str | None = None, limit: int = 50) -> dict:
    slug = _tenant_slug(tenant_slug)
    normalized_status = (status or "").strip().lower()
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["notifications"] if item["tenantSlug"] == slug]
        if normalized_status:
            records = [item for item in records if item["status"] == normalized_status]
        records = sorted(records, key=lambda item: item["createdAt"], reverse=True)
        payload = _paginate(records, cursor, limit)
        payload["tenantSlug"] = slug
        return payload

    clauses = ["tenant.slug = %s"]
    params: list[object] = [slug]
    if normalized_status:
        clauses.append("notification.status = %s")
        params.append(normalized_status)

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                f"""
                SELECT notification.public_id, notification.user_public_id, notification.title, notification.body,
                       notification.severity, notification.channel, notification.status, notification.source_module,
                       notification.entity_kind, notification.entity_public_id, notification.created_at, notification.updated_at
                FROM notification.notifications AS notification
                JOIN identity.tenants AS tenant ON tenant.id = notification.tenant_id
                WHERE {" AND ".join(clauses)}
                ORDER BY notification.created_at DESC
                """,
                params,
            )
            items = [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "userPublicId": row["user_public_id"],
                    "title": row["title"],
                    "body": row["body"],
                    "severity": row["severity"],
                    "channel": row["channel"],
                    "status": row["status"],
                    "sourceModule": row["source_module"],
                    "entityKind": row["entity_kind"],
                    "entityPublicId": row["entity_public_id"],
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor_db.fetchall()
            ]
            payload = _paginate(items, cursor, limit)
            payload["tenantSlug"] = slug
            return payload


def transition_notification(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    status = (payload.get("status") or "").strip().lower()
    if status not in {"unread", "read", "archived"}:
        raise ValueError("notification_status_invalid")

    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["notifications"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                item["status"] = status
                item["updatedAt"] = utc_now()
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE notification.notifications AS notification
                SET status = %s
                FROM identity.tenants AS tenant
                WHERE tenant.id = notification.tenant_id
                  AND tenant.slug = %s
                  AND notification.public_id = %s
                RETURNING notification.public_id, notification.user_public_id, notification.title, notification.body,
                          notification.severity, notification.channel, notification.status, notification.source_module,
                          notification.entity_kind, notification.entity_public_id, notification.created_at, notification.updated_at
                """,
                (status, slug, public_id),
            )
            row = cursor.fetchone()
            connection.commit()
            if row is None:
                return None
            return {
                "publicId": row["public_id"],
                "tenantSlug": slug,
                "userPublicId": row["user_public_id"],
                "title": row["title"],
                "body": row["body"],
                "severity": row["severity"],
                "channel": row["channel"],
                "status": row["status"],
                "sourceModule": row["source_module"],
                "entityKind": row["entity_kind"],
                "entityPublicId": row["entity_public_id"],
                "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            }


def build_summary(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        items = [item for item in IN_MEMORY_STATE["notifications"] if item["tenantSlug"] == slug]
        return {
            "tenantSlug": slug,
            "summary": {
                "total": len(items),
                "unread": sum(1 for item in items if item["status"] == "unread"),
                "critical": sum(1 for item in items if item["severity"] == "critical"),
                "email": sum(1 for item in items if item["channel"] == "email"),
            },
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  count(*) AS total,
                  count(*) FILTER (WHERE notification.status = 'unread') AS unread_total,
                  count(*) FILTER (WHERE notification.severity = 'critical') AS critical_total,
                  count(*) FILTER (WHERE notification.channel = 'email') AS email_total
                FROM notification.notifications AS notification
                JOIN identity.tenants AS tenant ON tenant.id = notification.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            row = cursor.fetchone() or {}
            return {
                "tenantSlug": slug,
                "summary": {
                    "total": int(row.get("total", 0) or 0),
                    "unread": int(row.get("unread_total", 0) or 0),
                    "critical": int(row.get("critical_total", 0) or 0),
                    "email": int(row.get("email_total", 0) or 0),
                },
            }


def capability_catalog() -> dict:
    return {
        "service": settings.service_name,
        "repositoryDriver": settings.repository_driver,
        "capabilities": [
            {"key": "notification.center", "status": "ready"},
            {"key": "notification.preferences", "status": "ready"},
            {"key": "notification.internal-dispatch", "status": "ready"},
        ],
    }
