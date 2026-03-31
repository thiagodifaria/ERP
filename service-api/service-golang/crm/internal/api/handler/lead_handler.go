// LeadHandler exposes the first operational lead endpoints of the CRM.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type LeadHandler struct {
	listLeads         query.ListLeads
	getLeadSummary    query.GetLeadPipelineSummary
	getLeadByPublicID query.GetLeadByPublicID
	createLead        command.CreateLead
	updateLeadStatus  command.UpdateLeadStatus
}

func NewLeadHandler(
	listLeads query.ListLeads,
	getLeadSummary query.GetLeadPipelineSummary,
	getLeadByPublicID query.GetLeadByPublicID,
	createLead command.CreateLead,
	updateLeadStatus command.UpdateLeadStatus,
) LeadHandler {
	return LeadHandler{
		listLeads:         listLeads,
		getLeadSummary:    getLeadSummary,
		getLeadByPublicID: getLeadByPublicID,
		createLead:        createLead,
		updateLeadStatus:  updateLeadStatus,
	}
}

func (handler LeadHandler) List(writer http.ResponseWriter, request *http.Request) {
	leads := handler.listLeads.Execute(query.LeadFilters{
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
	summary := handler.getLeadSummary.Execute(query.LeadFilters{
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
	lead := handler.getLeadByPublicID.Execute(request.PathValue("publicId"))
	if lead == nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    "lead_not_found",
			Message: "Lead was not found.",
		})
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(mapLead(*lead))
}

func (handler LeadHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateLeadRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    "invalid_json",
			Message: "Request body is invalid.",
		})
		return
	}

	result := handler.createLead.Execute(command.CreateLeadInput{
		Name:        payload.Name,
		Email:       payload.Email,
		Source:      payload.Source,
		OwnerUserID: payload.OwnerUserID,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    result.ErrorCode,
			Message: result.ErrorText,
		})
	case result.Conflict:
		writer.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    result.ErrorCode,
			Message: result.ErrorText,
		})
	default:
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(mapLead(*result.Lead))
	}
}

func (handler LeadHandler) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateLeadStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    "invalid_json",
			Message: "Request body is invalid.",
		})
		return
	}

	result := handler.updateLeadStatus.Execute(command.UpdateLeadStatusInput{
		PublicID: request.PathValue("publicId"),
		Status:   payload.Status,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    result.ErrorCode,
			Message: result.ErrorText,
		})
	case result.NotFound:
		writer.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    result.ErrorCode,
			Message: result.ErrorText,
		})
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
