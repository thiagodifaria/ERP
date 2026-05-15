from fastapi.testclient import TestClient

from app.runtime import IN_MEMORY_STATE
from app.server import app


def test_banking_capabilities_and_boleto_lifecycle() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)

    capabilities = client.get("/api/banking/capabilities")
    assert capabilities.status_code == 200
    assert capabilities.json()["service"] == "banking"

    created = client.post(
        "/api/banking/boletos",
        json={
            "tenantSlug": "bootstrap-ops",
            "boletoNumber": "BOL-001",
            "payerName": "Cliente Teste",
            "amountCents": 150000,
            "status": "active",
        },
    )
    assert created.status_code == 200
    payload = created.json()
    assert payload["boletoNumber"] == "BOL-001"
    assert payload["status"] == "active"

    summary = client.get("/api/banking/reconciliation/summary?tenant_slug=bootstrap-ops")
    assert summary.status_code == 200
    assert summary.json()["summary"]["boletos"]["active"] == 1


def test_banking_cnab_statement_pix_refund_webhook_and_open_finance() -> None:
    for records in IN_MEMORY_STATE.values():
        records.clear()

    client = TestClient(app)
    statement = client.post(
        "/api/banking/bank-statements",
        json={
            "tenantSlug": "bootstrap-ops",
            "statementNumber": "STM-001",
            "bankAccountPublicId": "bank-account-1",
            "statementDate": "2026-05-15",
            "amountCents": 30000,
            "transactionCount": 1,
        },
    )
    assert statement.status_code == 200

    cnab = client.post(
        "/api/banking/cnab-files/parse-return",
        json={"tenantSlug": "bootstrap-ops", "fileNumber": "RET-001", "rawContent": "0" * 225 + "000000000030000"},
    )
    assert cnab.status_code == 200
    assert cnab.json()["status"] == "processed"

    refund = client.post(
        "/api/banking/pix-refunds",
        json={"tenantSlug": "bootstrap-ops", "refundId": "RF-001", "txid": "TX-1", "amountCents": 5000, "status": "queued"},
    )
    webhook = client.post(
        "/api/banking/pix-webhooks",
        json={"tenantSlug": "bootstrap-ops", "eventId": "EVT-1", "txid": "TX-1", "eventType": "pix.received", "status": "processed"},
    )
    open_finance = client.post(
        "/api/banking/open-finance-connections",
        json={"tenantSlug": "bootstrap-ops", "connectionId": "OF-1", "bankName": "Banco Demo", "status": "active"},
    )
    assert refund.status_code == 200
    assert webhook.status_code == 200
    assert open_finance.status_code == 200

    reconciliation = client.post(
        "/api/banking/reconciliations/run",
        json={
            "tenantSlug": "bootstrap-ops",
            "statementPublicId": statement.json()["publicId"],
            "expectedAmountCents": 30000,
        },
    )
    assert reconciliation.status_code == 200
    assert reconciliation.json()["status"] == "matched"
