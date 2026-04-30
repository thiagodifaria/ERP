"""Pulso operacional transversal entre os principais serviços do ERP."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import append_filter, tenant_filter


def build_service_pulse(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_service_pulse(tenant_slug)

    return build_static_service_pulse(tenant_slug)


def build_static_service_pulse(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "services": {
            "crm": {"totalLeads": 128, "captured": 18, "contacted": 42, "qualified": 37, "disqualified": 31},
            "sales": {"opportunitiesTotal": 34, "proposalsTotal": 21, "salesTotal": 12, "bookedRevenueCents": 1775000},
            "finance": {"receivablesOpen": 4, "receivablesPaid": 9, "payablesOpen": 2, "cashAccounts": 2, "currentBalanceCents": 1246000, "periodClosures": 4},
            "billing": {"activeSubscriptions": 11, "graceSubscriptions": 1, "suspendedSubscriptions": 1, "invoicesOpen": 2, "invoicesPaid": 9, "failedAttempts": 3},
            "engagement": {"campaignsTotal": 2, "activeCampaigns": 1, "templatesTotal": 3, "deliveriesTotal": 17, "deliveredDeliveries": 12, "failedDeliveries": 2, "convertedTouchpoints": 3},
            "rentals": {"contractsTotal": 12, "activeContracts": 10, "scheduledCharges": 18, "paidCharges": 11, "cancelledCharges": 3, "overdueCharges": 1},
            "workflowControl": {"activeDefinitions": 6, "runsRunning": 7, "runsCompleted": 31, "runsFailed": 2, "runsCancelled": 1},
            "workflowRuntime": {"totalExecutions": 44, "running": 4, "completed": 28, "failed": 8, "cancelled": 4},
            "webhookHub": {"totalEvents": 93, "forwarded": 87, "queued": 2, "processing": 1, "failed": 3},
        },
    }


def build_postgres_service_pulse(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        crm_metrics = fetch_crm_metrics(connection, tenant_slug)
        sales_metrics = fetch_sales_service_metrics(connection, tenant_slug)
        finance_metrics = fetch_finance_service_metrics(connection, tenant_slug)
        billing_metrics = fetch_billing_service_metrics(connection, tenant_slug)
        engagement_metrics = fetch_engagement_service_metrics(connection, tenant_slug)
        rentals_metrics = fetch_rentals_service_metrics(connection, tenant_slug)
        workflow_control_metrics = fetch_workflow_control_metrics(connection, tenant_slug)
        workflow_runtime_metrics = fetch_workflow_runtime_metrics(connection, tenant_slug)
        webhook_hub_metrics = fetch_webhook_hub_metrics(connection)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "services": {
            "crm": crm_metrics,
            "sales": sales_metrics,
            "finance": finance_metrics,
            "billing": billing_metrics,
            "engagement": engagement_metrics,
            "rentals": rentals_metrics,
            "workflowControl": workflow_control_metrics,
            "workflowRuntime": workflow_runtime_metrics,
            "webhookHub": webhook_hub_metrics,
        },
    }


def fetch_crm_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total_leads,
            count(*) FILTER (WHERE lead.status = 'captured') AS captured,
            count(*) FILTER (WHERE lead.status = 'contacted') AS contacted,
            count(*) FILTER (WHERE lead.status = 'qualified') AS qualified,
            count(*) FILTER (WHERE lead.status = 'disqualified') AS disqualified
        FROM crm.leads AS lead
        JOIN identity.tenants AS tenant ON tenant.id = lead.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "totalLeads": int(row.get("total_leads", 0) or 0),
        "captured": int(row.get("captured", 0) or 0),
        "contacted": int(row.get("contacted", 0) or 0),
        "qualified": int(row.get("qualified", 0) or 0),
        "disqualified": int(row.get("disqualified", 0) or 0),
    }


def fetch_sales_service_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM sales.opportunities
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS opportunities_total,
            (
                SELECT count(*)
                FROM sales.proposals
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS proposals_total,
            (
                SELECT count(*)
                FROM sales.sales
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS sales_total,
            (
                SELECT COALESCE(sum(amount_cents), 0)
                FROM sales.sales
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status <> 'cancelled'
            ) AS booked_revenue_cents
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 4)
        row = cursor.fetchone() or {}

    return {
        "opportunitiesTotal": int(row.get("opportunities_total", 0) or 0),
        "proposalsTotal": int(row.get("proposals_total", 0) or 0),
        "salesTotal": int(row.get("sales_total", 0) or 0),
        "bookedRevenueCents": int(row.get("booked_revenue_cents", 0) or 0),
    }


def fetch_engagement_service_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM engagement.campaigns
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS campaigns_total,
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
            ) AS templates_total,
            (
                SELECT count(*)
                FROM engagement.touchpoint_deliveries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS deliveries_total,
            (
                SELECT count(*)
                FROM engagement.touchpoint_deliveries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'delivered'
            ) AS delivered_deliveries,
            (
                SELECT count(*)
                FROM engagement.touchpoint_deliveries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'failed'
            ) AS failed_deliveries,
            (
                SELECT count(*)
                FROM engagement.touchpoints
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'converted'
            ) AS converted_touchpoints
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 7)
        row = cursor.fetchone() or {}

    return {
        "campaignsTotal": int(row.get("campaigns_total", 0) or 0),
        "activeCampaigns": int(row.get("active_campaigns", 0) or 0),
        "templatesTotal": int(row.get("templates_total", 0) or 0),
        "deliveriesTotal": int(row.get("deliveries_total", 0) or 0),
        "deliveredDeliveries": int(row.get("delivered_deliveries", 0) or 0),
        "failedDeliveries": int(row.get("failed_deliveries", 0) or 0),
        "convertedTouchpoints": int(row.get("converted_touchpoints", 0) or 0),
    }


def fetch_finance_service_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM finance.receivable_entries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'open'
            ) AS receivables_open,
            (
                SELECT count(*)
                FROM finance.receivable_entries
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'paid'
            ) AS receivables_paid,
            (
                SELECT count(*)
                FROM finance.payables
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'open'
            ) AS payables_open,
            (
                SELECT count(*)
                FROM finance.cash_accounts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'active'
            ) AS cash_accounts,
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
                SELECT count(*)
                FROM finance.period_closures
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS period_closures
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 8)
        row = cursor.fetchone() or {}

    return {
        "receivablesOpen": int(row.get("receivables_open", 0) or 0),
        "receivablesPaid": int(row.get("receivables_paid", 0) or 0),
        "payablesOpen": int(row.get("payables_open", 0) or 0),
        "cashAccounts": int(row.get("cash_accounts", 0) or 0),
        "currentBalanceCents": int(row.get("current_balance_cents", 0) or 0),
        "periodClosures": int(row.get("period_closures", 0) or 0),
    }


def fetch_billing_service_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
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
                FROM billing.payment_attempts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'failed'
            ) AS failed_attempts
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 6)
        row = cursor.fetchone() or {}

    return {
        "activeSubscriptions": int(row.get("active_subscriptions", 0) or 0),
        "graceSubscriptions": int(row.get("grace_subscriptions", 0) or 0),
        "suspendedSubscriptions": int(row.get("suspended_subscriptions", 0) or 0),
        "invoicesOpen": int(row.get("invoices_open", 0) or 0),
        "invoicesPaid": int(row.get("invoices_paid", 0) or 0),
        "failedAttempts": int(row.get("failed_attempts", 0) or 0),
    }


def fetch_rentals_service_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("slug = %s", tenant_slug)

    query = f"""
        SELECT
            (
                SELECT count(*)
                FROM rentals.contracts
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
            ) AS contracts_total,
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
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'cancelled'
            ) AS cancelled_charges,
            (
                SELECT count(*)
                FROM rentals.contract_charges
                WHERE tenant_id IN (SELECT id FROM identity.tenants {filter_sql})
                  AND status = 'scheduled'
                  AND due_date < CURRENT_DATE
            ) AS overdue_charges
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params * 6)
        row = cursor.fetchone() or {}

    return {
        "contractsTotal": int(row.get("contracts_total", 0) or 0),
        "activeContracts": int(row.get("active_contracts", 0) or 0),
        "scheduledCharges": int(row.get("scheduled_charges", 0) or 0),
        "paidCharges": int(row.get("paid_charges", 0) or 0),
        "cancelledCharges": int(row.get("cancelled_charges", 0) or 0),
        "overdueCharges": int(row.get("overdue_charges", 0) or 0),
    }


