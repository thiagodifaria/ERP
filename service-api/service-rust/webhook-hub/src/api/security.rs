use std::sync::OnceLock;

use axum::{
    extract::Request,
    http::{header, HeaderValue, StatusCode},
    middleware::Next,
    response::Response,
};
use base64::{engine::general_purpose::URL_SAFE_NO_PAD, Engine as _};
use chrono::Utc;
use hmac::{digest::subtle::ConstantTimeEq, Hmac, Mac};
use serde_json::Value;
use sha2::Sha256;

static OPENFGA_HTTP_CLIENT: OnceLock<reqwest::Client> = OnceLock::new();

pub async fn security_middleware(mut request: Request, next: Next) -> Response {
    if !security_enforced() || request.uri().path().starts_with("/health/") {
        return next.run(request).await;
    }
    if requires_correlation(request.method().as_str())
        && request
            .headers()
            .get("x-correlation-id")
            .and_then(|value| value.to_str().ok())
            .unwrap_or("")
            .trim()
            .is_empty()
    {
        return super::router::security_error(StatusCode::BAD_REQUEST, "correlation_id_required", "Mutation requests require X-Correlation-Id.");
    }
    if requires_json_body(request.method().as_str())
        && request
            .headers()
            .get(header::CONTENT_TYPE)
            .and_then(|value| value.to_str().ok())
            .map(|value| value.to_ascii_lowercase().starts_with("application/json"))
            != Some(true)
    {
        return super::router::security_error(StatusCode::UNSUPPORTED_MEDIA_TYPE, "content_type_required", "Mutation requests require application/json.");
    }
    let Some(auth) = authenticate(&request) else {
        return super::router::security_error(StatusCode::UNAUTHORIZED, "unauthorized", "Bearer token is invalid or missing.");
    };
    if let Ok(value) = HeaderValue::from_str(&auth.subject) {
        request.headers_mut().insert("x-erp-auth-subject", value);
    }
    if let Ok(value) = HeaderValue::from_str(&auth.tenant_slug) {
        request.headers_mut().insert("x-erp-auth-tenant", value);
    }
    if let Ok(value) = HeaderValue::from_str(&auth.scopes.join(" ")) {
        request.headers_mut().insert("x-erp-auth-scopes", value);
    }
    let request_method = request.method().as_str().to_string();
    if !authorize_openfga("webhook-hub", &request_method, &auth).await {
        return super::router::security_error(StatusCode::FORBIDDEN, "openfga_denied", "OpenFGA denied the request.");
    }
    next.run(request).await
}

struct AuthContext {
    subject: String,
    tenant_slug: String,
    scopes: Vec<String>,
}

fn security_enforced() -> bool {
    let mode = std::env::var("ERP_AUTH_ENFORCEMENT").unwrap_or_default().trim().to_lowercase();
    if matches!(mode.as_str(), "disabled" | "off" | "false") {
        return false;
    }
    if matches!(mode.as_str(), "enforced" | "strict" | "true") {
        return true;
    }
    let environment = std::env::var("ERP_ENV").unwrap_or_else(|_| "local".to_string()).trim().to_lowercase();
    !matches!(environment.as_str(), "" | "local" | "dev" | "development" | "test" | "testing")
}

fn authenticate(request: &Request) -> Option<AuthContext> {
    let authorization = request.headers().get(header::AUTHORIZATION)?.to_str().ok()?;
    if !authorization.to_lowercase().starts_with("bearer ") {
        return None;
    }
    let token = authorization["Bearer ".len()..].trim();
    let internal_token = std::env::var("ERP_INTERNAL_SERVICE_TOKEN").unwrap_or_default();
    if !internal_token.trim().is_empty() && constant_time_equals(token, internal_token.trim()) {
        return Some(AuthContext {
            subject: "service:internal".to_string(),
            tenant_slug: resolve_tenant(request),
            scopes: vec!["service".to_string()],
        });
    }
    let claims = verify_jwt(token)?;
    let subject = read_claim(&claims, "sub").or_else(|| read_claim(&claims, "user_public_id"))?;
    let tenant_slug = read_claim(&claims, "tenant_slug")
        .or_else(|| read_claim(&claims, "tenant"))
        .unwrap_or_else(|| resolve_tenant(request));
    let scopes = match claims.get("scope") {
        Some(Value::String(value)) => value.split_whitespace().map(str::to_string).collect(),
        Some(Value::Array(values)) => values.iter().filter_map(Value::as_str).map(str::to_string).collect(),
        _ => vec![],
    };
    Some(AuthContext { subject, tenant_slug, scopes })
}

