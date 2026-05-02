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
    routing::{get, post},
    Json, Router,
};
use chrono::Utc;
use tokio::sync::RwLock;
use tokio_postgres::Row;
use uuid::Uuid;

use crate::{
    config::{AppConfig, RepositoryDriver},
    infrastructure::persistence::connect_postgres,
};

pub fn build_router() -> Router {
    build_router_with_state(WebhookHubState::memory())
}

pub async fn build_router_from_config(
    config: &AppConfig,
) -> Result<Router, Box<dyn std::error::Error + Send + Sync>> {
    let state = WebhookHubState::from_config(config).await?;
    Ok(build_router_with_state(state))
}

fn build_router_with_state(state: WebhookHubState) -> Router {
    Router::new()
        .route("/health/live", get(live))
        .route("/health/ready", get(ready))
        .route("/health/details", get(details))
        .route("/api/webhook-hub/capabilities", get(get_capabilities))
        .route("/api/webhook-hub/events", get(list_events).post(create_event))
        .route("/api/webhook-hub/events/summary", get(get_event_summary))
        .route("/api/webhook-hub/events/dead-letter", get(list_dead_letter_events))
        .route("/api/webhook-hub/events/:public_id", get(get_event_by_public_id))
        .route("/api/webhook-hub/events/:public_id/transitions", get(list_event_transitions))
        .route("/api/webhook-hub/events/:public_id/validate", post(validate_event))
        .route("/api/webhook-hub/events/:public_id/queue", post(queue_event))
        .route("/api/webhook-hub/events/:public_id/process", post(process_event))
        .route("/api/webhook-hub/events/:public_id/forward", post(forward_event))
        .route("/api/webhook-hub/events/:public_id/fail", post(fail_event))
        .route("/api/webhook-hub/events/:public_id/dead-letter", post(dead_letter_event))
        .route("/api/webhook-hub/events/:public_id/requeue", post(requeue_event))
        .route("/api/webhook-hub/events/:public_id/reject", post(reject_event))
        .with_state(state)
}

async fn live() -> Json<HealthResponse> {
    Json(HealthResponse::new("webhook-hub", "live"))
}

async fn ready() -> Json<HealthResponse> {
    Json(HealthResponse::new("webhook-hub", "ready"))
}

async fn details(State(state): State<WebhookHubState>) -> Json<ReadinessResponse> {
    Json(ReadinessResponse::new("webhook-hub", "ready", state.dependencies()))
}

async fn get_capabilities(State(state): State<WebhookHubState>) -> Json<WebhookHubCapabilitiesResponse> {
    Json(state.capabilities())
}

async fn list_events(
    State(state): State<WebhookHubState>,
    Query(filters): Query<ListWebhookEventsQuery>,
) -> Result<Json<Vec<WebhookEvent>>, (StatusCode, Json<ErrorResponse>)> {
    let webhook_events = state
        .list_events_filtered(filters.provider, filters.event_type, filters.status)
        .await
        .map_err(internal_storage_error)?;

    Ok(Json(webhook_events))
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
        .map_err(internal_storage_error)?
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
        .await
        .map_err(internal_storage_error)?;

    Ok((StatusCode::CREATED, Json(webhook_event)))
}

async fn get_event_summary(
    State(state): State<WebhookHubState>,
) -> Result<Json<WebhookEventSummary>, (StatusCode, Json<ErrorResponse>)> {
    let summary = state.summary().await.map_err(internal_storage_error)?;
    Ok(Json(summary))
}

async fn list_dead_letter_events(
    State(state): State<WebhookHubState>,
    Query(filters): Query<ListDeadLetterEventsQuery>,
) -> Result<Json<Vec<WebhookEvent>>, (StatusCode, Json<ErrorResponse>)> {
    let webhook_events = state
        .list_events_filtered(filters.provider, None, Some("dead_letter".to_string()))
        .await
        .map_err(internal_storage_error)?;

    Ok(Json(webhook_events))
}

async fn get_event_by_public_id(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    match state
        .find_event_by_public_id(&public_id)
        .await
        .map_err(internal_storage_error)?
    {
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

async fn list_event_transitions(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<Vec<WebhookEventStatusTransition>>, (StatusCode, Json<ErrorResponse>)> {
    match state
        .list_event_transitions(&public_id)
        .await
        .map_err(internal_storage_error)?
    {
        Some(status_history) => Ok(Json(status_history)),
        None => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse::new(
                "webhook_event_not_found",
                "Webhook event was not found.",
            )),
        )),
    }
}

async fn validate_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    map_transition_result(
        state
            .transition_event_status(&public_id, "validated", &["received"])
            .await,
    )
}

