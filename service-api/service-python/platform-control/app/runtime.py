from __future__ import annotations

from datetime import datetime, timezone
import json
import os
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "entitlements": [],
    "metering": [],
    "jobs": [],
    "job_events": [],
    "quotas": [],
    "blocks": [],
    "provider_defaults": [],
}

CAPABILITY_CATALOG = [
    {"capabilityKey": "catalog.items", "module": "catalog", "defaultEnabled": True, "category": "product"},
    {"capabilityKey": "catalog.bulk", "module": "catalog", "defaultEnabled": True, "category": "integration"},
    {"capabilityKey": "support.cases", "module": "support", "defaultEnabled": False, "category": "service"},
    {"capabilityKey": "notifications.center", "module": "notification", "defaultEnabled": True, "category": "communication"},
    {"capabilityKey": "engagement.providers.meta_ads", "module": "engagement", "defaultEnabled": False, "category": "providers"},
    {"capabilityKey": "engagement.providers.resend", "module": "engagement", "defaultEnabled": False, "category": "providers"},
    {"capabilityKey": "engagement.providers.whatsapp_cloud", "module": "engagement", "defaultEnabled": False, "category": "providers"},
    {"capabilityKey": "billing.pix", "module": "billing", "defaultEnabled": False, "category": "payments"},
    {"capabilityKey": "billing.webhook_reconciliation", "module": "billing", "defaultEnabled": True, "category": "payments"},
    {"capabilityKey": "documents.external_storage", "module": "documents", "defaultEnabled": False, "category": "storage"},
    {"capabilityKey": "documents.digital_signature", "module": "documents", "defaultEnabled": False, "category": "documents"},
    {"capabilityKey": "crm.cnpj_enrichment", "module": "crm", "defaultEnabled": False, "category": "enrichment"},
    {"capabilityKey": "webhook_hub.outbound_webhooks", "module": "webhook-hub", "defaultEnabled": True, "category": "integration"},
]

PROVIDER_CATALOG = [
    {
        "capabilityKey": "billing.pix",
        "providerKey": "asaas",
        "providerType": "payment_gateway",
        "critical": True,
        "fallbackAllowed": False,
        "envKey": "BILLING_ASAAS_API_KEY",
        "defaultMode": "unconfigured",
    },
    {
        "capabilityKey": "billing.pix",
        "providerKey": "mercado_pago",
        "providerType": "payment_gateway",
        "critical": True,
        "fallbackAllowed": False,
        "envKey": "BILLING_MERCADO_PAGO_ACCESS_TOKEN",
        "defaultMode": "unconfigured",
    },
    {
        "capabilityKey": "billing.pix",
        "providerKey": "stripe",
        "providerType": "payment_gateway",
        "critical": True,
        "fallbackAllowed": False,
        "envKey": "BILLING_STRIPE_SECRET_KEY",
        "defaultMode": "unconfigured",
    },
    {
        "capabilityKey": "engagement.providers.resend",
        "providerKey": "resend",
        "providerType": "transactional_email",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "ENGAGEMENT_RESEND_API_KEY",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "engagement.providers.whatsapp_cloud",
        "providerKey": "whatsapp_cloud",
        "providerType": "messaging",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "ENGAGEMENT_WHATSAPP_ACCESS_TOKEN",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "engagement.providers.meta_ads",
        "providerKey": "meta_ads",
        "providerType": "ads",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "ENGAGEMENT_META_ADS_ACCESS_TOKEN",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "documents.external_storage",
        "providerKey": "local",
        "providerType": "storage",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": None,
        "defaultMode": "manual",
    },
    {
        "capabilityKey": "documents.external_storage",
        "providerKey": "s3_compatible",
        "providerType": "storage",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "DOCUMENTS_STORAGE_BUCKET",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "documents.external_storage",
        "providerKey": "cloudflare_r2",
        "providerType": "storage",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "DOCUMENTS_R2_ACCOUNT_ID",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "documents.digital_signature",
        "providerKey": "local",
        "providerType": "digital_signature",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": None,
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "documents.digital_signature",
        "providerKey": "clicksign",
        "providerType": "digital_signature",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "DOCUMENTS_CLICKSIGN_API_KEY",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "documents.digital_signature",
        "providerKey": "docusign",
        "providerType": "digital_signature",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "DOCUMENTS_DOCUSIGN_ACCESS_TOKEN",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "crm.cnpj_enrichment",
        "providerKey": "receita_ws",
        "providerType": "enrichment",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "CRM_CNPJ_PROVIDER_TOKEN",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "crm.cnpj_enrichment",
        "providerKey": "conecta",
        "providerType": "enrichment",
        "critical": False,
        "fallbackAllowed": True,
        "envKey": "CRM_CONECTA_CNPJ_API_KEY",
        "defaultMode": "fallback",
    },
    {
        "capabilityKey": "webhook_hub.outbound_webhooks",
        "providerKey": "tenant_outbound",
        "providerType": "outbound_webhook",
        "critical": True,
        "fallbackAllowed": False,
        "envKey": "WEBHOOK_HUB_OUTBOUND_SIGNING_SECRET",
        "defaultMode": "unconfigured",
    },
]


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def _is_env_configured(env_key: str | None) -> bool:
    if env_key is None:
        return False
    return (os.getenv(env_key, "").strip()) != ""


