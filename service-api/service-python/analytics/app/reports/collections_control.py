"""Relatorio executivo do plano de cobranca e recuperacao de receita."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_collections_control(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_collections_control(tenant_slug)

    return build_static_collections_control(tenant_slug)


def build_static_collections_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "portfolio": {
            "casesTotal": 12,
            "openCases": 3,
            "contactedCases": 2,
            "promisedCases": 2,
            "recoveredCases": 4,
            "defaultedCases": 1,
            "criticalCases": 4,
            "openAmountCents": 286000,
            "promisedAmountCents": 121000,
            "recoveredAmountCents": 308000,
        },
        "invoices": {
            "failedAttempts": 7,
            "invoicesInRecovery": 5,
            "overdueInvoices": 3,
            "overdueAmountCents": 286000,
            "averageFailedAttemptsPerCase": 2.3,
        },
        "promises": {
            "activePromises": 2,
            "promisesDueSoon": 1,
            "brokenPromises": 1,
            "promisesKept": 3,
        },
        "throughput": {
            "touchpoints": 18,
            "promiseActions": 6,
            "recoveries": 4,
            "recoveryRate": 0.3333,
            "averageRecoveryLagDays": 5.2,
        },
        "governance": {
            "pendingActions": 3,
            "nextActionsDue": 2,
            "oldestOpenCaseDays": 11,
            "lastResolvedAt": "2026-04-29T14:20:00+00:00",
        },
    }


def build_postgres_collections_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        portfolio = fetch_portfolio_metrics(connection, tenant_slug)
        invoice_metrics = fetch_invoice_recovery_metrics(connection, tenant_slug)
        promise_metrics = fetch_promise_metrics(connection, tenant_slug)
        throughput_metrics = fetch_throughput_metrics(connection, tenant_slug)
        governance_metrics = fetch_governance_metrics(connection, tenant_slug)

    total_cases = portfolio["casesTotal"]
    recoveries = throughput_metrics["recoveries"]
    recovery_rate = round(recoveries / total_cases, 4) if total_cases > 0 else 0.0
    average_failed_attempts = (
        round(invoice_metrics["failedAttempts"] / total_cases, 2) if total_cases > 0 else 0.0
    )

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "portfolio": portfolio,
        "invoices": {
            **invoice_metrics,
            "averageFailedAttemptsPerCase": average_failed_attempts,
        },
        "promises": promise_metrics,
        "throughput": {
            **throughput_metrics,
            "recoveryRate": recovery_rate,
        },
        "governance": governance_metrics,
    }


def fetch_portfolio_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS cases_total,
            count(*) FILTER (WHERE recovery.status = 'open') AS open_cases,
            count(*) FILTER (WHERE recovery.status = 'contacted') AS contacted_cases,
            count(*) FILTER (WHERE recovery.status = 'promised') AS promised_cases,
            count(*) FILTER (WHERE recovery.status = 'recovered') AS recovered_cases,
            count(*) FILTER (WHERE recovery.status = 'defaulted') AS defaulted_cases,
            count(*) FILTER (WHERE recovery.severity = 'critical') AS critical_cases,
            COALESCE(sum(invoice.amount_cents) FILTER (
                WHERE recovery.status IN ('open', 'contacted', 'promised')
            ), 0) AS open_amount_cents,
            COALESCE(sum(invoice.amount_cents) FILTER (WHERE recovery.status = 'promised'), 0) AS promised_amount_cents,
            COALESCE(sum(invoice.amount_cents) FILTER (WHERE recovery.status = 'recovered'), 0) AS recovered_amount_cents
        FROM billing.recovery_cases AS recovery
        JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
        JOIN billing.subscription_invoices AS invoice
          ON invoice.tenant_id = recovery.tenant_id
         AND invoice.public_id = recovery.invoice_public_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "casesTotal": int(row.get("cases_total", 0) or 0),
        "openCases": int(row.get("open_cases", 0) or 0),
        "contactedCases": int(row.get("contacted_cases", 0) or 0),
        "promisedCases": int(row.get("promised_cases", 0) or 0),
        "recoveredCases": int(row.get("recovered_cases", 0) or 0),
        "defaultedCases": int(row.get("defaulted_cases", 0) or 0),
        "criticalCases": int(row.get("critical_cases", 0) or 0),
        "openAmountCents": int(row.get("open_amount_cents", 0) or 0),
        "promisedAmountCents": int(row.get("promised_amount_cents", 0) or 0),
        "recoveredAmountCents": int(row.get("recovered_amount_cents", 0) or 0),
    }


def fetch_invoice_recovery_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM billing.payment_attempts AS attempt
                JOIN identity.tenants AS tenant ON tenant.id = attempt.tenant_id
                {filter_sql}
                  AND attempt.status = 'failed'
            ) AS failed_attempts,
            (
                SELECT count(*)
                FROM billing.recovery_cases AS recovery
                JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
                {filter_sql}
                  AND recovery.status IN ('open', 'contacted', 'promised')
            ) AS invoices_in_recovery,
            (
                SELECT count(*)
                FROM billing.subscription_invoices AS invoice
                JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id
                {filter_sql}
                  AND invoice.status = 'open'
                  AND invoice.due_date < timezone('utc', now())::date
            ) AS overdue_invoices,
            (
                SELECT COALESCE(sum(invoice.amount_cents), 0)
                FROM billing.subscription_invoices AS invoice
                JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id
                {filter_sql}
                  AND invoice.status = 'open'
                  AND invoice.due_date < timezone('utc', now())::date
            ) AS overdue_amount_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 4)
        row = cursor.fetchone() or {}

    return {
        "failedAttempts": int(row.get("failed_attempts", 0) or 0),
        "invoicesInRecovery": int(row.get("invoices_in_recovery", 0) or 0),
        "overdueInvoices": int(row.get("overdue_invoices", 0) or 0),
        "overdueAmountCents": int(row.get("overdue_amount_cents", 0) or 0),
    }


def fetch_promise_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) FILTER (WHERE recovery.status = 'promised') AS active_promises,
            count(*) FILTER (
                WHERE recovery.status = 'promised'
                  AND recovery.promised_payment_date IS NOT NULL
                  AND recovery.promised_payment_date <= timezone('utc', now())::date + INTERVAL '3 days'
            ) AS promises_due_soon,
            count(*) FILTER (
                WHERE recovery.status = 'promised'
                  AND recovery.promised_payment_date IS NOT NULL
                  AND recovery.promised_payment_date < timezone('utc', now())::date
            ) AS broken_promises,
            count(*) FILTER (
                WHERE recovery.status = 'recovered'
                  AND recovery.promised_payment_date IS NOT NULL
                  AND recovery.resolved_at IS NOT NULL
                  AND recovery.resolved_at::date <= recovery.promised_payment_date
            ) AS promises_kept
        FROM billing.recovery_cases AS recovery
        JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "activePromises": int(row.get("active_promises", 0) or 0),
        "promisesDueSoon": int(row.get("promises_due_soon", 0) or 0),
        "brokenPromises": int(row.get("broken_promises", 0) or 0),
        "promisesKept": int(row.get("promises_kept", 0) or 0),
    }


def fetch_throughput_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM billing.recovery_actions AS action
                JOIN identity.tenants AS tenant ON tenant.id = action.tenant_id
                {filter_sql}
                  AND action.action_code = 'touchpoint_registered'
            ) AS touchpoints,
            (
                SELECT count(*)
                FROM billing.recovery_actions AS action
                JOIN identity.tenants AS tenant ON tenant.id = action.tenant_id
                {filter_sql}
                  AND action.action_code = 'promise_registered'
            ) AS promise_actions,
            (
                SELECT count(*)
                FROM billing.recovery_actions AS action
                JOIN identity.tenants AS tenant ON tenant.id = action.tenant_id
                {filter_sql}
                  AND action.action_code = 'case_recovered'
            ) AS recoveries,
            (
                SELECT COALESCE(avg(EXTRACT(EPOCH FROM (recovery.resolved_at - recovery.created_at)) / 86400.0), 0)
                FROM billing.recovery_cases AS recovery
                JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
                {filter_sql}
                  AND recovery.status = 'recovered'
                  AND recovery.resolved_at IS NOT NULL
            ) AS average_recovery_lag_days
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 4)
        row = cursor.fetchone() or {}

    return {
        "touchpoints": int(row.get("touchpoints", 0) or 0),
        "promiseActions": int(row.get("promise_actions", 0) or 0),
        "recoveries": int(row.get("recoveries", 0) or 0),
        "averageRecoveryLagDays": round(float(row.get("average_recovery_lag_days", 0) or 0), 2),
    }


def fetch_governance_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) FILTER (
                WHERE recovery.status IN ('open', 'contacted', 'promised')
                  AND recovery.next_action_at IS NOT NULL
            ) AS pending_actions,
            count(*) FILTER (
                WHERE recovery.status IN ('open', 'contacted', 'promised')
                  AND recovery.next_action_at IS NOT NULL
                  AND recovery.next_action_at <= timezone('utc', now()) + INTERVAL '24 hours'
            ) AS next_actions_due,
            COALESCE(max((timezone('utc', now())::date - recovery.created_at::date)) FILTER (
                WHERE recovery.status IN ('open', 'contacted', 'promised')
            ), 0) AS oldest_open_case_days,
            max(recovery.resolved_at) AS last_resolved_at
        FROM billing.recovery_cases AS recovery
        JOIN identity.tenants AS tenant ON tenant.id = recovery.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "pendingActions": int(row.get("pending_actions", 0) or 0),
        "nextActionsDue": int(row.get("next_actions_due", 0) or 0),
        "oldestOpenCaseDays": int(row.get("oldest_open_case_days", 0) or 0),
        "lastResolvedAt": row.get("last_resolved_at").isoformat() if row.get("last_resolved_at") else None,
    }