async fn reject_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    map_transition_result(
        state
            .transition_event_status(&public_id, "rejected", &["received", "validated"])
            .await,
    )
}

async fn queue_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    map_transition_result(
        state
            .transition_event_status(&public_id, "queued", &["validated"])
            .await,
    )
}

async fn process_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    map_transition_result(
        state
            .transition_event_status(&public_id, "processing", &["queued"])
            .await,
    )
}

async fn forward_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    map_transition_result(
        state
            .transition_event_status(&public_id, "forwarded", &["processing"])
            .await,
    )
}

async fn fail_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    map_transition_result(
        state
            .transition_event_status(&public_id, "failed", &["queued", "processing"])
            .await,
    )
}

async fn dead_letter_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
    Json(payload): Json<MoveToDeadLetterRequest>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    match state
        .move_to_dead_letter(
            &public_id,
            payload.error_code,
            payload.error_message,
            &["validated", "queued", "processing", "failed"],
        )
        .await
    {
        Ok(webhook_event) => Ok(Json(webhook_event)),
        Err(TransitionError::NotFound) => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse::new(
                "webhook_event_not_found",
                "Webhook event was not found.",
            )),
        )),
        Err(TransitionError::InvalidTransition) => Err((
            StatusCode::CONFLICT,
            Json(ErrorResponse::new(
                "webhook_event_transition_invalid",
                "Webhook event cannot transition from the current status.",
            )),
        )),
        Err(TransitionError::Storage) => Err(internal_storage_error(StorageError)),
    }
}

async fn requeue_event(
    State(state): State<WebhookHubState>,
    Path(public_id): Path<String>,
    Json(payload): Json<RequeueWebhookEventRequest>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    match state
        .requeue_event(&public_id, payload.reason, &["failed", "dead_letter"])
        .await
    {
        Ok(webhook_event) => Ok(Json(webhook_event)),
        Err(TransitionError::NotFound) => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse::new(
                "webhook_event_not_found",
                "Webhook event was not found.",
            )),
        )),
        Err(TransitionError::InvalidTransition) => Err((
            StatusCode::CONFLICT,
            Json(ErrorResponse::new(
                "webhook_event_transition_invalid",
                "Webhook event cannot transition from the current status.",
            )),
        )),
        Err(TransitionError::Storage) => Err(internal_storage_error(StorageError)),
    }
}

fn map_transition_result(
    result: Result<WebhookEvent, TransitionError>,
) -> Result<Json<WebhookEvent>, (StatusCode, Json<ErrorResponse>)> {
    match result {
        Ok(webhook_event) => Ok(Json(webhook_event)),
        Err(TransitionError::NotFound) => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse::new(
                "webhook_event_not_found",
                "Webhook event was not found.",
            )),
        )),
        Err(TransitionError::InvalidTransition) => Err((
            StatusCode::CONFLICT,
            Json(ErrorResponse::new(
                "webhook_event_transition_invalid",
                "Webhook event cannot transition from the current status.",
            )),
        )),
        Err(TransitionError::Storage) => Err(internal_storage_error(StorageError)),
    }
}

#[derive(Clone)]
struct WebhookHubState {
    storage: StorageBackend,
    outbound_signing_ready: bool,
}

#[derive(Clone)]
enum StorageBackend {
    Memory(Arc<MemoryStorage>),
    Postgres(PostgresStorage),
}

#[derive(Default)]
struct MemoryStorage {
    next_id: AtomicU64,
    events: RwLock<Vec<WebhookEvent>>,
}

#[derive(Clone)]
struct PostgresStorage {
    client: Arc<tokio_postgres::Client>,
}

impl WebhookHubState {
    fn memory() -> Self {
        Self {
            storage: StorageBackend::Memory(Arc::new(MemoryStorage::default())),
            outbound_signing_ready: false,
        }
    }

    async fn from_config(
        config: &AppConfig,
    ) -> Result<Self, Box<dyn std::error::Error + Send + Sync>> {
        match config.repository_driver {
            RepositoryDriver::Memory => Ok(Self::memory()),
            RepositoryDriver::Postgres => {
                let client = connect_postgres(&config.postgres).await?;
                Ok(Self {
                    storage: StorageBackend::Postgres(PostgresStorage { client }),
                    outbound_signing_ready: !config.outbound_signing_secret.trim().is_empty(),
                })
            }
        }
    }

