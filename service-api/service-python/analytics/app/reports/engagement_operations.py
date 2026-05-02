"""Relatorio operacional do contexto engagement."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_engagement_operations(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_engagement_operations(tenant_slug)

    return build_static_engagement_operations(tenant_slug)


def build_static_engagement_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "campaigns": {
            "total": 2,
            "active": 1,
            "paused": 1,
            "budgetCents": 130000,
        },
        "templates": {
            "total": 3,
            "active": 2,
            "draft": 1,
            "archived": 0,
        },
        "touchpoints": {
            "total": 18,
            "queued": 2,
            "sent": 3,
            "delivered": 5,
            "responded": 4,
            "converted": 3,
            "failed": 1,
            "workflowDispatched": 11,
        },
        "deliveries": {
            "total": 17,
            "delivered": 12,
            "failed": 2,
            "deliveryRate": 0.7059,
            "failureRate": 0.1176,
            "byProvider": {
                "resend": 6,
                "whatsapp_cloud": 8,
                "telegram_bot": 1,
                "manual": 2,
            },
            "byStatus": {
                "queued": 1,
                "sent": 2,
                "delivered": 12,
                "failed": 2,
            },
        },
        "providers": {
            "configured": 3,
            "fallbackEnabled": 1,
            "inboundEvents": 7,
            "processedEvents": 6,
            "failedEvents": 1,
            "inboundLeads": 3,
            "workflowDispatches": 4,
            "responsesTracked": 2,
        },
        "governance": {
            "templateLinked": 15,
            "activeProviders": 4,
            "providerCapabilitiesReady": 3,
        },
    }


def build_postgres_engagement_operations(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        campaigns = fetch_campaign_metrics(connection, tenant_slug)
        templates = fetch_template_metrics(connection, tenant_slug)
        touchpoints = fetch_touchpoint_metrics(connection, tenant_slug)
        deliveries = fetch_delivery_metrics(connection, tenant_slug)
        providers = fetch_provider_event_metrics(connection, tenant_slug)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "campaigns": campaigns,
        "templates": templates,
        "touchpoints": touchpoints,
        "deliveries": deliveries,
        "providers": providers,
        "governance": {
            "templateLinked": deliveries["templateLinked"],
            "activeProviders": deliveries["activeProviders"],
            "providerCapabilitiesReady": providers["configured"],
        },
    }


def fetch_campaign_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE campaign.status = 'active') AS active,
            count(*) FILTER (WHERE campaign.status = 'paused') AS paused,
            COALESCE(sum(campaign.budget_cents), 0) AS budget_cents
        FROM engagement.campaigns AS campaign
        JOIN identity.tenants AS tenant ON tenant.id = campaign.tenant_id
        {filter_sql}
    """
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}
    return {
        "total": int(row.get("total", 0) or 0),
        "active": int(row.get("active", 0) or 0),
        "paused": int(row.get("paused", 0) or 0),
        "budgetCents": int(row.get("budget_cents", 0) or 0),
    }


def fetch_template_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE template.status = 'active') AS active,
            count(*) FILTER (WHERE template.status = 'draft') AS draft,
            count(*) FILTER (WHERE template.status = 'archived') AS archived
        FROM engagement.templates AS template
        JOIN identity.tenants AS tenant ON tenant.id = template.tenant_id
        {filter_sql}
    """
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}
    return {
        "total": int(row.get("total", 0) or 0),
        "active": int(row.get("active", 0) or 0),
        "draft": int(row.get("draft", 0) or 0),
        "archived": int(row.get("archived", 0) or 0),
    }


def fetch_touchpoint_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE touchpoint.status = 'queued') AS queued,
            count(*) FILTER (WHERE touchpoint.status = 'sent') AS sent,
            count(*) FILTER (WHERE touchpoint.status = 'delivered') AS delivered,
            count(*) FILTER (WHERE touchpoint.status = 'responded') AS responded,
            count(*) FILTER (WHERE touchpoint.status = 'converted') AS converted,
            count(*) FILTER (WHERE touchpoint.status = 'failed') AS failed,
            count(*) FILTER (WHERE touchpoint.last_workflow_run_public_id IS NOT NULL) AS workflow_dispatched
        FROM engagement.touchpoints AS touchpoint
        JOIN identity.tenants AS tenant ON tenant.id = touchpoint.tenant_id
        {filter_sql}
    """
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}
    return {
        "total": int(row.get("total", 0) or 0),
        "queued": int(row.get("queued", 0) or 0),
        "sent": int(row.get("sent", 0) or 0),
        "delivered": int(row.get("delivered", 0) or 0),
        "responded": int(row.get("responded", 0) or 0),
        "converted": int(row.get("converted", 0) or 0),
        "failed": int(row.get("failed", 0) or 0),
        "workflowDispatched": int(row.get("workflow_dispatched", 0) or 0),
    }


