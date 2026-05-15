from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect

ENTITIES = [
    ('locations', 'location', 'locationCode', 'name', ['warehouseCode', 'kind', 'active']),
    ('movements', 'movement', 'movementNumber', 'reason', ['sku', 'locationCode', 'quantity', 'movementType', 'unitCostCents', 'sourceService', 'sourcePublicId']),
    ('reservations', 'reservation', 'reservationNumber', 'status', ['sku', 'locationCode', 'quantity', 'expiresAt', 'sourceService', 'sourcePublicId']),
    ('cycle-counts', 'cycle_count', 'countNumber', 'status', ['sku', 'locationCode', 'countedQuantity', 'expectedQuantity', 'countedBy']),
    ('cost-layers', 'cost_layer', 'layerNumber', 'sku', ['locationCode', 'quantity', 'remainingQuantity', 'unitCostCents', 'costMethod']),
]
CAPABILITIES = ['warehouse_locations', 'stock_balances', 'stock_movements', 'reservations', 'cycle_counts', 'average_cost', 'fifo_cost_layers']
SCHEMA_NAME = "inventory"
IN_MEMORY_STATE: dict[str, list[dict]] = {entity[0]: [] for entity in ENTITIES}


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def find_tenant_id(cursor, slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def entity_config(collection: str):
    for item in ENTITIES:
        if item[0] == collection:
            return item
    raise ValueError("entity_not_found")


def capability_catalog() -> dict:
    return {"service": settings.service_name, "capabilities": [{"key": key, "status": "ready"} for key in CAPABILITIES]}


def list_records(collection: str, tenant: str | None = None, status: str | None = None) -> list[dict]:
    slug = tenant_slug(tenant)
    entity_config(collection)
    normalized_status = (status or "").strip().lower()
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE[collection] if item["tenantSlug"] == slug]
        if normalized_status:
            records = [item for item in records if item.get("status") == normalized_status]
        return records
    with connect() as connection:
        with connection.cursor() as cursor:
            clauses = ["tenant.slug = %s"]
            params: list[object] = [slug]
            if normalized_status:
                clauses.append("record.status = %s")
                params.append(normalized_status)
            query = f"""
                SELECT record.public_id, record.record_key, record.name, record.status, record.payload_json, record.created_at, record.updated_at
                FROM {SCHEMA_NAME}.records AS record
                JOIN identity.tenants AS tenant ON tenant.id = record.tenant_id
                WHERE record.collection = %s AND {' AND '.join(clauses)}
                ORDER BY record.created_at DESC
            """
            cursor.execute(query, [collection, *params])
            return [map_row(row, slug, collection) for row in cursor.fetchall()]


def create_record(collection: str, payload: dict) -> dict:
    _, singular, key_field, name_field, extra_fields = entity_config(collection)
    slug = tenant_slug(payload.get("tenantSlug"))
    key = str(payload.get(key_field) or payload.get("key") or uuid.uuid4()).strip()
    label = str(payload.get(name_field) or payload.get("name") or key).strip()
    if not key:
        raise ValueError(f"{singular}_key_required")
    if not label:
        raise ValueError(f"{singular}_name_required")
    if collection in {"movements", "reservations", "cycle-counts", "cost-layers"}:
        quantity = int(payload.get("quantity") or payload.get("countedQuantity") or payload.get("remainingQuantity") or 0)
        if quantity <= 0:
            raise ValueError(f"{singular}_quantity_required")
    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "collection": collection,
        key_field: key,
        name_field: label,
        "status": str(payload.get("status") or "draft").strip().lower(),
        "payload": {field: payload.get(field) for field in extra_fields if field in payload},
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }
    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE[collection].append(record)
        return record
    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = find_tenant_id(cursor, slug)
            query = f"""
                INSERT INTO {SCHEMA_NAME}.records (tenant_id, public_id, collection, record_key, name, status, payload_json)
                VALUES (%s, %s, %s, %s, %s, %s, %s::jsonb)
                RETURNING created_at, updated_at
            """
            cursor.execute(query, (tenant_id, record["publicId"], collection, key, label, record["status"], json.dumps(record["payload"])))
            row = cursor.fetchone()
            connection.commit()
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def get_record(collection: str, public_id: str, tenant: str | None = None) -> dict | None:
    slug = tenant_slug(tenant)
    entity_config(collection)
    if settings.repository_driver != "postgres":
        return next((item for item in IN_MEMORY_STATE[collection] if item["tenantSlug"] == slug and item["publicId"] == public_id), None)
    with connect() as connection:
        with connection.cursor() as cursor:
            query = f"""
                SELECT record.public_id, record.record_key, record.name, record.status, record.payload_json, record.created_at, record.updated_at
                FROM {SCHEMA_NAME}.records AS record
                JOIN identity.tenants AS tenant ON tenant.id = record.tenant_id
                WHERE record.collection = %s AND record.public_id = %s::uuid AND tenant.slug = %s
            """
            cursor.execute(query, (collection, public_id, slug))
            row = cursor.fetchone()
            return map_row(row, slug, collection) if row else None


