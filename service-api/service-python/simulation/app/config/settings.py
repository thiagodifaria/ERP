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
    postgres_password: str = os.getenv("SIMULATION_POSTGRES_PASSWORD", "erp")
    postgres_ssl_mode: str = os.getenv("SIMULATION_POSTGRES_SSL_MODE", "disable")


settings = Settings()
