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
        postgres_password=os.getenv("ERP_POSTGRES_PASSWORD", "erp"),
        postgres_ssl_mode=os.getenv("SEARCH_POSTGRES_SSL_MODE", "disable"),
    )


settings = load_settings()

