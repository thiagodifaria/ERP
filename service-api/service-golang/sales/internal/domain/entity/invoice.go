// Invoice representa o documento basico de cobranca originado de uma venda.
// O contexto sales segura esse ciclo inicial ate a extracao de billing amadurecer.
package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvoicePublicIDInvalid         = errors.New("invoice public id is invalid")
	ErrInvoiceSaleIDInvalid           = errors.New("invoice sale public id is invalid")
	ErrInvoiceNumberInvalid           = errors.New("invoice number is invalid")
	ErrInvoiceAmountCentsInvalid      = errors.New("invoice amount cents is invalid")
	ErrInvoiceDueDateInvalid          = errors.New("invoice due date is invalid")
	ErrInvoiceStatusInvalid           = errors.New("invoice status is invalid")
	ErrInvoiceStatusTransitionInvalid = errors.New("invoice status transition is invalid")
	ErrInvoicePaidAtInvalid           = errors.New("invoice paid at is invalid")
)

type Invoice struct {
	PublicID     string
	SalePublicID string
	Number       string
	Status       string
	AmountCents  int64
	DueDate      string
	PaidAt       string
}

func NewInvoice(publicID string, salePublicID string, number string, amountCents int64, dueDate string) (Invoice, error) {
	return restoreInvoice(publicID, salePublicID, number, amountCents, dueDate, "draft", "")
}

func RestoreInvoice(
	publicID string,
	salePublicID string,
	number string,
	amountCents int64,
	dueDate string,
	status string,
	paidAt string,
) (Invoice, error) {
	return restoreInvoice(publicID, salePublicID, number, amountCents, dueDate, status, paidAt)
}

func (invoice Invoice) TransitionTo(status string, changedAt time.Time) (Invoice, error) {
	targetStatus := normalizeInvoiceStatus(status)
	if !isValidInvoiceStatus(targetStatus) {
		return Invoice{}, ErrInvoiceStatusInvalid
	}

	currentStatus := normalizeInvoiceStatus(invoice.Status)
	if currentStatus == targetStatus {
		return invoice, nil
	}

	if !canTransitionInvoiceStatus(currentStatus, targetStatus) {
		return Invoice{}, ErrInvoiceStatusTransitionInvalid
	}

	invoice.Status = targetStatus
	if targetStatus == "paid" {
		invoice.PaidAt = changedAt.UTC().Format(time.RFC3339)
	}

	return invoice, nil
}

func (invoice Invoice) IsOverdue(reference time.Time) bool {
	if invoice.Status == "paid" || invoice.Status == "cancelled" {
		return false
	}

	dueDate, err := time.Parse("2006-01-02", invoice.DueDate)
	if err != nil {
		return false
	}

	today := time.Date(reference.UTC().Year(), reference.UTC().Month(), reference.UTC().Day(), 0, 0, 0, 0, time.UTC)
	return dueDate.Before(today)
}

func restoreInvoice(
	publicID string,
	salePublicID string,
	number string,
	amountCents int64,
	dueDate string,
	status string,
	paidAt string,
) (Invoice, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedSalePublicID := strings.TrimSpace(salePublicID)
	normalizedNumber := strings.ToUpper(strings.TrimSpace(number))
	normalizedDueDate := strings.TrimSpace(dueDate)
	normalizedStatus := normalizeInvoiceStatus(status)
	normalizedPaidAt := strings.TrimSpace(paidAt)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Invoice{}, ErrInvoicePublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedSalePublicID); err != nil {
		return Invoice{}, ErrInvoiceSaleIDInvalid
	}

	if normalizedNumber == "" {
		return Invoice{}, ErrInvoiceNumberInvalid
	}

	if amountCents <= 0 {
		return Invoice{}, ErrInvoiceAmountCentsInvalid
	}

	if _, err := time.Parse("2006-01-02", normalizedDueDate); err != nil {
		return Invoice{}, ErrInvoiceDueDateInvalid
	}

	if !isValidInvoiceStatus(normalizedStatus) {
		return Invoice{}, ErrInvoiceStatusInvalid
	}

	if normalizedPaidAt != "" {
		if _, err := time.Parse(time.RFC3339, normalizedPaidAt); err != nil {
			return Invoice{}, ErrInvoicePaidAtInvalid
		}
	}

	return Invoice{
		PublicID:     normalizedPublicID,
		SalePublicID: normalizedSalePublicID,
		Number:       normalizedNumber,
		Status:       normalizedStatus,
		AmountCents:  amountCents,
		DueDate:      normalizedDueDate,
		PaidAt:       normalizedPaidAt,
	}, nil
}

func normalizeInvoiceStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func isValidInvoiceStatus(status string) bool {
	switch status {
	case "draft", "sent", "paid", "cancelled":
		return true
	default:
		return false
	}
}

func canTransitionInvoiceStatus(currentStatus string, targetStatus string) bool {
	switch currentStatus {
	case "draft":
		return targetStatus == "sent" || targetStatus == "paid" || targetStatus == "cancelled"
	case "sent":
		return targetStatus == "paid" || targetStatus == "cancelled"
	default:
		return false
	}
}
