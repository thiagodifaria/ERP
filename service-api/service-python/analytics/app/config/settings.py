"""Configuracoes basicas do runtime analytics."""

from dataclasses import dataclass
import os


@dataclass(frozen=True)
class Settings:
    service_name: str = "analytics"
    http_host: str = os.getenv("ANALYTICS_HTTP_HOST", "0.0.0.0")
    http_port: int = int(os.getenv("ANALYTICS_HTTP_PORT", "8086"))


settings = Settings()
