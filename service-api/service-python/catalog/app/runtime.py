from __future__ import annotations

from datetime import datetime, timezone
import json
import uuid

from app.config.settings import settings
from app.infrastructure.postgres import connect


IN_MEMORY_STATE = {
    "categories": [],
    "items": [],
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


def _normalize_limit(limit: int | None, fallback: int = 50) -> int:
    if limit is None or limit <= 0:
        return fallback
    return min(limit, 100)


def _paginate(records: list[dict], cursor: str | None, limit: int, cursor_field: str = "publicId") -> dict:
    page_limit = _normalize_limit(limit)
    start_index = 0
    if cursor:
        for index, record in enumerate(records):
            if str(record.get(cursor_field, "")) == cursor:
                start_index = index + 1
                break
    items = records[start_index : start_index + page_limit]
    next_cursor = None
    if start_index + page_limit < len(records) and items:
        next_cursor = str(items[-1].get(cursor_field, ""))
    return {
        "items": items,
        "pageInfo": {
            "cursor": cursor,
            "limit": page_limit,
            "returned": len(items),
            "nextCursor": next_cursor,
            "hasMore": next_cursor is not None,
        },
    }


def list_categories(tenant_slug: str | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        records = [category for category in IN_MEMORY_STATE["categories"] if category["tenantSlug"] == slug]
        return sorted(records, key=lambda item: item["name"])

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT category.public_id, category.category_key, category.name, category.active, category.created_at
                FROM catalog.categories AS category
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
                    "key": row["category_key"],
                    "name": row["name"],
                    "active": row["active"],
                    "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
                for row in cursor.fetchall()
            ]


def list_categories_page(tenant_slug: str | None = None, cursor: str | None = None, limit: int = 50) -> dict:
    payload = _paginate(list_categories(tenant_slug), cursor, limit)
    payload["tenantSlug"] = _tenant_slug(tenant_slug)
    return payload


def create_category(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    category_key = (payload.get("key") or "").strip().lower().replace(" ", "-")
    name = (payload.get("name") or "").strip()
    if category_key == "":
        raise ValueError("category_key_required")
    if name == "":
        raise ValueError("category_name_required")

    category = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "key": category_key,
        "name": name,
        "active": bool(payload.get("active", True)),
        "createdAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["categories"].append(category)
        return category

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            cursor.execute(
                """
                INSERT INTO catalog.categories (tenant_id, public_id, category_key, name, active)
                VALUES (%s, %s, %s, %s, %s)
                RETURNING created_at
                """,
                (tenant_id, category["publicId"], category_key, name, category["active"]),
            )
            row = cursor.fetchone()
            connection.commit()
            category["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            return category


def list_items(tenant_slug: str | None = None, item_type: str | None = None, active: bool | None = None) -> list[dict]:
    slug = _tenant_slug(tenant_slug)
    normalized_item_type = (item_type or "").strip().lower()
    if settings.repository_driver != "postgres":
        items = [item for item in IN_MEMORY_STATE["items"] if item["tenantSlug"] == slug]
        if normalized_item_type != "":
            items = [item for item in items if item["itemType"] == normalized_item_type]
        if active is not None:
            items = [item for item in items if item["active"] is active]
        return sorted(items, key=lambda item: item["name"])

    clauses = ["tenant.slug = %s"]
    params: list[object] = [slug]
    if normalized_item_type != "":
        clauses.append("item.item_type = %s")
        params.append(normalized_item_type)
    if active is not None:
        clauses.append("item.active = %s")
        params.append(active)

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                f"""
                SELECT item.public_id, item.sku, item.name, item.item_type, item.unit_code, item.price_base_cents,
                       item.currency_code, item.active, item.version_number, item.attributes_json, item.created_at,
                       item.updated_at, category.public_id AS category_public_id, category.name AS category_name
                FROM catalog.items AS item
                JOIN identity.tenants AS tenant ON tenant.id = item.tenant_id
                LEFT JOIN catalog.categories AS category ON category.id = item.category_id
                WHERE {" AND ".join(clauses)}
                ORDER BY item.name
                """,
                params,
            )
            return [_map_item_row(row, slug) for row in cursor.fetchall()]


def list_items_page(tenant_slug: str | None = None, item_type: str | None = None, active: bool | None = None, cursor: str | None = None, limit: int = 50) -> dict:
    payload = _paginate(list_items(tenant_slug, item_type, active), cursor, limit)
    payload["tenantSlug"] = _tenant_slug(tenant_slug)
    return payload


def get_item(public_id: str, tenant_slug: str | None = None) -> dict | None:
    slug = _tenant_slug(tenant_slug)
    if settings.repository_driver != "postgres":
        for item in IN_MEMORY_STATE["items"]:
            if item["tenantSlug"] == slug and item["publicId"] == public_id:
                return item
        return None

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                SELECT item.public_id, item.sku, item.name, item.item_type, item.unit_code, item.price_base_cents,
                       item.currency_code, item.active, item.version_number, item.attributes_json, item.created_at,
                       item.updated_at, category.public_id AS category_public_id, category.name AS category_name
                FROM catalog.items AS item
                JOIN identity.tenants AS tenant ON tenant.id = item.tenant_id
                LEFT JOIN catalog.categories AS category ON category.id = item.category_id
                WHERE tenant.slug = %s AND item.public_id = %s
                """,
                (slug, public_id),
            )
            row = cursor.fetchone()
            return None if row is None else _map_item_row(row, slug)


def create_item(payload: dict) -> dict:
    slug = _tenant_slug(payload.get("tenantSlug"))
    sku = (payload.get("sku") or "").strip().upper()
    name = (payload.get("name") or "").strip()
    item_type = (payload.get("itemType") or "").strip().lower()
    unit_code = (payload.get("unitCode") or "").strip().lower()
    currency_code = (payload.get("currencyCode") or "BRL").strip().upper()
    category_public_id = (payload.get("categoryPublicId") or "").strip()
    attributes = payload.get("attributes") or {}
    if sku == "":
        raise ValueError("catalog_sku_required")
    if name == "":
        raise ValueError("catalog_name_required")
    if item_type not in {"product", "service"}:
        raise ValueError("catalog_item_type_invalid")
    if unit_code == "":
        raise ValueError("catalog_unit_code_required")

    item = {
        "publicId": str(uuid.uuid4()),
        "tenantSlug": slug,
        "sku": sku,
        "name": name,
        "itemType": item_type,
        "unitCode": unit_code,
        "priceBaseCents": int(payload.get("priceBaseCents") or 0),
        "currencyCode": currency_code,
        "active": bool(payload.get("active", True)),
        "versionNumber": 1,
        "attributes": attributes,
        "category": None,
        "createdAt": utc_now(),
        "updatedAt": utc_now(),
    }

    if settings.repository_driver != "postgres":
        IN_MEMORY_STATE["items"].append(item)
        return item

    with connect() as connection:
        with connection.cursor() as cursor:
            tenant_id = _find_tenant_id(cursor, slug)
            category_id = None
            category_name = None
            if category_public_id != "":
                cursor.execute(
                    """
                    SELECT category.id, category.name
                    FROM catalog.categories AS category
                    JOIN identity.tenants AS tenant ON tenant.id = category.tenant_id
                    WHERE tenant.slug = %s AND category.public_id = %s
                    """,
                    (slug, category_public_id),
                )
                category_row = cursor.fetchone()
                if category_row is None:
                    raise ValueError("catalog_category_not_found")
                category_id = int(category_row["id"])
                category_name = category_row["name"]

            cursor.execute(
                """
                INSERT INTO catalog.items (
                  tenant_id, category_id, public_id, sku, name, item_type, unit_code, price_base_cents,
                  currency_code, active, version_number, attributes_json
                )
                VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, 1, %s::jsonb)
                RETURNING created_at, updated_at
                """,
                (
                    tenant_id,
                    category_id,
                    item["publicId"],
                    sku,
                    name,
                    item_type,
                    unit_code,
                    item["priceBaseCents"],
                    currency_code,
                    item["active"],
                    json.dumps(attributes),
                ),
            )
            row = cursor.fetchone()
            connection.commit()
            item["createdAt"] = row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            item["updatedAt"] = row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
            if category_id is not None:
                item["category"] = {"publicId": category_public_id, "name": category_name}
            return item


def bulk_create_items(payload: dict) -> dict:
    items = payload.get("items") or []
    results: list[dict] = []
    errors: list[dict] = []
    for index, item in enumerate(items):
        try:
            results.append(create_item(item))
        except ValueError as error:
            errors.append({"index": index, "code": str(error), "message": "Catalog item payload is invalid."})

    return {
        "tenantSlug": _tenant_slug(payload.get("tenantSlug")),
        "results": results,
        "errors": errors,
        "summary": {
            "requested": len(items),
            "succeeded": len(results),
            "failed": len(errors),
            "partialSuccess": len(results) > 0 and len(errors) > 0,
        },
    }


def update_item(public_id: str, payload: dict) -> dict | None:
    current = get_item(public_id, payload.get("tenantSlug"))
    if current is None:
        return None

    slug = current["tenantSlug"]
    new_name = (payload.get("name") or current["name"]).strip()
    new_price = int(payload.get("priceBaseCents", current["priceBaseCents"]))
    new_active = bool(payload.get("active", current["active"]))
    new_attributes = payload.get("attributes", current["attributes"])

    if settings.repository_driver != "postgres":
        current["name"] = new_name
        current["priceBaseCents"] = new_price
        current["active"] = new_active
        current["attributes"] = new_attributes
        current["versionNumber"] = int(current["versionNumber"]) + 1
        current["updatedAt"] = utc_now()
        return current

    with connect() as connection:
        with connection.cursor() as cursor:
            cursor.execute(
                """
                UPDATE catalog.items AS item
                SET name = %s,
                    price_base_cents = %s,
                    active = %s,
                    attributes_json = %s::jsonb,
                    version_number = item.version_number + 1
                FROM identity.tenants AS tenant
                WHERE tenant.id = item.tenant_id
                  AND tenant.slug = %s
                  AND item.public_id = %s
                RETURNING item.public_id, item.sku, item.name, item.item_type, item.unit_code, item.price_base_cents,
                          item.currency_code, item.active, item.version_number, item.attributes_json, item.created_at,
                          item.updated_at
                """,
                (new_name, new_price, new_active, json.dumps(new_attributes), slug, public_id),
            )
            row = cursor.fetchone()
            if row is None:
                return None
            connection.commit()
            return _map_item_row(row, slug)


def capability_catalog() -> dict:
    return {
        "service": settings.service_name,
        "domains": ["product", "service"],
        "supportsVersioning": True,
        "supportsActivation": True,
        "supportsCategories": True,
        "supportsCursorPagination": True,
        "supportsBulk": True,
        "repositoryDriver": settings.repository_driver,
    }


def _map_item_row(row: dict, tenant_slug: str) -> dict:
    category = None
    if row.get("category_public_id") is not None:
        category = {
            "publicId": row["category_public_id"],
            "name": row["category_name"],
        }

    return {
        "publicId": row["public_id"],
        "tenantSlug": tenant_slug,
        "sku": row["sku"],
        "name": row["name"],
        "itemType": row["item_type"],
        "unitCode": row["unit_code"],
        "priceBaseCents": row["price_base_cents"],
        "currencyCode": row["currency_code"],
        "active": row["active"],
        "versionNumber": row["version_number"],
        "attributes": row["attributes_json"] or {},
        "category": category,
        "createdAt": row["created_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "updatedAt": row["updated_at"].astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
    }
