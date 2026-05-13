"""Bootstrap HTTP do servico fiscal."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import (
    build_compliance_summary,
    build_privacy_export_package,
    build_retention_execution,
    cancel_document,
    create_consent,
    create_correction_letter,
    create_document,
    create_invalidation,
    create_privacy_request,
    execute_privacy_request,
    execute_retention_execution,
    get_company_profile,
    get_document,
    get_privacy_request,
    list_audit_events,
    list_capabilities,
    list_consents,
    list_document_events,
    list_documents,
    list_privacy_requests,
    list_retention_policies,
    transition_consent,
    transition_privacy_request,
    upsert_company_profile,
    upsert_retention_policy,
)


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
        {"name": "fiscal-rules", "status": "ready"},
        {"name": "lgpd-operations", "status": "ready"},
        {"name": "retention-governance", "status": "ready"},
    ]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/fiscal/capabilities")
def capabilities() -> dict:
    return list_capabilities()


@app.get("/api/fiscal/companies/{company_public_id}/profile")
def company_profile(company_public_id: str, tenant_slug: str | None = None) -> dict:
    payload = get_company_profile(company_public_id, tenant_slug)
    if payload is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_company_profile_not_found", "message": "Fiscal company profile was not found."})
    return payload


@app.put("/api/fiscal/companies/{company_public_id}/profile")
def put_company_profile(company_public_id: str, payload: dict) -> dict:
    try:
        return upsert_company_profile(company_public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Fiscal company profile payload is invalid."}) from error


@app.get("/api/fiscal/companies/{company_public_id}/retention-policies")
def retention_policies(company_public_id: str, tenant_slug: str | None = None) -> list[dict]:
    return list_retention_policies(company_public_id, tenant_slug)


@app.get("/api/fiscal/companies/{company_public_id}/retention-execution")
def retention_execution(company_public_id: str, tenant_slug: str | None = None) -> dict:
    return build_retention_execution(company_public_id, tenant_slug)


@app.post("/api/fiscal/companies/{company_public_id}/retention-execution/execute")
def post_retention_execution(company_public_id: str, payload: dict) -> dict:
    try:
        return execute_retention_execution(company_public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Retention execution payload is invalid."}) from error


@app.put("/api/fiscal/companies/{company_public_id}/retention-policies/{data_domain}")
def put_retention_policy(company_public_id: str, data_domain: str, payload: dict) -> dict:
    try:
        return upsert_retention_policy(company_public_id, data_domain, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Fiscal retention policy payload is invalid."}) from error


@app.get("/api/fiscal/documents")
def documents(tenant_slug: str | None = None) -> list[dict]:
    return list_documents(tenant_slug)


@app.get("/api/fiscal/documents/{public_id}")
def fiscal_document_detail(public_id: str, tenant_slug: str | None = None) -> dict:
    payload = get_document(public_id, tenant_slug)
    if payload is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_document_not_found", "message": "Fiscal document was not found."})
    return payload


@app.post("/api/fiscal/documents")
def post_document(payload: dict) -> dict:
    try:
        return create_document(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Fiscal document payload is invalid."}) from error


@app.post("/api/fiscal/documents/{public_id}/cancel")
def post_cancel_document(public_id: str, payload: dict) -> dict:
    record = cancel_document(public_id, payload)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_document_not_found", "message": "Fiscal document was not found."})
    return record


@app.post("/api/fiscal/documents/{public_id}/correction-letter")
def post_correction_letter(public_id: str, payload: dict) -> dict:
    try:
        record = create_correction_letter(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Fiscal correction payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_document_not_found", "message": "Fiscal document was not found."})
    return record


@app.post("/api/fiscal/documents/{public_id}/invalidate")
def post_invalidation(public_id: str, payload: dict) -> dict:
    record = create_invalidation(public_id, payload)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_document_not_found", "message": "Fiscal document was not found."})
    return record


@app.get("/api/fiscal/documents/{public_id}/events")
def document_events(public_id: str, tenant_slug: str | None = None) -> list[dict]:
    return list_document_events(tenant_slug, public_id)


@app.post("/api/fiscal/privacy-requests")
def post_privacy_request(payload: dict) -> dict:
    try:
        return create_privacy_request(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Privacy request payload is invalid."}) from error


@app.get("/api/fiscal/privacy-requests")
def privacy_requests(tenant_slug: str | None = None) -> list[dict]:
    return list_privacy_requests(tenant_slug)


@app.get("/api/fiscal/privacy-requests/{public_id}")
def privacy_request_detail(public_id: str, tenant_slug: str | None = None) -> dict:
    payload = get_privacy_request(public_id, tenant_slug)
    if payload is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_privacy_request_not_found", "message": "Privacy request was not found."})
    return payload


@app.get("/api/fiscal/privacy-requests/{public_id}/export-package")
def privacy_export_package(public_id: str, tenant_slug: str | None = None) -> dict:
    payload = build_privacy_export_package(public_id, tenant_slug)
    if payload is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_privacy_request_not_found", "message": "Privacy request was not found."})
    return payload


@app.post("/api/fiscal/privacy-requests/{public_id}/execute")
def post_execute_privacy_request(public_id: str, payload: dict) -> dict:
    try:
        result = execute_privacy_request(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Privacy request execution payload is invalid."}) from error
    if result is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_privacy_request_not_found", "message": "Privacy request was not found."})
    return result


@app.patch("/api/fiscal/privacy-requests/{public_id}/status")
def patch_privacy_request(public_id: str, payload: dict) -> dict:
    try:
        record = transition_privacy_request(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Privacy request transition payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_privacy_request_not_found", "message": "Privacy request was not found."})
    return record


@app.post("/api/fiscal/consents")
def post_consent(payload: dict) -> dict:
    try:
        return create_consent(payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Consent payload is invalid."}) from error


@app.get("/api/fiscal/consents")
def consents(tenant_slug: str | None = None) -> list[dict]:
    return list_consents(tenant_slug)


@app.patch("/api/fiscal/consents/{public_id}")
def patch_consent(public_id: str, payload: dict) -> dict:
    try:
        record = transition_consent(public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Consent transition payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "fiscal_consent_not_found", "message": "Consent was not found."})
    return record


@app.get("/api/fiscal/audit-events")
def audit_events(tenant_slug: str | None = None) -> list[dict]:
    return list_audit_events(tenant_slug)


@app.get("/api/fiscal/compliance/summary")
def compliance_summary(tenant_slug: str | None = None) -> dict:
    return build_compliance_summary(tenant_slug)
