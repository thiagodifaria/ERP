"""Relatorio consolidado de hardening enterprise."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.adapter_catalog import build_adapter_catalog
from app.reports.load_benchmark import build_postgres_load_benchmark
from app.reports.platform_reliability import (
    fetch_billing_metrics,
    fetch_webhook_metrics,
    fetch_workflow_metrics,
)


def build_hardening_review(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_hardening_review(tenant_slug)

    return build_static_hardening_review(tenant_slug)


def build_static_hardening_review(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    adapter_catalog = build_adapter_catalog()
    reviews = {
        "security": {"status": "stable", "mfaEnabledUsers": 2, "activeSessions": 3, "auditEvents": 18},
        "observability": {"status": "stable", "dashboardsReady": True, "probeCoverage": 8},
        "retriesAndDlq": {"status": "attention", "pendingWebhookEvents": 1, "deadLetterEvents": 1, "failedPaymentAttempts": 2},
        "backupRestore": {"status": "stable", "validated": True},
        "slos": {"status": "attention", "webhookForwardingRate": 0.94, "workflowSuccessRate": 0.91, "billingRecoveryRate": 0.78},
        "multiTenant": {"status": "stable", "tenantAccessGuardrailsReady": True, "scopedServicesReviewed": 10},
        "failover": {"status": "stable", "requeueReady": True, "replayReady": True},
        "performance": {"status": "attention", "latestBenchmarkStatus": "attention", "throughputRps": 46.8, "p95LatencyMs": 332},
        "permissions": {"status": "stable", "openfgaReady": True, "sessionEnforcementReady": True, "auditTrailReady": True},
        "providerCapabilities": {
            "status": "attention",
            "configuredCapabilities": adapter_catalog["summary"]["configuredCapabilities"],
            "criticalProviderGaps": adapter_catalog["summary"]["criticalUnconfiguredCapabilities"],
            "fallbackCapabilities": adapter_catalog["summary"]["fallbackCapabilities"],
        },
        "contractGovernance": {
            "status": "stable",
            "httpSpecs": adapter_catalog["contracts"]["summary"]["httpSpecs"],
            "eventSchemas": adapter_catalog["contracts"]["summary"]["eventSchemas"],
            "adrRecorded": adapter_catalog["contracts"]["summary"]["adrRecorded"],
        },
    }

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "summary": summarize_reviews(reviews),
        "reviews": reviews,
    }


def build_postgres_hardening_review(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    adapter_catalog = build_adapter_catalog()

    with connect() as connection:
        webhook_metrics = fetch_webhook_metrics(connection)
        workflow_metrics = fetch_workflow_metrics(connection, tenant_slug)
        billing_metrics = fetch_billing_metrics(connection, tenant_slug)
        identity_metrics = fetch_identity_security_metrics(connection, tenant_slug)

    benchmark = build_postgres_load_benchmark(tenant_slug)
    latest_benchmark = benchmark.get("latest") or {}

    security_status = "stable" if identity_metrics["auditEvents"] > 0 and identity_metrics["mfaEnabledUsers"] > 0 else "attention"
    retry_status = "stable"
    if webhook_metrics["deadLetterEvents"] > 0 or billing_metrics["failedPaymentAttempts"] > 0:
        retry_status = "attention"
    if webhook_metrics["deadLetterEvents"] > 1:
        retry_status = "critical"

    slo_status = "stable"
    webhook_rate = round(webhook_metrics["forwardedEvents"] / webhook_metrics["totalEvents"], 4) if webhook_metrics["totalEvents"] > 0 else 1
    workflow_rate = round(workflow_metrics["completedExecutions"] / workflow_metrics["totalExecutions"], 4) if workflow_metrics["totalExecutions"] > 0 else 1
    billing_rate = round(billing_metrics["recoveredCases"] / billing_metrics["recoveryCases"], 4) if billing_metrics["recoveryCases"] > 0 else 1
    if webhook_rate < 1 or workflow_rate < 1 or billing_rate < 1:
        slo_status = "attention"
    if webhook_rate < 0.9 or workflow_rate < 0.9 or billing_rate < 0.75:
        slo_status = "critical"

    performance_status = latest_benchmark.get("status", "stable")
    provider_status = "stable" if adapter_catalog["summary"]["criticalUnconfiguredCapabilities"] == 0 else "attention"

    reviews = {
        "security": {
            "status": security_status,
            "mfaEnabledUsers": identity_metrics["mfaEnabledUsers"],
            "activeSessions": identity_metrics["activeSessions"],
            "auditEvents": identity_metrics["auditEvents"],
        },
        "observability": {
            "status": "stable",
            "dashboardsReady": True,
            "probeCoverage": 8,
        },
        "retriesAndDlq": {
            "status": retry_status,
            "pendingWebhookEvents": webhook_metrics["pendingWebhookEvents"],
            "deadLetterEvents": webhook_metrics["deadLetterEvents"],
            "failedPaymentAttempts": billing_metrics["failedPaymentAttempts"],
        },
        "backupRestore": {
            "status": "stable",
            "validated": True,
        },
        "slos": {
            "status": slo_status,
            "webhookForwardingRate": webhook_rate,
            "workflowSuccessRate": workflow_rate,
            "billingRecoveryRate": billing_rate,
        },
        "multiTenant": {
            "status": "stable",
            "tenantAccessGuardrailsReady": True,
            "scopedServicesReviewed": 10,
        },
        "failover": {
            "status": "stable" if webhook_metrics["deadLetterEvents"] <= 1 else "attention",
            "requeueReady": True,
            "replayReady": True,
        },
        "performance": {
            "status": performance_status,
            "latestBenchmarkStatus": performance_status,
            "throughputRps": float(latest_benchmark.get("throughputRps", 0.0) or 0.0),
            "p95LatencyMs": int(latest_benchmark.get("p95LatencyMs", 0) or 0),
        },
        "permissions": {
            "status": "stable",
            "openfgaReady": True,
            "sessionEnforcementReady": identity_metrics["activeSessions"] > 0,
            "auditTrailReady": identity_metrics["auditEvents"] > 0,
        },
        "providerCapabilities": {
            "status": provider_status,
            "configuredCapabilities": adapter_catalog["summary"]["configuredCapabilities"],
            "criticalProviderGaps": adapter_catalog["summary"]["criticalUnconfiguredCapabilities"],
            "fallbackCapabilities": adapter_catalog["summary"]["fallbackCapabilities"],
        },
        "contractGovernance": {
            "status": "stable",
            "httpSpecs": adapter_catalog["contracts"]["summary"]["httpSpecs"],
            "eventSchemas": adapter_catalog["contracts"]["summary"]["eventSchemas"],
            "adrRecorded": adapter_catalog["contracts"]["summary"]["adrRecorded"],
        },
    }

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "summary": summarize_reviews(reviews),
        "reviews": reviews,
    }


def fetch_identity_security_metrics(connection, tenant_slug: str | None) -> dict:
    params: list[str] = []
    where = ""
    if tenant_slug:
        params.append(tenant_slug)
        where = "WHERE tenant.slug = %s"

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM identity.user_security_profiles AS profile
                JOIN identity.users AS user_account ON user_account.id = profile.user_id
                JOIN identity.tenants AS tenant ON tenant.id = user_account.tenant_id
                {where}
                  AND profile.mfa_enabled = true
            ) AS mfa_enabled_users,
            (
                SELECT count(*)
                FROM identity.sessions AS session
                JOIN identity.tenants AS tenant ON tenant.id = session.tenant_id
                {where}
                  AND session.status = 'active'
            ) AS active_sessions,
            (
                SELECT count(*)
                FROM identity.security_audit_events AS audit
                JOIN identity.tenants AS tenant ON tenant.id = audit.tenant_id
                {where}
            ) AS audit_events
    """
    multiplier = 3 if tenant_slug else 1
    with connection.cursor() as cursor:
        cursor.execute(query, params * multiplier)
        row = cursor.fetchone() or {}

    return {
        "mfaEnabledUsers": int(row.get("mfa_enabled_users", 0) or 0),
        "activeSessions": int(row.get("active_sessions", 0) or 0),
        "auditEvents": int(row.get("audit_events", 0) or 0),
    }


def summarize_reviews(reviews: dict) -> dict:
    stable = 0
    attention = 0
    critical = 0

    for review in reviews.values():
        status = review.get("status", "stable")
        if status == "critical":
            critical += 1
        elif status == "attention":
            attention += 1
        else:
            stable += 1

    overall_status = "stable"
    if critical > 0:
        overall_status = "critical"
    elif attention > 0:
        overall_status = "attention"

    return {
        "status": overall_status,
        "stableChecks": stable,
        "attentionChecks": attention,
        "criticalChecks": critical,
    }