def fetch_workflow_control_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    definitions_query = f"""
        SELECT count(*) AS total
        FROM workflow_control.workflow_definitions AS definition
        JOIN identity.tenants AS tenant ON tenant.id = definition.tenant_id
        WHERE definition.status = 'active'
        {append_filter(filter_sql)}
    """

    runs_query = f"""
        SELECT
            count(*) FILTER (WHERE run.status = 'running') AS running,
            count(*) FILTER (WHERE run.status = 'completed') AS completed,
            count(*) FILTER (WHERE run.status = 'failed') AS failed,
            count(*) FILTER (WHERE run.status = 'cancelled') AS cancelled
        FROM workflow_control.workflow_runs AS run
        JOIN identity.tenants AS tenant ON tenant.id = run.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(definitions_query, params)
        definitions_row = cursor.fetchone() or {"total": 0}
        cursor.execute(runs_query, params)
        runs_row = cursor.fetchone() or {}

    return {
        "activeDefinitions": int(definitions_row.get("total", 0) or 0),
        "runsRunning": int(runs_row.get("running", 0) or 0),
        "runsCompleted": int(runs_row.get("completed", 0) or 0),
        "runsFailed": int(runs_row.get("failed", 0) or 0),
        "runsCancelled": int(runs_row.get("cancelled", 0) or 0),
    }


def fetch_workflow_runtime_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)

    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE execution.status = 'running') AS running,
            count(*) FILTER (WHERE execution.status = 'completed') AS completed,
            count(*) FILTER (WHERE execution.status = 'failed') AS failed,
            count(*) FILTER (WHERE execution.status = 'cancelled') AS cancelled
        FROM workflow_runtime.executions AS execution
        JOIN identity.tenants AS tenant ON tenant.id = execution.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "totalExecutions": int(row.get("total", 0) or 0),
        "running": int(row.get("running", 0) or 0),
        "completed": int(row.get("completed", 0) or 0),
        "failed": int(row.get("failed", 0) or 0),
        "cancelled": int(row.get("cancelled", 0) or 0),
    }


def fetch_webhook_hub_metrics(connection) -> dict:
    with connection.cursor() as cursor:
        cursor.execute(
            """
                SELECT
                    count(*) AS total,
                    count(*) FILTER (WHERE status = 'forwarded') AS forwarded,
                    count(*) FILTER (WHERE status = 'queued') AS queued,
                    count(*) FILTER (WHERE status = 'processing') AS processing,
                    count(*) FILTER (WHERE status = 'failed') AS failed
                FROM webhook_hub.webhook_events
            """
        )
        row = cursor.fetchone() or {}

    return {
        "totalEvents": int(row.get("total", 0) or 0),
        "forwarded": int(row.get("forwarded", 0) or 0),
        "queued": int(row.get("queued", 0) or 0),
        "processing": int(row.get("processing", 0) or 0),
        "failed": int(row.get("failed", 0) or 0),
    }
