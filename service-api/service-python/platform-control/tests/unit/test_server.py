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

    invalid_complete_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}/complete",
        json={"summary": "Cannot complete before start."},
    )
    start_response = client.post(f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}/start", json={"summary": "Started."})
    complete_response = client.post(f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}/complete", json={"summary": "Completed."})
    detail_response = client.get(f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{job_public_id}")
    detail_payload = detail_response.json()

    assert invalid_complete_response.status_code == 400
    assert invalid_complete_response.json()["detail"]["code"] == "lifecycle_job_transition_invalid"
    assert start_response.status_code == 200
    assert complete_response.status_code == 200
    assert detail_response.status_code == 200
    assert detail_payload["status"] == "completed"
    assert len(detail_payload["events"]) >= 3


def test_lifecycle_job_can_cancel_from_queue() -> None:
    created = client.post(
        "/api/platform-control/tenants/bootstrap-ops/lifecycle/onboarding",
        json={"requestedBy": "ops@bootstrap-ops.local", "payload": {"seed": "starter"}},
    ).json()

    cancel_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/lifecycle/jobs/{created['publicId']}/cancel",
        json={"summary": "Cancelled before execution."},
    )

    assert cancel_response.status_code == 200
    assert cancel_response.json()["status"] == "cancelled"
    assert any(event["status"] == "cancelled" for event in cancel_response.json()["events"])


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


def test_go_live_rollout_rejects_invalid_transition_and_supports_rollback() -> None:
    create_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/go-live/rollouts",
        json={
            "requestedBy": "ops@bootstrap-ops.local",
            "targetEnv": "production",
            "waveKey": "wave-rollback",
            "rollbackPlaybook": "docs/OPERACOES.md#rollback",
            "adoptionTargetPct": 75,
        },
    )
    rollout_public_id = create_response.json()["publicId"]

    invalid_complete_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/go-live/rollouts/{rollout_public_id}/complete",
        json={"summary": "Cannot complete before start."},
    )
    start_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/go-live/rollouts/{rollout_public_id}/start",
        json={"summary": "Wave started."},
    )
    rollback_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/go-live/rollouts/{rollout_public_id}/rollback",
        json={"summary": "Rollback executed from controlled wave."},
    )

    assert create_response.status_code == 200
    assert invalid_complete_response.status_code == 400
    assert invalid_complete_response.json()["detail"]["code"] == "go_live_rollout_transition_invalid"
    assert start_response.status_code == 200
    assert rollback_response.status_code == 200
    assert rollback_response.json()["status"] == "rolled_back"
    assert any(event["status"] == "rolled_back" for event in rollback_response.json()["events"])


def test_incident_command_center_supports_response_lifecycle() -> None:
    create_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/incidents",
        json={
            "title": "Gateway latency above SLO",
            "service": "edge",
            "severity": "sev2",
            "impact": "API consumers observe elevated p95 latency.",
            "owner": "sre@bootstrap-ops.local",
        },
    )
    incident_public_id = create_response.json()["publicId"]
    timeline_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/incidents/{incident_public_id}/timeline",
        json={"eventType": "mitigation", "summary": "Traffic shifted to fallback pool.", "actor": "sre@bootstrap-ops.local"},
    )
    action_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/incidents/{incident_public_id}/actions",
        json={"title": "Tune downstream timeout budget", "owner": "platform@bootstrap-ops.local"},
    )
    resolve_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/incidents/{incident_public_id}/resolve",
        json={"summary": "Latency returned to normal after pool rebalance.", "actor": "sre@bootstrap-ops.local"},
    )
    postmortem_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/incidents/{incident_public_id}/postmortem",
        json={
            "rootCause": "Capacity skew in fallback pool.",
            "impactSummary": "Short p95 latency breach.",
            "preventiveActions": ["Add pool skew alert"],
        },
    )
    detail_response = client.get(f"/api/platform-control/tenants/bootstrap-ops/incidents/{incident_public_id}")
    readiness_response = client.get("/api/platform-control/tenants/bootstrap-ops/incident-command/readiness")

    assert create_response.status_code == 200
    assert timeline_response.status_code == 200
    assert action_response.status_code == 200
    assert resolve_response.status_code == 200
    assert postmortem_response.status_code == 200
    assert detail_response.status_code == 200
    assert readiness_response.status_code == 200
    assert detail_response.json()["status"] == "resolved"
    assert len(detail_response.json()["timeline"]) >= 4
    assert len(detail_response.json()["actions"]) == 1
    assert detail_response.json()["postmortem"]["rootCause"] == "Capacity skew in fallback pool."
    assert "incident-registry" in readiness_response.json()["controls"]


