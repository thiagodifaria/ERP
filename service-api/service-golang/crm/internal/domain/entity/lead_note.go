// LeadNote registra contexto operacional ao redor de um lead sem alterar o agregado principal.
package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrLeadNotePublicIDInvalid     = errors.New("lead note public id is invalid")
	ErrLeadNoteLeadPublicIDInvalid = errors.New("lead note lead public id is invalid")
	ErrLeadNoteBodyRequired        = errors.New("lead note body is required")
)

type LeadNote struct {
	PublicID     string
	LeadPublicID string
	Body         string
	Category     string
	CreatedAt    time.Time
}

func NewLeadNote(
	publicID string,
	leadPublicID string,
	body string,
	category string,
	createdAt time.Time,
) (LeadNote, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedLeadPublicID := strings.TrimSpace(leadPublicID)
	normalizedBody := strings.TrimSpace(body)
	normalizedCategory := strings.ToLower(strings.TrimSpace(category))

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return LeadNote{}, ErrLeadNotePublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedLeadPublicID); err != nil {
		return LeadNote{}, ErrLeadNoteLeadPublicIDInvalid
	}

	if normalizedBody == "" {
		return LeadNote{}, ErrLeadNoteBodyRequired
	}

	if normalizedCategory == "" {
		normalizedCategory = "internal"
	}

	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	return LeadNote{
		PublicID:     normalizedPublicID,
		LeadPublicID: normalizedLeadPublicID,
		Body:         normalizedBody,
		Category:     normalizedCategory,
		CreatedAt:    createdAt.UTC(),
	}, nil
}
