//go:build contract

package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

type leadResponse struct {
	PublicID    string `json:"publicId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Source      string `json:"source"`
	Status      string `json:"status"`
	OwnerUserID string `json:"ownerUserId"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type readinessResponse struct {
	Service      string               `json:"service"`
	Status       string               `json:"status"`
	Dependencies []dependencyResponse `json:"dependencies"`
}

type leadNoteResponse struct {
	PublicID     string `json:"publicId"`
	LeadPublicID string `json:"leadPublicId"`
	Body         string `json:"body"`
	Category     string `json:"category"`
}

type relationshipEventResponse struct {
	PublicID          string `json:"publicId"`
	AggregateType     string `json:"aggregateType"`
	AggregatePublicID string `json:"aggregatePublicId"`
	EventCode         string `json:"eventCode"`
	Actor             string `json:"actor"`
	Summary           string `json:"summary"`
}

type attachmentResponse struct {
	PublicID      string `json:"publicId"`
	OwnerType     string `json:"ownerType"`
	OwnerPublicID string `json:"ownerPublicId"`
	FileName      string `json:"fileName"`
}

type outboxEventResponse struct {
	PublicID          string `json:"publicId"`
	AggregateType     string `json:"aggregateType"`
	AggregatePublicID string `json:"aggregatePublicId"`
	EventType         string `json:"eventType"`
}

type dependencyResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func TestLeadListContractShouldExposePublicFields(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodGet,
		"/api/crm/leads",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload []leadResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(payload) == 0 {
		t.Fatalf("expected at least one lead")
	}

	for _, lead := range payload {
		if strings.TrimSpace(lead.PublicID) == "" {
			t.Fatalf("expected public id")
		}

		if strings.TrimSpace(lead.Name) == "" {
			t.Fatalf("expected name")
		}

		if strings.TrimSpace(lead.Email) == "" {
			t.Fatalf("expected email")
		}

		if strings.TrimSpace(lead.Source) == "" {
			t.Fatalf("expected source")
		}

		if strings.TrimSpace(lead.Status) == "" {
			t.Fatalf("expected status")
		}
	}
}

func TestCreateLeadContractShouldReturnCreatedResource(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodPost,
		"/api/crm/leads",
		bytes.NewBufferString(`{"name":"Contract Lead","email":"CONTRACT.LEAD@EXAMPLE.COM","source":"meta-ads"}`),
	)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var payload leadResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if payload.Name != "Contract Lead" {
		t.Fatalf("expected name Contract Lead, got %s", payload.Name)
	}

	if payload.Email != "contract.lead@example.com" {
		t.Fatalf("expected normalized email, got %s", payload.Email)
	}

	if payload.Source != "meta-ads" {
		t.Fatalf("expected source meta-ads, got %s", payload.Source)
	}

	if payload.Status != "captured" {
		t.Fatalf("expected status captured, got %s", payload.Status)
	}
}

func TestLeadStatusContractShouldReturnUpdatedResource(t *testing.T) {
	router := newContractRouter()

	createResponse := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/crm/leads",
		bytes.NewBufferString(`{"name":"Contract Status","email":"contract.status@example.com","source":"organic"}`),
	)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createResponse.Code)
	}

	var created leadResponse
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	updateResponse := performRequest(
		t,
		router,
		http.MethodPatch,
		"/api/crm/leads/"+created.PublicID+"/status",
		bytes.NewBufferString(`{"status":"contacted"}`),
	)

	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateResponse.Code)
	}

	var updated leadResponse
	if err := json.Unmarshal(updateResponse.Body.Bytes(), &updated); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if updated.PublicID != created.PublicID {
		t.Fatalf("expected same public id, got %s", updated.PublicID)
	}

	if updated.Status != "contacted" {
		t.Fatalf("expected status contacted, got %s", updated.Status)
	}
}

func TestLeadProfileContractShouldReturnUpdatedResource(t *testing.T) {
	router := newContractRouter()

	createResponse := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/crm/leads",
		bytes.NewBufferString(`{"name":"Contract Profile","email":"contract.profile@example.com","source":"organic"}`),
	)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createResponse.Code)
	}

	var created leadResponse
	if err := json.Unmarshal(createResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	updateResponse := performRequest(
		t,
		router,
		http.MethodPatch,
		"/api/crm/leads/"+created.PublicID,
		bytes.NewBufferString(`{"name":"Contract Profile Prime","email":"CONTRACT.PROFILE.PRIME@EXAMPLE.COM","source":"instagram"}`),
	)

	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateResponse.Code)
	}

	var updated leadResponse
	if err := json.Unmarshal(updateResponse.Body.Bytes(), &updated); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if updated.Name != "Contract Profile Prime" {
		t.Fatalf("expected updated name Contract Profile Prime, got %s", updated.Name)
	}

	if updated.Email != "contract.profile.prime@example.com" {
		t.Fatalf("expected normalized updated email, got %s", updated.Email)
	}

	if updated.Source != "instagram" {
		t.Fatalf("expected updated source instagram, got %s", updated.Source)
	}
}

