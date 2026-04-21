package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/infrastructure/persistence"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/telemetry"
)

func TestRouterShouldCreateAndListAttachments(t *testing.T) {
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository())

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/attachments",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","ownerType":"crm.lead","ownerPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e70000341","fileName":"proposta.pdf","contentType":"application/pdf","storageKey":"leads/proposta.pdf","storageDriver":"manual","source":"crm","uploadedBy":"ops-user"}`),
	)
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, createRequest)

	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/documents/attachments?tenantSlug=bootstrap-ops&ownerType=crm.lead&ownerPublicId=0195e7a0-7a9c-7c1f-8a44-4a6e70000341", nil)
	listRecorder := httptest.NewRecorder()
	router.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRecorder.Code)
	}

	var payload []dto.AttachmentResponse
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(payload) != 1 || payload[0].FileName != "proposta.pdf" {
		t.Fatalf("expected one persisted attachment, got %+v", payload)
	}
}
