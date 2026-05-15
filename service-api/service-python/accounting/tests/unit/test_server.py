from fastapi.testclient import TestClient

from app.runtime import IN_MEMORY_STATE
from app.server import app


def test_accounting_capabilities_and_account_lifecycle() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)

    capabilities = client.get("/api/accounting/capabilities")
    assert capabilities.status_code == 200
    assert capabilities.json()["service"] == "accounting"

    created = client.post(
        "/api/accounting/accounts",
        json={
            "tenantSlug": "bootstrap-ops",
            "accountCode": "1.01.001",
            "accountName": "Caixa",
            "accountType": "asset",
            "normalBalance": "debit",
            "status": "active",
        },
    )
    assert created.status_code == 200
    payload = created.json()
    assert payload["accountCode"] == "1.01.001"
    assert payload["status"] == "active"

    summary = client.get("/api/accounting/statements/management-summary?tenant_slug=bootstrap-ops")
    assert summary.status_code == 200
    assert summary.json()["summary"]["accounts"]["active"] == 1


def test_accounting_posting_rules_create_immutable_balanced_ledger() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)
    client.post(
        "/api/accounting/accounts",
        json={"tenantSlug": "bootstrap-ops", "accountCode": "1.01", "accountName": "Caixa", "accountType": "asset"},
    )
    client.post(
        "/api/accounting/accounts",
        json={"tenantSlug": "bootstrap-ops", "accountCode": "4.01", "accountName": "Receita", "accountType": "revenue"},
    )
    rule = client.post(
        "/api/accounting/posting-rules",
        json={
            "tenantSlug": "bootstrap-ops",
            "ruleKey": "sales.invoice.issued",
            "description": "Venda faturada",
            "sourceService": "sales",
            "eventType": "invoice_issued",
            "debitAccountCode": "1.01",
            "creditAccountCode": "4.01",
            "status": "active",
        },
    )
    assert rule.status_code == 200

    posted = client.post(
        "/api/accounting/posting-rules/apply",
        json={"tenantSlug": "bootstrap-ops", "sourceService": "sales", "eventType": "invoice_issued", "amountCents": 12500},
    )
    assert posted.status_code == 200
    assert posted.json()["status"] == "posted"

    immutable = client.patch(
        f"/api/accounting/journal-entries/{posted.json()['publicId']}/status",
        json={"tenantSlug": "bootstrap-ops", "status": "cancelled"},
    )
    assert immutable.status_code == 400
    assert immutable.json()["detail"]["code"] == "journal_entry_immutable"

    ledger = client.get("/api/accounting/ledger?tenant_slug=bootstrap-ops")
    assert ledger.status_code == 200
    assert len(ledger.json()["lines"]) == 2

    dre = client.get("/api/accounting/statements/dre?tenant_slug=bootstrap-ops")
    assert dre.status_code == 200
    assert dre.json()["dre"]["revenueCents"] == 12500
