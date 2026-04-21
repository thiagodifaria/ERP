package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrCommissionPublicIDInvalid      = errors.New("commission public id is invalid")
	ErrCommissionSalePublicIDInvalid  = errors.New("commission sale public id is invalid")
	ErrCommissionRecipientInvalid     = errors.New("commission recipient is invalid")
	ErrCommissionRoleCodeRequired     = errors.New("commission role code is required")
	ErrCommissionRateBpsInvalid       = errors.New("commission rate bps is invalid")
	ErrCommissionAmountCentsInvalid   = errors.New("commission amount cents is invalid")
	ErrCommissionStatusInvalid        = errors.New("commission status is invalid")
	ErrCommissionStatusTransitionInvalid = errors.New("commission status transition is invalid")
)

type Commission struct {
	PublicID        string
	SalePublicID    string
	RecipientUserID string
	RoleCode        string
	RateBps         int
	AmountCents     int64
	Status          string
}

func NewCommission(publicID string, salePublicID string, recipientUserID string, roleCode string, rateBps int, saleAmountCents int64) (Commission, error) {
	amountCents := (saleAmountCents * int64(rateBps)) / 10000
	return restoreCommission(publicID, salePublicID, recipientUserID, roleCode, rateBps, amountCents, "pending")
}

func RestoreCommission(publicID string, salePublicID string, recipientUserID string, roleCode string, rateBps int, amountCents int64, status string) (Commission, error) {
	return restoreCommission(publicID, salePublicID, recipientUserID, roleCode, rateBps, amountCents, status)
}

func (commission Commission) TransitionTo(status string) (Commission, error) {
	targetStatus := strings.ToLower(strings.TrimSpace(status))
	if !isValidCommissionStatus(targetStatus) {
		return Commission{}, ErrCommissionStatusInvalid
	}

	if commission.Status == targetStatus {
		return commission, nil
	}

	if !canTransitionCommissionStatus(commission.Status, targetStatus) {
		return Commission{}, ErrCommissionStatusTransitionInvalid
	}

	commission.Status = targetStatus
	return commission, nil
}

func restoreCommission(publicID string, salePublicID string, recipientUserID string, roleCode string, rateBps int, amountCents int64, status string) (Commission, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedSalePublicID := strings.TrimSpace(salePublicID)
	normalizedRecipientUserID := strings.TrimSpace(recipientUserID)
	normalizedRoleCode := strings.ToLower(strings.TrimSpace(roleCode))
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Commission{}, ErrCommissionPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedSalePublicID); err != nil {
		return Commission{}, ErrCommissionSalePublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedRecipientUserID); err != nil {
		return Commission{}, ErrCommissionRecipientInvalid
	}

	if normalizedRoleCode == "" {
		return Commission{}, ErrCommissionRoleCodeRequired
	}

	if rateBps <= 0 || rateBps > 10000 {
		return Commission{}, ErrCommissionRateBpsInvalid
	}

	if amountCents <= 0 {
		return Commission{}, ErrCommissionAmountCentsInvalid
	}

	if !isValidCommissionStatus(normalizedStatus) {
		return Commission{}, ErrCommissionStatusInvalid
	}

	return Commission{
		PublicID:        normalizedPublicID,
		SalePublicID:    normalizedSalePublicID,
		RecipientUserID: normalizedRecipientUserID,
		RoleCode:        normalizedRoleCode,
		RateBps:         rateBps,
		AmountCents:     amountCents,
		Status:          normalizedStatus,
	}, nil
}

func isValidCommissionStatus(status string) bool {
	switch status {
	case "pending", "blocked", "released":
		return true
	default:
		return false
	}
}

func canTransitionCommissionStatus(currentStatus string, targetStatus string) bool {
	switch currentStatus {
	case "pending":
		return targetStatus == "blocked" || targetStatus == "released"
	case "blocked":
		return targetStatus == "released"
	default:
		return false
	}
}
