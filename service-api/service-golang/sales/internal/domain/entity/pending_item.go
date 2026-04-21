package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPendingItemPublicIDInvalid   = errors.New("pending item public id is invalid")
	ErrPendingItemSalePublicIDInvalid = errors.New("pending item sale public id is invalid")
	ErrPendingItemCodeRequired      = errors.New("pending item code is required")
	ErrPendingItemSummaryRequired   = errors.New("pending item summary is required")
	ErrPendingItemStatusInvalid     = errors.New("pending item status is invalid")
	ErrPendingItemStatusTransitionInvalid = errors.New("pending item status transition is invalid")
)

type PendingItem struct {
	PublicID   string
	SalePublicID string
	Code       string
	Summary    string
	Status     string
	ResolvedAt string
}

func NewPendingItem(publicID string, salePublicID string, code string, summary string) (PendingItem, error) {
	return restorePendingItem(publicID, salePublicID, code, summary, "open", "")
}

func RestorePendingItem(publicID string, salePublicID string, code string, summary string, status string, resolvedAt string) (PendingItem, error) {
	return restorePendingItem(publicID, salePublicID, code, summary, status, resolvedAt)
}

func (item PendingItem) Resolve(resolvedAt time.Time) (PendingItem, error) {
	if item.Status != "open" {
		return PendingItem{}, ErrPendingItemStatusTransitionInvalid
	}

	item.Status = "resolved"
	item.ResolvedAt = resolvedAt.UTC().Format(time.RFC3339)
	return item, nil
}

func (item PendingItem) Cancel() PendingItem {
	item.Status = "cancelled"
	item.ResolvedAt = ""
	return item
}

func restorePendingItem(publicID string, salePublicID string, code string, summary string, status string, resolvedAt string) (PendingItem, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedSalePublicID := strings.TrimSpace(salePublicID)
	normalizedCode := strings.ToLower(strings.TrimSpace(code))
	normalizedSummary := strings.TrimSpace(summary)
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	normalizedResolvedAt := strings.TrimSpace(resolvedAt)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return PendingItem{}, ErrPendingItemPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedSalePublicID); err != nil {
		return PendingItem{}, ErrPendingItemSalePublicIDInvalid
	}

	if normalizedCode == "" {
		return PendingItem{}, ErrPendingItemCodeRequired
	}

	if normalizedSummary == "" {
		return PendingItem{}, ErrPendingItemSummaryRequired
	}

	switch normalizedStatus {
	case "open", "resolved", "cancelled":
	default:
		return PendingItem{}, ErrPendingItemStatusInvalid
	}

	if normalizedResolvedAt != "" {
		if _, err := time.Parse(time.RFC3339, normalizedResolvedAt); err != nil {
			return PendingItem{}, ErrPendingItemStatusInvalid
		}
	}

	return PendingItem{
		PublicID:     normalizedPublicID,
		SalePublicID: normalizedSalePublicID,
		Code:         normalizedCode,
		Summary:      normalizedSummary,
		Status:       normalizedStatus,
		ResolvedAt:   normalizedResolvedAt,
	}, nil
}
