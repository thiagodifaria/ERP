// Sale representa a venda fechada a partir de uma proposta aceita.
// Regras basicas de faturamento e cancelamento devem nascer aqui.
package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrSalePublicIDInvalid         = errors.New("sale public id is invalid")
	ErrSaleOpportunityIDInvalid    = errors.New("sale opportunity public id is invalid")
	ErrSaleProposalIDInvalid       = errors.New("sale proposal public id is invalid")
	ErrSaleAmountCentsInvalid      = errors.New("sale amount cents is invalid")
	ErrSaleStatusInvalid           = errors.New("sale status is invalid")
	ErrSaleStatusTransitionInvalid = errors.New("sale status transition is invalid")
)

type Sale struct {
	PublicID            string
	OpportunityPublicID string
	ProposalPublicID    string
	Status              string
	AmountCents         int64
}

func NewSale(publicID string, opportunityPublicID string, proposalPublicID string, amountCents int64) (Sale, error) {
	return restoreSale(publicID, opportunityPublicID, proposalPublicID, amountCents, "active")
}

func RestoreSale(publicID string, opportunityPublicID string, proposalPublicID string, amountCents int64, status string) (Sale, error) {
	return restoreSale(publicID, opportunityPublicID, proposalPublicID, amountCents, status)
}

func (sale Sale) TransitionTo(status string) (Sale, error) {
	targetStatus := normalizeSaleStatus(status)
	if !isValidSaleStatus(targetStatus) {
		return Sale{}, ErrSaleStatusInvalid
	}

	currentStatus := normalizeSaleStatus(sale.Status)
	if currentStatus == targetStatus {
		return sale, nil
	}

	if !canTransitionSaleStatus(currentStatus, targetStatus) {
		return Sale{}, ErrSaleStatusTransitionInvalid
	}

	sale.Status = targetStatus
	return sale, nil
}

func restoreSale(publicID string, opportunityPublicID string, proposalPublicID string, amountCents int64, status string) (Sale, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedOpportunityPublicID := strings.TrimSpace(opportunityPublicID)
	normalizedProposalPublicID := strings.TrimSpace(proposalPublicID)
	normalizedStatus := normalizeSaleStatus(status)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Sale{}, ErrSalePublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedOpportunityPublicID); err != nil {
		return Sale{}, ErrSaleOpportunityIDInvalid
	}

	if _, err := uuid.Parse(normalizedProposalPublicID); err != nil {
		return Sale{}, ErrSaleProposalIDInvalid
	}

	if amountCents <= 0 {
		return Sale{}, ErrSaleAmountCentsInvalid
	}

	if !isValidSaleStatus(normalizedStatus) {
		return Sale{}, ErrSaleStatusInvalid
	}

	return Sale{
		PublicID:            normalizedPublicID,
		OpportunityPublicID: normalizedOpportunityPublicID,
		ProposalPublicID:    normalizedProposalPublicID,
		Status:              normalizedStatus,
		AmountCents:         amountCents,
	}, nil
}

func normalizeSaleStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func isValidSaleStatus(status string) bool {
	switch status {
	case "active", "invoiced", "cancelled":
		return true
	default:
		return false
	}
}

func canTransitionSaleStatus(currentStatus string, targetStatus string) bool {
	switch currentStatus {
	case "active":
		return targetStatus == "invoiced" || targetStatus == "cancelled"
	case "invoiced":
		return targetStatus == "cancelled"
	default:
		return false
	}
}
