"""Cobertura inicial do bootstrap e do catalogo do simulation."""

from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_live_route_returns_service_health() -> None:
    response = client.get("/health/live")

    assert response.status_code == 200
    assert response.json() == {"service": "simulation", "status": "live"}


def test_details_route_exposes_simulation_dependencies() -> None:
    response = client.get("/health/details")
    payload = response.json()

    assert response.status_code == 200
    assert payload["status"] == "ready"
    assert any(dependency["name"] == "scenario-engine" for dependency in payload["dependencies"])
    assert any(dependency["name"] == "scenario-store" for dependency in payload["dependencies"])
    assert any(dependency["name"] == "cost-model" for dependency in payload["dependencies"])


def test_catalog_route_returns_supported_scenarios() -> None:
    response = client.get("/api/simulation/scenarios/catalog")
    payload = response.json()

    assert response.status_code == 200
    assert payload["service"] == "simulation"
    assert payload["scenarioTypes"][0]["key"] == "operational-load"
    assert "leadMultiplier" in payload["scenarioTypes"][0]["inputs"]


def test_operational_load_route_returns_projection_payload() -> None:
    response = client.post(
        "/api/simulation/scenarios/operational-load",
        json={
            "tenantSlug": "bootstrap-ops",
            "leadMultiplier": 1.8,
            "automationMultiplier": 1.4,
            "webhookMultiplier": 1.2,
            "teamSize": 4,
            "planningHorizonDays": 30,
            "storageGrowthGb": 2,
        },
    )
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["projection"]["monthlyOperationsProjected"] > 0
    assert payload["projection"]["costBreakdown"]["estimatedMonthlyCostCents"] > 0


def test_load_benchmark_route_returns_capacity_payload() -> None:
    response = client.post(
        "/api/simulation/benchmarks/load",
        json={
            "tenantSlug": "bootstrap-ops",
            "leadMultiplier": 1.8,
            "automationMultiplier": 1.4,
            "webhookMultiplier": 1.2,
            "sampleSize": 120,
        },
    )
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["results"]["sampleSize"] == 120
    assert payload["results"]["avgLatencyMs"] > 0
    assert payload["results"]["throughputRps"] > 0
