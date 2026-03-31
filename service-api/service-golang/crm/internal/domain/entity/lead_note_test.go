package entity

import (
	"testing"
	"time"
)

const (
	testLeadNotePublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e70000031"
)

func TestNewLeadNoteShouldNormalizeAndCreateNote(t *testing.T) {
	createdAt := time.Date(2026, time.March, 31, 13, 59, 0, 0, time.UTC)

	note, err := NewLeadNote(
		testLeadNotePublicID,
		testLeadPublicID,
		"  Cliente pediu retorno apos almoco.  ",
		"  Follow-Up  ",
		createdAt,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if note.Body != "Cliente pediu retorno apos almoco." {
		t.Fatalf("expected normalized body, got %s", note.Body)
	}

	if note.Category != "follow-up" {
		t.Fatalf("expected normalized category follow-up, got %s", note.Category)
	}

	if !note.CreatedAt.Equal(createdAt) {
		t.Fatalf("expected createdAt to be preserved")
	}
}

func TestNewLeadNoteShouldApplyDefaultCategory(t *testing.T) {
	note, err := NewLeadNote(
		testLeadNotePublicID,
		testLeadPublicID,
		"Primeira observacao interna",
		"   ",
		time.Time{},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if note.Category != "internal" {
		t.Fatalf("expected default category internal, got %s", note.Category)
	}
}

func TestNewLeadNoteShouldRejectInvalidPublicID(t *testing.T) {
	_, err := NewLeadNote(
		"note-public-id",
		testLeadPublicID,
		"Primeira observacao interna",
		"internal",
		time.Now().UTC(),
	)

	if err != ErrLeadNotePublicIDInvalid {
		t.Fatalf("expected ErrLeadNotePublicIDInvalid, got %v", err)
	}
}

func TestNewLeadNoteShouldRejectInvalidLeadPublicID(t *testing.T) {
	_, err := NewLeadNote(
		testLeadNotePublicID,
		"lead-public-id",
		"Primeira observacao interna",
		"internal",
		time.Now().UTC(),
	)

	if err != ErrLeadNoteLeadPublicIDInvalid {
		t.Fatalf("expected ErrLeadNoteLeadPublicIDInvalid, got %v", err)
	}
}

func TestNewLeadNoteShouldRejectBlankBody(t *testing.T) {
	_, err := NewLeadNote(
		testLeadNotePublicID,
		testLeadPublicID,
		"   ",
		"internal",
		time.Now().UTC(),
	)

	if err != ErrLeadNoteBodyRequired {
		t.Fatalf("expected ErrLeadNoteBodyRequired, got %v", err)
	}
}
