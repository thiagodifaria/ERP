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
    extract::{Path, Query, State},
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
        .route("/api/webhook-hub/events/:public_id", get(get_event_by_public_id))
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

async fn list_events(
    State(state): State<WebhookHubState>,
    Query(filters): Query<ListWebhookEventsQuery>,
) -> Json<Vec<WebhookEvent>> {
    Json(
        state
            .list_events_filtered(filters.provider, filters.event_type, filters.status)
            .await,
    )
}

async fn create_event(
    State(state): State<WebhookHubState>,
    Json(payload): Json<CreateWebhookEventRequest>,
) -> Result<(StatusCode, Json<WebhookEvent>), (StatusCode, Json<ErrorResponse>)> {
    let provider = payload.provider.trim().to_lowercase();
    let event_type = payload.event_type.trim().to_lowercase();
    let external_id = payload.external_id.trim().to_string();

    if provider.is_empty()
        || event_type.is_empty()
        || external_id.is_empty()
    {
        return Err((
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse::new(
                "webhook_event_payload_invalid",
                "Webhook event payload is invalid.",
            )),
        ));
    }

    if state
        .find_event_by_provider_and_external_id(&provider, &external_id)
        .await
        .is_some()
    {
        return Err((
            StatusCode::CONFLICT,
            Json(ErrorResponse::new(
                "webhook_event_conflict",
                "Webhook event with this provider and external id already exists.",
            )),
        ));
    }

    let webhook_event = state
        .push_event(provider, event_type, external_id, payload.payload_summary)
        .await;

    Ok((StatusCode::CREATED, Json(webhook_event)))
}

async fn get_event_summary(State(state): State<WebhookHubState>) -> Json<WebhookEventSummary> {
    Json(state.summary().await)
}

async fn get_event_by_public_id(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    match state.find_event_by_public_id(&public_id).await {
        Some(webhook_event) => Ok(Json(webhook_event)),
        None => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse::new(
                "webhook_event_not_found",
                "Webhook event was not found.",
            )),
        )),
    }
}

#[derive(Clone, Default)]
struct WebhookHubState {
    next_id: Arc<AtomicU64>,
    events: Arc<RwLock<Vec<WebhookEvent>>>,
}

impl WebhookHubState {
    async fn list_events_filtered(
        &self,
        provider: Option<String>,
        event_type: Option<String>,
        status: Option<String>,
    ) -> Vec<WebhookEvent> {
        let provider = provider.map(|value| value.trim().to_lowercase());
        let event_type = event_type.map(|value| value.trim().to_lowercase());
        let status = status.map(|value| value.trim().to_lowercase());

        self.events
            .read()
            .await
            .iter()
            .filter(|event| {
                if let Some(provider_filter) = &provider {
                    if &event.provider != provider_filter {
                        return false;
                    }
                }

                if let Some(event_type_filter) = &event_type {
                    if &event.event_type != event_type_filter {
                        return false;
                    }
                }

                if let Some(status_filter) = &status {
                    if &event.status != status_filter {
                        return false;
                    }
                }

                true
            })
            .cloned()
            .collect()
    }

    async fn find_event_by_public_id(&self, public_id: &str) -> Option<WebhookEvent> {
        self.events
            .read()
            .await
            .iter()
            .find(|event| event.public_id == public_id)
            .cloned()
    }

    async fn find_event_by_provider_and_external_id(
        &self,
        provider: &str,
        external_id: &str,
    ) -> Option<WebhookEvent> {
        self.events
            .read()
            .await
            .iter()
            .find(|event| event.provider == provider && event.external_id == external_id)
            .cloned()
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
            provider,
            event_type,
            external_id,
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

#[derive(Default, serde::Deserialize)]
struct ListWebhookEventsQuery {
    provider: Option<String>,
    event_type: Option<String>,
    status: Option<String>,
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

    #[tokio::test]
    async fn webhook_event_detail_should_return_created_resource_by_public_id() {
        let app = build_router();

        let create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"shopify","event_type":"order.created","external_id":"ord_001","payload_summary":"Pedido criado."}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(create_response.status(), StatusCode::CREATED);

        let create_body = create_response.into_body().collect().await.unwrap().to_bytes();
        let create_payload: Value = serde_json::from_slice(&create_body).unwrap();
        let public_id = create_payload["public_id"].as_str().unwrap();

        let detail_response = app
            .oneshot(
                Request::builder()
                    .uri(format!("/api/webhook-hub/events/{public_id}"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(detail_response.status(), StatusCode::OK);

        let detail_body = detail_response.into_body().collect().await.unwrap().to_bytes();
        let detail_payload: Value = serde_json::from_slice(&detail_body).unwrap();

        assert_eq!(detail_payload["public_id"], public_id);
        assert_eq!(detail_payload["provider"], "shopify");
        assert_eq!(detail_payload["event_type"], "order.created");
        assert_eq!(detail_payload["external_id"], "ord_001");
    }

    #[tokio::test]
    async fn webhook_event_list_should_support_provider_event_type_and_status_filters() {
        let app = build_router();

        let first_create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"stripe","event_type":"payment.succeeded","external_id":"evt_filter_001"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(first_create_response.status(), StatusCode::CREATED);

        let second_create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"tiny","event_type":"lead.created","external_id":"evt_filter_002"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(second_create_response.status(), StatusCode::CREATED);

        let filtered_response = app
            .oneshot(
                Request::builder()
                    .uri("/api/webhook-hub/events?provider=stripe&event_type=payment.succeeded&status=received")
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(filtered_response.status(), StatusCode::OK);

        let filtered_body = filtered_response.into_body().collect().await.unwrap().to_bytes();
        let filtered_payload: Value = serde_json::from_slice(&filtered_body).unwrap();

        assert_eq!(filtered_payload.as_array().unwrap().len(), 1);
        assert_eq!(filtered_payload[0]["provider"], "stripe");
        assert_eq!(filtered_payload[0]["event_type"], "payment.succeeded");
        assert_eq!(filtered_payload[0]["status"], "received");
    }

    #[tokio::test]
    async fn webhook_event_create_should_reject_duplicate_provider_and_external_id() {
        let app = build_router();

        let first_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"stripe","event_type":"payment.failed","external_id":"evt_duplicate_001"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(first_response.status(), StatusCode::CREATED);

        let duplicate_response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"stripe","event_type":"payment.failed","external_id":"evt_duplicate_001"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(duplicate_response.status(), StatusCode::CONFLICT);

        let duplicate_body = duplicate_response.into_body().collect().await.unwrap().to_bytes();
        let duplicate_payload: Value = serde_json::from_slice(&duplicate_body).unwrap();

        assert_eq!(duplicate_payload["code"], "webhook_event_conflict");
    }
}
