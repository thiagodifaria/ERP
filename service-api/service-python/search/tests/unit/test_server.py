from fastapi.testclient import TestClient

from app.runtime import IN_MEMORY_STATE
from app.server import app


def test_search_query_redacts_sensitive_data_and_audits() -> None:
    IN_MEMORY_STATE["query_audit_events"].clear()
    client = TestClient(app)

    response = client.get("/api/search/query?tenant_slug=bootstrap-ops&q=northwind&actor=analyst")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["total"] >= 1
    assert payload["items"][0]["redaction"]["applied"] is True
    assert payload["auditEvent"]["actor"] == "analyst"

    audit = client.get("/api/search/audit-events?tenant_slug=bootstrap-ops")
    assert audit.status_code == 200
    assert audit.json()["items"][0]["query"] == "northwind"


def test_search_discovery_case_legal_hold_and_export() -> None:
    client = TestClient(app)

    created = client.post("/api/search/discovery-cases", json={"tenantSlug": "bootstrap-ops", "title": "Auditoria fiscal", "owner": "legal"})
    assert created.status_code == 200
    case_id = created.json()["publicId"]

    item = client.post(
        f"/api/search/discovery-cases/{case_id}/items",
        json={"tenantSlug": "bootstrap-ops", "entityType": "documents.attachment", "entityPublicId": "doc-contract-001"},
    )
    assert item.status_code == 200

    hold = client.post("/api/search/legal-holds", json={"tenantSlug": "bootstrap-ops", "reason": "Fiscal audit"})
    assert hold.status_code == 200

    export = client.post("/api/search/exports", json={"tenantSlug": "bootstrap-ops", "query": "contract", "requestedBy": "legal"})
    assert export.status_code == 200
    assert export.json()["status"] == "held"
    assert export.json()["legalHoldCount"] >= 1

