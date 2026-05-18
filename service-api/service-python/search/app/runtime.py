from __future__ import annotations

from datetime import datetime, timezone
import json
import re
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect

SCHEMA_NAME = "search"
SENSITIVE_PATTERN = re.compile(r"([\\w.%-]+@[\\w.-]+\\.[A-Za-z]{2,}|\\b\\d{3}\\.?\\d{3}\\.?\\d{3}-?\\d{2}\\b|\\b\\d{2}\\.?\\d{3}\\.?\\d{3}/?\\d{4}-?\\d{2}\\b|Bearer\\s+[A-Za-z0-9._-]+|access[_-]?token[:=][A-Za-z0-9._-]+)", re.IGNORECASE)

IN_MEMORY_STATE: dict[str, list[dict]] = {
    "index_entries": [],
    "saved_queries": [],
    "query_audit_events": [],
    "discovery_cases": [],
    "discovery_case_items": [],
    "legal_holds": [],
    "export_requests": [],
}

BOOTSTRAP_ENTRIES = [
    {
        "tenantSlug": "bootstrap-ops",
        "entityType": "crm.lead",
        "entityPublicId": "lead-northwind-001",
        "title": "Lead Northwind ERP expansion",
        "summary": "Lead enterprise com contato financeiro e interesse em billing, workflow e documents.",
        "content": "northwind cfo@example.com CNPJ 12.345.678/0001-90 billing workflow documents",
        "classification": "restricted",
        "tags": ["crm", "lead", "billing"],
        "metadata": {"source": "crm", "owner": "sales"},
    },
    {
        "tenantSlug": "bootstrap-ops",
        "entityType": "documents.attachment",
        "entityPublicId": "doc-contract-001",
        "title": "Contrato de assinatura enterprise",
        "summary": "Contrato com retencao fiscal e assinatura digital pendente.",
        "content": "contrato assinatura digital fiscal retention access link",
        "classification": "confidential",
        "tags": ["documents", "contract", "fiscal"],
        "metadata": {"retention": "long-term", "storage": "r2"},
    },
    {
        "tenantSlug": "bootstrap-ops",
        "entityType": "support.case",
        "entityPublicId": "case-sla-001",
        "title": "Caso SLA em atencao",
        "summary": "Atendimento com SLA proximo do vencimento e dependencia de notificacao.",
        "content": "support sla notification incident customer success",
        "classification": "internal",
        "tags": ["support", "sla"],
        "metadata": {"queue": "enterprise", "severity": "medium"},
    },
]

