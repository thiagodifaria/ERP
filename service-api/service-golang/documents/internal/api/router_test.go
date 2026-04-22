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

func TestRouterShouldGetArchiveAndIssueAccessLinks(t *testing.T) {
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository())

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/attachments",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","ownerType":"crm.customer","ownerPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e70000342","fileName":"contrato.pdf","contentType":"application/pdf","storageKey":"customers/contrato.pdf","storageDriver":"r2","source":"crm","uploadedBy":"ops-user","fileSizeBytes":4096,"checksumSha256":"ABC123","visibility":"restricted","retentionDays":45}`),
	)
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, createRequest)

	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createRecorder.Code)
	}

	var created dto.AttachmentResponse
	if err := json.Unmarshal(createRecorder.Body.Bytes(), &created); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	detailRequest := httptest.NewRequest(http.MethodGet, "/api/documents/attachments/"+created.PublicID+"?tenantSlug=bootstrap-ops", nil)
	detailRecorder := httptest.NewRecorder()
	router.ServeHTTP(detailRecorder, detailRequest)

	if detailRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, detailRecorder.Code)
	}

	accessLinkRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/attachments/"+created.PublicID+"/access-links",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","expiresInSeconds":600}`),
	)
	accessLinkRequest.Host = "documents.local"
	accessLinkRecorder := httptest.NewRecorder()
	router.ServeHTTP(accessLinkRecorder, accessLinkRequest)

	if accessLinkRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, accessLinkRecorder.Code)
	}

	archiveRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/attachments/"+created.PublicID+"/archive",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","reason":"fim-do-contrato"}`),
	)
	archiveRecorder := httptest.NewRecorder()
	router.ServeHTTP(archiveRecorder, archiveRequest)

	if archiveRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, archiveRecorder.Code)
	}

	archivedListRequest := httptest.NewRequest(http.MethodGet, "/api/documents/attachments?tenantSlug=bootstrap-ops&archived=true", nil)
	archivedListRecorder := httptest.NewRecorder()
	router.ServeHTTP(archivedListRecorder, archivedListRequest)

	if archivedListRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, archivedListRecorder.Code)
	}

	var archivedPayload []dto.AttachmentResponse
	if err := json.Unmarshal(archivedListRecorder.Body.Bytes(), &archivedPayload); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	if len(archivedPayload) != 1 || archivedPayload[0].ArchivedAt == nil || archivedPayload[0].ArchiveReason != "fim-do-contrato" {
		t.Fatalf("expected archived attachment, got %+v", archivedPayload)
	}
}
