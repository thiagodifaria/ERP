"""Bootstrap HTTP do servico supplier."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import (
    build_summary,
    bulk_create_suppliers,
    capability_catalog,
    create_supplier,
    export_suppliers,
    get_supplier,
    list_categories,
    list_suppliers,
    update_supplier,
    upsert_category,
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
        {"name": "supplier-directory", "status": "ready"},
        {"name": "payables-profile", "status": "ready"},
        {"name": "cnpj-enrichment-ready", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/supplier/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/supplier/categories")
def categories(tenant_slug: str | None = None) -> list[dict]:
    return list_categories(tenant_slug)


@app.put("/api/supplier/categories/{category_key}")
def put_category(category_key: str, payload: dict) -> dict:
    try:
        return upsert_category(category_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Supplier category payload is invalid."}) from error


@app.get("/api/supplier/suppliers")
def suppliers(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_suppliers(tenant_slug, status)


@app.get("/api/supplier/suppliers/export")
def export(tenant_slug: str | None = None, status: str | None = None) -> dict:
    return export_suppliers(tenant_slug, status)


@app.post("/api/supplier/suppliers")
def post_supplier(payload: dict) -> dict:
    try:
        return create_supplier(payload)
    except ValueError as error:
        status_code = 404 if str(error) in {"tenant_not_found", "supplier_category_not_found"} else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Supplier payload is invalid."}) from error


@app.post("/api/supplier/suppliers/bulk")
def post_suppliers_bulk(payload: dict) -> dict:
    try:
        return bulk_create_suppliers(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Supplier bulk payload is invalid."}) from error


@app.get("/api/supplier/suppliers/summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)


@app.get("/api/supplier/suppliers/{public_id}")
def supplier_detail(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_supplier(public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "supplier_not_found", "message": "Supplier was not found."})
    return record


@app.patch("/api/supplier/suppliers/{public_id}")
def patch_supplier(public_id: str, payload: dict) -> dict:
    record = update_supplier(public_id, payload)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "supplier_not_found", "message": "Supplier was not found."})
    return record
