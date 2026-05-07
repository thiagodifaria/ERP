from __future__ import annotations

from datetime import datetime, timedelta, timezone
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "queues": [],
    "cases": [],
    "events": [],
}

QUEUE_TEMPLATES = [
    {"queueKey": "general", "name": "General", "slaTargetHours": 24, "active": True},
    {"queueKey": "billing", "name": "Billing", "slaTargetHours": 8, "active": True},
    {"queueKey": "technical", "name": "Technical", "slaTargetHours": 4, "active": True},
]


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _parse_timestamp(value: str) -> datetime:
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def _tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def _normalize_limit(limit: int | None, fallback: int = 50) -> int:
    if limit is None or limit <= 0:
        return fallback
    return min(limit, 100)


def _paginate(records: list[dict], cursor: str | None, limit: int) -> dict:
    page_limit = _normalize_limit(limit)
    start_index = 0
    if cursor:
        for index, record in enumerate(records):
            if record.get("publicId") == cursor:
                start_index = index + 1
                break
    items = records[start_index : start_index + page_limit]
    next_cursor = None
    if start_index + page_limit < len(records) and items:
        next_cursor = str(items[-1]["publicId"])
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


def _find_tenant_id(cursor, tenant_slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (tenant_slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def _ensure_bootstrap_queues_in_memory(tenant_slug: str) -> None:
    if any(queue["tenantSlug"] == tenant_slug for queue in IN_MEMORY_STATE["queues"]):
        return
    for template in QUEUE_TEMPLATES:
        IN_MEMORY_STATE["queues"].append(
            {
                "publicId": str(uuid.uuid4()),
                "tenantSlug": tenant_slug,
                **template,
                "createdAt": utc_now(),
                "updatedAt": utc_now(),
            }
        )


def _default_sla_hours(priority: str) -> int:
    return {
        "low": 48,
        "medium": 24,
        "high": 8,
        "urgent": 4,
    }.get(priority, 24)


def _append_event(case_public_id: str, tenant_slug: str, event_type: str, summary: str) -> dict:
    event = {
        "publicId": str(uuid.uuid4()),
        "casePublicId": case_public_id,
        "tenantSlug": tenant_slug,
        "eventType": event_type,
        "summary": summary,
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["events"].append(event)
    return event


def list_queues(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        _ensure_bootstrap_queues_in_memory(slug)
        return sorted([queue for queue in IN_MEMORY_STATE["queues"] if queue["tenantSlug"] == slug], key=lambda item: item["name"])

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                INSERT INTO support.queues (tenant_id, public_id, queue_key, name, sla_target_hours, active)
                SELECT tenant.id, gen_random_uuid(), template.queue_key, template.name, template.sla_target_hours, TRUE
                FROM identity.tenants AS tenant,
                     (VALUES ('general', 'General', 24), ('billing', 'Billing', 8), ('technical', 'Technical', 4))
                     AS template(queue_key, name, sla_target_hours)
                WHERE tenant.slug = %s
                  AND NOT EXISTS (
                    SELECT 1
                    FROM support.queues AS queue
                    WHERE queue.tenant_id = tenant.id AND queue.queue_key = template.queue_key
                  )
                """,
                (slug,),
            )
            cursor.execute(
                """
                SELECT queue.public_id, queue.queue_key, queue.name, queue.sla_target_hours, queue.active, queue.created_at, queue.updated_at
                FROM support.queues AS queue
                JOIN identity.tenants AS tenant ON tenant.id = queue.tenant_id
                WHERE tenant.slug = %s
                ORDER BY queue.name
                """,
                (slug,),
            )
            rows = cursor.fetchall()
            connection.commit()
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "queueKey": row["queue_key"],
                    "name": row["name"],
                    "slaTargetHours": int(row["sla_target_hours"]),
                    "active": row["active"],
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in rows
            ]


def upsert_queue(queue_key: str, payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    normalized_key = queue_key.strip().lower()
    name = (payload.get("name") or "").strip()
    if normalized_key == "":
        raise ValueError("support_queue_key_required")
    if name == "":
        raise ValueError("support_queue_name_required")

    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "queueKey": normalized_key,
        "name": name,
        "slaTargetHours": int(payload.get("slaTargetHours") or 24),
        "active": bool(payload.get("active", True)),
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        _ensure_bootstrap_queues_in_memory(slug)
        for index, queue in enumerate(IN_MEMORY_STATE["queues"]):
            if queue["tenantSlug"] == slug and queue["queueKey"] == normalized_key:
                record["publicId"] = queue["publicId"]
                record["createdAt"] = queue["createdAt"]
                IN_MEMORY_STATE["queues"][index] = record
                return record
        IN_MEMORY_STATE["queues"].append(record)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO support.queues (tenant_id, public_id, queue_key, name, sla_target_hours, active)
                VALUES (%s, %s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, queue_key)
                DO UPDATE SET
                  name = EXCLUDED.name,
                  sla_target_hours = EXCLUDED.sla_target_hours,
                  active = EXCLUDED.active
                RETURNING public_id, created_at, updated_at
                """,
                (tenant_id, record["publicId"], normalized_key, name, record["slaTargetHours"], record["active"]),
            )
            row = cursor.fetchone()
            connection.commit()
            record["publicId"] = row["public_id"]
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def _map_case_row(row: dict, tenant_slug: str, events: list[dict]) -> dict:
    return {
        "publicId": row["public_id"],
        "tenantSlug": tenant_slug,
        "caseKey": row["case_key"],
        "subject": row["subject"],
        "status": row["status"],
        "priority": row["priority"],
        "queueKey": row["queue_key"],
        "ownerUserId": row["owner_user_id"],
        "sourceKind": row["source_kind"],
        "entityKind": row["entity_kind"],
        "entityPublicId": row["entity_public_id"],
        "slaDueAt": row["sla_due_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ") if row["sla_due_at"] else None,
        "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "events": events,
    }


def list_cases(tenant_slug: str | None = None, status: str | None = None, priority: str | None = None, cursor: str | None = None, limit: int = 50) -> dict:
    slug = _tenant_slug(tenant_slug)
    normalized_status = (status or "").strip().lower()
    normalized_priority = (priority or "").strip().lower()

    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["cases"] if item["tenantSlug"] == slug]
        if normalized_status:
            records = [item for item in records if item["status"] == normalized_status]
        if normalized_priority:
            records = [item for item in records if item["priority"] == normalized_priority]
        records = sorted(records, key=lambda item: item["createdAt"], reverse=True)
        payload = _paginate(records, cursor, limit)
        payload["tenantSlug"] = slug
        return payload

    clauses = ["tenant.slug = %s"]
    params: list[object] = [slug]
    if normalized_status:
        clauses.append("support_case.status = %s")
        params.append(normalized_status)
    if normalized_priority:
        clauses.append("support_case.priority = %s")
        params.append(normalized_priority)

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                f"""
                SELECT support_case.public_id, support_case.case_key, support_case.subject, support_case.status, support_case.priority,
                       queue.queue_key, support_case.owner_user_id, support_case.source_kind, support_case.entity_kind,
                       support_case.entity_public_id, support_case.sla_due_at, support_case.created_at, support_case.updated_at
                FROM support.cases AS support_case
                JOIN identity.tenants AS tenant ON tenant.id = support_case.tenant_id
                JOIN support.queues AS queue ON queue.id = support_case.queue_id
                WHERE {" AND ".join(clauses)}
                ORDER BY support_case.created_at DESC
                """,
                params,
            )
            rows = cursor_db.fetchall()
            records = [get_case(row["public_id"], slug) for row in rows]
            payload = _paginate([record for record in records if record is not None], cursor, limit)
            payload["tenantSlug"] = slug
            return payload


def create_case(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    subject = (payload.get("subject") or "").strip()
    queue_key = (payload.get("queueKey") or "").strip().lower()
    priority = (payload.get("priority") or "medium").strip().lower()
    if subject == "":
        raise ValueError("support_case_subject_required")
    if queue_key == "":
        raise ValueError("support_queue_key_required")
    if priority not in {"low", "medium", "high", "urgent"}:
        raise ValueError("support_case_priority_invalid")

    case_key = f"SUP-{uuid.uuid4().hex[:8].upper()}"
    sla_hours = int(payload.get("slaTargetHours") or _default_sla_hours(priority))
    sla_due_at = (datetime.now(timezone.utc) + timedelta(hours=sla_hours)).strftime("%Y-%m-%dT%H:%M:%SZ")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "caseKey": case_key,
        "subject": subject,
        "status": "open",
        "priority": priority,
        "queueKey": queue_key,
        "ownerUserId": (payload.get("ownerUserId") or "").strip() or None,
        "sourceKind": (payload.get("sourceKind") or "manual").strip().lower(),
        "entityKind": (payload.get("entityKind") or "").strip().lower() or None,
        "entityPublicId": (payload.get("entityPublicId") or "").strip() or None,
        "slaDueAt": sla_due_at,
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
        "events": [],
    }

    if settings.repository_driver != "postgres":
        _ensure_bootstrap_queues_in_memory(slug)
        if not any(queue["tenantSlug"] == slug and queue["queueKey"] == queue_key for queue in IN_MEMORY_STATE["queues"]):
            raise ValueError("support_queue_not_found")
        IN_MEMORY_STATE["cases"].append(record)
        record["events"] = [_append_event(record["publicId"], slug, "created", "Case created.")]
        return record

    with connect() as connection:
        with connection.cursor() as cursor_db:
            tenant_id = _find_tenant_id(cursor_db, slug)
            cursor_db.execute(
                "SELECT id FROM support.queues WHERE tenant_id = %s AND queue_key = %s",
                (tenant_id, queue_key),
            )
            queue_row = cursor_db.fetchone()
            if queue_row is None:
                raise ValueError("support_queue_not_found")
            queue_id = int(queue_row["id"])
            cursor_db.execute(
                """
                INSERT INTO support.cases (
                  tenant_id, queue_id, public_id, case_key, subject, status, priority, owner_user_id,
                  source_kind, entity_kind, entity_public_id, sla_due_at
                )
                VALUES (%s, %s, %s, %s, %s, 'open', %s, %s, %s, %s, %s, %s)
                RETURNING created_at, updated_at
                """,
                (
                    tenant_id,
                    queue_id,
                    record["publicId"],
                    case_key,
                    subject,
                    priority,
                    record["ownerUserId"],
                    record["sourceKind"],
                    record["entityKind"],
                    record["entityPublicId"],
                    _parse_timestamp(sla_due_at),
                ),
            )
            case_row = cursor_db.fetchone()
            cursor_db.execute(
                """
                INSERT INTO support.case_events (case_id, public_id, event_type, summary)
                SELECT id, %s, 'created', 'Case created.'
                FROM support.cases
                WHERE public_id = %s
                """,
                (str(uuid.uuid4()), record["publicId"]),
            )
            connection.commit()
            record["createdAt"] = case_row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = case_row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return get_case(record["publicId"], slug) or record


def bulk_create_cases(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    items = payload.get("items") or []
    if not isinstance(items, list):
        raise ValueError("support_bulk_items_invalid")

    results: list[dict] = []
    succeeded = 0
    failed = 0
    for index, item in enumerate(items):
        candidate = {"tenantSlug": slug, **item}
        try:
            created = create_case(candidate)
            succeeded += 1
            results.append({"index": index, "status": "created", "case": created})
        except ValueError as error:
            failed += 1
            results.append(
                {
                    "index": index,
                    "status": "failed",
                    "errorCode": str(error),
                    "message": "Support case payload is invalid.",
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


def get_case(public_id: str, tenant_slug: str | None = None) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        for record in IN_MEMORY_STATE["cases"]:
            if record["tenantSlug"] == slug and record["publicId"] == public_id:
                record["events"] = [
                    event for event in IN_MEMORY_STATE["events"] if event["tenantSlug"] == slug and event["casePublicId"] == public_id
                ]
                return record
        return None

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                SELECT support_case.public_id, support_case.case_key, support_case.subject, support_case.status, support_case.priority,
                       queue.queue_key, support_case.owner_user_id, support_case.source_kind, support_case.entity_kind,
                       support_case.entity_public_id, support_case.sla_due_at, support_case.created_at, support_case.updated_at
                FROM support.cases AS support_case
                JOIN identity.tenants AS tenant ON tenant.id = support_case.tenant_id
                JOIN support.queues AS queue ON queue.id = support_case.queue_id
                WHERE tenant.slug = %s AND support_case.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor_db.fetchone()
            if row is None:
                return None
            cursor_db.execute(
                """
                SELECT event.public_id, event.event_type, event.summary, event.created_at
                FROM support.case_events AS event
                JOIN support.cases AS support_case ON support_case.id = event.case_id
                WHERE support_case.public_id = %s
                ORDER BY event.created_at
                """,
                (public_id,),
            )
            events = [
                {
                    "publicId": event["public_id"],
                    "eventType": event["event_type"],
                    "summary": event["summary"],
                    "createdAt": event["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for event in cursor_db.fetchall()
            ]
            return _map_case_row(row, slug, events)


def transition_case(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    status = (payload.get("status") or "").strip().lower()
    if status not in {"open", "in_progress", "waiting_customer", "resolved", "closed"}:
        raise ValueError("support_case_status_invalid")
    summary = (payload.get("summary") or f"Case moved to {status}.").strip()

    if settings.repository_driver != "postgres":
        for record in IN_MEMORY_STATE["cases"]:
            if record["tenantSlug"] == slug and record["publicId"] == public_id:
                record["status"] = status
                record["updatedAt"] = utc_now()
                _append_event(public_id, slug, "status_changed", summary)
                return get_case(public_id, slug)
        return None

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                UPDATE support.cases AS support_case
                SET status = %s
                FROM identity.tenants AS tenant
                WHERE tenant.id = support_case.tenant_id
                  AND tenant.slug = %s
                  AND support_case.public_id = %s
                RETURNING support_case.id
                """,
                (status, slug, public_id),
            )
            row = cursor_db.fetchone()
            if row is None:
                return None
            cursor_db.execute(
                "INSERT INTO support.case_events (case_id, public_id, event_type, summary) VALUES (%s, %s, 'status_changed', %s)",
                (row["id"], str(uuid.uuid4()), summary),
            )
            connection.commit()
            return get_case(public_id, slug)


def add_case_comment(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    message = (payload.get("message") or "").strip()
    if message == "":
        raise ValueError("support_case_comment_required")

    if settings.repository_driver != "postgres":
        for record in IN_MEMORY_STATE["cases"]:
            if record["tenantSlug"] == slug and record["publicId"] == public_id:
                _append_event(public_id, slug, "comment", message)
                return get_case(public_id, slug)
        return None

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                SELECT support_case.id
                FROM support.cases AS support_case
                JOIN identity.tenants AS tenant ON tenant.id = support_case.tenant_id
                WHERE tenant.slug = %s AND support_case.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor_db.fetchone()
            if row is None:
                return None
            cursor_db.execute(
                "INSERT INTO support.case_events (case_id, public_id, event_type, summary) VALUES (%s, %s, 'comment', %s)",
                (row["id"], str(uuid.uuid4()), message),
            )
            connection.commit()
            return get_case(public_id, slug)


def export_cases(tenant_slug: str | None = None, status: str | None = None, priority: str | None = None) -> dict:
    listing = list_cases(tenant_slug, status, priority, None, 1000)
    return {
        "tenantSlug": listing["tenantSlug"],
        "exported": len(listing["items"]),
        "items": listing["items"],
        "pageInfo": listing["pageInfo"],
    }


def build_summary(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        cases = [item for item in IN_MEMORY_STATE["cases"] if item["tenantSlug"] == slug]
        overdue = sum(
            1
            for item in cases
            if item["status"] not in {"resolved", "closed"} and _parse_timestamp(item["slaDueAt"]) < datetime.now(timezone.utc)
        )
        return {
            "tenantSlug": slug,
            "summary": {
                "total": len(cases),
                "open": sum(1 for item in cases if item["status"] == "open"),
                "inProgress": sum(1 for item in cases if item["status"] == "in_progress"),
                "resolved": sum(1 for item in cases if item["status"] == "resolved"),
                "overdue": overdue,
            },
            "byPriority": {
                "low": sum(1 for item in cases if item["priority"] == "low"),
                "medium": sum(1 for item in cases if item["priority"] == "medium"),
                "high": sum(1 for item in cases if item["priority"] == "high"),
                "urgent": sum(1 for item in cases if item["priority"] == "urgent"),
            },
        }

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                SELECT
                  count(*) AS total,
                  count(*) FILTER (WHERE support_case.status = 'open') AS open_total,
                  count(*) FILTER (WHERE support_case.status = 'in_progress') AS in_progress_total,
                  count(*) FILTER (WHERE support_case.status = 'resolved') AS resolved_total,
                  count(*) FILTER (
                    WHERE support_case.status NOT IN ('resolved', 'closed') AND support_case.sla_due_at < NOW()
                  ) AS overdue_total,
                  count(*) FILTER (WHERE support_case.priority = 'low') AS low_total,
                  count(*) FILTER (WHERE support_case.priority = 'medium') AS medium_total,
                  count(*) FILTER (WHERE support_case.priority = 'high') AS high_total,
                  count(*) FILTER (WHERE support_case.priority = 'urgent') AS urgent_total
                FROM support.cases AS support_case
                JOIN identity.tenants AS tenant ON tenant.id = support_case.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            row = cursor_db.fetchone() or {}
            return {
                "tenantSlug": slug,
                "summary": {
                    "total": int(row.get("total", 0) or 0),
                    "open": int(row.get("open_total", 0) or 0),
                    "inProgress": int(row.get("in_progress_total", 0) or 0),
                    "resolved": int(row.get("resolved_total", 0) or 0),
                    "overdue": int(row.get("overdue_total", 0) or 0),
                },
                "byPriority": {
                    "low": int(row.get("low_total", 0) or 0),
                    "medium": int(row.get("medium_total", 0) or 0),
                    "high": int(row.get("high_total", 0) or 0),
                    "urgent": int(row.get("urgent_total", 0) or 0),
                },
            }


def capability_catalog() -> dict:
    return {
        "service": settings.service_name,
        "repositoryDriver": settings.repository_driver,
        "capabilities": [
            {"key": "support.cases", "status": "ready"},
            {"key": "support.sla", "status": "ready"},
            {"key": "support.case-history", "status": "ready"},
        ],
    }
