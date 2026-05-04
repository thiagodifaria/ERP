"""Catalogo executivo de adapters externos e contracts compartilhados."""

from datetime import datetime, timezone
from pathlib import Path

from app.config.settings import settings
from app.reports.contract_governance import (
    FALLBACK_ADRS,
    FALLBACK_EVENT_SCHEMAS,
    FALLBACK_HTTP_SPECS,
)
from app.reports.repo_root import try_resolve_repo_root


def build_adapter_catalog() -> dict:
    engagement = build_engagement_capabilities()
    billing = build_billing_capabilities()
    documents = build_document_storage_capabilities()
    signing = build_document_signing_capabilities()
    enrichment = build_enrichment_capabilities()
    webhook_hub = build_webhook_hub_capabilities()
    contracts = build_contract_catalog()

    configured = (
        engagement["summary"]["configured"]
        + billing["summary"]["configured"]
        + documents["summary"]["configured"]
        + signing["summary"]["configured"]
        + enrichment["summary"]["configured"]
    )
    fallback = (
        engagement["summary"]["fallback"]
        + billing["summary"]["fallback"]
        + documents["summary"]["fallback"]
        + signing["summary"]["fallback"]
        + enrichment["summary"]["fallback"]
    )
    critical_unconfigured = (
        billing["summary"]["criticalUnconfigured"]
        + webhook_hub["summary"]["criticalUnconfigured"]
    )

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
        "documentSigning": signing,
        "crmEnrichment": enrichment,
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


def build_document_signing_capabilities() -> dict:
    capabilities = [
        capability(
            provider="local",
            scope="digital_signature",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="clicksign",
            scope="digital_signature",
            configured=bool(settings.documents_clicksign_api_key.strip()),
            credential_key="DOCUMENTS_CLICKSIGN_API_KEY",
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="docusign",
            scope="digital_signature",
            configured=bool(settings.documents_docusign_access_token.strip()),
            credential_key="DOCUMENTS_DOCUSIGN_ACCESS_TOKEN",
            fallback_viable=True,
            critical=False,
        ),
    ]
    return summarize_capabilities("document-signing", capabilities)


def build_enrichment_capabilities() -> dict:
    capabilities = [
        capability(
            provider="local",
            scope="cnpj_enrichment",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="receita_ws",
            scope="cnpj_enrichment",
            configured=bool(settings.crm_cnpj_provider_token.strip()),
            credential_key="CRM_CNPJ_PROVIDER_TOKEN",
            fallback_viable=True,
            critical=False,
        ),
        capability(
            provider="conecta_gov",
            scope="cnpj_enrichment",
            configured=bool(settings.crm_conecta_cnpj_api_key.strip()),
            credential_key="CRM_CONECTA_CNPJ_API_KEY",
            fallback_viable=True,
            critical=False,
        ),
    ]
    return summarize_capabilities("crm-enrichment", capabilities)


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
    root = try_resolve_repo_root(Path(__file__))
    if root is not None:
        http_specs = [item.name for item in sorted((root / "contracts" / "http").glob("*.yaml"))]
        event_schemas = [item.name for item in sorted((root / "contracts" / "events").glob("*.json"))]
        adr_docs = [item.name for item in sorted((root / "docs").glob("ADR-*.md"))]
        api_portal_ready = (root / "docs" / "API_PORTAL.md").exists()
        registry_ready = (root / "contracts" / "registry.json").exists()
    else:
        http_specs = FALLBACK_HTTP_SPECS
        event_schemas = FALLBACK_EVENT_SCHEMAS
        adr_docs = FALLBACK_ADRS
        api_portal_ready = True
        registry_ready = True

    return {
        "service": "contracts",
        "summary": {
            "artifacts": len(http_specs) + len(event_schemas) + len(adr_docs) + int(api_portal_ready) + int(registry_ready),
            "httpSpecs": len(http_specs),
            "eventSchemas": len(event_schemas),
            "adrRecorded": len(adr_docs) > 0,
            "centralUiReady": api_portal_ready,
            "schemaRegistryReady": registry_ready,
        },
        "httpSpecs": http_specs,
        "eventSchemas": event_schemas,
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
