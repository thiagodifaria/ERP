package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func (handler LeadHandler) resolveRepositories(writer http.ResponseWriter, request *http.Request, bodyTenantSlug string) (repository.TenantRepositorySet, string) {
	return resolveTenantRepositories(writer, request, bodyTenantSlug, handler.repositories)
}

func (handler LeadNoteHandler) resolveRepositories(writer http.ResponseWriter, request *http.Request) (repository.TenantRepositorySet, string) {
	return resolveTenantRepositories(writer, request, "", handler.repositories)
}

func (handler CustomerHandler) resolveRepositories(writer http.ResponseWriter, request *http.Request) (repository.TenantRepositorySet, string) {
	return resolveTenantRepositories(writer, request, "", handler.repositories)
}

func resolveTenantRepositories(
	writer http.ResponseWriter,
	request *http.Request,
	bodyTenantSlug string,
	factory repository.TenantRepositoryFactory,
) (repository.TenantRepositorySet, string) {
	tenantSlug := strings.TrimSpace(request.URL.Query().Get("tenantSlug"))
	if tenantSlug == "" {
		tenantSlug = strings.TrimSpace(bodyTenantSlug)
	}
	if tenantSlug == "" {
		tenantSlug = factory.BootstrapTenantSlug()
	}

	bundle, err := factory.ForTenant(tenantSlug)
	if err != nil {
		writeNotFound(writer, "tenant_not_found", "Tenant was not found.")
		return repository.TenantRepositorySet{}, ""
	}

	return bundle, tenantSlug
}

func writeBadRequest(writer http.ResponseWriter, code string, message string) {
	writeErrorResponse(writer, http.StatusBadRequest, code, message)
}

func writeNotFound(writer http.ResponseWriter, code string, message string) {
	writeErrorResponse(writer, http.StatusNotFound, code, message)
}

func writeErrorResponse(writer http.ResponseWriter, status int, code string, message string) {
	writeJSON(writer, status, dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
