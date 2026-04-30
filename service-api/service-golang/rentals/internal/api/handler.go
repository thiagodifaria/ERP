package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
)

type ContractHandler struct {
	repository  repository.ContractRepository
	attachments repository.AttachmentGateway
}

func NewContractHandler(repository repository.ContractRepository, attachments repository.AttachmentGateway) ContractHandler {
	return ContractHandler{repository: repository, attachments: attachments}
}

func (handler ContractHandler) ListContracts(writer http.ResponseWriter, request *http.Request) {
	contracts := handler.repository.List(repository.ContractFilters{
		TenantSlug:       request.URL.Query().Get("tenantSlug"),
		Status:           request.URL.Query().Get("status"),
		CustomerPublicID: request.URL.Query().Get("customerPublicId"),
	})

	response := make([]ContractResponse, 0, len(contracts))
	for _, contract := range contracts {
		response = append(response, mapContract(contract))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler ContractHandler) GetSummary(writer http.ResponseWriter, request *http.Request) {
	summary := handler.repository.Summary(request.URL.Query().Get("tenantSlug"))
	writeJSON(writer, http.StatusOK, SummaryResponse{
		TenantSlug:           summary.TenantSlug,
		TotalContracts:       summary.TotalContracts,
		ActiveContracts:      summary.ActiveContracts,
		TerminatedContracts:  summary.TerminatedContracts,
		ScheduledCharges:     summary.ScheduledCharges,
		PaidCharges:          summary.PaidCharges,
		CancelledCharges:     summary.CancelledCharges,
		Adjustments:          summary.Adjustments,
		HistoryEvents:        summary.HistoryEvents,
		PendingOutbox:        summary.PendingOutbox,
		ScheduledAmountCents: summary.ScheduledAmountCents,
		PaidAmountCents:      summary.PaidAmountCents,
		CancelledAmountCents: summary.CancelledAmountCents,
	})
}

func (handler ContractHandler) CreateContract(writer http.ResponseWriter, request *http.Request) {
	var payload CreateContractRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	startsAt, ok := parseDate(payload.StartsAt)
	if !ok {
		writeError(writer, http.StatusBadRequest, "invalid_contract_dates", "Contract dates are invalid.")
		return
	}
	endsAt, ok := parseDate(payload.EndsAt)
	if !ok {
		writeError(writer, http.StatusBadRequest, "invalid_contract_dates", "Contract dates are invalid.")
		return
	}

	contract, err := entity.NewContract(
		uuid.NewString(),
		payload.TenantSlug,
		payload.CustomerPublicID,
		payload.Title,
		payload.PropertyCode,
		payload.CurrencyCode,
		payload.AmountCents,
		payload.BillingDay,
		startsAt,
		endsAt,
		"active",
		nil,
		"",
		time.Now().UTC(),
		time.Now().UTC(),
	)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_contract", err.Error())
		return
	}

	charges := contract.BuildCharges()
	event, outbox, err := contract.BuildCreatedArtifacts(payload.RecordedBy, charges)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_contract", err.Error())
		return
	}

	created := handler.repository.Create(contract, charges, event, outbox)
	writeJSON(writer, http.StatusCreated, mapContract(created))
}

func (handler ContractHandler) GetContract(writer http.ResponseWriter, request *http.Request) {
	contract, ok := handler.repository.FindByPublicID(request.URL.Query().Get("tenantSlug"), request.PathValue("publicId"))
	if !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapContract(contract))
}

