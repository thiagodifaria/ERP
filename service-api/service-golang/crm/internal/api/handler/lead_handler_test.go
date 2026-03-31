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
  repository := persistence.NewInMemoryLeadRepository()
  handler := NewLeadHandler(
    query.NewListLeads(repository),
    query.NewGetLeadByPublicID(repository),
    command.NewCreateLead(repository),
    command.NewUpdateLeadStatus(repository),
  )
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

func TestCreateShouldReturnCreatedLead(t *testing.T) {
  repository := persistence.NewInMemoryLeadRepository()
  handler := NewLeadHandler(
    query.NewListLeads(repository),
    query.NewGetLeadByPublicID(repository),
    command.NewCreateLead(repository),
    command.NewUpdateLeadStatus(repository),
  )
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
  repository := persistence.NewInMemoryLeadRepository()
  handler := NewLeadHandler(
    query.NewListLeads(repository),
    query.NewGetLeadByPublicID(repository),
    command.NewCreateLead(repository),
    command.NewUpdateLeadStatus(repository),
  )
  body := bytes.NewBufferString(`{"name":"Duplicate Lead","email":"lead@bootstrap-ops.local","source":"manual","ownerUserId":"owner-duplicate"}`)
  request := httptest.NewRequest(http.MethodPost, "/api/crm/leads", body)
  recorder := httptest.NewRecorder()

  handler.Create(recorder, request)

  if recorder.Code != http.StatusConflict {
    t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
  }
}

func TestGetByPublicIDShouldReturnLead(t *testing.T) {
  repository := persistence.NewInMemoryLeadRepository()
  handler := NewLeadHandler(
    query.NewListLeads(repository),
    query.NewGetLeadByPublicID(repository),
    command.NewCreateLead(repository),
    command.NewUpdateLeadStatus(repository),
  )
  request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/lead-bootstrap-ops", nil)
  request.SetPathValue("publicId", "lead-bootstrap-ops")
  recorder := httptest.NewRecorder()

  handler.GetByPublicID(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
  }
}

func TestUpdateStatusShouldReturnUpdatedLead(t *testing.T) {
  repository := persistence.NewInMemoryLeadRepository()
  handler := NewLeadHandler(
    query.NewListLeads(repository),
    query.NewGetLeadByPublicID(repository),
    command.NewCreateLead(repository),
    command.NewUpdateLeadStatus(repository),
  )
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
