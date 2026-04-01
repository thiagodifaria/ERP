package entity

import "testing"

func TestNewSaleShouldDefaultToActive(t *testing.T) {
	sale, err := NewSale(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000041",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000042",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000043",
		120000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sale.Status != "active" {
		t.Fatalf("expected status active, got %s", sale.Status)
	}
}

func TestSaleShouldAllowInvoicingOnlyOnce(t *testing.T) {
	sale, _ := RestoreSale(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000051",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000052",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000053",
		120000,
		"active",
	)

	invoiced, err := sale.TransitionTo("invoiced")
	if err != nil {
		t.Fatalf("unexpected transition error: %v", err)
	}

	if _, err := invoiced.TransitionTo("active"); err == nil {
		t.Fatalf("expected reverse transition to fail")
	}
}
