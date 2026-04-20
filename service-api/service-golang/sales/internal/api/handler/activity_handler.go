package handler

import (
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

type ActivityHandler struct {
	listEvents query.ListCommercialEventsByAggregate
	listOutbox query.ListPendingOutboxEvents
}

func NewActivityHandler(
	listEvents query.ListCommercialEventsByAggregate,
	listOutbox query.ListPendingOutboxEvents,
) ActivityHandler {
	return ActivityHandler{
		listEvents: listEvents,
		listOutbox: listOutbox,
	}
}

func (handler ActivityHandler) ListOpportunityHistory(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, mapCommercialEvents(handler.listEvents.Execute("opportunity", request.PathValue("publicId"))))
}

func (handler ActivityHandler) ListProposalHistory(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, mapCommercialEvents(handler.listEvents.Execute("proposal", request.PathValue("publicId"))))
}

func (handler ActivityHandler) ListSaleHistory(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, mapCommercialEvents(handler.listEvents.Execute("sale", request.PathValue("publicId"))))
}

func (handler ActivityHandler) ListInvoiceHistory(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, mapCommercialEvents(handler.listEvents.Execute("invoice", request.PathValue("publicId"))))
}

func (handler ActivityHandler) ListPendingOutbox(writer http.ResponseWriter, request *http.Request) {
	response := make([]dto.OutboxEventResponse, 0)
	for _, event := range handler.listOutbox.Execute(100) {
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

func mapCommercialEvents(events []entity.CommercialEvent) []dto.CommercialEventResponse {
	response := make([]dto.CommercialEventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, dto.CommercialEventResponse{
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
