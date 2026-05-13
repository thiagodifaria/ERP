"""Runtime API security middleware shared by Python services."""

from __future__ import annotations

import base64
import hashlib
import hmac
import json
import os
import time
import urllib.error
import urllib.request
from typing import Any

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse


def install_security_middleware(app: FastAPI, service_name: str) -> None:
    @app.middleware("http")
    async def api_security(request: Request, call_next):  # type: ignore[no-untyped-def]
        if not _security_enforced() or request.url.path.startswith("/health/"):
            return await call_next(request)

        if _requires_correlation(request.method) and not request.headers.get("x-correlation-id"):
            return _error(400, "correlation_id_required", "Mutation requests require X-Correlation-Id.")

        auth = _authenticate(request)
        if auth is None:
            return _error(401, "unauthorized", "Bearer token is invalid or missing.")

        subject, tenant_slug, scopes = auth
        request.scope["headers"].append((b"x-erp-auth-subject", subject.encode()))
        request.scope["headers"].append((b"x-erp-auth-tenant", tenant_slug.encode()))
        request.scope["headers"].append((b"x-erp-auth-scopes", " ".join(scopes).encode()))

        if not _authorize_openfga(service_name, request, subject, tenant_slug):
            return _error(403, "openfga_denied", "OpenFGA denied the request.")

        return await call_next(request)


def _security_enforced() -> bool:
    mode = os.getenv("ERP_AUTH_ENFORCEMENT", "").strip().lower()
    if mode in {"disabled", "off", "false"}:
        return False
    if mode in {"enforced", "strict", "true"}:
        return True
    environment = os.getenv("ERP_ENV", "local").strip().lower()
    return environment not in {"", "local", "dev", "development", "test", "testing"}


def _authenticate(request: Request) -> tuple[str, str, list[str]] | None:
    header = request.headers.get("authorization", "")
    if not header.lower().startswith("bearer "):
        return None
    token = header[7:].strip()
    internal_token = os.getenv("ERP_INTERNAL_SERVICE_TOKEN", "").strip()
    if internal_token and hmac.compare_digest(token, internal_token):
        return "service:internal", _resolve_tenant(request), ["service"]

    claims = _verify_jwt(token)
    if claims is None:
        return None
    subject = str(claims.get("sub") or claims.get("user_public_id") or "")
    tenant_slug = str(claims.get("tenant_slug") or claims.get("tenant") or _resolve_tenant(request))
    scopes = claims.get("scope", [])
    if isinstance(scopes, str):
        parsed_scopes = scopes.split()
    elif isinstance(scopes, list):
        parsed_scopes = [str(scope) for scope in scopes]
    else:
        parsed_scopes = []
    return (subject, tenant_slug, parsed_scopes) if subject else None


def _verify_jwt(token: str) -> dict[str, Any] | None:
    secret = os.getenv("ERP_JWT_HS256_SECRET", "")
    parts = token.split(".")
    if not secret or len(parts) != 3:
        return None
    try:
        header = json.loads(_base64url_decode(parts[0]))
        if header.get("alg") != "HS256":
            return None
        expected = hmac.new(secret.encode(), f"{parts[0]}.{parts[1]}".encode(), hashlib.sha256).digest()
        if not hmac.compare_digest(parts[2], _base64url_encode(expected)):
            return None
        claims = json.loads(_base64url_decode(parts[1]))
    except (ValueError, json.JSONDecodeError):
        return None
    expires_at = claims.get("exp")
    if isinstance(expires_at, (int, float)) and expires_at <= time.time():
        return None
    return claims


def _authorize_openfga(service_name: str, request: Request, subject: str, tenant_slug: str) -> bool:
    if os.getenv("ERP_OPENFGA_ENFORCEMENT", "").lower() != "true":
        return True
    base_url = os.getenv("OPENFGA_BASE_URL", "").rstrip("/")
    store_id = os.getenv("OPENFGA_STORE_ID", "")
    if not base_url or not store_id:
        return False
    relation = "write" if _requires_correlation(request.method) else "read"
    target_object = f"tenant:{_normalize(tenant_slug)}" if tenant_slug else f"service:{_normalize(service_name)}"
    user = subject if subject.startswith("service:") else f"user:{subject}"
    payload: dict[str, Any] = {
        "tuple_key": {"user": user, "relation": relation, "object": target_object},
    }
    if model_id := os.getenv("OPENFGA_AUTHORIZATION_MODEL_ID"):
        payload["authorization_model_id"] = model_id
    body = json.dumps(payload).encode()
    api_request = urllib.request.Request(
        f"{base_url}/stores/{store_id}/check",
        data=body,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(api_request, timeout=2) as response:
            result = json.loads(response.read().decode())
            return bool(result.get("allowed"))
    except (urllib.error.URLError, TimeoutError, json.JSONDecodeError):
        return False


def _resolve_tenant(request: Request) -> str:
    return request.headers.get("x-tenant-slug") or request.headers.get("x-erp-tenant-slug") or request.query_params.get("tenant_slug", "")


def _requires_correlation(method: str) -> bool:
    return method.upper() not in {"GET", "HEAD", "OPTIONS"}


def _base64url_decode(value: str) -> bytes:
    return base64.urlsafe_b64decode(value + "=" * (-len(value) % 4))


def _base64url_encode(value: bytes) -> str:
    return base64.urlsafe_b64encode(value).decode().rstrip("=")


def _normalize(value: str) -> str:
    return value.strip().lower().replace(" ", "-")


def _error(status_code: int, code: str, message: str) -> JSONResponse:
    return JSONResponse(status_code=status_code, content={"code": code, "message": message})
