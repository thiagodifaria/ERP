"""Relatorio consolidado da operacao de faturamento e cobranca."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_revenue_operations(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_revenue_operations(tenant_slug)

    return build_static_revenue_operations(tenant_slug)


def build_static_revenue_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"
    booked_revenue_cents = 1775000
    paid_amount_cents = 845000

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "sales": {
            "total": 12,
            "active": 8,
            "invoiced": 4,
            "bookedRevenueCents": booked_revenue_cents,
        },
        "invoices": {
            "total": 9,
            "openAmountCents": 930000,
            "paidAmountCents": paid_amount_cents,
            "overdueAmountCents": 125000,
            "overdueCount": 1,
            "byStatus": {
                "draft": 1,
                "sent": 4,
                "paid": 4,
                "cancelled": 0,
            },
        },
        "collections": {
            "invoiceCoverageRate": 0.75,
            "collectionRate": round(paid_amount_cents / booked_revenue_cents, 4),
            "averageTicketCents": 197222,
        },
        "risk": {
            "invoicesDueSoon": 2,
            "overdueInvoices": 1,
            "overdueAmountCents": 125000,
        },
    }


def build_postgres_revenue_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        sales = fetch_sales_metrics(connection, tenant_slug)
        invoices = fetch_invoice_metrics(connection, tenant_slug)
        due_soon = fetch_due_soon_count(connection, tenant_slug)

    invoice_coverage_rate = round(invoices["total"] / sales["total"], 4) if sales["total"] > 0 else 0.0
    collection_rate = (
        round(invoices["paidAmountCents"] / sales["bookedRevenueCents"], 4) if sales["bookedRevenueCents"] > 0 else 0.0
    )
    average_ticket_cents = int(sales["bookedRevenueCents"] / sales["total"]) if sales["total"] > 0 else 0

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "sales": sales,
        "invoices": invoices,
        "collections": {
            "invoiceCoverageRate": invoice_coverage_rate,
            "collectionRate": collection_rate,
            "averageTicketCents": average_ticket_cents,
        },
        "risk": {
            "invoicesDueSoon": due_soon,
            "overdueInvoices": invoices["overdueCount"],
            "overdueAmountCents": invoices["overdueAmountCents"],
        },
    }


def fetch_sales_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE sale.status = 'active') AS active,
            count(*) FILTER (WHERE sale.status = 'invoiced') AS invoiced,
            COALESCE(sum(sale.amount_cents) FILTER (WHERE sale.status <> 'cancelled'), 0) AS booked_revenue_cents
        FROM sales.sales AS sale
        JOIN identity.tenants AS tenant ON tenant.id = sale.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "active": int(row.get("active", 0) or 0),
        "invoiced": int(row.get("invoiced", 0) or 0),
        "bookedRevenueCents": int(row.get("booked_revenue_cents", 0) or 0),
    }


def fetch_invoice_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            COALESCE(sum(invoice.amount_cents) FILTER (WHERE invoice.status NOT IN ('paid', 'cancelled')), 0) AS open_amount_cents,
            COALESCE(sum(invoice.amount_cents) FILTER (WHERE invoice.status = 'paid'), 0) AS paid_amount_cents,
            COALESCE(sum(invoice.amount_cents) FILTER (
                WHERE invoice.status NOT IN ('paid', 'cancelled')
                  AND invoice.due_date < timezone('utc', now())::date
            ), 0) AS overdue_amount_cents,
            count(*) FILTER (
                WHERE invoice.status NOT IN ('paid', 'cancelled')
                  AND invoice.due_date < timezone('utc', now())::date
            ) AS overdue_count,
            count(*) FILTER (WHERE invoice.status = 'draft') AS draft,
            count(*) FILTER (WHERE invoice.status = 'sent') AS sent,
            count(*) FILTER (WHERE invoice.status = 'paid') AS paid,
            count(*) FILTER (WHERE invoice.status = 'cancelled') AS cancelled
        FROM sales.invoices AS invoice
        JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "openAmountCents": int(row.get("open_amount_cents", 0) or 0),
        "paidAmountCents": int(row.get("paid_amount_cents", 0) or 0),
        "overdueAmountCents": int(row.get("overdue_amount_cents", 0) or 0),
        "overdueCount": int(row.get("overdue_count", 0) or 0),
        "byStatus": {
            "draft": int(row.get("draft", 0) or 0),
            "sent": int(row.get("sent", 0) or 0),
            "paid": int(row.get("paid", 0) or 0),
            "cancelled": int(row.get("cancelled", 0) or 0),
        },
    }


def fetch_due_soon_count(connection, tenant_slug: str | None) -> int:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total
        FROM sales.invoices AS invoice
        JOIN identity.tenants AS tenant ON tenant.id = invoice.tenant_id
        WHERE invoice.status NOT IN ('paid', 'cancelled')
          AND invoice.due_date >= timezone('utc', now())::date
          AND invoice.due_date <= timezone('utc', now())::date + INTERVAL '7 days'
          {f"AND {filter_sql.removeprefix('WHERE ')}" if filter_sql else ""}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return int(row.get("total", 0) or 0)
