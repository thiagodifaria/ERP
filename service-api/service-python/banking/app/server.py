"""Bootstrap HTTP do servico banking."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import build_summary, capability_catalog, create_record, get_record, list_records, parse_cnab_return, reconcile_statement, transition_record

app = FastAPI(title=settings.service_name)
install_security_middleware(app, settings.service_name)


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    dependencies = [{"name": "banking-runtime", "status": "ready"}]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/banking/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/banking/reconciliation/summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)


@app.post("/api/banking/cnab-files/parse-return")
def post_cnab_parse_return(payload: dict) -> dict:
    try:
        return parse_cnab_return(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking CNAB payload is invalid."}) from error


@app.post("/api/banking/reconciliations/run")
def post_reconciliation_run(payload: dict) -> dict:
    try:
        return reconcile_statement(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking reconciliation payload is invalid."}) from error



@app.get("/api/banking/bank-accounts")
def list_bank_accounts(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("bank-accounts", tenant_slug, status)


@app.post("/api/banking/bank-accounts")
def post_bank_account(payload: dict) -> dict:
    try:
        return create_record("bank-accounts", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/bank-accounts/{public_id}")
def get_bank_account(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("bank-accounts", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "bank_account_not_found", "message": "Bank Account was not found."})
    return record


@app.patch("/api/banking/bank-accounts/{public_id}/status")
def patch_bank_account_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("bank-accounts", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "bank_account_not_found", "message": "Bank Account was not found."})
    return record



@app.get("/api/banking/cnab-files")
def list_cnab_files(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("cnab-files", tenant_slug, status)


@app.post("/api/banking/cnab-files")
def post_cnab_file(payload: dict) -> dict:
    try:
        return create_record("cnab-files", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/cnab-files/{public_id}")
def get_cnab_file(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("cnab-files", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "cnab_file_not_found", "message": "Cnab File was not found."})
    return record


@app.patch("/api/banking/cnab-files/{public_id}/status")
def patch_cnab_file_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("cnab-files", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "cnab_file_not_found", "message": "Cnab File was not found."})
    return record



@app.get("/api/banking/bank-statements")
def list_bank_statements(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("bank-statements", tenant_slug, status)


@app.post("/api/banking/bank-statements")
def post_bank_statement(payload: dict) -> dict:
    try:
        return create_record("bank-statements", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/bank-statements/{public_id}")
def get_bank_statement(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("bank-statements", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "bank_statement_not_found", "message": "Bank Statement was not found."})
    return record



@app.get("/api/banking/boletos")
def list_boletos(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("boletos", tenant_slug, status)


@app.post("/api/banking/boletos")
def post_boleto(payload: dict) -> dict:
    try:
        return create_record("boletos", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/boletos/{public_id}")
def get_boleto(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("boletos", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "boleto_not_found", "message": "Boleto was not found."})
    return record


@app.patch("/api/banking/boletos/{public_id}/status")
def patch_boleto_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("boletos", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "boleto_not_found", "message": "Boleto was not found."})
    return record



@app.get("/api/banking/pix-charges")
def list_pix_charges(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("pix-charges", tenant_slug, status)


@app.post("/api/banking/pix-charges")
def post_pix_charge(payload: dict) -> dict:
    try:
        return create_record("pix-charges", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/pix-charges/{public_id}")
def get_pix_charge(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("pix-charges", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "pix_charge_not_found", "message": "Pix Charge was not found."})
    return record


@app.patch("/api/banking/pix-charges/{public_id}/status")
def patch_pix_charge_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("pix-charges", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "pix_charge_not_found", "message": "Pix Charge was not found."})
    return record



@app.get("/api/banking/pix-refunds")
def list_pix_refunds(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("pix-refunds", tenant_slug, status)


@app.post("/api/banking/pix-refunds")
def post_pix_refund(payload: dict) -> dict:
    try:
        return create_record("pix-refunds", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/pix-webhooks")
def list_pix_webhooks(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("pix-webhooks", tenant_slug, status)


@app.post("/api/banking/pix-webhooks")
def post_pix_webhook(payload: dict) -> dict:
    try:
        return create_record("pix-webhooks", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/open-finance-connections")
def list_open_finance_connections(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("open-finance-connections", tenant_slug, status)


@app.post("/api/banking/open-finance-connections")
def post_open_finance_connection(payload: dict) -> dict:
    try:
        return create_record("open-finance-connections", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error



@app.get("/api/banking/reconciliations")
def list_reconciliations(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("reconciliations", tenant_slug, status)


@app.post("/api/banking/reconciliations")
def post_reconciliation(payload: dict) -> dict:
    try:
        return create_record("reconciliations", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "banking payload is invalid."}) from error


@app.get("/api/banking/reconciliations/{public_id}")
def get_reconciliation(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("reconciliations", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "reconciliation_not_found", "message": "Reconciliation was not found."})
    return record


@app.patch("/api/banking/reconciliations/{public_id}/status")
def patch_reconciliation_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("reconciliations", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "banking status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "reconciliation_not_found", "message": "Reconciliation was not found."})
    return record
