package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/infrastructure/integration"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/telemetry"
)

func TestRouterShouldCreateAdjustTerminateAndAttachContract(t *testing.T) {
	router := NewRouter(telemetry.New("rentals-test"), persistence.NewInMemoryContractRepository(), integration.NewInMemoryDocumentsGateway())

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/rentals/contracts",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","customerPublicId":"0196445e-8d3a-78ae-8c74-bc7c4a400022","title":"Contrato Torre Leste","propertyCode":"TL-01","currencyCode":"BRL","amountCents":150000,"billingDay":5,"startsAt":"2026-05-01","endsAt":"2026-07-31","recordedBy":"ops-user"}`),
	)
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, createRequest)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
	}

	var created ContractResponse
	if err := json.Unmarshal(createRecorder.Body.Bytes(), &created); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	adjustRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/rentals/contracts/"+created.PublicID+"/adjustments",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","effectiveAt":"2026-06-01","newAmountCents":165000,"reason":"Reajuste anual","recordedBy":"ops-user"}`),
	)
	adjustRecorder := httptest.NewRecorder()
	router.ServeHTTP(adjustRecorder, adjustRequest)
	if adjustRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, adjustRecorder.Code)
	}

	attachmentRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/rentals/contracts/"+created.PublicID+"/attachments?tenantSlug=bootstrap-ops",
		bytes.NewBufferString(`{"fileName":"contrato.pdf","contentType":"application/pdf","storageKey":"rentals/contrato.pdf","storageDriver":"manual","source":"rentals","uploadedBy":"ops-user"}`),
	)
	attachmentRecorder := httptest.NewRecorder()
	router.ServeHTTP(attachmentRecorder, attachmentRequest)
	if attachmentRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, attachmentRecorder.Code)
	}

	terminateRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/rentals/contracts/"+created.PublicID+"/terminate",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","effectiveAt":"2026-06-15","reason":"Rescisao amigavel","recordedBy":"ops-user"}`),
	)
	terminateRecorder := httptest.NewRecorder()
	router.ServeHTTP(terminateRecorder, terminateRequest)
	if terminateRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, terminateRecorder.Code)
	}

	summaryRequest := httptest.NewRequest(http.MethodGet, "/api/rentals/contracts/summary?tenantSlug=bootstrap-ops", nil)
	summaryRecorder := httptest.NewRecorder()
	router.ServeHTTP(summaryRecorder, summaryRequest)
	if summaryRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, summaryRecorder.Code)
	}

	chargesRequest := httptest.NewRequest(http.MethodGet, "/api/rentals/contracts/"+created.PublicID+"/charges?tenantSlug=bootstrap-ops", nil)
	chargesRecorder := httptest.NewRecorder()
	router.ServeHTTP(chargesRecorder, chargesRequest)
	if chargesRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, chargesRecorder.Code)
	}

	var charges []ChargeResponse
	if err := json.Unmarshal(chargesRecorder.Body.Bytes(), &charges); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if len(charges) == 0 {
		t.Fatalf("expected charges to be generated")
	}

	payChargeRequest := httptest.NewRequest(
		http.MethodPatch,
		"/api/rentals/contracts/"+created.PublicID+"/charges/"+charges[0].PublicID+"/status",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","status":"paid","recordedBy":"ops-user","paidAt":"2026-05-06T10:30:00Z","paymentReference":"pix-rent-001"}`),
	)
	payChargeRecorder := httptest.NewRecorder()
	router.ServeHTTP(payChargeRecorder, payChargeRequest)
	if payChargeRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, payChargeRecorder.Code)
	}

	historyRequest := httptest.NewRequest(http.MethodGet, "/api/rentals/contracts/"+created.PublicID+"/history?tenantSlug=bootstrap-ops", nil)
	historyRecorder := httptest.NewRecorder()
	router.ServeHTTP(historyRecorder, historyRequest)
	if historyRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, historyRecorder.Code)
	}

	if !bytes.Contains(historyRecorder.Body.Bytes(), []byte(`"eventCode":"charge_paid"`)) {
		t.Fatalf("expected history to include charge_paid event")
	}
}
