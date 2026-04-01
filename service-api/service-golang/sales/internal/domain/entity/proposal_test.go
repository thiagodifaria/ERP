package entity

import "testing"

func TestNewProposalShouldDefaultToDraft(t *testing.T) {
	proposal, err := NewProposal(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000021",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000022",
		"Annual Proposal",
		120000,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if proposal.Status != "draft" {
		t.Fatalf("expected status draft, got %s", proposal.Status)
	}
}

func TestProposalShouldRestrictTransitions(t *testing.T) {
	proposal, _ := RestoreProposal(
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000031",
		"0195e7a0-7a9c-7c1f-8a44-4a6e72000032",
		"Annual Proposal",
		120000,
		"sent",
	)

	accepted, err := proposal.TransitionTo("accepted")
	if err != nil {
		t.Fatalf("unexpected transition error: %v", err)
	}

	if accepted.Status != "accepted" {
		t.Fatalf("expected status accepted, got %s", accepted.Status)
	}

	if _, err := accepted.TransitionTo("draft"); err == nil {
		t.Fatalf("expected invalid transition to fail")
	}
}
