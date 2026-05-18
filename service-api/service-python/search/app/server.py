"""Bootstrap HTTP do servico search."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import (
    add_case_item,
    build_facets,
    capability_catalog,
    create_discovery_case,
    create_export_request,
    create_legal_hold,
    create_saved_query,
    list_audit_events,
    query_index,
)
from app.security import install_security_middleware

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
        {"name": "search-index", "status": "ready"},
        {"name": "ediscovery", "status": "ready"},
        {"name": "query-audit", "status": "ready"},
        {"name": "redaction", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/search/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/search/query")
def search_query(tenant_slug: str | None = None, q: str | None = None, entity_type: str | None = None, actor: str | None = None, include_sensitive: bool = False) -> dict:
    return query_index(tenant_slug, q, entity_type, actor, include_sensitive)


@app.get("/api/search/facets")
def facets(tenant_slug: str | None = None) -> dict:
    return build_facets(tenant_slug)


@app.post("/api/search/saved-queries")
def saved_query(payload: dict) -> dict:
    try:
        return create_saved_query(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Saved query payload is invalid."}) from error


@app.get("/api/search/audit-events")
def audit_events(tenant_slug: str | None = None) -> dict:
    return list_audit_events(tenant_slug)


@app.post("/api/search/discovery-cases")
def discovery_case(payload: dict) -> dict:
    try:
        return create_discovery_case(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Discovery case payload is invalid."}) from error


@app.post("/api/search/discovery-cases/{public_id}/items")
def discovery_case_item(public_id: str, payload: dict) -> dict:
    try:
        return add_case_item(public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "discovery_case_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Discovery case item payload is invalid."}) from error


@app.post("/api/search/legal-holds")
def legal_hold(payload: dict) -> dict:
    try:
        return create_legal_hold(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Legal hold payload is invalid."}) from error


@app.post("/api/search/exports")
def export_request(payload: dict) -> dict:
    try:
        return create_export_request(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Export request payload is invalid."}) from error

