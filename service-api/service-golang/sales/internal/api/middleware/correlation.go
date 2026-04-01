// Correlation garante um identificador minimo para rastrear requests.
// Integracoes reais de tracing podem substituir este fallback depois.
package middleware

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/telemetry"
)

func WithCorrelation(logger *telemetry.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		correlationID := request.Header.Get("X-Correlation-Id")
		if correlationID == "" {
			correlationID = "pending-correlation"
		}

		writer.Header().Set("X-Correlation-Id", correlationID)
		logger.Printf("request method=%s path=%s correlation_id=%s", request.Method, request.URL.Path, correlationID)
		next.ServeHTTP(writer, request)
	})
}
