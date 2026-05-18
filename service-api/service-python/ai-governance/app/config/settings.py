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
    openai_api_key: str
    openai_model: str


def load_settings() -> Settings:
    return Settings(
        service_name=os.getenv("AI_GOVERNANCE_SERVICE_NAME", "ai-governance"),
        http_port=int(os.getenv("AI_GOVERNANCE_HTTP_PORT", "8108")),
        repository_driver=os.getenv("AI_GOVERNANCE_REPOSITORY_DRIVER", "memory"),
        bootstrap_tenant_slug=os.getenv("AI_GOVERNANCE_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
        postgres_host=os.getenv("AI_GOVERNANCE_POSTGRES_HOST", "service-postgresql"),
        postgres_port=int(os.getenv("AI_GOVERNANCE_POSTGRES_PORT", "5432")),
        postgres_db=os.getenv("ERP_POSTGRES_DB", "erp"),
        postgres_user=os.getenv("ERP_POSTGRES_USER", "erp"),
        postgres_password=os.getenv("ERP_POSTGRES_PASSWORD", "change-me-unsafe-local-only-postgres"),
        postgres_ssl_mode=os.getenv("AI_GOVERNANCE_POSTGRES_SSL_MODE", "disable"),
        openai_api_key=os.getenv("OPENAI_API_KEY", ""),
        openai_model=os.getenv("OPENAI_MODEL", "gpt-4.1-mini"),
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
