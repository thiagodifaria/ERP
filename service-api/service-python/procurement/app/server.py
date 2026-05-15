"""Bootstrap HTTP do servico procurement."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import approve_target, build_summary, capability_catalog, create_record, get_record, list_records, run_three_way_match, transition_record

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
    dependencies = [{"name": "procurement-runtime", "status": "ready"}]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/procurement/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/procurement/matching/summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)


@app.post("/api/procurement/approvals/apply")
def apply_approval(payload: dict) -> dict:
    try:
        return approve_target(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "approval_target_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "procurement approval payload is invalid."}) from error


@app.post("/api/procurement/matching/three-way")
def post_three_way_match(payload: dict) -> dict:
    try:
        return run_three_way_match(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "procurement matching payload is invalid."}) from error



@app.get("/api/procurement/requisitions")
def list_requisitions(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("requisitions", tenant_slug, status)


@app.post("/api/procurement/requisitions")
def post_requisition(payload: dict) -> dict:
    try:
        return create_record("requisitions", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "procurement payload is invalid."}) from error


@app.get("/api/procurement/requisitions/{public_id}")
def get_requisition(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("requisitions", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "requisition_not_found", "message": "Requisition was not found."})
    return record


@app.patch("/api/procurement/requisitions/{public_id}/status")
def patch_requisition_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("requisitions", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "procurement status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "requisition_not_found", "message": "Requisition was not found."})
    return record



@app.get("/api/procurement/quotations")
def list_quotations(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("quotations", tenant_slug, status)


@app.post("/api/procurement/quotations")
def post_quotation(payload: dict) -> dict:
    try:
        return create_record("quotations", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "procurement payload is invalid."}) from error


@app.get("/api/procurement/quotations/{public_id}")
def get_quotation(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("quotations", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "quotation_not_found", "message": "Quotation was not found."})
    return record


@app.patch("/api/procurement/quotations/{public_id}/status")
def patch_quotation_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("quotations", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "procurement status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "quotation_not_found", "message": "Quotation was not found."})
    return record



@app.get("/api/procurement/purchase-orders")
def list_purchase_orders(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("purchase-orders", tenant_slug, status)


@app.post("/api/procurement/purchase-orders")
def post_purchase_order(payload: dict) -> dict:
    try:
        return create_record("purchase-orders", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "procurement payload is invalid."}) from error


@app.get("/api/procurement/purchase-orders/{public_id}")
def get_purchase_order(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("purchase-orders", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "purchase_order_not_found", "message": "Purchase Order was not found."})
    return record


@app.patch("/api/procurement/purchase-orders/{public_id}/status")
def patch_purchase_order_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("purchase-orders", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "procurement status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "purchase_order_not_found", "message": "Purchase Order was not found."})
    return record



@app.get("/api/procurement/approvals")
def list_approvals(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("approvals", tenant_slug, status)


@app.post("/api/procurement/approvals")
def post_approval(payload: dict) -> dict:
    try:
        return create_record("approvals", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "procurement payload is invalid."}) from error


@app.get("/api/procurement/approvals/{public_id}")
def get_approval(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("approvals", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "approval_not_found", "message": "Approval was not found."})
    return record


@app.get("/api/procurement/receipts")
def list_receipts(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("receipts", tenant_slug, status)


@app.post("/api/procurement/receipts")
def post_receipt(payload: dict) -> dict:
    try:
        return create_record("receipts", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "procurement payload is invalid."}) from error


@app.get("/api/procurement/receipts/{public_id}")
def get_receipt(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("receipts", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "receipt_not_found", "message": "Receipt was not found."})
    return record


@app.patch("/api/procurement/receipts/{public_id}/status")
def patch_receipt_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("receipts", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "procurement status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "receipt_not_found", "message": "Receipt was not found."})
    return record


@app.get("/api/procurement/three-way-matches")
def list_three_way_matches(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("three-way-matches", tenant_slug, status)


@app.get("/api/procurement/three-way-matches/{public_id}")
def get_three_way_match(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("three-way-matches", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "three_way_match_not_found", "message": "Three Way Match was not found."})
    return record
