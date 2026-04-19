package entity

import (
	"testing"
	"time"
)

func TestNewInvoiceShouldDefaultToDraft(t *testing.T) {
	invoice, err := NewInvoice(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000061",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000062",
		"ACME-INV-0001",
		120000,
		"2026-05-20",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if invoice.Status != "draft" {
		t.Fatalf("expected status draft, got %s", invoice.Status)
	}
}

func TestInvoiceShouldStampPaidAtWhenMarkedPaid(t *testing.T) {
	invoice, _ := RestoreInvoice(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000071",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000072",
		"ACME-INV-0002",
		120000,
		"2026-05-20",
		"sent",
		"",
	)

	paid, err := invoice.TransitionTo("paid", time.Date(2026, 5, 21, 14, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected transition error: %v", err)
	}

	if paid.PaidAt == "" {
		t.Fatalf("expected paidAt to be filled")
	}

	if _, err := paid.TransitionTo("sent", time.Now().UTC()); err == nil {
		t.Fatalf("expected reverse transition to fail")
	}
}

func TestInvoiceShouldDetectOverdueOnlyWhenStillOpen(t *testing.T) {
	invoice, _ := RestoreInvoice(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000081",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000082",
		"ACME-INV-0003",
		120000,
		"2026-05-01",
		"sent",
		"",
	)

	if !invoice.IsOverdue(time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected invoice to be overdue")
	}

	paid, _ := invoice.TransitionTo("paid", time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC))
	if paid.IsOverdue(time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected paid invoice to stop being overdue")
	}
}
