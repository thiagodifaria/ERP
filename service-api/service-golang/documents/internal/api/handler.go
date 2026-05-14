package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type AttachmentHandler struct {
	listAttachments       query.ListAttachments
	createAttachment      command.CreateAttachment
	getAttachment         query.GetAttachment
	archiveAttachment     command.ArchiveAttachment
	listVersions          query.ListAttachmentVersions
	createVersion         command.CreateAttachmentVersion
	createUploadSession   command.CreateUploadSession
	getUploadSession      query.GetUploadSession
	completeUploadSession command.CompleteUploadSession
	attachmentRepository  repository.AttachmentRepository
	accessTokenSecret     string
}

func NewAttachmentHandler(attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository, accessTokenSecret string) AttachmentHandler {
	return AttachmentHandler{
		listAttachments:       query.NewListAttachments(attachmentRepository),
		createAttachment:      command.NewCreateAttachment(attachmentRepository),
		getAttachment:         query.NewGetAttachment(attachmentRepository),
		archiveAttachment:     command.NewArchiveAttachment(attachmentRepository),
		listVersions:          query.NewListAttachmentVersions(attachmentRepository),
		createVersion:         command.NewCreateAttachmentVersion(attachmentRepository),
		createUploadSession:   command.NewCreateUploadSession(uploadSessionRepository),
		getUploadSession:      query.NewGetUploadSession(uploadSessionRepository),
		completeUploadSession: command.NewCompleteUploadSession(uploadSessionRepository, attachmentRepository),
		attachmentRepository:  attachmentRepository,
		accessTokenSecret:     strings.TrimSpace(accessTokenSecret),
	}
}

