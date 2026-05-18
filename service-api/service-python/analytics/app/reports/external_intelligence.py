"""Relatorios de inteligencia externa e verificacao operacional."""

from datetime import datetime, timezone

from app.reports.adapter_catalog import build_adapter_catalog


def build_external_intelligence_readiness(tenant_slug: str | None = None) -> dict:
    catalog = build_adapter_catalog()
    domains = {
        "documentIntelligence": catalog["documentIntelligence"],
        "fiscalBrazil": catalog["fiscalBrazil"],
        "registryEnrichmentBrazil": catalog["registryEnrichmentBrazil"],
        "marketMacroRisk": catalog["marketMacroRisk"],
        "externalRiskFeed": catalog["externalRiskFeed"],
    }
    configured = sum(domain["summary"]["configured"] for domain in domains.values())
    unavailable = sum(
        1
        for domain in domains.values()
        for capability in domain["capabilities"]
        if capability["status"] == "unconfigured"
    )
    return {
        "tenantSlug": tenant_slug or "global",
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "release": "1.4.6",
        "readiness": {
            "status": "ready",
            "configuredProviders": configured,
            "unavailableProviders": unavailable,
            "policy": "Providers externos sao BYOK ou public APIs explicitamente marcadas; sem credencial, a capability nao e apresentada como produtiva.",
        },
        "domains": domains,
    }


def build_document_intelligence_readiness(tenant_slug: str | None = None) -> dict:
    domain = build_adapter_catalog()["documentIntelligence"]
    return {
        "tenantSlug": tenant_slug or "global",
        "domain": "document-intelligence",
        "summary": domain["summary"],
        "providers": domain["capabilities"],
        "supportedFlows": [
            "invoice_ocr",
            "fiscal_document_extraction",
            "contract_metadata_extraction",
            "purchase_order_matching_evidence",
        ],
        "fallbackPolicy": "Sem provider OCR configurado, documentos seguem fluxo manual/auditavel; extracao automatica fica indisponivel.",
    }


def build_fiscal_brazil_readiness(tenant_slug: str | None = None) -> dict:
    domain = build_adapter_catalog()["fiscalBrazil"]
    return {
        "tenantSlug": tenant_slug or "global",
        "domain": "fiscal-brazil",
        "summary": domain["summary"],
        "providers": domain["capabilities"],
        "controls": {
            "issuanceProviderDeclared": any(item["provider"] in {"focus_nfe", "enotas"} and item["configured"] for item in domain["capabilities"]),
            "certificatePostureDeclared": any(item["scope"] == "digital_certificate" and item["configured"] for item in domain["capabilities"]),
            "homologationRequired": True,
        },
    }


def build_brazil_registry_enrichment(tenant_slug: str | None = None) -> dict:
    domain = build_adapter_catalog()["registryEnrichmentBrazil"]
    return {
        "tenantSlug": tenant_slug or "global",
        "domain": "registry-enrichment-brazil",
        "summary": domain["summary"],
        "providers": domain["capabilities"],
        "checks": [
            {"checkKey": "cnpj_lookup", "providers": ["serpro_cnpj", "brasilapi", "receita_ws", "conecta_gov"], "targetServices": ["crm", "supplier"]},
            {"checkKey": "cep_lookup", "providers": ["viacep"], "targetServices": ["crm", "supplier", "documents"]},
            {"checkKey": "master_data_quality", "providers": ["brasilapi", "viacep"], "targetServices": ["analytics"]},
        ],
    }


def build_market_macro_risk(tenant_slug: str | None = None) -> dict:
    domain = build_adapter_catalog()["marketMacroRisk"]
    return {
        "tenantSlug": tenant_slug or "global",
        "domain": "market-macro-risk",
        "summary": domain["summary"],
        "providers": domain["capabilities"],
        "signals": [
            {"signalKey": "fx_reference", "provider": "bcb_ptax", "usage": "referencia cambial para exposicao financeira"},
            {"signalKey": "macro_rate", "provider": "bcb_sgs", "usage": "juros e series macro para risco financeiro"},
            {"signalKey": "market_quote", "provider": "alpha_vantage", "usage": "sinal de mercado para watchlists"},
            {"signalKey": "fx_latest", "provider": "fixer", "usage": "cotacao multi-moeda quando chave existir"},
        ],
    }


def build_external_risk_feed(tenant_slug: str | None = None) -> dict:
    domain = build_adapter_catalog()["externalRiskFeed"]
    return {
        "tenantSlug": tenant_slug or "global",
        "domain": "external-risk-feed",
        "summary": domain["summary"],
        "providers": domain["capabilities"],
        "watchlists": [
            {"watchlistKey": "customer-reputation", "sources": ["newsapi", "gdelt"], "targetServices": ["crm", "support", "analytics"]},
            {"watchlistKey": "supplier-risk", "sources": ["newsapi", "gdelt"], "targetServices": ["supplier", "procurement", "analytics"]},
            {"watchlistKey": "market-news-sentiment", "sources": ["alpha_vantage_news"], "targetServices": ["finance", "banking", "analytics"]},
        ],
        "policy": "Noticias entram como sinal operacional, nao como verdade de dominio; qualquer alerta critico deve virar evidencia auditavel antes de afetar decisao.",
    }
