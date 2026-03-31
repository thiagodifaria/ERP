package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type LeadNoteHandler struct {
	getLeadByPublicID query.GetLeadByPublicID
	listLeadNotes     query.ListLeadNotes
}

func NewLeadNoteHandler(
	getLeadByPublicID query.GetLeadByPublicID,
	listLeadNotes query.ListLeadNotes,
) LeadNoteHandler {
	return LeadNoteHandler{
		getLeadByPublicID: getLeadByPublicID,
		listLeadNotes:     listLeadNotes,
	}
}

func (handler LeadNoteHandler) List(writer http.ResponseWriter, request *http.Request) {
	leadPublicID := request.PathValue("publicId")
	lead := handler.getLeadByPublicID.Execute(leadPublicID)
	if lead == nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    "lead_not_found",
			Message: "Lead was not found.",
		})
		return
	}

	notes := handler.listLeadNotes.Execute(leadPublicID)
	response := make([]dto.LeadNoteResponse, 0, len(notes))
	for _, note := range notes {
		response = append(response, mapLeadNote(note))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
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
