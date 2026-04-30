package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrContractPublicIDInvalid    = errors.New("contract public id is invalid")
	ErrContractCustomerInvalid    = errors.New("contract customer public id is invalid")
	ErrContractTitleRequired      = errors.New("contract title is required")
	ErrContractAmountInvalid      = errors.New("contract amount is invalid")
	ErrContractBillingDayInvalid  = errors.New("contract billing day is invalid")
	ErrContractDatesInvalid       = errors.New("contract dates are invalid")
	ErrContractStatusInvalid      = errors.New("contract status is invalid")
	ErrContractAdjustmentInvalid  = errors.New("contract adjustment is invalid")
	ErrContractTerminationInvalid = errors.New("contract termination is invalid")
	ErrContractAlreadyTerminated  = errors.New("contract is already terminated")
	ErrContractReasonRequired     = errors.New("contract reason is required")
	ErrContractRecordedByRequired = errors.New("contract recorded by is required")
	ErrChargeStatusInvalid        = errors.New("charge status is invalid")
	ErrChargeTransitionInvalid    = errors.New("charge transition is invalid")
	ErrChargePaymentReference     = errors.New("charge payment reference is required")
	ErrChargePaidAtInvalid        = errors.New("charge paid at is invalid")
)

type Contract struct {
	PublicID          string
	TenantSlug        string
	CustomerPublicID  string
	Title             string
	PropertyCode      string
	CurrencyCode      string
	AmountCents       int64
	BillingDay        int
	StartsAt          time.Time
	EndsAt            time.Time
	Status            string
	TerminatedAt      *time.Time
	TerminationReason string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Charge struct {
	PublicID         string
	ContractPublicID string
	DueDate          time.Time
	AmountCents      int64
	Status           string
	PaidAt           *time.Time
	PaymentReference string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Adjustment struct {
	PublicID            string
	ContractPublicID    string
	EffectiveAt         time.Time
	PreviousAmountCents int64
	NewAmountCents      int64
	Reason              string
	RecordedBy          string
	CreatedAt           time.Time
}

type Event struct {
	PublicID         string
	ContractPublicID string
	EventCode        string
	Summary          string
	RecordedBy       string
	CreatedAt        time.Time
}

type OutboxEvent struct {
	PublicID          string
	AggregateType     string
	AggregatePublicID string
	EventType         string
	Payload           string
	Status            string
	CreatedAt         time.Time
}

func NewContract(
	publicID string,
	tenantSlug string,
	customerPublicID string,
	title string,
	propertyCode string,
	currencyCode string,
	amountCents int64,
	billingDay int,
	startsAt time.Time,
	endsAt time.Time,
	status string,
	terminatedAt *time.Time,
	terminationReason string,
	createdAt time.Time,
	updatedAt time.Time,
) (Contract, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Contract{}, ErrContractPublicIDInvalid
	}

	normalizedCustomerPublicID := strings.TrimSpace(customerPublicID)
	if _, err := uuid.Parse(normalizedCustomerPublicID); err != nil {
		return Contract{}, ErrContractCustomerInvalid
	}

	normalizedTitle := strings.TrimSpace(title)
	if normalizedTitle == "" {
		return Contract{}, ErrContractTitleRequired
	}
	if amountCents <= 0 {
		return Contract{}, ErrContractAmountInvalid
	}
	if billingDay < 1 || billingDay > 31 {
		return Contract{}, ErrContractBillingDayInvalid
	}

	normalizedStart := normalizeDate(startsAt)
	normalizedEnd := normalizeDate(endsAt)
	if normalizedStart.IsZero() || normalizedEnd.IsZero() || normalizedEnd.Before(normalizedStart) {
		return Contract{}, ErrContractDatesInvalid
	}

	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	if normalizedStatus == "" {
		normalizedStatus = "active"
	}
	if normalizedStatus != "active" && normalizedStatus != "terminated" {
		return Contract{}, ErrContractStatusInvalid
	}

	normalizedTenantSlug := strings.ToLower(strings.TrimSpace(tenantSlug))
	if normalizedTenantSlug == "" {
		normalizedTenantSlug = "bootstrap-ops"
	}

	normalizedCurrency := strings.ToUpper(strings.TrimSpace(currencyCode))
	if normalizedCurrency == "" {
		normalizedCurrency = "BRL"
	}

	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}

	var normalizedTerminatedAt *time.Time
	if terminatedAt != nil && !terminatedAt.IsZero() {
		value := terminatedAt.UTC()
		normalizedTerminatedAt = &value
	}

	return Contract{
		PublicID:          normalizedPublicID,
		TenantSlug:        normalizedTenantSlug,
		CustomerPublicID:  normalizedCustomerPublicID,
		Title:             normalizedTitle,
		PropertyCode:      strings.TrimSpace(propertyCode),
		CurrencyCode:      normalizedCurrency,
		AmountCents:       amountCents,
		BillingDay:        billingDay,
		StartsAt:          normalizedStart,
		EndsAt:            normalizedEnd,
		Status:            normalizedStatus,
		TerminatedAt:      normalizedTerminatedAt,
		TerminationReason: strings.TrimSpace(terminationReason),
		CreatedAt:         createdAt.UTC(),
		UpdatedAt:         updatedAt.UTC(),
	}, nil
}

