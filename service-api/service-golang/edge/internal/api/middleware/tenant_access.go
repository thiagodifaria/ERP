// Este middleware exige tenant e sessao antes de expor leituras operacionais de tenant.
package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/infrastructure/integration"
)

func WithTenantAccess(identityBaseURL string, resolver integration.TenantAccessResolver, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tenantSlug := strings.TrimSpace(request.URL.Query().Get("tenantSlug"))
		if tenantSlug == "" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(map[string]string{
				"code":    "tenant_slug_required",
				"message": "Tenant slug is required.",
			})
			return
		}

		sessionToken := extractBearerToken(request.Header.Get("Authorization"))
		if sessionToken == "" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(writer).Encode(map[string]string{
				"code":    "session_required",
				"message": "Session token is required.",
			})
			return
		}

		resolution, err := resolver.ResolveTenantAccess(request.Context(), identityBaseURL, tenantSlug, sessionToken)
		if err != nil {
			if accessError, ok := err.(integration.AccessResolutionError); ok {
				statusCode := accessError.StatusCode
				if statusCode == 0 {
					statusCode = http.StatusBadGateway
				}

				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(statusCode)
				_ = json.NewEncoder(writer).Encode(map[string]string{
					"code":    fallbackCode(accessError.Code, statusCode),
					"message": fallbackMessage(accessError.Message, statusCode),
				})
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadGateway)
			_ = json.NewEncoder(writer).Encode(map[string]string{
				"code":    "edge_identity_unavailable",
				"message": "Identity dependency is unavailable.",
			})
			return
		}

		request.Header.Set("X-ERP-Tenant-Slug", resolution.TenantSlug)
		request.Header.Set("X-ERP-User-Public-Id", resolution.UserPublicID)
		request.Header.Set("X-ERP-User-Roles", strings.Join(resolution.RoleCodes, ","))

		next.ServeHTTP(writer, request)
	})
}

func extractBearerToken(authorization string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(strings.ToLower(authorization), strings.ToLower(prefix)) {
		return ""
	}

	return strings.TrimSpace(authorization[len(prefix):])
}

func fallbackCode(code string, statusCode int) string {
	if code != "" {
		return code
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return "invalid_session"
	case http.StatusForbidden:
		return "tenant_scope_forbidden"
	case http.StatusNotFound:
		return "tenant_not_found"
	default:
		return "edge_identity_error"
	}
}

func fallbackMessage(message string, statusCode int) string {
	if message != "" {
		return message
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return "Session is invalid."
	case http.StatusForbidden:
		return "User does not have tenant access."
	case http.StatusNotFound:
		return "Tenant was not found."
	default:
		return "Identity dependency returned an unexpected response."
	}
}
