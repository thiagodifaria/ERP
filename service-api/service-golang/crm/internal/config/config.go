// Config centraliza as configuracoes de runtime do servico.
// Segredos e variaveis de ambiente devem entrar por aqui.
package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceName         string
	HTTPAddress         string
	RepositoryDriver    string
	BootstrapTenantSlug string
	DocumentsBaseURL    string
	PostgresHost        string
	PostgresPort        string
	PostgresDatabase    string
	PostgresUser        string
	PostgresPassword    string
	PostgresSSLMode     string
}

func Load() Config {
	return Config{
		ServiceName:         "crm",
		HTTPAddress:         envOrDefault("CRM_HTTP_ADDRESS", ":8083"),
		RepositoryDriver:    envOrDefault("CRM_REPOSITORY_DRIVER", "memory"),
		BootstrapTenantSlug: envOrDefault("CRM_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
		DocumentsBaseURL:    envOrDefault("CRM_DOCUMENTS_BASE_URL", ""),
		PostgresHost:        envOrDefault("CRM_POSTGRES_HOST", "localhost"),
		PostgresPort:        envOrDefault("CRM_POSTGRES_PORT", "5432"),
		PostgresDatabase:    envOrDefault("CRM_POSTGRES_DB", envOrDefault("ERP_POSTGRES_DB", "erp")),
		PostgresUser:        envOrDefault("CRM_POSTGRES_USER", envOrDefault("ERP_POSTGRES_USER", "erp")),
		PostgresPassword:    envOrDefault("CRM_POSTGRES_PASSWORD", envOrDefault("ERP_POSTGRES_PASSWORD", "erp")),
		PostgresSSLMode:     envOrDefault("CRM_POSTGRES_SSL_MODE", "disable"),
	}
}

func (cfg Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresSSLMode,
	)
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
