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

const (
	testOwnerAnaPublicID       = "0195e7a0-7a9c-7c1f-8a44-4a6e70000021"
	testOwnerCarolPublicID     = "0195e7a0-7a9c-7c1f-8a44-4a6e70000022"
	testOwnerDuplicatePublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e70000023"
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
		OwnerUserID: testOwnerAnaPublicID,
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
		"/api/crm/leads?status=contacted&q=ana&assigned=true&source=meta-ads&ownerUserId="+testOwnerAnaPublicID,
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

func TestSummaryShouldReturnPipelineSnapshot(t *testing.T) {
	repository := persistence.NewInMemoryLeadRepository()
	createLead := command.NewCreateLead(repository)
	updateLeadStatus := command.NewUpdateLeadStatus(repository)
	handler := newLeadHandlerForTest(repository)

	ana := createLead.Execute(command.CreateLeadInput{
		Name:        "Ana Souza",
		Email:       "ana@example.com",
		Source:      "meta-ads",
		OwnerUserID: testOwnerAnaPublicID,
	})
	bruno := createLead.Execute(command.CreateLeadInput{
		Name:   "Bruno Lima",
		Email:  "bruno@example.com",
		Source: "organic",
	})
	if ana.Lead == nil || bruno.Lead == nil {
		t.Fatalf("expected created leads for summary setup")
	}

	contacted := updateLeadStatus.Execute(command.UpdateLeadStatusInput{
		PublicID: ana.Lead.PublicID,
		Status:   "contacted",
	})
	qualified := updateLeadStatus.Execute(command.UpdateLeadStatusInput{
		PublicID: bruno.Lead.PublicID,
		Status:   "qualified",
	})
	if contacted.Lead == nil || qualified.Lead == nil {
		t.Fatalf("expected updated leads for summary setup")
	}

	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/summary", nil)
	recorder := httptest.NewRecorder()

	handler.Summary(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadSummaryResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Total != 3 {
		t.Fatalf("expected total 3, got %d", response.Total)
	}

	if response.Assigned != 2 || response.Unassigned != 1 {
		t.Fatalf("expected assigned/unassigned 2/1, got %d/%d", response.Assigned, response.Unassigned)
	}

	if response.ByStatus["captured"] != 1 || response.ByStatus["contacted"] != 1 || response.ByStatus["qualified"] != 1 {
		t.Fatalf("expected status buckets captured/contacted/qualified to be 1")
	}

	if response.BySource["manual"] != 1 || response.BySource["meta-ads"] != 1 || response.BySource["organic"] != 1 {
		t.Fatalf("expected source buckets manual/meta-ads/organic to be 1")
	}
}

func TestCreateShouldReturnCreatedLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"name":"Ana Souza","email":"ana@example.com","source":"meta-ads","ownerUserId":"` + testOwnerAnaPublicID + `"}`)
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

	if response.OwnerUserID != testOwnerAnaPublicID {
		t.Fatalf("expected owner %s, got %s", testOwnerAnaPublicID, response.OwnerUserID)
	}
}

func TestCreateShouldReturnConflictForDuplicateEmail(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"name":"Duplicate Lead","email":"lead@bootstrap-ops.local","source":"manual","ownerUserId":"` + testOwnerDuplicatePublicID + `"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/crm/leads", body)
	recorder := httptest.NewRecorder()

	handler.Create(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
}

func TestGetByPublicIDShouldReturnLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	request := httptest.NewRequest(http.MethodGet, "/api/crm/leads/"+persistence.BootstrapLeadPublicID, nil)
	request.SetPathValue("publicId", persistence.BootstrapLeadPublicID)
	recorder := httptest.NewRecorder()

	handler.GetByPublicID(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestUpdateStatusShouldReturnUpdatedLead(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"status":"contacted"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/status", body)
	request.SetPathValue("publicId", persistence.BootstrapLeadPublicID)
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

func TestUpdateOwnerShouldReturnUpdatedLead(t *testing.T) {
	repository := persistence.NewInMemoryLeadRepository()
	createLead := command.NewCreateLead(repository)
	handler := newLeadHandlerForTest(repository)

	createdLead := createLead.Execute(command.CreateLeadInput{
		Name:   "Carol Viana",
		Email:  "carol@example.com",
		Source: "manual",
	})
	if createdLead.Lead == nil {
		t.Fatalf("expected created lead, got error %s", createdLead.ErrorCode)
	}

	body := bytes.NewBufferString(`{"ownerUserId":"owner-carol"}`)
	body = bytes.NewBufferString(`{"ownerUserId":"` + testOwnerCarolPublicID + `"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/crm/leads/"+createdLead.Lead.PublicID+"/owner", body)
	request.SetPathValue("publicId", createdLead.Lead.PublicID)
	recorder := httptest.NewRecorder()

	handler.UpdateOwner(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.OwnerUserID != testOwnerCarolPublicID {
		t.Fatalf("expected owner %s, got %s", testOwnerCarolPublicID, response.OwnerUserID)
	}
}

func TestUpdateOwnerShouldAllowClearingOwner(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"ownerUserId":"   "}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/owner", body)
	request.SetPathValue("publicId", persistence.BootstrapLeadPublicID)
	recorder := httptest.NewRecorder()

	handler.UpdateOwner(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.LeadResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.OwnerUserID != "" {
		t.Fatalf("expected owner to be cleared, got %s", response.OwnerUserID)
	}
}

func TestUpdateOwnerShouldRejectInvalidUUID(t *testing.T) {
	handler := newLeadHandlerForTest(persistence.NewInMemoryLeadRepository())
	body := bytes.NewBufferString(`{"ownerUserId":"owner-carol"}`)
	request := httptest.NewRequest(http.MethodPatch, "/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/owner", body)
	request.SetPathValue("publicId", persistence.BootstrapLeadPublicID)
	recorder := httptest.NewRecorder()

	handler.UpdateOwner(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response dto.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Code != "invalid_owner_user_id" {
		t.Fatalf("expected invalid_owner_user_id, got %s", response.Code)
	}
}

func newLeadHandlerForTest(repository *persistence.InMemoryLeadRepository) LeadHandler {
	return NewLeadHandler(
		query.NewListLeads(repository),
		query.NewGetLeadPipelineSummary(repository),
		query.NewGetLeadByPublicID(repository),
		command.NewCreateLead(repository),
		command.NewUpdateLeadOwner(repository),
		command.NewUpdateLeadStatus(repository),
	)
}
