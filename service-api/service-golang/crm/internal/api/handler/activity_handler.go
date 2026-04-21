package handler

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type ActivityHandler struct {
	repositories repository.TenantRepositoryFactory
}

func NewActivityHandler(repositories repository.TenantRepositoryFactory) ActivityHandler {
	return ActivityHandler{repositories: repositories}
}

func (handler ActivityHandler) ListLeadHistory(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := resolveTenantRepositories(writer, request, "", handler.repositories)
	if bundle.LeadRepository == nil {
		return
	}

	lead := query.NewGetLeadByPublicID(bundle.LeadRepository).Execute(request.PathValue("publicId"))
	if lead == nil {
		writeNotFound(writer, "lead_not_found", "Lead was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapRelationshipEvents(query.NewListRelationshipEventsByAggregate(bundle.RelationshipEventRepository).Execute("lead", lead.PublicID)))
}

func (handler ActivityHandler) ListCustomerHistory(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := resolveTenantRepositories(writer, request, "", handler.repositories)
	if bundle.CustomerRepository == nil {
		return
	}

	customer := query.NewGetCustomerByPublicID(bundle.CustomerRepository).Execute(request.PathValue("publicId"))
	if customer == nil {
		writeNotFound(writer, "customer_not_found", "Customer was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapRelationshipEvents(query.NewListRelationshipEventsByAggregate(bundle.RelationshipEventRepository).Execute("customer", customer.PublicID)))
}

func (handler ActivityHandler) ListPendingOutbox(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := resolveTenantRepositories(writer, request, "", handler.repositories)
	if bundle.OutboxEventRepository == nil {
		return
	}

	response := make([]dto.OutboxEventResponse, 0)
	for _, event := range query.NewListPendingOutboxEvents(bundle.OutboxEventRepository).Execute(100) {
		response = append(response, dto.OutboxEventResponse{
			PublicID:          event.PublicID,
			AggregateType:     event.AggregateType,
			AggregatePublicID: event.AggregatePublicID,
			EventType:         event.EventType,
			Payload:           event.Payload,
			Status:            event.Status,
			CreatedAt:         event.CreatedAt,
			ProcessedAt:       event.ProcessedAt,
		})
	}

	writeJSON(writer, http.StatusOK, response)
}

func mapRelationshipEvents(events []entity.RelationshipEvent) []dto.RelationshipEventResponse {
	response := make([]dto.RelationshipEventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, dto.RelationshipEventResponse{
			PublicID:          event.PublicID,
			AggregateType:     event.AggregateType,
			AggregatePublicID: event.AggregatePublicID,
			EventCode:         event.EventCode,
			Actor:             event.Actor,
			Summary:           event.Summary,
			CreatedAt:         event.CreatedAt,
		})
	}

	return response
}