    fn dependencies(&self) -> Vec<DependencyHealth> {
        let postgresql_status = match &self.storage {
            StorageBackend::Memory(_) => "pending-runtime-wiring",
            StorageBackend::Postgres(_) => "ready",
        };

        vec![
            DependencyHealth::new("signature-validation", "ready"),
            DependencyHealth::new("postgresql", postgresql_status),
            DependencyHealth::new("redis", "pending-runtime-wiring"),
            DependencyHealth::new("dlq-engine", "ready"),
            DependencyHealth::new(
                "outbound-webhooks",
                if self.outbound_signing_ready {
                    "ready"
                } else {
                    "unconfigured"
                },
            ),
        ]
    }

    fn capabilities(&self) -> WebhookHubCapabilitiesResponse {
        WebhookHubCapabilitiesResponse {
            outbound_webhooks: OutboundWebhookCapability {
                signing_ready: self.outbound_signing_ready,
                status: if self.outbound_signing_ready {
                    "ready".to_string()
                } else {
                    "unconfigured".to_string()
                },
                supports_retries: true,
                supports_dlq: true,
                credential_key: "WEBHOOK_HUB_OUTBOUND_SIGNING_SECRET".to_string(),
            },
        }
    }

    async fn list_events_filtered(
        &self,
        provider: Option<String>,
        event_type: Option<String>,
        status: Option<String>,
    ) -> Result<Vec<WebhookEvent>, StorageError> {
        match &self.storage {
            StorageBackend::Memory(memory) => {
                let provider = provider.map(|value| value.trim().to_lowercase());
                let event_type = event_type.map(|value| value.trim().to_lowercase());
                let status = status.map(|value| value.trim().to_lowercase());

                Ok(memory
                    .events
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
                    .collect())
            }
            StorageBackend::Postgres(postgres) => {
                let provider_filter = provider.map(|value| value.trim().to_lowercase());
                let event_type_filter = event_type.map(|value| value.trim().to_lowercase());
                let status_filter = status.map(|value| value.trim().to_lowercase());
                let rows = postgres
                    .client
                    .query(
                        "
                        SELECT id, public_id, provider, event_type, external_id, payload_summary, status, retry_count, last_error_code, last_error_message, dead_lettered_at, received_at
                        FROM webhook_hub.webhook_events
                        WHERE ($1::text IS NULL OR provider = $1)
                          AND ($2::text IS NULL OR event_type = $2)
                          AND ($3::text IS NULL OR status = $3)
                        ORDER BY id ASC;
                        ",
                        &[&provider_filter, &event_type_filter, &status_filter],
                    )
                    .await
                    .map_err(|_| StorageError)?;

                let mut webhook_events = Vec::with_capacity(rows.len());
                for row in rows {
                    webhook_events.push(map_row_to_webhook_event(&postgres.client, row).await?);
                }

                Ok(webhook_events)
            }
        }
    }

    async fn find_event_by_public_id(&self, public_id: &str) -> Result<Option<WebhookEvent>, StorageError> {
        match &self.storage {
            StorageBackend::Memory(memory) => Ok(memory
                .events
                .read()
                .await
                .iter()
                .find(|event| event.public_id == public_id)
                .cloned()),
            StorageBackend::Postgres(postgres) => {
                let Ok(public_id) = Uuid::parse_str(public_id) else {
                    return Ok(None);
                };

                let row = postgres
                    .client
                    .query_opt(
                        "
                        SELECT id, public_id, provider, event_type, external_id, payload_summary, status, retry_count, last_error_code, last_error_message, dead_lettered_at, received_at
                        FROM webhook_hub.webhook_events
                        WHERE public_id = $1;
                        ",
                        &[&public_id],
                    )
                    .await
                    .map_err(|_| StorageError)?;

                match row {
                    Some(row) => Ok(Some(map_row_to_webhook_event(&postgres.client, row).await?)),
                    None => Ok(None),
                }
            }
        }
    }

    async fn find_event_by_provider_and_external_id(
        &self,
        provider: &str,
        external_id: &str,
    ) -> Result<Option<WebhookEvent>, StorageError> {
        match &self.storage {
            StorageBackend::Memory(memory) => Ok(memory
                .events
                .read()
                .await
                .iter()
                .find(|event| event.provider == provider && event.external_id == external_id)
                .cloned()),
            StorageBackend::Postgres(postgres) => {
                let row = postgres
                    .client
                    .query_opt(
                        "
                        SELECT id, public_id, provider, event_type, external_id, payload_summary, status, retry_count, last_error_code, last_error_message, dead_lettered_at, received_at
                        FROM webhook_hub.webhook_events
                        WHERE provider = $1 AND external_id = $2;
                        ",
                        &[&provider, &external_id],
                    )
                    .await
                    .map_err(|_| StorageError)?;

                match row {
                    Some(row) => Ok(Some(map_row_to_webhook_event(&postgres.client, row).await?)),
                    None => Ok(None),
                }
            }
        }
    }

