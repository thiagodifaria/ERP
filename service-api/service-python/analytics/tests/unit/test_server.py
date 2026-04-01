"""Cobertura inicial do bootstrap e do primeiro relatorio publico."""

from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_live_route_returns_service_health() -> None:
    response = client.get("/health/live")

    assert response.status_code == 200
    assert response.json() == {"service": "analytics", "status": "live"}


def test_details_route_exposes_analytics_dependencies() -> None:
    response = client.get("/health/details")
    payload = response.json()

    assert response.status_code == 200
    assert payload["status"] == "ready"
    assert any(dependency["name"] == "report-engine" for dependency in payload["dependencies"])
    assert any(dependency["name"] == "warehouse" for dependency in payload["dependencies"])


def test_pipeline_summary_returns_operational_payload() -> None:
    response = client.get("/api/analytics/reports/pipeline-summary?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["metrics"]["leadsCaptured"] == 128
    assert payload["bySource"]["whatsapp"] == 46
    assert payload["backlog"]["runningAutomations"] == 7


def test_service_pulse_returns_cross_service_payload() -> None:
    response = client.get("/api/analytics/reports/service-pulse?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["services"]["crm"]["totalLeads"] == 128
    assert payload["services"]["workflowControl"]["activeDefinitions"] == 6
    assert payload["services"]["webhookHub"]["forwarded"] == 87


def test_tenant_360_returns_tenant_operational_snapshot() -> None:
    response = client.get("/api/analytics/reports/tenant-360?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["identity"]["companies"] == 3
    assert payload["commercial"]["assignedLeads"] == 96
    assert payload["automation"]["workflowRuns"] == 41


def test_automation_board_returns_delivery_and_runtime_board() -> None:
    response = client.get("/api/analytics/reports/automation-board?tenant_slug=bootstrap-ops")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["catalog"]["definitionsActive"] == 6
    assert payload["runtime"]["completedExecutions"] == 28
    assert payload["delivery"]["forwarded"] == 87
