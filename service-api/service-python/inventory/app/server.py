"""Bootstrap HTTP do servico inventory."""

from fastapi import FastAPI, HTTPException

from app.config.settings import settings
from app.infrastructure.postgres import postgres_ready
from app.security import install_security_middleware
from app.runtime import build_balances, build_costing_summary, build_cycle_count_variances, build_summary, capability_catalog, create_record, get_record, list_records, transition_record

app = FastAPI(title=settings.service_name)
install_security_middleware(app, settings.service_name)


@app.get("/health/live")
def live() -> dict:
    return {"service": settings.service_name, "status": "live"}


@app.get("/health/ready")
def ready() -> dict:
    return {"service": settings.service_name, "status": "ready"}


@app.get("/health/details")
def details() -> dict:
    dependencies = [{"name": "inventory-runtime", "status": "ready"}]
    if settings.repository_driver == "postgres":
        dependencies.insert(0, {"name": "postgresql", "status": "ready" if postgres_ready() else "not_ready"})
    return {"service": settings.service_name, "status": "ready", "dependencies": dependencies}


@app.get("/api/inventory/capabilities")
def capabilities() -> dict:
    return capability_catalog()


@app.get("/api/inventory/summary")
def summary(tenant_slug: str | None = None) -> dict:
    return build_summary(tenant_slug)


@app.get("/api/inventory/balances")
def balances(tenant_slug: str | None = None, sku: str | None = None, location_code: str | None = None) -> dict:
    return build_balances(tenant_slug, sku, location_code)


@app.get("/api/inventory/costing/summary")
def costing_summary(tenant_slug: str | None = None, sku: str | None = None) -> dict:
    return build_costing_summary(tenant_slug, sku)


@app.get("/api/inventory/cycle-counts/variances")
def cycle_count_variances(tenant_slug: str | None = None) -> dict:
    return build_cycle_count_variances(tenant_slug)



@app.get("/api/inventory/locations")
def list_locations(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("locations", tenant_slug, status)


@app.post("/api/inventory/locations")
def post_location(payload: dict) -> dict:
    try:
        return create_record("locations", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "inventory payload is invalid."}) from error


@app.get("/api/inventory/locations/{public_id}")
def get_location(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("locations", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "location_not_found", "message": "Location was not found."})
    return record


@app.patch("/api/inventory/locations/{public_id}/status")
def patch_location_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("locations", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "inventory status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "location_not_found", "message": "Location was not found."})
    return record



@app.get("/api/inventory/movements")
def list_movements(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("movements", tenant_slug, status)


@app.post("/api/inventory/movements")
def post_movement(payload: dict) -> dict:
    try:
        return create_record("movements", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "inventory payload is invalid."}) from error


@app.get("/api/inventory/movements/{public_id}")
def get_movement(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("movements", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "movement_not_found", "message": "Movement was not found."})
    return record


@app.patch("/api/inventory/movements/{public_id}/status")
def patch_movement_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("movements", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "inventory status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "movement_not_found", "message": "Movement was not found."})
    return record



@app.get("/api/inventory/reservations")
def list_reservations(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("reservations", tenant_slug, status)


@app.post("/api/inventory/reservations")
def post_reservation(payload: dict) -> dict:
    try:
        return create_record("reservations", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "inventory payload is invalid."}) from error


@app.get("/api/inventory/reservations/{public_id}")
def get_reservation(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("reservations", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "reservation_not_found", "message": "Reservation was not found."})
    return record


@app.patch("/api/inventory/reservations/{public_id}/status")
def patch_reservation_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("reservations", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "inventory status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "reservation_not_found", "message": "Reservation was not found."})
    return record



@app.get("/api/inventory/cycle-counts")
def list_cycle_counts(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("cycle-counts", tenant_slug, status)


@app.post("/api/inventory/cycle-counts")
def post_cycle_count(payload: dict) -> dict:
    try:
        return create_record("cycle-counts", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "inventory payload is invalid."}) from error


@app.get("/api/inventory/cycle-counts/{public_id}")
def get_cycle_count(public_id: str, tenant_slug: str | None = None) -> dict:
    record = get_record("cycle-counts", public_id, tenant_slug)
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "cycle_count_not_found", "message": "Cycle Count was not found."})
    return record


@app.patch("/api/inventory/cycle-counts/{public_id}/status")
def patch_cycle_count_status(public_id: str, payload: dict) -> dict:
    try:
        record = transition_record("cycle-counts", public_id, payload)
    except ValueError as error:
        raise HTTPException(status_code=400, detail={"code": str(error), "message": "inventory status payload is invalid."}) from error
    if record is None:
        raise HTTPException(status_code=404, detail={"code": "cycle_count_not_found", "message": "Cycle Count was not found."})
    return record


@app.get("/api/inventory/cost-layers")
def list_cost_layers(tenant_slug: str | None = None, status: str | None = None) -> list[dict]:
    return list_records("cost-layers", tenant_slug, status)


@app.post("/api/inventory/cost-layers")
def post_cost_layer(payload: dict) -> dict:
    try:
        return create_record("cost-layers", payload)
    except ValueError as error:
        status_code = 404 if str(error) == "tenant_not_found" else 400
        raise HTTPException(status_code=status_code, detail={"code": str(error), "message": "inventory payload is invalid."}) from error
