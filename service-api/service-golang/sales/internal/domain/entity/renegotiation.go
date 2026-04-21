package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRenegotiationPublicIDInvalid = errors.New("renegotiation public id is invalid")
	ErrRenegotiationSalePublicIDInvalid = errors.New("renegotiation sale public id is invalid")
	ErrRenegotiationReasonRequired = errors.New("renegotiation reason is required")
	ErrRenegotiationAmountInvalid = errors.New("renegotiation amount is invalid")
)

type Renegotiation struct {
	PublicID            string
	SalePublicID        string
	Reason              string
	PreviousAmountCents int64
	NewAmountCents      int64
	Status              string
	AppliedAt           string
}

func NewAppliedRenegotiation(publicID string, salePublicID string, reason string, previousAmountCents int64, newAmountCents int64, appliedAt time.Time) (Renegotiation, error) {
	return restoreRenegotiation(publicID, salePublicID, reason, previousAmountCents, newAmountCents, "applied", appliedAt.UTC().Format(time.RFC3339))
}

func RestoreRenegotiation(publicID string, salePublicID string, reason string, previousAmountCents int64, newAmountCents int64, status string, appliedAt string) (Renegotiation, error) {
	return restoreRenegotiation(publicID, salePublicID, reason, previousAmountCents, newAmountCents, status, appliedAt)
}

func restoreRenegotiation(publicID string, salePublicID string, reason string, previousAmountCents int64, newAmountCents int64, status string, appliedAt string) (Renegotiation, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedSalePublicID := strings.TrimSpace(salePublicID)
	normalizedReason := strings.TrimSpace(reason)
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	normalizedAppliedAt := strings.TrimSpace(appliedAt)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Renegotiation{}, ErrRenegotiationPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedSalePublicID); err != nil {
		return Renegotiation{}, ErrRenegotiationSalePublicIDInvalid
	}

	if normalizedReason == "" {
		return Renegotiation{}, ErrRenegotiationReasonRequired
	}

	if previousAmountCents <= 0 || newAmountCents <= 0 || previousAmountCents == newAmountCents {
		return Renegotiation{}, ErrRenegotiationAmountInvalid
	}

	if normalizedStatus == "" {
		normalizedStatus = "applied"
	}

	if normalizedAppliedAt == "" {
		return Renegotiation{}, ErrRenegotiationAmountInvalid
	}

	if _, err := time.Parse(time.RFC3339, normalizedAppliedAt); err != nil {
		return Renegotiation{}, ErrRenegotiationAmountInvalid
	}

	return Renegotiation{
		PublicID:            normalizedPublicID,
		SalePublicID:        normalizedSalePublicID,
		Reason:              normalizedReason,
		PreviousAmountCents: previousAmountCents,
		NewAmountCents:      newAmountCents,
		Status:              normalizedStatus,
		AppliedAt:           normalizedAppliedAt,
	}, nil
}