def transition_record(collection: str, public_id: str, payload: dict) -> dict | None:
    slug = tenant_slug(payload.get("tenantSlug"))
    status = str(payload.get("status") or "").strip().lower()
    if not status:
        raise ValueError("status_required")
    if settings.repository_driver != "postgres":
        record = get_record(collection, public_id, slug)
        if record is None:
            return None
        record["status"] = status
        record["updatedAt"] = utc_now()
        return record
    with connect() as connection:
        with connection.cursor() as cursor:
            query = f"""
                UPDATE {SCHEMA_NAME}.records AS record
                SET status = %s, updated_at = NOW()
                FROM identity.tenants AS tenant
                WHERE tenant.id = record.tenant_id
                  AND record.collection = %s
                  AND record.public_id = %s::uuid
                  AND tenant.slug = %s
                RETURNING record.public_id, record.record_key, record.name, record.status, record.payload_json, record.created_at, record.updated_at
            """
            cursor.execute(query, (status, collection, public_id, slug))
            row = cursor.fetchone()
            connection.commit()
            return map_row(row, slug, collection) if row else None


def build_balances(tenant: str | None = None, sku: str | None = None, location_code: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    balances: dict[tuple[str, str], dict] = {}
    for movement in list_records("movements", slug):
        payload = movement.get("payload", {})
        item_sku = str(payload.get("sku") or "")
        item_location = str(payload.get("locationCode") or "")
        if sku and item_sku != sku:
            continue
        if location_code and item_location != location_code:
            continue
        quantity = int(payload.get("quantity") or 0)
        movement_type = str(payload.get("movementType") or "").lower()
        signed_quantity = -quantity if movement_type in {"out", "issue", "sale", "adjustment_negative"} else quantity
        key = (item_sku, item_location)
        current = balances.setdefault(key, {"sku": item_sku, "locationCode": item_location, "onHand": 0, "reserved": 0, "available": 0, "valuationCents": 0})
        current["onHand"] += signed_quantity
        current["valuationCents"] += signed_quantity * int(payload.get("unitCostCents") or 0)
    for reservation in list_records("reservations", slug):
        if reservation.get("status") not in {"active", "reserved", "approved"}:
            continue
        payload = reservation.get("payload", {})
        item_sku = str(payload.get("sku") or "")
        item_location = str(payload.get("locationCode") or "")
        if sku and item_sku != sku:
            continue
        if location_code and item_location != location_code:
            continue
        key = (item_sku, item_location)
        current = balances.setdefault(key, {"sku": item_sku, "locationCode": item_location, "onHand": 0, "reserved": 0, "available": 0, "valuationCents": 0})
        current["reserved"] += int(payload.get("quantity") or 0)
    for item in balances.values():
        item["available"] = item["onHand"] - item["reserved"]
        item["averageUnitCostCents"] = int(item["valuationCents"] / item["onHand"]) if item["onHand"] else 0
    return {"tenantSlug": slug, "generatedAt": utc_now(), "balances": list(balances.values())}


def build_costing_summary(tenant: str | None = None, sku: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    layers = [
        item
        for item in list_records("cost-layers", slug)
        if not sku or item.get("layerNumber") == sku or item.get("sku") == sku or item.get("payload", {}).get("sku") == sku
    ]
    total_quantity = sum(int(layer.get("payload", {}).get("remainingQuantity") or layer.get("payload", {}).get("quantity") or 0) for layer in layers)
    total_value = sum(
        int(layer.get("payload", {}).get("remainingQuantity") or layer.get("payload", {}).get("quantity") or 0)
        * int(layer.get("payload", {}).get("unitCostCents") or 0)
        for layer in layers
    )
    return {
        "tenantSlug": slug,
        "sku": sku,
        "method": "fifo+average",
        "layers": layers,
        "summary": {
            "remainingQuantity": total_quantity,
            "valuationCents": total_value,
            "averageUnitCostCents": int(total_value / total_quantity) if total_quantity else 0,
        },
    }


def build_cycle_count_variances(tenant: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    variances = []
    for count in list_records("cycle-counts", slug):
        payload = count.get("payload", {})
        counted = int(payload.get("countedQuantity") or 0)
        expected = int(payload.get("expectedQuantity") or 0)
        variances.append({**count, "varianceQuantity": counted - expected, "requiresAdjustment": counted != expected})
    return {"tenantSlug": slug, "generatedAt": utc_now(), "variances": variances}


def build_summary(tenant: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    summary = {}
    for collection, *_rest in ENTITIES:
        records = list_records(collection, slug)
        summary[collection.replace('-', '_')] = {"total": len(records), "active": sum(1 for item in records if item.get("status") in {"active", "posted", "approved", "matched", "closed"})}
    return {"tenantSlug": slug, "generatedAt": utc_now(), "summary": summary, "capabilities": CAPABILITIES}


def map_row(row: dict, slug: str, collection: str) -> dict:
    _, _singular, key_field, name_field, _extra_fields = entity_config(collection)
    payload = row["payload_json"] or {}
    return {
        "publicId": str(row["public_id"]),
        "tenantSlug": slug,
        "collection": collection,
        key_field: row["record_key"],
        name_field: row["name"],
        "status": row["status"],
        "payload": payload,
        "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
    }
