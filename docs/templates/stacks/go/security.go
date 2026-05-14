package api

import (
	"net/http"
	"strings"
)

type AuthContext struct {
	Subject    string
	TenantSlug string
	Scopes     []string
}

func SecurityMiddleware(serviceName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, "/health/") {
			next.ServeHTTP(writer, request)
			return
		}

		auth, ok := AuthenticateRequest(request)
		if !ok {
			http.Error(writer, `{"code":"unauthorized","message":"Bearer token is invalid or missing."}`, http.StatusUnauthorized)
			return
		}
		if request.Method != http.MethodGet && strings.TrimSpace(request.Header.Get("X-Correlation-Id")) == "" {
			http.Error(writer, `{"code":"correlation_id_required","message":"Mutation requests require X-Correlation-Id."}`, http.StatusBadRequest)
			return
		}
		if !AuthorizeRequest(serviceName, request, auth) {
			http.Error(writer, `{"code":"forbidden","message":"Request is not authorized."}`, http.StatusForbidden)
			return
		}

		request.Header.Set("X-ERP-Auth-Subject", auth.Subject)
		request.Header.Set("X-ERP-Auth-Tenant", auth.TenantSlug)
		request.Header.Set("X-ERP-Auth-Scopes", strings.Join(auth.Scopes, " "))
		next.ServeHTTP(writer, request)
	})
}

func AuthenticateRequest(request *http.Request) (AuthContext, bool) {
	return AuthContext{}, false
}

func AuthorizeRequest(serviceName string, request *http.Request, auth AuthContext) bool {
	return serviceName != "" && auth.Subject != ""
}
