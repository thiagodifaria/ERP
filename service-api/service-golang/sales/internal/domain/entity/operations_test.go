package entity

import "testing"

func TestNewCommissionShouldCalculateAmountFromRateBps(t *testing.T) {
	commission, err := NewCommission(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000111",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000112",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000113",
		"closer",
		750,
		200000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if commission.AmountCents != 15000 {
		t.Fatalf("expected 15000 cents, got %d", commission.AmountCents)
	}
}

func TestInstallmentShouldRejectInvalidSequence(t *testing.T) {
	if _, err := NewInstallment(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000121",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000122",
		0,
		50000,
		"2026-06-10",
	); err == nil {
		t.Fatalf("expected invalid sequence error")
	}
}
