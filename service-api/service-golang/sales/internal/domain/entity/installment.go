package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInstallmentPublicIDInvalid      = errors.New("installment public id is invalid")
	ErrInstallmentSalePublicIDInvalid  = errors.New("installment sale public id is invalid")
	ErrInstallmentSequenceInvalid      = errors.New("installment sequence is invalid")
	ErrInstallmentAmountCentsInvalid   = errors.New("installment amount cents is invalid")
	ErrInstallmentDueDateInvalid       = errors.New("installment due date is invalid")
	ErrInstallmentStatusInvalid        = errors.New("installment status is invalid")
	ErrInstallmentStatusTransitionInvalid = errors.New("installment status transition is invalid")
)

type Installment struct {
	PublicID       string
	SalePublicID   string
	SequenceNumber int
	AmountCents    int64
	DueDate        string
	Status         string
}

func NewInstallment(publicID string, salePublicID string, sequenceNumber int, amountCents int64, dueDate string) (Installment, error) {
	return restoreInstallment(publicID, salePublicID, sequenceNumber, amountCents, dueDate, "scheduled")
}

func RestoreInstallment(publicID string, salePublicID string, sequenceNumber int, amountCents int64, dueDate string, status string) (Installment, error) {
	return restoreInstallment(publicID, salePublicID, sequenceNumber, amountCents, dueDate, status)
}

func (installment Installment) TransitionTo(status string) (Installment, error) {
	targetStatus := strings.ToLower(strings.TrimSpace(status))
	if !isValidInstallmentStatus(targetStatus) {
		return Installment{}, ErrInstallmentStatusInvalid
	}

	if installment.Status == targetStatus {
		return installment, nil
	}

	if !canTransitionInstallmentStatus(installment.Status, targetStatus) {
		return Installment{}, ErrInstallmentStatusTransitionInvalid
	}

	installment.Status = targetStatus
	return installment, nil
}

func restoreInstallment(publicID string, salePublicID string, sequenceNumber int, amountCents int64, dueDate string, status string) (Installment, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedSalePublicID := strings.TrimSpace(salePublicID)
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	normalizedDueDate := strings.TrimSpace(dueDate)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Installment{}, ErrInstallmentPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedSalePublicID); err != nil {
		return Installment{}, ErrInstallmentSalePublicIDInvalid
	}

	if sequenceNumber <= 0 {
		return Installment{}, ErrInstallmentSequenceInvalid
	}

	if amountCents <= 0 {
		return Installment{}, ErrInstallmentAmountCentsInvalid
	}

	if _, err := time.Parse("2006-01-02", normalizedDueDate); err != nil {
		return Installment{}, ErrInstallmentDueDateInvalid
	}

	if !isValidInstallmentStatus(normalizedStatus) {
		return Installment{}, ErrInstallmentStatusInvalid
	}

	return Installment{
		PublicID:       normalizedPublicID,
		SalePublicID:   normalizedSalePublicID,
		SequenceNumber: sequenceNumber,
		AmountCents:    amountCents,
		DueDate:        normalizedDueDate,
		Status:         normalizedStatus,
	}, nil
}

func isValidInstallmentStatus(status string) bool {
	switch status {
	case "scheduled", "paid", "cancelled":
		return true
	default:
		return false
	}
}

func canTransitionInstallmentStatus(currentStatus string, targetStatus string) bool {
	switch currentStatus {
	case "scheduled":
		return targetStatus == "paid" || targetStatus == "cancelled"
	default:
		return false
	}
}
