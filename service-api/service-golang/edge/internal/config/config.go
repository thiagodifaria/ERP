// Config centraliza as configuracoes de runtime do servico.
// Segredos e variaveis de ambiente devem entrar por aqui.
package config

import (
	"os"
	"time"
)

type Config struct {
	ServiceName            string
	HTTPAddress            string
	DownstreamTimeout      time.Duration
	IdentityBaseURL        string
	CRMBaseURL             string
	WorkflowControlBaseURL string
	WorkflowRuntimeBaseURL string
	AnalyticsBaseURL       string
	WebhookHubBaseURL      string
}

func Load() Config {
	return Config{
		ServiceName:            "edge",
		HTTPAddress:            envOrDefault("EDGE_HTTP_ADDRESS", ":8080"),
		DownstreamTimeout:      durationOrDefault("EDGE_DOWNSTREAM_TIMEOUT", 1500*time.Millisecond),
		IdentityBaseURL:        envOrDefault("EDGE_IDENTITY_BASE_URL", "http://identity:8080"),
		CRMBaseURL:             envOrDefault("EDGE_CRM_BASE_URL", "http://crm:8083"),
		WorkflowControlBaseURL: envOrDefault("EDGE_WORKFLOW_CONTROL_BASE_URL", "http://workflow-control:8084"),
		WorkflowRuntimeBaseURL: envOrDefault("EDGE_WORKFLOW_RUNTIME_BASE_URL", "http://workflow-runtime:8085"),
		AnalyticsBaseURL:       envOrDefault("EDGE_ANALYTICS_BASE_URL", "http://analytics:8086"),
		WebhookHubBaseURL:      envOrDefault("EDGE_WEBHOOK_HUB_BASE_URL", "http://webhook-hub:8082"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
