from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_capability_catalog_returns_saas_capabilities() -> None:
    response = client.get("/api/platform-control/capabilities/catalog")
    assert response.status_code == 200
    payload = response.json()
    assert any(item["capabilityKey"] == "catalog.items" for item in payload)
    assert any(item["capabilityKey"] == "documents.digital_signature" for item in payload)


def test_provider_defaults_and_readiness_expose_tenant_provider_governance() -> None:
    catalog_response = client.get("/api/platform-control/providers/catalog")
    put_response = client.put(
        "/api/platform-control/tenants/bootstrap-ops/provider-defaults/documents.digital_signature",
        json={"providerKey": "local", "mode": "fallback", "source": "tenant-bootstrap"},
    )
    list_response = client.get("/api/platform-control/tenants/bootstrap-ops/provider-defaults")
    readiness_response = client.get("/api/platform-control/tenants/bootstrap-ops/lifecycle/readiness")

    assert catalog_response.status_code == 200
    assert put_response.status_code == 200
    assert list_response.status_code == 200
    assert readiness_response.status_code == 200
    assert any(item["capabilityKey"] == "documents.digital_signature" for item in catalog_response.json())
    assert any(item["capabilityKey"] == "documents.digital_signature" for item in list_response.json()["items"])
    assert readiness_response.json()["providers"]["total"] >= 1


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


def test_feature_flags_routes_reuse_capability_governance() -> None:
    put_response = client.put(
        "/api/platform-control/tenants/bootstrap-ops/feature-flags/catalog.items",
        json={"enabled": True, "planCode": "growth", "limitValue": 100, "source": "feature-ops"},
    )
    list_response = client.get("/api/platform-control/tenants/bootstrap-ops/feature-flags?limit=10")

    assert put_response.status_code == 200
    assert list_response.status_code == 200
    payload = list_response.json()
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert any(item["flagKey"] == "catalog.items" for item in payload["items"])
    assert any(item["module"] == "catalog" for item in payload["items"])


def test_lifecycle_onboarding_returns_accepted_and_location() -> None:
    preview_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/lifecycle/onboarding/preview",
        json={"requestedBy": "ops@bootstrap-ops.local", "payload": {"seed": "starter"}},
    )
    response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/lifecycle/onboarding",
        json={"requestedBy": "ops@bootstrap-ops.local", "payload": {"seed": "starter"}},
        headers={"Idempotency-Key": "onboarding-bootstrap-ops-1"},
    )
    payload = response.json()

    assert preview_response.status_code == 200
    assert preview_response.json()["jobType"] == "onboarding"
    assert len(preview_response.json()["steps"]) >= 3
    assert response.status_code == 202
    assert response.headers["location"].endswith(payload["publicId"])
    assert payload["status"] == "queued"
    assert payload["idempotencyKey"] == "onboarding-bootstrap-ops-1"
    assert "preview" in payload["payload"]


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


def test_go_live_rollout_flow_returns_readiness_and_history() -> None:
    readiness_response = client.get("/api/platform-control/tenants/bootstrap-ops/go-live/readiness")
    client.put(
        "/api/platform-control/tenants/bootstrap-ops/provider-defaults/catalog.items",
        json={"providerKey": "internal", "mode": "manual", "source": "tenant-bootstrap"},
    )
    client.put(
        "/api/platform-control/tenants/bootstrap-ops/provider-defaults/documents.digital_signature",
        json={"providerKey": "local", "mode": "fallback", "source": "tenant-bootstrap"},
    )
    client.put(
        "/api/platform-control/tenants/bootstrap-ops/quotas/documents.storage_bytes",
        json={"metricUnit": "bytes", "limitValue": 2048, "enforcementMode": "soft", "source": "plan"},
    )
    client.post(
        "/api/platform-control/tenants/bootstrap-ops/metering/snapshots",
        json={"metricKey": "documents.storage_bytes", "metricUnit": "bytes", "quantity": 512, "source": "documents"},
    )
    create_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/go-live/rollouts",
        json={
            "requestedBy": "ops@bootstrap-ops.local",
            "targetEnv": "production",
            "waveKey": "wave-1",
            "rollbackPlaybook": "docs/OPERACOES.md#rollback",
            "adoptionTargetPct": 80,
        },
    )
    rollout_public_id = create_response.json()["publicId"]
    start_response = client.post(f"/api/platform-control/tenants/bootstrap-ops/go-live/rollouts/{rollout_public_id}/start", json={"summary": "Wave started."})
    complete_response = client.post(f"/api/platform-control/tenants/bootstrap-ops/go-live/rollouts/{rollout_public_id}/complete", json={"summary": "Wave completed."})
    detail_response = client.get(f"/api/platform-control/tenants/bootstrap-ops/go-live/rollouts/{rollout_public_id}")
    adoption_response = client.get("/api/platform-control/tenants/bootstrap-ops/go-live/adoption")
    bottlenecks_response = client.get("/api/platform-control/tenants/bootstrap-ops/go-live/bottlenecks")
    playbook_response = client.get("/api/platform-control/tenants/bootstrap-ops/go-live/playbook")
    adjustments_response = client.get("/api/platform-control/tenants/bootstrap-ops/go-live/adjustments")

    assert readiness_response.status_code == 200
    assert create_response.status_code == 200
    assert start_response.status_code == 200
    assert complete_response.status_code == 200
    assert detail_response.status_code == 200
    assert adoption_response.status_code == 200
    assert bottlenecks_response.status_code == 200
    assert playbook_response.status_code == 200
    assert adjustments_response.status_code == 200
    assert "rolloutReady" in readiness_response.json()
    assert detail_response.json()["status"] == "completed"
    assert len(detail_response.json()["events"]) >= 3
    assert "adoptionPct" in adoption_response.json()
    assert "items" in bottlenecks_response.json()
    assert len(playbook_response.json()["checklist"]) >= 4
    assert "items" in adjustments_response.json()


def test_go_live_adjustments_can_apply_quota_and_block_changes() -> None:
    client.put(
        "/api/platform-control/tenants/bootstrap-ops/quotas/messages.daily",
        json={"metricUnit": "messages", "limitValue": 10, "enforcementMode": "hard", "source": "plan"},
    )
    client.post(
        "/api/platform-control/tenants/bootstrap-ops/metering/snapshots",
        json={"metricKey": "messages.daily", "metricUnit": "messages", "quantity": 10, "source": "engagement"},
    )
    client.put(
        "/api/platform-control/tenants/bootstrap-ops/blocks/provider_manual_review",
        json={"active": True, "reason": "manual review", "scope": "tenant", "source": "ops"},
    )

    adjustments_response = client.get("/api/platform-control/tenants/bootstrap-ops/go-live/adjustments")
    quota_apply_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/go-live/adjustments/apply",
        json={
            "actionType": "increase_quota_limit",
            "metricKey": "messages.daily",
            "metricUnit": "messages",
            "newLimitValue": 20,
            "enforcementMode": "hard",
            "actor": "ops@bootstrap-ops.local",
        },
    )
    block_apply_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/go-live/adjustments/apply",
        json={
            "actionType": "disable_block",
            "blockKey": "provider_manual_review",
            "reason": "released",
            "scope": "tenant",
            "actor": "ops@bootstrap-ops.local",
        },
    )

    assert adjustments_response.status_code == 200
    assert quota_apply_response.status_code == 200
    assert block_apply_response.status_code == 200
    assert quota_apply_response.json()["result"]["limitValue"] == 20
    assert block_apply_response.json()["result"]["active"] is False
