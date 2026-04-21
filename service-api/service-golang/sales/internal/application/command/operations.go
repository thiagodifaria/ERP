package command

import (
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type OperationResult[T any] struct {
	Resource   T
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
	Conflict   bool
}

type CreateInstallmentSchedule struct {
	saleRepository        repository.SaleRepository
	installmentRepository repository.InstallmentRepository
	eventRepository       repository.CommercialEventRepository
	outboxRepository      repository.OutboxEventRepository
}

type CreateInstallmentScheduleInput struct {
	SalePublicID string
	Installments []InstallmentPlanInput
}

type InstallmentPlanInput struct {
	AmountCents int64
	DueDate     string
}

func NewCreateInstallmentSchedule(
	saleRepository repository.SaleRepository,
	installmentRepository repository.InstallmentRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) CreateInstallmentSchedule {
	return CreateInstallmentSchedule{
		saleRepository:        saleRepository,
		installmentRepository: installmentRepository,
		eventRepository:       eventRepository,
		outboxRepository:      outboxRepository,
	}
}

func (useCase CreateInstallmentSchedule) Execute(input CreateInstallmentScheduleInput) OperationResult[[]entity.Installment] {
	sale := useCase.saleRepository.FindByPublicID(input.SalePublicID)
	if sale == nil {
		return OperationResult[[]entity.Installment]{ErrorCode: "sale_not_found", ErrorText: "Sale was not found.", NotFound: true}
	}

	if sale.Status == "cancelled" {
		return OperationResult[[]entity.Installment]{ErrorCode: "sale_not_billable", ErrorText: "Cancelled sales cannot receive installments.", BadRequest: true}
	}

	if len(useCase.installmentRepository.ListBySalePublicID(sale.PublicID)) > 0 {
		return OperationResult[[]entity.Installment]{ErrorCode: "installments_already_defined", ErrorText: "Installments were already defined for this sale.", Conflict: true}
	}

	if len(input.Installments) == 0 {
		return OperationResult[[]entity.Installment]{ErrorCode: "invalid_installments", ErrorText: "At least one installment is required.", BadRequest: true}
	}

	var totalAmount int64
	created := make([]entity.Installment, 0, len(input.Installments))
	for index, plan := range input.Installments {
		installment, err := entity.NewInstallment(newPublicID(), sale.PublicID, index+1, plan.AmountCents, plan.DueDate)
		if err != nil {
			return OperationResult[[]entity.Installment]{ErrorCode: "invalid_installments", ErrorText: "Installment payload is invalid.", BadRequest: true}
		}

		totalAmount += installment.AmountCents
		created = append(created, useCase.installmentRepository.Save(installment))
	}

	if totalAmount != sale.AmountCents {
		return OperationResult[[]entity.Installment]{ErrorCode: "installment_total_mismatch", ErrorText: "Installment total must match sale amount.", BadRequest: true}
	}

	recordCommercialEvent(useCase.eventRepository, "sale", sale.PublicID, "sale_installments_defined", "sales", "Commercial installment schedule defined.")
	appendOutboxEvent(useCase.outboxRepository, "sale", sale.PublicID, "sale.installments_defined", map[string]any{
		"salePublicId":  sale.PublicID,
		"installments": len(created),
		"amountCents":   sale.AmountCents,
	})
	return OperationResult[[]entity.Installment]{Resource: created}
}

type CreateCommission struct {
	saleRepository       repository.SaleRepository
	commissionRepository repository.CommissionRepository
	eventRepository      repository.CommercialEventRepository
}

type CreateCommissionInput struct {
	SalePublicID    string
	RecipientUserID string
	RoleCode        string
	RateBps         int
}

func NewCreateCommission(
	saleRepository repository.SaleRepository,
	commissionRepository repository.CommissionRepository,
	eventRepository repository.CommercialEventRepository,
) CreateCommission {
	return CreateCommission{
		saleRepository:       saleRepository,
		commissionRepository: commissionRepository,
		eventRepository:      eventRepository,
	}
}

func (useCase CreateCommission) Execute(input CreateCommissionInput) OperationResult[entity.Commission] {
	sale := useCase.saleRepository.FindByPublicID(input.SalePublicID)
	if sale == nil {
		return OperationResult[entity.Commission]{ErrorCode: "sale_not_found", ErrorText: "Sale was not found.", NotFound: true}
	}

	commission, err := entity.NewCommission(newPublicID(), sale.PublicID, input.RecipientUserID, input.RoleCode, input.RateBps, sale.AmountCents)
	if err != nil {
		return OperationResult[entity.Commission]{ErrorCode: "invalid_commission", ErrorText: "Commission payload is invalid.", BadRequest: true}
	}

	var totalCommissionAmount int64
	for _, existing := range useCase.commissionRepository.ListBySalePublicID(sale.PublicID) {
		totalCommissionAmount += existing.AmountCents
	}
	if totalCommissionAmount+commission.AmountCents > sale.AmountCents {
		return OperationResult[entity.Commission]{ErrorCode: "commission_total_exceeded", ErrorText: "Commission footprint cannot exceed sale amount.", BadRequest: true}
	}

	saved := useCase.commissionRepository.Save(commission)
	recordCommercialEvent(useCase.eventRepository, "sale", sale.PublicID, "sale_commission_created", "sales", "Operational commission registered for sale.")
	return OperationResult[entity.Commission]{Resource: saved}
}

type UpdateCommissionStatus struct {
	commissionRepository repository.CommissionRepository
	eventRepository      repository.CommercialEventRepository
}

func NewUpdateCommissionStatus(commissionRepository repository.CommissionRepository, eventRepository repository.CommercialEventRepository) UpdateCommissionStatus {
	return UpdateCommissionStatus{commissionRepository: commissionRepository, eventRepository: eventRepository}
}

func (useCase UpdateCommissionStatus) Execute(salePublicID string, commissionPublicID string, status string) OperationResult[entity.Commission] {
	commission := useCase.commissionRepository.FindByPublicID(commissionPublicID)
	if commission == nil || commission.SalePublicID != salePublicID {
		return OperationResult[entity.Commission]{ErrorCode: "commission_not_found", ErrorText: "Commission was not found.", NotFound: true}
	}

	updated, err := commission.TransitionTo(status)
	if err != nil {
		return OperationResult[entity.Commission]{ErrorCode: "invalid_commission_transition", ErrorText: "Commission status transition is invalid.", BadRequest: true}
	}

	saved := useCase.commissionRepository.Update(updated)
	recordCommercialEvent(useCase.eventRepository, "sale", salePublicID, "sale_commission_status_changed", "sales", "Operational commission transitioned to "+saved.Status+".")
	return OperationResult[entity.Commission]{Resource: saved}
}

type CreatePendingItem struct {
	saleRepository       repository.SaleRepository
	pendingItemRepository repository.PendingItemRepository
	eventRepository      repository.CommercialEventRepository
}

func NewCreatePendingItem(
	saleRepository repository.SaleRepository,
	pendingItemRepository repository.PendingItemRepository,
	eventRepository repository.CommercialEventRepository,
) CreatePendingItem {
	return CreatePendingItem{
		saleRepository:        saleRepository,
		pendingItemRepository: pendingItemRepository,
		eventRepository:       eventRepository,
	}
}

func (useCase CreatePendingItem) Execute(salePublicID string, code string, summary string) OperationResult[entity.PendingItem] {
	sale := useCase.saleRepository.FindByPublicID(salePublicID)
	if sale == nil {
		return OperationResult[entity.PendingItem]{ErrorCode: "sale_not_found", ErrorText: "Sale was not found.", NotFound: true}
	}

	item, err := entity.NewPendingItem(newPublicID(), sale.PublicID, code, summary)
	if err != nil {
		return OperationResult[entity.PendingItem]{ErrorCode: "invalid_pending_item", ErrorText: "Pending item payload is invalid.", BadRequest: true}
	}

	saved := useCase.pendingItemRepository.Save(item)
	recordCommercialEvent(useCase.eventRepository, "sale", sale.PublicID, "sale_pending_item_created", "sales", "Operational pending item created for sale.")
	return OperationResult[entity.PendingItem]{Resource: saved}
}

type ResolvePendingItem struct {
	pendingItemRepository repository.PendingItemRepository
	eventRepository       repository.CommercialEventRepository
}

func NewResolvePendingItem(pendingItemRepository repository.PendingItemRepository, eventRepository repository.CommercialEventRepository) ResolvePendingItem {
	return ResolvePendingItem{pendingItemRepository: pendingItemRepository, eventRepository: eventRepository}
}

func (useCase ResolvePendingItem) Execute(salePublicID string, pendingItemPublicID string) OperationResult[entity.PendingItem] {
	item := useCase.pendingItemRepository.FindByPublicID(pendingItemPublicID)
	if item == nil || item.SalePublicID != salePublicID {
		return OperationResult[entity.PendingItem]{ErrorCode: "pending_item_not_found", ErrorText: "Pending item was not found.", NotFound: true}
	}

	resolved, err := item.Resolve(time.Now().UTC())
	if err != nil {
		return OperationResult[entity.PendingItem]{ErrorCode: "invalid_pending_item_transition", ErrorText: "Pending item cannot be resolved from the current state.", BadRequest: true}
	}

	saved := useCase.pendingItemRepository.Update(resolved)
	recordCommercialEvent(useCase.eventRepository, "sale", salePublicID, "sale_pending_item_resolved", "sales", "Operational pending item resolved.")
	return OperationResult[entity.PendingItem]{Resource: saved}
}

type ApplyRenegotiation struct {
	saleRepository          repository.SaleRepository
	invoiceRepository       repository.InvoiceRepository
	installmentRepository   repository.InstallmentRepository
	commissionRepository    repository.CommissionRepository
	renegotiationRepository repository.RenegotiationRepository
	eventRepository         repository.CommercialEventRepository
	outboxRepository        repository.OutboxEventRepository
}

func NewApplyRenegotiation(
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	installmentRepository repository.InstallmentRepository,
	commissionRepository repository.CommissionRepository,
	renegotiationRepository repository.RenegotiationRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) ApplyRenegotiation {
	return ApplyRenegotiation{
		saleRepository:          saleRepository,
		invoiceRepository:       invoiceRepository,
		installmentRepository:   installmentRepository,
		commissionRepository:    commissionRepository,
		renegotiationRepository: renegotiationRepository,
		eventRepository:         eventRepository,
		outboxRepository:        outboxRepository,
	}
}

func (useCase ApplyRenegotiation) Execute(salePublicID string, reason string, newAmountCents int64) OperationResult[entity.Renegotiation] {
	sale := useCase.saleRepository.FindByPublicID(salePublicID)
	if sale == nil {
		return OperationResult[entity.Renegotiation]{ErrorCode: "sale_not_found", ErrorText: "Sale was not found.", NotFound: true}
	}

	if sale.Status != "active" || useCase.invoiceRepository.FindBySalePublicID(sale.PublicID) != nil || len(useCase.installmentRepository.ListBySalePublicID(sale.PublicID)) > 0 || len(useCase.commissionRepository.ListBySalePublicID(sale.PublicID)) > 0 {
		return OperationResult[entity.Renegotiation]{ErrorCode: "sale_not_renegotiable", ErrorText: "Sale cannot be renegotiated after billing, installments or commissions are defined.", BadRequest: true}
	}

	updatedSale, err := sale.ReviseAmount(newAmountCents)
	if err != nil {
		return OperationResult[entity.Renegotiation]{ErrorCode: "invalid_renegotiation", ErrorText: "Renegotiation amount is invalid.", BadRequest: true}
	}

	savedSale := useCase.saleRepository.Update(updatedSale)
	renegotiation, err := entity.NewAppliedRenegotiation(newPublicID(), sale.PublicID, reason, sale.AmountCents, newAmountCents, time.Now().UTC())
	if err != nil {
		return OperationResult[entity.Renegotiation]{ErrorCode: "invalid_renegotiation", ErrorText: "Renegotiation payload is invalid.", BadRequest: true}
	}

	saved := useCase.renegotiationRepository.Save(renegotiation)
	recordCommercialEvent(useCase.eventRepository, "sale", sale.PublicID, "sale_renegotiated", "sales", "Sale renegotiated to a new commercial amount.")
	appendOutboxEvent(useCase.outboxRepository, "sale", savedSale.PublicID, "sale.renegotiated", map[string]any{
		"salePublicId":        savedSale.PublicID,
		"previousAmountCents": sale.AmountCents,
		"amountCents":         savedSale.AmountCents,
		"reason":              reason,
	})
	return OperationResult[entity.Renegotiation]{Resource: saved}
}

type CancelSale struct {
	saleRepository        repository.SaleRepository
	invoiceRepository     repository.InvoiceRepository
	installmentRepository repository.InstallmentRepository
	commissionRepository  repository.CommissionRepository
	pendingItemRepository repository.PendingItemRepository
	eventRepository       repository.CommercialEventRepository
	outboxRepository      repository.OutboxEventRepository
}

func NewCancelSale(
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	installmentRepository repository.InstallmentRepository,
	commissionRepository repository.CommissionRepository,
	pendingItemRepository repository.PendingItemRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) CancelSale {
	return CancelSale{
		saleRepository:        saleRepository,
		invoiceRepository:     invoiceRepository,
		installmentRepository: installmentRepository,
		commissionRepository:  commissionRepository,
		pendingItemRepository: pendingItemRepository,
		eventRepository:       eventRepository,
		outboxRepository:      outboxRepository,
	}
}

func (useCase CancelSale) Execute(salePublicID string, reason string) OperationResult[entity.Sale] {
	sale := useCase.saleRepository.FindByPublicID(salePublicID)
	if sale == nil {
		return OperationResult[entity.Sale]{ErrorCode: "sale_not_found", ErrorText: "Sale was not found.", NotFound: true}
	}

	invoice := useCase.invoiceRepository.FindBySalePublicID(sale.PublicID)
	if invoice != nil && invoice.Status == "paid" {
		return OperationResult[entity.Sale]{ErrorCode: "sale_cancellation_blocked", ErrorText: "Paid sales cannot be cancelled directly.", BadRequest: true}
	}

	cancelledSale, err := sale.TransitionTo("cancelled")
	if err != nil {
		return OperationResult[entity.Sale]{ErrorCode: "invalid_sale_status_transition", ErrorText: "Sale status transition is invalid.", BadRequest: true}
	}

	savedSale := useCase.saleRepository.Update(cancelledSale)
	for _, installment := range useCase.installmentRepository.ListBySalePublicID(savedSale.PublicID) {
		if installment.Status == "scheduled" {
			cancelledInstallment, transitionErr := installment.TransitionTo("cancelled")
			if transitionErr == nil {
				useCase.installmentRepository.Update(cancelledInstallment)
			}
		}
	}

	for _, commission := range useCase.commissionRepository.ListBySalePublicID(savedSale.PublicID) {
		if commission.Status == "pending" {
			blockedCommission, transitionErr := commission.TransitionTo("blocked")
			if transitionErr == nil {
				useCase.commissionRepository.Update(blockedCommission)
			}
		}
	}

	for _, item := range useCase.pendingItemRepository.ListBySalePublicID(savedSale.PublicID) {
		if item.Status == "open" {
			useCase.pendingItemRepository.Update(item.Cancel())
		}
	}

	summary := "Sale cancelled."
	if strings.TrimSpace(reason) != "" {
		summary = "Sale cancelled: " + strings.TrimSpace(reason)
	}

	recordCommercialEvent(useCase.eventRepository, "sale", savedSale.PublicID, "sale_cancelled", "sales", summary)
	appendOutboxEvent(useCase.outboxRepository, "sale", savedSale.PublicID, "sale.status_changed", map[string]any{
		"salePublicId": savedSale.PublicID,
		"status":       savedSale.Status,
		"amountCents":  savedSale.AmountCents,
		"reason":       strings.TrimSpace(reason),
	})
	return OperationResult[entity.Sale]{Resource: savedSale}
}
