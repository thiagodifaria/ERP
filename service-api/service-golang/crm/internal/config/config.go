// Config centraliza as configuracoes de runtime do servico.
// Segredos e variaveis de ambiente devem entrar por aqui.
package config

import "os"

type Config struct {
	ServiceName string
	HTTPAddress string
}

func Load() Config {
	return Config{
		ServiceName: "crm",
		HTTPAddress: envOrDefault("CRM_HTTP_ADDRESS", ":8083"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
