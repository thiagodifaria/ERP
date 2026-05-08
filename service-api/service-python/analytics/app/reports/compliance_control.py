"""Painel de fiscal, privacidade e governanca de dados."""

from datetime import datetime, timezone
import os

from app.config.settings import settings
from app.infrastructure.postgres import connect


def _provider_ready() -> bool:
    return any(
        (
            os.getenv("FISCAL_FOCUS_NFE_API_KEY", "").strip(),
            os.getenv("FISCAL_ENOTAS_API_KEY", "").strip(),
            os.getenv("FISCAL_CERTIFICATE_A1_PATH", "").strip(),
            os.getenv("FISCAL_CERTIFICATE_A3_PROVIDER", "").strip(),
        )
    )


def build_compliance_control(tenant_slug: str | None = None) -> dict:
    slug = tenant_slug or "bootstrap-ops"
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": slug,
            "generatedAt": datetime.now(timezone.utc).isoformat(),
            "fiscal": {"documents": 4, "cancelled": 1, "profiles": 1},
            "documentOperations": {"events": 6, "providerReady": _provider_ready(), "anonymized": 1, "retentionReview": 1},
            "privacy": {"requests": 2, "pending": 1, "completed": 1, "exportPackages": 1},
            "consents": {"granted": 3, "revoked": 1},
            "retention": {"policies": 5, "restrictedDocuments": 10, "executions": 1},
            "audit": {"events": 12, "sensitiveOperations": 4},
            "readiness": {"status": "stable", "fiscalReady": True, "privacyReady": True, "retentionReady": True, "providerReady": _provider_ready()},
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  (SELECT count(*) FROM fiscal.documents AS document JOIN identity.tenants AS tenant_doc ON tenant_doc.id = document.tenant_id WHERE tenant_doc.slug = %s) AS fiscal_documents,
                  (SELECT count(*) FROM fiscal.documents AS document JOIN identity.tenants AS tenant_doc ON tenant_doc.id = document.tenant_id WHERE tenant_doc.slug = %s AND document.status = 'cancelled') AS fiscal_cancelled,
                  (SELECT count(*) FROM fiscal.company_profiles AS profile JOIN identity.tenants AS tenant_profile ON tenant_profile.id = profile.tenant_id WHERE tenant_profile.slug = %s) AS fiscal_profiles,
                  (SELECT count(*) FROM fiscal.document_events AS document_event JOIN identity.tenants AS tenant_event ON tenant_event.id = document_event.tenant_id WHERE tenant_event.slug = %s) AS document_events,
                  (SELECT count(*) FROM fiscal.privacy_requests AS request JOIN identity.tenants AS tenant_req ON tenant_req.id = request.tenant_id WHERE tenant_req.slug = %s) AS privacy_requests,
                  (SELECT count(*) FROM fiscal.privacy_requests AS request JOIN identity.tenants AS tenant_req ON tenant_req.id = request.tenant_id WHERE tenant_req.slug = %s AND request.status <> 'completed') AS privacy_pending,
                  (SELECT count(*) FROM fiscal.privacy_requests AS request JOIN identity.tenants AS tenant_req ON tenant_req.id = request.tenant_id WHERE tenant_req.slug = %s AND request.status = 'completed') AS privacy_completed,
                  (SELECT count(*) FROM fiscal.consents AS consent JOIN identity.tenants AS tenant_consent ON tenant_consent.id = consent.tenant_id WHERE tenant_consent.slug = %s AND consent.status = 'granted') AS consents_granted,
                  (SELECT count(*) FROM fiscal.consents AS consent JOIN identity.tenants AS tenant_consent ON tenant_consent.id = consent.tenant_id WHERE tenant_consent.slug = %s AND consent.status = 'revoked') AS consents_revoked,
                  (SELECT count(*) FROM fiscal.retention_policies AS policy JOIN identity.tenants AS tenant_pol ON tenant_pol.id = policy.tenant_id WHERE tenant_pol.slug = %s) AS retention_policies,
                  (SELECT count(*) FROM documents.attachments AS attachment JOIN identity.tenants AS tenant_att ON tenant_att.id = attachment.tenant_id WHERE tenant_att.slug = %s AND attachment.visibility = 'restricted') AS restricted_documents,
                  (SELECT count(*) FROM fiscal.document_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s AND event.event_type IN ('anonymized', 'deletion_controlled', 'retention_anonymized')) AS anonymized_documents,
                  (SELECT count(*) FROM fiscal.document_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s AND event.event_type = 'retention_review') AS retention_review_events,
                  (SELECT count(*) FROM fiscal.audit_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s AND event.category = 'privacy_request_execution') AS privacy_export_packages,
                  (SELECT count(*) FROM fiscal.audit_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s AND event.category = 'retention_execution') AS retention_executions,
                  (SELECT count(*) FROM fiscal.audit_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s) AS audit_events,
                  (SELECT count(*) FROM fiscal.audit_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s AND event.category IN ('privacy_request_execution', 'retention_execution', 'consent', 'company_profile', 'fiscal_document')) AS sensitive_operations
                """,
                (slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug),
            )
            row = cursor.fetchone() or {}

    status = "stable" if int(row.get("privacy_pending", 0) or 0) == 0 else "attention"
    return {
        "tenantSlug": slug,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "fiscal": {
            "documents": int(row.get("fiscal_documents", 0) or 0),
            "cancelled": int(row.get("fiscal_cancelled", 0) or 0),
            "profiles": int(row.get("fiscal_profiles", 0) or 0),
        },
        "documentOperations": {
            "events": int(row.get("document_events", 0) or 0),
            "providerReady": _provider_ready(),
            "anonymized": int(row.get("anonymized_documents", 0) or 0),
            "retentionReview": int(row.get("retention_review_events", 0) or 0),
        },
        "privacy": {
            "requests": int(row.get("privacy_requests", 0) or 0),
            "pending": int(row.get("privacy_pending", 0) or 0),
            "completed": int(row.get("privacy_completed", 0) or 0),
            "exportPackages": int(row.get("privacy_export_packages", 0) or 0),
        },
        "consents": {
            "granted": int(row.get("consents_granted", 0) or 0),
            "revoked": int(row.get("consents_revoked", 0) or 0),
        },
        "retention": {
            "policies": int(row.get("retention_policies", 0) or 0),
            "restrictedDocuments": int(row.get("restricted_documents", 0) or 0),
            "executions": int(row.get("retention_executions", 0) or 0),
        },
        "audit": {
            "events": int(row.get("audit_events", 0) or 0),
            "sensitiveOperations": int(row.get("sensitive_operations", 0) or 0),
        },
        "readiness": {
            "status": status,
            "fiscalReady": int(row.get("fiscal_profiles", 0) or 0) > 0,
            "privacyReady": int(row.get("privacy_requests", 0) or 0) > 0,
            "retentionReady": int(row.get("retention_policies", 0) or 0) > 0,
            "providerReady": _provider_ready(),
        },
    }
