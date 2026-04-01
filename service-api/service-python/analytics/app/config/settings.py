"""Configuracoes basicas do runtime analytics."""

from dataclasses import dataclass
import os


@dataclass(frozen=True)
class Settings:
    service_name: str = "analytics"
    http_host: str = os.getenv("ANALYTICS_HTTP_HOST", "0.0.0.0")
    http_port: int = int(os.getenv("ANALYTICS_HTTP_PORT", "8086"))
    repository_driver: str = os.getenv("ANALYTICS_REPOSITORY_DRIVER", "memory")
    postgres_host: str = os.getenv("ANALYTICS_POSTGRES_HOST", "service-postgresql")
    postgres_port: int = int(os.getenv("ANALYTICS_POSTGRES_PORT", "5432"))
    postgres_db: str = os.getenv("ANALYTICS_POSTGRES_DB", "erp")
    postgres_user: str = os.getenv("ANALYTICS_POSTGRES_USER", "erp")
    postgres_password: str = os.getenv("ANALYTICS_POSTGRES_PASSWORD", "erp")
    postgres_ssl_mode: str = os.getenv("ANALYTICS_POSTGRES_SSL_MODE", "disable")


settings = Settings()
