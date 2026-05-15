from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect

ENTITIES = [
    ('accounts', 'account', 'accountCode', 'accountName', ['accountType', 'normalBalance', 'parentAccountCode']),
    ('cost-centers', 'cost_center', 'costCenterCode', 'costCenterName', ['parentCostCenterCode', 'managerUserId']),
    ('journal-entries', 'journal_entry', 'entryNumber', 'description', ['entryDate', 'sourceDocument', 'sourceService', 'sourcePublicId', 'totalDebitCents', 'totalCreditCents', 'lines', 'status']),
    ('posting-rules', 'posting_rule', 'ruleKey', 'description', ['sourceService', 'eventType', 'debitAccountCode', 'creditAccountCode', 'costCenterCode', 'active']),
    ('period-closes', 'period_close', 'periodKey', 'status', ['closedBy', 'closedAt', 'notes']),
]
CAPABILITIES = ['chart_of_accounts', 'cost_centers', 'immutable_journal_entries', 'posting_rules', 'period_close', 'management_statements', 'general_ledger']
SCHEMA_NAME = "accounting"
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
    if collection == "journal-entries":
        debit = int(payload.get("totalDebitCents") or 0)
        credit = int(payload.get("totalCreditCents") or 0)
        lines = payload.get("lines") if isinstance(payload.get("lines"), list) else []
        if lines:
            debit = sum(int(line.get("debitCents") or 0) for line in lines)
            credit = sum(int(line.get("creditCents") or 0) for line in lines)
        if debit <= 0 or credit <= 0 or debit != credit:
            raise ValueError("journal_entry_unbalanced")
        payload = {**payload, "totalDebitCents": debit, "totalCreditCents": credit, "status": payload.get("status") or "posted"}
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
    if collection == "journal-entries":
        raise ValueError("journal_entry_immutable")
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


def build_general_ledger(tenant: str | None = None, account_code: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    entries = list_records("journal-entries", slug)
    lines: list[dict] = []
    balance = 0
    for entry in entries:
        entry_lines = entry.get("payload", {}).get("lines") or []
        if not entry_lines:
            entry_lines = [
                {
                    "accountCode": entry.get("payload", {}).get("debitAccountCode", "unclassified"),
                    "debitCents": entry.get("payload", {}).get("totalDebitCents", 0),
                    "creditCents": 0,
                },
                {
                    "accountCode": entry.get("payload", {}).get("creditAccountCode", "unclassified"),
                    "debitCents": 0,
                    "creditCents": entry.get("payload", {}).get("totalCreditCents", 0),
                },
            ]
        for line in entry_lines:
            if account_code and line.get("accountCode") != account_code:
                continue
            debit = int(line.get("debitCents") or 0)
            credit = int(line.get("creditCents") or 0)
            balance += debit - credit
            lines.append(
                {
                    "entryPublicId": entry["publicId"],
                    "entryNumber": entry.get("entryNumber"),
                    "accountCode": line.get("accountCode"),
                    "costCenterCode": line.get("costCenterCode"),
                    "debitCents": debit,
                    "creditCents": credit,
                    "runningBalanceCents": balance,
                }
            )
    return {"tenantSlug": slug, "accountCode": account_code, "balanceCents": balance, "lines": lines}


def build_financial_statement(statement: str, tenant: str | None = None) -> dict:
    slug = tenant_slug(tenant)
    accounts = {item.get("accountCode"): item for item in list_records("accounts", slug)}
    ledger = build_general_ledger(slug)
    buckets: dict[str, int] = {}
    for line in ledger["lines"]:
        account = accounts.get(line.get("accountCode"), {})
        account_type = str(account.get("payload", {}).get("accountType") or "unclassified")
        buckets[account_type] = buckets.get(account_type, 0) + int(line["debitCents"]) - int(line["creditCents"])
    return {
        "tenantSlug": slug,
        "statement": statement,
        "generatedAt": utc_now(),
        "buckets": buckets,
        "dre": {
            "revenueCents": abs(buckets.get("revenue", 0)),
            "expenseCents": abs(buckets.get("expense", 0)),
            "resultCents": abs(buckets.get("revenue", 0)) - abs(buckets.get("expense", 0)),
        },
        "balanceSheet": {
            "assetCents": buckets.get("asset", 0),
            "liabilityCents": abs(buckets.get("liability", 0)),
            "equityCents": abs(buckets.get("equity", 0)),
        },
    }


def post_source_event(payload: dict) -> dict:
    slug = tenant_slug(payload.get("tenantSlug"))
    source_service = str(payload.get("sourceService") or "").strip()
    event_type = str(payload.get("eventType") or "").strip()
    amount = int(payload.get("amountCents") or 0)
    if not source_service or not event_type or amount <= 0:
        raise ValueError("posting_source_invalid")
    rules = [
        rule
        for rule in list_records("posting-rules", slug, "active")
        if rule.get("payload", {}).get("sourceService") == source_service
        and rule.get("payload", {}).get("eventType") == event_type
    ]
    if not rules:
        raise ValueError("posting_rule_not_found")
    rule = rules[0]
    rule_payload = rule.get("payload", {})
    return create_record(
        "journal-entries",
        {
            "tenantSlug": slug,
            "entryNumber": payload.get("entryNumber") or f"{source_service}-{event_type}-{uuid.uuid4()}",
            "description": payload.get("description") or f"{source_service}.{event_type}",
            "sourceService": source_service,
            "sourcePublicId": payload.get("sourcePublicId"),
            "totalDebitCents": amount,
            "totalCreditCents": amount,
            "lines": [
                {
                    "accountCode": rule_payload.get("debitAccountCode"),
                    "costCenterCode": rule_payload.get("costCenterCode"),
                    "debitCents": amount,
                    "creditCents": 0,
                },
                {
                    "accountCode": rule_payload.get("creditAccountCode"),
                    "costCenterCode": rule_payload.get("costCenterCode"),
                    "debitCents": 0,
                    "creditCents": amount,
                },
            ],
            "status": "posted",
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
