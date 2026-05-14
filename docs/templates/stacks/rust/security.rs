use axum::{
    extract::Request,
    http::StatusCode,
    middleware::Next,
    response::{IntoResponse, Response},
    Json,
};
use serde_json::json;

pub async fn security_middleware(mut request: Request, next: Next) -> Response {
    if request.uri().path().starts_with("/health/") {
        return next.run(request).await;
    }

    let Some(auth) = authenticate_request(&request) else {
        return error(StatusCode::UNAUTHORIZED, "unauthorized", "Bearer token is invalid or missing.");
    };
    if request.method().as_str() != "GET" && request.headers().get("x-correlation-id").is_none() {
        return error(StatusCode::BAD_REQUEST, "correlation_id_required", "Mutation requests require X-Correlation-Id.");
    }
    if !authorize_request(&request, &auth) {
        return error(StatusCode::FORBIDDEN, "forbidden", "Request is not authorized.");
    }

    request.headers_mut().insert("x-erp-auth-subject", auth.subject.parse().unwrap());
    next.run(request).await
}

pub struct AuthContext {
    pub subject: String,
    pub tenant_slug: String,
    pub scopes: Vec<String>,
}

fn authenticate_request(_request: &Request) -> Option<AuthContext> {
    None
}

fn authorize_request(_request: &Request, auth: &AuthContext) -> bool {
    !auth.subject.is_empty()
}

fn error(status: StatusCode, code: &str, message: &str) -> Response {
    (status, Json(json!({ "code": code, "message": message }))).into_response()
}
