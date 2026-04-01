"""Relatorio consolidado da jornada comercial de lead ate venda."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import append_filter, tenant_filter


def build_sales_journey(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_sales_journey(tenant_slug)

    return build_static_sales_journey(tenant_slug)


def build_static_sales_journey(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "funnel": {
            "leadsCaptured": 128,
            "leadsWithOpportunity": 34,
            "opportunitiesWithProposal": 21,
            "proposalsConverted": 12,
            "salesWon": 12,
            "leadToSaleConversionRate": 0.0938,
        },
        "opportunities": {
            "total": 34,
            "totalAmountCents": 4820000,
            "byStage": {"qualified": 7, "proposal": 9, "negotiation": 6, "won": 12, "lost": 0},
        },
        "proposals": {
            "total": 21,
            "totalAmountCents": 3110000,
            "byStatus": {"draft": 4, "sent": 5, "accepted": 12, "rejected": 0},
        },
        "sales": {
            "total": 12,
            "bookedRevenueCents": 1775000,
            "byStatus": {"active": 8, "invoiced": 4, "cancelled": 0},
        },
        "automation": {
            "controlRuns": 17,
            "controlCompleted": 9,
            "runtimeExecutions": 15,
            "runtimeCompleted": 11,
        },
    }


def build_postgres_sales_journey(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        funnel = fetch_funnel_metrics(connection, tenant_slug)
        opportunities = fetch_opportunity_metrics(connection, tenant_slug)
        proposals = fetch_proposal_metrics(connection, tenant_slug)
        sales = fetch_sales_metrics(connection, tenant_slug)
        automation = fetch_automation_metrics(connection, tenant_slug)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "funnel": funnel,
        "opportunities": opportunities,
        "proposals": proposals,
        "sales": sales,
        "automation": automation,
    }


def fetch_funnel_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM crm.leads
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS leads_captured,
            (
                SELECT count(DISTINCT lead_public_id)
                FROM sales.opportunities
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS leads_with_opportunity,
            (
                SELECT count(DISTINCT opportunity_id)
                FROM sales.proposals
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS opportunities_with_proposal,
            (
                SELECT count(DISTINCT proposal_id)
                FROM sales.sales
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS proposals_converted,
            (
                SELECT count(*)
                FROM sales.sales
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS sales_total
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 5)
        row = cursor.fetchone() or {}

    leads_captured = int(row.get("leads_captured", 0) or 0)
    sales_total = int(row.get("sales_total", 0) or 0)

    return {
        "leadsCaptured": leads_captured,
        "leadsWithOpportunity": int(row.get("leads_with_opportunity", 0) or 0),
        "opportunitiesWithProposal": int(row.get("opportunities_with_proposal", 0) or 0),
        "proposalsConverted": int(row.get("proposals_converted", 0) or 0),
        "salesWon": sales_total,
        "leadToSaleConversionRate": round(sales_total / leads_captured, 4) if leads_captured > 0 else 0.0,
    }


def fetch_opportunity_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            COALESCE(sum(opportunity.amount_cents), 0) AS total_amount_cents,
            count(*) FILTER (WHERE opportunity.stage = 'qualified') AS qualified,
            count(*) FILTER (WHERE opportunity.stage = 'proposal') AS proposal,
            count(*) FILTER (WHERE opportunity.stage = 'negotiation') AS negotiation,
            count(*) FILTER (WHERE opportunity.stage = 'won') AS won,
            count(*) FILTER (WHERE opportunity.stage = 'lost') AS lost
        FROM sales.opportunities AS opportunity
        JOIN identity.tenants AS tenant ON tenant.id = opportunity.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "totalAmountCents": int(row.get("total_amount_cents", 0) or 0),
        "byStage": {
            "qualified": int(row.get("qualified", 0) or 0),
            "proposal": int(row.get("proposal", 0) or 0),
            "negotiation": int(row.get("negotiation", 0) or 0),
            "won": int(row.get("won", 0) or 0),
            "lost": int(row.get("lost", 0) or 0),
        },
    }


def fetch_proposal_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            COALESCE(sum(proposal.amount_cents), 0) AS total_amount_cents,
            count(*) FILTER (WHERE proposal.status = 'draft') AS draft,
            count(*) FILTER (WHERE proposal.status = 'sent') AS sent,
            count(*) FILTER (WHERE proposal.status = 'accepted') AS accepted,
            count(*) FILTER (WHERE proposal.status = 'rejected') AS rejected
        FROM sales.proposals AS proposal
        JOIN identity.tenants AS tenant ON tenant.id = proposal.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "totalAmountCents": int(row.get("total_amount_cents", 0) or 0),
        "byStatus": {
            "draft": int(row.get("draft", 0) or 0),
            "sent": int(row.get("sent", 0) or 0),
            "accepted": int(row.get("accepted", 0) or 0),
            "rejected": int(row.get("rejected", 0) or 0),
        },
    }


def fetch_sales_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            COALESCE(sum(sale.amount_cents) FILTER (WHERE sale.status <> 'cancelled'), 0) AS booked_revenue_cents,
            count(*) FILTER (WHERE sale.status = 'active') AS active,
            count(*) FILTER (WHERE sale.status = 'invoiced') AS invoiced,
            count(*) FILTER (WHERE sale.status = 'cancelled') AS cancelled
        FROM sales.sales AS sale
        JOIN identity.tenants AS tenant ON tenant.id = sale.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "total": int(row.get("total", 0) or 0),
        "bookedRevenueCents": int(row.get("booked_revenue_cents", 0) or 0),
        "byStatus": {
            "active": int(row.get("active", 0) or 0),
            "invoiced": int(row.get("invoiced", 0) or 0),
            "cancelled": int(row.get("cancelled", 0) or 0),
        },
    }


def fetch_automation_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    control_query = f"""
        SELECT
            count(DISTINCT run.id) AS total,
            count(DISTINCT run.id) FILTER (WHERE run.status = 'completed') AS completed
        FROM workflow_control.workflow_runs AS run
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        JOIN sales.opportunities AS opportunity
          ON opportunity.tenant_id = run.tenant_id
         AND opportunity.lead_public_id = run.subject_public_id
        WHERE run.subject_type = 'crm.lead'
        {append_filter(filter_sql)}
    """

    runtime_query = f"""
        SELECT
            count(DISTINCT execution.id) AS total,
            count(DISTINCT execution.id) FILTER (WHERE execution.status = 'completed') AS completed
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        JOIN sales.opportunities AS opportunity
          ON opportunity.tenant_id = execution.tenant_id
         AND opportunity.lead_public_id = execution.subject_public_id
        WHERE execution.subject_type = 'crm.lead'
        {append_filter(filter_sql)}
    """

    with connection.cursor() as cursor:
        cursor.execute(control_query, params)
        control_row = cursor.fetchone() or {}
        cursor.execute(runtime_query, params)
        runtime_row = cursor.fetchone() or {}

    return {
        "controlRuns": int(control_row.get("total", 0) or 0),
        "controlCompleted": int(control_row.get("completed", 0) or 0),
        "runtimeExecutions": int(runtime_row.get("total", 0) or 0),
        "runtimeCompleted": int(runtime_row.get("completed", 0) or 0),
    }
