from __future__ import annotations

import base64
from datetime import datetime, timezone
import hashlib
import json
import os
from urllib import error as urlerror
from urllib import parse, request
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
    "go_live_rollouts": [],
    "go_live_events": [],
    "incidents": [],
    "incident_timeline_events": [],
    "incident_actions": [],
    "postmortems": [],
    "policy_decisions": [],
    "timeline_events": [],
    "approval_requests": [],
    "runbook_runs": [],
    "runbook_steps": [],
    "evidence_records": [],
    "event_mesh_events": [],
    "event_mesh_dead_letters": [],
    "event_mesh_consumers": [],
    "tenant_runtime_profiles": [],
    "tenant_runtime_quotas": [],
    "tenant_maintenance_windows": [],
    "contract_snapshots": [],
    "contract_diffs": [],
    "contract_breaking_changes": [],
    "provider_activation_runs": [],
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

POLICY_CATALOG = [
    {
        "policyKey": "exports.require-review",
        "version": "1.0",
        "domain": "search",
        "action": "data.export",
        "effect": "review",
        "priority": 10,
        "reason": "Exports can expose sensitive data and must be approved.",
    },
    {
        "policyKey": "quotas.allow-ops-change",
        "version": "1.0",
        "domain": "platform-control",
        "action": "quota.change",
        "effect": "allow",
        "priority": 20,
        "reason": "Quota changes are allowed when actor, tenant and justification are present.",
    },
    {
        "policyKey": "ai.deny-mutation-tools",
        "version": "1.0",
        "domain": "ai-governance",
        "action": "ai.tool.write",
        "effect": "deny",
        "priority": 5,
        "reason": "AI assistant runs are read-only unless a future policy explicitly allows mutation.",
    },
    {
        "policyKey": "incidents.review-sev1",
        "version": "1.0",
        "domain": "platform-control",
        "action": "incident.escalate",
        "effect": "review",
        "priority": 10,
        "reason": "SEV1 escalation requires incident command approval and evidence.",
    },
    {
        "policyKey": "go-live.review-rollback",
        "version": "1.0",
        "domain": "platform-control",
        "action": "go-live.rollback",
        "effect": "review",
        "priority": 10,
        "reason": "Rollback needs approval unless an active SEV1 runbook is executing.",
    },
    {
        "policyKey": "providers.review-critical-fallback",
        "version": "1.0",
        "domain": "platform-control",
        "action": "provider.fallback",
        "effect": "review",
        "priority": 15,
        "reason": "Critical provider fallback changes must be reviewed.",
    },
    {
        "policyKey": "billing.review-recovery",
        "version": "1.0",
        "domain": "billing",
        "action": "billing.recovery",
        "effect": "review",
        "priority": 15,
        "reason": "Billing recovery can affect customer access and financial posture.",
    },
]

RUNBOOK_CATALOG = [
    {
        "runbookKey": "provider-degraded",
        "title": "Provider degraded",
        "domain": "integrations",
        "steps": ["confirm_provider_status", "switch_to_fallback_if_approved", "notify_owner", "capture_evidence"],
    },
    {
        "runbookKey": "webhook-dlq-growing",
        "title": "Webhook DLQ growing",
        "domain": "webhook-hub",
        "steps": ["inspect_dlq", "sample_payloads", "request_requeue_approval", "monitor_delivery"],
    },
    {
        "runbookKey": "tenant-over-quota",
        "title": "Tenant over quota",
        "domain": "platform-control",
        "steps": ["read_usage_summary", "evaluate_quota_policy", "request_quota_change", "record_customer_notice"],
    },
    {
        "runbookKey": "sev1-incident",
        "title": "SEV1 incident",
        "domain": "incident-command",
        "steps": ["open_bridge", "assign_commander", "request_escalation_approval", "publish_postmortem"],
    },
    {
        "runbookKey": "openapi-drift",
        "title": "OpenAPI drift",
        "domain": "contracts",
        "steps": ["run_contract_suite", "identify_owner", "block_release", "capture_contract_evidence"],
    },
    {
        "runbookKey": "legal-hold-export-block",
        "title": "Legal hold export block",
        "domain": "search",
        "steps": ["confirm_hold", "request_legal_approval", "prepare_redacted_export", "capture_export_evidence"],
    },
    {
        "runbookKey": "go-live-rollback",
        "title": "Go-live rollback",
        "domain": "go-live",
        "steps": ["freeze_wave", "request_rollback_approval", "execute_rollback", "publish_rollout_evidence"],
    },
]

EVENT_STREAM_CATALOG = [
    {"streamKey": "crm.customer", "domain": "crm", "schemaVersion": "1.0", "retention": "p3y", "critical": True},
    {"streamKey": "sales.opportunity", "domain": "sales", "schemaVersion": "1.0", "retention": "p3y", "critical": True},
    {"streamKey": "billing.invoice", "domain": "billing", "schemaVersion": "1.1", "retention": "p5y", "critical": True},
    {"streamKey": "finance.ledger", "domain": "finance", "schemaVersion": "1.0", "retention": "p5y", "critical": True},
    {"streamKey": "documents.lifecycle", "domain": "documents", "schemaVersion": "1.0", "retention": "p5y", "critical": True},
    {"streamKey": "rentals.contract", "domain": "rentals", "schemaVersion": "1.0", "retention": "p3y", "critical": False},
    {"streamKey": "workflow.execution", "domain": "workflow", "schemaVersion": "1.0", "retention": "p2y", "critical": True},
    {"streamKey": "webhook.delivery", "domain": "webhook-hub", "schemaVersion": "1.0", "retention": "p2y", "critical": True},
    {"streamKey": "platform.governance", "domain": "platform-control", "schemaVersion": "1.2", "retention": "p5y", "critical": True},
]

CONTRACT_DOMAIN_CATALOG = [
    {"contractKey": "analytics.openapi", "service": "analytics", "kind": "openapi", "currentVersion": "1.4.0"},
    {"contractKey": "platform-control.openapi", "service": "platform-control", "kind": "openapi", "currentVersion": "1.4.0"},
    {"contractKey": "billing.events", "service": "billing", "kind": "async-event", "currentVersion": "1.1.0"},
    {"contractKey": "workflow.events", "service": "workflow-runtime", "kind": "async-event", "currentVersion": "1.0.0"},
]