func (contract Contract) BuildCharges() []Charge {
	charges := make([]Charge, 0)
	now := time.Now().UTC()
	cursor := time.Date(contract.StartsAt.Year(), contract.StartsAt.Month(), 1, 0, 0, 0, 0, time.UTC)
	limit := time.Date(contract.EndsAt.Year(), contract.EndsAt.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !cursor.After(limit) {
		dueDate := dueDateForMonth(cursor.Year(), cursor.Month(), contract.BillingDay)
		if dueDate.Before(contract.StartsAt) || dueDate.After(contract.EndsAt) {
			cursor = cursor.AddDate(0, 1, 0)
			continue
		}

		charges = append(charges, Charge{
			PublicID:         uuid.NewString(),
			ContractPublicID: contract.PublicID,
			DueDate:          dueDate,
			AmountCents:      contract.AmountCents,
			Status:           "scheduled",
			CreatedAt:        now,
			UpdatedAt:        now,
		})
		cursor = cursor.AddDate(0, 1, 0)
	}

	return charges
}

func (contract Contract) BuildCreatedArtifacts(recordedBy string, charges []Charge) (Event, OutboxEvent, error) {
	if strings.TrimSpace(recordedBy) == "" {
		return Event{}, OutboxEvent{}, ErrContractRecordedByRequired
	}

	event := Event{
		PublicID:         uuid.NewString(),
		ContractPublicID: contract.PublicID,
		EventCode:        "contract_created",
		Summary:          fmt.Sprintf("Contrato criado com %d cobrancas recorrentes.", len(charges)),
		RecordedBy:       strings.TrimSpace(recordedBy),
		CreatedAt:        time.Now().UTC(),
	}

	payload, _ := json.Marshal(map[string]any{
		"contractPublicId": contract.PublicID,
		"customerPublicId": contract.CustomerPublicID,
		"status":           contract.Status,
		"amountCents":      contract.AmountCents,
		"charges":          len(charges),
	})

	outbox := OutboxEvent{
		PublicID:          uuid.NewString(),
		AggregateType:     "rental.contract",
		AggregatePublicID: contract.PublicID,
		EventType:         "rental.contract.created",
		Payload:           string(payload),
		Status:            "pending",
		CreatedAt:         time.Now().UTC(),
	}

	return event, outbox, nil
}

func (contract Contract) ApplyAdjustment(effectiveAt time.Time, newAmountCents int64, reason string, recordedBy string, charges []Charge) (Contract, Adjustment, []Charge, Event, OutboxEvent, error) {
	if contract.Status != "active" {
		return Contract{}, Adjustment{}, nil, Event{}, OutboxEvent{}, ErrContractAlreadyTerminated
	}
	if strings.TrimSpace(recordedBy) == "" {
		return Contract{}, Adjustment{}, nil, Event{}, OutboxEvent{}, ErrContractRecordedByRequired
	}
	if strings.TrimSpace(reason) == "" {
		return Contract{}, Adjustment{}, nil, Event{}, OutboxEvent{}, ErrContractReasonRequired
	}
	if newAmountCents <= 0 || newAmountCents == contract.AmountCents {
		return Contract{}, Adjustment{}, nil, Event{}, OutboxEvent{}, ErrContractAdjustmentInvalid
	}

	normalizedEffectiveAt := normalizeDate(effectiveAt)
	if normalizedEffectiveAt.IsZero() || normalizedEffectiveAt.After(contract.EndsAt) {
		return Contract{}, Adjustment{}, nil, Event{}, OutboxEvent{}, ErrContractAdjustmentInvalid
	}

	updatedContract := contract
	updatedContract.AmountCents = newAmountCents
	updatedContract.UpdatedAt = time.Now().UTC()

	updatedCharges := cloneCharges(charges)
	updatedCount := 0
	for index := range updatedCharges {
		if updatedCharges[index].Status != "scheduled" {
			continue
		}
		if updatedCharges[index].DueDate.Before(normalizedEffectiveAt) {
			continue
		}
		updatedCharges[index].AmountCents = newAmountCents
		updatedCharges[index].UpdatedAt = time.Now().UTC()
		updatedCount++
	}

	if updatedCount == 0 {
		return Contract{}, Adjustment{}, nil, Event{}, OutboxEvent{}, ErrContractAdjustmentInvalid
	}

	adjustment := Adjustment{
		PublicID:            uuid.NewString(),
		ContractPublicID:    contract.PublicID,
		EffectiveAt:         normalizedEffectiveAt,
		PreviousAmountCents: contract.AmountCents,
		NewAmountCents:      newAmountCents,
		Reason:              strings.TrimSpace(reason),
		RecordedBy:          strings.TrimSpace(recordedBy),
		CreatedAt:           time.Now().UTC(),
	}

	event := Event{
		PublicID:         uuid.NewString(),
		ContractPublicID: contract.PublicID,
		EventCode:        "contract_adjusted",
		Summary:          fmt.Sprintf("Contrato reajustado para %d centavos com efeito em %s.", newAmountCents, normalizedEffectiveAt.Format("2006-01-02")),
		RecordedBy:       strings.TrimSpace(recordedBy),
		CreatedAt:        time.Now().UTC(),
	}

	payload, _ := json.Marshal(map[string]any{
		"contractPublicId": contract.PublicID,
		"effectiveAt":      normalizedEffectiveAt.Format("2006-01-02"),
		"newAmountCents":   newAmountCents,
		"updatedCharges":   updatedCount,
	})

	outbox := OutboxEvent{
		PublicID:          uuid.NewString(),
		AggregateType:     "rental.contract",
		AggregatePublicID: contract.PublicID,
		EventType:         "rental.contract.adjusted",
		Payload:           string(payload),
		Status:            "pending",
		CreatedAt:         time.Now().UTC(),
	}

	return updatedContract, adjustment, updatedCharges, event, outbox, nil
}

func (contract Contract) Terminate(effectiveAt time.Time, reason string, recordedBy string, charges []Charge) (Contract, []Charge, Event, OutboxEvent, error) {
	if contract.Status != "active" {
		return Contract{}, nil, Event{}, OutboxEvent{}, ErrContractAlreadyTerminated
	}
	if strings.TrimSpace(recordedBy) == "" {
		return Contract{}, nil, Event{}, OutboxEvent{}, ErrContractRecordedByRequired
	}
	if strings.TrimSpace(reason) == "" {
		return Contract{}, nil, Event{}, OutboxEvent{}, ErrContractReasonRequired
	}

	normalizedEffectiveAt := normalizeDate(effectiveAt)
	if normalizedEffectiveAt.IsZero() || normalizedEffectiveAt.Before(contract.StartsAt) || normalizedEffectiveAt.After(contract.EndsAt) {
		return Contract{}, nil, Event{}, OutboxEvent{}, ErrContractTerminationInvalid
	}

	updatedContract := contract
	updatedContract.Status = "terminated"
	updatedContract.EndsAt = normalizedEffectiveAt
	updatedContract.UpdatedAt = time.Now().UTC()
	updatedContract.TerminationReason = strings.TrimSpace(reason)
	updatedContract.TerminatedAt = &updatedContract.UpdatedAt

	updatedCharges := cloneCharges(charges)
	cancelledCount := 0
	for index := range updatedCharges {
		if updatedCharges[index].Status != "scheduled" {
			continue
		}
		if !updatedCharges[index].DueDate.After(normalizedEffectiveAt) {
			continue
		}
		updatedCharges[index].Status = "cancelled"
		updatedCharges[index].UpdatedAt = time.Now().UTC()
		cancelledCount++
	}

	event := Event{
		PublicID:         uuid.NewString(),
		ContractPublicID: contract.PublicID,
		EventCode:        "contract_terminated",
		Summary:          fmt.Sprintf("Contrato rescindido com efeito em %s e %d cobrancas canceladas.", normalizedEffectiveAt.Format("2006-01-02"), cancelledCount),
		RecordedBy:       strings.TrimSpace(recordedBy),
		CreatedAt:        time.Now().UTC(),
	}

	payload, _ := json.Marshal(map[string]any{
		"contractPublicId": contract.PublicID,
		"effectiveAt":      normalizedEffectiveAt.Format("2006-01-02"),
		"cancelledCharges": cancelledCount,
	})

	outbox := OutboxEvent{
		PublicID:          uuid.NewString(),
		AggregateType:     "rental.contract",
		AggregatePublicID: contract.PublicID,
		EventType:         "rental.contract.terminated",
		Payload:           string(payload),
		Status:            "pending",
		CreatedAt:         time.Now().UTC(),
	}

	return updatedContract, updatedCharges, event, outbox, nil
}

func (charge Charge) UpdateStatus(nextStatus string, recordedBy string, paidAt time.Time, paymentReference string) (Charge, Event, OutboxEvent, error) {
	if strings.TrimSpace(recordedBy) == "" {
		return Charge{}, Event{}, OutboxEvent{}, ErrContractRecordedByRequired
	}

	normalizedStatus := strings.ToLower(strings.TrimSpace(nextStatus))
	if normalizedStatus != "paid" && normalizedStatus != "cancelled" {
		return Charge{}, Event{}, OutboxEvent{}, ErrChargeStatusInvalid
	}
	if charge.Status != "scheduled" {
		return Charge{}, Event{}, OutboxEvent{}, ErrChargeTransitionInvalid
	}

	updatedCharge := charge
	updatedCharge.Status = normalizedStatus
	updatedCharge.UpdatedAt = time.Now().UTC()
	updatedCharge.PaidAt = nil
	updatedCharge.PaymentReference = ""

	summary := fmt.Sprintf("Cobranca com vencimento em %s cancelada manualmente.", charge.DueDate.Format("2006-01-02"))
	eventCode := "charge_cancelled"
	eventType := "rental.contract.charge_cancelled"
	payload := map[string]any{
		"contractPublicId": charge.ContractPublicID,
		"chargePublicId":   charge.PublicID,
		"status":           normalizedStatus,
		"amountCents":      charge.AmountCents,
		"dueDate":          charge.DueDate.Format("2006-01-02"),
	}

	if normalizedStatus == "paid" {
		if paidAt.IsZero() {
			return Charge{}, Event{}, OutboxEvent{}, ErrChargePaidAtInvalid
		}

		normalizedPaidAt := paidAt.UTC()
		reference := strings.TrimSpace(paymentReference)
		if reference == "" {
			return Charge{}, Event{}, OutboxEvent{}, ErrChargePaymentReference
		}

		updatedCharge.PaidAt = &normalizedPaidAt
		updatedCharge.PaymentReference = reference
		summary = fmt.Sprintf("Cobranca marcada como paga com referencia %s.", reference)
		eventCode = "charge_paid"
		eventType = "rental.contract.charge_paid"
		payload["paidAt"] = normalizedPaidAt.Format(time.RFC3339)
		payload["paymentReference"] = reference
	}

	event := Event{
		PublicID:         uuid.NewString(),
		ContractPublicID: charge.ContractPublicID,
		EventCode:        eventCode,
		Summary:          summary,
		RecordedBy:       strings.TrimSpace(recordedBy),
		CreatedAt:        time.Now().UTC(),
	}

	rawPayload, _ := json.Marshal(payload)
	outbox := OutboxEvent{
		PublicID:          uuid.NewString(),
		AggregateType:     "rental.contract",
		AggregatePublicID: charge.ContractPublicID,
		EventType:         eventType,
		Payload:           string(rawPayload),
		Status:            "pending",
		CreatedAt:         time.Now().UTC(),
	}

	return updatedCharge, event, outbox, nil
}

func cloneCharges(charges []Charge) []Charge {
	response := make([]Charge, len(charges))
	copy(response, charges)
	return response
}

func normalizeDate(value time.Time) time.Time {
	if value.IsZero() {
		return time.Time{}
	}

	return time.Date(value.UTC().Year(), value.UTC().Month(), value.UTC().Day(), 0, 0, 0, 0, time.UTC)
}

func dueDateForMonth(year int, month time.Month, billingDay int) time.Time {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	day := billingDay
	if day > lastDay {
		day = lastDay
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
