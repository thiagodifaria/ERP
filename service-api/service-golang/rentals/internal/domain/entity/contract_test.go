package entity

import (
	"testing"
	"time"
)

func TestContractShouldGenerateRecurringCharges(t *testing.T) {
	contract, err := NewContract(
		"0196445e-8d3a-78ae-8c74-bc7c4a400001",
		"bootstrap-ops",
		"0196445e-8d3a-78ae-8c74-bc7c4a400002",
		"Contrato Torre Sul",
		"TS-101",
		"BRL",
		150000,
		5,
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 7, 31, 0, 0, 0, 0, time.UTC),
		"active",
		nil,
		"",
		time.Time{},
		time.Time{},
	)
	if err != nil {
		t.Fatalf("unexpected contract error: %v", err)
	}

	charges := contract.BuildCharges()
	if len(charges) != 3 {
		t.Fatalf("expected 3 charges, got %d", len(charges))
	}
	if charges[0].DueDate.Format("2006-01-02") != "2026-05-05" || charges[2].DueDate.Format("2006-01-02") != "2026-07-05" {
		t.Fatalf("unexpected due dates: %+v", charges)
	}
}

func TestContractShouldApplyAdjustmentAndTermination(t *testing.T) {
	contract, err := NewContract(
		"0196445e-8d3a-78ae-8c74-bc7c4a400011",
		"bootstrap-ops",
		"0196445e-8d3a-78ae-8c74-bc7c4a400012",
		"Contrato Torre Norte",
		"TN-201",
		"BRL",
		120000,
		10,
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
		"active",
		nil,
		"",
		time.Time{},
		time.Time{},
	)
	if err != nil {
		t.Fatalf("unexpected contract error: %v", err)
	}

	charges := contract.BuildCharges()
	updatedContract, adjustment, updatedCharges, _, _, err := contract.ApplyAdjustment(
		time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		135000,
		"Reajuste anual",
		"ops-user",
		charges,
	)
	if err != nil {
		t.Fatalf("unexpected adjustment error: %v", err)
	}
	if adjustment.NewAmountCents != 135000 || updatedContract.AmountCents != 135000 {
		t.Fatalf("unexpected adjustment outcome: %+v %+v", adjustment, updatedContract)
	}
	if updatedCharges[0].AmountCents != 120000 || updatedCharges[1].AmountCents != 135000 || updatedCharges[2].AmountCents != 135000 {
		t.Fatalf("unexpected adjusted charges: %+v", updatedCharges)
	}

	terminatedContract, terminatedCharges, _, _, err := updatedContract.Terminate(
		time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC),
		"Rescisao amigavel",
		"ops-user",
		updatedCharges,
	)
	if err != nil {
		t.Fatalf("unexpected termination error: %v", err)
	}
	if terminatedContract.Status != "terminated" {
		t.Fatalf("expected terminated contract, got %+v", terminatedContract)
	}
	if terminatedCharges[2].Status != "cancelled" {
		t.Fatalf("expected future charge cancellation, got %+v", terminatedCharges)
	}
}
