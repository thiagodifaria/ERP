package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/command"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/application/query"
)

type OperationsHandler struct {
	listInstallments    query.ListInstallmentsBySale
	createInstallments  command.CreateInstallmentSchedule
	listCommissions     query.ListCommissionsBySale
	createCommission    command.CreateCommission
	updateCommission    command.UpdateCommissionStatus
	listPendingItems    query.ListPendingItemsBySale
	createPendingItem   command.CreatePendingItem
	resolvePendingItem  command.ResolvePendingItem
	listRenegotiations  query.ListRenegotiationsBySale
	applyRenegotiation  command.ApplyRenegotiation
	cancelSale          command.CancelSale
}

func NewOperationsHandler(
	listInstallments query.ListInstallmentsBySale,
	createInstallments command.CreateInstallmentSchedule,
	listCommissions query.ListCommissionsBySale,
	createCommission command.CreateCommission,
	updateCommission command.UpdateCommissionStatus,
	listPendingItems query.ListPendingItemsBySale,
	createPendingItem command.CreatePendingItem,
	resolvePendingItem command.ResolvePendingItem,
	listRenegotiations query.ListRenegotiationsBySale,
	applyRenegotiation command.ApplyRenegotiation,
	cancelSale command.CancelSale,
) OperationsHandler {
	return OperationsHandler{
		listInstallments:   listInstallments,
		createInstallments: createInstallments,
		listCommissions:    listCommissions,
		createCommission:   createCommission,
		updateCommission:   updateCommission,
		listPendingItems:   listPendingItems,
		createPendingItem:  createPendingItem,
		resolvePendingItem: resolvePendingItem,
		listRenegotiations: listRenegotiations,
		applyRenegotiation: applyRenegotiation,
		cancelSale:         cancelSale,
	}
}

func (handler OperationsHandler) ListInstallments(writer http.ResponseWriter, request *http.Request) {
	installments := handler.listInstallments.Execute(request.PathValue("publicId"))
	response := make([]dto.InstallmentResponse, 0, len(installments))
	for _, installment := range installments {
		response = append(response, mapInstallment(installment))
	}
	writeJSON(writer, http.StatusOK, response)
}

func (handler OperationsHandler) CreateInstallments(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateInstallmentsRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	input := command.CreateInstallmentScheduleInput{SalePublicID: request.PathValue("publicId")}
	for _, item := range payload.Installments {
		input.Installments = append(input.Installments, command.InstallmentPlanInput{AmountCents: item.AmountCents, DueDate: item.DueDate})
	}

	result := handler.createInstallments.Execute(input)
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}
	if result.Conflict {
		writeError(writer, http.StatusConflict, result.ErrorCode, result.ErrorText)
		return
	}

	response := make([]dto.InstallmentResponse, 0, len(result.Resource))
	for _, installment := range result.Resource {
		response = append(response, mapInstallment(installment))
	}
	writeJSON(writer, http.StatusCreated, response)
}

func (handler OperationsHandler) ListCommissions(writer http.ResponseWriter, request *http.Request) {
	commissions := handler.listCommissions.Execute(request.PathValue("publicId"))
	response := make([]dto.CommissionResponse, 0, len(commissions))
	for _, commission := range commissions {
		response = append(response, mapCommission(commission))
	}
	writeJSON(writer, http.StatusOK, response)
}

func (handler OperationsHandler) CreateCommission(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateCommissionRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createCommission.Execute(command.CreateCommissionInput{
		SalePublicID:    request.PathValue("publicId"),
		RecipientUserID: payload.RecipientUserID,
		RoleCode:        payload.RoleCode,
		RateBps:         payload.RateBps,
	})
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}
	writeJSON(writer, http.StatusCreated, mapCommission(result.Resource))
}

func (handler OperationsHandler) BlockCommission(writer http.ResponseWriter, request *http.Request) {
	handler.writeCommissionTransition(writer, request, "blocked")
}

func (handler OperationsHandler) ReleaseCommission(writer http.ResponseWriter, request *http.Request) {
	handler.writeCommissionTransition(writer, request, "released")
}

func (handler OperationsHandler) writeCommissionTransition(writer http.ResponseWriter, request *http.Request, status string) {
	result := handler.updateCommission.Execute(request.PathValue("publicId"), request.PathValue("commissionPublicId"), status)
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}
	writeJSON(writer, http.StatusOK, mapCommission(result.Resource))
}

func (handler OperationsHandler) ListPendingItems(writer http.ResponseWriter, request *http.Request) {
	items := handler.listPendingItems.Execute(request.PathValue("publicId"))
	response := make([]dto.PendingItemResponse, 0, len(items))
	for _, item := range items {
		response = append(response, mapPendingItem(item))
	}
	writeJSON(writer, http.StatusOK, response)
}

func (handler OperationsHandler) CreatePendingItem(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreatePendingItemRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.createPendingItem.Execute(request.PathValue("publicId"), payload.Code, payload.Summary)
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}
	writeJSON(writer, http.StatusCreated, mapPendingItem(result.Resource))
}

func (handler OperationsHandler) ResolvePendingItem(writer http.ResponseWriter, request *http.Request) {
	result := handler.resolvePendingItem.Execute(request.PathValue("publicId"), request.PathValue("pendingItemPublicId"))
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}
	writeJSON(writer, http.StatusOK, mapPendingItem(result.Resource))
}

func (handler OperationsHandler) ListRenegotiations(writer http.ResponseWriter, request *http.Request) {
	renegotiations := handler.listRenegotiations.Execute(request.PathValue("publicId"))
	response := make([]dto.RenegotiationResponse, 0, len(renegotiations))
	for _, renegotiation := range renegotiations {
		response = append(response, mapRenegotiation(renegotiation))
	}
	writeJSON(writer, http.StatusOK, response)
}

func (handler OperationsHandler) ApplyRenegotiation(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CreateRenegotiationRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.applyRenegotiation.Execute(request.PathValue("publicId"), payload.Reason, payload.NewAmountCents)
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}
	writeJSON(writer, http.StatusCreated, mapRenegotiation(result.Resource))
}

func (handler OperationsHandler) CancelSale(writer http.ResponseWriter, request *http.Request) {
	var payload dto.CancelSaleRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	result := handler.cancelSale.Execute(request.PathValue("publicId"), payload.Reason)
	if result.BadRequest {
		writeError(writer, http.StatusBadRequest, result.ErrorCode, result.ErrorText)
		return
	}
	if result.NotFound {
		writeError(writer, http.StatusNotFound, result.ErrorCode, result.ErrorText)
		return
	}

	writeJSON(writer, http.StatusOK, dto.SaleCancellationResponse{
		Sale:   mapSale(result.Resource),
		Reason: payload.Reason,
	})
}
