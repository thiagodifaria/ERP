package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type AttachmentHandler struct {
	listAttachments   query.ListAttachments
	createAttachment  command.CreateAttachment
	getAttachment     query.GetAttachment
	archiveAttachment command.ArchiveAttachment
	createUploadSession   command.CreateUploadSession
	getUploadSession      query.GetUploadSession
	completeUploadSession command.CompleteUploadSession
	accessTokenSecret     string
}

func NewAttachmentHandler(attachmentRepository repository.AttachmentRepository, uploadSessionRepository repository.UploadSessionRepository, accessTokenSecret string) AttachmentHandler {
	return AttachmentHandler{
		listAttachments:       query.NewListAttachments(attachmentRepository),
		createAttachment:      command.NewCreateAttachment(attachmentRepository),
		getAttachment:         query.NewGetAttachment(attachmentRepository),
		archiveAttachment:     command.NewArchiveAttachment(attachmentRepository),
		createUploadSession:   command.NewCreateUploadSession(uploadSessionRepository),
		getUploadSession:      query.NewGetUploadSession(uploadSessionRepository),
		completeUploadSession: command.NewCompleteUploadSession(uploadSessionRepository, attachmentRepository),
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

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
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
		writeDocumentsError(writer, http.StatusUnauthorized, "invalid_access_token", "Access token is invalid.")
		return
	}
	if time.Now().UTC().After(expiresAt) {
		writeDocumentsError(writer, http.StatusUnauthorized, "access_token_expired", "Access token is expired.")
		return
	}

	attachment, ok := handler.getAttachment.Execute(tenantSlug, attachmentPublicID)
	if !ok {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}
	if attachment.ArchivedAt != nil {
		writeDocumentsError(writer, http.StatusGone, "attachment_archived", "Attachment is archived.")
		return
	}

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
		ArchiveReason:  attachment.ArchiveReason,
		ArchivedAt:     attachment.ArchivedAt,
		CreatedAt:      attachment.CreatedAt,
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

func (handler AttachmentHandler) tokenSecret() string {
	if handler.accessTokenSecret == "" {
		return "documents-local-secret"
	}

	return handler.accessTokenSecret
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
