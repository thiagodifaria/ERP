"""Catalogo executivo de adapters externos e contracts compartilhados."""

from datetime import datetime, timezone

from app.config.settings import settings


def build_adapter_catalog() -> dict:
    engagement = build_engagement_capabilities()
    billing = build_billing_capabilities()
    documents = build_document_storage_capabilities()
    webhook_hub = build_webhook_hub_capabilities()
    contracts = build_contract_catalog()

    configured = engagement["summary"]["configured"] + billing["summary"]["configured"] + documents["summary"]["configured"]
    fallback = engagement["summary"]["fallback"] + billing["summary"]["fallback"] + documents["summary"]["fallback"]
    critical_unconfigured = billing["summary"]["criticalUnconfigured"] + webhook_hub["summary"]["criticalUnconfigured"]

    return {
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "summary": {
            "configuredCapabilities": configured,
            "fallbackCapabilities": fallback,
            "criticalUnconfiguredCapabilities": critical_unconfigured,
            "contractArtifacts": contracts["summary"]["artifacts"],
        },
        "engagement": engagement,
        "billing": billing,
        "documents": documents,
        "webhookHub": webhook_hub,
        "contracts": contracts,
    }


def build_engagement_capabilities() -> dict:
    providers = [
        capability(
            provider="resend",
            scope="email",
            configured=bool(settings.engagement_resend_api_key.strip()),
            credential_key="ENGAGEMENT_RESEND_API_KEY",
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="whatsapp_cloud",
            scope="messaging",
            configured=bool(settings.engagement_whatsapp_access_token.strip()),
            credential_key="ENGAGEMENT_WHATSAPP_ACCESS_TOKEN",
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="telegram_bot",
            scope="messaging",
            configured=bool(settings.engagement_telegram_bot_token.strip()),
            credential_key="ENGAGEMENT_TELEGRAM_BOT_TOKEN",
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="meta_ads",
            scope="ads",
            configured=bool(settings.engagement_meta_ads_access_token.strip()),
            credential_key="ENGAGEMENT_META_ADS_ACCESS_TOKEN",
            fallback_viable=True,
            critical=False,
        ),
    ]
    return summarize_capabilities("engagement", providers)


def build_billing_capabilities() -> dict:
    gateways = [
        capability(
            provider="asaas",
            scope="payments",
            configured=bool(settings.billing_asaas_api_key.strip()),
            credential_key="BILLING_ASAAS_API_KEY",
            fallback_viable=False,
            critical=True,
        ),
        capability(
            provider="stripe_pix",
            scope="payments",
            configured=bool(settings.billing_stripe_secret_key.strip()),
            credential_key="BILLING_STRIPE_SECRET_KEY",
            fallback_viable=False,
            critical=True,
        ),
        capability(
            provider="mercado_pago_pix",
            scope="payments",
            configured=bool(settings.billing_mercado_pago_access_token.strip()),
            credential_key="BILLING_MERCADO_PAGO_ACCESS_TOKEN",
            fallback_viable=False,
            critical=True,
        ),
    ]
    return summarize_capabilities("billing", gateways)


def build_document_storage_capabilities() -> dict:
    storage = [
        capability(
            provider="local",
            scope="storage",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="s3_compatible",
            scope="storage",
            configured=(
                settings.documents_storage_driver in {"s3", "s3_compatible"}
                or bool(settings.documents_storage_bucket.strip() and settings.documents_storage_endpoint.strip())
            ),
            credential_key="DOCUMENTS_STORAGE_BUCKET",
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="cloudflare_r2",
            scope="storage",
            configured=(
                settings.documents_storage_driver == "r2"
                or bool(settings.documents_r2_account_id.strip() and settings.documents_r2_bucket.strip())
            ),
            credential_key="DOCUMENTS_R2_ACCOUNT_ID",
            fallback_viable=True,
            critical=False,
        ),
    ]
    return summarize_capabilities("documents", storage)


def build_webhook_hub_capabilities() -> dict:
    outbound_signing_ready = bool(settings.webhook_hub_outbound_signing_secret.strip())

    return {
        "service": "webhook-hub",
        "summary": {
            "configured": int(outbound_signing_ready),
            "fallback": 0,
            "criticalUnconfigured": int(not outbound_signing_ready),
        },
        "controls": {
            "outboundWebhookSigningReady": outbound_signing_ready,
            "retriesReady": True,
            "dlqReady": True,
        },
    }


def build_contract_catalog() -> dict:
    return {
        "service": "contracts",
        "summary": {
            "artifacts": 11,
            "httpSpecs": 7,
            "eventSchemas": 4,
            "adrRecorded": True,
            "centralUiReady": True,
        }
    }


def capability(
    *,
    provider: str,
    scope: str,
    configured: bool,
    credential_key: str | None,
    fallback_viable: bool,
    critical: bool,
) -> dict:
    if configured:
        mode = "configured"
        status = "ready"
    elif fallback_viable:
        mode = "fallback"
        status = "fallback"
    else:
        mode = "unconfigured"
        status = "unconfigured"

    return {
        "provider": provider,
        "scope": scope,
        "configured": configured,
        "credentialKey": credential_key,
        "critical": critical,
        "fallbackViable": fallback_viable,
        "mode": mode,
        "status": status,
    }


def summarize_capabilities(service: str, capabilities: list[dict]) -> dict:
    configured = sum(1 for item in capabilities if item["configured"])
    fallback = sum(1 for item in capabilities if item["mode"] == "fallback")
    critical_unconfigured = sum(1 for item in capabilities if item["critical"] and item["status"] == "unconfigured")

    return {
        "service": service,
        "summary": {
            "configured": configured,
            "fallback": fallback,
            "criticalUnconfigured": critical_unconfigured,
        },
        "capabilities": capabilities,
    }
