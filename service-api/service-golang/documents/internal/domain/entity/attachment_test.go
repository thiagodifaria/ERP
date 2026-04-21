package entity

import "testing"

func TestNewAttachmentShouldNormalizeDefaults(t *testing.T) {
	attachment, err := NewAttachment(
		"0195e7a0-7a9c-7c1f-8a44-4a6e70000331",
		"BOOTSTRAP-OPS",
		"crm.lead",
		"0195e7a0-7a9c-7c1f-8a44-4a6e70000341",
		" proposta.pdf ",
		"",
		" leads/proposta.pdf ",
		"",
		"",
		"",
		Attachment{}.CreatedAt,
	)
	if err != nil {
		t.Fatalf("expected attachment creation to succeed, got %v", err)
	}

	if attachment.TenantSlug != "bootstrap-ops" {
		t.Fatalf("expected normalized tenant slug, got %s", attachment.TenantSlug)
	}
	if attachment.ContentType != "application/octet-stream" {
		t.Fatalf("expected default content type, got %s", attachment.ContentType)
	}
}
