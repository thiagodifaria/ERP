package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

type jwtClaims struct {
	Subject    string          `json:"sub"`
	UserID     string          `json:"user_public_id"`
	TenantSlug string          `json:"tenant_slug"`
	Tenant     string          `json:"tenant"`
	Scope      json.RawMessage `json:"scope"`
	ExpiresAt  int64           `json:"exp"`
}

func WithSecurity(serviceName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !securityEnforced() || strings.HasPrefix(request.URL.Path, "/health/") {
			next.ServeHTTP(writer, request)
			return
		}

		if requiresCorrelation(request) && strings.TrimSpace(request.Header.Get("X-Correlation-Id")) == "" {
			writeSecurityError(writer, http.StatusBadRequest, "correlation_id_required", "Mutation requests require X-Correlation-Id.")
			return
		}

		subject, tenantSlug, scopes, ok := authenticateRequest(request)
		if !ok {
			writeSecurityError(writer, http.StatusUnauthorized, "unauthorized", "Bearer token is invalid or missing.")
			return
		}

		request.Header.Set("X-ERP-Auth-Subject", subject)
		request.Header.Set("X-ERP-Auth-Tenant", tenantSlug)
		request.Header.Set("X-ERP-Auth-Scopes", strings.Join(scopes, " "))

		if !authorizeOpenFGA(serviceName, request, subject, tenantSlug) {
			writeSecurityError(writer, http.StatusForbidden, "openfga_denied", "OpenFGA denied the request.")
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func securityEnforced() bool {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("ERP_AUTH_ENFORCEMENT")))
	if mode == "disabled" || mode == "off" || mode == "false" {
		return false
	}
	if mode == "enforced" || mode == "strict" || mode == "true" {
		return true
	}
	environment := strings.ToLower(strings.TrimSpace(os.Getenv("ERP_ENV")))
	return environment != "" && environment != "local" && environment != "dev" && environment != "development" && environment != "test" && environment != "testing"
}

func requiresCorrelation(request *http.Request) bool {
	return request.Method != http.MethodGet && request.Method != http.MethodHead && request.Method != http.MethodOptions
}

func authenticateRequest(request *http.Request) (string, string, []string, bool) {
	header := request.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return "", "", nil, false
	}

	token := strings.TrimSpace(header[len("Bearer "):])
	if internalToken := strings.TrimSpace(os.Getenv("ERP_INTERNAL_SERVICE_TOKEN")); internalToken != "" && subtle.ConstantTimeCompare([]byte(token), []byte(internalToken)) == 1 {
		return "service:internal", resolveTenant(request), []string{"service"}, true
	}

	claims, ok := verifyJWT(token)
	if !ok {
		return "", "", nil, false
	}
	subject := claims.Subject
	if subject == "" {
		subject = claims.UserID
	}
	tenantSlug := claims.TenantSlug
	if tenantSlug == "" {
		tenantSlug = claims.Tenant
	}
	if tenantSlug == "" {
		tenantSlug = resolveTenant(request)
	}
	return subject, tenantSlug, parseScopes(claims.Scope), subject != ""
}

func verifyJWT(token string) (jwtClaims, bool) {
	secret := os.Getenv("ERP_JWT_HS256_SECRET")
	parts := strings.Split(token, ".")
	if secret == "" || len(parts) != 3 {
		return jwtClaims{}, false
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return jwtClaims{}, false
	}
	var header struct {
		Algorithm string `json:"alg"`
	}
	if json.Unmarshal(headerBytes, &header) != nil || header.Algorithm != "HS256" {
		return jwtClaims{}, false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(parts[0] + "." + parts[1]))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(parts[2]), []byte(expected)) != 1 {
		return jwtClaims{}, false
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return jwtClaims{}, false
	}
	var claims jwtClaims
	if json.Unmarshal(payloadBytes, &claims) != nil {
		return jwtClaims{}, false
	}
	if claims.ExpiresAt > 0 && time.Unix(claims.ExpiresAt, 0).Before(time.Now().UTC()) {
		return jwtClaims{}, false
	}
	return claims, true
}

func authorizeOpenFGA(serviceName string, request *http.Request, subject string, tenantSlug string) bool {
	if strings.ToLower(os.Getenv("ERP_OPENFGA_ENFORCEMENT")) != "true" {
		return true
	}
	baseURL := strings.TrimRight(os.Getenv("OPENFGA_BASE_URL"), "/")
	storeID := os.Getenv("OPENFGA_STORE_ID")
	if baseURL == "" || storeID == "" {
		return false
	}
	relation := "read"
	if requiresCorrelation(request) {
		relation = "write"
	}
	object := "service:" + normalizeObject(serviceName)
	if tenantSlug != "" {
		object = "tenant:" + normalizeObject(tenantSlug)
	}
	user := subject
	if !strings.HasPrefix(user, "service:") {
		user = "user:" + user
	}
	payload := map[string]any{
		"tuple_key": map[string]string{
			"user":     user,
			"relation": relation,
			"object":   object,
		},
	}
	if modelID := os.Getenv("OPENFGA_AUTHORIZATION_MODEL_ID"); modelID != "" {
		payload["authorization_model_id"] = modelID
	}
	body, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Post(baseURL+"/stores/"+storeID+"/check", "application/json", bytes.NewReader(body))
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return false
	}
	var result struct {
		Allowed bool `json:"allowed"`
	}
	return json.NewDecoder(response.Body).Decode(&result) == nil && result.Allowed
}

func resolveTenant(request *http.Request) string {
	for _, header := range []string{"X-Tenant-Slug", "X-ERP-Tenant-Slug"} {
		if value := strings.TrimSpace(request.Header.Get(header)); value != "" {
			return value
		}
	}
	return strings.TrimSpace(request.URL.Query().Get("tenant_slug"))
}

func parseScopes(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return strings.Fields(text)
	}
	var values []string
	if json.Unmarshal(raw, &values) == nil {
		return values
	}
	return nil
}

func normalizeObject(value string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(value)), " ", "-")
}

func writeSecurityError(writer http.ResponseWriter, status int, code string, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(map[string]string{"code": code, "message": message})
}