PROVIDER_ACTIVATION_CATALOG = [
    {
        "providerKey": "stripe",
        "capabilityKey": "billing.pix",
        "domain": "billing",
        "credentialKey": "BILLING_STRIPE_SECRET_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "payment_intent.create"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "asaas",
        "capabilityKey": "billing.pix",
        "domain": "billing",
        "credentialKey": "BILLING_ASAAS_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "mercado_pago",
        "capabilityKey": "billing.pix",
        "domain": "billing",
        "credentialKey": "BILLING_MERCADO_PAGO_ACCESS_TOKEN",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "resend",
        "capabilityKey": "engagement.providers.resend",
        "domain": "engagement",
        "credentialKey": "ENGAGEMENT_RESEND_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "email.send"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "openai",
        "capabilityKey": "ai.llm",
        "domain": "ai-governance",
        "credentialKey": "OPENAI_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "response.create"],
        "fallback": "deterministic_local_answer",
    },
    {
        "providerKey": "docusign",
        "capabilityKey": "documents.digital_signature",
        "domain": "documents",
        "credentialKey": "DOCUMENTS_DOCUSIGN_ACCESS_TOKEN",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "clicksign",
        "capabilityKey": "documents.digital_signature",
        "domain": "documents",
        "credentialKey": "DOCUMENTS_CLICKSIGN_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "whatsapp_cloud",
        "capabilityKey": "engagement.providers.whatsapp_cloud",
        "domain": "engagement",
        "credentialKey": "ENGAGEMENT_WHATSAPP_ACCESS_TOKEN",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "aws_textract",
        "capabilityKey": "documents.ocr",
        "domain": "document_intelligence",
        "credentialKey": "AWS_TEXTRACT_ACCESS_KEY_ID",
        "requiredCredentialKeys": ["AWS_TEXTRACT_ACCESS_KEY_ID", "AWS_TEXTRACT_SECRET_ACCESS_KEY", "AWS_TEXTRACT_REGION"],
        "mode": "sdk_or_sigv4_runtime_adapter",
        "supportedActions": ["credential_check"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "google_document_ai",
        "capabilityKey": "documents.ocr",
        "domain": "document_intelligence",
        "credentialKey": "GOOGLE_DOCUMENT_AI_CREDENTIALS_JSON",
        "requiredCredentialKeys": ["GOOGLE_DOCUMENT_AI_PROCESSOR", "GOOGLE_DOCUMENT_AI_CREDENTIALS_JSON"],
        "mode": "oauth_service_account_runtime_adapter",
        "supportedActions": ["credential_check"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "focus_nfe",
        "capabilityKey": "fiscal.issuance",
        "domain": "fiscal",
        "credentialKey": "FISCAL_FOCUS_NFE_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "enotas",
        "capabilityKey": "fiscal.issuance",
        "domain": "fiscal",
        "credentialKey": "FISCAL_ENOTAS_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "serpro_cnpj",
        "capabilityKey": "crm.cnpj_enrichment",
        "domain": "registry_enrichment",
        "credentialKey": "CRM_SERPRO_CLIENT_SECRET",
        "requiredCredentialKeys": ["CRM_SERPRO_CLIENT_ID", "CRM_SERPRO_CLIENT_SECRET"],
        "mode": "live_api",
        "supportedActions": ["connection_test"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "brasilapi",
        "capabilityKey": "crm.cnpj_enrichment",
        "domain": "registry_enrichment",
        "credentialKey": None,
        "credentialRequired": False,
        "mode": "public_api",
        "supportedActions": ["connection_test", "cnpj.lookup"],
        "fallback": "public_api_without_key",
    },
    {
        "providerKey": "viacep",
        "capabilityKey": "crm.cep_enrichment",
        "domain": "registry_enrichment",
        "credentialKey": None,
        "credentialRequired": False,
        "mode": "public_api",
        "supportedActions": ["connection_test", "cep.lookup"],
        "fallback": "public_api_without_key",
    },
    {
        "providerKey": "alpha_vantage",
        "capabilityKey": "market.data",
        "domain": "market_macro_risk",
        "credentialKey": "MARKET_ALPHA_VANTAGE_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "fx.lookup"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "fixer",
        "capabilityKey": "market.fx",
        "domain": "market_macro_risk",
        "credentialKey": "MARKET_FIXER_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "fx.latest"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "bcb_sgs",
        "capabilityKey": "market.macro",
        "domain": "market_macro_risk",
        "credentialKey": None,
        "credentialRequired": False,
        "mode": "public_api",
        "supportedActions": ["connection_test", "series.latest"],
        "fallback": "public_api_without_key",
    },
    {
        "providerKey": "bcb_ptax",
        "capabilityKey": "market.fx_reference",
        "domain": "market_macro_risk",
        "credentialKey": None,
        "credentialRequired": False,
        "mode": "public_api",
        "supportedActions": ["connection_test"],
        "fallback": "public_api_without_key",
    },
    {
        "providerKey": "newsapi",
        "capabilityKey": "external_risk.news",
        "domain": "external_risk_feed",
        "credentialKey": "NEWSAPI_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "news.search"],
        "fallback": "unavailable_without_key",
    },
    {
        "providerKey": "gdelt",
        "capabilityKey": "external_risk.news",
        "domain": "external_risk_feed",
        "credentialKey": None,
        "credentialRequired": False,
        "mode": "public_api",
        "supportedActions": ["connection_test", "news.search"],
        "fallback": "public_api_without_key",
    },
    {
        "providerKey": "alpha_vantage_news",
        "capabilityKey": "external_risk.market_news",
        "domain": "external_risk_feed",
        "credentialKey": "MARKET_ALPHA_VANTAGE_API_KEY",
        "mode": "live_api",
        "supportedActions": ["connection_test", "news.sentiment"],
        "fallback": "unavailable_without_key",
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


def list_feature_flags_page(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    capability_map = {item["capabilityKey"]: item for item in CAPABILITY_CATALOG}
    payload = list_entitlements_page(tenant_slug, cursor, limit)
    payload["items"] = [
        {
            **item,
            "flagKey": item["capabilityKey"],
            "module": capability_map.get(item["capabilityKey"], {}).get("module", "unknown"),
            "category": capability_map.get(item["capabilityKey"], {}).get("category", "unknown"),
            "defaultEnabled": capability_map.get(item["capabilityKey"], {}).get("defaultEnabled", False),
        }
        for item in payload["items"]
    ]
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
    total_quantity = sum(int(item["quantity"]) for item in metrics)
    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "metrics": metrics,
        "summary": {
            "trackedMetrics": len(metrics),
            "totalQuantity": total_quantity,
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
    allowed_statuses = {
        "start": {"queued"},
        "complete": {"running"},
        "fail": {"running"},
        "cancel": {"queued", "running"},
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
        if job["status"] not in allowed_statuses[action]:
            raise ValueError("lifecycle_job_transition_invalid")
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
                SELECT job.id, job.status
                FROM identity.tenants AS tenant
                JOIN platform_control.lifecycle_jobs AS job ON tenant.id = job.tenant_id
                WHERE tenant.slug = %s
                  AND job.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                raise ValueError("lifecycle_job_not_found")
            if row["status"] not in allowed_statuses[action]:
                raise ValueError("lifecycle_job_transition_invalid")
            cursor.execute(
                """
                UPDATE platform_control.lifecycle_jobs AS job
                SET status = %s,
                    started_at = CASE WHEN %s = 'running' THEN NOW() ELSE job.started_at END,
                    completed_at = CASE WHEN %s = 'completed' THEN NOW() ELSE job.completed_at END,
                    failed_at = CASE WHEN %s = 'failed' THEN NOW() ELSE job.failed_at END,
                    cancelled_at = CASE WHEN %s = 'cancelled' THEN NOW() ELSE job.cancelled_at END,
                    failure_reason = CASE WHEN %s = 'failed' THEN %s ELSE job.failure_reason END
                WHERE job.id = %s
                RETURNING job.id
                """,
                (next_status, next_status, next_status, next_status, next_status, next_status, failure_reason, row["id"]),
            )
            updated = cursor.fetchone()
            if updated is None:
                raise ValueError("lifecycle_job_not_found")
            cursor.execute(
                """
                INSERT INTO platform_control.lifecycle_job_events (job_id, public_id, status, summary)
                VALUES (%s, %s, %s, %s)
                """,
                (updated["id"], str(uuid.uuid4()), next_status, summary),
            )
            connection.commit()

    job = get_lifecycle_job(slug, public_id)
    if job is None:
        raise ValueError("lifecycle_job_not_found")
    return job


def build_go_live_readiness(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    lifecycle = build_lifecycle_readiness(slug)
    usage = build_usage_summary(slug)
    provider_defaults = list_provider_defaults(slug)
    blocks = list_blocks(slug)

    critical_unconfigured = sum(1 for item in provider_defaults if item["critical"] and item["mode"] in {"unconfigured", "disabled"})
    active_blocks = sum(1 for item in blocks if item["active"])
    tracked_metrics = usage["summary"]["trackedMetrics"]
    usage_total = usage["summary"]["totalQuantity"]

    status = "stable"
    risks: list[str] = []
    if critical_unconfigured > 0:
        status = "attention"
        risks.append("critical_providers_unconfigured")
    if active_blocks > 0:
        status = "attention"
        risks.append("tenant_blocks_active")
    if tracked_metrics == 0:
        status = "attention"
        risks.append("adoption_metrics_not_started")

    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "status": status,
        "rolloutReady": status == "stable",
        "metricsObserved": tracked_metrics > 0,
        "rollbackReady": True,
        "lifecycleReady": lifecycle["readiness"]["status"] == "stable",
        "criticalProvidersUnconfigured": critical_unconfigured,
        "activeBlocks": active_blocks,
        "adoptionSignals": {
            "trackedMetrics": tracked_metrics,
            "totalQuantity": usage_total,
            "providerDefaults": len(provider_defaults),
        },
        "risks": risks,
    }


def build_go_live_adoption(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    usage = build_usage_summary(slug)
    rollouts = list_go_live_rollouts(slug, limit=200)["items"]
    latest_rollout = rollouts[0] if rollouts else None

    metrics = usage["metrics"]
    tracked_metrics = usage["summary"]["trackedMetrics"]
    healthy_metrics = sum(1 for item in metrics if item["status"] == "ok")
    attention_metrics = sum(1 for item in metrics if item["status"] == "attention")
    limited_metrics = sum(1 for item in metrics if item["status"] == "limit_reached")
    observed_usage = usage["summary"]["totalQuantity"]

    if tracked_metrics == 0:
        adoption_score = 0
    else:
        baseline = healthy_metrics / tracked_metrics
        pressure_penalty = ((attention_metrics * 0.15) + (limited_metrics * 0.35)) / tracked_metrics
        volume_bonus = 0.15 if observed_usage > 0 else 0.0
        adoption_score = max(0.0, min(1.0, baseline - pressure_penalty + volume_bonus))

    target_pct = int((latest_rollout or {}).get("adoptionTargetPct") or 70)
    adoption_pct = int(round(adoption_score * 100))
    gap_pct = max(target_pct - adoption_pct, 0)

    status = "tracking"
    if tracked_metrics == 0:
        status = "not_started"
    elif adoption_pct >= target_pct:
        status = "on_track"
    elif gap_pct <= 15:
        status = "attention"
    else:
        status = "off_track"

    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "status": status,
        "adoptionPct": adoption_pct,
        "targetPct": target_pct,
        "gapPct": gap_pct,
        "latestRollout": None
        if latest_rollout is None
        else {
            "publicId": latest_rollout["publicId"],
            "waveKey": latest_rollout["waveKey"],
            "status": latest_rollout["status"],
            "targetEnv": latest_rollout["targetEnv"],
        },
        "signals": {
            "trackedMetrics": tracked_metrics,
            "healthyMetrics": healthy_metrics,
            "attentionMetrics": attention_metrics,
            "limitReachedMetrics": limited_metrics,
            "observedUsage": observed_usage,
        },
    }


def build_go_live_bottlenecks(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    readiness = build_go_live_readiness(slug)
    lifecycle = build_lifecycle_readiness(slug)
    usage = build_usage_summary(slug)

    bottlenecks: list[dict] = []

    for capability_key in lifecycle["readiness"]["missingCriticalProviders"]:
        bottlenecks.append(
            {
                "category": "provider",
                "severity": "critical",
                "code": "critical_provider_missing",
                "summary": f"Critical provider for {capability_key} is not configured.",
            }
        )

    for block in lifecycle["blocks"]["items"]:
        bottlenecks.append(
            {
                "category": "tenant_block",
                "severity": "critical" if block["scope"] == "tenant" else "attention",
                "code": "tenant_block_active",
                "summary": f"Active block {block['blockKey']} is limiting rollout progression.",
                "metadata": {"scope": block["scope"], "reason": block["reason"]},
            }
        )

    for metric in usage["metrics"]:
        if metric["status"] == "limit_reached":
            bottlenecks.append(
                {
                    "category": "quota",
                    "severity": "critical",
                    "code": "quota_limit_reached",
                    "summary": f"Quota limit reached for {metric['metricKey']}.",
                    "metadata": {"metricKey": metric["metricKey"], "limitValue": metric["limitValue"], "quantity": metric["quantity"]},
                }
            )
        elif metric["status"] == "attention":
            bottlenecks.append(
                {
                    "category": "quota",
                    "severity": "attention",
                    "code": "quota_near_limit",
                    "summary": f"Quota near limit for {metric['metricKey']}.",
                    "metadata": {"metricKey": metric["metricKey"], "limitValue": metric["limitValue"], "quantity": metric["quantity"]},
                }
            )

    status = "clear"
    if any(item["severity"] == "critical" for item in bottlenecks):
        status = "blocked"
    elif bottlenecks:
        status = "attention"

    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "status": status,
        "summary": {
            "total": len(bottlenecks),
            "critical": sum(1 for item in bottlenecks if item["severity"] == "critical"),
            "attention": sum(1 for item in bottlenecks if item["severity"] == "attention"),
        },
        "items": bottlenecks,
    }


def build_go_live_playbook(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    readiness = build_go_live_readiness(slug)
    adoption = build_go_live_adoption(slug)
    bottlenecks = build_go_live_bottlenecks(slug)
    rollouts = list_go_live_rollouts(slug, limit=20)["items"]
    latest_rollout = rollouts[0] if rollouts else None

    checklist = [
        {
            "key": "validate-readiness",
            "status": "ready" if readiness["rolloutReady"] else "attention",
            "summary": "Validate providers, quotas and tenant blocks before the wave.",
        },
        {
            "key": "observe-adoption",
            "status": "ready" if adoption["status"] in {"on_track", "tracking"} else "attention",
            "summary": "Observe adoption signals, tracked metrics and usage growth after release.",
        },
        {
            "key": "remove-bottlenecks",
            "status": "ready" if bottlenecks["summary"]["critical"] == 0 else "blocked",
            "summary": "Resolve critical provider gaps, hard blocks and quota incidents before broad rollout.",
        },
        {
            "key": "rollback",
            "status": "ready",
            "summary": (latest_rollout or {}).get("rollbackPlaybook") or "docs/OPERACOES.md#rollback",
        },
    ]

    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "status": "ready" if readiness["rolloutReady"] and bottlenecks["summary"]["critical"] == 0 else "attention",
        "latestRollout": latest_rollout,
        "readiness": readiness,
        "adoption": adoption,
        "bottlenecks": bottlenecks,
        "checklist": checklist,
    }


def build_go_live_adjustments(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    usage = build_usage_summary(slug)
    lifecycle = build_lifecycle_readiness(slug)
    adoption = build_go_live_adoption(slug)

    quota_map = {item["metricKey"]: item for item in list_quotas(slug)}
    recommendations: list[dict] = []

    for metric in usage["metrics"]:
        if metric["status"] in {"attention", "limit_reached"}:
            current_limit = int(metric["limitValue"] or 0)
            suggested_limit = current_limit
            if current_limit > 0:
                suggested_limit = max(current_limit + 1, int(round(current_limit * 1.25)))
            recommendations.append(
                {
                    "actionType": "increase_quota_limit",
                    "severity": "critical" if metric["status"] == "limit_reached" else "attention",
                    "summary": f"Increase quota for {metric['metricKey']} to reduce rollout pressure.",
                    "applySupported": current_limit > 0,
                    "payload": {
                        "metricKey": metric["metricKey"],
                        "metricUnit": metric["metricUnit"] or quota_map.get(metric["metricKey"], {}).get("metricUnit", ""),
                        "newLimitValue": suggested_limit,
                        "enforcementMode": metric["enforcementMode"],
                    },
                }
            )

    for block in lifecycle["blocks"]["items"]:
        recommendations.append(
            {
                "actionType": "disable_block",
                "severity": "critical" if block["scope"] == "tenant" else "attention",
                "summary": f"Disable active block {block['blockKey']} before widening rollout.",
                "applySupported": True,
                "payload": {
                    "blockKey": block["blockKey"],
                    "reason": block["reason"],
                    "scope": block["scope"],
                },
            }
        )

    provider_defaults = list_provider_defaults(slug)
    for item in provider_defaults:
        if item["critical"] and item["mode"] in {"unconfigured", "disabled"} and item["fallbackAllowed"]:
            recommendations.append(
                {
                    "actionType": "set_provider_mode",
                    "severity": "attention",
                    "summary": f"Enable safe fallback for {item['capabilityKey']} while provider credentials are unavailable.",
                    "applySupported": True,
                    "payload": {
                        "capabilityKey": item["capabilityKey"],
                        "providerKey": item["providerKey"],
                        "mode": "fallback",
                    },
                }
            )

    if adoption["status"] in {"not_started", "off_track"}:
        recommendations.append(
            {
                "actionType": "review_rollout_wave",
                "severity": "attention",
                "summary": "Review rollout wave sizing and adoption target before progressing.",
                "applySupported": False,
                "payload": {
                    "targetPct": adoption["targetPct"],
                    "adoptionPct": adoption["adoptionPct"],
                    "gapPct": adoption["gapPct"],
                },
            }
        )

    return {
        "tenantSlug": slug,
        "generatedAt": utc_now(),
        "summary": {
            "recommended": len(recommendations),
            "applySupported": sum(1 for item in recommendations if item["applySupported"]),
            "critical": sum(1 for item in recommendations if item["severity"] == "critical"),
            "attention": sum(1 for item in recommendations if item["severity"] == "attention"),
        },
        "items": recommendations,
    }


def apply_go_live_adjustment(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    action_type = str(payload.get("actionType") or "").strip()
    actor = str(payload.get("actor") or "ops@erp.local").strip() or "ops@erp.local"

    if action_type == "increase_quota_limit":
        metric_key = str(payload.get("metricKey") or "").strip()
        metric_unit = str(payload.get("metricUnit") or "").strip()
        new_limit_value = int(payload.get("newLimitValue") or 0)
        if metric_key == "":
            raise ValueError("metric_key_required")
        if metric_unit == "":
            raise ValueError("metric_unit_required")
        if new_limit_value <= 0:
            raise ValueError("new_limit_value_required")
        updated = upsert_quota(
            slug,
            metric_key,
            {
                "metricUnit": metric_unit,
                "limitValue": new_limit_value,
                "enforcementMode": str(payload.get("enforcementMode") or "soft"),
                "source": "go-live-adjustment",
            },
        )
        return {
            "tenantSlug": slug,
            "actionType": action_type,
            "actor": actor,
            "appliedAt": utc_now(),
            "result": updated,
        }

    if action_type == "disable_block":
        block_key = str(payload.get("blockKey") or "").strip()
        if block_key == "":
            raise ValueError("block_key_required")
        updated = upsert_block(
            slug,
            block_key,
            {
                "active": False,
                "reason": str(payload.get("reason") or "released by go-live adjustment"),
                "scope": str(payload.get("scope") or "tenant"),
                "source": "go-live-adjustment",
            },
        )
        return {
            "tenantSlug": slug,
            "actionType": action_type,
            "actor": actor,
            "appliedAt": utc_now(),
            "result": updated,
        }

    if action_type == "set_provider_mode":
        capability_key = str(payload.get("capabilityKey") or "").strip()
        provider_key = str(payload.get("providerKey") or "").strip()
        mode = str(payload.get("mode") or "").strip()
        if capability_key == "":
            raise ValueError("capability_key_required")
        if provider_key == "":
            raise ValueError("provider_key_required")
        if mode == "":
            raise ValueError("provider_mode_invalid")
        updated = upsert_provider_default(
            slug,
            capability_key,
            {
                "providerKey": provider_key,
                "mode": mode,
                "source": "go-live-adjustment",
            },
        )
        return {
            "tenantSlug": slug,
            "actionType": action_type,
            "actor": actor,
            "appliedAt": utc_now(),
            "result": updated,
        }

    raise ValueError("go_live_adjustment_action_invalid")


def create_go_live_rollout(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    target_env = (payload.get("targetEnv") or "production").strip().lower()
    wave_key = (payload.get("waveKey") or "wave-1").strip().lower()
    requested_by = (payload.get("requestedBy") or "").strip()
    rollback_playbook = (payload.get("rollbackPlaybook") or "docs/OPERACOES.md#rollback").strip()
    adoption_target_pct = int(payload.get("adoptionTargetPct") or 70)
    if requested_by == "":
        raise ValueError("rollout_requested_by_required")

    readiness = build_go_live_readiness(slug)
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "targetEnv": target_env,
        "waveKey": wave_key,
        "requestedBy": requested_by,
        "status": "planned",
        "rollbackPlaybook": rollback_playbook,
        "adoptionTargetPct": adoption_target_pct,
        "createdAt": utc_now(),
        "startedAt": None,
        "completedAt": None,
        "rolledBackAt": None,
        "readinessSnapshot": readiness,
        "events": [],
    }

    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["go_live_rollouts"].append(record)
        IN_MEMORY_STATE["go_live_events"].append(
            {
                "publicId": str(uuid.uuid4()),
                "rolloutPublicId": record["publicId"],
                "tenantSlug": slug,
                "status": "planned",
                "summary": "Rollout planned.",
                "createdAt": utc_now(),
            }
        )
        return get_go_live_rollout(slug, record["publicId"]) or record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO platform_control.go_live_rollouts (
                  tenant_id, public_id, target_env, wave_key, status, requested_by, rollback_playbook,
                  adoption_target_pct, readiness_json
                )
                VALUES (%s, %s, %s, %s, 'planned', %s, %s, %s, %s::jsonb)
                RETURNING created_at
                """,
                (tenant_id, record["publicId"], target_env, wave_key, requested_by, rollback_playbook, adoption_target_pct, json.dumps(readiness)),
            )
            row = cursor.fetchone()
            cursor.execute(
                """
                INSERT INTO platform_control.go_live_rollout_events (rollout_id, public_id, status, summary)
                SELECT id, %s, 'planned', 'Rollout planned.'
                FROM platform_control.go_live_rollouts
                WHERE public_id = %s
                """,
                (str(uuid.uuid4()), record["publicId"]),
            )
            connection.commit()
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return get_go_live_rollout(slug, record["publicId"]) or record


def list_go_live_rollouts(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        records = sorted(
            [item for item in IN_MEMORY_STATE["go_live_rollouts"] if item["tenantSlug"] == slug],
            key=lambda item: item["createdAt"],
            reverse=True,
        )
        payload = _paginate(records, cursor, limit)
        payload["tenantSlug"] = slug
        return payload

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                SELECT rollout.public_id
                FROM platform_control.go_live_rollouts AS rollout
                JOIN identity.tenants AS tenant ON tenant.id = rollout.tenant_id
                WHERE tenant.slug = %s
                ORDER BY rollout.created_at DESC
                """,
                (slug,),
            )
            records = [get_go_live_rollout(slug, row["public_id"]) for row in cursor_db.fetchall()]
            payload = _paginate([record for record in records if record is not None], cursor, limit)
            payload["tenantSlug"] = slug
            return payload


def get_go_live_rollout(tenant_slug: str, public_id: str) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        for rollout in IN_MEMORY_STATE["go_live_rollouts"]:
            if rollout["tenantSlug"] == slug and rollout["publicId"] == public_id:
                rollout["events"] = [
                    event for event in IN_MEMORY_STATE["go_live_events"] if event["tenantSlug"] == slug and event["rolloutPublicId"] == public_id
                ]
                return rollout
        return None

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                SELECT rollout.public_id,
                       rollout.target_env,
                       rollout.wave_key,
                       rollout.status,
                       rollout.requested_by,
                       rollout.rollback_playbook,
                       rollout.adoption_target_pct,
                       rollout.readiness_json,
                       rollout.created_at,
                       rollout.started_at,
                       rollout.completed_at,
                       rollout.rolled_back_at
                FROM platform_control.go_live_rollouts AS rollout
                JOIN identity.tenants AS tenant ON tenant.id = rollout.tenant_id
                WHERE tenant.slug = %s AND rollout.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor_db.fetchone()
            if row is None:
                return None
            cursor_db.execute(
                """
                SELECT event.public_id, event.status, event.summary, event.created_at
                FROM platform_control.go_live_rollout_events AS event
                JOIN platform_control.go_live_rollouts AS rollout ON rollout.id = event.rollout_id
                WHERE rollout.public_id = %s
                ORDER BY event.created_at
                """,
                (public_id,),
            )
            events = [
                {
                    "publicId": event["public_id"],
                    "status": event["status"],
                    "summary": event["summary"],
                    "createdAt": event["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for event in cursor_db.fetchall()
            ]
            return {
                "publicId": row["public_id"],
                "tenantSlug": slug,
                "targetEnv": row["target_env"],
                "waveKey": row["wave_key"],
                "status": row["status"],
                "requestedBy": row["requested_by"],
                "rollbackPlaybook": row["rollback_playbook"],
                "adoptionTargetPct": int(row["adoption_target_pct"]),
                "readinessSnapshot": row["readiness_json"] or {},
                "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "startedAt": None if row["started_at"] is None else row["started_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "completedAt": None if row["completed_at"] is None else row["completed_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "rolledBackAt": None if row["rolled_back_at"] is None else row["rolled_back_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "events": events,
            }


def transition_go_live_rollout(tenant_slug: str, public_id: str, action: str, payload: dict | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    payload = payload or {}
    transition_map = {
        "start": ("running", "Rollout started."),
        "complete": ("completed", "Rollout completed."),
        "rollback": ("rolled_back", "Rollback executed."),
    }
    allowed_statuses = {
        "start": {"planned"},
        "complete": {"running"},
        "rollback": {"running", "completed"},
    }
    if action not in transition_map:
        raise ValueError("rollout_action_invalid")
    next_status, fallback_summary = transition_map[action]
    summary = (payload.get("summary") or fallback_summary).strip()

    if settings.repository_driver != "postgres":
        for rollout in IN_MEMORY_STATE["go_live_rollouts"]:
            if rollout["tenantSlug"] == slug and rollout["publicId"] == public_id:
                if rollout["status"] not in allowed_statuses[action]:
                    raise ValueError("go_live_rollout_transition_invalid")
                rollout["status"] = next_status
                timestamp = utc_now()
                if action == "start":
                    rollout["startedAt"] = timestamp
                elif action == "complete":
                    rollout["completedAt"] = timestamp
                elif action == "rollback":
                    rollout["rolledBackAt"] = timestamp
                IN_MEMORY_STATE["go_live_events"].append(
                    {
                        "publicId": str(uuid.uuid4()),
                        "rolloutPublicId": public_id,
                        "tenantSlug": slug,
                        "status": next_status,
                        "summary": summary,
                        "createdAt": timestamp,
                    }
                )
                return get_go_live_rollout(slug, public_id) or rollout
        raise ValueError("go_live_rollout_not_found")

    with connect() as connection:
        with connection.cursor() as cursor_db:
            cursor_db.execute(
                """
                SELECT rollout.id, rollout.status
                FROM identity.tenants AS tenant
                JOIN platform_control.go_live_rollouts AS rollout ON tenant.id = rollout.tenant_id
                WHERE tenant.slug = %s
                  AND rollout.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor_db.fetchone()
            if row is None:
                raise ValueError("go_live_rollout_not_found")
            if row["status"] not in allowed_statuses[action]:
                raise ValueError("go_live_rollout_transition_invalid")
            update_column = {
                "start": "started_at",
                "complete": "completed_at",
                "rollback": "rolled_back_at",
            }[action]
            cursor_db.execute(
                f"""
                UPDATE platform_control.go_live_rollouts
                SET status = %s, {update_column} = NOW()
                WHERE id = %s
                """,
                (next_status, row["id"]),
            )
            cursor_db.execute(
                """
                INSERT INTO platform_control.go_live_rollout_events (rollout_id, public_id, status, summary)
                VALUES (%s, %s, %s, %s)
                """,
                (row["id"], str(uuid.uuid4()), next_status, summary),
            )
            connection.commit()
            result = get_go_live_rollout(slug, public_id)
            if result is None:
                raise ValueError("go_live_rollout_not_found")
            return result


def list_incidents(tenant_slug: str | None = None, status: str | None = None, severity: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug or settings.bootstrap_tenant_slug)
    normalized_status = (status or "").strip().lower()
    normalized_severity = (severity or "").strip().lower()
    records = [item for item in IN_MEMORY_STATE["incidents"] if item["tenantSlug"] == slug]
    if normalized_status:
        records = [item for item in records if item["status"] == normalized_status]
    if normalized_severity:
        records = [item for item in records if item["severity"] == normalized_severity]
    records = sorted(records, key=lambda item: item["createdAt"], reverse=True)
    return {
        "tenantSlug": slug,
        "summary": {
            "total": len(records),
            "open": sum(1 for item in records if item["status"] in {"open", "investigating", "mitigating"}),
            "critical": sum(1 for item in records if item["severity"] == "sev1"),
            "resolved": sum(1 for item in records if item["status"] == "resolved"),
        },
        "items": records,
    }


def create_incident(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    title = str(payload.get("title") or "").strip()
    service = str(payload.get("service") or "").strip()
    severity = str(payload.get("severity") or "sev3").strip().lower()
    if not title:
        raise ValueError("incident_title_required")
    if not service:
        raise ValueError("incident_service_required")
    if severity not in {"sev1", "sev2", "sev3", "sev4"}:
        raise ValueError("incident_severity_invalid")
    incident = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "title": title,
        "service": service,
        "severity": severity,
        "status": "open",
        "impact": str(payload.get("impact") or "under investigation"),
        "owner": str(payload.get("owner") or "ops@erp.local"),
        "startedAt": payload.get("startedAt") or utc_now(),
        "resolvedAt": None,
        "createdAt": utc_now(),
        "timeline": [],
        "actions": [],
        "postmortem": None,
    }
    IN_MEMORY_STATE["incidents"].append(incident)
    append_incident_timeline(slug, incident["publicId"], {"eventType": "created", "summary": "Incident opened.", "actor": incident["owner"]})
    return get_incident(slug, incident["publicId"]) or incident


def get_incident(tenant_slug: str, public_id: str) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    incident = next((item for item in IN_MEMORY_STATE["incidents"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
    if incident is None:
        return None
    incident["timeline"] = [item for item in IN_MEMORY_STATE["incident_timeline_events"] if item["tenantSlug"] == slug and item["incidentPublicId"] == public_id]
    incident["actions"] = [item for item in IN_MEMORY_STATE["incident_actions"] if item["tenantSlug"] == slug and item["incidentPublicId"] == public_id]
    incident["postmortem"] = next((item for item in IN_MEMORY_STATE["postmortems"] if item["tenantSlug"] == slug and item["incidentPublicId"] == public_id), None)
    return incident


def append_incident_timeline(tenant_slug: str, public_id: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    incident = next((item for item in IN_MEMORY_STATE["incidents"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
    if incident is None:
        raise ValueError("incident_not_found")
    summary = str(payload.get("summary") or "").strip()
    if not summary:
        raise ValueError("incident_timeline_summary_required")
    event = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "incidentPublicId": public_id,
        "eventType": str(payload.get("eventType") or "note"),
        "summary": summary,
        "actor": str(payload.get("actor") or "ops@erp.local"),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["incident_timeline_events"].append(event)
    return event


def create_incident_action(tenant_slug: str, public_id: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    incident = get_incident(slug, public_id)
    if incident is None:
        raise ValueError("incident_not_found")
    title = str(payload.get("title") or "").strip()
    if not title:
        raise ValueError("incident_action_title_required")
    action = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "incidentPublicId": public_id,
        "title": title,
        "status": str(payload.get("status") or "open"),
        "owner": str(payload.get("owner") or incident["owner"]),
        "dueAt": payload.get("dueAt"),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["incident_actions"].append(action)
    append_incident_timeline(slug, public_id, {"eventType": "action_created", "summary": f"Action created: {title}", "actor": action["owner"]})
    return action


def resolve_incident(tenant_slug: str, public_id: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    incident = get_incident(slug, public_id)
    if incident is None:
        raise ValueError("incident_not_found")
    summary = str(payload.get("summary") or "").strip()
    if not summary:
        raise ValueError("incident_resolution_summary_required")
    incident["status"] = "resolved"
    incident["resolvedAt"] = utc_now()
    append_incident_timeline(slug, public_id, {"eventType": "resolved", "summary": summary, "actor": payload.get("actor") or incident["owner"]})
    return get_incident(slug, public_id) or incident


def create_postmortem(tenant_slug: str, public_id: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    incident = get_incident(slug, public_id)
    if incident is None:
        raise ValueError("incident_not_found")
    root_cause = str(payload.get("rootCause") or "").strip()
    if not root_cause:
        raise ValueError("postmortem_root_cause_required")
    postmortem = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "incidentPublicId": public_id,
        "rootCause": root_cause,
        "impactSummary": str(payload.get("impactSummary") or incident["impact"]),
        "preventiveActions": payload.get("preventiveActions") if isinstance(payload.get("preventiveActions"), list) else [],
        "evidence": payload.get("evidence") if isinstance(payload.get("evidence"), list) else [],
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["postmortems"].append(postmortem)
    append_incident_timeline(slug, public_id, {"eventType": "postmortem_created", "summary": "Postmortem created.", "actor": payload.get("actor") or incident["owner"]})
    return postmortem


def build_incident_command_readiness(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug or settings.bootstrap_tenant_slug)
    incidents = list_incidents(slug)["summary"]
    return {
        "acceptanceReady": incidents["critical"] == 0 or incidents["resolved"] >= incidents["critical"],
        "controls": ["incident-registry", "append-only-timeline", "action-tracking", "postmortem"],
        "summary": incidents,
    }


def list_policy_catalog() -> dict:
    return {"items": sorted(POLICY_CATALOG, key=lambda item: (item["domain"], item["priority"], item["policyKey"]))}


def evaluate_policy(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    action = str(payload.get("action") or "").strip()
    domain = str(payload.get("domain") or "").strip()
    actor = str(payload.get("actor") or "").strip()
    if not action:
        raise ValueError("policy_action_required")
    if not domain:
        raise ValueError("policy_domain_required")
    if not actor:
        raise ValueError("policy_actor_required")

    context = payload.get("context") if isinstance(payload.get("context"), dict) else {}
    candidates = [item for item in POLICY_CATALOG if item["action"] == action or item["domain"] == domain]
    selected = sorted(candidates, key=lambda item: item["priority"])[0] if candidates else None
    effect = selected["effect"] if selected else "allow"
    if action.startswith("ai.tool.") and context.get("toolMode") == "write":
        effect = "deny"
    if context.get("severity") == "sev1" or context.get("legalHoldActive") is True:
        effect = "review" if effect != "deny" else "deny"

    decision = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "policyKey": selected["policyKey"] if selected else "default.allow",
        "policyVersion": selected["version"] if selected else "1.0",
        "domain": domain,
        "action": action,
        "effect": effect,
        "decision": {"allow": "allow", "review": "review", "deny": "deny"}[effect],
        "actor": actor,
        "reason": selected["reason"] if selected else "No restrictive policy matched.",
        "context": context,
        "evaluatedAt": utc_now(),
    }
    IN_MEMORY_STATE["policy_decisions"].append(decision)
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "policy.decision",
            "entityPublicId": decision["publicId"],
            "eventType": f"policy.{effect}",
            "actor": actor,
            "summary": f"Policy {decision['policyKey']} returned {effect} for {action}.",
            "severity": "warning" if effect == "review" else ("critical" if effect == "deny" else "info"),
            "metadata": {"domain": domain, "action": action},
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "policy-decision",
            "entityType": "policy.decision",
            "entityPublicId": decision["publicId"],
            "actor": actor,
            "classification": "internal",
            "payload": decision,
            "retention": "p2y",
        },
    )
    return decision


def list_policy_decisions(tenant_slug: str, action: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["policy_decisions"] if item["tenantSlug"] == slug]
    if action:
        records = [item for item in records if item["action"] == action]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["evaluatedAt"], reverse=True)}


def record_timeline_event(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    source_service = str(payload.get("sourceService") or "").strip()
    entity_type = str(payload.get("entityType") or "").strip()
    entity_public_id = str(payload.get("entityPublicId") or "").strip()
    event_type = str(payload.get("eventType") or "").strip()
    actor = str(payload.get("actor") or "system:platform-control").strip()
    summary = str(payload.get("summary") or "").strip()
    if not source_service:
        raise ValueError("timeline_source_service_required")
    if not entity_type:
        raise ValueError("timeline_entity_type_required")
    if not entity_public_id:
        raise ValueError("timeline_entity_public_id_required")
    if not event_type:
        raise ValueError("timeline_event_type_required")
    if not summary:
        raise ValueError("timeline_summary_required")
    event = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "sourceService": source_service,
        "entityType": entity_type,
        "entityPublicId": entity_public_id,
        "eventType": event_type,
        "severity": str(payload.get("severity") or "info"),
        "actor": actor,
        "correlationId": str(payload.get("correlationId") or uuid.uuid4()),
        "summary": summary,
        "metadata": payload.get("metadata") if isinstance(payload.get("metadata"), dict) else {},
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["timeline_events"].append(event)
    return event


def list_timeline_events(tenant_slug: str, entity_type: str | None = None, entity_public_id: str | None = None, limit: int = 50) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["timeline_events"] if item["tenantSlug"] == slug]
    if entity_type:
        records = [item for item in records if item["entityType"] == entity_type]
    if entity_public_id:
        records = [item for item in records if item["entityPublicId"] == entity_public_id]
    records = sorted(records, key=lambda item: item["createdAt"], reverse=True)[: _normalize_limit(limit)]
    return {"tenantSlug": slug, "items": records, "summary": {"total": len(records)}}


def create_approval_request(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    command_type = str(payload.get("commandType") or "").strip()
    requested_by = str(payload.get("requestedBy") or "").strip()
    justification = str(payload.get("justification") or "").strip()
    if not command_type:
        raise ValueError("approval_command_type_required")
    if not requested_by:
        raise ValueError("approval_requested_by_required")
    if not justification:
        raise ValueError("approval_justification_required")
    policy_decision = evaluate_policy(
        slug,
        {
            "domain": str(payload.get("domain") or "platform-control"),
            "action": command_type,
            "actor": requested_by,
            "context": payload.get("commandPayload") if isinstance(payload.get("commandPayload"), dict) else {},
        },
    )
    approval = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "commandType": command_type,
        "domain": str(payload.get("domain") or "platform-control"),
        "status": "requested" if policy_decision["decision"] != "deny" else "rejected",
        "requestedBy": requested_by,
        "approvedBy": None,
        "rejectedBy": None,
        "executedBy": None,
        "justification": justification,
        "policyDecision": policy_decision,
        "commandPayload": payload.get("commandPayload") if isinstance(payload.get("commandPayload"), dict) else {},
        "createdAt": utc_now(),
        "decidedAt": None,
        "executedAt": None,
    }
    IN_MEMORY_STATE["approval_requests"].append(approval)
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "approval.request",
            "entityPublicId": approval["publicId"],
            "eventType": f"approval.{approval['status']}",
            "actor": requested_by,
            "summary": f"Approval requested for {command_type}.",
            "severity": "warning" if approval["status"] == "requested" else "critical",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "approval-request",
            "entityType": "approval.request",
            "entityPublicId": approval["publicId"],
            "actor": requested_by,
            "classification": "internal",
            "payload": approval,
            "retention": "p2y",
        },
    )
    return approval


def list_approval_requests(tenant_slug: str, status: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["approval_requests"] if item["tenantSlug"] == slug]
    if status:
        records = [item for item in records if item["status"] == status]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["createdAt"], reverse=True)}


def transition_approval_request(tenant_slug: str, public_id: str, action: str, payload: dict | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    payload = payload or {}
    approval = next((item for item in IN_MEMORY_STATE["approval_requests"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
    if approval is None:
        raise ValueError("approval_request_not_found")
    actor = str(payload.get("actor") or "").strip()
    if not actor:
        raise ValueError("approval_actor_required")
    now = utc_now()
    if action == "approve":
        if approval["status"] != "requested":
            raise ValueError("approval_transition_invalid")
        approval["status"] = "approved"
        approval["approvedBy"] = actor
        approval["decidedAt"] = now
    elif action == "reject":
        if approval["status"] != "requested":
            raise ValueError("approval_transition_invalid")
        approval["status"] = "rejected"
        approval["rejectedBy"] = actor
        approval["decidedAt"] = now
    elif action == "execute":
        if approval["status"] != "approved":
            raise ValueError("approval_transition_invalid")
        approval["status"] = "executed"
        approval["executedBy"] = actor
        approval["executedAt"] = now
    elif action == "cancel":
        if approval["status"] not in {"requested", "approved"}:
            raise ValueError("approval_transition_invalid")
        approval["status"] = "cancelled"
        approval["decidedAt"] = now
    else:
        raise ValueError("approval_action_invalid")
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "approval.request",
            "entityPublicId": public_id,
            "eventType": f"approval.{approval['status']}",
            "actor": actor,
            "summary": f"Approval {action} executed.",
            "severity": "info",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": f"approval-{approval['status']}",
            "entityType": "approval.request",
            "entityPublicId": public_id,
            "actor": actor,
            "classification": "internal",
            "payload": approval,
            "retention": "p2y",
        },
    )
    return approval


def list_runbook_catalog() -> dict:
    return {"items": RUNBOOK_CATALOG}


def create_runbook_run(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    runbook_key = str(payload.get("runbookKey") or "").strip()
    requested_by = str(payload.get("requestedBy") or "").strip()
    if not runbook_key:
        raise ValueError("runbook_key_required")
    if not requested_by:
        raise ValueError("runbook_requested_by_required")
    template = next((item for item in RUNBOOK_CATALOG if item["runbookKey"] == runbook_key), None)
    if template is None:
        raise ValueError("runbook_not_found")
    approval = create_approval_request(
        slug,
        {
            "commandType": f"runbook.{runbook_key}",
            "domain": template["domain"],
            "requestedBy": requested_by,
            "justification": str(payload.get("justification") or f"Runbook {runbook_key} execution."),
            "commandPayload": {"runbookKey": runbook_key, "severity": payload.get("severity")},
        },
    )
    status = "waiting_approval" if approval["status"] == "requested" else "failed"
    run = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "runbookKey": runbook_key,
        "title": template["title"],
        "domain": template["domain"],
        "status": status,
        "requestedBy": requested_by,
        "approvalPublicId": approval["publicId"],
        "context": payload.get("context") if isinstance(payload.get("context"), dict) else {},
        "createdAt": utc_now(),
        "startedAt": None,
        "completedAt": None,
        "steps": [],
    }
    for index, step_key in enumerate(template["steps"], start=1):
        run["steps"].append(
            {
                "publicId": str(uuid.uuid4()),
                "tenantSlug": slug,
                "runPublicId": run["publicId"],
                "sequence": index,
                "stepKey": step_key,
                "status": "pending",
                "summary": step_key.replace("_", " "),
            }
        )
    IN_MEMORY_STATE["runbook_runs"].append(run)
    IN_MEMORY_STATE["runbook_steps"].extend(run["steps"])
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "runbook.run",
            "entityPublicId": run["publicId"],
            "eventType": "runbook.created",
            "actor": requested_by,
            "summary": f"Runbook {runbook_key} created.",
            "severity": "warning",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "runbook-run",
            "entityType": "runbook.run",
            "entityPublicId": run["publicId"],
            "actor": requested_by,
            "classification": "internal",
            "payload": run,
            "retention": "p2y",
        },
    )
    return run


def list_runbook_runs(tenant_slug: str, status: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["runbook_runs"] if item["tenantSlug"] == slug]
    if status:
        records = [item for item in records if item["status"] == status]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["createdAt"], reverse=True)}


def transition_runbook_run(tenant_slug: str, public_id: str, action: str, payload: dict | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    payload = payload or {}
    run = next((item for item in IN_MEMORY_STATE["runbook_runs"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
    if run is None:
        raise ValueError("runbook_run_not_found")
    actor = str(payload.get("actor") or run["requestedBy"])
    now = utc_now()
    if action == "start":
        if run["status"] not in {"waiting_approval", "planned"}:
            raise ValueError("runbook_transition_invalid")
        approval = next((item for item in IN_MEMORY_STATE["approval_requests"] if item["publicId"] == run["approvalPublicId"]), None)
        if approval and approval["status"] not in {"approved", "executed"}:
            raise ValueError("runbook_approval_required")
        run["status"] = "running"
        run["startedAt"] = now
        if run["steps"]:
            run["steps"][0]["status"] = "running"
    elif action == "complete-step":
        if run["status"] != "running":
            raise ValueError("runbook_transition_invalid")
        pending = next((step for step in run["steps"] if step["status"] in {"pending", "running"}), None)
        if pending is None:
            raise ValueError("runbook_step_not_found")
        pending["status"] = "completed"
        next_step = next((step for step in run["steps"] if step["status"] == "pending"), None)
        if next_step:
            next_step["status"] = "running"
        else:
            run["status"] = "completed"
            run["completedAt"] = now
    elif action == "cancel":
        if run["status"] in {"completed", "failed"}:
            raise ValueError("runbook_transition_invalid")
        run["status"] = "cancelled"
        run["completedAt"] = now
    else:
        raise ValueError("runbook_action_invalid")
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "runbook.run",
            "entityPublicId": public_id,
            "eventType": f"runbook.{run['status']}",
            "actor": actor,
            "summary": f"Runbook action {action} executed.",
            "severity": "info",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": f"runbook-{run['status']}",
            "entityType": "runbook.run",
            "entityPublicId": public_id,
            "actor": actor,
            "classification": "internal",
            "payload": run,
            "retention": "p2y",
        },
    )
    return run


def register_evidence(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    source_service = str(payload.get("sourceService") or "").strip()
    evidence_type = str(payload.get("evidenceType") or "").strip()
    entity_type = str(payload.get("entityType") or "").strip()
    entity_public_id = str(payload.get("entityPublicId") or "").strip()
    actor = str(payload.get("actor") or "system:platform-control").strip()
    if not source_service:
        raise ValueError("evidence_source_service_required")
    if not evidence_type:
        raise ValueError("evidence_type_required")
    if not entity_type:
        raise ValueError("evidence_entity_type_required")
    if not entity_public_id:
        raise ValueError("evidence_entity_public_id_required")
    evidence_payload = payload.get("payload") if isinstance(payload.get("payload"), dict) else {}
    canonical = json.dumps(evidence_payload, sort_keys=True, separators=(",", ":"))
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "sourceService": source_service,
        "evidenceType": evidence_type,
        "entityType": entity_type,
        "entityPublicId": entity_public_id,
        "actor": actor,
        "classification": str(payload.get("classification") or "internal"),
        "retention": str(payload.get("retention") or "p2y"),
        "payload": evidence_payload,
        "payloadHash": hashlib.sha256(canonical.encode("utf-8")).hexdigest(),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["evidence_records"].append(record)
    return record


def list_evidence_records(tenant_slug: str, evidence_type: str | None = None, entity_type: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["evidence_records"] if item["tenantSlug"] == slug]
    if evidence_type:
        records = [item for item in records if item["evidenceType"] == evidence_type]
    if entity_type:
        records = [item for item in records if item["entityType"] == entity_type]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["createdAt"], reverse=True)}


def get_evidence_record(tenant_slug: str, public_id: str) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    return next((item for item in IN_MEMORY_STATE["evidence_records"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)


def build_autonomous_governance_readiness(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug or settings.bootstrap_tenant_slug)
    return {
        "acceptanceReady": True,
        "controls": [
            "policy-decision-center",
            "operational-timeline",
            "command-approvals",
            "runbook-automation",
            "audit-evidence-vault",
        ],
        "summary": {
            "policies": len(POLICY_CATALOG),
            "policyDecisions": len([item for item in IN_MEMORY_STATE["policy_decisions"] if item["tenantSlug"] == slug]),
            "timelineEvents": len([item for item in IN_MEMORY_STATE["timeline_events"] if item["tenantSlug"] == slug]),
            "approvals": len([item for item in IN_MEMORY_STATE["approval_requests"] if item["tenantSlug"] == slug]),
            "runbooks": len(RUNBOOK_CATALOG),
            "evidenceRecords": len([item for item in IN_MEMORY_STATE["evidence_records"] if item["tenantSlug"] == slug]),
        },
    }


def _payload_hash(payload: dict) -> str:
    canonical = json.dumps(payload, sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(canonical.encode("utf-8")).hexdigest()


def list_event_mesh_catalog() -> dict:
    return {
        "version": "1.2.0",
        "streams": EVENT_STREAM_CATALOG,
        "summary": {
            "streams": len(EVENT_STREAM_CATALOG),
            "critical": len([item for item in EVENT_STREAM_CATALOG if item["critical"]]),
            "retentionPolicies": sorted({item["retention"] for item in EVENT_STREAM_CATALOG}),
        },
    }


def record_event_mesh_event(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    stream_key = str(payload.get("streamKey") or "").strip()
    event_type = str(payload.get("eventType") or "").strip()
    producer = str(payload.get("producer") or "").strip()
    if not stream_key:
        raise ValueError("event_stream_key_required")
    if not event_type:
        raise ValueError("event_type_required")
    if not producer:
        raise ValueError("event_producer_required")
    stream = next((item for item in EVENT_STREAM_CATALOG if item["streamKey"] == stream_key), None)
    if stream is None:
        raise ValueError("event_stream_not_found")
    event_payload = payload.get("payload") if isinstance(payload.get("payload"), dict) else {}
    status = str(payload.get("status") or "published")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "streamKey": stream_key,
        "eventType": event_type,
        "schemaVersion": str(payload.get("schemaVersion") or stream["schemaVersion"]),
        "producer": producer,
        "consumer": payload.get("consumer"),
        "correlationId": str(payload.get("correlationId") or uuid.uuid4()),
        "causationId": payload.get("causationId"),
        "status": status,
        "payload": event_payload,
        "payloadHash": _payload_hash(event_payload),
        "occurredAt": str(payload.get("occurredAt") or utc_now()),
    }
    IN_MEMORY_STATE["event_mesh_events"].append(record)
    if status in {"failed", "dead_letter"}:
        dead_letter = {
            "publicId": str(uuid.uuid4()),
            "tenantSlug": slug,
            "eventPublicId": record["publicId"],
            "streamKey": stream_key,
            "eventType": event_type,
            "producer": producer,
            "reason": str(payload.get("reason") or "consumer_failure"),
            "status": "waiting_replay",
            "attempts": int(payload.get("attempts") or 1),
            "payloadHash": record["payloadHash"],
            "createdAt": utc_now(),
            "replayedAt": None,
        }
        IN_MEMORY_STATE["event_mesh_dead_letters"].append(dead_letter)
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "event.mesh",
            "entityPublicId": record["publicId"],
            "eventType": f"event.{status}",
            "actor": producer,
            "summary": f"{event_type} published on {stream_key}.",
            "severity": "warning" if status in {"failed", "dead_letter"} else "info",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "event-mesh-event",
            "entityType": "event.mesh",
            "entityPublicId": record["publicId"],
            "actor": producer,
            "classification": "internal",
            "payload": record,
            "retention": stream["retention"],
        },
    )
    return record


def list_event_mesh_events(tenant_slug: str, stream_key: str | None = None, status: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["event_mesh_events"] if item["tenantSlug"] == slug]
    if stream_key:
        records = [item for item in records if item["streamKey"] == stream_key]
    if status:
        records = [item for item in records if item["status"] == status]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["occurredAt"], reverse=True)}


def list_event_mesh_dead_letters(tenant_slug: str, status: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["event_mesh_dead_letters"] if item["tenantSlug"] == slug]
    if status:
        records = [item for item in records if item["status"] == status]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["createdAt"], reverse=True)}


def replay_event_mesh_dead_letter(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    payload = payload or {}
    dead_letter = next((item for item in IN_MEMORY_STATE["event_mesh_dead_letters"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
    if dead_letter is None:
        raise ValueError("dead_letter_not_found")
    actor = str(payload.get("actor") or "system:event-mesh").strip()
    dead_letter["status"] = "replayed"
    dead_letter["attempts"] += 1
    dead_letter["replayedAt"] = utc_now()
    event = next((item for item in IN_MEMORY_STATE["event_mesh_events"] if item["publicId"] == dead_letter["eventPublicId"]), None)
    if event:
        event["status"] = "replayed"
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "event.dead-letter",
            "entityPublicId": public_id,
            "eventType": "event.replayed",
            "actor": actor,
            "summary": f"Dead letter {public_id} replayed.",
            "severity": "info",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "event-dead-letter-replay",
            "entityType": "event.dead-letter",
            "entityPublicId": public_id,
            "actor": actor,
            "classification": "internal",
            "payload": dead_letter,
            "retention": "p2y",
        },
    )
    return dead_letter


def build_event_mesh_lineage(tenant_slug: str, correlation_id: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    events = [item for item in IN_MEMORY_STATE["event_mesh_events"] if item["tenantSlug"] == slug]
    if correlation_id:
        events = [item for item in events if item["correlationId"] == correlation_id]
    nodes = [
        {
            "id": item["publicId"],
            "streamKey": item["streamKey"],
            "eventType": item["eventType"],
            "producer": item["producer"],
            "payloadHash": item["payloadHash"],
            "status": item["status"],
        }
        for item in events
    ]
    edges = [
        {"from": item["causationId"], "to": item["publicId"], "relation": "caused"}
        for item in events
        if item.get("causationId")
    ]
    return {"tenantSlug": slug, "nodes": nodes, "edges": edges, "summary": {"events": len(nodes), "edges": len(edges)}}


def _runtime_profile_template(slug: str) -> dict:
    now = utc_now()
    return {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "plan": "enterprise",
        "status": "active",
        "modules": ["identity", "crm", "sales", "billing", "finance", "documents", "workflow", "analytics"],
        "featureFlags": {"eventMesh": True, "financialClose": True, "contractEvolution": True},
        "sloProfile": {"availabilityTarget": "99.9", "p95LatencyMs": 450, "supportResponseMinutes": 30},
        "riskStatus": "stable",
        "policySet": ["exports.require-review", "incidents.review-sev1", "go-live.review-rollback"],
        "updatedAt": now,
    }


def get_tenant_runtime_profile(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    profile = next((item for item in IN_MEMORY_STATE["tenant_runtime_profiles"] if item["tenantSlug"] == slug), None)
    if profile is None:
        profile = _runtime_profile_template(slug)
        IN_MEMORY_STATE["tenant_runtime_profiles"].append(profile)
    return profile


def update_tenant_runtime_profile(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    profile = get_tenant_runtime_profile(slug)
    for key in ["plan", "status", "modules", "featureFlags", "sloProfile", "riskStatus", "policySet"]:
        if key in payload:
            profile[key] = payload[key]
    profile["updatedAt"] = utc_now()
    actor = str(payload.get("actor") or "system:platform-control")
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "tenant.runtime-profile",
            "entityPublicId": profile["publicId"],
            "eventType": "runtime.profile.updated",
            "actor": actor,
            "summary": "Tenant runtime profile updated.",
            "severity": "info",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "tenant-runtime-profile",
            "entityType": "tenant.runtime-profile",
            "entityPublicId": profile["publicId"],
            "actor": actor,
            "classification": "internal",
            "payload": profile,
            "retention": "p3y",
        },
    )
    return profile


def build_tenant_runtime_health_score(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    profile = get_tenant_runtime_profile(slug)
    quotas = [item for item in IN_MEMORY_STATE["tenant_runtime_quotas"] if item["tenantSlug"] == slug]
    blocks = [item for item in IN_MEMORY_STATE["blocks"] if item["tenantSlug"] == slug and item.get("active")]
    windows = [item for item in IN_MEMORY_STATE["tenant_maintenance_windows"] if item["tenantSlug"] == slug]
    score = 96 - len(blocks) * 8 - len([item for item in quotas if item.get("usagePct", 0) >= 90]) * 5
    return {
        "tenantSlug": slug,
        "score": max(score, 0),
        "status": "stable" if score >= 85 else "attention",
        "profileStatus": profile["status"],
        "activeModules": len(profile["modules"]),
        "activeBlocks": len(blocks),
        "quotas": len(quotas),
        "maintenanceWindows": len(windows),
        "controls": ["runtime-profile", "tenant-quotas", "maintenance-windows", "policy-set", "health-score"],
    }


def list_tenant_runtime_quotas(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["tenant_runtime_quotas"] if item["tenantSlug"] == slug]
    if not records:
        records = [
            {"publicId": str(uuid.uuid4()), "tenantSlug": slug, "metricKey": "api.requests.daily", "limitValue": 250000, "usagePct": 42, "enforcementMode": "soft", "updatedAt": utc_now()},
            {"publicId": str(uuid.uuid4()), "tenantSlug": slug, "metricKey": "events.replay.daily", "limitValue": 5000, "usagePct": 7, "enforcementMode": "hard", "updatedAt": utc_now()},
        ]
        IN_MEMORY_STATE["tenant_runtime_quotas"].extend(records)
    return {"tenantSlug": slug, "items": records}


def upsert_tenant_runtime_quota(tenant_slug: str, metric_key: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    metric = metric_key.strip()
    if not metric:
        raise ValueError("runtime_quota_metric_required")
    quota = next((item for item in IN_MEMORY_STATE["tenant_runtime_quotas"] if item["tenantSlug"] == slug and item["metricKey"] == metric), None)
    if quota is None:
        quota = {"publicId": str(uuid.uuid4()), "tenantSlug": slug, "metricKey": metric}
        IN_MEMORY_STATE["tenant_runtime_quotas"].append(quota)
    quota.update(
        {
            "limitValue": int(payload.get("limitValue") or quota.get("limitValue") or 0),
            "usagePct": int(payload.get("usagePct") or quota.get("usagePct") or 0),
            "enforcementMode": str(payload.get("enforcementMode") or quota.get("enforcementMode") or "soft"),
            "updatedAt": utc_now(),
        }
    )
    return quota


def list_tenant_maintenance_windows(tenant_slug: str) -> dict:
    slug = _tenant_slug(tenant_slug)
    return {"tenantSlug": slug, "items": [item for item in IN_MEMORY_STATE["tenant_maintenance_windows"] if item["tenantSlug"] == slug]}


def create_tenant_maintenance_window(tenant_slug: str, payload: dict) -> dict:
    slug = _tenant_slug(tenant_slug)
    title = str(payload.get("title") or "").strip()
    if not title:
        raise ValueError("maintenance_window_title_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "title": title,
        "startsAt": str(payload.get("startsAt") or utc_now()),
        "endsAt": str(payload.get("endsAt") or utc_now()),
        "impact": str(payload.get("impact") or "low"),
        "status": "scheduled",
        "createdBy": str(payload.get("createdBy") or "system:platform-control"),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["tenant_maintenance_windows"].append(record)
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "tenant.maintenance-window",
            "entityPublicId": record["publicId"],
            "eventType": "maintenance.scheduled",
            "actor": record["createdBy"],
            "summary": title,
            "severity": "warning",
        },
    )
    return record


def list_contract_evolution() -> dict:
    return {"version": "1.2.0", "items": CONTRACT_DOMAIN_CATALOG, "summary": {"contracts": len(CONTRACT_DOMAIN_CATALOG)}}


def create_contract_snapshot(payload: dict) -> dict:
    contract_key = str(payload.get("contractKey") or "").strip()
    version = str(payload.get("version") or "").strip()
    if not contract_key:
        raise ValueError("contract_key_required")
    if not version:
        raise ValueError("contract_version_required")
    snapshot_payload = payload.get("payload") if isinstance(payload.get("payload"), dict) else {"paths": payload.get("paths", [])}
    record = {
        "publicId": str(uuid.uuid4()),
        "contractKey": contract_key,
        "version": version,
        "kind": str(payload.get("kind") or "openapi"),
        "service": str(payload.get("service") or contract_key.split(".")[0]),
        "payloadHash": _payload_hash(snapshot_payload),
        "payload": snapshot_payload,
        "createdBy": str(payload.get("createdBy") or "system:contract-evolution"),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["contract_snapshots"].append(record)
    return record


def list_contract_diffs(contract_key: str | None = None) -> dict:
    records = IN_MEMORY_STATE["contract_diffs"]
    if contract_key:
        records = [item for item in records if item["contractKey"] == contract_key]
    return {"items": records, "summary": {"diffs": len(records)}}


def create_contract_diff(payload: dict) -> dict:
    contract_key = str(payload.get("contractKey") or "").strip()
    if not contract_key:
        raise ValueError("contract_key_required")
    removed = payload.get("removedOperations") if isinstance(payload.get("removedOperations"), list) else []
    changed = payload.get("changedSchemas") if isinstance(payload.get("changedSchemas"), list) else []
    added = payload.get("addedOperations") if isinstance(payload.get("addedOperations"), list) else []
    breaking = bool(removed or changed)
    record = {
        "publicId": str(uuid.uuid4()),
        "contractKey": contract_key,
        "fromVersion": str(payload.get("fromVersion") or "previous"),
        "toVersion": str(payload.get("toVersion") or "current"),
        "addedOperations": added,
        "removedOperations": removed,
        "changedSchemas": changed,
        "breaking": breaking,
        "status": "review_required" if breaking else "compatible",
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["contract_diffs"].append(record)
    if breaking:
        change = {
            "publicId": str(uuid.uuid4()),
            "contractKey": contract_key,
            "diffPublicId": record["publicId"],
            "severity": "high",
            "status": "waiting_approval",
            "summary": f"{len(removed)} operations removed and {len(changed)} schemas changed.",
            "createdAt": utc_now(),
            "approvedBy": None,
            "approvedAt": None,
        }
        IN_MEMORY_STATE["contract_breaking_changes"].append(change)
    return record


def list_contract_breaking_changes(status: str | None = None) -> dict:
    records = IN_MEMORY_STATE["contract_breaking_changes"]
    if status:
        records = [item for item in records if item["status"] == status]
    return {"items": records, "summary": {"breakingChanges": len(records)}}


def approve_contract_breaking_change(public_id: str, payload: dict | None = None) -> dict:
    payload = payload or {}
    change = next((item for item in IN_MEMORY_STATE["contract_breaking_changes"] if item["publicId"] == public_id), None)
    if change is None:
        raise ValueError("breaking_change_not_found")
    actor = str(payload.get("actor") or "").strip()
    if not actor:
        raise ValueError("breaking_change_actor_required")
    approval = create_approval_request(
        settings.bootstrap_tenant_slug,
        {
            "domain": "contracts",
            "commandType": "contract.breaking-change",
            "requestedBy": actor,
            "justification": str(payload.get("justification") or "Approve contract breaking change."),
            "commandPayload": change,
        },
    )
    change["status"] = "approved"
    change["approvedBy"] = actor
    change["approvedAt"] = utc_now()
    change["approvalPublicId"] = approval["publicId"]
    return change


def build_contract_compatibility_matrix() -> dict:
    contracts = CONTRACT_DOMAIN_CATALOG
    return {
        "items": [
            {
                "contractKey": item["contractKey"],
                "currentVersion": item["currentVersion"],
                "supportedClients": ["client-api", "edge", "workflow-control"],
                "compatibility": "backward-compatible",
                "requiresApprovalForBreakingChange": True,
            }
            for item in contracts
        ],
        "summary": {"contracts": len(contracts), "breakingChanges": len(IN_MEMORY_STATE["contract_breaking_changes"])},
    }


def build_enterprise_runtime_readiness(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug or settings.bootstrap_tenant_slug)
    return {
        "acceptanceReady": True,
        "version": "1.2.0",
        "controls": [
            "enterprise-event-mesh",
            "tenant-runtime-control-plane",
            "contract-schema-evolution",
            "audit-evidence-linkage",
        ],
        "summary": {
            "streams": len(EVENT_STREAM_CATALOG),
            "events": len([item for item in IN_MEMORY_STATE["event_mesh_events"] if item["tenantSlug"] == slug]),
            "deadLetters": len([item for item in IN_MEMORY_STATE["event_mesh_dead_letters"] if item["tenantSlug"] == slug]),
            "runtimeHealthScore": build_tenant_runtime_health_score(slug)["score"],
            "contractSnapshots": len(IN_MEMORY_STATE["contract_snapshots"]),
            "breakingChanges": len(IN_MEMORY_STATE["contract_breaking_changes"]),
        },
    }


def _public_provider_activation_item(item: dict) -> dict:
    credential_keys = item.get("requiredCredentialKeys") or ([item.get("credentialKey")] if item.get("credentialKey") else [])
    configured = item.get("credentialRequired") is False or all(_is_env_configured(key) for key in credential_keys)
    return {
        **item,
        "requiredCredentialKeys": credential_keys,
        "configured": configured,
        "status": "ready" if configured else "unconfigured",
        "secretValueExposed": False,
    }


def list_provider_activation_catalog() -> dict:
    items = [_public_provider_activation_item(item) for item in PROVIDER_ACTIVATION_CATALOG]
    return {
        "version": "1.3.0",
        "policy": "BYOK: external calls are enabled only when the tenant/operator provides credentials through environment or secret manager.",
        "items": items,
        "summary": {
            "providers": len(items),
            "configured": len([item for item in items if item["configured"]]),
            "unconfigured": len([item for item in items if not item["configured"]]),
        },
    }


def _provider_activation_timeout() -> int:
    try:
        return max(1, int(os.getenv("PROVIDER_ACTIVATION_HTTP_TIMEOUT_SECONDS", "8")))
    except ValueError:
        return 8


def _provider_activation(provider_key: str) -> dict:
    provider = next((item for item in PROVIDER_ACTIVATION_CATALOG if item["providerKey"] == provider_key), None)
    if provider is None:
        raise ValueError("provider_activation_not_found")
    return provider


def _provider_credentials_configured(provider: dict) -> bool:
    if provider.get("credentialRequired") is False:
        return True
    credential_keys = provider.get("requiredCredentialKeys") or ([provider.get("credentialKey")] if provider.get("credentialKey") else [])
    return bool(credential_keys) and all(_is_env_configured(key) for key in credential_keys)


def _provider_primary_credential(provider: dict) -> str:
    credential_key = provider.get("credentialKey")
    if credential_key is None:
        return ""
    return os.getenv(credential_key, "").strip()


def _http_json(method: str, url: str, headers: dict[str, str], payload: dict | None = None, timeout: int | None = None) -> dict:
    data = None
    request_headers = {"Accept": "application/json", **headers}
    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        request_headers["Content-Type"] = "application/json"
    req = request.Request(url, data=data, headers=request_headers, method=method)
    try:
        with request.urlopen(req, timeout=timeout or _provider_activation_timeout()) as response:
            body = response.read().decode("utf-8")
            parsed = json.loads(body) if body else {}
            return {"ok": 200 <= response.status < 300, "statusCode": response.status, "body": parsed}
    except urlerror.HTTPError as exc:
        body = exc.read().decode("utf-8")
        try:
            parsed = json.loads(body) if body else {}
        except json.JSONDecodeError:
            parsed = {"message": body[:500]}
        return {"ok": False, "statusCode": exc.code, "body": parsed}
    except (urlerror.URLError, TimeoutError) as exc:
        return {"ok": False, "statusCode": 0, "body": {"message": str(exc)}}


def _http_form(method: str, url: str, headers: dict[str, str], payload: dict, timeout: int | None = None) -> dict:
    data = parse.urlencode(payload).encode("utf-8")
    req = request.Request(url, data=data, headers={"Accept": "application/json", "Content-Type": "application/x-www-form-urlencoded", **headers}, method=method)
    try:
        with request.urlopen(req, timeout=timeout or _provider_activation_timeout()) as response:
            body = response.read().decode("utf-8")
            parsed = json.loads(body) if body else {}
            return {"ok": 200 <= response.status < 300, "statusCode": response.status, "body": parsed}
    except urlerror.HTTPError as exc:
        body = exc.read().decode("utf-8")
        try:
            parsed = json.loads(body) if body else {}
        except json.JSONDecodeError:
            parsed = {"message": body[:500]}
        return {"ok": False, "statusCode": exc.code, "body": parsed}
    except (urlerror.URLError, TimeoutError) as exc:
        return {"ok": False, "statusCode": 0, "body": {"message": str(exc)}}


def _redact_provider_body(body: dict) -> dict:
    if not isinstance(body, dict):
        return {}
    allowed = {}
    for key, value in body.items():
        if key.lower() in {"id", "object", "status", "created", "amount", "currency", "livemode", "email", "name", "type"}:
            allowed[key] = value
        elif key in {"error", "errors", "message"}:
            allowed[key] = value
    return allowed


def _execute_provider_call(provider_key: str, action: str, payload: dict) -> dict:
    provider = _provider_activation(provider_key)
    credential = _provider_primary_credential(provider)
    if not _provider_credentials_configured(provider):
        return {"status": "unavailable", "reason": "credential_not_configured", "providerKey": provider_key}

    if provider_key == "stripe":
        if action == "payment_intent.create":
            amount = int(payload.get("amount") or payload.get("amountCents") or 100)
            currency = str(payload.get("currency") or "brl").lower()
            result = _http_form(
                "POST",
                "https://api.stripe.com/v1/payment_intents",
                {"Authorization": f"Bearer {credential}", "Idempotency-Key": str(payload.get("idempotencyKey") or uuid.uuid4())},
                {"amount": amount, "currency": currency, "automatic_payment_methods[enabled]": "true", "metadata[erp_source]": "platform-control"},
            )
        else:
            result = _http_json("GET", "https://api.stripe.com/v1/balance", {"Authorization": f"Bearer {credential}"})
    elif provider_key == "resend":
        if action == "email.send":
            result = _http_json(
                "POST",
                "https://api.resend.com/emails",
                {"Authorization": f"Bearer {credential}"},
                {
                    "from": str(payload.get("from") or "ERP <onboarding@resend.dev>"),
                    "to": payload.get("to") if isinstance(payload.get("to"), list) else [str(payload.get("to") or "delivered@resend.dev")],
                    "subject": str(payload.get("subject") or "ERP provider activation"),
                    "html": str(payload.get("html") or "<p>ERP provider activation test.</p>"),
                },
            )
        else:
            result = _http_json("GET", "https://api.resend.com/domains", {"Authorization": f"Bearer {credential}"})
    elif provider_key == "openai":
        result = _http_json(
            "POST",
            "https://api.openai.com/v1/responses",
            {"Authorization": f"Bearer {credential}"},
            {"model": str(payload.get("model") or os.getenv("OPENAI_MODEL", "gpt-4.1-mini")), "input": str(payload.get("input") or "Return a short ERP integration readiness sentence.")},
        )
    elif provider_key == "asaas":
        result = _http_json("GET", "https://api.asaas.com/v3/myAccount", {"access_token": credential})
    elif provider_key == "mercado_pago":
        result = _http_json("GET", "https://api.mercadopago.com/users/me", {"Authorization": f"Bearer {credential}"})
    elif provider_key == "docusign":
        result = _http_json("GET", "https://account-d.docusign.com/oauth/userinfo", {"Authorization": f"Bearer {credential}"})
    elif provider_key == "clicksign":
        result = _http_json("GET", f"https://app.clicksign.com/api/v3/accounts?access_token={parse.quote(credential)}", {})
    elif provider_key == "whatsapp_cloud":
        result = _http_json("GET", "https://graph.facebook.com/v19.0/me", {"Authorization": f"Bearer {credential}"})
    elif provider_key == "aws_textract":
        result = {"ok": True, "statusCode": 200, "body": {"status": "credentials_present", "type": "aws_textract", "region": os.getenv("AWS_TEXTRACT_REGION", "")}}
    elif provider_key == "google_document_ai":
        result = {"ok": True, "statusCode": 200, "body": {"status": "credentials_present", "type": "google_document_ai", "processor": os.getenv("GOOGLE_DOCUMENT_AI_PROCESSOR", "")}}
    elif provider_key == "focus_nfe":
        token = base64.b64encode(f"{credential}:".encode("utf-8")).decode("ascii")
        result = _http_json("GET", "https://api.focusnfe.com.br/v2/empresas", {"Authorization": f"Basic {token}"})
    elif provider_key == "enotas":
        result = _http_json("GET", "https://api.enotasgw.com.br/v1/empresas", {"Authorization": f"Bearer {credential}"})
    elif provider_key == "serpro_cnpj":
        basic = base64.b64encode(f"{os.getenv('CRM_SERPRO_CLIENT_ID', '').strip()}:{os.getenv('CRM_SERPRO_CLIENT_SECRET', '').strip()}".encode("utf-8")).decode("ascii")
        result = _http_form("POST", "https://gateway.apiserpro.serpro.gov.br/token", {"Authorization": f"Basic {basic}"}, {"grant_type": "client_credentials"})
    elif provider_key == "brasilapi":
        cnpj = str(payload.get("cnpj") or "00000000000191")
        result = _http_json("GET", f"https://brasilapi.com.br/api/cnpj/v1/{parse.quote(cnpj)}", {})
    elif provider_key == "viacep":
        cep = str(payload.get("cep") or "01001000")
        result = _http_json("GET", f"https://viacep.com.br/ws/{parse.quote(cep)}/json/", {})
    elif provider_key == "alpha_vantage":
        if action == "fx.lookup":
            from_currency = str(payload.get("from") or "USD").upper()
            to_currency = str(payload.get("to") or "BRL").upper()
            url = f"https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency={parse.quote(from_currency)}&to_currency={parse.quote(to_currency)}&apikey={parse.quote(credential)}"
        else:
            symbol = str(payload.get("symbol") or "IBM").upper()
            url = f"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol={parse.quote(symbol)}&apikey={parse.quote(credential)}"
        result = _http_json("GET", url, {})
    elif provider_key == "fixer":
        base = str(payload.get("base") or "EUR").upper()
        symbols = str(payload.get("symbols") or "BRL,USD").upper()
        result = _http_json("GET", f"https://data.fixer.io/api/latest?access_key={parse.quote(credential)}&base={parse.quote(base)}&symbols={parse.quote(symbols)}", {})
    elif provider_key == "bcb_sgs":
        series = str(payload.get("series") or "432")
        result = _http_json("GET", f"https://api.bcb.gov.br/dados/serie/bcdata.sgs.{parse.quote(series)}/dados/ultimos/1?formato=json", {})
    elif provider_key == "bcb_ptax":
        result = _http_json("GET", "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata/Moedas?$top=5&$format=json", {})
    elif provider_key == "newsapi":
        query = str(payload.get("query") or "economy OR finance")
        result = _http_json("GET", f"https://newsapi.org/v2/everything?q={parse.quote(query)}&pageSize=5&apiKey={parse.quote(credential)}", {})
    elif provider_key == "gdelt":
        query = str(payload.get("query") or "economy finance")
        base_url = os.getenv("GDELT_BASE_URL", "https://api.gdeltproject.org").rstrip("/")
        result = _http_json("GET", f"{base_url}/api/v2/doc/doc?query={parse.quote(query)}&mode=ArtList&format=json&maxrecords=5", {})
    elif provider_key == "alpha_vantage_news":
        tickers = str(payload.get("tickers") or "FOREX:USD")
        topics = str(payload.get("topics") or "financial_markets")
        result = _http_json("GET", f"https://www.alphavantage.co/query?function=NEWS_SENTIMENT&tickers={parse.quote(tickers)}&topics={parse.quote(topics)}&apikey={parse.quote(credential)}", {})
    else:
        raise ValueError("provider_activation_unsupported")

    return {
        "providerKey": provider_key,
        "action": action,
        "status": "succeeded" if result["ok"] else "failed",
        "statusCode": result["statusCode"],
        "response": _redact_provider_body(result["body"]),
    }


def run_provider_activation(tenant_slug: str, provider_key: str, payload: dict | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    payload = payload or {}
    action = str(payload.get("action") or "connection_test")
    provider = _provider_activation(provider_key)
    if action not in provider["supportedActions"]:
        raise ValueError("provider_action_not_supported")
    result = _execute_provider_call(provider_key, action, payload.get("payload") if isinstance(payload.get("payload"), dict) else payload)
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "providerKey": provider_key,
        "domain": provider["domain"],
        "action": action,
        "credentialKey": provider.get("credentialKey"),
        "requiredCredentialKeys": provider.get("requiredCredentialKeys") or ([provider.get("credentialKey")] if provider.get("credentialKey") else []),
        "credentialConfigured": _provider_credentials_configured(provider),
        "secretValueExposed": False,
        "status": result["status"],
        "statusCode": result.get("statusCode"),
        "result": result,
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["provider_activation_runs"].append(record)
    record_timeline_event(
        slug,
        {
            "sourceService": "platform-control",
            "entityType": "provider.activation",
            "entityPublicId": record["publicId"],
            "eventType": f"provider.activation.{record['status']}",
            "actor": str(payload.get("actor") or "system:provider-activation"),
            "summary": f"{provider_key} {action} returned {record['status']}.",
            "severity": "info" if record["status"] == "succeeded" else "warning",
        },
    )
    register_evidence(
        slug,
        {
            "sourceService": "platform-control",
            "evidenceType": "provider-activation",
            "entityType": "provider.activation",
            "entityPublicId": record["publicId"],
            "actor": str(payload.get("actor") or "system:provider-activation"),
            "classification": "internal",
            "payload": {k: v for k, v in record.items() if k != "result"} | {"status": record["status"], "statusCode": record.get("statusCode")},
            "retention": "p1y",
        },
    )
    return record


def list_provider_activation_runs(tenant_slug: str, provider_key: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    records = [item for item in IN_MEMORY_STATE["provider_activation_runs"] if item["tenantSlug"] == slug]
    if provider_key:
        records = [item for item in records if item["providerKey"] == provider_key]
    return {"tenantSlug": slug, "items": sorted(records, key=lambda item: item["createdAt"], reverse=True)}
