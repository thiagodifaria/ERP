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
	ErrSaleCustomerPublicIDInvalid = errors.New("sale customer public id is invalid")
	ErrSaleOwnerUserIDInvalid      = errors.New("sale owner user id is invalid")
	ErrSaleAmountCentsInvalid      = errors.New("sale amount cents is invalid")
	ErrSaleTypeInvalid             = errors.New("sale type is invalid")
	ErrSaleStatusInvalid           = errors.New("sale status is invalid")
	ErrSaleStatusTransitionInvalid = errors.New("sale status transition is invalid")
)

type Sale struct {
	PublicID            string
	OpportunityPublicID string
	ProposalPublicID    string
	CustomerPublicID    string
	OwnerUserID         string
	SaleType            string
	Status              string
	AmountCents         int64
}

func NewSale(publicID string, opportunityPublicID string, proposalPublicID string, customerPublicID string, ownerUserID string, saleType string, amountCents int64) (Sale, error) {
	return restoreSale(publicID, opportunityPublicID, proposalPublicID, customerPublicID, ownerUserID, saleType, amountCents, "active")
}

func RestoreSale(publicID string, opportunityPublicID string, proposalPublicID string, customerPublicID string, ownerUserID string, saleType string, amountCents int64, status string) (Sale, error) {
	return restoreSale(publicID, opportunityPublicID, proposalPublicID, customerPublicID, ownerUserID, saleType, amountCents, status)
}

func (sale Sale) ReviseAmount(amountCents int64) (Sale, error) {
	return restoreSale(sale.PublicID, sale.OpportunityPublicID, sale.ProposalPublicID, sale.CustomerPublicID, sale.OwnerUserID, sale.SaleType, amountCents, sale.Status)
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

func restoreSale(publicID string, opportunityPublicID string, proposalPublicID string, customerPublicID string, ownerUserID string, saleType string, amountCents int64, status string) (Sale, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedOpportunityPublicID := strings.TrimSpace(opportunityPublicID)
	normalizedProposalPublicID := strings.TrimSpace(proposalPublicID)
	normalizedCustomerPublicID := strings.TrimSpace(customerPublicID)
	normalizedOwnerUserID := strings.TrimSpace(ownerUserID)
	normalizedSaleType := normalizeSaleType(saleType)
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

	if _, err := uuid.Parse(normalizedCustomerPublicID); err != nil {
		return Sale{}, ErrSaleCustomerPublicIDInvalid
	}

	if normalizedOwnerUserID != "" {
		if _, err := uuid.Parse(normalizedOwnerUserID); err != nil {
			return Sale{}, ErrSaleOwnerUserIDInvalid
		}
	}

	if amountCents <= 0 {
		return Sale{}, ErrSaleAmountCentsInvalid
	}

	if !isValidSaleType(normalizedSaleType) {
		return Sale{}, ErrSaleTypeInvalid
	}

	if !isValidSaleStatus(normalizedStatus) {
		return Sale{}, ErrSaleStatusInvalid
	}

	return Sale{
		PublicID:            normalizedPublicID,
		OpportunityPublicID: normalizedOpportunityPublicID,
		ProposalPublicID:    normalizedProposalPublicID,
		CustomerPublicID:    normalizedCustomerPublicID,
		OwnerUserID:         normalizedOwnerUserID,
		SaleType:            normalizedSaleType,
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
