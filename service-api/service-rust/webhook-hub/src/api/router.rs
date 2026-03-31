// Router define as rotas publicas minimas do servico.
// Regras de protecao de entrada crescem a partir deste ponto.
use axum::{routing::get, Json, Router};

pub fn build_router() -> Router {
    Router::new()
        .route("/health/live", get(live))
        .route("/health/ready", get(ready))
        .route("/health/details", get(details))
}

async fn live() -> Json<HealthResponse> {
    Json(HealthResponse::new("webhook-hub", "live"))
}

async fn ready() -> Json<HealthResponse> {
    Json(HealthResponse::new("webhook-hub", "ready"))
}

async fn details() -> Json<ReadinessResponse> {
    Json(ReadinessResponse::new(
        "webhook-hub",
        "ready",
        vec![
            DependencyHealth::new("signature-validation", "ready"),
            DependencyHealth::new("postgresql", "pending-runtime-wiring"),
            DependencyHealth::new("redis", "pending-runtime-wiring"),
        ],
    ))
}

#[derive(serde::Serialize)]
struct HealthResponse {
    service: &'static str,
    status: &'static str,
}

impl HealthResponse {
    fn new(service: &'static str, status: &'static str) -> Self {
        Self { service, status }
    }
}

#[derive(serde::Serialize)]
struct ReadinessResponse {
    service: &'static str,
    status: &'static str,
    dependencies: Vec<DependencyHealth>,
}

impl ReadinessResponse {
    fn new(service: &'static str, status: &'static str, dependencies: Vec<DependencyHealth>) -> Self {
        Self {
            service,
            status,
            dependencies,
        }
    }
}

#[derive(serde::Serialize)]
struct DependencyHealth {
    name: &'static str,
    status: &'static str,
}

impl DependencyHealth {
    fn new(name: &'static str, status: &'static str) -> Self {
        Self { name, status }
    }
}

#[cfg(test)]
mod tests {
    use super::build_router;
    use axum::http::{Request, StatusCode};
    use tower::ServiceExt;

    #[tokio::test]
    async fn live_route_should_return_ok() {
        let app = build_router();
        let response = app
            .oneshot(Request::builder().uri("/health/live").body(axum::body::Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn ready_route_should_return_ok() {
        let app = build_router();
        let response = app
            .oneshot(Request::builder().uri("/health/ready").body(axum::body::Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn details_route_should_return_ok() {
        let app = build_router();
        let response = app
            .oneshot(Request::builder().uri("/health/details").body(axum::body::Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }
}
