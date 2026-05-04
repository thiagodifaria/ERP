package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/config"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type SigningHandler struct {
	attachments repository.AttachmentRepository
	cfg         config.Config
}

func NewSigningHandler(attachmentRepository repository.AttachmentRepository, cfg config.Config) SigningHandler {
	return SigningHandler{
		attachments: attachmentRepository,
		cfg:         cfg,
	}
}

func buildSigningCapabilities(cfg config.Config) []dto.SigningCapabilityResponse {
	clicksignConfigured := strings.TrimSpace(cfg.ClicksignAPIKey) != ""
	docusignConfigured := strings.TrimSpace(cfg.DocusignAccessToken) != ""

	return []dto.SigningCapabilityResponse{
		{
			Provider:          "local",
			Scope:             "digital_signature",
			Configured:        true,
			Mode:              "fallback",
			Status:            "fallback",
			FallbackViable:    true,
			SupportsContracts: true,
			SupportsProposals: true,
			SupportsRentals:   true,
			Notes:             []string{"Modo local simula assinatura e mantem o fluxo operacional para desenvolvimento e smoke."},
		},
		{
			Provider:          "clicksign",
			Scope:             "digital_signature",
			Configured:        clicksignConfigured,
			CredentialKey:     "DOCUMENTS_CLICKSIGN_API_KEY",
			Mode:              map[bool]string{true: "configured", false: "fallback"}[clicksignConfigured],
			Status:            map[bool]string{true: "ready", false: "fallback"}[clicksignConfigured],
			FallbackViable:    true,
			SupportsContracts: true,
			SupportsProposals: true,
			SupportsRentals:   true,
			Notes:             []string{map[bool]string{true: "Clicksign pronto para contratos, propostas e anexos juridicos.", false: "Sem chave Clicksign, o runtime continua com simulacao local sem quebrar o fluxo."}[clicksignConfigured]},
		},
		{
			Provider:          "docusign",
			Scope:             "digital_signature",
			Configured:        docusignConfigured,
			CredentialKey:     "DOCUMENTS_DOCUSIGN_ACCESS_TOKEN",
			Mode:              map[bool]string{true: "configured", false: "fallback"}[docusignConfigured],
			Status:            map[bool]string{true: "ready", false: "fallback"}[docusignConfigured],
			FallbackViable:    true,
			SupportsContracts: true,
			SupportsProposals: true,
			SupportsRentals:   true,
			Notes:             []string{map[bool]string{true: "DocuSign pronto como provider alternativo de assinatura.", false: "Quando DocuSign nao estiver configurado, o servico opera com fallback local."}[docusignConfigured]},
		},
	}
}

func SigningCapabilitiesForRuntime(cfg config.Config) http.HandlerFunc {
	capabilities := buildSigningCapabilities(cfg)

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		provider := strings.TrimSpace(request.PathValue("provider"))
		if provider == "" {
			writer.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(writer).Encode(capabilities)
			return
		}

		for _, capability := range capabilities {
			if capability.Provider == provider {
				writer.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(writer).Encode(capability)
				return
			}
		}

		writer.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(writer).Encode(map[string]string{
			"code":    "signing_capability_not_found",
			"message": "Signing capability was not found.",
		})
	}
}

func (handler SigningHandler) CreateRequest(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateSigningRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeDocumentsError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	provider := strings.TrimSpace(payload.Provider)
	if provider == "" {
		provider = strings.TrimSpace(handler.cfg.SigningProvider)
	}
	if provider == "" {
		provider = "local"
	}
	if strings.TrimSpace(payload.TenantSlug) == "" {
		writeDocumentsError(writer, http.StatusBadRequest, "tenant_slug_required", "Tenant slug is required.")
		return
	}
	if strings.TrimSpace(payload.AttachmentPublicID) == "" {
		writeDocumentsError(writer, http.StatusBadRequest, "attachment_public_id_required", "Attachment public id is required.")
		return
	}
	if strings.TrimSpace(payload.RequestedBy) == "" {
		writeDocumentsError(writer, http.StatusBadRequest, "requested_by_required", "Requested by is required.")
		return
	}
	if len(payload.Signers) == 0 {
		writeDocumentsError(writer, http.StatusBadRequest, "signers_required", "At least one signer is required.")
		return
	}

	attachment, ok := handler.attachments.FindByPublicID(payload.TenantSlug, payload.AttachmentPublicID)
	if !ok {
		writeDocumentsError(writer, http.StatusNotFound, "attachment_not_found", "Attachment was not found.")
		return
	}

	response := dto.SigningRequestResponse{
		PublicID:                 uuid.NewString(),
		TenantSlug:               payload.TenantSlug,
		AttachmentPublicID:       attachment.PublicID,
		Provider:                 provider,
		DocumentKind:             strings.TrimSpace(payload.DocumentKind),
		Status:                   "queued",
		RequestedBy:              payload.RequestedBy,
		Signers:                  payload.Signers,
		RelatedAggregate:         strings.TrimSpace(payload.RelatedAggregate),
		RelatedAggregatePublicID: strings.TrimSpace(payload.RelatedAggregateID),
		SigningURL:               fmt.Sprintf("https://signing.%s.local/requests/%s", provider, attachment.PublicID),
		CreatedAt:                time.Now().UTC(),
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(response)
}