func TestLeadNotesContractShouldExposeBootstrapHistory(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodGet,
		"/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/notes",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload []leadNoteResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(payload) == 0 {
		t.Fatalf("expected at least one lead note")
	}

	if payload[0].LeadPublicID != persistence.BootstrapLeadPublicID {
		t.Fatalf("expected bootstrap lead public id, got %s", payload[0].LeadPublicID)
	}
}

func TestCreateLeadNoteContractShouldReturnCreatedResource(t *testing.T) {
	router := newContractRouter()

	createLeadResponse := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/crm/leads",
		bytes.NewBufferString(`{"name":"Contract Note","email":"contract.note@example.com","source":"whatsapp"}`),
	)

	if createLeadResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createLeadResponse.Code)
	}

	var leadPayload leadResponse
	if err := json.Unmarshal(createLeadResponse.Body.Bytes(), &leadPayload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	createNoteResponse := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/crm/leads/"+leadPayload.PublicID+"/notes",
		bytes.NewBufferString(`{"body":"Cliente pediu comparativo com plano premium.","category":"follow-up"}`),
	)

	if createNoteResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createNoteResponse.Code)
	}

	var payload leadNoteResponse
	if err := json.Unmarshal(createNoteResponse.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if payload.Category != "follow-up" {
		t.Fatalf("expected category follow-up, got %s", payload.Category)
	}

	if payload.LeadPublicID != leadPayload.PublicID {
		t.Fatalf("expected created note lead public id %s, got %s", leadPayload.PublicID, payload.LeadPublicID)
	}
}

func TestLeadHistoryContractShouldExposeRelationshipEvents(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodGet,
		"/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/history",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload []relationshipEventResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(payload) == 0 {
		t.Fatalf("expected at least one relationship event")
	}
}

func TestLeadAttachmentContractShouldReturnCreatedResource(t *testing.T) {
	router := newContractRouter()

	response := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/attachments",
		bytes.NewBufferString(`{"fileName":"crm-brief.pdf","contentType":"application/pdf","storageKey":"crm/crm-brief.pdf","storageDriver":"manual","source":"crm","uploadedBy":"contract-test"}`),
	)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var payload attachmentResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if payload.OwnerType != "crm.lead" {
		t.Fatalf("expected owner type crm.lead, got %s", payload.OwnerType)
	}
}

func TestPendingOutboxContractShouldExposePublicShape(t *testing.T) {
	router := newContractRouter()

	createLeadResponse := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/crm/leads",
		bytes.NewBufferString(`{"name":"Contract Outbox","email":"contract.outbox@example.com","source":"manual"}`),
	)

	if createLeadResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createLeadResponse.Code)
	}

	response := performRequest(
		t,
		router,
		http.MethodGet,
		"/api/crm/outbox/pending",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload []outboxEventResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(payload) == 0 {
		t.Fatalf("expected at least one outbox event")
	}

	if strings.TrimSpace(payload[0].EventType) == "" {
		t.Fatalf("expected event type")
	}
}

func TestHealthDetailsContractShouldExposeDependencyShape(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodGet,
		"/health/details",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload readinessResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if payload.Service != "crm" {
		t.Fatalf("expected service crm, got %s", payload.Service)
	}

	if payload.Status != "ready" {
		t.Fatalf("expected status ready, got %s", payload.Status)
	}

	if len(payload.Dependencies) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(payload.Dependencies))
	}

	postgresql := payload.Dependencies[1]
	if postgresql.Name != "postgresql" || postgresql.Status != "ready" {
		t.Fatalf("expected postgresql dependency ready, got %s/%s", postgresql.Name, postgresql.Status)
	}
}

func TestErrorContractShouldExposeCodeAndMessage(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodPatch,
		"/api/crm/leads/"+persistence.BootstrapLeadPublicID+"/owner",
		bytes.NewBufferString(`{"ownerUserId":"invalid-owner"}`),
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	var payload errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if strings.TrimSpace(payload.Code) == "" {
		t.Fatalf("expected error code")
	}

	if strings.TrimSpace(payload.Message) == "" {
		t.Fatalf("expected error message")
	}
}

func performRequest(
	t *testing.T,
	router http.Handler,
	method string,
	path string,
	body *bytes.Buffer,
) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *bytes.Buffer
	if body == nil {
		requestBody = bytes.NewBuffer(nil)
	} else {
		requestBody = body
	}

	request := httptest.NewRequest(method, path, requestBody)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func newContractRouter() http.Handler {
	return api.NewRouterWithRuntime(
		telemetry.New("crm-contract"),
		persistence.NewInMemoryTenantRepositoryFactory("bootstrap-ops"),
		integration.NewInMemoryDocumentsGateway(),
		"postgres",
	)
}
