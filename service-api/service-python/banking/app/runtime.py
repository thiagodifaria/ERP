from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect

ENTITIES = [
    ('bank-accounts', 'bank_account', 'accountCode', 'bankName', ['branchNumber', 'accountNumber', 'pixKey', 'status']),
    ('boletos', 'boleto', 'boletoNumber', 'payerName', ['bankAccountPublicId', 'amountCents', 'dueDate', 'status', 'ourNumber']),
    ('cnab-files', 'cnab_file', 'fileNumber', 'fileType', ['bankAccountPublicId', 'direction', 'status', 'amountCents', 'layoutVersion', 'rawContent']),
    ('bank-statements', 'bank_statement', 'statementNumber', 'bankAccountPublicId', ['statementDate', 'amountCents', 'transactionCount', 'rawContent']),
    ('pix-charges', 'pix_charge', 'txid', 'description', ['bankAccountPublicId', 'amountCents', 'status', 'expiresAt']),
    ('pix-refunds', 'pix_refund', 'refundId', 'txid', ['amountCents', 'reason', 'status']),
    ('pix-webhooks', 'pix_webhook', 'eventId', 'txid', ['eventType', 'amountCents', 'status', 'payload']),
    ('open-finance-connections', 'open_finance_connection', 'connectionId', 'bankName', ['consentId', 'status', 'expiresAt']),
    ('reconciliations', 'reconciliation', 'reconciliationNumber', 'status', ['bankAccountPublicId', 'statementDate', 'matchedAmountCents', 'divergenceAmountCents', 'statementPublicId']),
]
CAPABILITIES = ['cnab_240', 'boleto_lifecycle', 'pix_cob', 'pix_refund', 'pix_webhooks', 'bank_statements', 'bank_reconciliation', 'open_finance']
SCHEMA_NAME = "banking"
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


def parse_cnab_return(payload: dict) -> dict:
    raw_content = str(payload.get("rawContent") or "")
    if not raw_content:
        raise ValueError("cnab_content_required")
    lines = [line for line in raw_content.splitlines() if line.strip()]
    amount = sum(int("".join(ch for ch in line[-15:] if ch.isdigit()) or 0) for line in lines)
    return create_record(
        "cnab-files",
        {
            "tenantSlug": payload.get("tenantSlug"),
            "fileNumber": payload.get("fileNumber") or f"CNAB-{uuid.uuid4()}",
            "fileType": payload.get("fileType") or "return",
            "direction": "return",
            "status": "processed",
            "amountCents": amount,
            "layoutVersion": payload.get("layoutVersion") or "240",
            "rawContent": raw_content,
        },
    )


def reconcile_statement(payload: dict) -> dict:
    statement_id = str(payload.get("statementPublicId") or "").strip()
    slug = tenant_slug(payload.get("tenantSlug"))
    statement = get_record("bank-statements", statement_id, slug) if statement_id else None
    statement_amount = int(payload.get("statementAmountCents") or (statement or {}).get("payload", {}).get("amountCents") or 0)
    expected_amount = int(payload.get("expectedAmountCents") or 0)
    divergence = statement_amount - expected_amount
    status = "matched" if divergence == 0 else "divergent"
    return create_record(
        "reconciliations",
        {
            "tenantSlug": slug,
            "reconciliationNumber": payload.get("reconciliationNumber") or f"REC-{uuid.uuid4()}",
            "status": status,
            "bankAccountPublicId": payload.get("bankAccountPublicId") or (statement or {}).get("bankAccountPublicId"),
            "statementPublicId": statement_id,
            "statementDate": payload.get("statementDate"),
            "matchedAmountCents": min(statement_amount, expected_amount),
            "divergenceAmountCents": abs(divergence),
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