for entry in BOOTSTRAP_ENTRIES:
    IN_MEMORY_STATE["index_entries"].append({**entry, "publicId": str(uuid.uuid4()), "indexedAt": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")})


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def find_tenant_id(cursor, slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def redact_text(value: str, allow_sensitive: bool = False) -> tuple[str, list[str]]:
    findings: list[str] = []
    if allow_sensitive:
        return value, findings

    def replace(match: re.Match) -> str:
        token = match.group(0)
        kind = "token" if "Bearer" in token or "token" in token.lower() else "identifier"
        findings.append(kind)
        return f"[REDACTED:{kind}]"

    return SENSITIVE_PATTERN.sub(replace, value), findings


def capability_catalog() -> dict:
    return {
        "service": settings.service_name,
        "capabilities": [
            {"key": "operational-search", "status": "ready"},
            {"key": "tenant-faceted-query", "status": "ready"},
            {"key": "ediscovery-cases", "status": "ready"},
            {"key": "legal-holds", "status": "ready"},
            {"key": "query-audit", "status": "ready"},
            {"key": "pii-redaction", "status": "ready"},
        ],
    }


def query_index(tenant: str | None, q: str | None, entity_type: str | None = None, actor: str | None = None, include_sensitive: bool = False) -> dict:
    slug = tenant_slug(tenant)
    query_text = (q or "").strip().lower()
    entity_filter = (entity_type or "").strip()
    records = _list_index_entries(slug)
    if entity_filter:
        records = [item for item in records if item["entityType"] == entity_filter]
    if query_text:
        records = [
            item
            for item in records
            if query_text in " ".join([item.get("title", ""), item.get("summary", ""), item.get("content", ""), " ".join(item.get("tags", []))]).lower()
        ]
    redacted = [_redact_entry(item, include_sensitive) for item in records[:50]]
    audit = _audit_query(slug, actor or "system:search", query_text or "*", len(redacted), include_sensitive)
    return {
        "tenantSlug": slug,
        "query": q or "*",
        "entityType": entity_filter or None,
        "total": len(redacted),
        "items": redacted,
        "auditEvent": audit,
    }


def build_facets(tenant: str | None) -> dict:
    slug = tenant_slug(tenant)
    records = _list_index_entries(slug)
    entity_counts: dict[str, int] = {}
    tag_counts: dict[str, int] = {}
    classification_counts: dict[str, int] = {}
    for item in records:
        entity_counts[item["entityType"]] = entity_counts.get(item["entityType"], 0) + 1
        classification_counts[item["classification"]] = classification_counts.get(item["classification"], 0) + 1
        for tag in item.get("tags", []):
            tag_counts[tag] = tag_counts.get(tag, 0) + 1
    return {"tenantSlug": slug, "entities": entity_counts, "tags": tag_counts, "classifications": classification_counts}


def create_saved_query(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    name = str(payload.get("name") or "").strip()
    query = str(payload.get("query") or "").strip()
    if not name:
        raise ValueError("saved_query_name_required")
    if not query:
        raise ValueError("saved_query_query_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "name": name,
        "query": query,
        "filters": payload.get("filters") or {},
        "createdBy": str(payload.get("createdBy") or "system:search"),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["saved_queries"].append(record)
    return record


def list_audit_events(tenant: str | None) -> dict:
    slug = tenant_slug(tenant)
    return {"tenantSlug": slug, "items": [item for item in IN_MEMORY_STATE["query_audit_events"] if item["tenantSlug"] == slug]}


def create_discovery_case(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    title = str(payload.get("title") or "").strip()
    if not title:
        raise ValueError("discovery_case_title_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "title": title,
        "status": "open",
        "owner": str(payload.get("owner") or "legal@erp.local"),
        "scope": payload.get("scope") or {},
        "createdAt": utc_now(),
        "items": [],
    }
    IN_MEMORY_STATE["discovery_cases"].append(record)
    return record


def add_case_item(case_public_id: str, payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    case = _find_case(slug, case_public_id)
    if case is None:
        raise ValueError("discovery_case_not_found")
    entity_public_id = str(payload.get("entityPublicId") or "").strip()
    entity_type = str(payload.get("entityType") or "").strip()
    if not entity_public_id or not entity_type:
        raise ValueError("discovery_case_item_required")
    item = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "casePublicId": case_public_id,
        "entityType": entity_type,
        "entityPublicId": entity_public_id,
        "reason": str(payload.get("reason") or "manual selection"),
        "addedAt": utc_now(),
    }
    IN_MEMORY_STATE["discovery_case_items"].append(item)
    case["items"].append(item)
    return item


def create_legal_hold(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    reason = str(payload.get("reason") or "").strip()
    if not reason:
        raise ValueError("legal_hold_reason_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "entityType": str(payload.get("entityType") or "tenant"),
        "entityPublicId": str(payload.get("entityPublicId") or slug),
        "reason": reason,
        "status": "active",
        "createdBy": str(payload.get("createdBy") or "legal@erp.local"),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["legal_holds"].append(record)
    return record


def create_export_request(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    query = str(payload.get("query") or "").strip()
    if not query:
        raise ValueError("export_query_required")
    holds = [item for item in IN_MEMORY_STATE["legal_holds"] if item["tenantSlug"] == slug and item["status"] == "active"]
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "query": query,
        "status": "held" if holds else "queued",
        "requestedBy": str(payload.get("requestedBy") or "ops@erp.local"),
        "format": str(payload.get("format") or "jsonl"),
        "legalHoldCount": len(holds),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["export_requests"].append(record)
    return record


def _list_index_entries(slug: str) -> list[dict]:
    if settings.repository_driver != "postgres":
        return [item for item in IN_MEMORY_STATE["index_entries"] if item["tenantSlug"] == slug]
    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT entry.public_id, entry.entity_type, entry.entity_public_id, entry.title, entry.summary,
                       entry.content_text, entry.classification, entry.tags_json, entry.metadata_json, entry.indexed_at
                FROM search.index_entries AS entry
                JOIN identity.tenants AS tenant ON tenant.id = entry.tenant_id
                WHERE tenant.slug = %s
                ORDER BY entry.indexed_at DESC
                """,
                (slug,),
            )
            return [
                {
                    "publicId": str(row["public_id"]),
                    "tenantSlug": slug,
                    "entityType": row["entity_type"],
                    "entityPublicId": row["entity_public_id"],
                    "title": row["title"],
                    "summary": row["summary"],
                    "content": row["content_text"],
                    "classification": row["classification"],
                    "tags": row["tags_json"] or [],
                    "metadata": row["metadata_json"] or {},
                    "indexedAt": row["indexed_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def _audit_query(slug: str, actor: str, query: str, result_count: int, include_sensitive: bool) -> dict:
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "actor": actor,
        "query": query,
        "resultCount": result_count,
        "sensitiveAccess": include_sensitive,
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["query_audit_events"].append(record)
    return record


def _redact_entry(item: dict, include_sensitive: bool) -> dict:
    summary, summary_findings = redact_text(item.get("summary", ""), include_sensitive)
    title, title_findings = redact_text(item.get("title", ""), include_sensitive)
    content, content_findings = redact_text(item.get("content", ""), include_sensitive)
    return {
        **item,
        "title": title,
        "summary": summary,
        "contentPreview": content[:180],
        "redaction": {
            "applied": not include_sensitive and bool(summary_findings or title_findings or content_findings),
            "findings": sorted(set(summary_findings + title_findings + content_findings)),
        },
    }


def _find_case(slug: str, public_id: str) -> dict | None:
    return next((item for item in IN_MEMORY_STATE["discovery_cases"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)

