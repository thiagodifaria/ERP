package entity

import "testing"

func TestNewLeadShouldNormalizeAndCreateLead(t *testing.T) {
  lead, err := NewLead(
    "lead-public-id",
    "  Ana Souza  ",
    "ANA@EXAMPLE.COM",
    "",
    "owner-public-id",
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
  _, err := NewLead("lead-public-id", "   ", "ana@example.com", "meta-ads", "owner-public-id")

  if err != ErrLeadNameRequired {
    t.Fatalf("expected ErrLeadNameRequired, got %v", err)
  }
}

func TestNewLeadShouldRejectInvalidEmail(t *testing.T) {
  _, err := NewLead("lead-public-id", "Ana Souza", "invalid-email", "meta-ads", "owner-public-id")

  if err != ErrLeadEmailInvalid {
    t.Fatalf("expected ErrLeadEmailInvalid, got %v", err)
  }
}

func TestTransitionToShouldMoveLeadToQualified(t *testing.T) {
  lead, err := NewLead("lead-public-id", "Ana Souza", "ana@example.com", "meta-ads", "owner-public-id")
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
  lead, err := NewLead("lead-public-id", "Ana Souza", "ana@example.com", "meta-ads", "owner-public-id")
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
