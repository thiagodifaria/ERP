from __future__ import annotations

from datetime import datetime, timezone
import json
import os
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "profiles": [],
    "documents": [],
    "documentEvents": [],
    "retentionPolicies": [],
    "consents": [],
    "privacyRequests": [],
    "auditEvents": [],
}


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def _find_tenant_id(cursor, tenant_slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (tenant_slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def _append_audit_event(tenant_slug: str, company_public_id: str, category: str, summary: str, actor: str, payload: dict | None = None) -> dict:
    event = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": tenant_slug,
        "companyPublicId": company_public_id,
        "category": category,
        "summary": summary,
        "actor": actor,
        "payload": payload or {},
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["auditEvents"].append(event)
    return event


def _append_document_event(
    tenant_slug: str,
    company_public_id: str,
    document_public_id: str,
    event_type: str,
    summary: str,
    actor: str,
    payload: dict | None = None,
) -> dict:
    event = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": tenant_slug,
        "companyPublicId": company_public_id,
        "documentPublicId": document_public_id,
        "eventType": event_type,
        "summary": summary,
        "actor": actor,
        "payload": payload or {},
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["documentEvents"].append(event)
    return event


def list_capabilities() -> dict:
    focus_nfe_configured = os.getenv("FISCAL_FOCUS_NFE_API_KEY", "").strip() != ""
    enotas_configured = os.getenv("FISCAL_ENOTAS_API_KEY", "").strip() != ""
    certificate_a1_configured = os.getenv("FISCAL_CERTIFICATE_A1_PATH", "").strip() != ""
    certificate_a3_configured = os.getenv("FISCAL_CERTIFICATE_A3_PROVIDER", "").strip() != ""
    provider_modes = [
        {
            "capabilityKey": "fiscal.nfe",
            "provider": "focus_nfe",
            "configured": focus_nfe_configured,
            "mode": "configured" if focus_nfe_configured else "unconfigured",
            "critical": True,
        },
        {
            "capabilityKey": "fiscal.nfse",
            "provider": "enotas",
            "configured": enotas_configured,
            "mode": "configured" if enotas_configured else "unconfigured",
            "critical": True,
        },
        {
            "capabilityKey": "fiscal.certificate.a1",
            "provider": "local_a1",
            "configured": certificate_a1_configured,
            "mode": "configured" if certificate_a1_configured else "fallback",
            "critical": False,
        },
        {
            "capabilityKey": "fiscal.certificate.a3",
            "provider": "smartcard",
            "configured": certificate_a3_configured,
            "mode": "configured" if certificate_a3_configured else "disabled",
            "critical": False,
        },
    ]
    return {"service": settings.service_name, "capabilities": provider_modes}


def _provider_readiness() -> dict:
    capabilities = list_capabilities()["capabilities"]
    configured = sum(1 for item in capabilities if item["configured"])
    critical_gaps = sum(1 for item in capabilities if item["critical"] and item["mode"] == "unconfigured")
    return {
        "configured": configured,
        "criticalGaps": critical_gaps,
        "providerReady": critical_gaps == 0,
    }


def upsert_company_profile(company_public_id: str, payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "companyPublicId": company_public_id,
        "taxRegime": (payload.get("taxRegime") or "simples_nacional").strip().lower(),
        "cnae": (payload.get("cnae") or "").strip(),
        "stateRegistration": (payload.get("stateRegistration") or "").strip(),
        "municipalRegistration": (payload.get("municipalRegistration") or "").strip(),
        "certificateMode": (payload.get("certificateMode") or "a1").strip().lower(),
        "certificateLabel": (payload.get("certificateLabel") or "local-fallback").strip(),
        "environmentMode": (payload.get("environmentMode") or "homologation").strip().lower(),
        "updatedAt": utc_now(),
    }
    if settings.repository_driver != "postgres":
        for index, item in enumerate(IN_MEMORY_STATE["profiles"]):
            if item["tenantSlug"] == slug and item["companyPublicId"] == company_public_id:
                record["publicId"] = item["publicId"]
                IN_MEMORY_STATE["profiles"][index] = record
                _append_audit_event(slug, company_public_id, "company_profile", "Fiscal profile updated.", payload.get("actor", "ops"))
                return record
        IN_MEMORY_STATE["profiles"].append(record)
        _append_audit_event(slug, company_public_id, "company_profile", "Fiscal profile created.", payload.get("actor", "ops"))
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO fiscal.company_profiles (
                  tenant_id, public_id, company_public_id, tax_regime, cnae, state_registration,
                  municipal_registration, certificate_mode, certificate_label, environment_mode
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, company_public_id)
                DO UPDATE SET
                  tax_regime = EXCLUDED.tax_regime,
                  cnae = EXCLUDED.cnae,
                  state_registration = EXCLUDED.state_registration,
                  municipal_registration = EXCLUDED.municipal_registration,
                  certificate_mode = EXCLUDED.certificate_mode,
                  certificate_label = EXCLUDED.certificate_label,
                  environment_mode = EXCLUDED.environment_mode
                RETURNING public_id, updated_at
                """,
                (
                    tenant_id,
                    record["publicId"],
                    company_public_id,
                    record["taxRegime"],
                    record["cnae"],
                    record["stateRegistration"],
                    record["municipalRegistration"],
                    record["certificateMode"],
                    record["certificateLabel"],
                    record["environmentMode"],
                ),
            )
            row = cursor.fetchone()
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'company_profile', %s, %s, %s::jsonb)
                """,
                (
                    tenant_id,
                    str(uuid.uuid4()),
                    company_public_id,
                    "Fiscal profile updated.",
                    payload.get("actor", "ops"),
                    json.dumps(payload),
                ),
            )
            connection.commit()
            record["publicId"] = row["public_id"]
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def get_company_profile(company_public_id: str, tenant_slug: str | None = None) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["profiles"]:
            if item["tenantSlug"] == slug and item["companyPublicId"] == company_public_id:
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT public_id, tax_regime, cnae, state_registration, municipal_registration,
                       certificate_mode, certificate_label, environment_mode, updated_at
                FROM fiscal.company_profiles AS profile
                JOIN identity.tenants AS tenant ON tenant.id = profile.tenant_id
                WHERE tenant.slug = %s AND profile.company_public_id = %s
                """,
                (slug, company_public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            return {
                "publicId": row["public_id"],
                "tenantSlug": slug,
                "companyPublicId": company_public_id,
                "taxRegime": row["tax_regime"],
                "cnae": row["cnae"],
                "stateRegistration": row["state_registration"],
                "municipalRegistration": row["municipal_registration"],
                "certificateMode": row["certificate_mode"],
                "certificateLabel": row["certificate_label"],
                "environmentMode": row["environment_mode"],
                "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            }


def upsert_retention_policy(company_public_id: str, data_domain: str, payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    domain = data_domain.strip().lower()
    if domain == "":
        raise ValueError("fiscal_retention_domain_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "companyPublicId": company_public_id,
        "dataDomain": domain,
        "classification": (payload.get("classification") or "internal").strip().lower(),
        "retentionDays": int(payload.get("retentionDays") or 365),
        "anonymizeAfterDays": int(payload.get("anonymizeAfterDays") or 730),
        "source": (payload.get("source") or "manual").strip().lower(),
        "updatedAt": utc_now(),
    }
    if settings.repository_driver != "postgres":
        for index, item in enumerate(IN_MEMORY_STATE["retentionPolicies"]):
            if item["tenantSlug"] == slug and item["companyPublicId"] == company_public_id and item["dataDomain"] == domain:
                record["publicId"] = item["publicId"]
                IN_MEMORY_STATE["retentionPolicies"][index] = record
                return record
        IN_MEMORY_STATE["retentionPolicies"].append(record)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO fiscal.retention_policies (
                  tenant_id, public_id, company_public_id, data_domain, classification, retention_days,
                  anonymize_after_days, source
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, company_public_id, data_domain)
                DO UPDATE SET
                  classification = EXCLUDED.classification,
                  retention_days = EXCLUDED.retention_days,
                  anonymize_after_days = EXCLUDED.anonymize_after_days,
                  source = EXCLUDED.source
                RETURNING public_id, updated_at
                """,
                (
                    tenant_id,
                    record["publicId"],
                    company_public_id,
                    domain,
                    record["classification"],
                    record["retentionDays"],
                    record["anonymizeAfterDays"],
                    record["source"],
                ),
            )
            row = cursor.fetchone()
            connection.commit()
            record["publicId"] = row["public_id"]
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def list_retention_policies(company_public_id: str, tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [
            item for item in IN_MEMORY_STATE["retentionPolicies"] if item["tenantSlug"] == slug and item["companyPublicId"] == company_public_id
        ]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT public_id, data_domain, classification, retention_days, anonymize_after_days, source, updated_at
                FROM fiscal.retention_policies AS policy
                JOIN identity.tenants AS tenant ON tenant.id = policy.tenant_id
                WHERE tenant.slug = %s AND policy.company_public_id = %s
                ORDER BY policy.data_domain
                """,
                (slug, company_public_id),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "companyPublicId": company_public_id,
                    "dataDomain": row["data_domain"],
                    "classification": row["classification"],
                    "retentionDays": int(row["retention_days"]),
                    "anonymizeAfterDays": int(row["anonymize_after_days"]),
                    "source": row["source"],
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def create_document(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    company_public_id = (payload.get("companyPublicId") or "").strip()
    if company_public_id == "":
        raise ValueError("fiscal_company_public_id_required")
    document_kind = (payload.get("documentKind") or "nfe").strip().lower()
    if document_kind not in {"nfe", "nfse"}:
        raise ValueError("fiscal_document_kind_invalid")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "companyPublicId": company_public_id,
        "documentKind": document_kind,
        "seriesCode": (payload.get("seriesCode") or "1").strip(),
        "numberCode": (payload.get("numberCode") or uuid.uuid4().hex[:6].upper()).strip(),
        "status": "issued",
        "customerPublicId": (payload.get("customerPublicId") or "").strip() or None,
        "amountCents": int(payload.get("amountCents") or 0),
        "providerKey": (payload.get("providerKey") or "local").strip().lower(),
        "issuedAt": utc_now(),
        "cancelledAt": None,
        "payload": payload.get("payload") or {},
    }
    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["documents"].append(record)
        _append_document_event(
            slug,
            company_public_id,
            record["publicId"],
            "issued",
            "Fiscal document issued.",
            payload.get("actor", "fiscal-bot"),
            payload,
        )
        _append_audit_event(slug, company_public_id, "fiscal_document", "Fiscal document issued.", payload.get("actor", "fiscal-bot"), payload)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO fiscal.documents (
                  tenant_id, public_id, company_public_id, document_kind, series_code, number_code, status,
                  customer_public_id, amount_cents, provider_key, payload_json, issued_at
                )
                VALUES (%s, %s, %s, %s, %s, %s, 'issued', %s, %s, %s, %s::jsonb, NOW())
                RETURNING issued_at
                """,
                (
                    tenant_id,
                    record["publicId"],
                    company_public_id,
                    document_kind,
                    record["seriesCode"],
                    record["numberCode"],
                    record["customerPublicId"],
                    record["amountCents"],
                    record["providerKey"],
                    json.dumps(record["payload"]),
                ),
            )
            row = cursor.fetchone()
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'fiscal_document', 'Fiscal document issued.', %s, %s::jsonb)
                """,
                (tenant_id, str(uuid.uuid4()), company_public_id, payload.get("actor", "fiscal-bot"), json.dumps(payload)),
            )
            cursor.execute(
                """
                INSERT INTO fiscal.document_events (
                  tenant_id, public_id, company_public_id, document_public_id, event_type, summary, actor, payload_json
                )
                VALUES (%s, %s, %s, %s, 'issued', 'Fiscal document issued.', %s, %s::jsonb)
                """,
                (tenant_id, str(uuid.uuid4()), company_public_id, record["publicId"], payload.get("actor", "fiscal-bot"), json.dumps(payload)),
            )
            connection.commit()
            record["issuedAt"] = row["issued_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def list_documents(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [item for item in IN_MEMORY_STATE["documents"] if item["tenantSlug"] == slug]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT public_id, company_public_id, document_kind, series_code, number_code, status, customer_public_id,
                       amount_cents, provider_key, issued_at, cancelled_at
                FROM fiscal.documents AS document
                JOIN identity.tenants AS tenant ON tenant.id = document.tenant_id
                WHERE tenant.slug = %s
                ORDER BY document.issued_at DESC
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "companyPublicId": row["company_public_id"],
                    "documentKind": row["document_kind"],
                    "seriesCode": row["series_code"],
                    "numberCode": row["number_code"],
                    "status": row["status"],
                    "customerPublicId": row["customer_public_id"],
                    "amountCents": int(row["amount_cents"]),
                    "providerKey": row["provider_key"],
                    "issuedAt": row["issued_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "cancelledAt": None if row["cancelled_at"] is None else row["cancelled_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def cancel_document(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    reason = (payload.get("reason") or "Cancelled by operator.").strip()
    actor = (payload.get("actor") or "ops").strip()
    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["documents"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                item["status"] = "cancelled"
                item["cancelledAt"] = utc_now()
                _append_document_event(slug, item["companyPublicId"], public_id, "cancelled", reason, actor, payload)
                _append_audit_event(slug, item["companyPublicId"], "fiscal_document_cancelled", reason, actor, payload)
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE fiscal.documents AS document
                SET status = 'cancelled', cancelled_at = NOW()
                FROM identity.tenants AS tenant
                WHERE tenant.id = document.tenant_id
                  AND tenant.slug = %s
                  AND document.public_id = %s
                RETURNING document.company_public_id, document.document_kind, document.series_code, document.number_code,
                          document.amount_cents, document.provider_key, document.issued_at, document.cancelled_at,
                          document.customer_public_id
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'fiscal_document_cancelled', %s, %s, %s::jsonb)
                """,
                (tenant_id, str(uuid.uuid4()), row["company_public_id"], reason, actor, json.dumps(payload)),
            )
            cursor.execute(
                """
                INSERT INTO fiscal.document_events (
                  tenant_id, public_id, company_public_id, document_public_id, event_type, summary, actor, payload_json
                )
                VALUES (%s, %s, %s, %s, 'cancelled', %s, %s, %s::jsonb)
                """,
                (tenant_id, str(uuid.uuid4()), row["company_public_id"], public_id, reason, actor, json.dumps(payload)),
            )
            connection.commit()
            return {
                "publicId": public_id,
                "tenantSlug": slug,
                "companyPublicId": row["company_public_id"],
                "documentKind": row["document_kind"],
                "seriesCode": row["series_code"],
                "numberCode": row["number_code"],
                "status": "cancelled",
                "customerPublicId": row["customer_public_id"],
                "amountCents": int(row["amount_cents"]),
                "providerKey": row["provider_key"],
                "issuedAt": row["issued_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "cancelledAt": row["cancelled_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            }


def create_correction_letter(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    actor = (payload.get("actor") or "fiscal-ops").strip()
    correction = (payload.get("correctionText") or "").strip()
    if correction == "":
        raise ValueError("fiscal_correction_text_required")

    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["documents"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                item["status"] = "corrected"
                _append_document_event(slug, item["companyPublicId"], public_id, "correction_letter", correction, actor, payload)
                _append_audit_event(slug, item["companyPublicId"], "fiscal_document_correction", correction, actor, payload)
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE fiscal.documents AS document
                SET status = 'corrected'
                FROM identity.tenants AS tenant
                WHERE tenant.id = document.tenant_id
                  AND tenant.slug = %s
                  AND document.public_id = %s
                RETURNING tenant.id AS tenant_id, document.company_public_id, document.document_kind, document.series_code,
                          document.number_code, document.amount_cents, document.provider_key, document.issued_at,
                          document.cancelled_at, document.customer_public_id
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'fiscal_document_correction', %s, %s, %s::jsonb)
                """,
                (row["tenant_id"], str(uuid.uuid4()), row["company_public_id"], correction, actor, json.dumps(payload)),
            )
            cursor.execute(
                """
                INSERT INTO fiscal.document_events (
                  tenant_id, public_id, company_public_id, document_public_id, event_type, summary, actor, payload_json
                )
                VALUES (%s, %s, %s, %s, 'correction_letter', %s, %s, %s::jsonb)
                """,
                (row["tenant_id"], str(uuid.uuid4()), row["company_public_id"], public_id, correction, actor, json.dumps(payload)),
            )
            connection.commit()
            return {
                "publicId": public_id,
                "tenantSlug": slug,
                "companyPublicId": row["company_public_id"],
                "documentKind": row["document_kind"],
                "seriesCode": row["series_code"],
                "numberCode": row["number_code"],
                "status": "corrected",
                "customerPublicId": row["customer_public_id"],
                "amountCents": int(row["amount_cents"]),
                "providerKey": row["provider_key"],
                "issuedAt": row["issued_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "cancelledAt": None if row["cancelled_at"] is None else row["cancelled_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            }


def create_invalidation(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    actor = (payload.get("actor") or "fiscal-ops").strip()
    reason = (payload.get("reason") or "Invalidated by fiscal operator.").strip()

    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["documents"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                item["status"] = "cancelled"
                item["cancelledAt"] = utc_now()
                _append_document_event(slug, item["companyPublicId"], public_id, "invalidation", reason, actor, payload)
                _append_audit_event(slug, item["companyPublicId"], "fiscal_document_invalidation", reason, actor, payload)
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE fiscal.documents AS document
                SET status = 'cancelled', cancelled_at = NOW()
                FROM identity.tenants AS tenant
                WHERE tenant.id = document.tenant_id
                  AND tenant.slug = %s
                  AND document.public_id = %s
                RETURNING tenant.id AS tenant_id, document.company_public_id, document.document_kind, document.series_code,
                          document.number_code, document.amount_cents, document.provider_key, document.issued_at,
                          document.cancelled_at, document.customer_public_id
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'fiscal_document_invalidation', %s, %s, %s::jsonb)
                """,
                (row["tenant_id"], str(uuid.uuid4()), row["company_public_id"], reason, actor, json.dumps(payload)),
            )
            cursor.execute(
                """
                INSERT INTO fiscal.document_events (
                  tenant_id, public_id, company_public_id, document_public_id, event_type, summary, actor, payload_json
                )
                VALUES (%s, %s, %s, %s, 'invalidation', %s, %s, %s::jsonb)
                """,
                (row["tenant_id"], str(uuid.uuid4()), row["company_public_id"], public_id, reason, actor, json.dumps(payload)),
            )
            connection.commit()
            return {
                "publicId": public_id,
                "tenantSlug": slug,
                "companyPublicId": row["company_public_id"],
                "documentKind": row["document_kind"],
                "seriesCode": row["series_code"],
                "numberCode": row["number_code"],
                "status": "cancelled",
                "customerPublicId": row["customer_public_id"],
                "amountCents": int(row["amount_cents"]),
                "providerKey": row["provider_key"],
                "issuedAt": row["issued_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                "cancelledAt": row["cancelled_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            }


def create_privacy_request(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    company_public_id = (payload.get("companyPublicId") or "").strip()
    request_type = (payload.get("requestType") or "").strip().lower()
    if request_type not in {"access", "portability", "anonymization", "deletion", "consent_revoke"}:
        raise ValueError("fiscal_privacy_request_type_invalid")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "companyPublicId": company_public_id,
        "requestType": request_type,
        "subjectKind": (payload.get("subjectKind") or "user").strip().lower(),
        "subjectPublicId": (payload.get("subjectPublicId") or "").strip(),
        "status": "received",
        "requestedBy": (payload.get("requestedBy") or "privacy@erp.local").strip(),
        "consentReference": (payload.get("consentReference") or "").strip() or None,
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }
    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["privacyRequests"].append(record)
        _append_audit_event(slug, company_public_id, "privacy_request", "Privacy request received.", record["requestedBy"], payload)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO fiscal.privacy_requests (
                  tenant_id, public_id, company_public_id, request_type, subject_kind, subject_public_id,
                  status, requested_by, consent_reference
                )
                VALUES (%s, %s, %s, %s, %s, %s, 'received', %s, %s)
                RETURNING created_at, updated_at
                """,
                (
                    tenant_id,
                    record["publicId"],
                    company_public_id,
                    request_type,
                    record["subjectKind"],
                    record["subjectPublicId"],
                    record["requestedBy"],
                    record["consentReference"],
                ),
            )
            row = cursor.fetchone()
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'privacy_request', 'Privacy request received.', %s, %s::jsonb)
                """,
                (tenant_id, str(uuid.uuid4()), company_public_id, record["requestedBy"], json.dumps(payload)),
            )
            connection.commit()
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def create_consent(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    company_public_id = (payload.get("companyPublicId") or "").strip()
    subject_public_id = (payload.get("subjectPublicId") or "").strip()
    purpose_key = (payload.get("purposeKey") or "").strip()
    if company_public_id == "":
        raise ValueError("fiscal_company_public_id_required")
    if subject_public_id == "":
        raise ValueError("fiscal_consent_subject_required")
    if purpose_key == "":
        raise ValueError("fiscal_consent_purpose_required")

    actor = (payload.get("actor") or "privacy@erp.local").strip()
    status = (payload.get("status") or "granted").strip().lower()
    if status not in {"granted", "revoked"}:
        raise ValueError("fiscal_consent_status_invalid")

    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "companyPublicId": company_public_id,
        "subjectKind": (payload.get("subjectKind") or "customer").strip().lower(),
        "subjectPublicId": subject_public_id,
        "purposeKey": purpose_key,
        "status": status,
        "source": (payload.get("source") or "ops").strip().lower(),
        "grantedAt": utc_now() if status == "granted" else None,
        "revokedAt": utc_now() if status == "revoked" else None,
        "updatedAt": utc_now(),
    }
    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["consents"].append(record)
        _append_audit_event(slug, company_public_id, "consent", f"Consent {status}.", actor, payload)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO fiscal.consents (
                  tenant_id, public_id, company_public_id, subject_kind, subject_public_id, purpose_key,
                  status, source, granted_at, revoked_at
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s,
                        CASE WHEN %s = 'granted' THEN NOW() ELSE NULL END,
                        CASE WHEN %s = 'revoked' THEN NOW() ELSE NULL END)
                RETURNING granted_at, revoked_at, updated_at
                """,
                (
                    tenant_id,
                    record["publicId"],
                    company_public_id,
                    record["subjectKind"],
                    subject_public_id,
                    purpose_key,
                    status,
                    record["source"],
                    status,
                    status,
                ),
            )
            row = cursor.fetchone()
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'consent', %s, %s, %s::jsonb)
                """,
                (tenant_id, str(uuid.uuid4()), company_public_id, f"Consent {status}.", actor, json.dumps(payload)),
            )
            connection.commit()
            record["grantedAt"] = None if row["granted_at"] is None else row["granted_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["revokedAt"] = None if row["revoked_at"] is None else row["revoked_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def transition_consent(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    status = (payload.get("status") or "").strip().lower()
    actor = (payload.get("actor") or "privacy@erp.local").strip()
    if status not in {"granted", "revoked"}:
        raise ValueError("fiscal_consent_status_invalid")

    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["consents"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                item["status"] = status
                if status == "granted":
                    item["grantedAt"] = utc_now()
                    item["revokedAt"] = None
                else:
                    item["revokedAt"] = utc_now()
                item["updatedAt"] = utc_now()
                _append_audit_event(slug, item["companyPublicId"], "consent", f"Consent {status}.", actor, payload)
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE fiscal.consents AS consent
                SET status = %s,
                    granted_at = CASE WHEN %s = 'granted' THEN NOW() ELSE consent.granted_at END,
                    revoked_at = CASE WHEN %s = 'revoked' THEN NOW() ELSE NULL END
                FROM identity.tenants AS tenant
                WHERE tenant.id = consent.tenant_id
                  AND tenant.slug = %s
                  AND consent.public_id = %s
                RETURNING tenant.id AS tenant_id, consent.company_public_id
                """,
                (status, status, status, slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'consent', %s, %s, %s::jsonb)
                """,
                (row["tenant_id"], str(uuid.uuid4()), row["company_public_id"], f"Consent {status}.", actor, json.dumps(payload)),
            )
            connection.commit()
            return next((item for item in list_consents(slug) if item["publicId"] == public_id), None)


def list_consents(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [item for item in IN_MEMORY_STATE["consents"] if item["tenantSlug"] == slug]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT consent.public_id,
                       consent.company_public_id,
                       consent.subject_kind,
                       consent.subject_public_id,
                       consent.purpose_key,
                       consent.status,
                       consent.source,
                       consent.granted_at,
                       consent.revoked_at,
                       consent.updated_at
                FROM fiscal.consents AS consent
                JOIN identity.tenants AS tenant ON tenant.id = consent.tenant_id
                WHERE tenant.slug = %s
                ORDER BY consent.updated_at DESC
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "companyPublicId": row["company_public_id"],
                    "subjectKind": row["subject_kind"],
                    "subjectPublicId": row["subject_public_id"],
                    "purposeKey": row["purpose_key"],
                    "status": row["status"],
                    "source": row["source"],
                    "grantedAt": None if row["granted_at"] is None else row["granted_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "revokedAt": None if row["revoked_at"] is None else row["revoked_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def list_privacy_requests(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [item for item in IN_MEMORY_STATE["privacyRequests"] if item["tenantSlug"] == slug]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT request.public_id,
                       request.company_public_id,
                       request.request_type,
                       request.subject_kind,
                       request.subject_public_id,
                       request.status,
                       request.requested_by,
                       request.consent_reference,
                       request.created_at,
                       request.updated_at
                FROM fiscal.privacy_requests AS request
                JOIN identity.tenants AS tenant ON tenant.id = request.tenant_id
                WHERE tenant.slug = %s
                ORDER BY request.created_at DESC
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "companyPublicId": row["company_public_id"],
                    "requestType": row["request_type"],
                    "subjectKind": row["subject_kind"],
                    "subjectPublicId": row["subject_public_id"],
                    "status": row["status"],
                    "requestedBy": row["requested_by"],
                    "consentReference": row["consent_reference"],
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def transition_privacy_request(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    status = (payload.get("status") or "").strip().lower()
    actor = (payload.get("actor") or "privacy@erp.local").strip()
    if status not in {"received", "processing", "completed", "denied"}:
        raise ValueError("fiscal_privacy_request_status_invalid")

    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["privacyRequests"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                item["status"] = status
                item["updatedAt"] = utc_now()
                _append_audit_event(slug, item["companyPublicId"], "privacy_request", f"Privacy request moved to {status}.", actor, payload)
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE fiscal.privacy_requests AS request
                SET status = %s
                FROM identity.tenants AS tenant
                WHERE tenant.id = request.tenant_id
                  AND tenant.slug = %s
                  AND request.public_id = %s
                RETURNING tenant.id AS tenant_id, request.company_public_id
                """,
                (status, slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            cursor.execute(
                """
                INSERT INTO fiscal.audit_events (tenant_id, public_id, company_public_id, category, summary, actor, payload_json)
                VALUES (%s, %s, %s, 'privacy_request', %s, %s, %s::jsonb)
                """,
                (row["tenant_id"], str(uuid.uuid4()), row["company_public_id"], f"Privacy request moved to {status}.", actor, json.dumps(payload)),
            )
            connection.commit()
            return next((item for item in list_privacy_requests(slug) if item["publicId"] == public_id), None)


def list_audit_events(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return [item for item in IN_MEMORY_STATE["auditEvents"] if item["tenantSlug"] == slug]

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT event.public_id,
                       event.company_public_id,
                       event.category,
                       event.summary,
                       event.actor,
                       event.payload_json,
                       event.created_at
                FROM fiscal.audit_events AS event
                JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
                WHERE tenant.slug = %s
                ORDER BY event.created_at DESC
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "companyPublicId": row["company_public_id"],
                    "category": row["category"],
                    "summary": row["summary"],
                    "actor": row["actor"],
                    "payload": row["payload_json"] or {},
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def list_document_events(tenant_slug: str | None = None, document_public_id: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        items = [item for item in IN_MEMORY_STATE["documentEvents"] if item["tenantSlug"] == slug]
        if document_public_id:
            items = [item for item in items if item["documentPublicId"] == document_public_id]
        return items

    clauses = ["tenant.slug = %s"]
    params: list[object] = [slug]
    if document_public_id:
        clauses.append("event.document_public_id = %s::uuid")
        params.append(document_public_id)

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                f"""
                SELECT event.public_id,
                       event.company_public_id,
                       event.document_public_id,
                       event.event_type,
                       event.summary,
                       event.actor,
                       event.payload_json,
                       event.created_at
                FROM fiscal.document_events AS event
                JOIN identity.tenants AS tenant ON tenant.id = event.tenant_id
                WHERE {" AND ".join(clauses)}
                ORDER BY event.created_at DESC
                """,
                params,
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "companyPublicId": row["company_public_id"],
                    "documentPublicId": str(row["document_public_id"]),
                    "eventType": row["event_type"],
                    "summary": row["summary"],
                    "actor": row["actor"],
                    "payload": row["payload_json"] or {},
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def build_compliance_summary(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    provider_readiness = _provider_readiness()
    if settings.repository_driver != "postgres":
        return {
            "tenantSlug": slug,
            "summary": {
                "fiscalDocuments": len([item for item in IN_MEMORY_STATE["documents"] if item["tenantSlug"] == slug]),
                "documentEvents": len([item for item in IN_MEMORY_STATE["documentEvents"] if item["tenantSlug"] == slug]),
                "consents": len([item for item in IN_MEMORY_STATE["consents"] if item["tenantSlug"] == slug]),
                "privacyRequests": len([item for item in IN_MEMORY_STATE["privacyRequests"] if item["tenantSlug"] == slug]),
                "retentionPolicies": len([item for item in IN_MEMORY_STATE["retentionPolicies"] if item["tenantSlug"] == slug]),
                "auditEvents": len([item for item in IN_MEMORY_STATE["auditEvents"] if item["tenantSlug"] == slug]),
            },
            "providers": provider_readiness,
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  (SELECT count(*) FROM fiscal.documents AS document JOIN identity.tenants AS tenant_doc ON tenant_doc.id = document.tenant_id WHERE tenant_doc.slug = %s) AS fiscal_documents,
                  (SELECT count(*) FROM fiscal.document_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s) AS document_events,
                  (SELECT count(*) FROM fiscal.consents AS consent JOIN identity.tenants AS tenant_consent ON tenant_consent.id = consent.tenant_id WHERE tenant_consent.slug = %s) AS consents,
                  (SELECT count(*) FROM fiscal.privacy_requests AS request JOIN identity.tenants AS tenant_req ON tenant_req.id = request.tenant_id WHERE tenant_req.slug = %s) AS privacy_requests,
                  (SELECT count(*) FROM fiscal.retention_policies AS policy JOIN identity.tenants AS tenant_pol ON tenant_pol.id = policy.tenant_id WHERE tenant_pol.slug = %s) AS retention_policies,
                  (SELECT count(*) FROM fiscal.audit_events AS event JOIN identity.tenants AS tenant_evt ON tenant_evt.id = event.tenant_id WHERE tenant_evt.slug = %s) AS audit_events
                """,
                (slug, slug, slug, slug, slug, slug),
            )
            row = cursor.fetchone() or {}
            return {
                "tenantSlug": slug,
                "summary": {
                    "fiscalDocuments": int(row.get("fiscal_documents", 0) or 0),
                    "documentEvents": int(row.get("document_events", 0) or 0),
                    "consents": int(row.get("consents", 0) or 0),
                    "privacyRequests": int(row.get("privacy_requests", 0) or 0),
                    "retentionPolicies": int(row.get("retention_policies", 0) or 0),
                    "auditEvents": int(row.get("audit_events", 0) or 0),
                },
                "providers": provider_readiness,
            }