func (handler ContractHandler) ListCharges(writer http.ResponseWriter, request *http.Request) {
	tenantSlug := request.URL.Query().Get("tenantSlug")
	publicID := request.PathValue("publicId")
	if _, ok := handler.repository.FindByPublicID(tenantSlug, publicID); !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	charges := handler.repository.ListCharges(tenantSlug, publicID, request.URL.Query().Get("status"))
	response := make([]ChargeResponse, 0, len(charges))
	for _, charge := range charges {
		response = append(response, mapCharge(charge))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler ContractHandler) UpdateChargeStatus(writer http.ResponseWriter, request *http.Request) {
	var payload UpdateChargeStatusRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	tenantSlug := payload.TenantSlug
	contractPublicID := request.PathValue("publicId")
	if _, ok := handler.repository.FindByPublicID(tenantSlug, contractPublicID); !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	charge, ok := handler.repository.FindChargeByPublicID(tenantSlug, contractPublicID, request.PathValue("chargePublicId"))
	if !ok {
		writeError(writer, http.StatusNotFound, "charge_not_found", "Charge was not found.")
		return
	}

	paidAt := time.Time{}
	if strings.TrimSpace(payload.PaidAt) != "" {
		var parsed bool
		paidAt, parsed = parseTimestamp(payload.PaidAt)
		if !parsed {
			writeError(writer, http.StatusBadRequest, "invalid_charge_status", "Charge paid at is invalid.")
			return
		}
	}

	updatedCharge, event, outbox, err := charge.UpdateStatus(payload.Status, payload.RecordedBy, paidAt, payload.PaymentReference)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_charge_status", err.Error())
		return
	}

	savedCharge, saved := handler.repository.SaveChargeStatus(tenantSlug, updatedCharge, event, outbox)
	if !saved {
		writeError(writer, http.StatusNotFound, "charge_not_found", "Charge was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapCharge(savedCharge))
}

func (handler ContractHandler) ListHistory(writer http.ResponseWriter, request *http.Request) {
	tenantSlug := request.URL.Query().Get("tenantSlug")
	publicID := request.PathValue("publicId")
	if _, ok := handler.repository.FindByPublicID(tenantSlug, publicID); !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	events := handler.repository.ListEvents(tenantSlug, publicID)
	response := make([]EventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, mapEvent(event))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler ContractHandler) ListAdjustments(writer http.ResponseWriter, request *http.Request) {
	tenantSlug := request.URL.Query().Get("tenantSlug")
	publicID := request.PathValue("publicId")
	if _, ok := handler.repository.FindByPublicID(tenantSlug, publicID); !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	adjustments := handler.repository.ListAdjustments(tenantSlug, publicID)
	response := make([]AdjustmentResponse, 0, len(adjustments))
	for _, adjustment := range adjustments {
		response = append(response, mapAdjustment(adjustment))
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler ContractHandler) CreateAdjustment(writer http.ResponseWriter, request *http.Request) {
	var payload CreateAdjustmentRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	tenantSlug := payload.TenantSlug
	publicID := request.PathValue("publicId")
	contract, ok := handler.repository.FindByPublicID(tenantSlug, publicID)
	if !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	effectiveAt, ok := parseDate(payload.EffectiveAt)
	if !ok {
		writeError(writer, http.StatusBadRequest, "invalid_adjustment", "Adjustment effective date is invalid.")
		return
	}

	charges := handler.repository.ListCharges(tenantSlug, publicID, "")
	updatedContract, adjustment, updatedCharges, event, outbox, err := contract.ApplyAdjustment(effectiveAt, payload.NewAmountCents, payload.Reason, payload.RecordedBy, charges)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_adjustment", err.Error())
		return
	}

	savedContract, saved := handler.repository.SaveAdjustment(tenantSlug, updatedContract, adjustment, updatedCharges, event, outbox)
	if !saved {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapContract(savedContract))
}

func (handler ContractHandler) TerminateContract(writer http.ResponseWriter, request *http.Request) {
	var payload TerminateContractRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	tenantSlug := payload.TenantSlug
	publicID := request.PathValue("publicId")
	contract, ok := handler.repository.FindByPublicID(tenantSlug, publicID)
	if !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	effectiveAt, ok := parseDate(payload.EffectiveAt)
	if !ok {
		writeError(writer, http.StatusBadRequest, "invalid_termination", "Termination effective date is invalid.")
		return
	}

	charges := handler.repository.ListCharges(tenantSlug, publicID, "")
	updatedContract, updatedCharges, event, outbox, err := contract.Terminate(effectiveAt, payload.Reason, payload.RecordedBy, charges)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_termination", err.Error())
		return
	}

	savedContract, saved := handler.repository.SaveTermination(tenantSlug, updatedContract, updatedCharges, event, outbox)
	if !saved {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	writeJSON(writer, http.StatusOK, mapContract(savedContract))
}

func (handler ContractHandler) ListAttachments(writer http.ResponseWriter, request *http.Request) {
	tenantSlug := request.URL.Query().Get("tenantSlug")
	publicID := request.PathValue("publicId")
	if _, ok := handler.repository.FindByPublicID(tenantSlug, publicID); !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	attachments, err := handler.attachments.List(resolveTenant(tenantSlug), "rentals.contract", publicID)
	if err != nil {
		writeError(writer, http.StatusBadGateway, "documents_unavailable", "Documents service is unavailable.")
		return
	}

	response := make([]AttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		response = append(response, AttachmentResponse{
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

func (handler ContractHandler) CreateAttachment(writer http.ResponseWriter, request *http.Request) {
	var payload CreateAttachmentRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_json", "Request body is invalid.")
		return
	}

	tenantSlug := request.URL.Query().Get("tenantSlug")
	publicID := request.PathValue("publicId")
	if _, ok := handler.repository.FindByPublicID(tenantSlug, publicID); !ok {
		writeError(writer, http.StatusNotFound, "contract_not_found", "Contract was not found.")
		return
	}

	attachment, err := handler.attachments.Create(repository.CreateAttachmentInput{
		TenantSlug:    resolveTenant(tenantSlug),
		OwnerType:     "rentals.contract",
		OwnerPublicID: publicID,
		FileName:      payload.FileName,
		ContentType:   payload.ContentType,
		StorageKey:    payload.StorageKey,
		StorageDriver: payload.StorageDriver,
		Source:        payload.Source,
		UploadedBy:    payload.UploadedBy,
	})
	if err != nil {
		writeError(writer, http.StatusBadGateway, "documents_unavailable", "Documents service is unavailable.")
		return
	}

	writeJSON(writer, http.StatusCreated, AttachmentResponse{
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

func parseDate(value string) (time.Time, bool) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, false
	}
	return parsed.UTC(), true
}

func parseTimestamp(value string) (time.Time, bool) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, false
	}
	return parsed.UTC(), true
}

func resolveTenant(tenantSlug string) string {
	normalized := strings.ToLower(strings.TrimSpace(tenantSlug))
	if normalized == "" {
		return "bootstrap-ops"
	}
	return normalized
}

func mapContract(contract entity.Contract) ContractResponse {
	return ContractResponse{
		PublicID:          contract.PublicID,
		TenantSlug:        contract.TenantSlug,
		CustomerPublicID:  contract.CustomerPublicID,
		Title:             contract.Title,
		PropertyCode:      contract.PropertyCode,
		CurrencyCode:      contract.CurrencyCode,
		AmountCents:       contract.AmountCents,
		BillingDay:        contract.BillingDay,
		StartsAt:          contract.StartsAt.Format("2006-01-02"),
		EndsAt:            contract.EndsAt.Format("2006-01-02"),
		Status:            contract.Status,
		TerminatedAt:      contract.TerminatedAt,
		TerminationReason: contract.TerminationReason,
		CreatedAt:         contract.CreatedAt,
		UpdatedAt:         contract.UpdatedAt,
	}
}

func mapCharge(charge entity.Charge) ChargeResponse {
	return ChargeResponse{
		PublicID:         charge.PublicID,
		ContractPublicID: charge.ContractPublicID,
		DueDate:          charge.DueDate.Format("2006-01-02"),
		AmountCents:      charge.AmountCents,
		Status:           charge.Status,
		PaidAt:           charge.PaidAt,
		PaymentReference: charge.PaymentReference,
		CreatedAt:        charge.CreatedAt,
		UpdatedAt:        charge.UpdatedAt,
	}
}

func mapAdjustment(adjustment entity.Adjustment) AdjustmentResponse {
	return AdjustmentResponse{
		PublicID:            adjustment.PublicID,
		ContractPublicID:    adjustment.ContractPublicID,
		EffectiveAt:         adjustment.EffectiveAt.Format("2006-01-02"),
		PreviousAmountCents: adjustment.PreviousAmountCents,
		NewAmountCents:      adjustment.NewAmountCents,
		Reason:              adjustment.Reason,
		RecordedBy:          adjustment.RecordedBy,
		CreatedAt:           adjustment.CreatedAt,
	}
}

func mapEvent(event entity.Event) EventResponse {
	return EventResponse{
		PublicID:         event.PublicID,
		ContractPublicID: event.ContractPublicID,
		EventCode:        event.EventCode,
		Summary:          event.Summary,
		RecordedBy:       event.RecordedBy,
		CreatedAt:        event.CreatedAt,
	}
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func writeError(writer http.ResponseWriter, status int, code string, message string) {
	writeJSON(writer, status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}
