"""Relatorio de confiabilidade e carga do fluxo de entrega de webhooks."""

from datetime import datetime, timezone

from app.config.settings import settings
from app.infrastructure.postgres import connect


def build_delivery_reliability(provider: str | None = None) -> dict:
    if settings.repository_driver == "postgres":
        return build_postgres_delivery_reliability(provider)

    return build_static_delivery_reliability(provider)


def build_static_delivery_reliability(provider: str | None = None) -> dict:
    return {
        "provider": provider,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "lifecycle": {
            "totalEvents": 93,
            "handledEvents": 90,
            "failedEvents": 3,
            "pendingEvents": 0,
            "avgTransitionsPerEvent": 4.7,
        },
        "statusFootprint": {
            "received": 93,
            "validated": 92,
            "queued": 90,
            "processing": 89,
            "forwarded": 87,
            "failed": 3,
            "rejected": 2,
            "dead_letter": 1,
        },
        "transitionLoad": {
            "received": 93,
            "validated": 92,
            "queued": 90,
            "processing": 89,
            "forwarded": 87,
            "failed": 3,
            "rejected": 2,
            "dead_letter": 1,
        },
        "providerLeaderboard": [
            {"provider": "stripe", "total": 61, "forwarded": 58, "failed": 2},
            {"provider": "meta", "total": 19, "forwarded": 18, "failed": 1},
            {"provider": "asaas", "total": 13, "forwarded": 11, "failed": 0},
        ],
    }


def build_postgres_delivery_reliability(provider: str | None = None) -> dict:
    provider_filter_sql, params = provider_filter(provider)

    with connect() as connection:
        lifecycle_metrics = fetch_lifecycle_metrics(connection, provider_filter_sql, params)
        status_footprint = fetch_status_footprint(connection, provider_filter_sql, params)
        transition_load = fetch_transition_load(connection, provider_filter_sql, params)
        provider_leaderboard = fetch_provider_leaderboard(connection, provider_filter_sql, params)

    return {
        "provider": provider.lower().strip() if provider else None,
        "generatedAt": datetime.now(timezone.utc).isoformat(),
        "dataSource": "postgresql",
        "lifecycle": lifecycle_metrics,
        "statusFootprint": status_footprint,
        "transitionLoad": transition_load,
        "providerLeaderboard": provider_leaderboard,
    }


def fetch_lifecycle_metrics(connection, provider_filter_sql: str, params: list[str]) -> dict:
    query = f"""
        SELECT
            count(*) AS total_events,
            count(*) FILTER (WHERE status IN ('forwarded', 'rejected', 'dead_letter')) AS handled_events,
            count(*) FILTER (WHERE status = 'failed') AS failed_events,
            count(*) FILTER (WHERE status IN ('received', 'validated', 'queued', 'processing')) AS pending_events,
            COALESCE(avg(transition_totals.total_transitions), 0) AS avg_transitions_per_event
        FROM webhook_hub.webhook_events AS event
        LEFT JOIN (
            SELECT webhook_event_id, count(*) AS total_transitions
            FROM webhook_hub.webhook_event_transitions
            GROUP BY webhook_event_id
        ) AS transition_totals ON transition_totals.webhook_event_id = event.id
        {provider_filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {
        "totalEvents": int(row.get("total_events", 0) or 0),
        "handledEvents": int(row.get("handled_events", 0) or 0),
        "failedEvents": int(row.get("failed_events", 0) or 0),
        "pendingEvents": int(row.get("pending_events", 0) or 0),
        "avgTransitionsPerEvent": round(float(row.get("avg_transitions_per_event", 0) or 0), 2),
    }


def fetch_status_footprint(connection, provider_filter_sql: str, params: list[str]) -> dict:
    query = f"""
        SELECT
            count(*) FILTER (WHERE status = 'received') AS received,
            count(*) FILTER (WHERE status = 'validated') AS validated,
            count(*) FILTER (WHERE status = 'queued') AS queued,
            count(*) FILTER (WHERE status = 'processing') AS processing,
            count(*) FILTER (WHERE status = 'forwarded') AS forwarded,
            count(*) FILTER (WHERE status = 'failed') AS failed,
            count(*) FILTER (WHERE status = 'rejected') AS rejected,
            count(*) FILTER (WHERE status = 'dead_letter') AS dead_letter
        FROM webhook_hub.webhook_events
        {provider_filter_sql}
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        row = cursor.fetchone() or {}

    return {status: int(row.get(status, 0) or 0) for status in ["received", "validated", "queued", "processing", "forwarded", "failed", "rejected", "dead_letter"]}


def fetch_transition_load(connection, provider_filter_sql: str, params: list[str]) -> dict:
    query = f"""
        SELECT
            transition.status,
            count(*) AS total
        FROM webhook_hub.webhook_event_transitions AS transition
        JOIN webhook_hub.webhook_events AS event ON event.id = transition.webhook_event_id
        {provider_filter_sql}
        GROUP BY transition.status
        ORDER BY transition.status ASC
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall()

    metrics = {status: 0 for status in ["received", "validated", "queued", "processing", "forwarded", "failed", "rejected", "dead_letter"]}
    for row in rows:
        metrics[row["status"]] = int(row["total"] or 0)

    return metrics


def fetch_provider_leaderboard(connection, provider_filter_sql: str, params: list[str]) -> list[dict]:
    query = f"""
        SELECT
            provider,
            count(*) AS total,
            count(*) FILTER (WHERE status = 'forwarded') AS forwarded,
            count(*) FILTER (WHERE status = 'failed') AS failed
        FROM webhook_hub.webhook_events
        {provider_filter_sql}
        GROUP BY provider
        ORDER BY total DESC, provider ASC
        LIMIT 5
    """

    with connection.cursor() as cursor:
        cursor.execute(query, params)
        rows = cursor.fetchall()

    return [
        {
            "provider": row["provider"],
            "total": int(row["total"] or 0),
            "forwarded": int(row["forwarded"] or 0),
            "failed": int(row["failed"] or 0),
        }
        for row in rows
    ]


def provider_filter(provider: str | None) -> tuple[str, list[str]]:
    normalized = (provider or "").strip().lower()
    if normalized:
        return "WHERE provider = %s", [normalized]
    return "", []
