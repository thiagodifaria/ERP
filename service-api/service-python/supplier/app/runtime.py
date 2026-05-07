from __future__ import annotations

from datetime import datetime, timezone
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "categories": [],
    "suppliers": [],
}


def utc_now() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _tenant_slug(value: str | None) -> str:
    normalized = (value or settings.bootstrap_tenant_slug).strip()
    if normalized == "":
        raise ValueError("tenant_slug_required")
    return normalized


def _find_tenant_id(cursor, tenant_slug: str) -> int:
    cursor.execute("SELECT id FROM identity.tenants WHERE slug = %s", (tenant_slug,))
    row = cursor.fetchone()
    if row is None:
        raise ValueError("tenant_not_found")
    return int(row["id"])


def list_categories(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        return sorted([item for item in IN_MEMORY_STATE["categories"] if item["tenantSlug"] == slug], key=lambda item: item["name"])

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT category.public_id, category.category_key, category.name, category.active, category.created_at, category.updated_at
                FROM supplier.categories AS category
                JOIN identity.tenants AS tenant ON tenant.id = category.tenant_id
                WHERE tenant.slug = %s
                ORDER BY category.name
                """,
                (slug,),
            )
            return [
                {
                    "publicId": row["public_id"],
                    "tenantSlug": slug,
                    "categoryKey": row["category_key"],
                    "name": row["name"],
                    "active": row["active"],
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def upsert_category(category_key: str, payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    normalized_key = category_key.strip().lower()
    name = (payload.get("name") or "").strip()
    if normalized_key == "":
        raise ValueError("supplier_category_key_required")
    if name == "":
        raise ValueError("supplier_category_name_required")

    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "categoryKey": normalized_key,
        "name": name,
        "active": bool(payload.get("active", True)),
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        for index, item in enumerate(IN_MEMORY_STATE["categories"]):
            if item["tenantSlug"] == slug and item["categoryKey"] == normalized_key:
                record["publicId"] = item["publicId"]
                record["createdAt"] = item["createdAt"]
                IN_MEMORY_STATE["categories"][index] = record
                return record
        IN_MEMORY_STATE["categories"].append(record)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO supplier.categories (tenant_id, public_id, category_key, name, active)
                VALUES (%s, %s, %s, %s, %s)
                ON CONFLICT (tenant_id, category_key)
                DO UPDATE SET
                  name = EXCLUDED.name,
                  active = EXCLUDED.active
                RETURNING public_id, created_at, updated_at
                """,
                (tenant_id, record["publicId"], normalized_key, name, record["active"]),
            )
            row = cursor.fetchone()
            connection.commit()
            record["publicId"] = row["public_id"]
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return record


def list_suppliers(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    normalized_status = (status or "").strip().lower()
    if settings.repository_driver != "postgres":
        records = [item for item in IN_MEMORY_STATE["suppliers"] if item["tenantSlug"] == slug]
        if normalized_status:
            records = [item for item in records if item["status"] == normalized_status]
        return sorted(records, key=lambda item: item["companyName"])

    clauses = ["tenant.slug = %s"]
    params: list[object] = [slug]
    if normalized_status:
        clauses.append("supplier.status = %s")
        params.append(normalized_status)

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                f"""
                SELECT supplier.public_id, supplier.company_name, supplier.trade_name, supplier.tax_id, supplier.status,
                       supplier.payable_term_days, supplier.bank_name, supplier.pix_key, supplier.contact_email,
                       supplier.created_at, supplier.updated_at, category.category_key, category.name AS category_name
                FROM supplier.suppliers AS supplier
                JOIN identity.tenants AS tenant ON tenant.id = supplier.tenant_id
                LEFT JOIN supplier.categories AS category ON category.id = supplier.category_id
                WHERE {" AND ".join(clauses)}
                ORDER BY supplier.company_name
                """,
                params,
            )
            return [_map_supplier_row(row, slug) for row in cursor.fetchall()]


