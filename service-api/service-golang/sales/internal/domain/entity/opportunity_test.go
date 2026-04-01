package entity

import "testing"

func TestNewOpportunityShouldNormalizeAndDefaultStage(t *testing.T) {
	opportunity, err := NewOpportunity(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000001",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000002",
		"  Enterprise Expansion  ",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000003",
		99000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opportunity.Stage != "qualified" {
		t.Fatalf("expected stage qualified, got %s", opportunity.Stage)
	}

	if opportunity.Title != "Enterprise Expansion" {
		t.Fatalf("expected trimmed title, got %s", opportunity.Title)
	}
}

func TestOpportunityShouldRespectTransitions(t *testing.T) {
	opportunity, _ := RestoreOpportunity(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000011",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000012",
		"Expansion",
		"",
		120000,
		"proposal",
	)

	updated, err := opportunity.TransitionTo("won")
	if err != nil {
		t.Fatalf("unexpected transition error: %v", err)
	}

	if updated.Stage != "won" {
		t.Fatalf("expected stage won, got %s", updated.Stage)
	}

	if _, err := updated.TransitionTo("proposal"); err == nil {
		t.Fatalf("expected terminal stage transition to fail")
	}
}