def _provider_catalog_item(capability_key: str, provider_key: str) -> dict:
    for item in PROVIDER_CATALOG:
        if item["capabilityKey"] == capability_key and item["providerKey"] == provider_key:
            return item
    raise ValueError("provider_catalog_item_not_found")


def _build_provider_default_template(capability_key: str) -> dict | None:
    candidates = [item for item in PROVIDER_CATALOG if item["capabilityKey"] == capability_key]
    if not candidates:
        return None

    def sort_key(item: dict) -> tuple[int, int, str]:
        priority = {
            "manual": 0,
            "fallback": 1,
            "configured": 2,
            "unconfigured": 3,
            "disabled": 4,
        }.get(item["defaultMode"], 9)
        env_priority = 0 if _is_env_configured(item.get("envKey")) else 1
        return (env_priority, priority, item["providerKey"])

    selected = min(candidates, key=sort_key)
    configured = _is_env_configured(selected.get("envKey"))
    mode = "configured" if configured else selected["defaultMode"]

    return {
        "publicId": str(uuid.uuid4()),
        "capabilityKey": selected["capabilityKey"],
        "providerKey": selected["providerKey"],
        "providerType": selected["providerType"],
        "mode": mode,
        "critical": selected["critical"],
        "fallbackAllowed": selected["fallbackAllowed"],
        "envKey": selected["envKey"],
        "source": "template-default",
        "configured": configured or mode in {"manual", "fallback"},
        "metadata": {
            "activation": "env",
            "supportsFallback": selected["fallbackAllowed"],
        },
        "updatedAt": utc_now(),
    }


def _merge_provider_defaults(records: list[dict]) -> list[dict]:
    merged = {item["capabilityKey"]: item for item in records}
    for capability in sorted({item["capabilityKey"] for item in PROVIDER_CATALOG}):
        merged.setdefault(capability, _build_provider_default_template(capability))
    return [item for _, item in sorted(merged.items(), key=lambda entry: entry[0])]


