"""Bootstrap HTTP do servico simulation."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import (
    build_catalog,
    create_load_benchmark,
    create_operational_load_scenario,
    get_operational_load_scenario,
    list_load_benchmarks,
    list_operational_load_scenarios,
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
        {"name": "scenario-engine", "status": "ready"},
        {"name": "load-model", "status": "ready"},
        {"name": "cost-model", "status": "ready"},
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


@app.post("/api/simulation/scenarios/operational-load")
def create_operational_load(payload: dict) -> dict:
    return create_operational_load_scenario(payload)


@app.get("/api/simulation/scenarios/runs")
def list_scenarios(tenant_slug: str | None = None) -> list[dict]:
    return list_operational_load_scenarios(tenant_slug)


@app.get("/api/simulation/scenarios/runs/{public_id}")
def get_scenario(public_id: str) -> dict:
    scenario = get_operational_load_scenario(public_id)
    if scenario is None:
        raise HTTPException(status_code=404, detail={"code": "simulation_run_not_found", "message": "Simulation run was not found."})
    return scenario


@app.post("/api/simulation/benchmarks/load")
def create_benchmark(payload: dict) -> dict:
    return create_load_benchmark(payload)


@app.get("/api/simulation/benchmarks/runs")
def list_benchmarks(tenant_slug: str | None = None) -> list[dict]:
    return list_load_benchmarks(tenant_slug)
