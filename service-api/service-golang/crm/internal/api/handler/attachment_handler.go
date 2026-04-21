package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type AttachmentHandler struct {
	repositories repository.TenantRepositoryFactory
	gateway      repository.AttachmentGateway
}

func NewAttachmentHandler(repositories repository.TenantRepositoryFactory, gateway repository.AttachmentGateway) AttachmentHandler {
	return AttachmentHandler{
		repositories: repositories,
		gateway:      gateway,
	}
}

func (handler AttachmentHandler) ListLeadAttachments(writer http.ResponseWriter, request *http.Request) {
	handler.listAttachmentsForOwner(writer, request, "lead", "crm.lead")
}

func (handler AttachmentHandler) CreateLeadAttachment(writer http.ResponseWriter, request *http.Request) {
	handler.createAttachmentForOwner(writer, request, "lead", "crm.lead")
}

func (handler AttachmentHandler) ListCustomerAttachments(writer http.ResponseWriter, request *http.Request) {
	handler.listAttachmentsForOwner(writer, request, "customer", "crm.customer")
}

func (handler AttachmentHandler) CreateCustomerAttachment(writer http.ResponseWriter, request *http.Request) {
	handler.createAttachmentForOwner(writer, request, "customer", "crm.customer")
}

func (handler AttachmentHandler) listAttachmentsForOwner(writer http.ResponseWriter, request *http.Request, aggregateType string, ownerType string) {
	bundle, tenantSlug := resolveTenantRepositories(writer, request, "", handler.repositories)
	if tenantSlug == "" {
		return
	}

	ownerPublicID, ok := handler.validateAggregate(writer, request, bundle, aggregateType)
	if !ok {
		return
	}

	attachments, err := handler.gateway.List(tenantSlug, ownerType, ownerPublicID)
	if err != nil {
		writeErrorResponse(writer, http.StatusBadGateway, "documents_unavailable", "Documents service is unavailable.")
		return
	}

	response := make([]dto.AttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		response = append(response, dto.AttachmentResponse{
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
		})
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler AttachmentHandler) createAttachmentForOwner(writer http.ResponseWriter, request *http.Request, aggregateType string, ownerType string) {
	var payload dto.CreateAttachmentRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeBadRequest(writer, "invalid_json", "Request body is invalid.")
		return
	}

	bundle, tenantSlug := resolveTenantRepositories(writer, request, "", handler.repositories)
	if tenantSlug == "" {
		return
	}

	ownerPublicID, ok := handler.validateAggregate(writer, request, bundle, aggregateType)
	if !ok {
		return
	}

	attachment, err := handler.gateway.Create(repository.CreateAttachmentInput{
		TenantSlug:    tenantSlug,
		OwnerType:     ownerType,
		OwnerPublicID: ownerPublicID,
		FileName:      payload.FileName,
		ContentType:   payload.ContentType,
		StorageKey:    payload.StorageKey,
		StorageDriver: payload.StorageDriver,
		Source:        payload.Source,
		UploadedBy:    payload.UploadedBy,
	})
	if err != nil {
		writeErrorResponse(writer, http.StatusBadGateway, "documents_unavailable", "Documents service is unavailable.")
		return
	}

	writeJSON(writer, http.StatusCreated, dto.AttachmentResponse{
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
	})
}

func (handler AttachmentHandler) validateAggregate(writer http.ResponseWriter, request *http.Request, bundle repository.TenantRepositorySet, aggregateType string) (string, bool) {
	publicID := request.PathValue("publicId")
	switch aggregateType {
	case "lead":
		lead := query.NewGetLeadByPublicID(bundle.LeadRepository).Execute(publicID)
		if lead == nil {
			writeNotFound(writer, "lead_not_found", "Lead was not found.")
			return "", false
		}
		return lead.PublicID, true
	case "customer":
		customer := query.NewGetCustomerByPublicID(bundle.CustomerRepository).Execute(publicID)
		if customer == nil {
			writeNotFound(writer, "customer_not_found", "Customer was not found.")
			return "", false
		}
		return customer.PublicID, true
	default:
		writeBadRequest(writer, "invalid_attachment_owner", "Attachment owner is invalid.")
		return "", false
	}
}
