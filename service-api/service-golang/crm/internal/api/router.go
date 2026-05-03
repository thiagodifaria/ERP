// Router define as rotas publicas do servico crm.
// Composicao de regras de negocio nao deve acontecer aqui.
package api

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/handler"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/middleware"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/telemetry"
)

func NewRouter(
	logger *telemetry.Logger,
	repositories repository.TenantRepositoryFactory,
	attachmentGateway repository.AttachmentGateway,
) http.Handler {
	mux := http.NewServeMux()
	leadHandler := handler.NewLeadHandler(repositories)
	customerHandler := handler.NewCustomerHandler(repositories)
	leadNoteHandler := handler.NewLeadNoteHandler(repositories)
	activityHandler := handler.NewActivityHandler(repositories)
	attachmentHandler := handler.NewAttachmentHandler(repositories, attachmentGateway)
	intelligenceHandler := handler.NewIntelligenceHandler(repositories)

	mux.HandleFunc("/health/live", handler.Live)
	mux.HandleFunc("/health/ready", handler.Ready)
	mux.HandleFunc("/health/details", handler.DetailsForRuntime("memory"))
	mux.HandleFunc("GET /api/crm/leads/summary", leadHandler.Summary)
	mux.HandleFunc("GET /api/crm/leads", leadHandler.List)
	mux.HandleFunc("POST /api/crm/leads", leadHandler.Create)
	mux.HandleFunc("GET /api/crm/leads/{publicId}", leadHandler.GetByPublicID)
	mux.HandleFunc("POST /api/crm/leads/{publicId}/convert", leadHandler.Convert)
	mux.HandleFunc("GET /api/crm/leads/{publicId}/history", activityHandler.ListLeadHistory)
	mux.HandleFunc("GET /api/crm/leads/{publicId}/notes", leadNoteHandler.List)
	mux.HandleFunc("POST /api/crm/leads/{publicId}/notes", leadNoteHandler.Create)
	mux.HandleFunc("GET /api/crm/leads/{publicId}/attachments", attachmentHandler.ListLeadAttachments)
	mux.HandleFunc("POST /api/crm/leads/{publicId}/attachments", attachmentHandler.CreateLeadAttachment)
	mux.HandleFunc("GET /api/crm/customers", customerHandler.List)
	mux.HandleFunc("GET /api/crm/customers/{publicId}", customerHandler.GetByPublicID)
	mux.HandleFunc("GET /api/crm/customers/{publicId}/history", activityHandler.ListCustomerHistory)
	mux.HandleFunc("GET /api/crm/customers/{publicId}/attachments", attachmentHandler.ListCustomerAttachments)
	mux.HandleFunc("POST /api/crm/customers/{publicId}/attachments", attachmentHandler.CreateCustomerAttachment)
	mux.HandleFunc("GET /api/crm/outbox/pending", activityHandler.ListPendingOutbox)
	mux.HandleFunc("GET /api/crm/pipeline/config", intelligenceHandler.GetPipelineConfig)
	mux.HandleFunc("PUT /api/crm/pipeline/config", intelligenceHandler.UpsertPipelineConfig)
	mux.HandleFunc("GET /api/crm/leads/intelligence/summary", intelligenceHandler.Summary)
	mux.HandleFunc("PATCH /api/crm/leads/{publicId}", leadHandler.UpdateProfile)
	mux.HandleFunc("PATCH /api/crm/leads/{publicId}/owner", leadHandler.UpdateOwner)
	mux.HandleFunc("PATCH /api/crm/leads/{publicId}/status", leadHandler.UpdateStatus)

	return middleware.WithCorrelation(logger, mux)
}

func NewRouterWithRuntime(
	logger *telemetry.Logger,
	repositories repository.TenantRepositoryFactory,
	attachmentGateway repository.AttachmentGateway,
	repositoryDriver string,
) http.Handler {
	mux := http.NewServeMux()
	leadHandler := handler.NewLeadHandler(repositories)
	customerHandler := handler.NewCustomerHandler(repositories)
	leadNoteHandler := handler.NewLeadNoteHandler(repositories)
	activityHandler := handler.NewActivityHandler(repositories)
	attachmentHandler := handler.NewAttachmentHandler(repositories, attachmentGateway)
	intelligenceHandler := handler.NewIntelligenceHandler(repositories)

	mux.HandleFunc("/health/live", handler.Live)
	mux.HandleFunc("/health/ready", handler.Ready)
	mux.HandleFunc("/health/details", handler.DetailsForRuntime(repositoryDriver))
	mux.HandleFunc("GET /api/crm/leads/summary", leadHandler.Summary)
	mux.HandleFunc("GET /api/crm/leads", leadHandler.List)
	mux.HandleFunc("POST /api/crm/leads", leadHandler.Create)
	mux.HandleFunc("GET /api/crm/leads/{publicId}", leadHandler.GetByPublicID)
	mux.HandleFunc("POST /api/crm/leads/{publicId}/convert", leadHandler.Convert)
	mux.HandleFunc("GET /api/crm/leads/{publicId}/history", activityHandler.ListLeadHistory)
	mux.HandleFunc("GET /api/crm/leads/{publicId}/notes", leadNoteHandler.List)
	mux.HandleFunc("POST /api/crm/leads/{publicId}/notes", leadNoteHandler.Create)
	mux.HandleFunc("GET /api/crm/leads/{publicId}/attachments", attachmentHandler.ListLeadAttachments)
	mux.HandleFunc("POST /api/crm/leads/{publicId}/attachments", attachmentHandler.CreateLeadAttachment)
	mux.HandleFunc("GET /api/crm/customers", customerHandler.List)
	mux.HandleFunc("GET /api/crm/customers/{publicId}", customerHandler.GetByPublicID)
	mux.HandleFunc("GET /api/crm/customers/{publicId}/history", activityHandler.ListCustomerHistory)
	mux.HandleFunc("GET /api/crm/customers/{publicId}/attachments", attachmentHandler.ListCustomerAttachments)
	mux.HandleFunc("POST /api/crm/customers/{publicId}/attachments", attachmentHandler.CreateCustomerAttachment)
	mux.HandleFunc("GET /api/crm/outbox/pending", activityHandler.ListPendingOutbox)
	mux.HandleFunc("GET /api/crm/pipeline/config", intelligenceHandler.GetPipelineConfig)
	mux.HandleFunc("PUT /api/crm/pipeline/config", intelligenceHandler.UpsertPipelineConfig)
	mux.HandleFunc("GET /api/crm/leads/intelligence/summary", intelligenceHandler.Summary)
	mux.HandleFunc("PATCH /api/crm/leads/{publicId}", leadHandler.UpdateProfile)
	mux.HandleFunc("PATCH /api/crm/leads/{publicId}/owner", leadHandler.UpdateOwner)
	mux.HandleFunc("PATCH /api/crm/leads/{publicId}/status", leadHandler.UpdateStatus)

	return middleware.WithCorrelation(logger, mux)
}
