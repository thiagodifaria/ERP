"""Bootstrap HTTP do servico notification."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import (
    bulk_create_notifications,
    build_summary,
    capability_catalog,
    create_notification,
    get_preference,
    list_notifications,
    transition_notification,
    upsert_preference,
)


app = FastAPI(title=settings.service_name)


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    dependencies = [
        {"name": "notification-center", "status": "ready"},
        {"name": "preferences", "status": "ready"},
        {"name": "dispatch", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/notification/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/notification/preferences/{user_public_id}")
def preference(user_public_id: str, tenant_slug: str | None = None) -> dict:
    return get_preference(user_public_id, tenant_slug)


@app.put("/api/notification/preferences/{user_public_id}")
def put_preference(user_public_id: str, payload: dict) -> dict:
    try:
        return upsert_preference(user_public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Notification preference payload is invalid."}) from error


@app.get("/api/notification/center")
def center(tenant_slug: str | None = None, status: str | None = None, cursor: str | None = None, limit: int = 50) -> dict:
    return list_notifications(tenant_slug, status, cursor, limit)


@app.post("/api/notification/center")
def post_notification(payload: dict) -> dict:
    try:
        return create_notification(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Notification payload is invalid."}) from error


@app.post("/api/notification/center/bulk")
def post_notification_bulk(payload: dict) -> dict:
    try:
        return bulk_create_notifications(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Notification bulk payload is invalid."}) from error


@app.patch("/api/notification/center/{public_id}/status")
def patch_notification(public_id: str, payload: dict) -> dict:
    try:
        record = transition_notification(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Notification transition payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "notification_not_found", "message": "Notification was not found."})
    return record


@app.get("/api/notification/summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)
