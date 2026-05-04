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
    engagement_resend_api_key: str = os.getenv("ENGAGEMENT_RESEND_API_KEY", "")
    engagement_whatsapp_access_token: str = os.getenv("ENGAGEMENT_WHATSAPP_ACCESS_TOKEN", "")
    engagement_telegram_bot_token: str = os.getenv("ENGAGEMENT_TELEGRAM_BOT_TOKEN", "")
    engagement_meta_ads_access_token: str = os.getenv("ENGAGEMENT_META_ADS_ACCESS_TOKEN", "")
    billing_asaas_api_key: str = os.getenv("BILLING_ASAAS_API_KEY", "")
    billing_stripe_secret_key: str = os.getenv("BILLING_STRIPE_SECRET_KEY", "")
    billing_mercado_pago_access_token: str = os.getenv("BILLING_MERCADO_PAGO_ACCESS_TOKEN", "")
    documents_storage_driver: str = os.getenv("DOCUMENTS_STORAGE_DRIVER", "local")
    documents_storage_bucket: str = os.getenv("DOCUMENTS_STORAGE_BUCKET", "")
    documents_storage_endpoint: str = os.getenv("DOCUMENTS_STORAGE_ENDPOINT", "")
    documents_r2_account_id: str = os.getenv("DOCUMENTS_R2_ACCOUNT_ID", "")
    documents_r2_bucket: str = os.getenv("DOCUMENTS_R2_BUCKET", "")
    documents_clicksign_api_key: str = os.getenv("DOCUMENTS_CLICKSIGN_API_KEY", "")
    documents_docusign_access_token: str = os.getenv("DOCUMENTS_DOCUSIGN_ACCESS_TOKEN", "")
    crm_cnpj_provider_token: str = os.getenv("CRM_CNPJ_PROVIDER_TOKEN", "")
    crm_conecta_cnpj_api_key: str = os.getenv("CRM_CONECTA_CNPJ_API_KEY", "")
    webhook_hub_outbound_signing_secret: str = os.getenv("WEBHOOK_HUB_OUTBOUND_SIGNING_SECRET", "")


settings = Settings()
