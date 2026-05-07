"""Leitura de maturidade comercial e operacional do relacionamento."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_relationship_intelligence(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": slug,
            "generatedAt": datetime.now(timezone.utc).isoformat(),
            "scoring": {"average": 68, "hot": 3, "warm": 4, "cold": 2},
            "pipeline": {"configs": 1, "stages": 5, "autoScoring": True, "territoryRules": 1, "approvalPolicies": 1},
            "support": {"openCases": 3, "overdueCases": 1, "slaTrackedCases": 3},
            "conversations": {"threads": 4, "participants": 4, "channels": 3},
            "bulkOperations": {"importsReady": True, "exportsReady": True, "partialSuccessTracking": True},
            "approvals": {"policies": 1, "auditReady": True},
            "territories": {"rules": 1, "assignmentReady": True},
            "forecast": {"bookedRevenueCents": 1775000, "weightedPipelineCents": 2314000, "confidence": "attention", "scenarioCount": 2},
            "readiness": {"status": "attention", "slaReady": True, "forecastReady": True, "relationshipSignalsReady": True, "bulkReady": True},
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  (SELECT count(*) FROM crm.leads AS lead JOIN identity.tenants AS tenant_lead ON tenant_lead.id = lead.tenant_id WHERE tenant_lead.slug = %s) AS leads_total,
                  (SELECT count(*) FROM crm.leads AS lead JOIN identity.tenants AS tenant_lead ON tenant_lead.id = lead.tenant_id WHERE tenant_lead.slug = %s AND lead.status = 'qualified') AS leads_qualified,
                  (SELECT count(*) FROM crm.pipeline_configs AS config JOIN identity.tenants AS tenant_cfg ON tenant_cfg.id = config.tenant_id WHERE tenant_cfg.slug = %s) AS pipeline_configs,
                  (SELECT coalesce(sum(jsonb_array_length(config.stages_json)), 0) FROM crm.pipeline_configs AS config JOIN identity.tenants AS tenant_cfg ON tenant_cfg.id = config.tenant_id WHERE tenant_cfg.slug = %s) AS pipeline_stages,
                  (SELECT count(*) FROM crm.pipeline_configs AS config JOIN identity.tenants AS tenant_cfg ON tenant_cfg.id = config.tenant_id WHERE tenant_cfg.slug = %s AND config.auto_scoring = TRUE) AS pipeline_auto_scoring,
                  (SELECT coalesce(sum(jsonb_array_length(config.territory_rules_json)), 0) FROM crm.pipeline_configs AS config JOIN identity.tenants AS tenant_cfg ON tenant_cfg.id = config.tenant_id WHERE tenant_cfg.slug = %s) AS pipeline_territory_rules,
                  (SELECT coalesce(sum(jsonb_array_length(config.approval_policies_json)), 0) FROM crm.pipeline_configs AS config JOIN identity.tenants AS tenant_cfg ON tenant_cfg.id = config.tenant_id WHERE tenant_cfg.slug = %s) AS pipeline_approval_policies,
                  (SELECT count(*) FROM support.cases AS support_case JOIN identity.tenants AS tenant_support ON tenant_support.id = support_case.tenant_id WHERE tenant_support.slug = %s AND support_case.status IN ('open', 'in_progress', 'waiting_customer')) AS support_open,
                  (SELECT count(*) FROM support.cases AS support_case JOIN identity.tenants AS tenant_support ON tenant_support.id = support_case.tenant_id WHERE tenant_support.slug = %s AND support_case.status NOT IN ('resolved', 'closed') AND support_case.sla_due_at < NOW()) AS support_overdue,
                  (SELECT count(DISTINCT touchpoint.thread_public_id) FROM engagement.touchpoints AS touchpoint JOIN identity.tenants AS tenant_touchpoint ON tenant_touchpoint.id = touchpoint.tenant_id WHERE tenant_touchpoint.slug = %s) AS conversation_threads,
                  (SELECT count(DISTINCT touchpoint.participant_kind || ':' || touchpoint.participant_public_id::text) FROM engagement.touchpoints AS touchpoint JOIN identity.tenants AS tenant_touchpoint ON tenant_touchpoint.id = touchpoint.tenant_id WHERE tenant_touchpoint.slug = %s) AS conversation_participants,
                  (SELECT count(DISTINCT touchpoint.channel) FROM engagement.touchpoints AS touchpoint JOIN identity.tenants AS tenant_touchpoint ON tenant_touchpoint.id = touchpoint.tenant_id WHERE tenant_touchpoint.slug = %s) AS conversation_channels,
                  (SELECT coalesce(sum(opportunity.amount_cents), 0) FROM sales.opportunities AS opportunity JOIN identity.tenants AS tenant_sales ON tenant_sales.id = opportunity.tenant_id WHERE tenant_sales.slug = %s AND opportunity.stage NOT IN ('lost')) AS weighted_pipeline_cents,
                  (SELECT coalesce(sum(sale.amount_cents), 0) FROM sales.sales AS sale JOIN identity.tenants AS tenant_sale ON tenant_sale.id = sale.tenant_id WHERE tenant_sale.slug = %s AND sale.status <> 'cancelled') AS booked_revenue_cents
                """,
                (slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug),
            )
            row = cursor.fetchone() or {}

    leads_total = int(row.get("leads_total", 0) or 0)
    leads_qualified = int(row.get("leads_qualified", 0) or 0)
    average_score = 0 if leads_total == 0 else min(100, 30 + int((leads_qualified / max(leads_total, 1)) * 70))
    confidence = "stable" if int(row.get("support_overdue", 0) or 0) == 0 else "attention"

    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "scoring": {
            "average": average_score,
            "hot": leads_qualified,
            "warm": max(leads_total - leads_qualified, 0),
            "cold": max(leads_total - leads_qualified * 2, 0),
        },
        "pipeline": {
            "configs": int(row.get("pipeline_configs", 0) or 0),
            "stages": int(row.get("pipeline_stages", 0) or 0),
            "autoScoring": int(row.get("pipeline_auto_scoring", 0) or 0) > 0,
            "territoryRules": int(row.get("pipeline_territory_rules", 0) or 0),
            "approvalPolicies": int(row.get("pipeline_approval_policies", 0) or 0),
        },
        "support": {
            "openCases": int(row.get("support_open", 0) or 0),
            "overdueCases": int(row.get("support_overdue", 0) or 0),
            "slaTrackedCases": int(row.get("support_open", 0) or 0),
        },
        "conversations": {
            "threads": int(row.get("conversation_threads", 0) or 0),
            "participants": int(row.get("conversation_participants", 0) or 0),
            "channels": int(row.get("conversation_channels", 0) or 0),
        },
        "bulkOperations": {
            "importsReady": True,
            "exportsReady": True,
            "partialSuccessTracking": True,
        },
        "approvals": {
            "policies": int(row.get("pipeline_approval_policies", 0) or 0),
            "auditReady": int(row.get("pipeline_approval_policies", 0) or 0) > 0,
        },
        "territories": {
            "rules": int(row.get("pipeline_territory_rules", 0) or 0),
            "assignmentReady": int(row.get("pipeline_territory_rules", 0) or 0) > 0,
        },
        "forecast": {
            "bookedRevenueCents": int(row.get("booked_revenue_cents", 0) or 0),
            "weightedPipelineCents": int(row.get("weighted_pipeline_cents", 0) or 0),
            "confidence": confidence,
            "scenarioCount": 2 if int(row.get("weighted_pipeline_cents", 0) or 0) > 0 else 0,
        },
        "readiness": {
            "status": confidence,
            "slaReady": True,
            "forecastReady": int(row.get("weighted_pipeline_cents", 0) or 0) > 0,
            "relationshipSignalsReady": leads_total > 0,
            "bulkReady": True,
        },
    }
