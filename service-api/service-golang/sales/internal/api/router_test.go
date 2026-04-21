package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/telemetry"
)

func newTestRouter() http.Handler {
	return NewRouter(
		telemetry.New("sales-test"),
		persistence.NewInMemoryOpportunityRepository(),
		persistence.NewInMemoryProposalRepository(),
		persistence.NewInMemorySaleRepository(),
		persistence.NewInMemoryInvoiceRepository(),
		persistence.NewInMemoryInstallmentRepository(),
		persistence.NewInMemoryCommissionRepository(),
		persistence.NewInMemoryPendingItemRepository(),
		persistence.NewInMemoryRenegotiationRepository(),
		persistence.NewInMemoryCommercialEventRepository(),
		persistence.NewInMemoryOutboxEventRepository(),
	)
}

func TestRouterShouldExposeHealthDetails(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/details", nil)
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Header().Get("X-Correlation-Id") != "pending-correlation" {
		t.Fatalf("expected fallback correlation id header")
	}

	var response dto.ReadinessResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Service != "sales" {
		t.Fatalf("expected service sales, got %s", response.Service)
	}
}

func TestRouterShouldExposeOpportunitySummary(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/sales/opportunities/summary", nil)
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var response dto.OpportunitySummaryResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if response.Total != 1 {
		t.Fatalf("expected total 1, got %d", response.Total)
	}
}

func TestRouterShouldCreateOpportunityAndConvertProposal(t *testing.T) {
	router := newTestRouter()

	createOpportunityRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/sales/opportunities",
		bytes.NewBufferString(`{"leadPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000001","customerPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000003","title":"Runtime Expansion","saleType":"upsell","ownerUserId":"0195e7a0-7a9c-7c1f-8a44-4a6e73000002","amountCents":99000}`),
	)
	createOpportunityResponse := httptest.NewRecorder()
	router.ServeHTTP(createOpportunityResponse, createOpportunityRequest)

	if createOpportunityResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createOpportunityResponse.Code)
	}

	var createdOpportunity dto.OpportunityResponse
	if err := json.Unmarshal(createOpportunityResponse.Body.Bytes(), &createdOpportunity); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if createdOpportunity.CustomerPublicID != "0195e7a0-7a9c-7c1f-8a44-4a6e73000003" || createdOpportunity.SaleType != "upsell" {
		t.Fatalf("expected customer and sale type to be persisted in opportunity payload")
	}

	createProposalRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/sales/opportunities/"+createdOpportunity.PublicID+"/proposals",
		bytes.NewBufferString(`{"title":"Runtime Expansion Proposal","amountCents":99000}`),
	)
	createProposalResponse := httptest.NewRecorder()
	router.ServeHTTP(createProposalResponse, createProposalRequest)

	if createProposalResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createProposalResponse.Code)
	}

	var createdProposal dto.ProposalResponse
	if err := json.Unmarshal(createProposalResponse.Body.Bytes(), &createdProposal); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	updateProposalRequest := httptest.NewRequest(
		http.MethodPatch,
		"/api/sales/proposals/"+createdProposal.PublicID+"/status",
		bytes.NewBufferString(`{"status":"sent"}`),
	)
	updateProposalResponse := httptest.NewRecorder()
	router.ServeHTTP(updateProposalResponse, updateProposalRequest)

	if updateProposalResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateProposalResponse.Code)
	}

	convertRequest := httptest.NewRequest(http.MethodPost, "/api/sales/proposals/"+createdProposal.PublicID+"/convert", nil)
	convertResponse := httptest.NewRecorder()
	router.ServeHTTP(convertResponse, convertRequest)

	if convertResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, convertResponse.Code)
	}

	var createdSale dto.SaleResponse
	if err := json.Unmarshal(convertResponse.Body.Bytes(), &createdSale); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if createdSale.Status != "active" {
		t.Fatalf("expected sale status active, got %s", createdSale.Status)
	}

	if createdSale.CustomerPublicID != createdOpportunity.CustomerPublicID || createdSale.SaleType != "upsell" {
		t.Fatalf("expected sale to inherit customer and sale type from opportunity")
	}

	createInvoiceRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/sales/sales/"+createdSale.PublicID+"/invoice",
		bytes.NewBufferString(`{"number":"RUNTIME-INV-0001","dueDate":"2026-05-20"}`),
	)
	createInvoiceResponse := httptest.NewRecorder()
	router.ServeHTTP(createInvoiceResponse, createInvoiceRequest)

	if createInvoiceResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createInvoiceResponse.Code)
	}

	var createdInvoice dto.InvoiceResponse
	if err := json.Unmarshal(createInvoiceResponse.Body.Bytes(), &createdInvoice); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	updateInvoiceRequest := httptest.NewRequest(
		http.MethodPatch,
		"/api/sales/invoices/"+createdInvoice.PublicID+"/status",
		bytes.NewBufferString(`{"status":"paid"}`),
	)
	updateInvoiceResponse := httptest.NewRecorder()
	router.ServeHTTP(updateInvoiceResponse, updateInvoiceRequest)

	if updateInvoiceResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateInvoiceResponse.Code)
	}

	invoiceSummaryRequest := httptest.NewRequest(http.MethodGet, "/api/sales/invoices/summary", nil)
	invoiceSummaryResponse := httptest.NewRecorder()
	router.ServeHTTP(invoiceSummaryResponse, invoiceSummaryRequest)

	if invoiceSummaryResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, invoiceSummaryResponse.Code)
	}

	var summary dto.InvoiceSummaryResponse
	if err := json.Unmarshal(invoiceSummaryResponse.Body.Bytes(), &summary); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if summary.PaidAmountCents < createdInvoice.AmountCents {
		t.Fatalf("expected paid amount to include runtime invoice")
	}
}
