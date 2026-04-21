// OpportunityHandler expoe a trilha operacional de oportunidades do Sales MVP.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
)

type OpportunityHandler struct {
	listOpportunities        query.ListOpportunities
	getOpportunitySummary    query.GetOpportunitySummary
	getOpportunityByPublicID query.GetOpportunityByPublicID
	createOpportunity        command.CreateOpportunity
	updateOpportunity        command.UpdateOpportunityProfile
	updateOpportunityStage   command.UpdateOpportunityStage
}

func NewOpportunityHandler(
	listOpportunities query.ListOpportunities,
	getOpportunitySummary query.GetOpportunitySummary,
	getOpportunityByPublicID query.GetOpportunityByPublicID,
	createOpportunity command.CreateOpportunity,
	updateOpportunity command.UpdateOpportunityProfile,
	updateOpportunityStage command.UpdateOpportunityStage,
) OpportunityHandler {
	return OpportunityHandler{
		listOpportunities:        listOpportunities,
		getOpportunitySummary:    getOpportunitySummary,
		getOpportunityByPublicID: getOpportunityByPublicID,
		createOpportunity:        createOpportunity,
		updateOpportunity:        updateOpportunity,
		updateOpportunityStage:   updateOpportunityStage,
	}
}

func (handler OpportunityHandler) List(writer http.ResponseWriter, request *http.Request) {
	opportunities := handler.listOpportunities.Execute(query.OpportunityFilters{
		Stage:            request.URL.Query().Get("stage"),
		LeadPublicID:     request.URL.Query().Get("leadPublicId"),
		CustomerPublicID: request.URL.Query().Get("customerPublicId"),
		SaleType:         request.URL.Query().Get("saleType"),
		OwnerUserID:      request.URL.Query().Get("ownerUserId"),
		Search:           request.URL.Query().Get("q"),
	})
	response := make([]dto.OpportunityResponse, 0, len(opportunities))

	for _, opportunity := range opportunities {
		response = append(response, mapOpportunity(opportunity))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler OpportunityHandler) Summary(writer http.ResponseWriter, request *http.Request) {
	summary := handler.getOpportunitySummary.Execute(query.OpportunityFilters{
		Stage:            request.URL.Query().Get("stage"),
		LeadPublicID:     request.URL.Query().Get("leadPublicId"),
		CustomerPublicID: request.URL.Query().Get("customerPublicId"),
		SaleType:         request.URL.Query().Get("saleType"),
		OwnerUserID:      request.URL.Query().Get("ownerUserId"),
		Search:           request.URL.Query().Get("q"),
	})

	writeJSON(writer, http.StatusOK, dto.OpportunitySummaryResponse{
		Total:            summary.Total,
		TotalAmountCents: summary.TotalAmountCents,
		ByStage:          summary.ByStage,
	})
}

func (handler OpportunityHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	opportunity := handler.getOpportunityByPublicID.Execute(request.PathValue("publicId"))
	if opportunity == nil {
		writeError(writer, http.StatusNotFound, "opportunity_not_found", "Opportunity was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapOpportunity(*opportunity))
}

func (handler OpportunityHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateOpportunityRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createOpportunity.Execute(command.CreateOpportunityInput{
		LeadPublicID:     payload.LeadPublicID,
		CustomerPublicID: payload.CustomerPublicID,
		Title:            payload.Title,
		SaleType:         payload.SaleType,
		OwnerUserID:      payload.OwnerUserID,
		AmountCents:      payload.AmountCents,
	})

	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}

	writeJSON(writer, http.StatusCreated, mapOpportunity(*result.Opportunity))
}

func (handler OpportunityHandler) Update(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateOpportunityRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.updateOpportunity.Execute(command.UpdateOpportunityProfileInput{
		PublicID:         request.PathValue("publicId"),
		CustomerPublicID: payload.CustomerPublicID,
		Title:            payload.Title,
		SaleType:         payload.SaleType,
		OwnerUserID:      payload.OwnerUserID,
		AmountCents:      payload.AmountCents,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusOK, mapOpportunity(*result.Opportunity))
	}
}

func (handler OpportunityHandler) UpdateStage(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateOpportunityStageRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.updateOpportunityStage.Execute(command.UpdateOpportunityStageInput{
		PublicID: request.PathValue("publicId"),
		Stage:    payload.Stage,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusOK, mapOpportunity(*result.Opportunity))
	}
}
