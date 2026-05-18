from __future__ import annotations

from datetime import datetime, timezone
import json
import re
from urllib import error as urlerror
from urllib import request
import uuid

from app.config.settings import settings

SENSITIVE_PATTERN = re.compile(r"([\\w.%-]+@[\\w.-]+\\.[A-Za-z]{2,}|\\b\\d{3}\\.?\\d{3}\\.?\\d{3}-?\\d{2}\\b|\\b\\d{2}\\.?\\d{3}\\.?\\d{3}/?\\d{4}-?\\d{2}\\b|Bearer\\s+[A-Za-z0-9._-]+|access[_-]?token[:=][A-Za-z0-9._-]+)", re.IGNORECASE)

TOOL_REGISTRY = [
    {
        "toolKey": "search.query",
        "service": "search",
        "mode": "read",
        "capability": "operational-search",
        "description": "Consulta o indice operacional tenant-aware.",
    },
    {
        "toolKey": "analytics.metric.lookup",
        "service": "analytics",
        "mode": "read",
        "capability": "semantic-metrics",
        "description": "Consulta definicoes e snapshots de metricas.",
    },
    {
        "toolKey": "platform.incident.summary",
        "service": "platform-control",
        "mode": "read",
        "capability": "incident-command",
        "description": "Resume incidentes ativos e postmortems.",
    },
]

POLICIES = [
    {
        "policyKey": "read-only-default",
        "effect": "allow",
        "mode": "read",
        "requiresTenant": True,
        "requiresAudit": True,
        "mutationAllowed": False,
    },
    {
        "policyKey": "deny-mutation",
        "effect": "deny",
        "mode": "write",
        "requiresTenant": True,
        "requiresAudit": True,
        "mutationAllowed": False,
    },
]

IN_MEMORY_STATE: dict[str, list[dict]] = {
    "assistant_runs": [],
    "assistant_run_actions": [],
    "prompt_audit_events": [],
    "redaction_events": [],
}


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def list_tools() -> dict:
    return {
        "service": settings.service_name,
        "provider": {
            "providerKey": "openai",
            "configured": bool(settings.openai_api_key.strip()),
            "credentialKey": "OPENAI_API_KEY",
            "model": settings.openai_model,
            "fallback": "deterministic_local_answer",
        },
        "items": TOOL_REGISTRY,
    }


def list_policies() -> dict:
    return {"service": settings.service_name, "items": POLICIES}


def redact_text(value: str) -> tuple[str, list[str]]:
    findings: list[str] = []

    def replace(match: re.Match) -> str:
        token = match.group(0)
        kind = "token" if "Bearer" in token or "token" in token.lower() else "identifier"
        findings.append(kind)
        return f"[REDACTED:{kind}]"

    return SENSITIVE_PATTERN.sub(replace, value), sorted(set(findings))


def preview_redaction(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    text = str(payload.get("text") or "")
    redacted, findings = redact_text(text)
    event = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "findings": findings,
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["redaction_events"].append(event)
    return {"tenantSlug": slug, "redactedText": redacted, "findings": findings, "auditEvent": event}


def create_assistant_run(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    actor = str(payload.get("actor") or "").strip()
    prompt = str(payload.get("prompt") or "").strip()
    requested_tools = payload.get("tools") if isinstance(payload.get("tools"), list) else []
    if not actor:
        raise ValueError("actor_required")
    if not prompt:
        raise ValueError("prompt_required")
    redacted_prompt, findings = redact_text(prompt)
    tool_map = {tool["toolKey"]: tool for tool in TOOL_REGISTRY}
    actions: list[dict] = []
    denied: list[str] = []
    for tool_key in requested_tools:
        tool = tool_map.get(str(tool_key))
        if tool is None:
            denied.append(str(tool_key))
            continue
        if tool["mode"] != "read":
            denied.append(tool["toolKey"])
            continue
        actions.append(
            {
                "publicId": str(uuid.uuid4()),
                "toolKey": tool["toolKey"],
                "mode": tool["mode"],
                "status": "allowed",
                "summary": f"Read-only tool {tool['toolKey']} approved for tenant {slug}.",
            }
        )
    answer = _deterministic_answer(redacted_prompt, actions, denied)
    provider = "deterministic"
    if settings.openai_api_key.strip() and not denied:
        provider_answer = _openai_response(redacted_prompt, actions)
        if provider_answer["status"] == "succeeded":
            answer = provider_answer["answer"]
            provider = "openai"
        else:
            answer["providerFallbackReason"] = provider_answer["reason"]
    run = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "actor": actor,
        "mode": "read-only",
        "status": "completed" if not denied else "completed_with_denials",
        "provider": provider,
        "model": settings.openai_model if provider == "openai" else "deterministic-local",
        "prompt": redacted_prompt,
        "redactionFindings": findings,
        "answer": answer,
        "actions": actions,
        "deniedTools": denied,
        "createdAt": utc_now(),
    }
    audit = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "actor": actor,
        "runPublicId": run["publicId"],
        "redactionFindings": findings,
        "toolCount": len(actions),
        "deniedToolCount": len(denied),
        "createdAt": utc_now(),
    }
    IN_MEMORY_STATE["assistant_runs"].append(run)
    IN_MEMORY_STATE["assistant_run_actions"].extend(actions)
    IN_MEMORY_STATE["prompt_audit_events"].append(audit)
    run["auditEvent"] = audit
    return run


def get_assistant_run(public_id: str, tenant: str | None = None) -> dict | None:
    slug = tenant_slug(tenant)
    return next((item for item in IN_MEMORY_STATE["assistant_runs"] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)


def list_audit_events(tenant: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    return {"tenantSlug": slug, "items": [item for item in IN_MEMORY_STATE["prompt_audit_events"] if item["tenantSlug"] == slug]}


def _deterministic_answer(prompt: str, actions: list[dict], denied: list[str]) -> dict:
    return {
        "summary": "Resposta governada gerada em modo deterministico e somente leitura.",
        "promptPreview": prompt[:180],
        "toolActions": len(actions),
        "deniedTools": denied,
        "nextSteps": [
            "Use search.query para recuperar evidencias tenant-aware.",
            "Use analytics.metric.lookup para conferir metricas versionadas.",
            "Abra incidente apenas via platform-control, fora da camada AI read-only.",
        ],
    }


def _openai_response(prompt: str, actions: list[dict]) -> dict:
    payload = {
        "model": settings.openai_model,
        "input": (
            "You are a read-only ERP operational analyst. "
            "Do not invent data. Summarize the following redacted prompt and approved read-only tools. "
            f"Prompt: {prompt}\nTools: {[item['toolKey'] for item in actions]}"
        ),
    }
    req = request.Request(
        "https://api.openai.com/v1/responses",
        data=json.dumps(payload).encode("utf-8"),
        headers={
            "Authorization": f"Bearer {settings.openai_api_key.strip()}",
            "Content-Type": "application/json",
            "Accept": "application/json",
        },
        method="POST",
    )
    try:
        with request.urlopen(req, timeout=20) as response:
            body = json.loads(response.read().decode("utf-8"))
            text = body.get("output_text")
            if not text:
                text = "LLM response returned without output_text; inspect provider response in external logs."
            return {
                "status": "succeeded",
                "answer": {
                    "summary": text,
                    "toolActions": len(actions),
                    "deniedTools": [],
                    "provider": "openai",
                    "model": settings.openai_model,
                },
            }
    except urlerror.HTTPError as exc:
        return {"status": "failed", "reason": f"openai_http_{exc.code}"}
    except (urlerror.URLError, TimeoutError) as exc:
        return {"status": "failed", "reason": str(exc)}