    async fn list_event_transitions(
        &self,
        public_id: &str,
    ) -> Result<Option<Vec<WebhookEventStatusTransition>>, StorageError> {
        match &self.storage {
            StorageBackend::Memory(memory) => Ok(memory
                .events
                .read()
                .await
                .iter()
                .find(|event| event.public_id == public_id)
                .map(|event| event.status_history.clone())),
            StorageBackend::Postgres(postgres) => {
                let Ok(public_id) = Uuid::parse_str(public_id) else {
                    return Ok(None);
                };

                let rows = postgres
                    .client
                    .query(
                        "
                        SELECT transition.status, transition.changed_at
                        FROM webhook_hub.webhook_event_transitions AS transition
                        INNER JOIN webhook_hub.webhook_events AS event
                          ON event.id = transition.webhook_event_id
                        WHERE event.public_id = $1
                        ORDER BY transition.changed_at ASC, transition.id ASC;
                        ",
                        &[&public_id],
                    )
                    .await
                    .map_err(|_| StorageError)?;

                if rows.is_empty() {
                    return Ok(None);
                }

                Ok(Some(
                    rows.into_iter()
                        .map(map_row_to_transition)
                        .collect(),
                ))
            }
        }
    }

    async fn push_event(
        &self,
        provider: String,
        event_type: String,
        external_id: String,
        payload_summary: Option<String>,
    ) -> Result<WebhookEvent, StorageError> {
        match &self.storage {
            StorageBackend::Memory(memory) => {
                let received_at = Utc::now().to_rfc3339();
                let webhook_event = WebhookEvent {
                    id: memory.next_id.fetch_add(1, Ordering::SeqCst) + 1,
                    public_id: Uuid::new_v4().to_string(),
                    provider,
                    event_type,
                    external_id,
                    payload_summary: payload_summary.map(|value| value.trim().to_string()),
                    status: "received".to_string(),
                    retry_count: 0,
                    last_error_code: None,
                    last_error_message: None,
                    dead_lettered_at: None,
                    received_at: received_at.clone(),
                    status_history: vec![WebhookEventStatusTransition::new("received", received_at)],
                };

                memory.events.write().await.push(webhook_event.clone());
                Ok(webhook_event)
            }
            StorageBackend::Postgres(postgres) => {
                let public_id = Uuid::new_v4();
                let received_at = Utc::now();
                let payload_summary = payload_summary.map(|value| value.trim().to_string());

                let row = postgres
                    .client
                    .query_one(
                        "
                        INSERT INTO webhook_hub.webhook_events (
                          public_id, provider, event_type, external_id, payload_summary, status, received_at
                        )
                        VALUES ($1, $2, $3, $4, $5, 'received', $6)
                        RETURNING id, public_id, provider, event_type, external_id, payload_summary, status, retry_count, last_error_code, last_error_message, dead_lettered_at, received_at;
                        ",
                        &[&public_id, &provider, &event_type, &external_id, &payload_summary, &received_at],
                    )
                    .await
                    .map_err(|_| StorageError)?;

                let event_id = row.get::<_, i64>("id");
                postgres
                    .client
                    .execute(
                        "
                        INSERT INTO webhook_hub.webhook_event_transitions (webhook_event_id, status, changed_at)
                        VALUES ($1, 'received', $2);
                        ",
                        &[&event_id, &received_at],
                    )
                    .await
                    .map_err(|_| StorageError)?;

                let event = map_row_to_webhook_event(&postgres.client, row).await?;
                Ok(event)
            }
        }
    }