def _map_supplier_row(row: dict, tenant_slug: str) -> dict:
    return {
        "publicId": row["public_id"],
        "tenantSlug": tenant_slug,
        "companyName": row["company_name"],
        "tradeName": row["trade_name"],
        "taxId": row["tax_id"],
        "status": row["status"],
        "payableTermDays": int(row["payable_term_days"]),
        "bankName": row["bank_name"],
        "pixKey": row["pix_key"],
        "contactEmail": row["contact_email"],
        "category": None
        if row["category_key"] is None
        else {"categoryKey": row["category_key"], "name": row["category_name"]},
        "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
    }


def create_supplier(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    company_name = (payload.get("companyName") or "").strip()
    tax_id = (payload.get("taxId") or "").strip()
    category_key = (payload.get("categoryKey") or "").strip().lower()
    if company_name == "":
        raise ValueError("supplier_company_name_required")
    if tax_id == "":
        raise ValueError("supplier_tax_id_required")

    record = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "companyName": company_name,
        "tradeName": (payload.get("tradeName") or company_name).strip(),
        "taxId": tax_id,
        "status": (payload.get("status") or "active").strip().lower(),
        "payableTermDays": int(payload.get("payableTermDays") or 30),
        "bankName": (payload.get("bankName") or "").strip(),
        "pixKey": (payload.get("pixKey") or "").strip(),
        "contactEmail": (payload.get("contactEmail") or "").strip(),
        "category": {"categoryKey": category_key} if category_key else None,
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        if category_key and not any(item["tenantSlug"] == slug and item["categoryKey"] == category_key for item in IN_MEMORY_STATE["categories"]):
            raise ValueError("supplier_category_not_found")
        IN_MEMORY_STATE["suppliers"].append(record)
        return record

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            category_id = None
            category_name = None
            if category_key:
                cursor.execute(
                    """
                    SELECT id, name
                    FROM supplier.categories
                    WHERE tenant_id = %s AND category_key = %s
                    """,
                    (tenant_id, category_key),
                )
                category = cursor.fetchone()
                if category is None:
                    raise ValueError("supplier_category_not_found")
                category_id = int(category["id"])
                category_name = category["name"]

            cursor.execute(
                """
                INSERT INTO supplier.suppliers (
                  tenant_id, category_id, public_id, company_name, trade_name, tax_id, status,
                  payable_term_days, bank_name, pix_key, contact_email
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                RETURNING created_at, updated_at
                """,
                (
                    tenant_id,
                    category_id,
                    record["publicId"],
                    record["companyName"],
                    record["tradeName"],
                    record["taxId"],
                    record["status"],
                    record["payableTermDays"],
                    record["bankName"],
                    record["pixKey"],
                    record["contactEmail"],
                ),
            )
            row = cursor.fetchone()
            connection.commit()
            record["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            record["category"] = None if category_key == "" else {"categoryKey": category_key, "name": category_name}
            return record


def bulk_create_suppliers(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    items = payload.get("items") or []
    if not isinstance(items, list):
        raise ValueError("supplier_bulk_items_invalid")

    results: list[dict] = []
    succeeded = 0
    failed = 0
    for index, item in enumerate(items):
        candidate = {"tenantSlug": slug, **item}
        try:
            supplier = create_supplier(candidate)
            succeeded += 1
            results.append({"index": index, "status": "created", "supplier": supplier})
        except ValueError as error:
            failed += 1
            results.append(
                {
                    "index": index,
                    "status": "failed",
                    "errorCode": str(error),
                    "message": "Supplier payload is invalid.",
                }
            )

    return {
        "tenantSlug": slug,
        "summary": {
            "requested": len(items),
            "succeeded": succeeded,
            "failed": failed,
            "partialSuccess": failed > 0,
        },
        "items": results,
    }


def get_supplier(public_id: str, tenant_slug: str | None = None) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        for record in IN_MEMORY_STATE["suppliers"]:
            if record["tenantSlug"] == slug and record["publicId"] == public_id:
                return record
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT supplier.public_id, supplier.company_name, supplier.trade_name, supplier.tax_id, supplier.status,
                       supplier.payable_term_days, supplier.bank_name, supplier.pix_key, supplier.contact_email,
                       supplier.created_at, supplier.updated_at, category.category_key, category.name AS category_name
                FROM supplier.suppliers AS supplier
                JOIN identity.tenants AS tenant ON tenant.id = supplier.tenant_id
                LEFT JOIN supplier.categories AS category ON category.id = supplier.category_id
                WHERE tenant.slug = %s AND supplier.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            return None if row is None else _map_supplier_row(row, slug)


def update_supplier(public_id: str, payload: dict) -> dict | None:
    slug = _tenant_slug(payload.get("tenantSlug"))
    if settings.repository_driver != "postgres":
        for index, record in enumerate(IN_MEMORY_STATE["suppliers"]):
            if record["tenantSlug"] == slug and record["publicId"] == public_id:
                updated = {**record}
                for key in ("companyName", "tradeName", "status", "bankName", "pixKey", "contactEmail"):
                    if key in payload:
                        updated[key] = payload[key]
                if "payableTermDays" in payload:
                    updated["payableTermDays"] = int(payload["payableTermDays"])
                updated["updatedAt"] = utc_now()
                IN_MEMORY_STATE["suppliers"][index] = updated
                return updated
        return None

    current = get_supplier(public_id, slug)
    if current is None:
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE supplier.suppliers AS supplier
                SET company_name = %s,
                    trade_name = %s,
                    status = %s,
                    payable_term_days = %s,
                    bank_name = %s,
                    pix_key = %s,
                    contact_email = %s
                FROM identity.tenants AS tenant
                WHERE tenant.id = supplier.tenant_id
                  AND tenant.slug = %s
                  AND supplier.public_id = %s
                """,
                (
                    payload.get("companyName", current["companyName"]),
                    payload.get("tradeName", current["tradeName"]),
                    payload.get("status", current["status"]),
                    int(payload.get("payableTermDays", current["payableTermDays"])),
                    payload.get("bankName", current["bankName"]),
                    payload.get("pixKey", current["pixKey"]),
                    payload.get("contactEmail", current["contactEmail"]),
                    slug,
                    public_id,
                ),
            )
            connection.commit()
            return get_supplier(public_id, slug)


def export_suppliers(tenant_slug: str | None = None, status: str | None = None) -> dict:
    records = list_suppliers(tenant_slug, status)
    slug = _tenant_slug(tenant_slug)
    return {
        "tenantSlug": slug,
        "exported": len(records),
        "items": records,
    }


def build_summary(tenant_slug: str | None = None) -> dict:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        suppliers = [item for item in IN_MEMORY_STATE["suppliers"] if item["tenantSlug"] == slug]
        categories = [item for item in IN_MEMORY_STATE["categories"] if item["tenantSlug"] == slug]
        return {
            "tenantSlug": slug,
            "summary": {
                "suppliersTotal": len(suppliers),
                "active": sum(1 for item in suppliers if item["status"] == "active"),
                "watchlist": sum(1 for item in suppliers if item["status"] == "watchlist"),
                "categories": len(categories),
            }
        }

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT
                  count(*) AS suppliers_total,
                  count(*) FILTER (WHERE supplier.status = 'active') AS active_total,
                  count(*) FILTER (WHERE supplier.status = 'watchlist') AS watchlist_total
                FROM supplier.suppliers AS supplier
                JOIN identity.tenants AS tenant ON tenant.id = supplier.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            supplier_row = cursor.fetchone() or {}
            cursor.execute(
                """
                SELECT count(*) AS categories_total
                FROM supplier.categories AS category
                JOIN identity.tenants AS tenant ON tenant.id = category.tenant_id
                WHERE tenant.slug = %s
                """,
                (slug,),
            )
            category_row = cursor.fetchone() or {}
            return {
                "tenantSlug": slug,
                "summary": {
                    "suppliersTotal": int(supplier_row.get("suppliers_total", 0) or 0),
                    "active": int(supplier_row.get("active_total", 0) or 0),
                    "watchlist": int(supplier_row.get("watchlist_total", 0) or 0),
                    "categories": int(category_row.get("categories_total", 0) or 0),
                },
            }


def capability_catalog() -> dict:
    return {
        "service": settings.service_name,
        "repositoryDriver": settings.repository_driver,
        "capabilities": [
            {"key": "supplier.directory", "status": "ready"},
            {"key": "supplier.payables-profile", "status": "ready"},
            {"key": "supplier.cnpj-enrichment-ready", "status": "ready"},
        ],
    }
