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
	PostgresHost        string
	PostgresPort        string
	PostgresDatabase    string
	PostgresUser        string
	PostgresPassword    string
	PostgresSSLMode     string
}

func Load() Config {
	return Config{
		ServiceName:         "documents",
		HTTPAddress:         envOrDefault("DOCUMENTS_HTTP_ADDRESS", ":8086"),
		RepositoryDriver:    envOrDefault("DOCUMENTS_REPOSITORY_DRIVER", "memory"),
		BootstrapTenantSlug: envOrDefault("DOCUMENTS_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
		PostgresHost:        envOrDefault("DOCUMENTS_POSTGRES_HOST", "localhost"),
		PostgresPort:        envOrDefault("DOCUMENTS_POSTGRES_PORT", "5432"),
		PostgresDatabase:    envOrDefault("DOCUMENTS_POSTGRES_DB", envOrDefault("ERP_POSTGRES_DB", "erp")),
		PostgresUser:        envOrDefault("DOCUMENTS_POSTGRES_USER", envOrDefault("ERP_POSTGRES_USER", "erp")),
		PostgresPassword:    envOrDefault("DOCUMENTS_POSTGRES_PASSWORD", envOrDefault("ERP_POSTGRES_PASSWORD", "erp")),
		PostgresSSLMode:     envOrDefault("DOCUMENTS_POSTGRES_SSL_MODE", "disable"),
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
