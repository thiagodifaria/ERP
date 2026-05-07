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
        service_name=os.getenv("SUPPORT_SERVICE_NAME", "support"),
        http_port=int(os.getenv("SUPPORT_HTTP_PORT", "8099")),
        repository_driver=os.getenv("SUPPORT_REPOSITORY_DRIVER", "memory"),
        bootstrap_tenant_slug=os.getenv("SUPPORT_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
        postgres_host=os.getenv("SUPPORT_POSTGRES_HOST", "service-postgresql"),
        postgres_port=int(os.getenv("SUPPORT_POSTGRES_PORT", "5432")),
        postgres_db=os.getenv("ERP_POSTGRES_DB", "erp"),
        postgres_user=os.getenv("ERP_POSTGRES_USER", "erp"),
        postgres_password=os.getenv("ERP_POSTGRES_PASSWORD", "erp"),
        postgres_ssl_mode=os.getenv("SUPPORT_POSTGRES_SSL_MODE", "disable"),
    )


settings = load_settings()
