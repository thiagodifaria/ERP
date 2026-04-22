package entity

import (
	"testing"
	"time"
)

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
		0,
		"",
		"",
		0,
		"",
		nil,
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
	if attachment.Visibility != "internal" {
		t.Fatalf("expected default visibility, got %s", attachment.Visibility)
	}
	if attachment.RetentionDays != 365 {
		t.Fatalf("expected default retention days, got %d", attachment.RetentionDays)
	}
}

func TestAttachmentArchiveShouldPersistReasonAndTimestamp(t *testing.T) {
	attachment, err := NewAttachment(
		"0195e7a0-7a9c-7c1f-8a44-4a6e70000331",
		"bootstrap-ops",
		"crm.lead",
		"0195e7a0-7a9c-7c1f-8a44-4a6e70000341",
		"proposta.pdf",
		"application/pdf",
		"leads/proposta.pdf",
		"manual",
		"crm",
		"ops-user",
		2048,
		"ABCDEF",
		"restricted",
		30,
		"",
		nil,
		Attachment{}.CreatedAt,
	)
	if err != nil {
		t.Fatalf("expected attachment creation to succeed, got %v", err)
	}

	archivedAt := attachment.CreatedAt.Add(time.Hour)
	archived := attachment.Archive("expirado", archivedAt)

	if archived.ArchivedAt == nil || !archived.ArchivedAt.Equal(archivedAt.UTC()) {
		t.Fatalf("expected archived timestamp to be preserved, got %+v", archived.ArchivedAt)
	}
	if archived.ArchiveReason != "expirado" {
		t.Fatalf("expected archive reason to be preserved, got %s", archived.ArchiveReason)
	}
}
