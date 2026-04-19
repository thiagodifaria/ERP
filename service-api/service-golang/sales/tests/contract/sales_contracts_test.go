//go:build contract

package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/telemetry"
)

type readinessResponse struct {
	Service      string               `json:"service"`
	Status       string               `json:"status"`
	Dependencies []dependencyResponse `json:"dependencies"`
}

type dependencyResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type opportunityResponse struct {
	PublicID     string `json:"publicId"`
	LeadPublicID string `json:"leadPublicId"`
	Title        string `json:"title"`
	Stage        string `json:"stage"`
	OwnerUserID  string `json:"ownerUserId"`
	AmountCents  int64  `json:"amountCents"`
}

type proposalResponse struct {
	PublicID            string `json:"publicId"`
	OpportunityPublicID string `json:"opportunityPublicId"`
	Title               string `json:"title"`
	Status              string `json:"status"`
	AmountCents         int64  `json:"amountCents"`
}

type saleResponse struct {
	PublicID            string `json:"publicId"`
	OpportunityPublicID string `json:"opportunityPublicId"`
	ProposalPublicID    string `json:"proposalPublicId"`
	Status              string `json:"status"`
	AmountCents         int64  `json:"amountCents"`
}

type invoiceResponse struct {
	PublicID     string `json:"publicId"`
	SalePublicID string `json:"salePublicId"`
	Number       string `json:"number"`
	Status       string `json:"status"`
	AmountCents  int64  `json:"amountCents"`
	DueDate      string `json:"dueDate"`
	PaidAt       string `json:"paidAt"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestOpportunityListContractShouldExposePublicFields(t *testing.T) {
	response := performRequest(t, newContractRouter(), http.MethodGet, "/api/sales/opportunities", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload []opportunityResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(payload) == 0 {
		t.Fatalf("expected at least one opportunity")
	}
}

func TestCreateOpportunityContractShouldReturnCreatedResource(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodPost,
		"/api/sales/opportunities",
		bytes.NewBufferString(`{"leadPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000011","title":"Contract Expansion","ownerUserId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000012","amountCents":88000}`),
	)
	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	var payload opportunityResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if payload.Stage != "qualified" {
		t.Fatalf("expected stage qualified, got %s", payload.Stage)
	}
}

func TestProposalConversionContractShouldReturnSale(t *testing.T) {
	router := newContractRouter()

	createOpportunity := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/sales/opportunities",
		bytes.NewBufferString(`{"leadPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000021","title":"Contract Expansion","ownerUserId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000022","amountCents":91000}`),
	)
	if createOpportunity.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createOpportunity.Code)
	}

	var opportunity opportunityResponse
	_ = json.Unmarshal(createOpportunity.Body.Bytes(), &opportunity)

	createProposal := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/sales/opportunities/"+opportunity.PublicID+"/proposals",
		bytes.NewBufferString(`{"title":"Contract Proposal","amountCents":91000}`),
	)
	if createProposal.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createProposal.Code)
	}

	var proposal proposalResponse
	_ = json.Unmarshal(createProposal.Body.Bytes(), &proposal)

	updateProposal := performRequest(
		t,
		router,
		http.MethodPatch,
		"/api/sales/proposals/"+proposal.PublicID+"/status",
		bytes.NewBufferString(`{"status":"sent"}`),
	)
	if updateProposal.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateProposal.Code)
	}

	convertProposal := performRequest(t, router, http.MethodPost, "/api/sales/proposals/"+proposal.PublicID+"/convert", nil)
	if convertProposal.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, convertProposal.Code)
	}

	var sale saleResponse
	if err := json.Unmarshal(convertProposal.Body.Bytes(), &sale); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if sale.Status != "active" {
		t.Fatalf("expected status active, got %s", sale.Status)
	}
}

func TestInvoiceLifecycleContractShouldExposeBillableSaleShape(t *testing.T) {
	router := newContractRouter()

	createInvoice := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/sales/sales/"+persistence.BootstrapSalePublicID+"/invoice",
		bytes.NewBufferString(`{"number":"BOOTSTRAP-INV-0002","dueDate":"2026-05-25"}`),
	)
	if createInvoice.Code != http.StatusConflict {
		t.Fatalf("expected bootstrap sale to reject duplicate invoice, got %d", createInvoice.Code)
	}

	createOpportunity := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/sales/opportunities",
		bytes.NewBufferString(`{"leadPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000031","title":"Contract Billing","ownerUserId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000032","amountCents":121000}`),
	)
	if createOpportunity.Code != http.StatusCreated {
		t.Fatalf("expected opportunity create success, got %d", createOpportunity.Code)
	}

	var opportunity opportunityResponse
	_ = json.Unmarshal(createOpportunity.Body.Bytes(), &opportunity)

	createProposal := performRequest(
		t,
		router,
		http.MethodPost,
		"/api/sales/opportunities/"+opportunity.PublicID+"/proposals",
		bytes.NewBufferString(`{"title":"Contract Billing Proposal","amountCents":121000}`),
	)
	if createProposal.Code != http.StatusCreated {
		t.Fatalf("expected proposal create success, got %d", createProposal.Code)
	}

	var proposal proposalResponse
	_ = json.Unmarshal(createProposal.Body.Bytes(), &proposal)

	performRequest(
		t,
		router,
		http.MethodPatch,
		"/api/sales/proposals/"+proposal.PublicID+"/status",
		bytes.NewBufferString(`{"status":"sent"}`),
	)
	convertProposal := performRequest(t, router, http.MethodPost, "/api/sales/proposals/"+proposal.PublicID+"/convert", nil)
	if convertProposal.Code != http.StatusCreated {
		t.Fatalf("expected convert success, got %d", convertProposal.Code)
	}

	var sale saleResponse
	_ = json.Unmarshal(convertProposal.Body.Bytes(), &sale)

	createInvoice = performRequest(
		t,
		router,
		http.MethodPost,
		"/api/sales/sales/"+sale.PublicID+"/invoice",
		bytes.NewBufferString(`{"number":"CONTRACT-INV-0001","dueDate":"2026-05-25"}`),
	)
	if createInvoice.Code != http.StatusCreated {
		t.Fatalf("expected invoice create success, got %d", createInvoice.Code)
	}

	var invoice invoiceResponse
	if err := json.Unmarshal(createInvoice.Body.Bytes(), &invoice); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if invoice.Status != "draft" || invoice.Number != "CONTRACT-INV-0001" {
		t.Fatalf("unexpected invoice payload: %+v", invoice)
	}

	payInvoice := performRequest(
		t,
		router,
		http.MethodPatch,
		"/api/sales/invoices/"+invoice.PublicID+"/status",
		bytes.NewBufferString(`{"status":"paid"}`),
	)
	if payInvoice.Code != http.StatusOK {
		t.Fatalf("expected paid status update success, got %d", payInvoice.Code)
	}

	var paid invoiceResponse
	if err := json.Unmarshal(payInvoice.Body.Bytes(), &paid); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if paid.PaidAt == "" {
		t.Fatalf("expected paidAt to be exposed")
	}
}

func TestHealthDetailsContractShouldExposeDependencyShape(t *testing.T) {
	response := performRequest(t, newContractRouter(), http.MethodGet, "/health/details", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload readinessResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if payload.Service != "sales" || payload.Status != "ready" {
		t.Fatalf("unexpected readiness payload: %+v", payload)
	}

	if len(payload.Dependencies) != 3 {
		t.Fatalf("expected 3 dependencies, got %d", len(payload.Dependencies))
	}
}

func TestErrorContractShouldExposeCodeAndMessage(t *testing.T) {
	response := performRequest(
		t,
		newContractRouter(),
		http.MethodPost,
		"/api/sales/opportunities",
		bytes.NewBufferString(`{"leadPublicId":"invalid","title":"","amountCents":0}`),
	)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	var payload errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if strings.TrimSpace(payload.Code) == "" || strings.TrimSpace(payload.Message) == "" {
		t.Fatalf("expected populated error payload")
	}
}

func performRequest(t *testing.T, router http.Handler, method string, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
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
		telemetry.New("sales-contract"),
		persistence.NewInMemoryOpportunityRepository(),
		persistence.NewInMemoryProposalRepository(),
		persistence.NewInMemorySaleRepository(),
		persistence.NewInMemoryInvoiceRepository(),
		"postgres",
	)
}
