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
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository(), persistence.NewInMemoryUploadSessionRepository())

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
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository(), persistence.NewInMemoryUploadSessionRepository())

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

	var accessLink dto.AccessLinkResponse
	if err := json.Unmarshal(accessLinkRecorder.Body.Bytes(), &accessLink); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	downloadRequest := httptest.NewRequest(http.MethodGet, accessLink.AccessURL, nil)
	downloadRecorder := httptest.NewRecorder()
	router.ServeHTTP(downloadRecorder, downloadRequest)

	if downloadRecorder.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, downloadRecorder.Code)
	}
	if downloadRecorder.Header().Get("Location") == "" {
		t.Fatalf("expected redirect location header")
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

func TestRouterShouldCreateAndListAttachmentVersions(t *testing.T) {
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository(), persistence.NewInMemoryUploadSessionRepository())

	createRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/attachments",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","ownerType":"rentals.contract","ownerPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e70000351","fileName":"contrato-v1.pdf","contentType":"application/pdf","storageKey":"rentals/contrato-v1.pdf","storageDriver":"manual","source":"rentals","uploadedBy":"ops-user","fileSizeBytes":4096,"checksumSha256":"V1"}`),
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
	if created.CurrentVersion != 1 || created.VersionCount != 1 {
		t.Fatalf("expected initial version counters, got %+v", created)
	}

	createVersionRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/attachments/"+created.PublicID+"/versions",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","fileName":"contrato-v2.pdf","contentType":"application/pdf","storageKey":"rentals/contrato-v2.pdf","storageDriver":"manual","source":"rentals","uploadedBy":"ops-user","fileSizeBytes":5120,"checksumSha256":"V2"}`),
	)
	createVersionRecorder := httptest.NewRecorder()
	router.ServeHTTP(createVersionRecorder, createVersionRequest)
	if createVersionRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createVersionRecorder.Code)
	}

	if !bytes.Contains(createVersionRecorder.Body.Bytes(), []byte(`"versionNumber":2`)) {
		t.Fatalf("expected second version payload, got %s", createVersionRecorder.Body.String())
	}
	if !bytes.Contains(createVersionRecorder.Body.Bytes(), []byte(`"currentVersion":2`)) {
		t.Fatalf("expected attachment currentVersion updated, got %s", createVersionRecorder.Body.String())
	}

	listVersionsRequest := httptest.NewRequest(http.MethodGet, "/api/documents/attachments/"+created.PublicID+"/versions?tenantSlug=bootstrap-ops", nil)
	listVersionsRecorder := httptest.NewRecorder()
	router.ServeHTTP(listVersionsRecorder, listVersionsRequest)
	if listVersionsRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listVersionsRecorder.Code)
	}

	var versions []dto.AttachmentVersionResponse
	if err := json.Unmarshal(listVersionsRecorder.Body.Bytes(), &versions); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if len(versions) != 2 || versions[0].VersionNumber != 2 || versions[1].VersionNumber != 1 {
		t.Fatalf("expected two ordered versions, got %+v", versions)
	}
}

func TestRouterShouldCreateAndCompleteUploadSession(t *testing.T) {
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository(), persistence.NewInMemoryUploadSessionRepository())

	createSessionRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/upload-sessions",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","ownerType":"crm.customer","ownerPublicId":"0195e7a0-7a9c-7c1f-8a44-4a6e70000342","fileName":"evidencia.png","contentType":"image/png","storageDriver":"local","source":"crm","requestedBy":"documents-smoke","visibility":"restricted","retentionDays":120}`),
	)
	createSessionRequest.Host = "documents.local"
	createSessionRecorder := httptest.NewRecorder()
	router.ServeHTTP(createSessionRecorder, createSessionRequest)

	if createSessionRecorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createSessionRecorder.Code)
	}

	var session dto.UploadSessionResponse
	if err := json.Unmarshal(createSessionRecorder.Body.Bytes(), &session); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}

	detailRequest := httptest.NewRequest(http.MethodGet, "/api/documents/upload-sessions/"+session.PublicID+"?tenantSlug=bootstrap-ops", nil)
	detailRecorder := httptest.NewRecorder()
	router.ServeHTTP(detailRecorder, detailRequest)

	if detailRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, detailRecorder.Code)
	}

	completeRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/documents/upload-sessions/"+session.PublicID+"/complete",
		bytes.NewBufferString(`{"tenantSlug":"bootstrap-ops","uploadedBy":"documents-smoke","fileSizeBytes":10240,"checksumSha256":"AA11BB22"}`),
	)
	completeRecorder := httptest.NewRecorder()
	router.ServeHTTP(completeRecorder, completeRequest)

	if completeRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, completeRecorder.Code)
	}

	if !bytes.Contains(completeRecorder.Body.Bytes(), []byte(`"status":"completed"`)) {
		t.Fatalf("expected completed upload session payload, got %s", completeRecorder.Body.String())
	}
	if !bytes.Contains(completeRecorder.Body.Bytes(), []byte(`"fileSizeBytes":10240`)) {
		t.Fatalf("expected attachment payload in completion response, got %s", completeRecorder.Body.String())
	}
}

func TestRouterShouldExposeStorageCapabilitiesAndRuntimeDetails(t *testing.T) {
	router := NewRouter(telemetry.New("documents-test"), persistence.NewInMemoryAttachmentRepository(), persistence.NewInMemoryUploadSessionRepository())

	capabilitiesRequest := httptest.NewRequest(http.MethodGet, "/api/documents/storage/capabilities", nil)
	capabilitiesRecorder := httptest.NewRecorder()
	router.ServeHTTP(capabilitiesRecorder, capabilitiesRequest)

	if capabilitiesRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, capabilitiesRecorder.Code)
	}

	if !bytes.Contains(capabilitiesRecorder.Body.Bytes(), []byte(`"provider":"local"`)) {
		t.Fatalf("expected local storage capability, got %s", capabilitiesRecorder.Body.String())
	}
	if !bytes.Contains(capabilitiesRecorder.Body.Bytes(), []byte(`"provider":"s3_compatible"`)) {
		t.Fatalf("expected s3-compatible capability, got %s", capabilitiesRecorder.Body.String())
	}

	detailRequest := httptest.NewRequest(http.MethodGet, "/health/details", nil)
	detailRecorder := httptest.NewRecorder()
	router.ServeHTTP(detailRecorder, detailRequest)

	if detailRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, detailRecorder.Code)
	}

	if !bytes.Contains(detailRecorder.Body.Bytes(), []byte(`"name":"storage:local"`)) {
		t.Fatalf("expected storage dependency in readiness details, got %s", detailRecorder.Body.String())
	}
}
