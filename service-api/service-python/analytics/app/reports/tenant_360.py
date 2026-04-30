"""Relatorio consolidado do tenant cruzando identidade, comercial e automacao."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import append_filter, tenant_filter


def build_tenant_360(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_tenant_360(tenant_slug)

    return build_static_tenant_360(tenant_slug)


def build_static_tenant_360(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "identity": {
            "companies": 3,
            "users": 18,
            "teams": 6,
            "roles": 5,
            "teamMemberships": 14,
        },
        "commercial": {
            "leads": 128,
            "qualifiedLeads": 37,
            "assignedLeads": 96,
            "leadNotes": 52,
            "opportunities": 34,
            "proposals": 21,
            "sales": 12,
            "bookedRevenueCents": 1775000,
        },
        "engagement": {
            "campaigns": 2,
            "activeCampaigns": 1,
            "templates": 3,
            "touchpoints": 18,
            "deliveries": 17,
            "deliveredDeliveries": 12,
            "convertedTouchpoints": 3,
            "failedDeliveries": 2,
        },
        "rentals": {
            "contracts": 12,
            "activeContracts": 10,
            "scheduledCharges": 18,
            "paidCharges": 11,
            "attachments": 7,
            "outstandingAmountCents": 3015000,
        },
        "automation": {
            "activeDefinitions": 6,
            "workflowRuns": 41,
            "workflowRunEvents": 84,
            "runtimeExecutions": 44,
            "runtimeCompleted": 28,
            "runtimeFailed": 8,
            "runtimeCancelled": 4,
        },
    }


def build_postgres_tenant_360(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        identity_metrics = fetch_identity_metrics(connection, tenant_slug)
        commercial_metrics = fetch_commercial_metrics(connection, tenant_slug)
        engagement_metrics = fetch_engagement_metrics(connection, tenant_slug)
        rentals_metrics = fetch_rentals_metrics(connection, tenant_slug)
        automation_metrics = fetch_automation_metrics(connection, tenant_slug)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "identity": identity_metrics,
        "commercial": commercial_metrics,
        "engagement": engagement_metrics,
        "rentals": rentals_metrics,
        "automation": automation_metrics,
    }


def fetch_identity_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM identity.companies
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS companies,
            (
                SELECT count(*)
                FROM identity.users
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS users,
            (
                SELECT count(*)
                FROM identity.teams
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS teams,
            (
                SELECT count(*)
                FROM identity.roles
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS roles,
            (
                SELECT count(*)
                FROM identity.team_memberships
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS team_memberships
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 5)
        row = cursor.fetchone() or {}

    return {
        "companies": int(row.get("companies", 0) or 0),
        "users": int(row.get("users", 0) or 0),
        "teams": int(row.get("teams", 0) or 0),
        "roles": int(row.get("roles", 0) or 0),
        "teamMemberships": int(row.get("team_memberships", 0) or 0),
    }


def fetch_commercial_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    tenant_subquery_filter_sql, tenant_subquery_params = tenant_filter("slug = %s", tenant_slug)

    leads_query = f"""
        SELECT
            count(*) AS leads,
            count(*) FILTER (WHERE lead.status = 'qualified') AS qualified_leads,
            count(*) FILTER (WHERE lead.owner_user_public_id IS NOT NULL) AS assigned_leads
        FROM crm.leads AS lead
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    notes_query = f"""
        SELECT count(*) AS lead_notes
        FROM crm.lead_notes AS note
        JOIN crm.leads AS lead ON lead.id = note.lead_id
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    sales_query = f"""
        SELECT
            (
                SELECT count(*)
                FROM sales.opportunities
                WHERE tenant_id IN (SELECT id FROM identity.tenants {tenant_subquery_filter_sql})
            ) AS opportunities,
            (
                SELECT count(*)
                FROM sales.proposals
                WHERE tenant_id IN (SELECT id FROM identity.tenants {tenant_subquery_filter_sql})
            ) AS proposals,
            (
                SELECT count(*)
                FROM sales.sales
                WHERE tenant_id IN (SELECT id FROM identity.tenants {tenant_subquery_filter_sql})
            ) AS sales,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM sales.sales
                WHERE tenant_id IN (SELECT id FROM identity.tenants {tenant_subquery_filter_sql})
                  AND status <> 'cancelled'
            ) AS booked_revenue_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(leads_query, params)
        leads_row = cursor.fetchone() or {}
        cursor.execute(notes_query, params)
        notes_row = cursor.fetchone() or {}
        cursor.execute(sales_query, tenant_subquery_params * 4)
        sales_row = cursor.fetchone() or {}

    return {
        "leads": int(leads_row.get("leads", 0) or 0),
        "qualifiedLeads": int(leads_row.get("qualified_leads", 0) or 0),
        "assignedLeads": int(leads_row.get("assigned_leads", 0) or 0),
        "leadNotes": int(notes_row.get("lead_notes", 0) or 0),
        "opportunities": int(sales_row.get("opportunities", 0) or 0),
        "proposals": int(sales_row.get("proposals", 0) or 0),
        "sales": int(sales_row.get("sales", 0) or 0),
        "bookedRevenueCents": int(sales_row.get("booked_revenue_cents", 0) or 0),
    }


