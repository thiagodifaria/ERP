from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse


def install_security(app: FastAPI, service_name: str) -> None:
    @app.middleware("http")
    async def security_middleware(request: Request, call_next):
        if request.url.path.startswith("/health/"):
            return await call_next(request)

        auth = authenticate_request(request)
        if auth is None:
            return JSONResponse(status_code=401, content={"code": "unauthorized", "message": "Bearer token is invalid or missing."})
        if request.method not in {"GET", "HEAD", "OPTIONS"} and not request.headers.get("x-correlation-id"):
            return JSONResponse(status_code=400, content={"code": "correlation_id_required", "message": "Mutation requests require X-Correlation-Id."})
        if not authorize_request(service_name, request, auth):
            return JSONResponse(status_code=403, content={"code": "forbidden", "message": "Request is not authorized."})

        request.scope["erp_auth"] = auth
        return await call_next(request)


def authenticate_request(request: Request) -> dict | None:
    return None


def authorize_request(service_name: str, request: Request, auth: dict) -> bool:
    return bool(service_name and auth.get("subject"))
