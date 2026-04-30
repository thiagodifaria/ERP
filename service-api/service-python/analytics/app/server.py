"""Bootstrap HTTP do servico analytics."""

from fastapi import FastAPI

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.reports.automation_board import build_automation_board
from app.reports.cost_estimator import build_cost_estimator
from app.reports.delivery_reliability import build_delivery_reliability
from app.reports.engagement_operations import build_engagement_operations
from app.reports.load_benchmark import build_load_benchmark
from app.reports.pipeline_summary import build_pipeline_summary
from app.reports.rental_operations import build_rental_operations
from app.reports.revenue_operations import build_revenue_operations
from app.reports.sales_journey import build_sales_journey
from app.reports.service_pulse import build_service_pulse
from app.reports.tenant_360 import build_tenant_360
from app.reports.workflow_definition_health import build_workflow_definition_health


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
        {"name": "forecast-model", "status": "ready"},
        {"name": "simulation-catalog", "status": "ready"},
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


@app.get("/api/analytics/reports/sales-journey")
def sales_journey(tenant_slug: str | None = None) -> dict:
    return build_sales_journey(tenant_slug)


@app.get("/api/analytics/reports/tenant-360")
def tenant_360(tenant_slug: str | None = None) -> dict:
    return build_tenant_360(tenant_slug)


@app.get("/api/analytics/reports/automation-board")
def automation_board(tenant_slug: str | None = None) -> dict:
    return build_automation_board(tenant_slug)


@app.get("/api/analytics/reports/workflow-definition-health")
def workflow_definition_health(tenant_slug: str | None = None) -> dict:
    return build_workflow_definition_health(tenant_slug)


@app.get("/api/analytics/reports/delivery-reliability")
def delivery_reliability(provider: str | None = None) -> dict:
    return build_delivery_reliability(provider)


@app.get("/api/analytics/reports/engagement-operations")
def engagement_operations(tenant_slug: str | None = None) -> dict:
    return build_engagement_operations(tenant_slug)


@app.get("/api/analytics/reports/revenue-operations")
def revenue_operations(tenant_slug: str | None = None) -> dict:
    return build_revenue_operations(tenant_slug)


@app.get("/api/analytics/reports/rental-operations")
def rental_operations(tenant_slug: str | None = None) -> dict:
    return build_rental_operations(tenant_slug)


@app.get("/api/analytics/reports/cost-estimator")
def cost_estimator(tenant_slug: str | None = None) -> dict:
    return build_cost_estimator(tenant_slug)


@app.get("/api/analytics/reports/load-benchmark")
def load_benchmark(tenant_slug: str | None = None) -> dict:
    return build_load_benchmark(tenant_slug)
