"""Relatorio executivo de prontidao das integracoes externas."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect
from app.reports.pipeline_summary import tenant_filter


def build_integration_readiness(tenant_slug: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_integration_readiness(tenant_slug)

    return build_static_integration_readiness(tenant_slug)


def build_static_integration_readiness(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "providers": {
            "configured": 4,
            "activeInboundProviders": 2,
            "activeOutboundProviders": 3,
            "fallbackEnabled": True,
        },
        "flows": {
            "inboundLeads": 3,
            "workflowDispatches": 4,
            "responsesTracked": 2,
            "conversionsTracked": 1,
            "processedProviderEvents": 6,
            "failedProviderEvents": 1,
            "businessEntityLinkedEvents": 6,
        },
        "webhookHub": {
            "pendingEvents": 1,
            "deadLetterEvents": 1,
            "requeueReady": True,
        },
        "readiness": {
            "status": "attention",
            "leadIntakeReady": True,
            "workflowDispatchReady": True,
            "callbackTraceabilityReady": True,
            "businessEntityLinkageReady": True,
            "externalAdaptersPrepared": True,
            "openProviderRisks": 1,
        },
    }


def build_postgres_integration_readiness(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    with connect() as connection:
        provider_metrics = fetch_provider_metrics(connection, tenant_slug)
        webhook_metrics = fetch_webhook_metrics(connection)

    readiness_status = "stable"
    if webhook_metrics["deadLetterEvents"] > 0 or provider_metrics["failedProviderEvents"] > 0:
        readiness_status = "attention"
    if webhook_metrics["deadLetterEvents"] > 1 or provider_metrics["failedProviderEvents"] > 1:
        readiness_status = "critical"

    open_provider_risks = int(webhook_metrics["deadLetterEvents"] > 0) + int(provider_metrics["failedProviderEvents"] > 0)

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "providers": {
            "configured": provider_metrics["configuredProviders"],
            "activeInboundProviders": provider_metrics["activeInboundProviders"],
            "activeOutboundProviders": provider_metrics["activeOutboundProviders"],
            "fallbackEnabled": provider_metrics["fallbackEnabled"],
        },
        "flows": {
            "inboundLeads": provider_metrics["inboundLeads"],
            "workflowDispatches": provider_metrics["workflowDispatches"],
            "responsesTracked": provider_metrics["responsesTracked"],
            "conversionsTracked": provider_metrics["conversionsTracked"],
            "processedProviderEvents": provider_metrics["processedProviderEvents"],
            "failedProviderEvents": provider_metrics["failedProviderEvents"],
            "businessEntityLinkedEvents": provider_metrics["businessEntityLinkedEvents"],
        },
        "webhookHub": {
            "pendingEvents": webhook_metrics["pendingEvents"],
            "deadLetterEvents": webhook_metrics["deadLetterEvents"],
            "requeueReady": True,
        },
        "readiness": {
            "status": readiness_status,
            "leadIntakeReady": provider_metrics["inboundLeads"] > 0,
            "workflowDispatchReady": provider_metrics["workflowDispatches"] > 0,
            "callbackTraceabilityReady": (provider_metrics["responsesTracked"] + provider_metrics["conversionsTracked"]) > 0,
            "businessEntityLinkageReady": provider_metrics["businessEntityLinkedEvents"] > 0,
            "externalAdaptersPrepared": provider_metrics["configuredProviders"] >= 3,
            "openProviderRisks": open_provider_risks,
        },
    }


def fetch_provider_metrics(connection, tenant_slug: str | None) -> dict:
    filter_sql, params = tenant_filter("tenant.slug = %s", tenant_slug)
    query = f"""
        SELECT
            count(*) FILTER (WHERE event.provider IN ('resend', 'whatsapp_cloud', 'telegram_bot', 'meta_ads')) AS configured_providers_events,
            count(DISTINCT event.provider) FILTER (WHERE event.direction = 'inbound') AS active_inbound_providers,
            count(DISTINCT event.provider) FILTER (WHERE event.direction = 'outbound') AS active_outbound_providers,
            count(*) FILTER (WHERE event.provider = 'manual') AS fallback_events,
            count(*) FILTER (WHERE event.event_type = 'lead.ingested') AS inbound_leads,
            count(*) FILTER (WHERE event.event_type = 'workflow.dispatched') AS workflow_dispatches,
            count(*) FILTER (WHERE event.event_type = 'delivery.responded') AS responses_tracked,
            count(*) FILTER (WHERE event.event_type = 'delivery.converted') AS conversions_tracked,
            count(*) FILTER (WHERE event.status = 'processed') AS processed_provider_events,
            count(*) FILTER (WHERE event.status = 'failed') AS failed_provider_events,
            count(*) FILTER (WHERE event.business_entity_public_id IS NOT NULL) AS business_linked_events,
            count(*) FILTER (WHERE event.provider = 'resend') AS resend_events,
            count(*) FILTER (WHERE event.provider = 'whatsapp_cloud') AS whatsapp_events,
            count(*) FILTER (WHERE event.provider = 'telegram_bot') AS telegram_events,
            count(*) FILTER (WHERE event.provider = 'meta_ads') AS meta_ads_events
        FROM engagement.provider_events AS event
        JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
        {filter_sql}
    """
    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    configured_providers = 0
    for key in ("resend_events", "whatsapp_events", "telegram_events", "meta_ads_events"):
        if int(row.get(key, 0) or 0) > 0:
            configured_providers += 1

    return {
        "configuredProviders": configured_providers,
        "activeInboundProviders": int(row.get("active_inbound_providers", 0) or 0),
        "activeOutboundProviders": int(row.get("active_outbound_providers", 0) or 0),
        "fallbackEnabled": int(row.get("fallback_events", 0) or 0) > 0,
        "inboundLeads": int(row.get("inbound_leads", 0) or 0),
        "workflowDispatches": int(row.get("workflow_dispatches", 0) or 0),
        "responsesTracked": int(row.get("responses_tracked", 0) or 0),
        "conversionsTracked": int(row.get("conversions_tracked", 0) or 0),
        "processedProviderEvents": int(row.get("processed_provider_events", 0) or 0),
        "failedProviderEvents": int(row.get("failed_provider_events", 0) or 0),
        "businessEntityLinkedEvents": int(row.get("business_linked_events", 0) or 0),
    }


def fetch_webhook_metrics(connection) -> dict:
    with connection.cursor() as cursor:
        cursor.execute(
            """
                SELECT
                    count(*) FILTER (WHERE status IN ('validated', 'queued', 'processing')) AS pending_events,
                    count(*) FILTER (WHERE status = 'dead_letter') AS dead_letter_events
                FROM webhook_hub.webhook_events
            """
        )
        row = cursor.fetchone() or {}

    return {
        "pendingEvents": int(row.get("pending_events", 0) or 0),
        "deadLetterEvents": int(row.get("dead_letter_events", 0) or 0),
    }
