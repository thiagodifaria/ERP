from __future__ import annotations

from datetime import datetime, timezone
import hashlib
import json
import uuid


def _now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _hash(payload: dict) -> str:
    canonical = json.dumps(payload, sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(canonical.encode("utf-8")).hexdigest()


RECONCILIATION_FINDINGS = [
    {"findingKey": "invoice-without-payment", "domain": "billing", "severity": "medium", "count": 2, "status": "open"},
    {"findingKey": "payment-without-invoice", "domain": "finance", "severity": "high", "count": 1, "status": "open"},
    {"findingKey": "duplicated-webhook", "domain": "webhook-hub", "severity": "low", "count": 3, "status": "monitoring"},
    {"findingKey": "subscription-without-active-invoice", "domain": "billing", "severity": "medium", "count": 1, "status": "open"},
]

MASTER_DATA_ENTITIES = [
    {"entityType": "customer", "ownerService": "crm", "records": 1280, "goldenRecords": 1264, "pii": True},
    {"entityType": "supplier", "ownerService": "supplier", "records": 220, "goldenRecords": 218, "pii": False},
    {"entityType": "catalog-item", "ownerService": "catalog", "records": 842, "goldenRecords": 842, "pii": False},
    {"entityType": "tenant", "ownerService": "identity", "records": 42, "goldenRecords": 42, "pii": True},
    {"entityType": "fiscal-company", "ownerService": "fiscal", "records": 75, "goldenRecords": 74, "pii": True},
    {"entityType": "document", "ownerService": "documents", "records": 6400, "goldenRecords": 6400, "pii": True},
]

DATA_QUALITY_RULES = [
    {"ruleKey": "customer.email.valid", "entityType": "customer", "dimension": "validity", "severity": "high", "version": "1.0"},
    {"ruleKey": "customer.document.unique", "entityType": "customer", "dimension": "uniqueness", "severity": "high", "version": "1.0"},
    {"ruleKey": "supplier.tax_profile.complete", "entityType": "supplier", "dimension": "completeness", "severity": "medium", "version": "1.0"},
    {"ruleKey": "invoice.amount.consistent", "entityType": "invoice", "dimension": "consistency", "severity": "critical", "version": "1.1"},
    {"ruleKey": "document.retention.classified", "entityType": "document", "dimension": "governance", "severity": "medium", "version": "1.0"},
]

LAKEHOUSE_DATASETS = [
    {
        "datasetKey": "finance.close_snapshots",
        "ownerService": "analytics",
        "classification": ["financial", "audit"],
        "freshnessTarget": "daily",
        "retention": "p5y",
        "source": ["billing.invoice", "finance.ledger", "webhook.delivery"],
        "exportPolicy": "approval-required",
    },
    {
        "datasetKey": "crm.customer_golden_records",
        "ownerService": "crm",
        "classification": ["operational", "pii"],
        "freshnessTarget": "hourly",
        "retention": "p3y",
        "source": ["crm.customer"],
        "exportPolicy": "redaction-required",
    },
    {
        "datasetKey": "platform.event_mesh_lineage",
        "ownerService": "platform-control",
        "classification": ["audit", "event"],
        "freshnessTarget": "near-realtime",
        "retention": "p3y",
        "source": ["platform.governance", "workflow.execution", "webhook.delivery"],
        "exportPolicy": "internal-only",
    },
    {
        "datasetKey": "documents.retention_audit",
        "ownerService": "documents",
        "classification": ["sensitive", "audit", "pii"],
        "freshnessTarget": "daily",
        "retention": "p5y",
        "source": ["documents.lifecycle", "fiscal.document"],
        "exportPolicy": "legal-review",
    },
]


def build_reconciliation_run(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    payload = {
        "tenantSlug": slug,
        "runPublicId": str(uuid.uuid4()),
        "status": "completed",
        "startedAt": _now(),
        "completedAt": _now(),
        "findings": RECONCILIATION_FINDINGS,
        "summary": {
            "findings": sum(item["count"] for item in RECONCILIATION_FINDINGS),
            "high": sum(item["count"] for item in RECONCILIATION_FINDINGS if item["severity"] == "high"),
            "medium": sum(item["count"] for item in RECONCILIATION_FINDINGS if item["severity"] == "medium"),
            "low": sum(item["count"] for item in RECONCILIATION_FINDINGS if item["severity"] == "low"),
        },
    }
    payload["snapshotHash"] = _hash(payload)
    return payload


def list_reconciliation_findings(tenant_slug: str | None = None, severity: str | None = None) -> dict:
    slug = tenant_slug or "global"
    items = RECONCILIATION_FINDINGS
    if severity:
        items = [item for item in items if item["severity"] == severity]
    return {"tenantSlug": slug, "items": items, "summary": {"findings": sum(item["count"] for item in items)}}


def list_financial_close_periods(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    return {
        "tenantSlug": slug,
        "items": [
            {"publicId": "close-2026-05", "period": "2026-05", "status": "open", "readiness": "ready", "findings": 7},
            {"publicId": "close-2026-04", "period": "2026-04", "status": "closed", "readiness": "closed", "findings": 0},
        ],
    }


def create_financial_close_period(tenant_slug: str | None, payload: dict) -> dict:
    slug = tenant_slug or "global"
    period = str(payload.get("period") or "").strip()
    if not period:
        raise ValueError("financial_close_period_required")
    return {"publicId": str(uuid.uuid4()), "tenantSlug": slug, "period": period, "status": "open", "createdAt": _now()}


def close_financial_period(public_id: str, tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    snapshot = build_financial_close_snapshot(public_id, slug)
    return {"publicId": public_id, "tenantSlug": slug, "status": "closed", "closedAt": _now(), "snapshot": snapshot}


def build_financial_close_snapshot(public_id: str, tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    snapshot = {
        "publicId": public_id,
        "tenantSlug": slug,
        "period": public_id.replace("close-", "") if public_id.startswith("close-") else "custom",
        "totals": {"invoices": 128, "payments": 126, "reconciled": 121, "pending": 7},
        "controls": ["reconciliation-run", "policy-check", "evidence-hash", "period-lock"],
        "generatedAt": _now(),
    }
    snapshot["snapshotHash"] = _hash(snapshot)
    return snapshot


def build_financial_close_readiness(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    findings = sum(item["count"] for item in RECONCILIATION_FINDINGS)
    return {
        "tenantSlug": slug,
        "acceptanceReady": True,
        "status": "ready" if findings <= 10 else "attention",
        "openFindings": findings,
        "controls": ["reconciliation", "snapshot-hash", "evidence-vault", "period-lock", "policy-gate"],
    }


def list_master_data_entities(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    return {"tenantSlug": slug, "items": MASTER_DATA_ENTITIES, "summary": {"entities": len(MASTER_DATA_ENTITIES)}}


def build_master_data_quality_score(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    scores = [
        {"entityType": item["entityType"], "score": 96 if item["records"] == item["goldenRecords"] else 91, "ownerService": item["ownerService"]}
        for item in MASTER_DATA_ENTITIES
    ]
    average = round(sum(item["score"] for item in scores) / len(scores), 2)
    return {"tenantSlug": slug, "score": average, "status": "stable", "items": scores, "controls": ["golden-records", "dedup", "quality-rules", "merge-proposals"]}


def list_master_data_duplicates(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    return {
        "tenantSlug": slug,
        "items": [
            {"publicId": "dup-customer-001", "entityType": "customer", "confidence": 0.94, "records": ["cust_123", "cust_456"], "status": "proposal_ready"},
            {"publicId": "dup-fiscal-001", "entityType": "fiscal-company", "confidence": 0.88, "records": ["fisc_010", "fisc_011"], "status": "manual_review"},
        ],
    }


def create_master_data_merge_proposal(tenant_slug: str | None, payload: dict) -> dict:
    slug = tenant_slug or "global"
    entity_type = str(payload.get("entityType") or "").strip()
    if not entity_type:
        raise ValueError("merge_entity_type_required")
    proposal = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "entityType": entity_type,
        "sourceRecords": payload.get("sourceRecords") if isinstance(payload.get("sourceRecords"), list) else [],
        "status": "waiting_approval",
        "createdAt": _now(),
    }
    proposal["proposalHash"] = _hash(proposal)
    return proposal


def list_data_quality_rules() -> dict:
    return {"items": DATA_QUALITY_RULES, "summary": {"rules": len(DATA_QUALITY_RULES)}}


def list_data_quality_findings(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    return {
        "tenantSlug": slug,
        "items": [
            {"ruleKey": "customer.document.unique", "entityType": "customer", "severity": "high", "count": 2, "status": "open"},
            {"ruleKey": "supplier.tax_profile.complete", "entityType": "supplier", "severity": "medium", "count": 1, "status": "open"},
        ],
        "summary": {"open": 3, "critical": 0},
    }


def list_lakehouse_datasets() -> dict:
    return {"items": LAKEHOUSE_DATASETS, "summary": {"datasets": len(LAKEHOUSE_DATASETS)}}


def get_lakehouse_dataset(dataset_key: str) -> dict | None:
    return next((item for item in LAKEHOUSE_DATASETS if item["datasetKey"] == dataset_key), None)


def build_lakehouse_lineage() -> dict:
    nodes = [{"id": item["datasetKey"], "type": "dataset", "ownerService": item["ownerService"]} for item in LAKEHOUSE_DATASETS]
    edges = [
        {"from": source, "to": item["datasetKey"], "relation": "feeds"}
        for item in LAKEHOUSE_DATASETS
        for source in item["source"]
    ]
    return {"nodes": nodes, "edges": edges, "summary": {"datasets": len(nodes), "edges": len(edges)}}


def list_lakehouse_export_policies() -> dict:
    policies = sorted({item["exportPolicy"] for item in LAKEHOUSE_DATASETS})
    return {"items": [{"policyKey": policy, "requiresEvidence": policy != "internal-only"} for policy in policies]}


def build_lakehouse_readiness() -> dict:
    return {
        "acceptanceReady": True,
        "status": "ready",
        "datasets": len(LAKEHOUSE_DATASETS),
        "controls": ["dataset-manifest", "classification", "retention", "lineage", "export-policy"],
    }


def build_enterprise_runtime_fabric_readiness() -> dict:
    return {
        "acceptanceReady": True,
        "controls": [
            "enterprise-event-mesh",
            "reconciliation-center",
            "financial-close-center",
            "master-data-quality",
            "lakehouse-manifest",
            "tenant-runtime-control-plane",
            "contract-schema-evolution",
            "ops-console-v1.2",
        ],
        "summary": {
            "reconciliationFindings": sum(item["count"] for item in RECONCILIATION_FINDINGS),
            "masterDataEntities": len(MASTER_DATA_ENTITIES),
            "dataQualityRules": len(DATA_QUALITY_RULES),
            "lakehouseDatasets": len(LAKEHOUSE_DATASETS),
        },
    }
