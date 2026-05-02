export type EngagementConfig = {
  serviceName: "engagement";
  repositoryDriver: "memory" | "postgres";
  bootstrapTenantSlug: string;
  crmBaseUrl: string;
  resendApiKey: string;
  whatsappAccessToken: string;
  telegramBotToken: string;
  metaAdsAccessToken: string;
  postgresHost: string;
  postgresPort: string;
  postgresDatabase: string;
  postgresUser: string;
  postgresPassword: string;
  postgresSslMode: string;
};

function ensureRepositoryDriver(value: string): "memory" | "postgres" {
  const normalizedValue = value.trim().toLowerCase();

  if (normalizedValue !== "memory" && normalizedValue !== "postgres") {
    throw new Error("engagement_repository_driver_invalid");
  }

  return normalizedValue;
}

function envOrDefault(key: string, fallback: string): string {
  return process.env[key] ?? fallback;
}

export function loadConfig(): EngagementConfig {
  return {
    serviceName: "engagement",
    repositoryDriver: ensureRepositoryDriver(envOrDefault("ENGAGEMENT_REPOSITORY_DRIVER", "memory")),
    bootstrapTenantSlug: envOrDefault("ENGAGEMENT_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
    crmBaseUrl: envOrDefault("ENGAGEMENT_CRM_BASE_URL", "http://localhost:8083"),
    resendApiKey: envOrDefault("ENGAGEMENT_RESEND_API_KEY", ""),
    whatsappAccessToken: envOrDefault("ENGAGEMENT_WHATSAPP_ACCESS_TOKEN", ""),
    telegramBotToken: envOrDefault("ENGAGEMENT_TELEGRAM_BOT_TOKEN", ""),
    metaAdsAccessToken: envOrDefault("ENGAGEMENT_META_ADS_ACCESS_TOKEN", ""),
    postgresHost: envOrDefault("ENGAGEMENT_POSTGRES_HOST", envOrDefault("ERP_POSTGRES_HOST", "localhost")),
    postgresPort: envOrDefault("ENGAGEMENT_POSTGRES_PORT", "5432"),
    postgresDatabase: envOrDefault("ENGAGEMENT_POSTGRES_DB", envOrDefault("ERP_POSTGRES_DB", "erp")),
    postgresUser: envOrDefault("ENGAGEMENT_POSTGRES_USER", envOrDefault("ERP_POSTGRES_USER", "erp")),
    postgresPassword: envOrDefault("ENGAGEMENT_POSTGRES_PASSWORD", envOrDefault("ERP_POSTGRES_PASSWORD", "erp")),
    postgresSslMode: envOrDefault("ENGAGEMENT_POSTGRES_SSL_MODE", "disable")
  };
}
