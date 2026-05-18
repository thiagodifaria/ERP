"""Configuracoes basicas do runtime simulation."""

from dataclasses import dataclass
import os


@dataclass(frozen=True)
class Settings:
    service_name: str = "simulation"
    http_host: str = os.getenv("SIMULATION_HTTP_HOST", "0.0.0.0")
    http_port: int = int(os.getenv("SIMULATION_HTTP_PORT", "8094"))
    repository_driver: str = os.getenv("SIMULATION_REPOSITORY_DRIVER", "memory")
    postgres_host: str = os.getenv("SIMULATION_POSTGRES_HOST", "service-postgresql")
    postgres_port: int = int(os.getenv("SIMULATION_POSTGRES_PORT", "5432"))
    postgres_db: str = os.getenv("SIMULATION_POSTGRES_DB", "erp")
    postgres_user: str = os.getenv("SIMULATION_POSTGRES_USER", "erp")
    postgres_password: str = os.getenv("SIMULATION_POSTGRES_PASSWORD", "change-me-unsafe-local-only-postgres")
    postgres_ssl_mode: str = os.getenv("SIMULATION_POSTGRES_SSL_MODE", "disable")


def _validate_settings(value: Settings) -> Settings:
    environment = os.getenv("ERP_ENV", "local").strip().lower()
    if environment not in {"", "local", "dev", "development", "test", "testing"}:
        if value.repository_driver != "postgres":
            raise RuntimeError(f"{value.service_name}_requires_postgres_outside_local")
        if value.postgres_password in {"", "erp", "admin"} or value.postgres_password.startswith("change-me-unsafe-local-only"):
            raise RuntimeError(f"{value.service_name}_requires_non_local_postgres_password")
    return value

settings = _validate_settings(Settings())
