package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/persistence"
)

func TestLeadNoteListShouldReturnBootstrapNote(t *testing.T) {
	leadRepository := persistence.NewInMemoryLeadRepository()
	leadNoteRepository := persistence.NewInMemoryLeadNoteRepository()
	handler := NewLeadNoteHandler(
		query.NewGetLeadByPublicID(leadRepository),
		query.NewListLeadNotes(leadNoteRepository),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/notes", nil)
	request.SetPathValue("publicId", persistence.BootstrapLeadPublicID)
	recorder := httptest.NewRecorder()

	handler.List(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []dto.LeadNoteResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 note, got %d", len(response))
	}

	if response[0].Category != "qualification" {
		t.Fatalf("expected category qualification, got %s", response[0].Category)
	}
}

func TestLeadNoteListShouldReturnNotFoundForUnknownLead(t *testing.T) {
	leadRepository := persistence.NewInMemoryLeadRepository()
	leadNoteRepository := persistence.NewInMemoryLeadNoteRepository()
	handler := NewLeadNoteHandler(
		query.NewGetLeadByPublicID(leadRepository),
		query.NewListLeadNotes(leadNoteRepository),
	)
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/missing/notes", nil)
	request.SetPathValue("publicId", "0195e7a0-7a9c-7c1f-8a44-4a6e70009999")
	recorder := httptest.NewRecorder()

	handler.List(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}