    async fn transition_event_status(
        &self,
        public_id: &str,
        next_status: &str,
        allowed_statuses: &[&str],
    ) -> Result<WebhookEvent, TransitionError> {
        match &self.storage {
            StorageBackend::Memory(memory) => {
                let mut events = memory.events.write().await;
                let webhook_event = events
                    .iter_mut()
                    .find(|event| event.public_id == public_id)
                    .ok_or(TransitionError::NotFound)?;

                if !allowed_statuses.iter().any(|status| webhook_event.status == *status) {
                    return Err(TransitionError::InvalidTransition);
                }

                webhook_event.status = next_status.to_string();
                webhook_event
                    .status_history
                    .push(WebhookEventStatusTransition::new(next_status, Utc::now().to_rfc3339()));
                Ok(webhook_event.clone())
            }
            StorageBackend::Postgres(postgres) => {
                let Ok(public_id) = Uuid::parse_str(public_id) else {
                    return Err(TransitionError::NotFound);
                };

                let row = postgres
                    .client
                    .query_opt(
                        "
                        SELECT id, status
                        FROM webhook_hub.webhook_events
                        WHERE public_id = $1
                        FOR UPDATE;
                        ",
                        &[&public_id],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;

                let Some(row) = row else {
                    return Err(TransitionError::NotFound);
                };

                let event_id = row.get::<_, i64>("id");
                let current_status = row.get::<_, String>("status");
                if !allowed_statuses.iter().any(|status| current_status == *status) {
                    return Err(TransitionError::InvalidTransition);
                }

                let changed_at = Utc::now();
                postgres
                    .client
                    .execute(
                        "
                        UPDATE webhook_hub.webhook_events
                        SET status = $2
                        WHERE id = $1;
                        ",
                        &[&event_id, &next_status],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;
                postgres
                    .client
                    .execute(
                        "
                        INSERT INTO webhook_hub.webhook_event_transitions (webhook_event_id, status, changed_at)
                        VALUES ($1, $2, $3);
                        ",
                        &[&event_id, &next_status, &changed_at],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;

                self.find_event_by_public_id(&public_id.to_string())
                    .await
                    .map_err(|_| TransitionError::Storage)?
                    .ok_or(TransitionError::NotFound)
            }
        }
    }

    async fn move_to_dead_letter(
        &self,
        public_id: &str,
        error_code: Option<String>,
        error_message: Option<String>,
        allowed_statuses: &[&str],
    ) -> Result<WebhookEvent, TransitionError> {
        match &self.storage {
            StorageBackend::Memory(memory) => {
                let mut events = memory.events.write().await;
                let webhook_event = events
                    .iter_mut()
                    .find(|event| event.public_id == public_id)
                    .ok_or(TransitionError::NotFound)?;

                if !allowed_statuses.iter().any(|status| webhook_event.status == *status) {
                    return Err(TransitionError::InvalidTransition);
                }

                webhook_event.status = "dead_letter".to_string();
                webhook_event.retry_count += 1;
                webhook_event.last_error_code = error_code
                    .map(|value| value.trim().to_lowercase())
                    .filter(|value| !value.is_empty());
                webhook_event.last_error_message = error_message
                    .map(|value| value.trim().to_string())
                    .filter(|value| !value.is_empty());
                webhook_event.dead_lettered_at = Some(Utc::now().to_rfc3339());
                webhook_event.status_history.push(WebhookEventStatusTransition::new(
                    "dead_letter",
                    Utc::now().to_rfc3339(),
                ));
                Ok(webhook_event.clone())
            }
            StorageBackend::Postgres(postgres) => {
                let Ok(public_id) = Uuid::parse_str(public_id) else {
                    return Err(TransitionError::NotFound);
                };

                let row = postgres
                    .client
                    .query_opt(
                        "
                        SELECT id, status, retry_count
                        FROM webhook_hub.webhook_events
                        WHERE public_id = $1
                        FOR UPDATE;
                        ",
                        &[&public_id],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;

                let Some(row) = row else {
                    return Err(TransitionError::NotFound);
                };

                let event_id = row.get::<_, i64>("id");
                let current_status = row.get::<_, String>("status");
                if !allowed_statuses.iter().any(|status| current_status == *status) {
                    return Err(TransitionError::InvalidTransition);
                }

                let retry_count = row.get::<_, i32>("retry_count") + 1;
                let normalized_error_code = error_code
                    .map(|value| value.trim().to_lowercase())
                    .filter(|value| !value.is_empty());
                let normalized_error_message = error_message
                    .map(|value| value.trim().to_string())
                    .filter(|value| !value.is_empty());
                let dead_lettered_at = Utc::now();
                postgres
                    .client
                    .execute(
                        "
                        UPDATE webhook_hub.webhook_events
                        SET status = 'dead_letter',
                            retry_count = $2,
                            last_error_code = $3,
                            last_error_message = $4,
                            dead_lettered_at = $5
                        WHERE id = $1;
                        ",
                        &[&event_id, &retry_count, &normalized_error_code, &normalized_error_message, &dead_lettered_at],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;
                postgres
                    .client
                    .execute(
                        "
                        INSERT INTO webhook_hub.webhook_event_transitions (webhook_event_id, status, changed_at)
                        VALUES ($1, 'dead_letter', $2);
                        ",
                        &[&event_id, &dead_lettered_at],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;

                self.find_event_by_public_id(&public_id.to_string())
                    .await
                    .map_err(|_| TransitionError::Storage)?
                    .ok_or(TransitionError::NotFound)
            }
        }
    }

    async fn requeue_event(
        &self,
        public_id: &str,
        reason: Option<String>,
        allowed_statuses: &[&str],
    ) -> Result<WebhookEvent, TransitionError> {
        match &self.storage {
            StorageBackend::Memory(memory) => {
                let mut events = memory.events.write().await;
                let webhook_event = events
                    .iter_mut()
                    .find(|event| event.public_id == public_id)
                    .ok_or(TransitionError::NotFound)?;

                if !allowed_statuses.iter().any(|status| webhook_event.status == *status) {
                    return Err(TransitionError::InvalidTransition);
                }

                webhook_event.status = "queued".to_string();
                webhook_event.last_error_message = reason
                    .map(|value| value.trim().to_string())
                    .filter(|value| !value.is_empty())
                    .or(Some("Webhook event requeued after operator intervention.".to_string()));
                webhook_event.dead_lettered_at = None;
                webhook_event.status_history.push(WebhookEventStatusTransition::new(
                    "queued",
                    Utc::now().to_rfc3339(),
                ));
                Ok(webhook_event.clone())
            }
            StorageBackend::Postgres(postgres) => {
                let Ok(public_id) = Uuid::parse_str(public_id) else {
                    return Err(TransitionError::NotFound);
                };

                let row = postgres
                    .client
                    .query_opt(
                        "
                        SELECT id, status
                        FROM webhook_hub.webhook_events
                        WHERE public_id = $1
                        FOR UPDATE;
                        ",
                        &[&public_id],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;

                let Some(row) = row else {
                    return Err(TransitionError::NotFound);
                };

                let event_id = row.get::<_, i64>("id");
                let current_status = row.get::<_, String>("status");
                if !allowed_statuses.iter().any(|status| current_status == *status) {
                    return Err(TransitionError::InvalidTransition);
                }

                let changed_at = Utc::now();
                let normalized_reason = reason
                    .map(|value| value.trim().to_string())
                    .filter(|value| !value.is_empty())
                    .or(Some("Webhook event requeued after operator intervention.".to_string()));
                postgres
                    .client
                    .execute(
                        "
                        UPDATE webhook_hub.webhook_events
                        SET status = 'queued',
                            last_error_message = $2,
                            dead_lettered_at = NULL
                        WHERE id = $1;
                        ",
                        &[&event_id, &normalized_reason],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;
                postgres
                    .client
                    .execute(
                        "
                        INSERT INTO webhook_hub.webhook_event_transitions (webhook_event_id, status, changed_at)
                        VALUES ($1, 'queued', $2);
                        ",
                        &[&event_id, &changed_at],
                    )
                    .await
                    .map_err(|_| TransitionError::Storage)?;

                self.find_event_by_public_id(&public_id.to_string())
                    .await
                    .map_err(|_| TransitionError::Storage)?
                    .ok_or(TransitionError::NotFound)
            }
        }
    }

    async fn summary(&self) -> Result<WebhookEventSummary, StorageError> {
        let events = self.list_events_filtered(None, None, None).await?;
        let mut by_provider = BTreeMap::new();
        let mut by_status = BTreeMap::new();

        for event in events.iter() {
            *by_provider.entry(event.provider.clone()).or_insert(0) += 1;
            *by_status.entry(event.status.clone()).or_insert(0) += 1;
        }

        Ok(WebhookEventSummary {
            total: events.len() as u64,
            received: events
                .iter()
                .filter(|event| event.status == "received")
                .count() as u64,
            pending_delivery: events
                .iter()
                .filter(|event| matches!(event.status.as_str(), "validated" | "queued" | "processing"))
                .count() as u64,
            handled: events
                .iter()
                .filter(|event| matches!(event.status.as_str(), "forwarded" | "failed" | "rejected" | "dead_letter"))
                .count() as u64,
            dead_letter: events
                .iter()
                .filter(|event| event.status == "dead_letter")
                .count() as u64,
            by_provider,
            by_status,
        })
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
struct WebhookHubCapabilitiesResponse {
    outbound_webhooks: OutboundWebhookCapability,
}

#[derive(serde::Serialize)]
struct OutboundWebhookCapability {
    signing_ready: bool,
    status: String,
    supports_retries: bool,
    supports_dlq: bool,
    credential_key: String,
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
    retry_count: i32,
    last_error_code: Option<String>,
    last_error_message: Option<String>,
    dead_lettered_at: Option<String>,
    received_at: String,
    status_history: Vec<WebhookEventStatusTransition>,
}

#[derive(Clone, serde::Serialize)]
struct WebhookEventStatusTransition {
    status: String,
    changed_at: String,
}

impl WebhookEventStatusTransition {
    fn new(status: &str, changed_at: String) -> Self {
        Self {
            status: status.to_string(),
            changed_at,
        }
    }
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

#[derive(Default, serde::Deserialize)]
struct ListDeadLetterEventsQuery {
    provider: Option<String>,
}

#[derive(serde::Serialize)]
struct WebhookEventSummary {
    total: u64,
    received: u64,
    pending_delivery: u64,
    handled: u64,
    dead_letter: u64,
    by_provider: BTreeMap<String, u64>,
    by_status: BTreeMap<String, u64>,
}

#[derive(Default, serde::Deserialize)]
struct MoveToDeadLetterRequest {
    error_code: Option<String>,
    error_message: Option<String>,
}

#[derive(Default, serde::Deserialize)]
struct RequeueWebhookEventRequest {
    reason: Option<String>,
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

enum TransitionError {
    NotFound,
    InvalidTransition,
    Storage,
}

#[derive(Debug)]
struct StorageError;

fn internal_storage_error(_error: StorageError) -> (StatusCode, Json<ErrorResponse>) {
    (
        StatusCode::INTERNAL_SERVER_ERROR,
        Json(ErrorResponse::new(
            "webhook_event_storage_error",
            "Webhook event storage operation failed.",
        )),
    )
}

async fn map_row_to_webhook_event(
    client: &tokio_postgres::Client,
    row: Row,
) -> Result<WebhookEvent, StorageError> {
    let event_id = row.get::<_, i64>("id");
    let status_history_rows = client
        .query(
            "
            SELECT status, changed_at
            FROM webhook_hub.webhook_event_transitions
            WHERE webhook_event_id = $1
            ORDER BY changed_at ASC, id ASC;
            ",
            &[&event_id],
        )
        .await
        .map_err(|_| StorageError)?;

    let status_history = status_history_rows
        .into_iter()
        .map(map_row_to_transition)
        .collect();

    Ok(WebhookEvent {
        id: event_id as u64,
        public_id: row.get::<_, Uuid>("public_id").to_string(),
        provider: row.get("provider"),
        event_type: row.get("event_type"),
        external_id: row.get("external_id"),
        payload_summary: row.get("payload_summary"),
        status: row.get("status"),
        retry_count: row.try_get("retry_count").unwrap_or(0),
        last_error_code: row.try_get("last_error_code").ok().flatten(),
        last_error_message: row.try_get("last_error_message").ok().flatten(),
        dead_lettered_at: row
            .try_get::<_, Option<chrono::DateTime<Utc>>>("dead_lettered_at")
            .ok()
            .flatten()
            .map(|value| value.to_rfc3339()),
        received_at: row
            .get::<_, chrono::DateTime<Utc>>("received_at")
            .to_rfc3339(),
        status_history,
    })
}

fn map_row_to_transition(row: Row) -> WebhookEventStatusTransition {
    WebhookEventStatusTransition::new(
        row.get::<_, String>("status").as_str(),
        row.get::<_, chrono::DateTime<Utc>>("changed_at").to_rfc3339(),
    )
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
    async fn capabilities_route_should_return_outbound_webhook_controls() {
        let app = build_router();
        let response = app
            .oneshot(Request::builder().uri("/api/webhook-hub/capabilities").body(axum::body::Body::empty()).unwrap())
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = response.into_body().collect().await.unwrap().to_bytes();
        let payload: Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(payload["outbound_webhooks"]["status"], "unconfigured");
        assert_eq!(payload["outbound_webhooks"]["supports_dlq"], true);
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
            .clone()
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
        assert_eq!(summary_payload["pending_delivery"], 0);
        assert_eq!(summary_payload["handled"], 0);
        assert_eq!(summary_payload["by_provider"]["tiny"], 1);
        assert_eq!(summary_payload["by_status"]["received"], 1);
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
        assert_eq!(detail_payload["status_history"][0]["status"], "received");
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

    #[tokio::test]
    async fn webhook_event_should_validate_and_reject_through_explicit_routes() {
        let app = build_router();

        let create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"stripe","event_type":"payment.pending","external_id":"evt_transition_001"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(create_response.status(), StatusCode::CREATED);

        let create_body = create_response.into_body().collect().await.unwrap().to_bytes();
        let create_payload: Value = serde_json::from_slice(&create_body).unwrap();
        let public_id = create_payload["public_id"].as_str().unwrap();

        let validate_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/validate"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(validate_response.status(), StatusCode::OK);

        let validate_body = validate_response.into_body().collect().await.unwrap().to_bytes();
        let validate_payload: Value = serde_json::from_slice(&validate_body).unwrap();
        assert_eq!(validate_payload["status"], "validated");

        let reject_response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/reject"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(reject_response.status(), StatusCode::OK);

        let reject_body = reject_response.into_body().collect().await.unwrap().to_bytes();
        let reject_payload: Value = serde_json::from_slice(&reject_body).unwrap();
        assert_eq!(reject_payload["status"], "rejected");
    }

    #[tokio::test]
    async fn webhook_event_should_block_invalid_status_transitions() {
        let app = build_router();

        let create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"tiny","event_type":"lead.created","external_id":"evt_transition_002"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(create_response.status(), StatusCode::CREATED);

        let create_body = create_response.into_body().collect().await.unwrap().to_bytes();
        let create_payload: Value = serde_json::from_slice(&create_body).unwrap();
        let public_id = create_payload["public_id"].as_str().unwrap();

        let first_validate_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/validate"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(first_validate_response.status(), StatusCode::OK);

        let second_validate_response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/validate"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(second_validate_response.status(), StatusCode::CONFLICT);

        let second_validate_body = second_validate_response.into_body().collect().await.unwrap().to_bytes();
        let second_validate_payload: Value = serde_json::from_slice(&second_validate_body).unwrap();
        assert_eq!(
            second_validate_payload["code"],
            "webhook_event_transition_invalid"
        );
    }

    #[tokio::test]
    async fn webhook_event_should_follow_queue_process_and_forward_path() {
        let app = build_router();

        let create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"meta","event_type":"lead.generated","external_id":"evt_transition_003"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(create_response.status(), StatusCode::CREATED);

        let create_body = create_response.into_body().collect().await.unwrap().to_bytes();
        let create_payload: Value = serde_json::from_slice(&create_body).unwrap();
        let public_id = create_payload["public_id"].as_str().unwrap();

        let validate_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/validate"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(validate_response.status(), StatusCode::OK);

        let queue_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/queue"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(queue_response.status(), StatusCode::OK);

        let process_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/process"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(process_response.status(), StatusCode::OK);

        let forward_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/forward"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(forward_response.status(), StatusCode::OK);

        let forward_body = forward_response.into_body().collect().await.unwrap().to_bytes();
        let forward_payload: Value = serde_json::from_slice(&forward_body).unwrap();
        assert_eq!(forward_payload["status"], "forwarded");

        let summary_response = app
            .clone()
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
        assert_eq!(summary_payload["handled"], 1);
        assert_eq!(summary_payload["by_status"]["forwarded"], 1);

        let transitions_response = app
            .oneshot(
                Request::builder()
                    .uri(format!("/api/webhook-hub/events/{public_id}/transitions"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(transitions_response.status(), StatusCode::OK);

        let transitions_body = transitions_response.into_body().collect().await.unwrap().to_bytes();
        let transitions_payload: Value = serde_json::from_slice(&transitions_body).unwrap();
        assert_eq!(transitions_payload.as_array().unwrap().len(), 5);
        assert_eq!(transitions_payload[0]["status"], "received");
        assert_eq!(transitions_payload[1]["status"], "validated");
        assert_eq!(transitions_payload[2]["status"], "queued");
        assert_eq!(transitions_payload[3]["status"], "processing");
        assert_eq!(transitions_payload[4]["status"], "forwarded");
    }

    #[tokio::test]
    async fn webhook_event_should_fail_after_being_queued() {
        let app = build_router();

        let create_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/api/webhook-hub/events")
                    .header("content-type", "application/json")
                    .body(axum::body::Body::from(
                        r#"{"provider":"hotmart","event_type":"invoice.failed","external_id":"evt_transition_004"}"#,
                    ))
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(create_response.status(), StatusCode::CREATED);

        let create_body = create_response.into_body().collect().await.unwrap().to_bytes();
        let create_payload: Value = serde_json::from_slice(&create_body).unwrap();
        let public_id = create_payload["public_id"].as_str().unwrap();

        let validate_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/validate"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(validate_response.status(), StatusCode::OK);

        let queue_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/queue"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(queue_response.status(), StatusCode::OK);

        let fail_response = app
            .clone()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri(format!("/api/webhook-hub/events/{public_id}/fail"))
                    .body(axum::body::Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();
        assert_eq!(fail_response.status(), StatusCode::OK);

        let fail_body = fail_response.into_body().collect().await.unwrap().to_bytes();
        let fail_payload: Value = serde_json::from_slice(&fail_body).unwrap();
        assert_eq!(fail_payload["status"], "failed");

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
        assert_eq!(summary_payload["handled"], 1);
        assert_eq!(summary_payload["by_status"]["failed"], 1);
    }
}
