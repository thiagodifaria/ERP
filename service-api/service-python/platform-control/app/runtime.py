from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "entitlements": [],
    "metering": [],
    "jobs": [],
}

CAPABILITY_CATALOG = [
    {"capabilityKey": "catalog.items", "module": "catalog", "defaultEnabled": True},
    {"capabilityKey": "support.cases", "module": "support", "defaultEnabled": False},
    {"capabilityKey": "notifications.center", "module": "notification", "defaultEnabled": True},
    {"capabilityKey": "engagement.providers.meta_ads", "module": "engagement", "defaultEnabled": False},
    {"capabilityKey": "billing.pix", "module": "billing", "defaultEnabled": False},
    {"capabilityKey": "documents.external_storage", "module": "documents", "defaultEnabled": False},
]


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


def list_capability_catalog() -> list[dict]:
    return CAPABILITY_CATALOG


def list_entitlements(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [item for item in IN_MEMORY_STATE["entitlements"] if item["tenantSlug"] == slug]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT entitlement.public_id, entitlement.capability_key, entitlement.enabled, entitlement.plan_code,
                       entitlement.limit_value, entitlement.source, entitlement.updated_at
                FROM platform_control.entitlements AS entitlement
                JOIN identity.tenants AS tenant ON tenant.id = entitlement.tenant_id
                WHERE tenant.slug = %s
                ORDER BY entitlement.capability_key
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "capabilityKey": row["capability_key"],
                    "enabled": row["enabled"],
                    "planCode": row["plan_code"],
                    "limitValue": row["limit_value"],
                    "source": row["source"],
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def upsert_entitlement(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    normalized_key = capability_key.strip()
    if normalized_key == "":
        raise ValueError("capability_key_required")

    entitlement = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "capabilityKey": normalized_key,
        "enabled": bool(payload.get("enabled", True)),
        "planCode": (payload.get("planCode") or "custom").strip() or "custom",
        "limitValue": int(payload.get("limitValue", 0) or 0),
        "source": (payload.get("source") or "manual").strip() or "manual",
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        existing = next((item for item in IN_MEMORY_STATE["entitlements"] if item["tenantSlug"] == slug and item["capabilityKey"] == normalized_key), None)
        if existing is not None:
            existing.update(entitlement)
            entitlement["publicId"] = existing["publicId"]
            return existing
        IN_MEMORY_STATE["entitlements"].append(entitlement)
        return entitlement

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.entitlements (
                  tenant_id, public_id, capability_key, enabled, plan_code, limit_value, source
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, capability_key)
                DO UPDATE SET
                  enabled = EXCLUDED.enabled,
                  plan_code = EXCLUDED.plan_code,
                  limit_value = EXCLUDED.limit_value,
                  source = EXCLUDED.source,
                  updated_at = NOW()
                RETURNING public_id, updated_at
                """,
                (tenant_id, entitlement["publicId"], normalized_key, entitlement["enabled"], entitlement["planCode"], entitlement["limitValue"], entitlement["source"]),
            )
            row = cursor.fetchone()
            connection.commit()
            entitlement["publicId"] = row["public_id"]
            entitlement["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return entitlement


def list_metering(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        snapshots = [item for item in IN_MEMORY_STATE["metering"] if item["tenantSlug"] == slug]
        return {"tenantSlug": slug, "snapshots": snapshots, "summary": {"metrics": len(snapshots), "totalQuantity": sum(item["quantity"] for item in snapshots)}}

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT metric_key, metric_unit, quantity, source, captured_at
                FROM platform_control.usage_snapshots AS snapshot
                JOIN identity.tenants AS tenant ON tenant.id = snapshot.tenant_id
                WHERE tenant.slug = %s
                ORDER BY captured_at DESC
                LIMIT 50
                """,
                (slug,),
            )
            rows = cursor.fetchall()
            snapshots = [
                {
                    "metricKey": row["metric_key"],
                    "metricUnit": row["metric_unit"],
                    "quantity": row["quantity"],
                    "source": row["source"],
                    "capturedAt": row["captured_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in rows
            ]
            return {"tenantSlug": slug, "snapshots": snapshots, "summary": {"metrics": len(snapshots), "totalQuantity": sum(item["quantity"] for item in snapshots)}}


def create_metering_snapshot(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    metric_key = (payload.get("metricKey") or "").strip()
    metric_unit = (payload.get("metricUnit") or "").strip()
    source = (payload.get("source") or "manual").strip() or "manual"
    quantity = int(payload.get("quantity", 0) or 0)
    if metric_key == "":
        raise ValueError("metric_key_required")
    if metric_unit == "":
        raise ValueError("metric_unit_required")

    snapshot = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "metricKey": metric_key,
        "metricUnit": metric_unit,
        "quantity": quantity,
        "source": source,
        "capturedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["metering"].append(snapshot)
        return snapshot

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.usage_snapshots (
                  tenant_id, public_id, metric_key, metric_unit, quantity, source
                )
                VALUES (%s, %s, %s, %s, %s, %s)
                RETURNING captured_at
                """,
                (tenant_id, snapshot["publicId"], metric_key, metric_unit, quantity, source),
            )
            row = cursor.fetchone()
            connection.commit()
            snapshot["capturedAt"] = row["captured_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return snapshot


def create_lifecycle_job(tenant_slug: str, job_type: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    requested_by = (payload.get("requestedBy") or "").strip()
    if requested_by == "":
        raise ValueError("requested_by_required")
    if job_type not in {"onboarding", "offboarding"}:
        raise ValueError("lifecycle_job_type_invalid")

    job = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "jobType": job_type,
        "status": "queued",
        "requestedBy": requested_by,
        "payload": payload.get("payload") or {},
        "createdAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["jobs"].append(job)
        return job

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.lifecycle_jobs (
                  tenant_id, public_id, job_type, status, requested_by, payload_json
                )
                VALUES (%s, %s, %s, 'queued', %s, %s::jsonb)
                RETURNING created_at
                """,
                (tenant_id, job["publicId"], job_type, requested_by, json.dumps(payload.get("payload") or {})),
            )
            row = cursor.fetchone()
            connection.commit()
            job["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return job


def list_lifecycle_jobs(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [job for job in IN_MEMORY_STATE["jobs"] if job["tenantSlug"] == slug]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT job.public_id, job.job_type, job.status, job.requested_by, job.payload_json, job.created_at, job.completed_at
                FROM platform_control.lifecycle_jobs AS job
                JOIN identity.tenants AS tenant ON tenant.id = job.tenant_id
                WHERE tenant.slug = %s
                ORDER BY job.created_at DESC
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "jobType": row["job_type"],
                    "status": row["status"],
                    "requestedBy": row["requested_by"],
                    "payload": row["payload_json"] or {},
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "completedAt": None if row["completed_at"] is None else row["completed_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]
