// Router define as rotas publicas minimas do servico.
// Regras de protecao de entrada crescem a partir deste ponto.
use axum::{routing::get, Json, Router};

pub fn build_router() -> Router {
    Router::new()
        .route("/health/live", get(live))
        .route("/health/ready", get(ready))
}

async fn live() -> Json<HealthResponse> {
    Json(HealthResponse::new("webhook-hub", "live"))
}

async fn ready() -> Json<HealthResponse> {
    Json(HealthResponse::new("webhook-hub", "ready"))
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
