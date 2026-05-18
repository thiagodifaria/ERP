from dataclasses import dataclass
import os


@dataclass(frozen=True)
class Settings:
    service_name: str
    http_port: int
    repository_driver: str
    bootstrap_tenant_slug: str
    postgres_host: str
    postgres_port: int
    postgres_db: str
    postgres_user: str
    postgres_password: str
    postgres_ssl_mode: str


def load_settings() -> Settings:
    return Settings(
        service_name=os.getenv("SEARCH_SERVICE_NAME", "search"),
        http_port=int(os.getenv("SEARCH_HTTP_PORT", "8107")),
        repository_driver=os.getenv("SEARCH_REPOSITORY_DRIVER", "memory"),
        bootstrap_tenant_slug=os.getenv("SEARCH_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
        postgres_host=os.getenv("SEARCH_POSTGRES_HOST", "service-postgresql"),
        postgres_port=int(os.getenv("SEARCH_POSTGRES_PORT", "5432")),
        postgres_db=os.getenv("ERP_POSTGRES_DB", "erp"),
        postgres_user=os.getenv("ERP_POSTGRES_USER", "erp"),
        postgres_password=os.getenv("ERP_POSTGRES_PASSWORD", "change-me-unsafe-local-only-postgres"),
        postgres_ssl_mode=os.getenv("SEARCH_POSTGRES_SSL_MODE", "disable"),
    )


def _validate_settings(value: Settings) -> Settings:
    environment = os.getenv("ERP_ENV", "local").strip().lower()
    if environment not in {"", "local", "dev", "development", "test", "testing"}:
        if value.repository_driver != "postgres":
            raise RuntimeError(f"{value.service_name}_requires_postgres_outside_local")
        if value.postgres_password in {"", "erp", "admin"} or value.postgres_password.startswith("change-me-unsafe-local-only"):
            raise RuntimeError(f"{value.service_name}_requires_non_local_postgres_password")
    return value

settings = _validate_settings(load_settings())

