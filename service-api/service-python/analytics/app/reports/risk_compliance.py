"""Risk and compliance scoring for the version 1.1.0 governance controls."""

from __future__ import annotations

from datetime import datetime, timezone


RISK_DIMENSIONS = [
    {
        "dimension": "security",
        "score": 82,
        "status": "stable",
        "signals": ["JWT/OpenFGA middleware", "gateway hardening", "restricted workloads"],
        "recommendation": "Keep CodeQL, dependency review and strict secret posture enabled.",
    },
    {
        "dimension": "privacy",
        "score": 78,
        "status": "attention",
        "signals": ["LGPD inventory", "search redaction", "legal hold controls"],
        "recommendation": "Expand automated retention execution evidence for all document classes.",
    },
    {
        "dimension": "financial",
        "score": 84,
        "status": "stable",
        "signals": ["idempotent payment attempts", "ledger controls", "period close"],
        "recommendation": "Connect provider reconciliation evidence to the evidence vault.",
    },
    {
        "dimension": "operational",
        "score": 88,
        "status": "stable",
        "signals": ["incident command", "runbook automation", "go-live readiness"],
        "recommendation": "Attach every SEV1 runbook to a postmortem and timeline event.",
    },
    {
        "dimension": "provider",
        "score": 73,
        "status": "attention",
        "signals": ["provider defaults", "fallback modes", "adapter catalog"],
        "recommendation": "Replace local fallback providers with configured production credentials by tenant.",
    },
    {
        "dimension": "contract",
        "score": 91,
        "status": "stable",
        "signals": ["OpenAPI registry", "route-vs-contract checks", "client-api generation"],
        "recommendation": "Block releases when route coverage diverges from OpenAPI.",
    },
    {
        "dimension": "incident",
        "score": 86,
        "status": "stable",
        "signals": ["timeline", "action tracking", "postmortem readiness"],
        "recommendation": "Drive escalation through approval policies and evidence snapshots.",
    },
    {
        "dimension": "ai",
        "score": 80,
        "status": "stable",
        "signals": ["tool allowlist", "read-only default", "prompt redaction"],
        "recommendation": "Require approval before enabling any mutating AI tool.",
    },
]


DOMAIN_SCORES = [
    {"domain": "platform-control", "score": 89, "status": "stable", "drivers": ["policies", "approvals", "runbooks", "evidence"]},
    {"domain": "analytics", "score": 87, "status": "stable", "drivers": ["semantic metrics", "risk scoring", "readiness"]},
    {"domain": "search", "score": 81, "status": "stable", "drivers": ["redaction", "legal hold", "exports"]},
    {"domain": "ai-governance", "score": 80, "status": "stable", "drivers": ["tool policies", "audit", "redaction"]},
    {"domain": "billing-finance", "score": 83, "status": "stable", "drivers": ["idempotency", "ledger", "recovery"]},
    {"domain": "providers", "score": 73, "status": "attention", "drivers": ["fallback", "unconfigured critical credentials"]},
]


SERVICE_SCORES = [
    {"service": "platform-control", "score": 90, "status": "stable"},
    {"service": "analytics", "score": 88, "status": "stable"},
    {"service": "search", "score": 82, "status": "stable"},
    {"service": "ai-governance", "score": 80, "status": "stable"},
    {"service": "webhook-hub", "score": 78, "status": "attention"},
    {"service": "billing", "score": 82, "status": "stable"},
]


def _average_score(items: list[dict]) -> int:
    return round(sum(int(item["score"]) for item in items) / max(len(items), 1))


def build_tenant_risk_score(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    score = _average_score(RISK_DIMENSIONS)
    status = "stable" if score >= 80 else "attention"
    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "score": score,
        "status": status,
        "dimensions": RISK_DIMENSIONS,
        "controls": ["policy-decision-center", "approval-workflows", "runbook-automation", "audit-evidence-vault"],
    }


def build_domain_scores(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    return {
        "tenantSlug": slug,
        "items": DOMAIN_SCORES,
        "summary": {
            "domains": len(DOMAIN_SCORES),
            "stable": sum(1 for item in DOMAIN_SCORES if item["status"] == "stable"),
            "attention": sum(1 for item in DOMAIN_SCORES if item["status"] == "attention"),
        },
    }


def build_service_scores(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    return {
        "tenantSlug": slug,
        "items": SERVICE_SCORES,
        "summary": {
            "services": len(SERVICE_SCORES),
            "averageScore": _average_score(SERVICE_SCORES),
        },
    }


def build_compliance_posture(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    return {
        "tenantSlug": slug,
        "status": "stable",
        "requirements": [
            {"key": "lgpd-purpose", "status": "covered", "evidence": "docs/SEGURANCA.md"},
            {"key": "auditability", "status": "covered", "evidence": "platform-control evidence vault"},
            {"key": "approval-for-sensitive-command", "status": "covered", "evidence": "command approval workflows"},
            {"key": "ai-read-only-default", "status": "covered", "evidence": "ai-governance policies"},
            {"key": "contract-drift-block", "status": "covered", "evidence": "scripts/test.sh contract"},
        ],
    }


def build_risk_recommendations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    recommendations = [
        {"priority": "high", "domain": "providers", "summary": "Replace fallback providers with configured credentials for critical capabilities."},
        {"priority": "high", "domain": "approval", "summary": "Require approvals for export, rollback, quota and provider fallback commands."},
        {"priority": "medium", "domain": "runbooks", "summary": "Attach SEV1 and go-live rollback runbooks to evidence vault snapshots."},
        {"priority": "medium", "domain": "ai", "summary": "Keep AI tools read-only until a dedicated mutation approval policy exists."},
    ]
    return {"tenantSlug": slug, "items": recommendations, "summary": {"total": len(recommendations), "high": 2}}


def build_risk_readiness() -> dict:
    tenant_score = build_tenant_risk_score()
    domain_scores = build_domain_scores()
    return {
        "acceptanceReady": tenant_score["score"] >= 80 and domain_scores["summary"]["attention"] <= 1,
        "controls": ["tenant-risk-score", "domain-risk-score", "service-risk-score", "compliance-posture", "recommendations"],
        "score": tenant_score["score"],
    }
