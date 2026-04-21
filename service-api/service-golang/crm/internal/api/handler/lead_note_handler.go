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

type LeadNoteHandler struct {
	repositories repository.TenantRepositoryFactory
}

func NewLeadNoteHandler(repositories repository.TenantRepositoryFactory) LeadNoteHandler {
	return LeadNoteHandler{repositories: repositories}
}

func (handler LeadNoteHandler) List(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request)
	if bundle.LeadRepository == nil {
		return
	}

	leadPublicID := request.PathValue("publicId")
	lead := query.NewGetLeadByPublicID(bundle.LeadRepository).Execute(leadPublicID)
	if lead == nil {
		writeNotFound(writer, "lead_not_found", "Lead was not found.")
		return
	}

	notes := query.NewListLeadNotes(bundle.LeadNoteRepository).Execute(leadPublicID)
	response := make([]dto.LeadNoteResponse, 0, len(notes))
	for _, note := range notes {
		response = append(response, mapLeadNote(note))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler LeadNoteHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateLeadNoteRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(writer, "invalid_json", "Request body is invalid.")
		return
	}

	bundle, _ := handler.resolveRepositories(writer, request)
	if bundle.LeadRepository == nil {
		return
	}

	result := command.NewCreateLeadNote(bundle.LeadRepository, bundle.LeadNoteRepository).Execute(command.CreateLeadNoteInput{
		LeadPublicID: request.PathValue("publicId"),
		Body:         payload.Body,
		Category:     payload.Category,
	})

	writer.Header().Set("Content-Type", "application/json")

	switch {
	case result.BadRequest:
		writeErrorResponse(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeErrorResponse(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writer.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(writer).Encode(mapLeadNote(*result.Note))
	}
}

func mapLeadNote(note entity.LeadNote) dto.LeadNoteResponse {
	return dto.LeadNoteResponse{
		PublicID:     note.PublicID,
		LeadPublicID: note.LeadPublicID,
		Body:         note.Body,
		Category:     note.Category,
		CreatedAt:    note.CreatedAt,
	}
}