fn verify_jwt(token: &str) -> Option<Value> {
    let secret = std::env::var("ERP_JWT_HS256_SECRET").ok()?;
    if secret.trim().is_empty() {
        return None;
    }
    let parts: Vec<&str> = token.split('.').collect();
    if parts.len() != 3 {
        return None;
    }
    let header_value: Value = serde_json::from_slice(&URL_SAFE_NO_PAD.decode(parts[0]).ok()?).ok()?;
    if header_value.get("alg").and_then(Value::as_str) != Some("HS256") {
        return None;
    }
    type HmacSha256 = Hmac<Sha256>;
    let mut mac = HmacSha256::new_from_slice(secret.as_bytes()).ok()?;
    mac.update(format!("{}.{}", parts[0], parts[1]).as_bytes());
    if URL_SAFE_NO_PAD.encode(mac.finalize().into_bytes()) != parts[2] {
        return None;
    }
    let claims: Value = serde_json::from_slice(&URL_SAFE_NO_PAD.decode(parts[1]).ok()?).ok()?;
    if let Some(expires_at) = claims.get("exp").and_then(Value::as_i64) {
        if expires_at <= Utc::now().timestamp() {
            return None;
        }
    }
    Some(claims)
}

async fn authorize_openfga(service_name: &str, method: &str, auth: &AuthContext) -> bool {
    if std::env::var("ERP_OPENFGA_ENFORCEMENT").unwrap_or_default().to_lowercase() != "true" {
        return true;
    }
    let base_url = std::env::var("OPENFGA_BASE_URL").unwrap_or_default().trim_end_matches('/').to_string();
    let store_id = std::env::var("OPENFGA_STORE_ID").unwrap_or_default();
    if base_url.is_empty() || store_id.is_empty() {
        return false;
    }
    let relation = if requires_correlation(method) { "write" } else { "read" };
    let object = if auth.tenant_slug.is_empty() {
        format!("service:{}", normalize_object(service_name))
    } else {
        format!("tenant:{}", normalize_object(&auth.tenant_slug))
    };
    let mut payload = serde_json::json!({
        "tuple_key": {
            "user": if auth.subject.starts_with("service:") { auth.subject.clone() } else { format!("user:{}", auth.subject) },
            "relation": relation,
            "object": object
        }
    });
    if let Ok(model_id) = std::env::var("OPENFGA_AUTHORIZATION_MODEL_ID") {
        payload["authorization_model_id"] = Value::String(model_id);
    }
    let client = OPENFGA_HTTP_CLIENT.get_or_init(reqwest::Client::new);
    let Ok(response) = client
        .post(format!("{}/stores/{}/check", base_url, store_id))
        .json(&payload)
        .timeout(std::time::Duration::from_secs(2))
        .send()
        .await
    else {
        return false;
    };
    if !response.status().is_success() {
        return false;
    }
    response
        .json::<Value>()
        .await
        .ok()
        .and_then(|value| value.get("allowed").and_then(Value::as_bool))
        .unwrap_or(false)
}

fn constant_time_equals(left: &str, right: &str) -> bool {
    let left_bytes = left.as_bytes();
    let right_bytes = right.as_bytes();
    left_bytes.len() == right_bytes.len() && left_bytes.ct_eq(right_bytes).into()
}

fn requires_correlation(method: &str) -> bool {
    !matches!(method, "GET" | "HEAD" | "OPTIONS")
}

fn requires_json_body(method: &str) -> bool {
    matches!(method, "POST" | "PUT" | "PATCH")
}

fn resolve_tenant(request: &Request) -> String {
    for header_name in ["x-tenant-slug", "x-erp-tenant-slug"] {
        if let Some(value) = request.headers().get(header_name).and_then(|value| value.to_str().ok()) {
            if !value.trim().is_empty() {
                return value.trim().to_string();
            }
        }
    }
    request
        .uri()
        .query()
        .and_then(|query| query.split('&').find_map(|part| part.strip_prefix("tenant_slug=")))
        .unwrap_or("")
        .to_string()
}

fn read_claim(value: &Value, name: &str) -> Option<String> {
    value.get(name).and_then(Value::as_str).map(str::to_string)
}

fn normalize_object(value: &str) -> String {
    value.trim().to_lowercase().replace(' ', "-")
}
