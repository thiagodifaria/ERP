package entity

import "testing"

func TestNewCustomerFromLeadShouldProjectQualifiedLead(t *testing.T) {
	lead, _ := NewLead("0195e7a0-7a9c-7c1f-8a44-4a6e70000111", "CRM Customer", "crm.customer@example.com", "manual", "0195e7a0-7a9c-7c1f-8a44-4a6e70000112")
	qualified, _ := lead.TransitionTo("qualified")

	customer, err := NewCustomerFromLead("0195e7a0-7a9c-7c1f-8a44-4a6e70000113", qualified)
	if err != nil {
		t.Fatalf("expected customer conversion to succeed, got %v", err)
	}

	if customer.Status != "active" {
		t.Fatalf("expected active customer, got %s", customer.Status)
	}

	if customer.LeadPublicID != qualified.PublicID {
		t.Fatalf("expected lead linkage to be preserved")
	}
}
