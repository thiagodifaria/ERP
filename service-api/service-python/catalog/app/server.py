"""Bootstrap HTTP do servico catalog."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import capability_catalog, create_category, create_item, get_item, list_categories, list_items, update_item


app = FastAPI(title=settings.service_name)


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    dependencies = [{"name": "catalog-store", "status": "ready"}]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/catalog/capabilities")
def get_capabilities() -> dict:
    return capability_catalog()


@app.get("/api/catalog/categories")
def categories(tenant_slug: str | None = None) -> list[dict]:
    return list_categories(tenant_slug)


@app.post("/api/catalog/categories")
def post_category(payload: dict) -> dict:
    try:
        return create_category(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Catalog category payload is invalid."}) from error


@app.get("/api/catalog/items")
def items(tenant_slug: str | None = None, item_type: str | None = None, active: bool | None = None) -> list[dict]:
    return list_items(tenant_slug, item_type, active)


@app.post("/api/catalog/items")
def post_item(payload: dict) -> dict:
    try:
        return create_item(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "catalog_category_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Catalog item payload is invalid."}) from error


@app.get("/api/catalog/items/{public_id}")
def item_detail(public_id: str, tenant_slug: str | None = None) -> dict:
    item = get_item(public_id, tenant_slug)
    if item is None:
        raise HTTPException(status_code=404, detail={"code": "catalog_item_not_found", "message": "Catalog item was not found."})
    return item


@app.patch("/api/catalog/items/{public_id}")
def patch_item(public_id: str, payload: dict) -> dict:
    item = update_item(public_id, payload)
    if item is None:
        raise HTTPException(status_code=404, detail={"code": "catalog_item_not_found", "message": "Catalog item was not found."})
    return item
