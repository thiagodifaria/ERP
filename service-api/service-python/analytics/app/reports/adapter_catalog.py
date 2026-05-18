"""Catalogo executivo de adapters externos e contracts compartilhados."""

from datetime import datetime, timezone
from pathlib import Path

from app.config.settings import settings
from app.reports.contract_governance import (
    FALLBACK_DECISION_DOCS,
    FALLBACK_EVENT_SCHEMAS,
    FALLBACK_HTTP_SPECS,
    resolve_architecture_decision_docs,
)
from app.reports.repo_root import try_resolve_repo_root


def build_adapter_catalog() -> dict:
    engagement = build_engagement_capabilities()
    billing = build_billing_capabilities()
    documents = build_document_storage_capabilities()
    signing = build_document_signing_capabilities()
    enrichment = build_enrichment_capabilities()
    webhook_hub = build_webhook_hub_capabilities()
    ai_governance = build_ai_governance_capabilities()
    document_intelligence = build_document_intelligence_capabilities()
    fiscal_brazil = build_fiscal_brazil_capabilities()
    registry_enrichment = build_registry_enrichment_capabilities()
    market_data = build_market_data_capabilities()
    external_risk_feed = build_external_risk_feed_capabilities()
    contracts = build_contract_catalog()

    configured = (
        engagement["summary"]["configured"]
        + billing["summary"]["configured"]
        + documents["summary"]["configured"]
        + signing["summary"]["configured"]
        + enrichment["summary"]["configured"]
        + ai_governance["summary"]["configured"]
        + document_intelligence["summary"]["configured"]
        + fiscal_brazil["summary"]["configured"]
        + registry_enrichment["summary"]["configured"]
        + market_data["summary"]["configured"]
        + external_risk_feed["summary"]["configured"]
    )
    fallback = (
        engagement["summary"]["fallback"]
        + billing["summary"]["fallback"]
        + documents["summary"]["fallback"]
        + signing["summary"]["fallback"]
        + enrichment["summary"]["fallback"]
        + ai_governance["summary"]["fallback"]
        + document_intelligence["summary"]["fallback"]
        + fiscal_brazil["summary"]["fallback"]
        + registry_enrichment["summary"]["fallback"]
        + market_data["summary"]["fallback"]
        + external_risk_feed["summary"]["fallback"]
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
        "aiGovernance": ai_governance,
        "documentIntelligence": document_intelligence,
        "fiscalBrazil": fiscal_brazil,
        "registryEnrichmentBrazil": registry_enrichment,
        "marketMacroRisk": market_data,
        "externalRiskFeed": external_risk_feed,
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


def build_ai_governance_capabilities() -> dict:
    capabilities = [
        capability(
            provider="openai",
            scope="llm_assistant",
            configured=bool(settings.openai_api_key.strip()),
            credential_key="OPENAI_API_KEY",
            fallback_viable=True,
            critical=False,
            model=settings.openai_model,
        )
    ]
    report = summarize_capabilities("ai-governance", capabilities)
    report["policy"] = "LLM externo so executa em modo BYOK; sem OPENAI_API_KEY o assistente permanece deterministico e somente leitura."
    return report


def build_document_intelligence_capabilities() -> dict:
    capabilities = [
        capability(
            provider="aws_textract",
            scope="ocr_document_intelligence",
            configured=bool(settings.aws_textract_access_key_id.strip() and settings.aws_textract_secret_access_key.strip() and settings.aws_textract_region.strip()),
            credential_key="AWS_TEXTRACT_ACCESS_KEY_ID",
            fallback_viable=False,
            critical=False,
            modeDetail="sdk_or_sigv4_required_for_runtime_extraction",
        ),
        capability(
            provider="google_document_ai",
            scope="ocr_document_intelligence",
            configured=bool(settings.google_document_ai_processor.strip() and settings.google_document_ai_credentials_json.strip()),
            credential_key="GOOGLE_DOCUMENT_AI_CREDENTIALS_JSON",
            fallback_viable=False,
            critical=False,
            modeDetail="oauth_service_account_required_for_runtime_extraction",
        ),
    ]
    return summarize_capabilities("document-intelligence", capabilities)


def build_fiscal_brazil_capabilities() -> dict:
    capabilities = [
        capability(
            provider="focus_nfe",
            scope="fiscal_issuance",
            configured=bool(settings.fiscal_focus_nfe_api_key.strip()),
            credential_key="FISCAL_FOCUS_NFE_API_KEY",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="enotas",
            scope="fiscal_issuance",
            configured=bool(settings.fiscal_enotas_api_key.strip()),
            credential_key="FISCAL_ENOTAS_API_KEY",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="certificate_a1",
            scope="digital_certificate",
            configured=bool(settings.fiscal_certificate_a1_secret.strip()),
            credential_key="FISCAL_CERTIFICATE_A1_SECRET",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="certificate_a3",
            scope="digital_certificate",
            configured=bool(settings.fiscal_certificate_a3_provider.strip()),
            credential_key="FISCAL_CERTIFICATE_A3_PROVIDER",
            fallback_viable=False,
            critical=False,
        ),
    ]
    return summarize_capabilities("fiscal-brazil", capabilities)


def build_registry_enrichment_capabilities() -> dict:
    capabilities = [
        capability(
            provider="serpro_cnpj",
            scope="cnpj_enrichment",
            configured=bool(settings.crm_serpro_client_id.strip() and settings.crm_serpro_client_secret.strip()),
            credential_key="CRM_SERPRO_CLIENT_SECRET",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="brasilapi",
            scope="cnpj_enrichment",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
            publicApi=True,
        ),
        capability(
            provider="viacep",
            scope="cep_enrichment",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
            publicApi=True,
        ),
    ]
    return summarize_capabilities("registry-enrichment-brazil", capabilities)


def build_market_data_capabilities() -> dict:
    capabilities = [
        capability(
            provider="alpha_vantage",
            scope="market_data",
            configured=bool(settings.market_alpha_vantage_api_key.strip()),
            credential_key="MARKET_ALPHA_VANTAGE_API_KEY",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="fixer",
            scope="fx_rates",
            configured=bool(settings.market_fixer_api_key.strip()),
            credential_key="MARKET_FIXER_API_KEY",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="bcb_sgs",
            scope="macro_rates",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
            publicApi=True,
        ),
        capability(
            provider="bcb_ptax",
            scope="fx_reference",
            configured=True,
            credential_key=None,
            fallback_viable=True,
            critical=False,
            publicApi=True,
        ),
    ]
    return summarize_capabilities("market-macro-risk", capabilities)


def build_external_risk_feed_capabilities() -> dict:
    capabilities = [
        capability(
            provider="newsapi",
            scope="news_risk_feed",
            configured=bool(settings.newsapi_key.strip()),
            credential_key="NEWSAPI_KEY",
            fallback_viable=False,
            critical=False,
        ),
        capability(
            provider="gdelt",
            scope="news_risk_feed",
            configured=bool(settings.gdelt_base_url.strip()),
            credential_key=None,
            fallback_viable=True,
            critical=False,
            publicApi=True,
        ),
        capability(
            provider="alpha_vantage_news",
            scope="market_news_sentiment",
            configured=bool(settings.market_alpha_vantage_api_key.strip()),
            credential_key="MARKET_ALPHA_VANTAGE_API_KEY",
            fallback_viable=False,
            critical=False,
        ),
    ]
    return summarize_capabilities("external-risk-feed", capabilities)


def build_contract_catalog() -> dict:
    root = try_resolve_repo_root(Path(__file__))
    if root is not None:
        http_specs = [item.name for item in sorted((root / "docs" / "contracts" / "http").glob("*.yaml"))]
        event_schemas = [item.name for item in sorted((root / "docs" / "contracts" / "events").glob("*.json"))]
        decision_docs = resolve_architecture_decision_docs(root)
        api_portal_ready = (root / "docs" / "contracts" / "portal" / "index.html").exists()
        registry_ready = (root / "docs" / "contracts" / "registry.json").exists() and (root / "docs" / "contracts" / "schema-registry.json").exists()
    else:
        http_specs = FALLBACK_HTTP_SPECS
        event_schemas = FALLBACK_EVENT_SCHEMAS
        decision_docs = FALLBACK_DECISION_DOCS
        api_portal_ready = True
        registry_ready = True

    return {
        "service": "contracts",
        "summary": {
            "artifacts": len(http_specs) + len(event_schemas) + len(decision_docs) + int(api_portal_ready) + int(registry_ready),
            "httpSpecs": len(http_specs),
            "eventSchemas": len(event_schemas),
            "adrRecorded": len(decision_docs) > 0,
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
    **extra: object,
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
    } | extra


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