def fetch_engagement_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM engagement.campaigns
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS campaigns,
            (
                SELECT count(*)
                FROM engagement.campaigns
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS active_campaigns,
            (
                SELECT count(*)
                FROM engagement.templates
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS templates,
            (
                SELECT count(*)
                FROM engagement.touchpoints
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS touchpoints,
            (
                SELECT count(*)
                FROM engagement.touchpoint_deliveries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS deliveries,
            (
                SELECT count(*)
                FROM engagement.touchpoint_deliveries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'delivered'
            ) AS delivered_deliveries,
            (
                SELECT count(*)
                FROM engagement.touchpoints
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'converted'
            ) AS converted_touchpoints,
            (
                SELECT count(*)
                FROM engagement.touchpoint_deliveries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'failed'
            ) AS failed_deliveries
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 8)
        row = cursor.fetchone() or {}

    return {
        "campaigns": int(row.get("campaigns", 0) or 0),
        "activeCampaigns": int(row.get("active_campaigns", 0) or 0),
        "templates": int(row.get("templates", 0) or 0),
        "touchpoints": int(row.get("touchpoints", 0) or 0),
        "deliveries": int(row.get("deliveries", 0) or 0),
        "deliveredDeliveries": int(row.get("delivered_deliveries", 0) or 0),
        "convertedTouchpoints": int(row.get("converted_touchpoints", 0) or 0),
        "failedDeliveries": int(row.get("failed_deliveries", 0) or 0),
    }


def fetch_automation_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    definitions_query = f"""
        SELECT count(*) AS active_definitions
        FROM workflow_control.workflow_definitions AS definition
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        WHERE definition.status = 'active'
        {append_filter(filter_sql)}
    """

    runs_query = f"""
        SELECT
            count(*) AS workflow_runs,
            count(*) FILTER (WHERE run.status = 'completed') AS runs_completed,
            count(*) FILTER (WHERE run.status = 'failed') AS runs_failed,
            count(*) FILTER (WHERE run.status = 'cancelled') AS runs_cancelled
        FROM workflow_control.workflow_runs AS run
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    events_query = f"""
        SELECT count(*) AS workflow_run_events
        FROM workflow_control.workflow_run_events AS event
        JOIN workflow_control.workflow_runs AS run ON run.id = event.workflow_run_id
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    runtime_query = f"""
        SELECT
            count(*) AS runtime_executions,
            count(*) FILTER (WHERE execution.status = 'completed') AS runtime_completed,
            count(*) FILTER (WHERE execution.status = 'failed') AS runtime_failed,
            count(*) FILTER (WHERE execution.status = 'cancelled') AS runtime_cancelled
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(definitions_query, params)
        definitions_row = cursor.fetchone() or {}
        cursor.execute(runs_query, params)
        runs_row = cursor.fetchone() or {}
        cursor.execute(events_query, params)
        events_row = cursor.fetchone() or {}
        cursor.execute(runtime_query, params)
        runtime_row = cursor.fetchone() or {}

    return {
        "activeDefinitions": int(definitions_row.get("active_definitions", 0) or 0),
        "workflowRuns": int(runs_row.get("workflow_runs", 0) or 0),
        "workflowRunEvents": int(events_row.get("workflow_run_events", 0) or 0),
        "runtimeExecutions": int(runtime_row.get("runtime_executions", 0) or 0),
        "runtimeCompleted": int(runtime_row.get("runtime_completed", 0) or 0),
        "runtimeFailed": int(runtime_row.get("runtime_failed", 0) or 0),
        "runtimeCancelled": int(runtime_row.get("runtime_cancelled", 0) or 0),
    }


def fetch_rentals_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS contracts,
            (
                SELECT count(*)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS active_contracts,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
            ) AS scheduled_charges,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'paid'
            ) AS paid_charges,
            (
                SELECT count(*)
                FROM documents.attachments
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND owner_type = 'rentals.contract'
            ) AS attachments,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
            ) AS outstanding_amount_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 6)
        row = cursor.fetchone() or {}

    return {
        "contracts": int(row.get("contracts", 0) or 0),
        "activeContracts": int(row.get("active_contracts", 0) or 0),
        "scheduledCharges": int(row.get("scheduled_charges", 0) or 0),
        "paidCharges": int(row.get("paid_charges", 0) or 0),
        "attachments": int(row.get("attachments", 0) or 0),
        "outstandingAmountCents": int(row.get("outstanding_amount_cents", 0) or 0),
    }
