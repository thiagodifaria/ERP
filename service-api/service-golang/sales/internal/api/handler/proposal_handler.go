// ProposalHandler expoe a trilha publica de propostas do Sales MVP.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
)

type ProposalHandler struct {
	listProposals         query.ListProposalsByOpportunity
	getProposalByPublicID query.GetProposalByPublicID
	createProposal        command.CreateProposal
	updateProposalStatus  command.UpdateProposalStatus
	convertProposal       command.ConvertProposalToSale
}

func NewProposalHandler(
	listProposals query.ListProposalsByOpportunity,
	getProposalByPublicID query.GetProposalByPublicID,
	createProposal command.CreateProposal,
	updateProposalStatus command.UpdateProposalStatus,
	convertProposal command.ConvertProposalToSale,
) ProposalHandler {
	return ProposalHandler{
		listProposals:         listProposals,
		getProposalByPublicID: getProposalByPublicID,
		createProposal:        createProposal,
		updateProposalStatus:  updateProposalStatus,
		convertProposal:       convertProposal,
	}
}

func (handler ProposalHandler) ListByOpportunity(writer http.ResponseWriter, request *http.Request) {
	proposals := handler.listProposals.Execute(request.PathValue("publicId"))
	response := make([]dto.ProposalResponse, 0, len(proposals))

	for _, proposal := range proposals {
		response = append(response, mapProposal(proposal))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler ProposalHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateProposalRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createProposal.Execute(command.CreateProposalInput{
		OpportunityPublicID: request.PathValue("publicId"),
		Title:               payload.Title,
		AmountCents:         payload.AmountCents,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusCreated, mapProposal(*result.Proposal))
	}
}

func (handler ProposalHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	proposal := handler.getProposalByPublicID.Execute(request.PathValue("publicId"))
	if proposal == nil {
		writeError(writer, http.StatusNotFound, "proposal_not_found", "Proposal was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapProposal(*proposal))
}

func (handler ProposalHandler) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateProposalStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.updateProposalStatus.Execute(command.UpdateProposalStatusInput{
		PublicID: request.PathValue("publicId"),
		Status:   payload.Status,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusOK, mapProposal(*result.Proposal))
	}
}

func (handler ProposalHandler) Convert(writer http.ResponseWriter, request *http.Request) {
	result := handler.convertProposal.Execute(command.ConvertProposalToSaleInput{
		ProposalPublicID: request.PathValue("publicId"),
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	case result.Conflict:
		writeError(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusCreated, mapSale(*result.Sale))
	}
}
