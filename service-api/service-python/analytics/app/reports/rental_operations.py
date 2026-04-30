"""Relatorio operacional do contexto de locacoes."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_rental_operations(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_rental_operations(tenant_slug)

    return build_static_rental_operations(tenant_slug)


def build_static_rental_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "contracts": {
            "total": 12,
            "active": 10,
            "terminated": 2,
            "monthlyRecurringAmountCents": 2450000,
        },
        "charges": {
            "scheduled": 18,
            "paid": 11,
            "cancelled": 3,
            "outstandingAmountCents": 3015000,
            "collectedAmountCents": 1665000,
            "cancelledAmountCents": 450000,
            "dueSoon": 2,
            "overdue": 1,
        },
        "governance": {
            "adjustments": 4,
            "historyEvents": 19,
            "pendingOutbox": 3,
            "attachments": 7,
        },
    }


def build_postgres_rental_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        contracts = fetch_contract_metrics(connection, tenant_slug)
        charges = fetch_charge_metrics(connection, tenant_slug)
        governance = fetch_governance_metrics(connection, tenant_slug)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "contracts": contracts,
        "charges": charges,
        "governance": governance,
    }


def fetch_contract_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS total,
            (
                SELECT count(*)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS active,
            (
                SELECT count(*)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'terminated'
            ) AS terminated,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS monthly_recurring_amount_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 4)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "active": int(row.get("active", 0) or 0),
        "terminated": int(row.get("terminated", 0) or 0),
        "monthlyRecurringAmountCents": int(row.get("monthly_recurring_amount_cents", 0) or 0),
    }


def fetch_charge_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
            ) AS scheduled,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'paid'
            ) AS paid,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'cancelled'
            ) AS cancelled,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
            ) AS outstanding_amount_cents,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'paid'
            ) AS collected_amount_cents,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'cancelled'
            ) AS cancelled_amount_cents,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
                  AND due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '15 days'
            ) AS due_soon,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
                  AND due_date < CURRENT_DATE
            ) AS overdue
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 8)
        row = cursor.fetchone() or {}

    return {
        "scheduled": int(row.get("scheduled", 0) or 0),
        "paid": int(row.get("paid", 0) or 0),
        "cancelled": int(row.get("cancelled", 0) or 0),
        "outstandingAmountCents": int(row.get("outstanding_amount_cents", 0) or 0),
        "collectedAmountCents": int(row.get("collected_amount_cents", 0) or 0),
        "cancelledAmountCents": int(row.get("cancelled_amount_cents", 0) or 0),
        "dueSoon": int(row.get("due_soon", 0) or 0),
        "overdue": int(row.get("overdue", 0) or 0),
    }


def fetch_governance_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM rentals.contract_adjustments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS adjustments,
            (
                SELECT count(*)
                FROM rentals.contract_events
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS history_events,
            (
                SELECT count(*)
                FROM rentals.outbox_events
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'pending'
            ) AS pending_outbox,
            (
                SELECT count(*)
                FROM documents.attachments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND owner_type = 'rentals.contract'
            ) AS attachments
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 4)
        row = cursor.fetchone() or {}

    return {
        "adjustments": int(row.get("adjustments", 0) or 0),
        "historyEvents": int(row.get("history_events", 0) or 0),
        "pendingOutbox": int(row.get("pending_outbox", 0) or 0),
        "attachments": int(row.get("attachments", 0) or 0),
    }
