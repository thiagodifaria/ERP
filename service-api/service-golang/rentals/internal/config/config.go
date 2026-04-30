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
		ServiceName:         "rentals",
		HTTPAddress:         envOrDefault("RENTALS_HTTP_ADDRESS", ":8096"),
		RepositoryDriver:    envOrDefault("RENTALS_REPOSITORY_DRIVER", "memory"),
		BootstrapTenantSlug: envOrDefault("RENTALS_BOOTSTRAP_TENANT_SLUG", "bootstrap-ops"),
		DocumentsBaseURL:    envOrDefault("RENTALS_DOCUMENTS_BASE_URL", ""),
		PostgresHost:        envOrDefault("RENTALS_POSTGRES_HOST", "localhost"),
		PostgresPort:        envOrDefault("RENTALS_POSTGRES_PORT", "5432"),
		PostgresDatabase:    envOrDefault("RENTALS_POSTGRES_DB", envOrDefault("ERP_POSTGRES_DB", "erp")),
		PostgresUser:        envOrDefault("RENTALS_POSTGRES_USER", envOrDefault("ERP_POSTGRES_USER", "erp")),
		PostgresPassword:    envOrDefault("RENTALS_POSTGRES_PASSWORD", envOrDefault("ERP_POSTGRES_PASSWORD", "erp")),
		PostgresSSLMode:     envOrDefault("RENTALS_POSTGRES_SSL_MODE", "disable"),
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