func (handler AttachmentHandler) List(writer http.ResponseWriter, request *http.Request) {
	attachments := handler.listAttachments.Execute(repository.AttachmentFilters{
		TenantSlug:    request.URL.Query().Get("tenantSlug"),
		OwnerType:     request.URL.Query().Get("ownerType"),
		OwnerPublicID: request.URL.Query().Get("ownerPublicId"),
		Source:        request.URL.Query().Get("source"),
		Visibility:    request.URL.Query().Get("visibility"),
		Archived:      request.URL.Query().Get("archived"),
	})

	response := make([]dto.AttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		response = append(response, mapAttachment(attachment))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler AttachmentHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateAttachmentRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createAttachment.Execute(command.CreateAttachmentInput{
		TenantSlug:     payload.TenantSlug,
		OwnerType:      payload.OwnerType,
		OwnerPublicID:  payload.OwnerPublicID,
		FileName:       payload.FileName,
		ContentType:    payload.ContentType,
		StorageKey:     payload.StorageKey,
		StorageDriver:  payload.StorageDriver,
		Source:         payload.Source,
		UploadedBy:     payload.UploadedBy,
		FileSizeBytes:  payload.FileSizeBytes,
		ChecksumSHA256: payload.ChecksumSHA256,
		Visibility:     payload.Visibility,
		RetentionDays:  payload.RetentionDays,
	})

	if result.BadRequest {
		writeDocumentsError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(mapAttachment(*result.Attachment))
}

func (handler AttachmentHandler) Get(writer http.ResponseWriter, request *http.Request) {
	attachment, ok := handler.getAttachment.Execute(request.URL.Query().Get("tenantSlug"), request.PathValue("publicId"))
	if !ok {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(mapAttachment(*attachment))
}

func (handler AttachmentHandler) Archive(writer http.ResponseWriter, request *http.Request) {
	var payload dto.ArchiveAttachmentRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.archiveAttachment.Execute(command.ArchiveAttachmentInput{
		TenantSlug: payload.TenantSlug,
		PublicID:   request.PathValue("publicId"),
		Reason:     payload.Reason,
	})
	if result.BadRequest {
		writeDocumentsError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(mapAttachment(*result.Attachment))
}

func (handler AttachmentHandler) ListVersions(writer http.ResponseWriter, request *http.Request) {
	versions := handler.listVersions.Execute(request.URL.Query().Get("tenantSlug"), request.PathValue("publicId"))
	response := make([]dto.AttachmentVersionResponse, 0, len(versions))
	for _, version := range versions {
		response = append(response, mapAttachmentVersion(version))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler AttachmentHandler) CreateVersion(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateAttachmentVersionRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createVersion.Execute(command.CreateAttachmentVersionInput{
		TenantSlug:         payload.TenantSlug,
		AttachmentPublicID: request.PathValue("publicId"),
		FileName:           payload.FileName,
		ContentType:        payload.ContentType,
		StorageKey:         payload.StorageKey,
		StorageDriver:      payload.StorageDriver,
		Source:             payload.Source,
		UploadedBy:         payload.UploadedBy,
		FileSizeBytes:      payload.FileSizeBytes,
		ChecksumSHA256:     payload.ChecksumSHA256,
	})
	if result.BadRequest {
		writeDocumentsError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(struct {
		Attachment dto.AttachmentResponse        `json:"attachment"`
		Version    dto.AttachmentVersionResponse `json:"version"`
	}{
		Attachment: mapAttachment(*result.Attachment),
		Version:    mapAttachmentVersion(*result.Version),
	})
}

func (handler AttachmentHandler) CreateAccessLink(writer http.ResponseWriter, request *http.Request) {
	var payload dto.AccessLinkRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	attachment, ok := handler.getAttachment.Execute(payload.TenantSlug, request.PathValue("publicId"))
	if !ok {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}
	if attachment.ArchivedAt != nil {
		writeDocumentsError(writer, http.StatusConflict, "attachment_archived", "Attachment is archived and cannot issue new access links.")
		return
	}

	expiresInSeconds := payload.ExpiresInSeconds
	if expiresInSeconds <= 0 {
		expiresInSeconds = 900
	}
	expiresAt := time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)
	scheme := strings.TrimSpace(request.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = "http"
		if request.TLS != nil {
			scheme = "https"
		}
	}
	host := strings.TrimSpace(request.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = request.Host
	}
	if host == "" {
		host = "documents.local"
	}
	baseURL := scheme + "://" + strings.TrimRight(host, "/")
	token := handler.newSignedAccessToken(attachment.PublicID, attachment.TenantSlug, expiresAt)

	response := dto.AccessLinkResponse{
		AttachmentPublicID: attachment.PublicID,
		TenantSlug:         attachment.TenantSlug,
		StorageDriver:      attachment.StorageDriver,
		StorageKey:         attachment.StorageKey,
		AccessURL:          fmt.Sprintf("%s/api/documents/attachments/%s/download?accessToken=%s&expiresAt=%s", baseURL, attachment.PublicID, token, strconv.FormatInt(expiresAt.Unix(), 10)),
		ExpiresAt:          expiresAt,
		AccessMode:         "read",
	}
	handler.recordDocumentAuditEvent(attachment.TenantSlug, attachment.PublicID, "access_link_created", resolveDocumentsActor(request, payload.RequestedBy), "expires_at="+expiresAt.Format(time.RFC3339), request.Header.Get("X-Correlation-Id"))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler AttachmentHandler) RevokeAccessLink(writer http.ResponseWriter, request *http.Request) {
	var payload dto.AccessLinkRevocationRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	attachmentPublicID, tenantSlug, _, valid := handler.parseSignedAccessToken(payload.AccessToken)
	if !valid || attachmentPublicID != strings.TrimSpace(request.PathValue("publicId")) {
		writeDocumentsError(writer, http.StatusUnauthorized, "invalid_access_token", "Access token is invalid.")
		return
	}
	if strings.TrimSpace(payload.TenantSlug) != "" && !strings.EqualFold(strings.TrimSpace(payload.TenantSlug), tenantSlug) {
		writeDocumentsError(writer, http.StatusForbidden, "tenant_mismatch", "Access token does not belong to the requested tenant.")
		return
	}

	revokedAt := time.Now().UTC()
	reason := strings.TrimSpace(payload.Reason)
	if reason == "" {
		reason = "manual_revocation"
	}
	actor := resolveDocumentsActor(request, payload.RevokedBy)
	revokedAt = handler.attachmentRepository.RevokeAccessToken(tenantSlug, attachmentPublicID, hashAccessToken(payload.AccessToken), reason, actor, request.Header.Get("X-Correlation-Id"))
	handler.recordDocumentAuditEvent(tenantSlug, attachmentPublicID, "access_link_revoked", actor, reason, request.Header.Get("X-Correlation-Id"))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(dto.AccessLinkRevocationResponse{
		AttachmentPublicID: attachmentPublicID,
		TenantSlug:         tenantSlug,
		Revoked:            true,
		Reason:             reason,
		RevokedBy:          actor,
		RevokedAt:          revokedAt,
	})
}

func (handler AttachmentHandler) ListAuditEvents(writer http.ResponseWriter, request *http.Request) {
	tenantSlug := strings.TrimSpace(request.URL.Query().Get("tenantSlug"))
	attachmentPublicID := strings.TrimSpace(request.URL.Query().Get("attachmentPublicId"))

	events := handler.attachmentRepository.ListDocumentAuditEvents(repository.DocumentAuditEventFilters{
		TenantSlug:         tenantSlug,
		AttachmentPublicID: attachmentPublicID,
	})

	response := make([]dto.DocumentAuditEventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, dto.DocumentAuditEventResponse{
			PublicID:           event.PublicID,
			TenantSlug:         event.TenantSlug,
			AttachmentPublicID: event.AttachmentPublicID,
			EventCode:          event.EventCode,
			Actor:              event.Actor,
			Reason:             event.Reason,
			CorrelationID:      event.CorrelationID,
			CreatedAt:          event.CreatedAt,
		})
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler AttachmentHandler) CreateUploadSession(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateUploadSessionRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createUploadSession.Execute(command.CreateUploadSessionInput{
		TenantSlug:       payload.TenantSlug,
		OwnerType:        payload.OwnerType,
		OwnerPublicID:    payload.OwnerPublicID,
		FileName:         payload.FileName,
		ContentType:      payload.ContentType,
		StorageKey:       payload.StorageKey,
		StorageDriver:    payload.StorageDriver,
		Source:           payload.Source,
		RequestedBy:      payload.RequestedBy,
		Visibility:       payload.Visibility,
		RetentionDays:    payload.RetentionDays,
		ExpiresInSeconds: payload.ExpiresInSeconds,
	})
	if result.BadRequest {
		writeDocumentsError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(mapUploadSession(*result.UploadSession, request))
}

func (handler AttachmentHandler) GetUploadSession(writer http.ResponseWriter, request *http.Request) {
	session, ok := handler.getUploadSession.Execute(request.URL.Query().Get("tenantSlug"), request.PathValue("publicId"))
	if !ok {
		writeDocumentsError(writer, http.StatusNotFound, "upload_session_not_found", "Upload session was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(mapUploadSession(*session, request))
}

func (handler AttachmentHandler) CompleteUploadSession(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CompleteUploadSessionRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}
	if suspiciousUpload(payload.ChecksumSHA256) {
		handler.recordDocumentAuditEvent(payload.TenantSlug, request.PathValue("publicId"), "upload_malware_blocked", resolveDocumentsActor(request, payload.UploadedBy), "malware_signature_detected", request.Header.Get("X-Correlation-Id"))
		writeDocumentsError(writer, http.StatusUnprocessableEntity, "malware_detected", "Upload failed malware scan and was blocked.")
		return
	}

	result := handler.completeUploadSession.Execute(command.CompleteUploadSessionInput{
		TenantSlug:     payload.TenantSlug,
		PublicID:       request.PathValue("publicId"),
		UploadedBy:     payload.UploadedBy,
		FileSizeBytes:  payload.FileSizeBytes,
		ChecksumSHA256: payload.ChecksumSHA256,
	})
	if result.BadRequest {
		writeDocumentsError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.Conflict {
		writeDocumentsError(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeDocumentsError(writer, http.StatusNotFound, "upload_session_not_found", "Upload session was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(struct {
		UploadSession dto.UploadSessionResponse `json:"uploadSession"`
		Attachment    dto.AttachmentResponse    `json:"attachment"`
	}{
		UploadSession: mapUploadSession(*result.UploadSession, request),
		Attachment:    mapAttachment(*result.Attachment),
	})
}

func (handler AttachmentHandler) Download(writer http.ResponseWriter, request *http.Request) {
	accessToken := strings.TrimSpace(request.URL.Query().Get("accessToken"))
	attachmentPublicID, tenantSlug, expiresAt, valid := handler.parseSignedAccessToken(accessToken)
	if !valid || attachmentPublicID != strings.TrimSpace(request.PathValue("publicId")) {
		handler.recordDocumentAuditEvent(tenantSlug, request.PathValue("publicId"), "access_link_denied", resolveDocumentsActor(request, ""), "invalid_token", request.Header.Get("X-Correlation-Id"))
		writeDocumentsError(writer, http.StatusUnauthorized, "invalid_access_token", "Access token is invalid.")
		return
	}
	if handler.attachmentRepository.IsAccessTokenRevoked(hashAccessToken(accessToken)) {
		handler.recordDocumentAuditEvent(tenantSlug, attachmentPublicID, "access_link_denied", resolveDocumentsActor(request, ""), "revoked_token", request.Header.Get("X-Correlation-Id"))
		writeDocumentsError(writer, http.StatusUnauthorized, "access_token_revoked", "Access token has been revoked.")
		return
	}
	if time.Now().UTC().After(expiresAt) {
		handler.recordDocumentAuditEvent(tenantSlug, attachmentPublicID, "access_link_denied", resolveDocumentsActor(request, ""), "expired_token", request.Header.Get("X-Correlation-Id"))
		writeDocumentsError(writer, http.StatusUnauthorized, "access_token_expired", "Access token is expired.")
		return
	}

	attachment, ok := handler.getAttachment.Execute(tenantSlug, attachmentPublicID)
	if !ok {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}
	if attachment.ArchivedAt != nil {
		handler.recordDocumentAuditEvent(tenantSlug, attachmentPublicID, "download_denied", resolveDocumentsActor(request, ""), "attachment_archived", request.Header.Get("X-Correlation-Id"))
		writeDocumentsError(writer, http.StatusGone, "attachment_archived", "Attachment is archived.")
		return
	}
	if retentionExpired(*attachment, time.Now().UTC()) {
		handler.recordDocumentAuditEvent(tenantSlug, attachmentPublicID, "download_denied", resolveDocumentsActor(request, ""), "retention_expired", request.Header.Get("X-Correlation-Id"))
		writeDocumentsError(writer, http.StatusGone, "attachment_retention_expired", "Attachment retention window has expired.")
		return
	}

	handler.recordDocumentAuditEvent(tenantSlug, attachmentPublicID, "download_redirected", resolveDocumentsActor(request, ""), "temporary_redirect", request.Header.Get("X-Correlation-Id"))
	writer.Header().Set("Location", buildStorageRedirectURL(request, *attachment, accessToken, expiresAt))
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func mapAttachment(attachment entity.Attachment) dto.AttachmentResponse {
	return dto.AttachmentResponse{
		PublicID:       attachment.PublicID,
		TenantSlug:     attachment.TenantSlug,
		OwnerType:      attachment.OwnerType,
		OwnerPublicID:  attachment.OwnerPublicID,
		FileName:       attachment.FileName,
		ContentType:    attachment.ContentType,
		StorageKey:     attachment.StorageKey,
		StorageDriver:  attachment.StorageDriver,
		Source:         attachment.Source,
		UploadedBy:     attachment.UploadedBy,
		FileSizeBytes:  attachment.FileSizeBytes,
		ChecksumSHA256: attachment.ChecksumSHA256,
		Visibility:     attachment.Visibility,
		RetentionDays:  attachment.RetentionDays,
		CurrentVersion: attachment.CurrentVersion,
		VersionCount:   attachment.VersionCount,
		ArchiveReason:  attachment.ArchiveReason,
		ArchivedAt:     attachment.ArchivedAt,
		CreatedAt:      attachment.CreatedAt,
	}
}

func mapAttachmentVersion(version entity.AttachmentVersion) dto.AttachmentVersionResponse {
	return dto.AttachmentVersionResponse{
		PublicID:           version.PublicID,
		TenantSlug:         version.TenantSlug,
		AttachmentPublicID: version.AttachmentPublicID,
		VersionNumber:      version.VersionNumber,
		FileName:           version.FileName,
		ContentType:        version.ContentType,
		StorageKey:         version.StorageKey,
		StorageDriver:      version.StorageDriver,
		Source:             version.Source,
		UploadedBy:         version.UploadedBy,
		FileSizeBytes:      version.FileSizeBytes,
		ChecksumSHA256:     version.ChecksumSHA256,
		CreatedAt:          version.CreatedAt,
	}
}

func mapUploadSession(session entity.UploadSession, request *http.Request) dto.UploadSessionResponse {
	return dto.UploadSessionResponse{
		PublicID:           session.PublicID,
		TenantSlug:         session.TenantSlug,
		OwnerType:          session.OwnerType,
		OwnerPublicID:      session.OwnerPublicID,
		FileName:           session.FileName,
		ContentType:        session.ContentType,
		StorageKey:         session.StorageKey,
		StorageDriver:      session.StorageDriver,
		Source:             session.Source,
		RequestedBy:        session.RequestedBy,
		Visibility:         session.Visibility,
		RetentionDays:      session.RetentionDays,
		Status:             session.Status,
		AttachmentPublicID: session.AttachmentPublicID,
		UploadURL:          buildUploadURL(request, session.PublicID),
		ExpiresAt:          session.ExpiresAt,
		CompletedAt:        session.CompletedAt,
		CreatedAt:          session.CreatedAt,
	}
}

func buildUploadURL(request *http.Request, publicID string) string {
	baseURL := inferBaseURL(request)
	return fmt.Sprintf("%s/api/documents/upload-sessions/%s/complete", baseURL, publicID)
}

func inferBaseURL(request *http.Request) string {
	scheme := strings.TrimSpace(request.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = "http"
		if request.TLS != nil {
			scheme = "https"
		}
	}
	host := strings.TrimSpace(request.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = request.Host
	}
	if host == "" {
		host = "documents.local"
	}

	return scheme + "://" + strings.TrimRight(host, "/")
}
func (handler AttachmentHandler) newSignedAccessToken(attachmentPublicID string, tenantSlug string, expiresAt time.Time) string {
	raw := fmt.Sprintf("%s|%s|%d", strings.TrimSpace(attachmentPublicID), strings.ToLower(strings.TrimSpace(tenantSlug)), expiresAt.Unix())
	mac := hmac.New(sha256.New, []byte(handler.tokenSecret()))
	_, _ = mac.Write([]byte(raw))
	signature := fmt.Sprintf("%x", mac.Sum(nil))
	payload := raw + "|" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(payload))
}

func (handler AttachmentHandler) parseSignedAccessToken(token string) (string, string, time.Time, bool) {
	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(token))
	if err != nil {
		return "", "", time.Time{}, false
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 4 {
		return "", "", time.Time{}, false
	}

	publicID := strings.TrimSpace(parts[0])
	tenantSlug := strings.ToLower(strings.TrimSpace(parts[1]))
	expiresUnix, err := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
	if err != nil {
		return "", "", time.Time{}, false
	}

	raw := fmt.Sprintf("%s|%s|%d", publicID, tenantSlug, expiresUnix)
	mac := hmac.New(sha256.New, []byte(handler.tokenSecret()))
	_, _ = mac.Write([]byte(raw))
	expectedSignature := fmt.Sprintf("%x", mac.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(strings.TrimSpace(parts[3]))) {
		return "", "", time.Time{}, false
	}

	return publicID, tenantSlug, time.Unix(expiresUnix, 0).UTC(), true
}

func hashAccessToken(token string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return fmt.Sprintf("%x", digest[:])
}

func (handler AttachmentHandler) recordDocumentAuditEvent(tenantSlug string, attachmentPublicID string, eventCode string, actor string, reason string, correlationID string) {
	if actor == "" {
		actor = "system"
	}
	_ = handler.attachmentRepository.RecordDocumentAuditEvent(repository.DocumentAuditEvent{
		PublicID:           uuid.NewString(),
		TenantSlug:         strings.ToLower(strings.TrimSpace(tenantSlug)),
		AttachmentPublicID: strings.TrimSpace(attachmentPublicID),
		EventCode:          strings.TrimSpace(eventCode),
		Actor:              strings.TrimSpace(actor),
		Reason:             strings.TrimSpace(reason),
		CorrelationID:      strings.TrimSpace(correlationID),
		CreatedAt:          time.Now().UTC(),
	})
}

func resolveDocumentsActor(request *http.Request, explicitActor string) string {
	if strings.TrimSpace(explicitActor) != "" {
		return strings.TrimSpace(explicitActor)
	}
	for _, header := range []string{"X-ERP-Auth-Subject", "X-Actor", "X-User-Public-Id"} {
		if value := strings.TrimSpace(request.Header.Get(header)); value != "" {
			return value
		}
	}
	return "system"
}

func suspiciousUpload(checksum string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(checksum))
	return strings.Contains(normalized, "EICAR") || strings.Contains(normalized, "MALWARE")
}

func retentionExpired(attachment entity.Attachment, now time.Time) bool {
	if attachment.RetentionDays <= 0 {
		return false
	}
	return attachment.CreatedAt.AddDate(0, 0, attachment.RetentionDays).Before(now)
}

func (handler AttachmentHandler) tokenSecret() string {
	if strings.TrimSpace(handler.accessTokenSecret) == "" && localRuntimeMode() {
		return "documents-local-secret"
	}
	if strings.TrimSpace(handler.accessTokenSecret) == "" {
		panic("DOCUMENTS_ACCESS_TOKEN_SECRET is required outside local/test environments")
	}

	return handler.accessTokenSecret
}

func localRuntimeMode() bool {
	environment := strings.ToLower(strings.TrimSpace(os.Getenv("ERP_ENV")))
	return environment == "" || environment == "local" || environment == "dev" || environment == "development" || environment == "test"
}

func buildStorageRedirectURL(request *http.Request, attachment entity.Attachment, accessToken string, expiresAt time.Time) string {
	baseURL := inferBaseURL(request)
	storagePath := fmt.Sprintf(
		"%s/files/%s/%s?tenantSlug=%s&attachmentPublicId=%s&accessToken=%s&expiresAt=%d",
		baseURL,
		url.PathEscape(strings.ToLower(strings.TrimSpace(attachment.StorageDriver))),
		url.PathEscape(attachment.StorageKey),
		url.QueryEscape(attachment.TenantSlug),
		url.QueryEscape(attachment.PublicID),
		url.QueryEscape(accessToken),
		expiresAt.Unix(),
	)

	return storagePath
}

func writeDocumentsError(writer http.ResponseWriter, status int, code string, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
