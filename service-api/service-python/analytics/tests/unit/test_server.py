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
