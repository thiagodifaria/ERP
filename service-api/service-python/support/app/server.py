"""Bootstrap HTTP do servico support."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import (
    add_case_comment,
    bulk_create_cases,
    build_summary,
    capability_catalog,
    create_case,
    export_cases,
    get_case,
    list_cases,
    list_queues,
    transition_case,
    upsert_queue,
)


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
    dependencies = [
        {"name": "case-management", "status": "ready"},
        {"name": "sla-engine", "status": "ready"},
        {"name": "event-history", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/support/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/support/queues")
def queues(tenant_slug: str | None = None) -> list[dict]:
    return list_queues(tenant_slug)


@app.put("/api/support/queues/{queue_key}")
def put_queue(queue_key: str, payload: dict) -> dict:
    try:
        return upsert_queue(queue_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Support queue payload is invalid."}) from error


@app.get("/api/support/cases")
def cases(tenant_slug: str | None = None, status: str | None = None, priority: str | None = None, cursor: str | None = None, limit: int = 50) -> dict:
    return list_cases(tenant_slug, status, priority, cursor, limit)


@app.get("/api/support/cases/export")
def export(tenant_slug: str | None = None, status: str | None = None, priority: str | None = None) -> dict:
    return export_cases(tenant_slug, status, priority)


@app.post("/api/support/cases")
def post_case(payload: dict) -> dict:
    try:
        return create_case(payload)
    except ValueError as error:
        status_code = 404 if str(error) in {"tenant_not_found", "support_queue_not_found"} else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Support case payload is invalid."}) from error


@app.post("/api/support/cases/bulk")
def post_cases_bulk(payload: dict) -> dict:
    try:
        return bulk_create_cases(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Support bulk payload is invalid."}) from error


@app.get("/api/support/cases/summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)


@app.get("/api/support/cases/{public_id}")
def case_detail(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_case(public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "support_case_not_found", "message": "Support case was not found."})
    return record


@app.patch("/api/support/cases/{public_id}/status")
def patch_case_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_case(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Support case status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "support_case_not_found", "message": "Support case was not found."})
    return record


@app.post("/api/support/cases/{public_id}/comments")
def post_case_comment(public_id: str, payload: dict) -> dict:
    try:
        record = add_case_comment(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Support case comment payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "support_case_not_found", "message": "Support case was not found."})
    return record
