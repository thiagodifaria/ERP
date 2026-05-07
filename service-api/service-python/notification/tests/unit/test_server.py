from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_preferences_and_notification_center_flow() -> None:
    preference_response = client.put(
        "/api/notification/preferences/user-001",
        json={"tenantSlug": "bootstrap-ops", "inAppEnabled": True, "emailEnabled": True, "quietHours": {"from": "23:00", "to": "06:00"}},
    )
    create_response = client.post(
        "/api/notification/center",
        json={
            "tenantSlug": "bootstrap-ops",
            "userPublicId": "user-001",
            "title": "Payment recovery needs review",
            "body": "A billing recovery case escalated to urgent.",
            "severity": "critical",
            "channel": "email",
            "sourceModule": "billing",
            "entityKind": "billing.recovery_case",
            "entityPublicId": "recovery-001",
        },
    )
    notification_public_id = create_response.json()["publicId"]
    transition_response = client.patch(
        f"/api/notification/center/{notification_public_id}/status",
        json={"tenantSlug": "bootstrap-ops", "status": "read"},
    )
    summary_response = client.get("/api/notification/summary?tenant_slug=bootstrap-ops")

    assert preference_response.status_code == 200
    assert create_response.status_code == 200
    assert transition_response.status_code == 200
    assert summary_response.status_code == 200
    assert transition_response.json()["status"] == "read"
    assert summary_response.json()["summary"]["total"] >= 1


def test_capabilities_and_center_listing_routes_return_operational_payload() -> None:
    capabilities_response = client.get("/api/notification/capabilities")
    center_response = client.get("/api/notification/center?tenant_slug=bootstrap-ops")

    assert capabilities_response.status_code == 200
    assert center_response.status_code == 200
    assert any(item["key"] == "notification.center" for item in capabilities_response.json()["capabilities"])
    assert "pageInfo" in center_response.json()


def test_bulk_notification_route_returns_partial_success_shape() -> None:
    bulk_response = client.post(
        "/api/notification/center/bulk",
        json={
            "tenantSlug": "bootstrap-ops",
            "items": [
                {
                    "userPublicId": "user-001",
                    "title": "Support case escalated",
                    "body": "Case moved to urgent queue.",
                    "severity": "warning",
                    "channel": "in_app",
                    "sourceModule": "support",
                    "entityKind": "support.case",
                    "entityPublicId": "case-001",
                },
                {
                    "userPublicId": "user-002",
                    "title": "",
                    "body": "Invalid notification example.",
                    "severity": "critical",
                    "channel": "email",
                    "sourceModule": "support",
                },
            ],
        },
    )

    assert bulk_response.status_code == 200
    assert bulk_response.json()["summary"]["requested"] == 2
    assert bulk_response.json()["summary"]["succeeded"] == 1
    assert bulk_response.json()["summary"]["failed"] == 1
    assert bulk_response.json()["summary"]["partialSuccess"] is True
