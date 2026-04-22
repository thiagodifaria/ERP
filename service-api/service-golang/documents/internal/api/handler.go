package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
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
}

func NewAttachmentHandler(repository repository.AttachmentRepository) AttachmentHandler {
	return AttachmentHandler{
		listAttachments:   query.NewListAttachments(repository),
		createAttachment:  command.NewCreateAttachment(repository),
		getAttachment:     query.NewGetAttachment(repository),
		archiveAttachment: command.NewArchiveAttachment(repository),
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
	token := newAccessToken()
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

func newAccessToken() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "documents-access-token"
	}

	return fmt.Sprintf("%x", raw)
}

func writeDocumentsError(writer http.ResponseWriter, status int, code string, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
