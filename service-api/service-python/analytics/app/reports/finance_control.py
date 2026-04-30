"""Relatorio executivo cruzando financeiro operacional, tesouraria e billing."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_finance_control(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_finance_control(tenant_slug)

    return build_static_finance_control(tenant_slug)


def build_static_finance_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "treasury": {
            "accountsTotal": 3,
            "activeAccounts": 2,
            "movementsTotal": 19,
            "currentBalanceCents": 1246000,
            "inflowCents": 845000,
            "outflowCents": 214000,
        },
        "receivables": {
            "total": 14,
            "open": 4,
            "paid": 9,
            "cancelled": 1,
            "openAmountCents": 312000,
            "paidAmountCents": 845000,
            "overdueCount": 1,
        },
        "payables": {
            "total": 6,
            "open": 2,
            "paid": 4,
            "cancelled": 0,
            "openAmountCents": 88000,
            "paidAmountCents": 173000,
            "dueSoonCount": 1,
        },
        "billing": {
            "plansTotal": 3,
            "activePlans": 3,
            "activeSubscriptions": 11,
            "graceSubscriptions": 1,
            "suspendedSubscriptions": 1,
            "invoicesOpen": 2,
            "invoicesPaid": 9,
            "invoicesFailed": 1,
            "failedAttempts": 3,
            "succeededAttempts": 9,
            "monthlyRecurringRevenueCents": 539000,
        },
        "profitability": {
            "collectedReceivablesCents": 845000,
            "releasedCommissionsCents": 91000,
            "operatingCostCents": 128000,
            "paidPayablesCents": 173000,
            "netOperationalMarginCents": 453000,
            "netTreasuryMovementCents": 631000,
        },
        "governance": {
            "periodClosures": 4,
            "subscriptionEvents": 23,
            "cashMovements": 19,
            "failedPaymentAttempts": 3,
        },
    }


def build_postgres_finance_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        treasury = fetch_treasury_metrics(connection, tenant_slug)
        receivables = fetch_receivable_metrics(connection, tenant_slug)
        payables = fetch_payable_metrics(connection, tenant_slug)
        billing = fetch_billing_metrics(connection, tenant_slug)
        governance = fetch_governance_metrics(connection, tenant_slug)

    profitability = {
        "collectedReceivablesCents": receivables["paidAmountCents"],
        "releasedCommissionsCents": governance["releasedCommissionsCents"],
        "operatingCostCents": governance["operatingCostCents"],
        "paidPayablesCents": payables["paidAmountCents"],
        "netOperationalMarginCents": (
            receivables["paidAmountCents"]
            - governance["releasedCommissionsCents"]
            - governance["operatingCostCents"]
            - payables["paidAmountCents"]
        ),
        "netTreasuryMovementCents": treasury["inflowCents"] - treasury["outflowCents"],
    }

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "treasury": treasury,
        "receivables": receivables,
        "payables": payables,
        "billing": billing,
        "profitability": profitability,
        "governance": governance,
    }


def fetch_treasury_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM finance.cash_accounts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS accounts_total,
            (
                SELECT count(*)
                FROM finance.cash_accounts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS active_accounts,
            (
                SELECT count(*)
                FROM finance.cash_movements
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS movements_total,
            (
                SELECT COALESCE(sum(opening_balance_cents), 0)
                FROM finance.cash_accounts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) +
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM finance.cash_movements
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND direction = 'inflow'
            ) -
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM finance.cash_movements
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND direction = 'outflow'
            ) AS current_balance_cents,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM finance.cash_movements
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND direction = 'inflow'
            ) AS inflow_cents,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM finance.cash_movements
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND direction = 'outflow'
            ) AS outflow_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 8)
        row = cursor.fetchone() or {}

    return {
        "accountsTotal": int(row.get("accounts_total", 0) or 0),
        "activeAccounts": int(row.get("active_accounts", 0) or 0),
        "movementsTotal": int(row.get("movements_total", 0) or 0),
        "currentBalanceCents": int(row.get("current_balance_cents", 0) or 0),
        "inflowCents": int(row.get("inflow_cents", 0) or 0),
        "outflowCents": int(row.get("outflow_cents", 0) or 0),
    }


def fetch_receivable_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE receivable.status = 'open') AS open,
            count(*) FILTER (WHERE receivable.status = 'paid') AS paid,
            count(*) FILTER (WHERE receivable.status = 'cancelled') AS cancelled,
            COALESCE(sum(receivable.amount_cents) FILTER (WHERE receivable.status = 'open'), 0) AS open_amount_cents,
            COALESCE(sum(receivable.amount_cents) FILTER (WHERE receivable.status = 'paid'), 0) AS paid_amount_cents,
            count(*) FILTER (
                WHERE receivable.status = 'open'
                  AND receivable.due_date < timezone('utc', now())::date
            ) AS overdue_count
        FROM finance.receivable_entries AS receivable
        JOIN identity.tenants AS tenant ON tenant.id = receivable.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "open": int(row.get("open", 0) or 0),
        "paid": int(row.get("paid", 0) or 0),
        "cancelled": int(row.get("cancelled", 0) or 0),
        "openAmountCents": int(row.get("open_amount_cents", 0) or 0),
        "paidAmountCents": int(row.get("paid_amount_cents", 0) or 0),
        "overdueCount": int(row.get("overdue_count", 0) or 0),
    }


def fetch_payable_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE payable.status = 'open') AS open,
            count(*) FILTER (WHERE payable.status = 'paid') AS paid,
            count(*) FILTER (WHERE payable.status = 'cancelled') AS cancelled,
            COALESCE(sum(payable.amount_cents) FILTER (WHERE payable.status = 'open'), 0) AS open_amount_cents,
            COALESCE(sum(payable.amount_cents) FILTER (WHERE payable.status = 'paid'), 0) AS paid_amount_cents,
            count(*) FILTER (
                WHERE payable.status = 'open'
                  AND payable.due_date >= timezone('utc', now())::date
                  AND payable.due_date <= timezone('utc', now())::date + INTERVAL '7 days'
            ) AS due_soon_count
        FROM finance.payables AS payable
        JOIN identity.tenants AS tenant ON tenant.id = payable.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "open": int(row.get("open", 0) or 0),
        "paid": int(row.get("paid", 0) or 0),
        "cancelled": int(row.get("cancelled", 0) or 0),
        "openAmountCents": int(row.get("open_amount_cents", 0) or 0),
        "paidAmountCents": int(row.get("paid_amount_cents", 0) or 0),
        "dueSoonCount": int(row.get("due_soon_count", 0) or 0),
    }


def fetch_billing_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (SELECT count(*) FROM billing.plans) AS plans_total,
            (SELECT count(*) FROM billing.plans WHERE active) AS active_plans,
            (
                SELECT count(*)
                FROM billing.subscriptions
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS active_subscriptions,
            (
                SELECT count(*)
                FROM billing.subscriptions
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'grace_period'
            ) AS grace_subscriptions,
            (
                SELECT count(*)
                FROM billing.subscriptions
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'suspended'
            ) AS suspended_subscriptions,
            (
                SELECT count(*)
                FROM billing.subscription_invoices
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'open'
            ) AS invoices_open,
            (
                SELECT count(*)
                FROM billing.subscription_invoices
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'paid'
            ) AS invoices_paid,
            (
                SELECT count(*)
                FROM billing.subscription_invoices
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'failed'
            ) AS invoices_failed,
            (
                SELECT count(*)
                FROM billing.payment_attempts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'failed'
            ) AS failed_attempts,
            (
                SELECT count(*)
                FROM billing.payment_attempts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'succeeded'
            ) AS succeeded_attempts,
            (
                SELECT COALESCE(sum(plan.amount_cents), 0)
                FROM billing.subscriptions AS subscription
                JOIN billing.plans AS plan ON plan.id = subscription.plan_id
                WHERE subscription.tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND subscription.status = 'active'
            ) AS monthly_recurring_revenue_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 9)
        row = cursor.fetchone() or {}

    return {
        "plansTotal": int(row.get("plans_total", 0) or 0),
        "activePlans": int(row.get("active_plans", 0) or 0),
        "activeSubscriptions": int(row.get("active_subscriptions", 0) or 0),
        "graceSubscriptions": int(row.get("grace_subscriptions", 0) or 0),
        "suspendedSubscriptions": int(row.get("suspended_subscriptions", 0) or 0),
        "invoicesOpen": int(row.get("invoices_open", 0) or 0),
        "invoicesPaid": int(row.get("invoices_paid", 0) or 0),
        "invoicesFailed": int(row.get("invoices_failed", 0) or 0),
        "failedAttempts": int(row.get("failed_attempts", 0) or 0),
        "succeededAttempts": int(row.get("succeeded_attempts", 0) or 0),
        "monthlyRecurringRevenueCents": int(row.get("monthly_recurring_revenue_cents", 0) or 0),
    }


def fetch_governance_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM finance.period_closures
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS period_closures,
            (
                SELECT count(*)
                FROM billing.subscription_events
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS subscription_events,
            (
                SELECT count(*)
                FROM finance.cash_movements
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS cash_movements,
            (
                SELECT count(*)
                FROM billing.payment_attempts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'failed'
            ) AS failed_payment_attempts,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM finance.commission_entries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'released'
            ) AS released_commissions_cents,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM finance.cost_entries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS operating_cost_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 6)
        row = cursor.fetchone() or {}

    return {
        "periodClosures": int(row.get("period_closures", 0) or 0),
        "subscriptionEvents": int(row.get("subscription_events", 0) or 0),
        "cashMovements": int(row.get("cash_movements", 0) or 0),
        "failedPaymentAttempts": int(row.get("failed_payment_attempts", 0) or 0),
        "releasedCommissionsCents": int(row.get("released_commissions_cents", 0) or 0),
        "operatingCostCents": int(row.get("operating_cost_cents", 0) or 0),
    }
