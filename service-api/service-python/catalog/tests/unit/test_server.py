from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_capabilities_returns_catalog_shape() -> None:
    response = client.get("/api/catalog/capabilities")
    assert response.status_code == 200
    payload = response.json()
    assert payload["supportsCategories"] is True
    assert payload["supportsCursorPagination"] is True
    assert payload["supportsBulk"] is True
    assert "product" in payload["domains"]


def test_items_bulk_returns_partial_success_shape() -> None:
    response = client.post(
        "/api/catalog/items/bulk",
        json={
            "tenantSlug": "bootstrap-ops",
            "items": [
                {
                    "tenantSlug": "bootstrap-ops",
                    "sku": "ERP-001",
                    "name": "ERP Core",
                    "itemType": "service",
                    "unitCode": "license",
                    "priceBaseCents": 9900,
                },
                {
                    "tenantSlug": "bootstrap-ops",
                    "name": "Broken Item",
                    "itemType": "invalid",
                    "unitCode": "unit",
                },
            ],
        },
    )
    payload = response.json()

    assert response.status_code == 200
    assert payload["summary"]["requested"] == 2
    assert payload["summary"]["succeeded"] == 1
    assert payload["summary"]["failed"] == 1


def test_categories_page_returns_cursor_shape() -> None:
    client.post("/api/catalog/categories", json={"tenantSlug": "bootstrap-ops", "key": "licenses", "name": "Licenses"})
    response = client.get("/api/catalog/categories/page?tenant_slug=bootstrap-ops&limit=10")
    payload = response.json()

    assert response.status_code == 200
    assert payload["tenantSlug"] == "bootstrap-ops"
    assert payload["pageInfo"]["returned"] >= 1
