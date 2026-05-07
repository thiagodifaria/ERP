from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_queue_case_and_summary_flow() -> None:
    queue_response = client.put(
        "/api/support/queues/backoffice",
        json={"tenantSlug": "bootstrap-ops", "name": "Backoffice", "slaTargetHours": 12, "active": True},
    )
    case_response = client.post(
        "/api/support/cases",
        json={
            "tenantSlug": "bootstrap-ops",
            "subject": "Customer cannot generate invoice.",
            "queueKey": "backoffice",
            "priority": "high",
            "ownerUserId": "owner-user",
            "sourceKind": "crm",
            "entityKind": "crm.customer",
            "entityPublicId": "customer-001",
        },
    )
    case_public_id = case_response.json()["publicId"]
    status_response = client.patch(
        f"/api/support/cases/{case_public_id}/status",
        json={"tenantSlug": "bootstrap-ops", "status": "in_progress", "summary": "Assigned to finance support."},
    )
    comment_response = client.post(
        f"/api/support/cases/{case_public_id}/comments",
        json={"tenantSlug": "bootstrap-ops", "message": "Waiting payment gateway callback evidence."},
    )
    summary_response = client.get("/api/support/cases/summary?tenant_slug=bootstrap-ops")

    assert queue_response.status_code == 200
    assert case_response.status_code == 200
    assert status_response.status_code == 200
    assert comment_response.status_code == 200
    assert summary_response.status_code == 200
    assert summary_response.json()["summary"]["total"] >= 1
    assert summary_response.json()["byPriority"]["high"] >= 1
    assert any(event["eventType"] == "comment" for event in comment_response.json()["events"])


def test_capabilities_and_case_listing_routes_return_operational_payload() -> None:
    capabilities_response = client.get("/api/support/capabilities")
    queues_response = client.get("/api/support/queues?tenant_slug=bootstrap-ops")
    list_response = client.get("/api/support/cases?tenant_slug=bootstrap-ops&limit=10")

    assert capabilities_response.status_code == 200
    assert queues_response.status_code == 200
    assert list_response.status_code == 200
    assert any(item["key"] == "support.sla" for item in capabilities_response.json()["capabilities"])
    assert len(queues_response.json()) >= 1
    assert "pageInfo" in list_response.json()


def test_bulk_and_export_routes_return_partial_success_shape() -> None:
    bulk_response = client.post(
        "/api/support/cases/bulk",
        json={
            "tenantSlug": "bootstrap-ops",
            "items": [
                {
                    "subject": "Need VAT adjustment.",
                    "queueKey": "billing",
                    "priority": "medium",
                    "sourceKind": "sales",
                    "entityKind": "sales.invoice",
                    "entityPublicId": "invoice-001",
                },
                {
                    "subject": "",
                    "queueKey": "technical",
                },
            ],
        },
    )
    export_response = client.get("/api/support/cases/export?tenant_slug=bootstrap-ops")

    assert bulk_response.status_code == 200
    assert bulk_response.json()["summary"]["requested"] == 2
    assert bulk_response.json()["summary"]["succeeded"] == 1
    assert bulk_response.json()["summary"]["failed"] == 1
    assert bulk_response.json()["summary"]["partialSuccess"] is True
    assert export_response.status_code == 200
    assert export_response.json()["exported"] >= 1