def test_autonomous_governance_policy_approval_runbook_timeline_and_evidence() -> None:
    policy_catalog = client.get("/api/platform-control/policies/catalog")
    decision_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/policies/evaluate",
        json={
            "domain": "search",
            "action": "data.export",
            "actor": "legal@bootstrap-ops.local",
            "context": {"legalHoldActive": True},
        },
    )
    approval_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/approvals",
        json={
            "domain": "platform-control",
            "commandType": "quota.change",
            "requestedBy": "ops@bootstrap-ops.local",
            "justification": "Temporary expansion for onboarding wave.",
            "commandPayload": {"metricKey": "documents.storage_bytes", "newLimitValue": 4096},
        },
    )
    approval_public_id = approval_response.json()["publicId"]
    approve_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/approvals/{approval_public_id}/approve",
        json={"actor": "owner@bootstrap-ops.local"},
    )
    execute_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/approvals/{approval_public_id}/execute",
        json={"actor": "ops@bootstrap-ops.local"},
    )
    runbook_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/runbooks",
        json={
            "runbookKey": "tenant-over-quota",
            "requestedBy": "ops@bootstrap-ops.local",
            "justification": "Tenant crossed quota threshold.",
        },
    )
    runbook_public_id = runbook_response.json()["publicId"]
    runbook_approval_id = runbook_response.json()["approvalPublicId"]
    client.post(
        f"/api/platform-control/tenants/bootstrap-ops/approvals/{runbook_approval_id}/approve",
        json={"actor": "owner@bootstrap-ops.local"},
    )
    start_runbook_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/runbooks/{runbook_public_id}/start",
        json={"actor": "ops@bootstrap-ops.local"},
    )
    complete_step_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/runbooks/{runbook_public_id}/complete-step",
        json={"actor": "ops@bootstrap-ops.local"},
    )
    timeline_response = client.get("/api/platform-control/tenants/bootstrap-ops/timeline?entity_type=runbook.run")
    evidence_response = client.get("/api/platform-control/tenants/bootstrap-ops/evidence")
    readiness_response = client.get("/api/platform-control/tenants/bootstrap-ops/autonomous-governance/readiness")

    assert policy_catalog.status_code == 200
    assert any(item["policyKey"] == "exports.require-review" for item in policy_catalog.json()["items"])
    assert decision_response.status_code == 200
    assert decision_response.json()["decision"] == "review"
    assert approval_response.status_code == 200
    assert approval_response.json()["status"] == "requested"
    assert approve_response.status_code == 200
    assert approve_response.json()["status"] == "approved"
    assert execute_response.status_code == 200
    assert execute_response.json()["status"] == "executed"
    assert runbook_response.status_code == 200
    assert runbook_response.json()["status"] == "waiting_approval"
    assert start_runbook_response.status_code == 200
    assert start_runbook_response.json()["status"] == "running"
    assert complete_step_response.status_code == 200
    assert any(step["status"] == "completed" for step in complete_step_response.json()["steps"])
    assert timeline_response.status_code == 200
    assert timeline_response.json()["items"]
    assert evidence_response.status_code == 200
    assert any(item["evidenceType"] == "policy-decision" for item in evidence_response.json()["items"])
    assert readiness_response.status_code == 200
    assert "audit-evidence-vault" in readiness_response.json()["controls"]


