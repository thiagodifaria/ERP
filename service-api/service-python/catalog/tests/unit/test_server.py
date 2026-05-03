from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_capabilities_returns_catalog_shape() -> None:
    response = client.get("/api/catalog/capabilities")
    assert response.status_code == 200
    payload = response.json()
    assert payload["supportsCategories"] is True
    assert "product" in payload["domains"]