def fetch_delivery_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE delivery.status = 'delivered') AS delivered,
            count(*) FILTER (WHERE delivery.status = 'failed') AS failed,
            count(*) FILTER (WHERE delivery.status = 'queued') AS queued,
            count(*) FILTER (WHERE delivery.status = 'sent') AS sent,
            count(*) FILTER (WHERE delivery.provider = 'resend') AS resend,
            count(*) FILTER (WHERE delivery.provider = 'whatsapp_cloud') AS whatsapp_cloud,
            count(*) FILTER (WHERE delivery.provider = 'telegram_bot') AS telegram_bot,
            count(*) FILTER (WHERE delivery.provider = 'manual') AS manual,
            count(*) FILTER (WHERE delivery.template_id IS NOT NULL) AS template_linked,
            count(DISTINCT delivery.provider) AS active_providers
        FROM engagement.touchpoint_deliveries AS delivery
        JOIN identity.tenants AS tenant ON tenant.id = delivery.tenant_id
        {filter_sql}
    """
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    total = int(row.get("total", 0) or 0)
    delivered = int(row.get("delivered", 0) or 0)
    failed = int(row.get("failed", 0) or 0)
    return {
        "total": total,
        "delivered": delivered,
        "failed": failed,
        "deliveryRate": round(delivered / total, 4) if total > 0 else 0,
        "failureRate": round(failed / total, 4) if total > 0 else 0,
        "byProvider": {
            "resend": int(row.get("resend", 0) or 0),
            "whatsapp_cloud": int(row.get("whatsapp_cloud", 0) or 0),
            "telegram_bot": int(row.get("telegram_bot", 0) or 0),
            "manual": int(row.get("manual", 0) or 0),
        },
        "byStatus": {
            "queued": int(row.get("queued", 0) or 0),
            "sent": int(row.get("sent", 0) or 0),
            "delivered": delivered,
            "failed": failed,
        },
        "templateLinked": int(row.get("template_linked", 0) or 0),
        "activeProviders": int(row.get("active_providers", 0) or 0),
    }


def fetch_provider_event_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    query = f"""
        SELECT
            count(*) AS total,
            count(*) FILTER (WHERE event.direction = 'inbound') AS inbound_events,
            count(*) FILTER (WHERE event.direction = 'outbound') AS outbound_events,
            count(*) FILTER (WHERE event.status = 'processed') AS processed_events,
            count(*) FILTER (WHERE event.status = 'failed') AS failed_events,
            count(*) FILTER (WHERE event.event_type = 'lead.ingested') AS inbound_leads,
            count(*) FILTER (WHERE event.event_type = 'workflow.dispatched') AS workflow_dispatches,
            count(*) FILTER (WHERE event.event_type = 'delivery.responded') AS responses_tracked,
            count(*) FILTER (WHERE event.provider = 'resend') AS resend_configured,
            count(*) FILTER (WHERE event.provider = 'whatsapp_cloud') AS whatsapp_configured,
            count(*) FILTER (WHERE event.provider = 'telegram_bot') AS telegram_configured,
            count(*) FILTER (WHERE event.provider = 'manual') AS fallback_events
        FROM engagement.provider_events AS event
        JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
        {filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    configured = 0
    for key in ("resend_configured", "whatsapp_configured", "telegram_configured"):
        if int(row.get(key, 0) or 0) > 0:
            configured += 1

    return {
        "total": int(row.get("total", 0) or 0),
        "inboundEvents": int(row.get("inbound_events", 0) or 0),
        "outboundEvents": int(row.get("outbound_events", 0) or 0),
        "processedEvents": int(row.get("processed_events", 0) or 0),
        "failedEvents": int(row.get("failed_events", 0) or 0),
        "inboundLeads": int(row.get("inbound_leads", 0) or 0),
        "workflowDispatches": int(row.get("workflow_dispatches", 0) or 0),
        "responsesTracked": int(row.get("responses_tracked", 0) or 0),
        "configured": configured,
        "fallbackEnabled": 1 if int(row.get("fallback_events", 0) or 0) > 0 else 0,
    }
