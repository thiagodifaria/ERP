from fastapi.testclient import TestClient

from app.runtime import IN_MEMORY_STATE
from app.server import app


def test_ai_governance_tools_policies_and_read_only_run() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)

    tools = client.get("/api/ai-governance/tools")
    assert tools.status_code == 200
    assert tools.json()["items"][0]["mode"] == "read"
    assert tools.json()["provider"]["providerKey"] == "openai"
    assert tools.json()["provider"]["configured"] is False
    assert tools.json()["provider"]["credentialKey"] == "OPENAI_API_KEY"

    run = client.post(
        "/api/ai-governance/runs",
        json={
            "tenantSlug": "bootstrap-ops",
            "actor": "analyst@erp.local",
            "prompt": "Resuma o lead cfo@example.com com token access_token=abc123",
            "tools": ["search.query", "analytics.metric.lookup", "write.invoice"],
        },
    )
    payload = run.json()
    assert run.status_code == 200
    assert payload["mode"] == "read-only"
    assert payload["status"] == "completed_with_denials"
    assert payload["provider"] == "deterministic"
    assert payload["model"] == "deterministic-local"
    assert "write.invoice" in payload["deniedTools"]
    assert payload["redactionFindings"] == ["identifier", "token"]

    detail = client.get(f"/api/ai-governance/runs/{payload['publicId']}?tenant_slug=bootstrap-ops")
    assert detail.status_code == 200


def test_ai_governance_redaction_preview_and_audit() -> None:
    client = TestClient(app)
    response = client.post("/api/ai-governance/redaction/preview", json={"tenantSlug": "bootstrap-ops", "text": "CPF 123.456.789-00 e email user@example.com"})

    assert response.status_code == 200
    assert "[REDACTED:identifier]" in response.json()["redactedText"]

    audit = client.get("/api/ai-governance/audit-events?tenant_slug=bootstrap-ops")
    assert audit.status_code == 200
