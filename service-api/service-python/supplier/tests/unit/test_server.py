from fastapi.testclient import TestClient

from app.server import app


client = TestClient(app)


def test_supplier_category_and_supplier_flow() -> None:
    category_response = client.put(
        "/api/supplier/categories/software",
        json={"tenantSlug": "bootstrap-ops", "name": "Software", "active": True},
    )
    supplier_response = client.post(
        "/api/supplier/suppliers",
        json={
            "tenantSlug": "bootstrap-ops",
            "companyName": "ERP Providers Ltda",
            "tradeName": "ERP Providers",
            "taxId": "12.345.678/0001-90",
            "categoryKey": "software",
            "status": "active",
            "payableTermDays": 15,
            "bankName": "Banco Example",
            "pixKey": "financeiro@erp.test",
            "contactEmail": "financeiro@erp.test",
        },
    )
    supplier_public_id = supplier_response.json()["publicId"]
    patch_response = client.patch(
        f"/api/supplier/suppliers/{supplier_public_id}",
        json={"tenantSlug": "bootstrap-ops", "status": "watchlist", "payableTermDays": 30},
    )
    summary_response = client.get("/api/supplier/suppliers/summary?tenant_slug=bootstrap-ops")

    assert category_response.status_code == 200
    assert supplier_response.status_code == 200
    assert patch_response.status_code == 200
    assert summary_response.status_code == 200
    assert patch_response.json()["status"] == "watchlist"
    assert summary_response.json()["summary"]["suppliersTotal"] >= 1


def test_capabilities_and_list_routes_return_operational_payload() -> None:
    capabilities_response = client.get("/api/supplier/capabilities")
    categories_response = client.get("/api/supplier/categories?tenant_slug=bootstrap-ops")
    list_response = client.get("/api/supplier/suppliers?tenant_slug=bootstrap-ops")

    assert capabilities_response.status_code == 200
    assert categories_response.status_code == 200
    assert list_response.status_code == 200
    assert any(item["key"] == "supplier.directory" for item in capabilities_response.json()["capabilities"])


def test_bulk_and_export_routes_return_partial_success_shape() -> None:
    client.put(
        "/api/supplier/categories/software",
        json={"tenantSlug": "bootstrap-ops", "name": "Software", "active": True},
    )
    bulk_response = client.post(
        "/api/supplier/suppliers/bulk",
        json={
            "tenantSlug": "bootstrap-ops",
            "items": [
                {
                    "companyName": "Cloud Freight LTDA",
                    "taxId": "11.222.333/0001-44",
                    "categoryKey": "software",
                    "status": "active",
                    "payableTermDays": 20,
                },
                {
                    "companyName": "",
                    "taxId": "11.222.333/0001-45",
                },
            ],
        },
    )
    export_response = client.get("/api/supplier/suppliers/export?tenant_slug=bootstrap-ops")

    assert bulk_response.status_code == 200
    assert bulk_response.json()["summary"]["requested"] == 2
    assert bulk_response.json()["summary"]["succeeded"] == 1
    assert bulk_response.json()["summary"]["failed"] == 1
    assert bulk_response.json()["summary"]["partialSuccess"] is True
    assert export_response.status_code == 200
    assert export_response.json()["exported"] >= 1
