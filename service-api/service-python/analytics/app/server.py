"""Bootstrap HTTP do servico analytics."""

from fastapi import FastAPI

from app.config.settings import settings
from app.reports.pipeline_summary import build_pipeline_summary


app = FastAPI(title=settings.service_name)


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    return {
        "service": settings.service_name,
        "status": "ready",
        "dependencies": [
            {"name": "report-engine", "status": "ready"},
            {"name": "warehouse", "status": "pending-runtime-wiring"},
            {"name": "forecast-model", "status": "pending-runtime-wiring"},
        ],
    }


@app.get("/api/analytics/reports/pipeline-summary")
def pipeline_summary(tenant_slug: str | None = None) -> dict:
    return build_pipeline_summary(tenant_slug)
