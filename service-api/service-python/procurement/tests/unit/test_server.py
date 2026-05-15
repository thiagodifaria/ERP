from fastapi.testclient import TestClient

from app.runtime import IN_MEMORY_STATE
from app.server import app


def test_procurement_capabilities_and_requisition_lifecycle() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)

    capabilities = client.get("/api/procurement/capabilities")
    assert capabilities.status_code == 200
    assert capabilities.json()["service"] == "procurement"

    created = client.post(
        "/api/procurement/requisitions",
        json={
            "tenantSlug": "bootstrap-ops",
            "requisitionNumber": "REQ-001",
            "title": "Compra de notebook",
            "requestedBy": "ops",
            "status": "approved",
        },
    )
    assert created.status_code == 200
    payload = created.json()
    assert payload["requisitionNumber"] == "REQ-001"
    assert payload["status"] == "approved"

    summary = client.get("/api/procurement/matching/summary?tenant_slug=bootstrap-ops")
    assert summary.status_code == 200
    assert summary.json()["summary"]["requisitions"]["active"] == 1


def test_procurement_approval_and_three_way_matching() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)
    order = client.post(
        "/api/procurement/purchase-orders",
        json={
            "tenantSlug": "bootstrap-ops",
            "orderNumber": "PO-001",
            "description": "Compra aprovada",
            "supplierPublicId": "supplier-1",
            "amountCents": 20000,
        },
    )
    assert order.status_code == 200

    approval = client.post(
        "/api/procurement/approvals/apply",
        json={
            "tenantSlug": "bootstrap-ops",
            "targetCollection": "purchase-orders",
            "targetPublicId": order.json()["publicId"],
            "approvedBy": "manager",
            "status": "approved",
        },
    )
    assert approval.status_code == 200
    assert approval.json()["target"]["status"] == "approved"

    receipt = client.post(
        "/api/procurement/receipts",
        json={
            "tenantSlug": "bootstrap-ops",
            "receiptNumber": "REC-001",
            "description": "Recebimento",
            "purchaseOrderPublicId": order.json()["publicId"],
            "amountCents": 20000,
        },
    )
    match = client.post(
        "/api/procurement/matching/three-way",
        json={
            "tenantSlug": "bootstrap-ops",
            "purchaseOrderPublicId": order.json()["publicId"],
            "receiptPublicId": receipt.json()["publicId"],
            "fiscalDocumentPublicId": "fiscal-doc-1",
            "orderAmountCents": 20000,
            "receiptAmountCents": 20000,
            "invoiceAmountCents": 20000,
        },
    )
    assert match.status_code == 200
    assert match.json()["status"] == "matched"
