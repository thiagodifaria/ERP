from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_fiscal_profile_document_and_privacy_flow() -> None:
    profile_response = client.put(
        "/api/fiscal/companies/company-001/profile",
        json={
            "tenantSlug": "bootstrap-ops",
            "taxRegime": "lucro_real",
            "cnae": "6201-5/01",
            "stateRegistration": "123456789",
            "municipalRegistration": "998877",
            "certificateMode": "a1",
            "certificateLabel": "erp-local-cert",
            "environmentMode": "homologation",
            "actor": "fiscal@erp.local",
        },
    )
    retention_response = client.put(
        "/api/fiscal/companies/company-001/retention-policies/documents",
        json={
            "tenantSlug": "bootstrap-ops",
            "classification": "restricted",
            "retentionDays": 1825,
            "anonymizeAfterDays": 3650,
            "source": "policy",
        },
    )
    document_response = client.post(
        "/api/fiscal/documents",
        json={
            "tenantSlug": "bootstrap-ops",
            "companyPublicId": "company-001",
            "documentKind": "nfe",
            "seriesCode": "1",
            "numberCode": "NF-1001",
            "customerPublicId": "customer-001",
            "amountCents": 159900,
            "providerKey": "local",
            "actor": "fiscal@erp.local",
        },
    )
    document_public_id = document_response.json()["publicId"]
    detail_response = client.get(f"/api/fiscal/documents/{document_public_id}?tenant_slug=bootstrap-ops")
    correction_response = client.post(
        f"/api/fiscal/documents/{document_public_id}/correction-letter",
        json={
            "tenantSlug": "bootstrap-ops",
            "correctionText": "Corrected tax nature and service code.",
            "actor": "fiscal@erp.local",
        },
    )
    cancel_response = client.post(
        f"/api/fiscal/documents/{document_public_id}/cancel",
        json={"tenantSlug": "bootstrap-ops", "reason": "Correction required.", "actor": "fiscal@erp.local"},
    )
    consent_response = client.post(
        "/api/fiscal/consents",
        json={
            "tenantSlug": "bootstrap-ops",
            "companyPublicId": "company-001",
            "subjectKind": "customer",
            "subjectPublicId": "customer-001",
            "purposeKey": "marketing.email",
            "status": "granted",
            "source": "crm",
            "actor": "dpo@erp.local",
        },
    )
    privacy_response = client.post(
        "/api/fiscal/privacy-requests",
        json={
            "tenantSlug": "bootstrap-ops",
            "companyPublicId": "company-001",
            "requestType": "anonymization",
            "subjectKind": "customer",
            "subjectPublicId": "customer-001",
            "requestedBy": "dpo@erp.local",
            "consentReference": "consent-001",
        },
    )
    privacy_public_id = privacy_response.json()["publicId"]
    privacy_detail_response = client.get(f"/api/fiscal/privacy-requests/{privacy_public_id}?tenant_slug=bootstrap-ops")
    privacy_transition_response = client.patch(
        f"/api/fiscal/privacy-requests/{privacy_public_id}/status",
        json={"tenantSlug": "bootstrap-ops", "status": "processing", "actor": "dpo@erp.local"},
    )
    consent_public_id = consent_response.json()["publicId"]
    consent_transition_response = client.patch(
        f"/api/fiscal/consents/{consent_public_id}",
        json={"tenantSlug": "bootstrap-ops", "status": "revoked", "actor": "dpo@erp.local"},
    )
    events_response = client.get(f"/api/fiscal/documents/{document_public_id}/events?tenant_slug=bootstrap-ops")
    summary_response = client.get("/api/fiscal/compliance/summary?tenant_slug=bootstrap-ops")
    export_response = client.get(f"/api/fiscal/privacy-requests/{privacy_public_id}/export-package?tenant_slug=bootstrap-ops")
    retention_execution_response = client.get("/api/fiscal/companies/company-001/retention-execution?tenant_slug=bootstrap-ops")
    execute_privacy_response = client.post(
        f"/api/fiscal/privacy-requests/{privacy_public_id}/execute",
        json={"tenantSlug": "bootstrap-ops", "actor": "dpo@erp.local"},
    )
    execute_retention_response = client.post(
        "/api/fiscal/companies/company-001/retention-execution/execute",
        json={"tenantSlug": "bootstrap-ops", "actor": "compliance@erp.local"},
    )

    assert profile_response.status_code == 200
    assert retention_response.status_code == 200
    assert document_response.status_code == 200
    assert detail_response.status_code == 200
    assert correction_response.status_code == 200
    assert cancel_response.status_code == 200
    assert consent_response.status_code == 200
    assert privacy_response.status_code == 200
    assert privacy_detail_response.status_code == 200
    assert privacy_transition_response.status_code == 200
    assert consent_transition_response.status_code == 200
    assert events_response.status_code == 200
    assert summary_response.status_code == 200
    assert export_response.status_code == 200
    assert retention_execution_response.status_code == 200
    assert execute_privacy_response.status_code == 200
    assert execute_retention_response.status_code == 200
    assert cancel_response.json()["status"] == "cancelled"
    assert detail_response.json()["providerKey"] == "local"
    assert privacy_transition_response.json()["status"] == "processing"
    assert privacy_detail_response.json()["requestType"] == "anonymization"
    assert consent_transition_response.json()["status"] == "revoked"
    assert summary_response.json()["summary"]["fiscalDocuments"] >= 1
    assert summary_response.json()["summary"]["consents"] >= 1
    assert summary_response.json()["summary"]["documentEvents"] >= 2
    assert "providerReady" in summary_response.json()["providers"]
    assert export_response.json()["summary"]["documents"] >= 1
    assert retention_execution_response.json()["summary"]["documentsTracked"] >= 1
    assert execute_privacy_response.json()["request"]["status"] == "completed"
    assert "execution" in execute_privacy_response.json()
    assert "summary" in execute_retention_response.json()


def test_capabilities_and_audit_routes_return_operational_payload() -> None:
    capabilities_response = client.get("/api/fiscal/capabilities")
    audit_response = client.get("/api/fiscal/audit-events?tenant_slug=bootstrap-ops")
    privacy_response = client.get("/api/fiscal/privacy-requests?tenant_slug=bootstrap-ops")
    consent_response = client.get("/api/fiscal/consents?tenant_slug=bootstrap-ops")

    assert capabilities_response.status_code == 200
    assert audit_response.status_code == 200
    assert privacy_response.status_code == 200
    assert consent_response.status_code == 200
    assert any(item["capabilityKey"] == "fiscal.nfe" for item in capabilities_response.json()["capabilities"])
    assert any(item["capabilityKey"] == "fiscal.certificate.a1" for item in capabilities_response.json()["capabilities"])
