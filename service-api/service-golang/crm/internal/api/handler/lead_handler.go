// LeadHandler exposes the first operational lead endpoints of the CRM.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type LeadHandler struct {
	repositories repository.TenantRepositoryFactory
}

func NewLeadHandler(repositories repository.TenantRepositoryFactory) LeadHandler {
	return LeadHandler{repositories: repositories}
}

func (handler LeadHandler) List(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	leads := query.NewListLeads(bundle.LeadRepository).Execute(query.LeadFilters{
		Status:      request.URL.Query().Get("status"),
		Source:      request.URL.Query().Get("source"),
		OwnerUserID: request.URL.Query().Get("ownerUserId"),
		Search:      request.URL.Query().Get("q"),
		Assigned:    request.URL.Query().Get("assigned"),
	})
	response := make([]dto.LeadResponse, 0, len(leads))

	for _, lead := range leads {
		response = append(response, mapLead(lead))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler LeadHandler) Summary(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	summary := query.NewGetLeadPipelineSummary(bundle.LeadRepository).Execute(query.LeadFilters{
		Status:      request.URL.Query().Get("status"),
		Source:      request.URL.Query().Get("source"),
		OwnerUserID: request.URL.Query().Get("ownerUserId"),
		Search:      request.URL.Query().Get("q"),
		Assigned:    request.URL.Query().Get("assigned"),
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.LeadSummaryResponse{
		Total:      summary.Total,
		Assigned:   summary.Assigned,
		Unassigned: summary.Unassigned,
		ByStatus:   summary.ByStatus,
		BySource:   summary.BySource,
	})
}

func (handler LeadHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	lead := query.NewGetLeadByPublicID(bundle.LeadRepository).Execute(request.PathValue("publicId"))
	if lead == nil {
		writeNotFound(writer, "lead_not_found", "Lead was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(mapLead(*lead))
}

func (handler LeadHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateLeadRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(writer, "invalid_json", "Request body is invalid.")
		return
	}

	bundle, _ := handler.resolveRepositories(writer, request, payload.TenantSlug)
	if bundle.LeadRepository == nil {
		return
	}

	result := command.NewCreateLead(
		bundle.LeadRepository,
		bundle.RelationshipEventRepository,
		bundle.OutboxEventRepository,
	).Execute(command.CreateLeadInput{
		Name:        payload.Name,
		Email:       payload.Email,
		Source:      payload.Source,
		OwnerUserID: payload.OwnerUserID,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writeErrorResponse(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.Conflict:
		writeErrorResponse(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
	default:
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(mapLead(*result.Lead))
	}
}

func (handler LeadHandler) Convert(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	result := command.NewConvertLeadToCustomer(
		bundle.LeadRepository,
		bundle.CustomerRepository,
		bundle.RelationshipEventRepository,
		bundle.OutboxEventRepository,
	).Execute(command.ConvertLeadToCustomerInput{
		LeadPublicID: request.PathValue("publicId"),
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writeErrorResponse(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.Conflict:
		writeErrorResponse(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeErrorResponse(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(dto.ConvertLeadResponse{
			Lead:     mapLead(*result.Lead),
			Customer: mapCustomer(*result.Customer),
		})
	}
}

func (handler LeadHandler) UpdateProfile(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateLeadProfileRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(writer, "invalid_json", "Request body is invalid.")
		return
	}

	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	result := command.NewUpdateLeadProfile(
		bundle.LeadRepository,
		bundle.RelationshipEventRepository,
		bundle.OutboxEventRepository,
	).Execute(command.UpdateLeadProfileInput{
		PublicID: request.PathValue("publicId"),
		Name:     payload.Name,
		Email:    payload.Email,
		Source:   payload.Source,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writeErrorResponse(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.Conflict:
		writeErrorResponse(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeErrorResponse(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(mapLead(*result.Lead))
	}
}

func (handler LeadHandler) UpdateOwner(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateLeadOwnerRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(writer, "invalid_json", "Request body is invalid.")
		return
	}

	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	result := command.NewUpdateLeadOwner(
		bundle.LeadRepository,
		bundle.RelationshipEventRepository,
		bundle.OutboxEventRepository,
	).Execute(command.UpdateLeadOwnerInput{
		PublicID:    request.PathValue("publicId"),
		OwnerUserID: payload.OwnerUserID,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writeErrorResponse(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeErrorResponse(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(mapLead(*result.Lead))
	}
}

func (handler LeadHandler) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateLeadStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(writer, "invalid_json", "Request body is invalid.")
		return
	}

	bundle, _ := handler.resolveRepositories(writer, request, "")
	if bundle.LeadRepository == nil {
		return
	}

	result := command.NewUpdateLeadStatus(
		bundle.LeadRepository,
		bundle.RelationshipEventRepository,
		bundle.OutboxEventRepository,
	).Execute(command.UpdateLeadStatusInput{
		PublicID: request.PathValue("publicId"),
		Status:   payload.Status,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writeErrorResponse(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeErrorResponse(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(mapLead(*result.Lead))
	}
}

func mapLead(lead entity.Lead) dto.LeadResponse {
	return dto.LeadResponse{
		PublicID:    lead.PublicID,
		Name:        lead.Name,
		Email:       lead.Email,
		Source:      lead.Source,
		Status:      lead.Status,
		OwnerUserID: lead.OwnerUserID,
	}
}
