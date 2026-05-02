package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func requiredTenantSlug(writer http.ResponseWriter, request *http.Request) (string, bool) {
	tenantSlug := strings.TrimSpace(request.URL.Query().Get("tenantSlug"))
	if tenantSlug == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(map[string]string{
			"code":    "tenant_slug_required",
			"message": "Tenant slug is required.",
		})
		return "", false
	}

	return tenantSlug, true
}

func respondAnalyticsDependencyError(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadGateway)
	_ = json.NewEncoder(writer).Encode(map[string]string{
		"code":    "edge_dependency_unavailable",
		"message": "A downstream analytics dependency is unavailable.",
	})
}

func readMapInt(payload map[string]any, path ...string) int {
	var current any = payload

	for _, key := range path {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return 0
		}

		nextValue, ok := currentMap[key]
		if !ok {
			return 0
		}

		current = nextValue
	}

	switch value := current.(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case json.Number:
		parsed, err := value.Int64()
		if err == nil {
			return int(parsed)
		}
	}

	return 0
}

func readMapFloat(payload map[string]any, path ...string) float64 {
	var current any = payload

	for _, key := range path {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return 0
		}

		nextValue, ok := currentMap[key]
		if !ok {
			return 0
		}

		current = nextValue
	}

	switch value := current.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case int:
		return float64(value)
	case int32:
		return float64(value)
	case int64:
		return float64(value)
	case json.Number:
		parsed, err := value.Float64()
		if err == nil {
			return parsed
		}
	}

	return 0
}

func readMapBool(payload map[string]any, path ...string) bool {
	var current any = payload

	for _, key := range path {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return false
		}

		nextValue, ok := currentMap[key]
		if !ok {
			return false
		}

		current = nextValue
	}

	value, ok := current.(bool)
	if !ok {
		return false
	}

	return value
}
