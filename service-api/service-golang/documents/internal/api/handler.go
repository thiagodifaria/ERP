package api

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type AttachmentHandler struct {
	listAttachments  query.ListAttachments
	createAttachment command.CreateAttachment
}

func NewAttachmentHandler(repository repository.AttachmentRepository) AttachmentHandler {
	return AttachmentHandler{
		listAttachments:  query.NewListAttachments(repository),
		createAttachment: command.NewCreateAttachment(repository),
	}
}

func (handler AttachmentHandler) List(writer http.ResponseWriter, request *http.Request) {
	attachments := handler.listAttachments.Execute(repository.AttachmentFilters{
		TenantSlug:    request.URL.Query().Get("tenantSlug"),
		OwnerType:     request.URL.Query().Get("ownerType"),
		OwnerPublicID: request.URL.Query().Get("ownerPublicId"),
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
		TenantSlug:    payload.TenantSlug,
		OwnerType:     payload.OwnerType,
		OwnerPublicID: payload.OwnerPublicID,
		FileName:      payload.FileName,
		ContentType:   payload.ContentType,
		StorageKey:    payload.StorageKey,
		StorageDriver: payload.StorageDriver,
		Source:        payload.Source,
		UploadedBy:    payload.UploadedBy,
	})

	if result.BadRequest {
		writeDocumentsError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(mapAttachment(*result.Attachment))
}

func mapAttachment(attachment entity.Attachment) dto.AttachmentResponse {
	return dto.AttachmentResponse{
		PublicID:      attachment.PublicID,
		TenantSlug:    attachment.TenantSlug,
		OwnerType:     attachment.OwnerType,
		OwnerPublicID: attachment.OwnerPublicID,
		FileName:      attachment.FileName,
		ContentType:   attachment.ContentType,
		StorageKey:    attachment.StorageKey,
		StorageDriver: attachment.StorageDriver,
		Source:        attachment.Source,
		UploadedBy:    attachment.UploadedBy,
		CreatedAt:     attachment.CreatedAt,
	}
}

func writeDocumentsError(writer http.ResponseWriter, status int, code string, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}
