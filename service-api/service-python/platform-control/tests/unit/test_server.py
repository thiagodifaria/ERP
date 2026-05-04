from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_capability_catalog_returns_saas_capabilities() -> None:
    response = client.get("/api/platform-control/capabilities/catalog")
    assert response.status_code == 200
    payload = response.json()
    assert any(item["capabilityKey"] == "catalog.items" for item in payload)
    assert any(item["capabilityKey"] == "documents.digital_signature" for item in payload)


def test_entitlements_bulk_returns_partial_success_shape() -> None:
    response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/entitlements/bulk",
        json={
            "items": [
                {"capabilityKey": "catalog.items", "enabled": True, "planCode": "growth", "limitValue": 100},
                {"enabled": False},
            ]
        },
    )
    payload = response.json()

    assert response.status_code == 200
    assert payload["summary"]["requested"] == 2
    assert payload["summary"]["succeeded"] == 1
    assert payload["summary"]["failed"] == 1
    assert payload["summary"]["partialSuccess"] is True


def test_lifecycle_onboarding_returns_accepted_and_location() -> None:
    response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/lifecycle/onboarding",
        json={"requestedBy": "ops@bootstrap-ops.local", "payload": {"seed": "starter"}},
        headers={"Idempotency-Key": "onboarding-bootstrap-ops-1"},
    )
    payload = response.json()

    assert response.status_code == 202
    assert response.headers["location"].endswith(payload["publicId"])
    assert payload["status"] == "queued"
    assert payload["idempotencyKey"] == "onboarding-bootstrap-ops-1"


def test_quota_and_usage_summary_routes_return_operational_payload() -> None:
    quota_response = client.put(
        "/api/platform-control/tenants/bootstrap-ops/quotas/documents.storage_bytes",
        json={"metricUnit": "bytes", "limitValue": 2048, "enforcementMode": "soft", "source": "plan"},
    )
    snapshot_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/metering/snapshots",
        json={"metricKey": "documents.storage_bytes", "metricUnit": "bytes", "quantity": 1024, "source": "documents"},
    )
    summary_response = client.get("/api/platform-control/tenants/bootstrap-ops/usage-summary")

    assert quota_response.status_code == 200
    assert snapshot_response.status_code == 200
    summary_payload = summary_response.json()
    assert summary_response.status_code == 200
    assert summary_payload["summary"]["trackedMetrics"] >= 1
    assert summary_payload["metrics"][0]["metricKey"] == "documents.storage_bytes"


def test_lifecycle_job_transitions_append_audit_events() -> None:
    created = client.post(
        "/api/platform-control/tenants/bootstrap-ops/lifecycle/offboarding",
        json={"requestedBy": "ops@bootstrap-ops.local", "payload": {"mode": "export-only"}},
    ).json()
    job_public_id = created["publicId"]

    start_response = client.post(f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}/start", json={"summary": "Started."})
    complete_response = client.post(f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}/complete", json={"summary": "Completed."})
    detail_response = client.get(f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}")
    detail_payload = detail_response.json()

    assert start_response.status_code == 200
    assert complete_response.status_code == 200
    assert detail_response.status_code == 200
    assert detail_payload["status"] == "completed"
    assert len(detail_payload["events"]) >= 3
