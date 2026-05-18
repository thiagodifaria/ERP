"""Bootstrap HTTP do servico platform-control."""

from fastapi import FastAPI, Header, HTTPException
from fastapi.responses import JSONResponse

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import (
    apply_go_live_adjustment,
    append_incident_timeline,
    build_autonomous_governance_readiness,
    build_contract_compatibility_matrix,
    build_enterprise_runtime_readiness,
    build_event_mesh_lineage,
    build_go_live_adoption,
    build_go_live_adjustments,
    build_go_live_bottlenecks,
    build_go_live_playbook,
    build_go_live_readiness,
    build_incident_command_readiness,
    build_lifecycle_preview,
    build_lifecycle_readiness,
    build_tenant_runtime_health_score,
    build_usage_summary,
    approve_contract_breaking_change,
    bulk_upsert_entitlements,
    bulk_upsert_quotas,
    create_contract_diff,
    create_contract_snapshot,
    create_tenant_maintenance_window,
    create_go_live_rollout,
    create_approval_request,
    create_incident,
    create_incident_action,
    create_lifecycle_job,
    create_metering_snapshot,
    create_postmortem,
    create_runbook_run,
    evaluate_policy,
    get_evidence_record,
    get_go_live_rollout,
    get_incident,
    get_lifecycle_job,
    get_tenant_runtime_profile,
    list_approval_requests,
    list_contract_breaking_changes,
    list_contract_diffs,
    list_contract_evolution,
    list_blocks,
    list_capability_catalog,
    list_entitlements_page,
    list_evidence_records,
    list_feature_flags_page,
    list_go_live_rollouts,
    list_incidents,
    list_lifecycle_jobs_page,
    list_metering_page,
    list_policy_catalog,
    list_policy_decisions,
    list_provider_catalog,
    list_provider_activation_catalog,
    list_provider_activation_runs,
    list_provider_defaults,
    list_quotas,
    list_runbook_catalog,
    list_runbook_runs,
    list_timeline_events,
    record_event_mesh_event,
    record_timeline_event,
    register_evidence,
    replay_event_mesh_dead_letter,
    run_provider_activation,
    resolve_incident,
    transition_go_live_rollout,
    transition_approval_request,
    transition_lifecycle_job,
    transition_runbook_run,
    update_tenant_runtime_profile,
    upsert_tenant_runtime_quota,
    upsert_block,
    upsert_entitlement,
    upsert_provider_default,
    upsert_quota,
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
        {"name": "capability-registry", "status": "ready"},
        {"name": "entitlements", "status": "ready"},
        {"name": "metering", "status": "ready"},
        {"name": "provider-defaults", "status": "ready"},
        {"name": "tenant-lifecycle", "status": "ready"},
        {"name": "lifecycle-readiness", "status": "ready"},
        {"name": "go-live-control", "status": "ready"},
        {"name": "incident-command", "status": "ready"},
        {"name": "policy-decision-center", "status": "ready"},
        {"name": "operational-timeline", "status": "ready"},
        {"name": "command-approvals", "status": "ready"},
        {"name": "runbook-automation", "status": "ready"},
        {"name": "audit-evidence-vault", "status": "ready"},
        {"name": "enterprise-event-mesh", "status": "ready"},
        {"name": "tenant-runtime-control-plane", "status": "ready"},
        {"name": "contract-schema-evolution", "status": "ready"},
        {"name": "external-provider-activation", "status": "ready"},
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


@app.get("/api/platform-control/providers/catalog")
def get_provider_catalog() -> list[dict]:
    return list_provider_catalog()


@app.get("/api/platform-control/providers/activation/catalog")
def get_provider_activation_catalog() -> dict:
    return list_provider_activation_catalog()


@app.get("/api/platform-control/tenants/{tenant_slug}/providers/activation/runs")
def provider_activation_runs(tenant_slug: str, provider_key: str | None = None) -> dict:
    return list_provider_activation_runs(tenant_slug, provider_key)


@app.post("/api/platform-control/tenants/{tenant_slug}/providers/activation/{provider_key}/test")
def post_provider_activation(tenant_slug: str, provider_key: str, payload: dict | None = None) -> dict:
    try:
        return run_provider_activation(tenant_slug, provider_key, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "provider_activation_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Provider activation payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/entitlements")
def entitlements(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    return list_entitlements_page(tenant_slug, cursor, limit)


@app.get("/api/platform-control/tenants/{tenant_slug}/feature-flags")
def feature_flags(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    return list_feature_flags_page(tenant_slug, cursor, limit)


@app.put("/api/platform-control/tenants/{tenant_slug}/entitlements/{capability_key}")
def put_entitlement(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    try:
        return upsert_entitlement(tenant_slug, capability_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Entitlement payload is invalid."}) from error


@app.put("/api/platform-control/tenants/{tenant_slug}/feature-flags/{capability_key}")
def put_feature_flag(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    try:
        return upsert_entitlement(tenant_slug, capability_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Feature flag payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/entitlements/bulk")
def post_bulk_entitlements(tenant_slug: str, payload: dict) -> dict:
    return bulk_upsert_entitlements(tenant_slug, payload)


@app.get("/api/platform-control/tenants/{tenant_slug}/provider-defaults")
def provider_defaults(tenant_slug: str) -> dict:
    return {"tenantSlug": tenant_slug, "items": list_provider_defaults(tenant_slug)}


@app.put("/api/platform-control/tenants/{tenant_slug}/provider-defaults/{capability_key}")
def put_provider_default(tenant_slug: str, capability_key: str, payload: dict) -> dict:
    try:
        return upsert_provider_default(tenant_slug, capability_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Provider default payload is invalid."}) from error


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


@app.get("/api/platform-control/tenants/{tenant_slug}/lifecycle/readiness")
def lifecycle_readiness(tenant_slug: str) -> dict:
    return build_lifecycle_readiness(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/readiness")
def go_live_readiness(tenant_slug: str) -> dict:
    return build_go_live_readiness(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/adoption")
def go_live_adoption(tenant_slug: str) -> dict:
    return build_go_live_adoption(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/bottlenecks")
def go_live_bottlenecks(tenant_slug: str) -> dict:
    return build_go_live_bottlenecks(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/playbook")
def go_live_playbook(tenant_slug: str) -> dict:
    return build_go_live_playbook(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/adjustments")
def go_live_adjustments(tenant_slug: str) -> dict:
    return build_go_live_adjustments(tenant_slug)


@app.post("/api/platform-control/tenants/{tenant_slug}/go-live/adjustments/apply")
def post_go_live_adjustment(tenant_slug: str, payload: dict) -> dict:
    try:
        return apply_go_live_adjustment(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Go-live adjustment payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/rollouts")
def go_live_rollouts(tenant_slug: str, cursor: str | None = None, limit: int = 50) -> dict:
    return list_go_live_rollouts(tenant_slug, cursor, limit)


@app.get("/api/platform-control/tenants/{tenant_slug}/go-live/rollouts/{public_id}")
def go_live_rollout_detail(tenant_slug: str, public_id: str) -> dict:
    rollout = get_go_live_rollout(tenant_slug, public_id)
    if rollout is None:
        raise HTTPException(status_code=404, detail={"code": "go_live_rollout_not_found", "message": "Go-live rollout was not found."})
    return rollout


@app.post("/api/platform-control/tenants/{tenant_slug}/go-live/rollouts")
def post_go_live_rollout(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_go_live_rollout(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Go-live rollout payload is invalid."}) from error


def _transition_go_live(tenant_slug: str, public_id: str, action: str, payload: dict | None = None) -> dict:
    try:
        return transition_go_live_rollout(tenant_slug, public_id, action, payload or {})
    except ValueError as error:
        if str(error) == "go_live_rollout_not_found":
            raise HTTPException(status_code=404, detail={"code": str(error), "message": "Go-live rollout was not found."}) from error
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Go-live rollout transition payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/go-live/rollouts/{public_id}/start")
def start_go_live_rollout(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_go_live(tenant_slug, public_id, "start", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/go-live/rollouts/{public_id}/complete")
def complete_go_live_rollout(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_go_live(tenant_slug, public_id, "complete", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/go-live/rollouts/{public_id}/rollback")
def rollback_go_live_rollout(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_go_live(tenant_slug, public_id, "rollback", payload)


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


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/onboarding/preview")
def onboarding_preview(tenant_slug: str, payload: dict) -> dict:
    return build_lifecycle_preview(tenant_slug, "onboarding", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/offboarding")
def offboarding(tenant_slug: str, payload: dict, idempotency_key: str | None = Header(default=None, alias="Idempotency-Key")) -> JSONResponse:
    return _queue_lifecycle_job(tenant_slug, "offboarding", payload, idempotency_key)


@app.post("/api/platform-control/tenants/{tenant_slug}/lifecycle/offboarding/preview")
def offboarding_preview(tenant_slug: str, payload: dict) -> dict:
    return build_lifecycle_preview(tenant_slug, "offboarding", payload)


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


@app.get("/api/platform-control/tenants/{tenant_slug}/incidents")
def incidents(tenant_slug: str, status: str | None = None, severity: str | None = None) -> dict:
    return list_incidents(tenant_slug, status, severity)


@app.post("/api/platform-control/tenants/{tenant_slug}/incidents")
def post_incident(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_incident({**payload, "tenantSlug": tenant_slug})
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Incident payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/incident-command/readiness")
def incident_readiness(tenant_slug: str) -> dict:
    return build_incident_command_readiness(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/incidents/{public_id}")
def incident_detail(tenant_slug: str, public_id: str) -> dict:
    incident = get_incident(tenant_slug, public_id)
    if incident is None:
        raise HTTPException(status_code=404, detail={"code": "incident_not_found", "message": "Incident was not found."})
    return incident


@app.post("/api/platform-control/tenants/{tenant_slug}/incidents/{public_id}/timeline")
def post_incident_timeline(tenant_slug: str, public_id: str, payload: dict) -> dict:
    try:
        return append_incident_timeline(tenant_slug, public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "incident_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Incident timeline payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/incidents/{public_id}/actions")
def post_incident_action(tenant_slug: str, public_id: str, payload: dict) -> dict:
    try:
        return create_incident_action(tenant_slug, public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "incident_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Incident action payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/incidents/{public_id}/resolve")
def post_incident_resolve(tenant_slug: str, public_id: str, payload: dict) -> dict:
    try:
        return resolve_incident(tenant_slug, public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "incident_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Incident resolution payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/incidents/{public_id}/postmortem")
def post_incident_postmortem(tenant_slug: str, public_id: str, payload: dict) -> dict:
    try:
        return create_postmortem(tenant_slug, public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "incident_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Postmortem payload is invalid."}) from error


@app.get("/api/platform-control/policies/catalog")
def policy_catalog() -> dict:
    return list_policy_catalog()


@app.post("/api/platform-control/tenants/{tenant_slug}/policies/evaluate")
def post_policy_evaluation(tenant_slug: str, payload: dict) -> dict:
    try:
        return evaluate_policy(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Policy evaluation payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/policies/decisions")
def policy_decisions(tenant_slug: str, action: str | None = None) -> dict:
    return list_policy_decisions(tenant_slug, action)


@app.get("/api/platform-control/tenants/{tenant_slug}/timeline")
def timeline(tenant_slug: str, entity_type: str | None = None, entity_public_id: str | None = None, limit: int = 50) -> dict:
    return list_timeline_events(tenant_slug, entity_type, entity_public_id, limit)


@app.post("/api/platform-control/tenants/{tenant_slug}/timeline")
def post_timeline_event(tenant_slug: str, payload: dict) -> dict:
    try:
        return record_timeline_event(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Timeline event payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/approvals")
def approvals(tenant_slug: str, status: str | None = None) -> dict:
    return list_approval_requests(tenant_slug, status)


@app.post("/api/platform-control/tenants/{tenant_slug}/approvals")
def post_approval(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_approval_request(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Approval payload is invalid."}) from error


def _transition_approval(tenant_slug: str, public_id: str, action: str, payload: dict | None = None) -> dict:
    try:
        return transition_approval_request(tenant_slug, public_id, action, payload or {})
    except ValueError as error:
        status_code = 404 if str(error) == "approval_request_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Approval transition payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/approvals/{public_id}/approve")
def approve_approval(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_approval(tenant_slug, public_id, "approve", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/approvals/{public_id}/reject")
def reject_approval(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_approval(tenant_slug, public_id, "reject", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/approvals/{public_id}/execute")
def execute_approval(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_approval(tenant_slug, public_id, "execute", payload)


@app.get("/api/platform-control/runbooks/catalog")
def runbook_catalog() -> dict:
    return list_runbook_catalog()


@app.get("/api/platform-control/tenants/{tenant_slug}/runbooks")
def runbooks(tenant_slug: str, status: str | None = None) -> dict:
    return list_runbook_runs(tenant_slug, status)


@app.post("/api/platform-control/tenants/{tenant_slug}/runbooks")
def post_runbook(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_runbook_run(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Runbook payload is invalid."}) from error


def _transition_runbook(tenant_slug: str, public_id: str, action: str, payload: dict | None = None) -> dict:
    try:
        return transition_runbook_run(tenant_slug, public_id, action, payload or {})
    except ValueError as error:
        status_code = 404 if str(error) == "runbook_run_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Runbook transition payload is invalid."}) from error


@app.post("/api/platform-control/tenants/{tenant_slug}/runbooks/{public_id}/start")
def start_runbook(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_runbook(tenant_slug, public_id, "start", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/runbooks/{public_id}/complete-step")
def complete_runbook_step(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_runbook(tenant_slug, public_id, "complete-step", payload)


@app.post("/api/platform-control/tenants/{tenant_slug}/runbooks/{public_id}/cancel")
def cancel_runbook(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    return _transition_runbook(tenant_slug, public_id, "cancel", payload)


@app.get("/api/platform-control/tenants/{tenant_slug}/evidence")
def evidence_records(tenant_slug: str, evidence_type: str | None = None, entity_type: str | None = None) -> dict:
    return list_evidence_records(tenant_slug, evidence_type, entity_type)


@app.post("/api/platform-control/tenants/{tenant_slug}/evidence")
def post_evidence(tenant_slug: str, payload: dict) -> dict:
    try:
        return register_evidence(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Evidence payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/evidence/{public_id}")
def evidence_detail(tenant_slug: str, public_id: str) -> dict:
    evidence = get_evidence_record(tenant_slug, public_id)
    if evidence is None:
        raise HTTPException(status_code=404, detail={"code": "evidence_not_found", "message": "Evidence record was not found."})
    return evidence


@app.get("/api/platform-control/tenants/{tenant_slug}/autonomous-governance/readiness")
def autonomous_governance_readiness(tenant_slug: str) -> dict:
    return build_autonomous_governance_readiness(tenant_slug)


@app.get("/api/platform-control/event-mesh/catalog")
def event_mesh_catalog() -> dict:
    return list_event_mesh_catalog()


@app.get("/api/platform-control/tenants/{tenant_slug}/event-mesh/events")
def event_mesh_events(tenant_slug: str, stream_key: str | None = None, status: str | None = None) -> dict:
    return list_event_mesh_events(tenant_slug, stream_key, status)


@app.post("/api/platform-control/tenants/{tenant_slug}/event-mesh/events")
def post_event_mesh_event(tenant_slug: str, payload: dict) -> dict:
    try:
        return record_event_mesh_event(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Event mesh payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/event-mesh/dead-letters")
def event_mesh_dead_letters(tenant_slug: str, status: str | None = None) -> dict:
    return list_event_mesh_dead_letters(tenant_slug, status)


@app.post("/api/platform-control/tenants/{tenant_slug}/event-mesh/dead-letters/{public_id}/replay")
def replay_dead_letter(tenant_slug: str, public_id: str, payload: dict | None = None) -> dict:
    try:
        return replay_event_mesh_dead_letter(tenant_slug, public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "dead_letter_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Dead letter replay payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/event-mesh/lineage")
def event_mesh_lineage(tenant_slug: str, correlation_id: str | None = None) -> dict:
    return build_event_mesh_lineage(tenant_slug, correlation_id)


@app.get("/api/platform-control/tenants/{tenant_slug}/runtime/profile")
def tenant_runtime_profile(tenant_slug: str) -> dict:
    return get_tenant_runtime_profile(tenant_slug)


@app.put("/api/platform-control/tenants/{tenant_slug}/runtime/profile")
def put_tenant_runtime_profile(tenant_slug: str, payload: dict) -> dict:
    return update_tenant_runtime_profile(tenant_slug, payload)


@app.get("/api/platform-control/tenants/{tenant_slug}/runtime/health-score")
def tenant_runtime_health_score(tenant_slug: str) -> dict:
    return build_tenant_runtime_health_score(tenant_slug)


@app.get("/api/platform-control/tenants/{tenant_slug}/runtime/quotas")
def tenant_runtime_quotas(tenant_slug: str) -> dict:
    return list_tenant_runtime_quotas(tenant_slug)


@app.put("/api/platform-control/tenants/{tenant_slug}/runtime/quotas/{metric_key}")
def put_tenant_runtime_quota(tenant_slug: str, metric_key: str, payload: dict) -> dict:
    try:
        return upsert_tenant_runtime_quota(tenant_slug, metric_key, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Runtime quota payload is invalid."}) from error


@app.get("/api/platform-control/tenants/{tenant_slug}/runtime/maintenance-windows")
def tenant_maintenance_windows(tenant_slug: str) -> dict:
    return list_tenant_maintenance_windows(tenant_slug)


@app.post("/api/platform-control/tenants/{tenant_slug}/runtime/maintenance-windows")
def post_tenant_maintenance_window(tenant_slug: str, payload: dict) -> dict:
    try:
        return create_tenant_maintenance_window(tenant_slug, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Maintenance window payload is invalid."}) from error


@app.get("/api/platform-control/contracts/evolution")
def contract_evolution() -> dict:
    return list_contract_evolution()


@app.post("/api/platform-control/contracts/evolution/snapshots")
def post_contract_snapshot(payload: dict) -> dict:
    try:
        return create_contract_snapshot(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Contract snapshot payload is invalid."}) from error


@app.get("/api/platform-control/contracts/evolution/diffs")
def contract_diffs(contract_key: str | None = None) -> dict:
    return list_contract_diffs(contract_key)


@app.post("/api/platform-control/contracts/evolution/diffs")
def post_contract_diff(payload: dict) -> dict:
    try:
        return create_contract_diff(payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "Contract diff payload is invalid."}) from error


@app.get("/api/platform-control/contracts/evolution/breaking-changes")
def contract_breaking_changes(status: str | None = None) -> dict:
    return list_contract_breaking_changes(status)


@app.post("/api/platform-control/contracts/evolution/breaking-changes/{public_id}/approve")
def approve_contract_breaking(public_id: str, payload: dict | None = None) -> dict:
    try:
        return approve_contract_breaking_change(public_id, payload)
    except ValueError as error:
        status_code = 404 if str(error) == "breaking_change_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "Breaking change approval payload is invalid."}) from error


@app.get("/api/platform-control/contracts/evolution/compatibility-matrix")
def contract_compatibility_matrix() -> dict:
    return build_contract_compatibility_matrix()


@app.get("/api/platform-control/tenants/{tenant_slug}/enterprise-runtime/readiness")
def enterprise_runtime_readiness(tenant_slug: str) -> dict:
    return build_enterprise_runtime_readiness(tenant_slug)
    list_event_mesh_catalog,
    list_event_mesh_dead_letters,
    list_event_mesh_events,
    list_tenant_maintenance_windows,
    list_tenant_runtime_quotas,
