"""Bootstrap HTTP do servico ai-governance."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.runtime import create_assistant_run, get_assistant_run, list_audit_events, list_policies, list_tools, preview_redaction
from app.security import install_security_middleware

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
        {"name": "tool-registry", "status": "ready"},
        {"name": "policy-engine", "status": "ready"},
        {"name": "prompt-audit", "status": "ready"},
        {"name": "redaction", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/ai-governance/tools")
def tools() -> dict:
    return list_tools()


@app.get("/api/ai-governance/policies")
def policies() -> dict:
    return list_policies()


@app.post("/api/ai-governance/runs")
def post_run(payload: dict) -> dict:
    try:
        return create_assistant_run(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Assistant run payload is invalid."}) from error


@app.get("/api/ai-governance/runs/{public_id}")
def get_run(public_id: str, tenant_slug: str | None = None) -> dict:
    run = get_assistant_run(public_id, tenant_slug)
    if run is None:
        raise HTTPException(status_code=404, detail={"code": "assistant_run_not_found", "message": "Assistant run was not found."})
    return run


@app.get("/api/ai-governance/audit-events")
def audit_events(tenant_slug: str | None = None) -> dict:
    return list_audit_events(tenant_slug)


@app.post("/api/ai-governance/redaction/preview")
def redaction_preview(payload: dict) -> dict:
    return preview_redaction(payload)

