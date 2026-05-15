"""Bootstrap HTTP do servico accounting."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import build_financial_statement, build_general_ledger, build_summary, capability_catalog, create_record, get_record, list_records, post_source_event, transition_record

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
    dependencies = [{"name": "accounting-runtime", "status": "ready"}]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/accounting/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/accounting/statements/management-summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)


@app.get("/api/accounting/ledger")
def ledger(tenant_slug: str | None = None, account_code: str | None = None) -> dict:
    return build_general_ledger(tenant_slug, account_code)


@app.get("/api/accounting/statements/{statement_kind}")
def statement(statement_kind: str, tenant_slug: str | None = None) -> dict:
    return build_financial_statement(statement_kind, tenant_slug)


@app.post("/api/accounting/posting-rules/apply")
def apply_posting_rule(payload: dict) -> dict:
    try:
        return post_source_event(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "posting_rule_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "accounting posting payload is invalid."}) from error



@app.get("/api/accounting/cost-centers")
def list_cost_centers(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("cost-centers", tenant_slug, status)


@app.post("/api/accounting/cost-centers")
def post_cost_center(payload: dict) -> dict:
    try:
        return create_record("cost-centers", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "accounting payload is invalid."}) from error


@app.get("/api/accounting/cost-centers/{public_id}")
def get_cost_center(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("cost-centers", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "cost_center_not_found", "message": "Cost Center was not found."})
    return record


@app.patch("/api/accounting/cost-centers/{public_id}/status")
def patch_cost_center_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("cost-centers", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "accounting status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "cost_center_not_found", "message": "Cost Center was not found."})
    return record



@app.get("/api/accounting/accounts")
def list_accounts(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("accounts", tenant_slug, status)


@app.post("/api/accounting/accounts")
def post_account(payload: dict) -> dict:
    try:
        return create_record("accounts", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "accounting payload is invalid."}) from error


@app.get("/api/accounting/accounts/{public_id}")
def get_account(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("accounts", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "account_not_found", "message": "Account was not found."})
    return record


@app.patch("/api/accounting/accounts/{public_id}/status")
def patch_account_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("accounts", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "accounting status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "account_not_found", "message": "Account was not found."})
    return record



@app.get("/api/accounting/journal-entries")
def list_journal_entries(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("journal-entries", tenant_slug, status)


@app.post("/api/accounting/journal-entries")
def post_journal_entry(payload: dict) -> dict:
    try:
        return create_record("journal-entries", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "accounting payload is invalid."}) from error


@app.get("/api/accounting/journal-entries/{public_id}")
def get_journal_entry(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("journal-entries", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "journal_entry_not_found", "message": "Journal Entry was not found."})
    return record


@app.patch("/api/accounting/journal-entries/{public_id}/status")
def patch_journal_entry_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("journal-entries", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "accounting status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "journal_entry_not_found", "message": "Journal Entry was not found."})
    return record



@app.get("/api/accounting/posting-rules")
def list_posting_rules(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("posting-rules", tenant_slug, status)


@app.post("/api/accounting/posting-rules")
def post_posting_rule(payload: dict) -> dict:
    try:
        return create_record("posting-rules", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "accounting payload is invalid."}) from error


@app.get("/api/accounting/posting-rules/{public_id}")
def get_posting_rule(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("posting-rules", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "posting_rule_not_found", "message": "Posting Rule was not found."})
    return record


@app.patch("/api/accounting/posting-rules/{public_id}/status")
def patch_posting_rule_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("posting-rules", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "accounting status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "posting_rule_not_found", "message": "Posting Rule was not found."})
    return record



@app.get("/api/accounting/period-closes")
def list_period_closes(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("period-closes", tenant_slug, status)


@app.post("/api/accounting/period-closes")
def post_period_close(payload: dict) -> dict:
    try:
        return create_record("period-closes", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "accounting payload is invalid."}) from error


@app.get("/api/accounting/period-closes/{public_id}")
def get_period_close(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("period-closes", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "period_close_not_found", "message": "Period Close was not found."})
    return record


@app.patch("/api/accounting/period-closes/{public_id}/status")
def patch_period_close_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("period-closes", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "accounting status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "period_close_not_found", "message": "Period Close was not found."})
    return record
