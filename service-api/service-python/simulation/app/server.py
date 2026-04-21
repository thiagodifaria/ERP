"""Bootstrap HTTP do servico simulation."""

from fastapi import FastAPI

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready


app = FastAPI(title=settings.service_name)


def build_catalog() -> dict:
    return {
        "service": settings.service_name,
        "scenarioTypes": [
            {
                "key": "operational-load",
                "name": "Operational load",
                "description": "Projects leads, workflows, webhooks and storage pressure from current tenant activity.",
                "inputs": [
                    "tenantSlug",
                    "leadMultiplier",
                    "automationMultiplier",
                    "webhookMultiplier",
                    "teamSize",
                    "planningHorizonDays",
                ],
            },
            {
                "key": "cost-baseline",
                "name": "Cost baseline",
                "description": "Defines the operational levers needed for future cost-estimation reports.",
                "inputs": [
                    "tenantSlug",
                    "storageGrowthGb",
                    "outboundEventsPerMonth",
                    "documentDownloadsPerMonth",
                ],
            },
        ],
    }


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    dependencies = [
        {"name": "scenario-engine", "status": "ready"},
        {"name": "load-model", "status": "ready"},
    ]

    if settings.repository_driver == "postgres":
        dependencies.insert(
            1,
            {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"},
        )
    else:
        dependencies.insert(1, {"name": "scenario-store", "status": "pending-runtime-wiring"})

    return {
        "service": settings.service_name,
        "status": "ready",
        "dependencies": dependencies,
    }


@app.get("/api/simulation/scenarios/catalog")
def scenario_catalog() -> dict:
    return build_catalog()
