// Este middleware garante um correlation id minimo na borda.
// Enriquecimento completo entra quando a observabilidade crescer.
package middleware

import (
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/telemetry"
)

func WithCorrelation(logger *telemetry.Logger, next http.Handler) http.Handler {
  return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
    correlationID := request.Header.Get("X-Correlation-Id")
    if correlationID == "" {
      correlationID = "pending-correlation"
    }

    writer.Header().Set("X-Correlation-Id", correlationID)
    logger.Printf("request %s %s correlation=%s", request.Method, request.URL.Path, correlationID)

    next.ServeHTTP(writer, request)
  })
}