def _find_tenant_id(cursor, tenant_slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (tenant_slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def _normalize_limit(limit: int | None, fallback: int = 50) -> int:
    if limit is None:
        return fallback
    if limit <= 0:
        return fallback
    return min(limit, 100)


def _paginate(records: list[dict], cursor: str | None, limit: int, cursor_field: str = "publicId") -> dict:
    page_limit = _normalize_limit(limit)
    start_index = 0
    if cursor:
        for index, record in enumerate(records):
            if str(record.get(cursor_field, "")) == cursor:
                start_index = index + 1
                break
    page_items = records[start_index : start_index + page_limit]
    next_cursor = None
    if start_index + page_limit < len(records) and page_items:
        next_cursor = str(page_items[-1].get(cursor_field, ""))
    return {
        "items": page_items,
        "pageInfo": {
            "cursor": cursor,
            "limit": page_limit,
            "returned": len(page_items),
            "nextCursor": next_cursor,
            "hasMore": next_cursor is not None,
        },
    }


def list_capability_catalog() -> list[dict]:
    return CAPABILITY_CATALOG


def list_provider_catalog() -> list[dict]:
    catalog = []
    for item in PROVIDER_CATALOG:
        configured = _is_env_configured(item.get("envKey"))
        mode = "configured" if configured else item["defaultMode"]
        catalog.append(
            {
                "capabilityKey": item["capabilityKey"],
                "providerKey": item["providerKey"],
                "providerType": item["providerType"],
                "critical": item["critical"],
                "fallbackAllowed": item["fallbackAllowed"],
                "envKey": item["envKey"],
                "mode": mode,
                "status": "ready" if configured else mode,
                "configured": configured or mode in {"manual", "fallback"},
            }
        )
    return sorted(catalog, key=lambda entry: (entry["capabilityKey"], entry["providerKey"]))


def list_provider_defaults(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["provider_defaults"] if item["tenantSlug"] == slug]
        return _merge_provider_defaults(records)

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT provider.public_id, provider.capability_key, provider.provider_key, provider.provider_type,
                       provider.mode, provider.critical, provider.fallback_allowed, provider.env_key,
                       provider.source, provider.configured, provider.metadata_json, provider.updated_at
                FROM platform_control.provider_defaults AS provider
                JOIN identity.tenants AS tenant ON tenant.id = provider.tenant_id
                WHERE tenant.slug = %s
                ORDER BY provider.capability_key
                """,
                (slug,),
            )
            rows = [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "capabilityKey": row["capability_key"],
                    "providerKey": row["provider_key"],
                    "providerType": row["provider_type"],
                    "mode": row["mode"],
                    "critical": row["critical"],
                    "fallbackAllowed": row["fallback_allowed"],
                    "envKey": row["env_key"],
                    "source": row["source"],
                    "configured": row["configured"],
                    "metadata": row["metadata_json"] or {},
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]
            return _merge_provider_defaults(rows)


def upsert_provider_default(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    normalized_capability = capability_key.strip()
    provider_key = str(payload.get("providerKey") or "").strip()
    if normalized_capability == "":
        raise ValueError("capability_key_required")
    if provider_key == "":
        raise ValueError("provider_key_required")

    catalog_item = _provider_catalog_item(normalized_capability, provider_key)
    configured = _is_env_configured(catalog_item.get("envKey"))
    requested_mode = str(payload.get("mode") or catalog_item["defaultMode"]).strip() or catalog_item["defaultMode"]
    if requested_mode not in {"configured", "fallback", "manual", "unconfigured", "disabled"}:
        raise ValueError("provider_mode_invalid")
    if requested_mode == "configured" and not configured:
        raise ValueError("provider_env_missing")

    provider_default = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "capabilityKey": normalized_capability,
        "providerKey": provider_key,
        "providerType": catalog_item["providerType"],
        "mode": requested_mode if not configured else "configured",
        "critical": catalog_item["critical"],
        "fallbackAllowed": catalog_item["fallbackAllowed"],
        "envKey": catalog_item["envKey"],
        "source": str(payload.get("source") or "tenant-default").strip() or "tenant-default",
        "configured": configured or requested_mode in {"manual", "fallback"},
        "metadata": payload.get("metadata") or {"activation": "env", "supportsFallback": catalog_item["fallbackAllowed"]},
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        existing = next(
            (item for item in IN_MEMORY_STATE["provider_defaults"] if item["tenantSlug"] == slug and item["capabilityKey"] == normalized_capability),
            None,
        )
        if existing is not None:
            existing.update(provider_default)
            provider_default["publicId"] = existing["publicId"]
            return existing
        IN_MEMORY_STATE["provider_defaults"].append(provider_default)
        return provider_default

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.provider_defaults (
                  tenant_id, public_id, capability_key, provider_key, provider_type, mode,
                  critical, fallback_allowed, env_key, source, configured, metadata_json
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s::jsonb)
                ON CONFLICT (tenant_id, capability_key)
                DO UPDATE SET
                  provider_key = EXCLUDED.provider_key,
                  provider_type = EXCLUDED.provider_type,
                  mode = EXCLUDED.mode,
                  critical = EXCLUDED.critical,
                  fallback_allowed = EXCLUDED.fallback_allowed,
                  env_key = EXCLUDED.env_key,
                  source = EXCLUDED.source,
                  configured = EXCLUDED.configured,
                  metadata_json = EXCLUDED.metadata_json,
                  updated_at = NOW()
                RETURNING public_id, updated_at
                """,
                (
                    tenant_id,
                    provider_default["publicId"],
                    normalized_capability,
                    provider_key,
                    provider_default["providerType"],
                    provider_default["mode"],
                    provider_default["critical"],
                    provider_default["fallbackAllowed"],
                    provider_default["envKey"],
                    provider_default["source"],
                    provider_default["configured"],
                    json.dumps(provider_default["metadata"]),
                ),
            )
            row = cursor.fetchone()
            connection.commit()
            provider_default["publicId"] = row["public_id"]
            provider_default["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return provider_default


def list_entitlements(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["entitlements"] if item["tenantSlug"] == slug]
        return sorted(records, key=lambda item: item["capabilityKey"])

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


def list_entitlements_page(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    payload = _paginate(list_entitlements(tenant_slug), cursor, limit)
    payload["tenantSlug"] = _tenant_slug(tenant_slug)
    return payload


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


def bulk_upsert_entitlements(tenant_slug: str, payload: dict) -> dict:
    items = payload.get("items") or []
    results: list[dict] = []
    errors: list[dict] = []
    for index, item in enumerate(items):
        capability_key = str(item.get("capabilityKey") or "").strip()
        if capability_key == "":
            errors.append({"index": index, "code": "capability_key_required", "message": "Capability key is required."})
            continue
        try:
            results.append(upsert_entitlement(tenant_slug, capability_key, item))
        except ValueError as error:
            errors.append({"index": index, "code": str(error), "message": "Entitlement payload is invalid."})

    return {
        "tenantSlug": _tenant_slug(tenant_slug),
        "results": results,
        "errors": errors,
        "summary": {
            "requested": len(items),
            "succeeded": len(results),
            "failed": len(errors),
            "partialSuccess": len(results) > 0 and len(errors) > 0,
        },
    }


def list_quotas(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["quotas"] if item["tenantSlug"] == slug]
        return sorted(records, key=lambda item: item["metricKey"])

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT quota.public_id, quota.metric_key, quota.metric_unit, quota.limit_value,
                       quota.enforcement_mode, quota.source, quota.updated_at
                FROM platform_control.quotas AS quota
                JOIN identity.tenants AS tenant ON tenant.id = quota.tenant_id
                WHERE tenant.slug = %s
                ORDER BY quota.metric_key
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "metricKey": row["metric_key"],
                    "metricUnit": row["metric_unit"],
                    "limitValue": row["limit_value"],
                    "enforcementMode": row["enforcement_mode"],
                    "source": row["source"],
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def upsert_quota(tenant_slug: str, metric_key: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    normalized_key = metric_key.strip()
    metric_unit = str(payload.get("metricUnit") or "").strip()
    if normalized_key == "":
        raise ValueError("metric_key_required")
    if metric_unit == "":
        raise ValueError("metric_unit_required")

    quota = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "metricKey": normalized_key,
        "metricUnit": metric_unit,
        "limitValue": int(payload.get("limitValue", 0) or 0),
        "enforcementMode": str(payload.get("enforcementMode") or "soft").strip() or "soft",
        "source": str(payload.get("source") or "manual").strip() or "manual",
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        existing = next((item for item in IN_MEMORY_STATE["quotas"] if item["tenantSlug"] == slug and item["metricKey"] == normalized_key), None)
        if existing is not None:
            existing.update(quota)
            quota["publicId"] = existing["publicId"]
            return existing
        IN_MEMORY_STATE["quotas"].append(quota)
        return quota

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.quotas (
                  tenant_id, public_id, metric_key, metric_unit, limit_value, enforcement_mode, source
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, metric_key)
                DO UPDATE SET
                  metric_unit = EXCLUDED.metric_unit,
                  limit_value = EXCLUDED.limit_value,
                  enforcement_mode = EXCLUDED.enforcement_mode,
                  source = EXCLUDED.source,
                  updated_at = NOW()
                RETURNING public_id, updated_at
                """,
                (tenant_id, quota["publicId"], normalized_key, metric_unit, quota["limitValue"], quota["enforcementMode"], quota["source"]),
            )
            row = cursor.fetchone()
            connection.commit()
            quota["publicId"] = row["public_id"]
            quota["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return quota


def bulk_upsert_quotas(tenant_slug: str, payload: dict) -> dict:
    items = payload.get("items") or []
    results: list[dict] = []
    errors: list[dict] = []
    for index, item in enumerate(items):
        metric_key = str(item.get("metricKey") or "").strip()
        if metric_key == "":
            errors.append({"index": index, "code": "metric_key_required", "message": "Metric key is required."})
            continue
        try:
            results.append(upsert_quota(tenant_slug, metric_key, item))
        except ValueError as error:
            errors.append({"index": index, "code": str(error), "message": "Quota payload is invalid."})

    return {
        "tenantSlug": _tenant_slug(tenant_slug),
        "results": results,
        "errors": errors,
        "summary": {
            "requested": len(items),
            "succeeded": len(results),
            "failed": len(errors),
            "partialSuccess": len(results) > 0 and len(errors) > 0,
        },
    }


def list_blocks(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["blocks"] if item["tenantSlug"] == slug]
        return sorted(records, key=lambda item: item["blockKey"])

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT block.public_id, block.block_key, block.active, block.reason, block.scope, block.source, block.updated_at
                FROM platform_control.tenant_blocks AS block
                JOIN identity.tenants AS tenant ON tenant.id = block.tenant_id
                WHERE tenant.slug = %s
                ORDER BY block.block_key
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "blockKey": row["block_key"],
                    "active": row["active"],
                    "reason": row["reason"],
                    "scope": row["scope"],
                    "source": row["source"],
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def upsert_block(tenant_slug: str, block_key: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    normalized_key = block_key.strip()
    if normalized_key == "":
        raise ValueError("block_key_required")

    block = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "blockKey": normalized_key,
        "active": bool(payload.get("active", False)),
        "reason": str(payload.get("reason") or "").strip(),
        "scope": str(payload.get("scope") or "tenant").strip() or "tenant",
        "source": str(payload.get("source") or "manual").strip() or "manual",
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        existing = next((item for item in IN_MEMORY_STATE["blocks"] if item["tenantSlug"] == slug and item["blockKey"] == normalized_key), None)
        if existing is not None:
            existing.update(block)
            block["publicId"] = existing["publicId"]
            return existing
        IN_MEMORY_STATE["blocks"].append(block)
        return block

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.tenant_blocks (
                  tenant_id, public_id, block_key, active, reason, scope, source
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, block_key)
                DO UPDATE SET
                  active = EXCLUDED.active,
                  reason = EXCLUDED.reason,
                  scope = EXCLUDED.scope,
                  source = EXCLUDED.source,
                  updated_at = NOW()
                RETURNING public_id, updated_at
                """,
                (tenant_id, block["publicId"], normalized_key, block["active"], block["reason"], block["scope"], block["source"]),
            )
            row = cursor.fetchone()
            connection.commit()
            block["publicId"] = row["public_id"]
            block["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return block


def list_metering(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        snapshots = [item for item in IN_MEMORY_STATE["metering"] if item["tenantSlug"] == slug]
        snapshots = sorted(snapshots, key=lambda item: (item["capturedAt"], item["publicId"]), reverse=True)
        return {"tenantSlug": slug, "snapshots": snapshots, "summary": {"metrics": len(snapshots), "totalQuantity": sum(item["quantity"] for item in snapshots)}}

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT snapshot.public_id, snapshot.metric_key, snapshot.metric_unit, snapshot.quantity, snapshot.source, snapshot.captured_at
                FROM platform_control.usage_snapshots AS snapshot
                JOIN identity.tenants AS tenant ON tenant.id = snapshot.tenant_id
                WHERE tenant.slug = %s
                ORDER BY snapshot.captured_at DESC, snapshot.public_id DESC
                LIMIT 200
                """,
                (slug,),
            )
            rows = cursor.fetchall()
            snapshots = [
                {
                    "publicId": row["public_id"],
                    "metricKey": row["metric_key"],
                    "metricUnit": row["metric_unit"],
                    "quantity": row["quantity"],
                    "source": row["source"],
                    "capturedAt": row["captured_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in rows
            ]
            return {"tenantSlug": slug, "snapshots": snapshots, "summary": {"metrics": len(snapshots), "totalQuantity": sum(item["quantity"] for item in snapshots)}}


def list_metering_page(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    data = list_metering(tenant_slug)
    paged = _paginate(data["snapshots"], cursor, limit)
    return {"tenantSlug": data["tenantSlug"], "summary": data["summary"], **paged}


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


def build_usage_summary(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    metering = list_metering(slug)
    quotas = list_quotas(slug)
    blocks = list_blocks(slug)

    aggregated: dict[str, dict] = {}
    for snapshot in metering["snapshots"]:
        current = aggregated.setdefault(
            snapshot["metricKey"],
            {
                "metricKey": snapshot["metricKey"],
                "metricUnit": snapshot["metricUnit"],
                "quantity": 0,
            },
        )
        current["quantity"] += int(snapshot["quantity"])

    quota_map = {item["metricKey"]: item for item in quotas}
    for metric_key, quota in quota_map.items():
        aggregated.setdefault(
            metric_key,
            {
                "metricKey": metric_key,
                "metricUnit": quota.get("metricUnit", ""),
                "quantity": 0,
            },
        )
    metrics = []
    for metric_key, aggregate in sorted(aggregated.items()):
        quota = quota_map.get(metric_key, {})
        limit_value = int(quota.get("limitValue", 0) or 0)
        quantity = int(aggregate["quantity"])
        remaining = None if limit_value <= 0 else max(limit_value - quantity, 0)
        utilization_rate = None if limit_value <= 0 else round(quantity / limit_value, 4)
        status = "ok"
        if limit_value > 0 and quantity >= limit_value:
            status = "limit_reached"
        elif limit_value > 0 and quantity >= int(limit_value * 0.85):
            status = "attention"
        metrics.append(
            {
                "metricKey": metric_key,
                "metricUnit": aggregate["metricUnit"],
                "quantity": quantity,
                "limitValue": limit_value,
                "remaining": remaining,
                "utilizationRate": utilization_rate,
                "enforcementMode": quota.get("enforcementMode", "soft"),
                "status": status,
            }
        )

    active_blocks = [item for item in blocks if item["active"]]
    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "metrics": metrics,
        "summary": {
            "trackedMetrics": len(metrics),
            "activeQuotas": len(quotas),
            "activeBlocks": len(active_blocks),
            "limitReached": sum(1 for item in metrics if item["status"] == "limit_reached"),
            "attention": sum(1 for item in metrics if item["status"] == "attention"),
        },
    }


def build_lifecycle_readiness(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    entitlements = list_entitlements(slug)
    quotas = list_quotas(slug)
    blocks = list_blocks(slug)
    provider_defaults = list_provider_defaults(slug)

    enabled_entitlements = [item for item in entitlements if item["enabled"]]
    active_blocks = [item for item in blocks if item["active"]]
    critical_unconfigured = [
        item for item in provider_defaults if item["critical"] and item["mode"] in {"unconfigured", "disabled"}
    ]
    attention_quotas = [item for item in quotas if item.get("enforcementMode") == "hard"]

    status = "stable"
    if critical_unconfigured or active_blocks:
        status = "attention"
    if len(critical_unconfigured) > 1:
        status = "critical"

    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "providers": {
            "total": len(provider_defaults),
            "configured": sum(1 for item in provider_defaults if item["mode"] == "configured"),
            "fallback": sum(1 for item in provider_defaults if item["mode"] == "fallback"),
            "manual": sum(1 for item in provider_defaults if item["mode"] == "manual"),
            "criticalUnconfigured": len(critical_unconfigured),
            "items": provider_defaults,
        },
        "entitlements": {
            "total": len(entitlements),
            "enabled": len(enabled_entitlements),
        },
        "quotas": {
            "total": len(quotas),
            "hardEnforcement": len(attention_quotas),
        },
        "blocks": {
            "active": len(active_blocks),
            "items": active_blocks,
        },
        "readiness": {
            "status": status,
            "onboardingReady": len(critical_unconfigured) == 0,
            "offboardingReady": True,
            "providerDefaultsReady": len(provider_defaults) > 0,
            "missingCriticalProviders": [item["capabilityKey"] for item in critical_unconfigured],
        },
    }


def build_lifecycle_preview(tenant_slug: str, job_type: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    readiness = build_lifecycle_readiness(slug)
    lifecycle_payload = payload.get("payload") or {}

    if job_type == "onboarding":
        steps = [
            {"key": "seed-tenant", "status": "ready", "summary": "Bootstrap de tenant e catalogos basicos."},
            {"key": "apply-entitlements", "status": "ready", "summary": "Aplicar entitlements e limites operacionais."},
            {
                "key": "configure-providers",
                "status": "ready" if readiness["readiness"]["onboardingReady"] else "attention",
                "summary": "Aplicar defaults de provider e validar credenciais criticas.",
            },
            {
                "key": "seed",
                "status": "ready",
                "summary": f"Seed opcional selecionado: {lifecycle_payload.get('seed', 'none')}.",
            },
        ]
    else:
        steps = [
            {"key": "export-data", "status": "ready", "summary": "Preparar exportacao de dados e artefatos."},
            {"key": "revoke-access", "status": "ready", "summary": "Revogar acessos e providers configurados."},
            {
                "key": "retention-plan",
                "status": "ready",
                "summary": f"Plano de retencao: {lifecycle_payload.get('mode', 'retention-default')}.",
            },
            {"key": "purge-controls", "status": "ready", "summary": "Aplicar purge seletivo conforme governanca."},
        ]

    return {
        "tenantSlug": slug,
        "jobType": job_type,
        "generatedAt": utc_now(),
        "steps": steps,
        "readiness": readiness["readiness"],
        "providers": {
            "recommendedDefaults": [
                {
                    "capabilityKey": item["capabilityKey"],
                    "providerKey": item["providerKey"],
                    "mode": item["mode"],
                }
                for item in readiness["providers"]["items"]
            ],
        },
    }


def create_lifecycle_job(tenant_slug: str, job_type: str, payload: dict, idempotency_key: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    requested_by = (payload.get("requestedBy") or "").strip()
    if requested_by == "":
        raise ValueError("requested_by_required")
    if job_type not in {"onboarding", "offboarding"}:
        raise ValueError("lifecycle_job_type_invalid")

    normalized_idempotency_key = (idempotency_key or "").strip() or None
    preview = build_lifecycle_preview(slug, job_type, payload)
    job_payload = {
        "input": payload.get("payload") or {},
        "preview": preview,
    }
    if settings.repository_driver != "postgres":
        if normalized_idempotency_key is not None:
            existing = next(
                (
                    item
                    for item in IN_MEMORY_STATE["jobs"]
                    if item["tenantSlug"] == slug and item["jobType"] == job_type and item.get("idempotencyKey") == normalized_idempotency_key
                ),
                None,
            )
            if existing is not None:
                return get_lifecycle_job(slug, existing["publicId"]) or existing

        job = {
            "publicId": str(uuid.uuid4()),
            "tenantSlug": slug,
            "jobType": job_type,
            "status": "queued",
            "requestedBy": requested_by,
            "idempotencyKey": normalized_idempotency_key,
            "payload": job_payload,
            "createdAt": utc_now(),
            "startedAt": None,
            "completedAt": None,
            "failedAt": None,
            "cancelledAt": None,
            "failureReason": None,
        }
        IN_MEMORY_STATE["jobs"].append(job)
        _append_job_event_memory(job["publicId"], slug, "queued", f"{job_type} requested")
        return get_lifecycle_job(slug, job["publicId"]) or job

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            if normalized_idempotency_key is not None:
                cursor.execute(
                    """
                    SELECT public_id
                    FROM platform_control.lifecycle_jobs
                    WHERE tenant_id = %s AND job_type = %s AND idempotency_key = %s
                    """,
                    (tenant_id, job_type, normalized_idempotency_key),
                )
                existing = cursor.fetchone()
                if existing is not None:
                    connection.commit()
                    return get_lifecycle_job(slug, existing["public_id"])

            public_id = str(uuid.uuid4())
            cursor.execute(
                """
                INSERT INTO platform_control.lifecycle_jobs (
                  tenant_id, public_id, job_type, status, requested_by, payload_json, idempotency_key
                )
                VALUES (%s, %s, %s, 'queued', %s, %s::jsonb, %s)
                RETURNING created_at
                """,
                (tenant_id, public_id, job_type, requested_by, json.dumps(job_payload), normalized_idempotency_key),
            )
            created_row = cursor.fetchone()
            cursor.execute(
                """
                INSERT INTO platform_control.lifecycle_job_events (
                  job_id, public_id, status, summary
                )
                SELECT id, %s, 'queued', %s
                FROM platform_control.lifecycle_jobs
                WHERE public_id = %s
                """,
                (str(uuid.uuid4()), f"{job_type} requested", public_id),
            )
            connection.commit()
            job = get_lifecycle_job(slug, public_id)
            if job is None:
                raise ValueError("lifecycle_job_not_found")
            job["createdAt"] = created_row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return job


def _append_job_event_memory(public_id: str, tenant_slug: str, status: str, summary: str) -> None:
    IN_MEMORY_STATE["job_events"].append(
        {
            "publicId": str(uuid.uuid4()),
            "tenantSlug": tenant_slug,
            "jobPublicId": public_id,
            "status": status,
            "summary": summary,
            "createdAt": utc_now(),
        }
    )


def list_lifecycle_jobs(tenant_slug: str) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        jobs = [job for job in IN_MEMORY_STATE["jobs"] if job["tenantSlug"] == slug]
        jobs = sorted(jobs, key=lambda item: (item["createdAt"], item["publicId"]), reverse=True)
        return [get_lifecycle_job(slug, job["publicId"]) or job for job in jobs]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT job.public_id
                FROM platform_control.lifecycle_jobs AS job
                JOIN identity.tenants AS tenant ON tenant.id = job.tenant_id
                WHERE tenant.slug = %s
                ORDER BY job.created_at DESC, job.public_id DESC
                """,
                (slug,),
            )
            rows = cursor.fetchall()
            return [job for row in rows if (job := get_lifecycle_job(slug, row["public_id"])) is not None]


def list_lifecycle_jobs_page(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    payload = _paginate(list_lifecycle_jobs(tenant_slug), cursor, limit)
    payload["tenantSlug"] = _tenant_slug(tenant_slug)
    return payload


def get_lifecycle_job(tenant_slug: str, public_id: str) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        job = next((item for item in IN_MEMORY_STATE["jobs"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
        if job is None:
            return None
        events = [
            event
            for event in IN_MEMORY_STATE["job_events"]
            if event["tenantSlug"] == slug and event["jobPublicId"] == public_id
        ]
        return {
            **job,
            "events": sorted(events, key=lambda item: item["createdAt"]),
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT job.public_id, job.job_type, job.status, job.requested_by, job.payload_json,
                       job.idempotency_key, job.created_at, job.started_at, job.completed_at,
                       job.failed_at, job.cancelled_at, job.failure_reason
                FROM platform_control.lifecycle_jobs AS job
                JOIN identity.tenants AS tenant ON tenant.id = job.tenant_id
                WHERE tenant.slug = %s AND job.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None

            cursor.execute(
                """
                SELECT event.public_id, event.status, event.summary, event.created_at
                FROM platform_control.lifecycle_job_events AS event
                JOIN platform_control.lifecycle_jobs AS job ON job.id = event.job_id
                WHERE job.public_id = %s
                ORDER BY event.created_at
                """,
                (public_id,),
            )
            events = [
                {
                    "publicId": str(event["public_id"]),
                    "status": event["status"],
                    "summary": event["summary"],
                    "createdAt": event["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for event in cursor.fetchall()
            ]

    return {
        "publicId": str(row["public_id"]),
        "tenantSlug": slug,
        "jobType": row["job_type"],
        "status": row["status"],
        "requestedBy": row["requested_by"],
        "payload": row["payload_json"] or {},
        "idempotencyKey": row["idempotency_key"],
        "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "startedAt": None if row["started_at"] is None else row["started_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "completedAt": None if row["completed_at"] is None else row["completed_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "failedAt": None if row["failed_at"] is None else row["failed_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "cancelledAt": None if row["cancelled_at"] is None else row["cancelled_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "failureReason": row["failure_reason"],
        "events": events,
    }


def transition_lifecycle_job(tenant_slug: str, public_id: str, action: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    action_map = {
        "start": ("running", "Job execution started."),
        "complete": ("completed", "Job execution completed."),
        "fail": ("failed", "Job execution failed."),
        "cancel": ("cancelled", "Job execution cancelled."),
    }
    if action not in action_map:
        raise ValueError("lifecycle_action_invalid")

    next_status, default_summary = action_map[action]
    summary = str(payload.get("summary") or default_summary).strip() or default_summary
    failure_reason = str(payload.get("failureReason") or "").strip() or None

    if settings.repository_driver != "postgres":
        job = next((item for item in IN_MEMORY_STATE["jobs"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
        if job is None:
            raise ValueError("lifecycle_job_not_found")
        job["status"] = next_status
        now = utc_now()
        if next_status == "running":
            job["startedAt"] = now
        elif next_status == "completed":
            job["completedAt"] = now
        elif next_status == "failed":
            job["failedAt"] = now
            job["failureReason"] = failure_reason
        elif next_status == "cancelled":
            job["cancelledAt"] = now
        _append_job_event_memory(public_id, slug, next_status, summary)
        return get_lifecycle_job(slug, public_id) or job

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE platform_control.lifecycle_jobs AS job
                SET status = %s,
                    started_at = CASE WHEN %s = 'running' THEN NOW() ELSE job.started_at END,
                    completed_at = CASE WHEN %s = 'completed' THEN NOW() ELSE job.completed_at END,
                    failed_at = CASE WHEN %s = 'failed' THEN NOW() ELSE job.failed_at END,
                    cancelled_at = CASE WHEN %s = 'cancelled' THEN NOW() ELSE job.cancelled_at END,
                    failure_reason = CASE WHEN %s = 'failed' THEN %s ELSE job.failure_reason END
                FROM identity.tenants AS tenant
                WHERE tenant.id = job.tenant_id
                  AND tenant.slug = %s
                  AND job.public_id = %s
                RETURNING job.id
                """,
                (next_status, next_status, next_status, next_status, next_status, next_status, failure_reason, slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                raise ValueError("lifecycle_job_not_found")
            cursor.execute(
                """
                INSERT INTO platform_control.lifecycle_job_events (job_id, public_id, status, summary)
                VALUES (%s, %s, %s, %s)
                """,
                (row["id"], str(uuid.uuid4()), next_status, summary),
            )
            connection.commit()

    job = get_lifecycle_job(slug, public_id)
    if job is None:
        raise ValueError("lifecycle_job_not_found")
    return job
