from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect

ENTITIES = [
    ('requisitions', 'requisition', 'requisitionNumber', 'description', ['requestedBy', 'costCenter', 'estimatedAmountCents', 'status']),
    ('quotations', 'quotation', 'quotationNumber', 'description', ['supplierPublicId', 'amountCents', 'validUntil', 'status']),
    ('purchase-orders', 'purchase_order', 'orderNumber', 'description', ['supplierPublicId', 'amountCents', 'status', 'expectedDeliveryDate']),
    ('approvals', 'approval', 'approvalNumber', 'status', ['targetCollection', 'targetPublicId', 'approvedBy', 'rejectedBy', 'reason']),
    ('receipts', 'receipt', 'receiptNumber', 'description', ['purchaseOrderPublicId', 'receivedBy', 'amountCents', 'status']),
    ('three-way-matches', 'three_way_match', 'matchNumber', 'status', ['purchaseOrderPublicId', 'receiptPublicId', 'fiscalDocumentPublicId', 'orderAmountCents', 'receiptAmountCents', 'invoiceAmountCents', 'divergenceCents']),
]
CAPABILITIES = ['purchase_requisitions', 'quotations', 'purchase_orders', 'approvals', 'receipts', 'three_way_match', 'payables_cost_integration']
SCHEMA_NAME = "procurement"
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


def approve_target(payload: dict) -> dict:
    target_collection = str(payload.get("targetCollection") or "").strip()
    target_public_id = str(payload.get("targetPublicId") or "").strip()
    status = str(payload.get("status") or "approved").strip().lower()
    if target_collection not in {"requisitions", "quotations", "purchase-orders", "receipts"} or not target_public_id:
        raise ValueError("approval_target_invalid")
    target = get_record(target_collection, target_public_id, payload.get("tenantSlug"))
    if target is None:
        raise ValueError("approval_target_not_found")
    approval = create_record(
        "approvals",
        {
            "tenantSlug": payload.get("tenantSlug"),
            "approvalNumber": payload.get("approvalNumber") or f"APR-{uuid.uuid4()}",
            "status": status,
            "targetCollection": target_collection,
            "targetPublicId": target_public_id,
            "approvedBy": payload.get("approvedBy"),
            "rejectedBy": payload.get("rejectedBy"),
            "reason": payload.get("reason"),
        },
    )
    transition_record(target_collection, target_public_id, {"tenantSlug": payload.get("tenantSlug"), "status": status})
    return {"approval": approval, "target": get_record(target_collection, target_public_id, payload.get("tenantSlug"))}


def run_three_way_match(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    order_id = str(payload.get("purchaseOrderPublicId") or "").strip()
    receipt_id = str(payload.get("receiptPublicId") or "").strip()
    fiscal_document_id = str(payload.get("fiscalDocumentPublicId") or "").strip()
    order = get_record("purchase-orders", order_id, slug) if order_id else None
    receipt = get_record("receipts", receipt_id, slug) if receipt_id else None
    if order is None or receipt is None or not fiscal_document_id:
        raise ValueError("three_way_match_target_invalid")
    order_amount = int(payload.get("orderAmountCents") or order.get("payload", {}).get("amountCents") or 0)
    receipt_amount = int(payload.get("receiptAmountCents") or receipt.get("payload", {}).get("amountCents") or order_amount)
    invoice_amount = int(payload.get("invoiceAmountCents") or receipt_amount)
    divergence = max(order_amount, receipt_amount, invoice_amount) - min(order_amount, receipt_amount, invoice_amount)
    tolerance = int(payload.get("toleranceCents") or 0)
    status = "matched" if divergence <= tolerance else "divergent"
    return create_record(
        "three-way-matches",
        {
            "tenantSlug": slug,
            "matchNumber": payload.get("matchNumber") or f"MATCH-{uuid.uuid4()}",
            "status": status,
            "purchaseOrderPublicId": order_id,
            "receiptPublicId": receipt_id,
            "fiscalDocumentPublicId": fiscal_document_id,
            "orderAmountCents": order_amount,
            "receiptAmountCents": receipt_amount,
            "invoiceAmountCents": invoice_amount,
            "divergenceCents": divergence,
        },
    )


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
