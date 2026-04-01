// Router define as rotas publicas minimas do servico.
// Regras de protecao de entrada crescem a partir deste ponto.
use std::{
    collections::BTreeMap,
    sync::{
        atomic::{AtomicU64, Ordering},
        Arc,
    },
};

use axum::{
    extract::State,
    http::StatusCode,
    routing::get,
    Json, Router,
};
use chrono::Utc;
use tokio::sync::RwLock;
use uuid::Uuid;

pub fn build_router() -> Router {
    let state = WebhookHubState::default();

    Router::new()
        .route("/health/live", get(live))
        .route("/health/ready", get(ready))
        .route("/health/details", get(details))
        .route("/api/webhook-hub/events", get(list_events).post(create_event))
        .route("/api/webhook-hub/events/summary", get(get_event_summary))
        .with_state(state)
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

async fn list_events(State(state): State<WebhookHubState>) -> Json<Vec<WebhookEvent>> {
    Json(state.list_events().await)
}

async fn create_event(
    State(state): State<WebhookHubState>,
    Json(payload): Json<CreateWebhookEventRequest>,
) -> Result<(StatusCode, Json<WebhookEvent>), (StatusCode, Json<ErrorResponse>)> {
    if payload.provider.trim().is_empty()
        || payload.event_type.trim().is_empty()
        || payload.external_id.trim().is_empty()
    {
        return Err((
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse::new(
                "webhook_event_payload_invalid",
                "Webhook event payload is invalid.",
            )),
        ));
    }

    let webhook_event = state
        .push_event(payload.provider, payload.event_type, payload.external_id, payload.payload_summary)
        .await;

    Ok((StatusCode::CREATED, Json(webhook_event)))
}

async fn get_event_summary(State(state): State<WebhookHubState>) -> Json<WebhookEventSummary> {
    Json(state.summary().await)
}

#[derive(Clone, Default)]
struct WebhookHubState {
    next_id: Arc<AtomicU64>,
    events: Arc<RwLock<Vec<WebhookEvent>>>,
}

impl WebhookHubState {
    async fn list_events(&self) -> Vec<WebhookEvent> {
        self.events.read().await.clone()
    }

    async fn push_event(
        &self,
        provider: String,
        event_type: String,
        external_id: String,
        payload_summary: Option<String>,
    ) -> WebhookEvent {
        let webhook_event = WebhookEvent {
            id: self.next_id.fetch_add(1, Ordering::SeqCst) + 1,
            public_id: Uuid::new_v4().to_string(),
            provider: provider.trim().to_lowercase(),
            event_type: event_type.trim().to_lowercase(),
            external_id: external_id.trim().to_string(),
            payload_summary: payload_summary.map(|value| value.trim().to_string()),
            status: "received".to_string(),
            received_at: Utc::now().to_rfc3339(),
        };

        self.events.write().await.push(webhook_event.clone());
        webhook_event
    }

    async fn summary(&self) -> WebhookEventSummary {
        let events = self.events.read().await;
        let mut by_provider = BTreeMap::new();

        for event in events.iter() {
            *by_provider.entry(event.provider.clone()).or_insert(0) += 1;
        }

        WebhookEventSummary {
            total: events.len() as u64,
            received: events
                .iter()
                .filter(|event| event.status == "received")
                .count() as u64,
            by_provider,
        }
    }
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

#[derive(Clone, serde::Serialize)]
struct WebhookEvent {
    id: u64,
    public_id: String,
    provider: String,
    event_type: String,
    external_id: String,
    payload_summary: Option<String>,
    status: String,
    received_at: String,
}

#[derive(serde::Deserialize)]
struct CreateWebhookEventRequest {
    provider: String,
    event_type: String,
    external_id: String,
    payload_summary: Option<String>,
}

#[derive(serde::Serialize)]
struct WebhookEventSummary {
    total: u64,
    received: u64,
    by_provider: BTreeMap<String, u64>,
}

#[derive(serde::Serialize)]
struct ErrorResponse {
    code: &'static str,
    message: &'static str,
}

impl ErrorResponse {
    fn new(code: &'static str, message: &'static str) -> Self {
        Self { code, message }
    }
}

#[cfg(test)]
mod tests {
    use super::build_router;
    use axum::http::{Request, StatusCode};
    use http_body_util::BodyExt;
    use serde_json::Value;
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

    #[tokio::test]
    async fn webhook_event_create_should_return_created_resource() {
        let app = build_router();
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"stripe","event_type":"payment.succeeded","external_id":"evt_123","payload_summary":"Pagamento confirmado."}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::CREATED);

        let body = response.into_body().collect().await.unwrap().to_bytes();
        let payload: Value = serde_json::from_slice(&body).unwrap();

        assert_eq!(payload["provider"], "stripe");
        assert_eq!(payload["event_type"], "payment.succeeded");
        assert_eq!(payload["external_id"], "evt_123");
        assert_eq!(payload["status"], "received");
    }

    #[tokio::test]
    async fn webhook_event_routes_should_list_and_summarize_live_ingestion() {
        let app = build_router();

        let create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"tiny","event_type":"lead.created","external_id":"wh_001","payload_summary":"Lead recebido do formulario externo."}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(create_response.status(), StatusCode::CREATED);

        let list_response = app
            .clone()
            .oneshot(Request::builder().uri("/api/webhook-hub/events").body(axum::body::Body::empty()).unwrap())
            .await
            .unwrap();
        assert_eq!(list_response.status(), StatusCode::OK);

        let list_body = list_response.into_body().collect().await.unwrap().to_bytes();
        let list_payload: Value = serde_json::from_slice(&list_body).unwrap();
        assert_eq!(list_payload.as_array().unwrap().len(), 1);
        assert_eq!(list_payload[0]["provider"], "tiny");

        let summary_response = app
            .oneshot(
                Request::builder()
                    .uri("/api/webhook-hub/events/summary")
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(summary_response.status(), StatusCode::OK);

        let summary_body = summary_response.into_body().collect().await.unwrap().to_bytes();
        let summary_payload: Value = serde_json::from_slice(&summary_body).unwrap();
        assert_eq!(summary_payload["total"], 1);
        assert_eq!(summary_payload["received"], 1);
        assert_eq!(summary_payload["by_provider"]["tiny"], 1);
    }
}
