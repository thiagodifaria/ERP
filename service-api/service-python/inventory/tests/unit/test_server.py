from fastapi.testclient import TestClient

from app.runtime import IN_MEMORY_STATE
from app.server import app


def test_inventory_capabilities_and_location_lifecycle() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)

    capabilities = client.get("/api/inventory/capabilities")
    assert capabilities.status_code == 200
    assert capabilities.json()["service"] == "inventory"

    created = client.post(
        "/api/inventory/locations",
        json={
            "tenantSlug": "bootstrap-ops",
            "locationCode": "MAIN",
            "locationName": "Deposito Principal",
            "warehouseCode": "WH-01",
            "status": "active",
        },
    )
    assert created.status_code == 200
    payload = created.json()
    assert payload["locationCode"] == "MAIN"
    assert payload["status"] == "active"

    summary = client.get("/api/inventory/summary?tenant_slug=bootstrap-ops")
    assert summary.status_code == 200
    assert summary.json()["summary"]["locations"]["active"] == 1


def test_inventory_balances_costing_and_cycle_count_variance() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)
    client.post(
        "/api/inventory/movements",
        json={
            "tenantSlug": "bootstrap-ops",
            "movementNumber": "MOV-IN",
            "reason": "purchase_receipt",
            "sku": "SKU-1",
            "locationCode": "MAIN",
            "quantity": 10,
            "movementType": "in",
            "unitCostCents": 1000,
        },
    )
    client.post(
        "/api/inventory/reservations",
        json={
            "tenantSlug": "bootstrap-ops",
            "reservationNumber": "RES-1",
            "sku": "SKU-1",
            "locationCode": "MAIN",
            "quantity": 3,
            "status": "active",
        },
    )
    client.post(
        "/api/inventory/cost-layers",
        json={
            "tenantSlug": "bootstrap-ops",
            "layerNumber": "L1",
            "sku": "SKU-1",
            "locationCode": "MAIN",
            "quantity": 10,
            "remainingQuantity": 10,
            "unitCostCents": 1000,
        },
    )
    client.post(
        "/api/inventory/cycle-counts",
        json={
            "tenantSlug": "bootstrap-ops",
            "countNumber": "CNT-1",
            "sku": "SKU-1",
            "locationCode": "MAIN",
            "countedQuantity": 9,
            "expectedQuantity": 10,
        },
    )

    balances = client.get("/api/inventory/balances?tenant_slug=bootstrap-ops&sku=SKU-1")
    assert balances.status_code == 200
    assert balances.json()["balances"][0]["available"] == 7

    costing = client.get("/api/inventory/costing/summary?tenant_slug=bootstrap-ops&sku=SKU-1")
    assert costing.status_code == 200
    assert costing.json()["summary"]["averageUnitCostCents"] == 1000

    variances = client.get("/api/inventory/cycle-counts/variances?tenant_slug=bootstrap-ops")
    assert variances.status_code == 200
    assert variances.json()["variances"][0]["varianceQuantity"] == -1
