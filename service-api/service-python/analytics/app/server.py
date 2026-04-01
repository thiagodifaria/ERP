"""Bootstrap HTTP do servico analytics."""

from fastapi import FastAPI

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.reports.automation_board import build_automation_board
from app.reports.pipeline_summary import build_pipeline_summary
from app.reports.service_pulse import build_service_pulse
from app.reports.tenant_360 import build_tenant_360


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
        {"name": "report-engine", "status": "ready"},
        {"name": "forecast-model", "status": "pending-runtime-wiring"},
    ]

    if settings.repository_driver == "postgres":
        dependencies.insert(
            1,
            {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"},
        )
    else:
        dependencies.insert(1, {"name": "warehouse", "status": "pending-runtime-wiring"})

    return {
        "service": settings.service_name,
        "status": "ready",
        "dependencies": dependencies,
    }


@app.get("/api/analytics/reports/pipeline-summary")
def pipeline_summary(tenant_slug: str | None = None) -> dict:
    return build_pipeline_summary(tenant_slug)


@app.get("/api/analytics/reports/service-pulse")
def service_pulse(tenant_slug: str | None = None) -> dict:
    return build_service_pulse(tenant_slug)


@app.get("/api/analytics/reports/tenant-360")
def tenant_360(tenant_slug: str | None = None) -> dict:
    return build_tenant_360(tenant_slug)


@app.get("/api/analytics/reports/automation-board")
def automation_board(tenant_slug: str | None = None) -> dict:
    return build_automation_board(tenant_slug)
