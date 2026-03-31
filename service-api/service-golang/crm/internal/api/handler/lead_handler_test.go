package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/persistence"
)

func TestListShouldReturnBootstrapLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads", nil)
	recorder := httptest.NewRecorder()

	handler.List(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 lead, got %d", len(response))
	}

	if response[0].Email != "lead@bootstrap-ops.local" {
		t.Fatalf("expected bootstrap lead email, got %s", response[0].Email)
	}
}

func TestListShouldFilterByStatusSearchAndAssignment(t *testing.T) {
	repository := persistence.NewInMemoryLeadRepository()
	createLead := command.NewCreateLead(repository)
	updateLeadStatus := command.NewUpdateLeadStatus(repository)
	handler := newLeadHandlerForTest(repository)

	createdLead := createLead.Execute(command.CreateLeadInput{
		Name:        "Ana Souza",
		Email:       "ana@example.com",
		Source:      "Meta-Ads",
		OwnerUserID: "owner-ana",
	})
	if createdLead.Lead == nil {
		t.Fatalf("expected created lead, got error %s", createdLead.ErrorCode)
	}

	updatedLead := updateLeadStatus.Execute(command.UpdateLeadStatusInput{
		PublicID: createdLead.Lead.PublicID,
		Status:   "contacted",
	})
	if updatedLead.Lead == nil {
		t.Fatalf("expected updated lead, got error %s", updatedLead.ErrorCode)
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/crm/leads?status=contacted&q=ana&assigned=true&source=meta-ads&ownerUserId=owner-ana",
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.List(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 filtered lead, got %d", len(response))
	}

	if response[0].Email != "ana@example.com" {
		t.Fatalf("expected filtered lead ana@example.com, got %s", response[0].Email)
	}
}

func TestListShouldFilterUnassignedLeads(t *testing.T) {
	repository := persistence.NewInMemoryLeadRepository()
	createLead := command.NewCreateLead(repository)
	handler := newLeadHandlerForTest(repository)

	createdLead := createLead.Execute(command.CreateLeadInput{
		Name:   "Bruno Lima",
		Email:  "bruno@example.com",
		Source: "organic",
	})
	if createdLead.Lead == nil {
		t.Fatalf("expected created lead, got error %s", createdLead.ErrorCode)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads?assigned=false", nil)
	recorder := httptest.NewRecorder()

	handler.List(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response []dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 unassigned lead, got %d", len(response))
	}

	if response[0].Email != "bruno@example.com" {
		t.Fatalf("expected unassigned lead bruno@example.com, got %s", response[0].Email)
	}
}

func TestCreateShouldReturnCreatedLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"name":"Ana Souza","email":"ana@example.com","source":"meta-ads","ownerUserId":"owner-ana"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/crm/leads", body)
	recorder := httptest.NewRecorder()

	handler.Create(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Email != "ana@example.com" {
		t.Fatalf("expected normalized email, got %s", response.Email)
	}

	if response.Status != "captured" {
		t.Fatalf("expected status captured, got %s", response.Status)
	}
}

func TestCreateShouldReturnConflictForDuplicateEmail(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"name":"Duplicate Lead","email":"lead@bootstrap-ops.local","source":"manual","ownerUserId":"owner-duplicate"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/crm/leads", body)
	recorder := httptest.NewRecorder()

	handler.Create(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
}

func TestGetByPublicIDShouldReturnLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/lead-bootstrap-ops", nil)
	request.SetPathValue("publicId", "lead-bootstrap-ops")
	recorder := httptest.NewRecorder()

	handler.GetByPublicID(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestUpdateStatusShouldReturnUpdatedLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"status":"contacted"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/crm/leads/lead-bootstrap-ops/status", body)
	request.SetPathValue("publicId", "lead-bootstrap-ops")
	recorder := httptest.NewRecorder()

	handler.UpdateStatus(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Status != "contacted" {
		t.Fatalf("expected status contacted, got %s", response.Status)
	}
}

func newLeadHandlerForTest(repository *persistence.InMemoryLeadRepository) LeadHandler {
	return NewLeadHandler(
		query.NewListLeads(repository),
		query.NewGetLeadByPublicID(repository),
		command.NewCreateLead(repository),
		command.NewUpdateLeadStatus(repository),
	)
}
