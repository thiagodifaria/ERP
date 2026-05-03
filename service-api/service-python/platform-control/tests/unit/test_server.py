from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_capability_catalog_returns_saas_capabilities() -> None:
    response = client.get("/api/platform-control/capabilities/catalog")
    assert response.status_code == 200
    payload = response.json()
    assert any(item["capabilityKey"] == "catalog.items" for item in payload)
