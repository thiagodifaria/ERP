// InvoiceHandler expoe a trilha publica de faturamento do Sales MVP.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
)

type InvoiceHandler struct {
	listInvoices         query.ListInvoices
	getInvoiceSummary    query.GetInvoiceSummary
	getInvoiceByPublicID query.GetInvoiceByPublicID
	createInvoice        command.CreateInvoice
	updateInvoiceStatus  command.UpdateInvoiceStatus
}

func NewInvoiceHandler(
	listInvoices query.ListInvoices,
	getInvoiceSummary query.GetInvoiceSummary,
	getInvoiceByPublicID query.GetInvoiceByPublicID,
	createInvoice command.CreateInvoice,
	updateInvoiceStatus command.UpdateInvoiceStatus,
) InvoiceHandler {
	return InvoiceHandler{
		listInvoices:         listInvoices,
		getInvoiceSummary:    getInvoiceSummary,
		getInvoiceByPublicID: getInvoiceByPublicID,
		createInvoice:        createInvoice,
		updateInvoiceStatus:  updateInvoiceStatus,
	}
}

func (handler InvoiceHandler) List(writer http.ResponseWriter, request *http.Request) {
	invoices := handler.listInvoices.Execute(query.InvoiceFilters{
		Status: request.URL.Query().Get("status"),
	})
	response := make([]dto.InvoiceResponse, 0, len(invoices))

	for _, invoice := range invoices {
		response = append(response, mapInvoice(invoice))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler InvoiceHandler) Summary(writer http.ResponseWriter, request *http.Request) {
	summary := handler.getInvoiceSummary.Execute(query.InvoiceFilters{
		Status: request.URL.Query().Get("status"),
	})

	writeJSON(writer, http.StatusOK, dto.InvoiceSummaryResponse{
		Total:              summary.Total,
		OpenAmountCents:    summary.OpenAmountCents,
		PaidAmountCents:    summary.PaidAmountCents,
		OverdueAmountCents: summary.OverdueAmountCents,
		OverdueCount:       summary.OverdueCount,
		ByStatus:           summary.ByStatus,
	})
}

func (handler InvoiceHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	invoice := handler.getInvoiceByPublicID.Execute(request.PathValue("publicId"))
	if invoice == nil {
		writeError(writer, http.StatusNotFound, "invoice_not_found", "Invoice was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapInvoice(*invoice))
}

func (handler InvoiceHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateInvoiceRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createInvoice.Execute(command.CreateInvoiceInput{
		SalePublicID: request.PathValue("publicId"),
		Number:       payload.Number,
		DueDate:      payload.DueDate,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	case result.Conflict:
		writeError(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusCreated, mapInvoice(*result.Invoice))
	}
}

func (handler InvoiceHandler) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateInvoiceStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.updateInvoiceStatus.Execute(command.UpdateInvoiceStatusInput{
		PublicID: request.PathValue("publicId"),
		Status:   payload.Status,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusOK, mapInvoice(*result.Invoice))
	}
}
