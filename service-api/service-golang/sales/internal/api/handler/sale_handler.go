// SaleHandler expoe a trilha publica de fechamento do Sales MVP.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
)

type SaleHandler struct {
	listSales         query.ListSales
	getSaleSummary    query.GetSaleSummary
	getSaleByPublicID query.GetSaleByPublicID
	updateSaleStatus  command.UpdateSaleStatus
}

func NewSaleHandler(
	listSales query.ListSales,
	getSaleSummary query.GetSaleSummary,
	getSaleByPublicID query.GetSaleByPublicID,
	updateSaleStatus command.UpdateSaleStatus,
) SaleHandler {
	return SaleHandler{
		listSales:         listSales,
		getSaleSummary:    getSaleSummary,
		getSaleByPublicID: getSaleByPublicID,
		updateSaleStatus:  updateSaleStatus,
	}
}

func (handler SaleHandler) List(writer http.ResponseWriter, request *http.Request) {
	sales := handler.listSales.Execute(query.SaleFilters{
		Status: request.URL.Query().Get("status"),
	})
	response := make([]dto.SaleResponse, 0, len(sales))

	for _, sale := range sales {
		response = append(response, mapSale(sale))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler SaleHandler) Summary(writer http.ResponseWriter, request *http.Request) {
	summary := handler.getSaleSummary.Execute(query.SaleFilters{
		Status: request.URL.Query().Get("status"),
	})

	writeJSON(writer, http.StatusOK, dto.SaleSummaryResponse{
		Total:              summary.Total,
		BookedRevenueCents: summary.BookedRevenueCents,
		ByStatus:           summary.ByStatus,
	})
}

func (handler SaleHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	sale := handler.getSaleByPublicID.Execute(request.PathValue("publicId"))
	if sale == nil {
		writeError(writer, http.StatusNotFound, "sale_not_found", "Sale was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapSale(*sale))
}

func (handler SaleHandler) UpdateStatus(writer http.ResponseWriter, request *http.Request) {
	var payload dto.UpdateSaleStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.updateSaleStatus.Execute(command.UpdateSaleStatusInput{
		PublicID: request.PathValue("publicId"),
		Status:   payload.Status,
	})

	switch {
	case result.BadRequest:
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
	case result.NotFound:
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
	default:
		writeJSON(writer, http.StatusOK, mapSale(*result.Sale))
	}
}
