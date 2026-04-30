package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/telemetry"
)

func withCorrelation(logger *telemetry.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		correlationID := request.Header.Get("X-Correlation-Id")
		if correlationID == "" {
			correlationID = uuid.Nil.String()
		}

		writer.Header().Set("X-Correlation-Id", correlationID)
		logger.Printf("rentals request correlation=%s method=%s path=%s", correlationID, request.Method, request.URL.Path)
		next.ServeHTTP(writer, request)
	})
}
