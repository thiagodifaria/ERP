package entity

import "testing"

const (
	testLeadPublicID      = "0195e7a0-7a9c-7c1f-8a44-4a6e70000001"
	testOwnerUserPublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e70000011"
)

func TestNewLeadShouldNormalizeAndCreateLead(t *testing.T) {
	lead, err := NewLead(
		testLeadPublicID,
		"  Ana Souza  ",
		"ANA@EXAMPLE.COM",
		"",
		testOwnerUserPublicID,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lead.Name != "Ana Souza" {
		t.Fatalf("expected normalized name, got %s", lead.Name)
	}

	if lead.Email != "ana@example.com" {
		t.Fatalf("expected normalized email, got %s", lead.Email)
	}

	if lead.Source != "manual" {
		t.Fatalf("expected default source manual, got %s", lead.Source)
	}

	if lead.Status != "captured" {
		t.Fatalf("expected status captured, got %s", lead.Status)
	}
}

func TestNewLeadShouldRejectBlankName(t *testing.T) {
	_, err := NewLead(testLeadPublicID, "   ", "ana@example.com", "meta-ads", testOwnerUserPublicID)

	if err != ErrLeadNameRequired {
		t.Fatalf("expected ErrLeadNameRequired, got %v", err)
	}
}

func TestNewLeadShouldRejectInvalidEmail(t *testing.T) {
	_, err := NewLead(testLeadPublicID, "Ana Souza", "invalid-email", "meta-ads", testOwnerUserPublicID)

	if err != ErrLeadEmailInvalid {
		t.Fatalf("expected ErrLeadEmailInvalid, got %v", err)
	}
}

func TestNewLeadShouldRejectInvalidPublicID(t *testing.T) {
	_, err := NewLead("lead-public-id", "Ana Souza", "ana@example.com", "meta-ads", testOwnerUserPublicID)

	if err != ErrLeadPublicIDInvalid {
		t.Fatalf("expected ErrLeadPublicIDInvalid, got %v", err)
	}
}

func TestNewLeadShouldRejectInvalidOwnerUserID(t *testing.T) {
	_, err := NewLead(testLeadPublicID, "Ana Souza", "ana@example.com", "meta-ads", "owner-public-id")

	if err != ErrLeadOwnerUserIDInvalid {
		t.Fatalf("expected ErrLeadOwnerUserIDInvalid, got %v", err)
	}
}

func TestTransitionToShouldMoveLeadToQualified(t *testing.T) {
	lead, err := NewLead(testLeadPublicID, "Ana Souza", "ana@example.com", "meta-ads", testOwnerUserPublicID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	transitionedLead, err := lead.TransitionTo("qualified")
	if err != nil {
		t.Fatalf("unexpected transition error: %v", err)
	}

	if transitionedLead.Status != "qualified" {
		t.Fatalf("expected status qualified, got %s", transitionedLead.Status)
	}
}

func TestTransitionToShouldRejectInvalidTransition(t *testing.T) {
	lead, err := NewLead(testLeadPublicID, "Ana Souza", "ana@example.com", "meta-ads", testOwnerUserPublicID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contactedLead, err := lead.TransitionTo("contacted")
	if err != nil {
		t.Fatalf("unexpected transition error: %v", err)
	}

	_, err = contactedLead.TransitionTo("captured")
	if err != ErrLeadStatusTransitionInvalid {
		t.Fatalf("expected ErrLeadStatusTransitionInvalid, got %v", err)
	}
}

func TestAssignOwnerShouldNormalizeValue(t *testing.T) {
	lead, err := NewLead(testLeadPublicID, "Ana Souza", "ana@example.com", "meta-ads", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assignedLead, err := lead.AssignOwner("  " + testOwnerUserPublicID + "  ")
	if err != nil {
		t.Fatalf("unexpected owner assignment error: %v", err)
	}

	if assignedLead.OwnerUserID != testOwnerUserPublicID {
		t.Fatalf("expected normalized owner %s, got %s", testOwnerUserPublicID, assignedLead.OwnerUserID)
	}

	unassignedLead, err := assignedLead.AssignOwner("   ")
	if err != nil {
		t.Fatalf("unexpected owner clearing error: %v", err)
	}

	if unassignedLead.OwnerUserID != "" {
		t.Fatalf("expected owner to be cleared, got %s", unassignedLead.OwnerUserID)
	}
}

func TestAssignOwnerShouldRejectInvalidValue(t *testing.T) {
	lead, err := NewLead(testLeadPublicID, "Ana Souza", "ana@example.com", "meta-ads", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = lead.AssignOwner("owner-ana")
	if err != ErrLeadOwnerUserIDInvalid {
		t.Fatalf("expected ErrLeadOwnerUserIDInvalid, got %v", err)
	}
}
