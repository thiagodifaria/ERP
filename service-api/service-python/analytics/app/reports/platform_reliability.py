"""Relatorio executivo de confiabilidade da plataforma."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_platform_reliability(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_platform_reliability(tenant_slug)

    return build_static_platform_reliability(tenant_slug)


def build_static_platform_reliability(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "stability": {
            "status": "attention",
            "pendingWebhookEvents": 1,
            "deadLetterEvents": 1,
            "failedWorkflowExecutions": 2,
            "criticalRecoveryCases": 1,
            "failedPaymentAttempts": 3,
        },
        "serviceLevelObjectives": {
            "webhookForwardingRate": 0.94,
            "workflowSuccessRate": 0.91,
            "billingRecoveryRate": 0.78,
            "status": "attention",
        },
        "safeguards": {
            "backupRestoreValidated": True,
            "observabilityReady": True,
            "permissionGuardrailsReady": True,
            "multiTenantReviewed": True,
            "dlqReady": True,
            "openCriticalRisks": 1,
        },
    }


def build_postgres_platform_reliability(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        webhook_metrics = fetch_webhook_metrics(connection)
        workflow_metrics = fetch_workflow_metrics(connection, tenant_slug)
        billing_metrics = fetch_billing_metrics(connection, tenant_slug)

    pending_webhook_events = webhook_metrics["pendingWebhookEvents"]
    dead_letter_events = webhook_metrics["deadLetterEvents"]
    failed_workflow_executions = workflow_metrics["failedWorkflowExecutions"]
    critical_recovery_cases = billing_metrics["criticalRecoveryCases"]
    failed_payment_attempts = billing_metrics["failedPaymentAttempts"]

    stability_status = "stable"
    if (
        dead_letter_events > 0
        or failed_workflow_executions > 0
        or critical_recovery_cases > 0
        or failed_payment_attempts > 1
    ):
        stability_status = "attention"
    if dead_letter_events > 1 or failed_workflow_executions > 1 or critical_recovery_cases > 1:
        stability_status = "critical"

    webhook_total = webhook_metrics["totalEvents"]
    workflow_total = workflow_metrics["totalExecutions"]
    recovery_total = billing_metrics["recoveryCases"]

    webhook_forwarding_rate = (
        round(webhook_metrics["forwardedEvents"] / webhook_total, 4) if webhook_total > 0 else 1
    )
    workflow_success_rate = (
        round(workflow_metrics["completedExecutions"] / workflow_total, 4) if workflow_total > 0 else 1
    )
    billing_recovery_rate = (
        round(billing_metrics["recoveredCases"] / recovery_total, 4) if recovery_total > 0 else 1
    )

    safeguards = {
        "backupRestoreValidated": True,
        "observabilityReady": True,
        "permissionGuardrailsReady": True,
        "multiTenantReviewed": True,
        "dlqReady": True,
        "openCriticalRisks": int(dead_letter_events > 0)
        + int(failed_workflow_executions > 1)
        + int(critical_recovery_cases > 0),
    }

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "stability": {
            "status": stability_status,
            "pendingWebhookEvents": pending_webhook_events,
            "deadLetterEvents": dead_letter_events,
            "failedWorkflowExecutions": failed_workflow_executions,
            "criticalRecoveryCases": critical_recovery_cases,
            "failedPaymentAttempts": failed_payment_attempts,
        },
        "serviceLevelObjectives": {
            "webhookForwardingRate": webhook_forwarding_rate,
            "workflowSuccessRate": workflow_success_rate,
            "billingRecoveryRate": billing_recovery_rate,
            "status": stability_status,
        },
        "safeguards": safeguards,
    }


def fetch_webhook_metrics(connection) -> dict:
    with connection.cursor() as cursor:
        cursor.execute(
            """
                SELECT
                    count(*) AS total_events,
                    count(*) FILTER (WHERE status = 'forwarded') AS forwarded_events,
                    count(*) FILTER (WHERE status IN ('validated', 'queued', 'processing')) AS pending_events,
                    count(*) FILTER (WHERE status = 'dead_letter') AS dead_letter_events
                FROM webhook_hub.webhook_events
            """
        )
        row = cursor.fetchone() or {}

    return {
        "totalEvents": int(row.get("total_events", 0) or 0),
        "forwardedEvents": int(row.get("forwarded_events", 0) or 0),
        "pendingWebhookEvents": int(row.get("pending_events", 0) or 0),
        "deadLetterEvents": int(row.get("dead_letter_events", 0) or 0),
    }


def fetch_workflow_metrics(connection, tenant_slug: str | None) -> dict:
    params: list[str] = []
    where = ""
    if tenant_slug:
        params.append(tenant_slug)
        where = "WHERE tenant.slug = %s"

    query = f"""
        SELECT
            count(*) AS total_executions,
            count(*) FILTER (WHERE execution.status = 'completed') AS completed_executions,
            count(*) FILTER (WHERE execution.status = 'failed') AS failed_workflow_executions
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {where}
    """
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "totalExecutions": int(row.get("total_executions", 0) or 0),
        "completedExecutions": int(row.get("completed_executions", 0) or 0),
        "failedWorkflowExecutions": int(row.get("failed_workflow_executions", 0) or 0),
    }


def fetch_billing_metrics(connection, tenant_slug: str | None) -> dict:
    params: list[str] = []
    where = ""
    if tenant_slug:
        params.append(tenant_slug)
        where = "WHERE tenant.slug = %s"

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM billing.recovery_cases AS recovery
                JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
                {where}
            ) AS recovery_cases,
            (
                SELECT count(*)
                FROM billing.recovery_cases AS recovery
                JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
                {where}
                  AND recovery.status = 'recovered'
            ) AS recovered_cases,
            (
                SELECT count(*)
                FROM billing.recovery_cases AS recovery
                JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
                {where}
                  AND recovery.severity = 'critical'
                  AND recovery.status <> 'recovered'
            ) AS critical_recovery_cases,
            (
                SELECT count(*)
                FROM billing.payment_attempts AS attempt
                JOIN billing.subscription_invoices AS invoice ON invoice.id = attempt.invoice_id
                JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id
                {where}
                  AND attempt.status = 'failed'
            ) AS failed_payment_attempts
    """
    multiplier = 4 if tenant_slug else 1
    with connection.cursor() as cursor:
        cursor.execute(query, params * multiplier)
        row = cursor.fetchone() or {}

    return {
        "recoveryCases": int(row.get("recovery_cases", 0) or 0),
        "recoveredCases": int(row.get("recovered_cases", 0) or 0),
        "criticalRecoveryCases": int(row.get("critical_recovery_cases", 0) or 0),
        "failedPaymentAttempts": int(row.get("failed_payment_attempts", 0) or 0),
    }
