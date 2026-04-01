"""Relatorio inicial de pipeline comercial para o plano analitico."""

from datetime import datetime, timezone


def build_pipeline_summary(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "global"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "metrics": {
            "leadsCaptured": 128,
            "leadsQualified": 37,
            "conversions": 12,
            "conversionRate": 0.3243,
        },
        "bySource": {
            "whatsapp": 46,
            "meta_ads": 39,
            "referral": 23,
            "landing_page": 20,
        },
        "backlog": {
            "pendingContact": 18,
            "runningAutomations": 7,
            "awaitingFinancialReview": 3,
        },
    }
