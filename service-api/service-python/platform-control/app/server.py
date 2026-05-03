"""Bootstrap HTTP do servico platform-control."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import (
    create_lifecycle_job,
    create_metering_snapshot,
    list_capability_catalog,
    list_entitlements,
    list_lifecycle_jobs,
    list_metering,
    upsert_entitlement,
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
        {"name": "capability-registry", "status": "ready"},
        {"name": "entitlements", "status": "ready"},
        {"name": "metering", "status": "ready"},
        {"name": "tenant-lifecycle", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/platform-control/capabilities/catalog")
def get_catalog() -> list[dict]:
    return list_capability_catalog()


@app.get("/api/platform-control/tenants/{tenant_slug}/entitlements")
def entitlements(tenant_slug: str) -> list[dict]:
    return list_entitlements(tenant_slug)


@app.put("/api/platform-control/tenants/{tenant_slug}/entitlements/{capability_key}")
def put_entitlement(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    try:
        return upsert_entitlement(tenant_slug, capability_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Entitlement payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/metering")
def metering(tenant_slug: str) -> dict:
    return list_metering(tenant_slug)


@app.post("/api/platform-control/tenants/{tenant_slug}/metering/snapshots")
def post_metering(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_metering_snapshot(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Metering snapshot payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs")
def lifecycle_jobs(tenant_slug: str) -> list[dict]:
    return list_lifecycle_jobs(tenant_slug)


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/onboarding")
def onboarding(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_lifecycle_job(tenant_slug, "onboarding", payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Onboarding payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/offboarding")
def offboarding(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_lifecycle_job(tenant_slug, "offboarding", payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Offboarding payload is invalid."}) from error
