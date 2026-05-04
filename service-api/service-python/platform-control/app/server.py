"""Bootstrap HTTP do servico platform-control."""

from fastapi import FastAPI, Header, HTTPException
from fastapi.responses import JSONResponse

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import (
    build_usage_summary,
    bulk_upsert_entitlements,
    bulk_upsert_quotas,
    create_lifecycle_job,
    create_metering_snapshot,
    get_lifecycle_job,
    list_blocks,
    list_capability_catalog,
    list_entitlements_page,
    list_lifecycle_jobs_page,
    list_metering_page,
    list_quotas,
    transition_lifecycle_job,
    upsert_block,
    upsert_entitlement,
    upsert_quota,
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
        {"name": "quotas", "status": "ready"},
        {"name": "tenant-blocks", "status": "ready"},
        {"name": "idempotency", "status": "ready"},
        {"name": "async-jobs", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/platform-control/capabilities/catalog")
def get_catalog() -> list[dict]:
    return list_capability_catalog()


@app.get("/api/platform-control/tenants/{tenant_slug}/entitlements")
def entitlements(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    return list_entitlements_page(tenant_slug, cursor, limit)


@app.put("/api/platform-control/tenants/{tenant_slug}/entitlements/{capability_key}")
def put_entitlement(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    try:
        return upsert_entitlement(tenant_slug, capability_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Entitlement payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/entitlements/bulk")
def post_bulk_entitlements(tenant_slug: str, payload: dict) -> dict:
    return bulk_upsert_entitlements(tenant_slug, payload)


@app.get("/api/platform-control/tenants/{tenant_slug}/quotas")
def quotas(tenant_slug: str) -> dict:
    return {"tenantSlug": tenant_slug, "items": list_quotas(tenant_slug)}


@app.put("/api/platform-control/tenants/{tenant_slug}/quotas/{metric_key}")
def put_quota(tenant_slug: str, metric_key: str, payload: dict) -> dict:
    try:
        return upsert_quota(tenant_slug, metric_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Quota payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/quotas/bulk")
def post_bulk_quotas(tenant_slug: str, payload: dict) -> dict:
    return bulk_upsert_quotas(tenant_slug, payload)


@app.get("/api/platform-control/tenants/{tenant_slug}/blocks")
def blocks(tenant_slug: str) -> dict:
    return {"tenantSlug": tenant_slug, "items": list_blocks(tenant_slug)}


@app.put("/api/platform-control/tenants/{tenant_slug}/blocks/{block_key}")
def put_block(tenant_slug: str, block_key: str, payload: dict) -> dict:
    try:
        return upsert_block(tenant_slug, block_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Block payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/metering")
def metering(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    return list_metering_page(tenant_slug, cursor, limit)


@app.post("/api/platform-control/tenants/{tenant_slug}/metering/snapshots")
def post_metering(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_metering_snapshot(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Metering snapshot payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/usage-summary")
def usage_summary(tenant_slug: str) -> dict:
    return build_usage_summary(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs")
def lifecycle_jobs(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    return list_lifecycle_jobs_page(tenant_slug, cursor, limit)


@app.get("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs/{public_id}")
def lifecycle_job_detail(tenant_slug: str, public_id: str) -> dict:
    job = get_lifecycle_job(tenant_slug, public_id)
    if job is None:
        raise HTTPException(status_code=404, detail={"code": "lifecycle_job_not_found", "message": "Lifecycle job was not found."})
    return job


def _queue_lifecycle_job(tenant_slug: str, job_type: str, payload: dict, idempotency_key: str | None) -> JSONResponse:
    try:
        job = create_lifecycle_job(tenant_slug, job_type, payload, idempotency_key)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": f"{job_type.capitalize()} payload is invalid."}) from error

    location = f"/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs/{job['publicId']}"
    return JSONResponse(status_code=202, content=job, headers={"Location": location})


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/onboarding")
def onboarding(tenant_slug: str, payload: dict, idempotency_key: str | None = Header(default=None, alias="Idempotency-Key")) -> JSONResponse:
    return _queue_lifecycle_job(tenant_slug, "onboarding", payload, idempotency_key)


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/offboarding")
def offboarding(tenant_slug: str, payload: dict, idempotency_key: str | None = Header(default=None, alias="Idempotency-Key")) -> JSONResponse:
    return _queue_lifecycle_job(tenant_slug, "offboarding", payload, idempotency_key)


def _transition_lifecycle_job(tenant_slug: str, public_id: str, action: str, payload: dict) -> dict:
    try:
        return transition_lifecycle_job(tenant_slug, public_id, action, payload)
    except ValueError as error:
        if str(error) == "lifecycle_job_not_found":
            raise HTTPException(status_code=404, detail={"code": str(error), "message": "Lifecycle job was not found."}) from error
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Lifecycle transition payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs/{public_id}/start")
def start_lifecycle_job(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_lifecycle_job(tenant_slug, public_id, "start", payload or {})


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs/{public_id}/complete")
def complete_lifecycle_job(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_lifecycle_job(tenant_slug, public_id, "complete", payload or {})


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs/{public_id}/fail")
def fail_lifecycle_job(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_lifecycle_job(tenant_slug, public_id, "fail", payload or {})


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/jobs/{public_id}/cancel")
def cancel_lifecycle_job(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_lifecycle_job(tenant_slug, public_id, "cancel", payload or {})