def test_enterprise_runtime_event_mesh_tenant_runtime_and_contract_evolution() -> None:
    catalog_response = client.get("/api/platform-control/event-mesh/catalog")
    event_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/event-mesh/events",
        json={
            "streamKey": "billing.invoice",
            "eventType": "billing.invoice.payment_failed",
            "producer": "billing",
            "status": "dead_letter",
            "reason": "consumer_timeout",
            "payload": {"invoicePublicId": "inv_123", "amountCents": 4900},
        },
    )
    dead_letters_response = client.get("/api/platform-control/tenants/bootstrap-ops/event-mesh/dead-letters")
    dead_letter_public_id = dead_letters_response.json()["items"][0]["publicId"]
    replay_response = client.post(
        f"/api/platform-control/tenants/bootstrap-ops/event-mesh/dead-letters/{dead_letter_public_id}/replay",
        json={"actor": "sre@bootstrap-ops.local"},
    )
    lineage_response = client.get("/api/platform-control/tenants/bootstrap-ops/event-mesh/lineage")

    profile_response = client.get("/api/platform-control/tenants/bootstrap-ops/runtime/profile")
    update_profile_response = client.put(
        "/api/platform-control/tenants/bootstrap-ops/runtime/profile",
        json={"actor": "ops@bootstrap-ops.local", "riskStatus": "stable", "modules": ["identity", "billing", "analytics"]},
    )
    quota_response = client.put(
        "/api/platform-control/tenants/bootstrap-ops/runtime/quotas/api.requests.daily",
        json={"limitValue": 300000, "usagePct": 41, "enforcementMode": "soft"},
    )
    health_response = client.get("/api/platform-control/tenants/bootstrap-ops/runtime/health-score")
    window_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/runtime/maintenance-windows",
        json={"title": "Schema registry migration", "createdBy": "ops@bootstrap-ops.local"},
    )

    contracts_response = client.get("/api/platform-control/contracts/evolution")
    snapshot_response = client.post(
        "/api/platform-control/contracts/evolution/snapshots",
        json={"contractKey": "platform-control.openapi", "version": "1.4.0", "payload": {"paths": ["/api/platform-control/event-mesh/catalog", "/api/platform-control/providers/activation/catalog"]}},
    )
    diff_response = client.post(
        "/api/platform-control/contracts/evolution/diffs",
        json={"contractKey": "platform-control.openapi", "fromVersion": "1.3.0", "toVersion": "1.4.0", "removedOperations": [], "changedSchemas": []},
    )
    matrix_response = client.get("/api/platform-control/contracts/evolution/compatibility-matrix")
    readiness_response = client.get("/api/platform-control/tenants/bootstrap-ops/enterprise-runtime/readiness")

    assert catalog_response.status_code == 200
    assert catalog_response.json()["summary"]["streams"] >= 8
    assert event_response.status_code == 200
    assert event_response.json()["payloadHash"]
    assert dead_letters_response.status_code == 200
    assert replay_response.status_code == 200
    assert replay_response.json()["status"] == "replayed"
    assert lineage_response.status_code == 200
    assert lineage_response.json()["summary"]["events"] >= 1
    assert profile_response.status_code == 200
    assert update_profile_response.status_code == 200
    assert update_profile_response.json()["modules"] == ["identity", "billing", "analytics"]
    assert quota_response.status_code == 200
    assert health_response.status_code == 200
    assert window_response.status_code == 200
    assert contracts_response.status_code == 200
    assert snapshot_response.status_code == 200
    assert diff_response.status_code == 200
    assert diff_response.json()["status"] == "compatible"
    assert matrix_response.status_code == 200
    assert readiness_response.status_code == 200
    assert "enterprise-event-mesh" in readiness_response.json()["controls"]


def test_external_provider_activation_catalog_and_missing_key_run() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)

    catalog_response = client.get("/api/platform-control/providers/activation/catalog")
    catalog = catalog_response.json()
    stripe = next(item for item in catalog["items"] if item["providerKey"] == "stripe")

    run_response = client.post(
        "/api/platform-control/tenants/bootstrap-ops/providers/activation/stripe/test",
        json={"actor": "ops@bootstrap-ops.local", "action": "connection_test"},
    )
    runs_response = client.get("/api/platform-control/tenants/bootstrap-ops/providers/activation/runs?provider_key=stripe")

    assert catalog_response.status_code == 200
    assert catalog["version"] == "1.3.0"
    assert stripe["credentialKey"] == "BILLING_STRIPE_SECRET_KEY"
    assert stripe["secretValueExposed"] is False
    assert run_response.status_code == 200
    assert run_response.json()["status"] == "unavailable"
    assert run_response.json()["credentialConfigured"] is False
    assert run_response.json()["secretValueExposed"] is False
    assert runs_response.status_code == 200
    assert runs_response.json()["items"][0]["providerKey"] == "stripe"


def test_external_intelligence_public_provider_activation() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)
    catalog_response = client.get("/api/platform-control/providers/activation/catalog")
    gdelt = next(item for item in catalog_response.json()["items"] if item["providerKey"] == "gdelt")
    viacep = next(item for item in catalog_response.json()["items"] if item["providerKey"] == "viacep")

    assert gdelt["credentialRequired"] is False
    assert gdelt["configured"] is True
    assert viacep["credentialRequired"] is False
    assert viacep["configured"] is True
