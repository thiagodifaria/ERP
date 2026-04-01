// Opportunity representa a oportunidade comercial aberta a partir de um lead.
// Transicoes de stage e validacoes essenciais devem nascer aqui.
package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrOpportunityPublicIDInvalid        = errors.New("opportunity public id is invalid")
	ErrOpportunityLeadPublicIDInvalid    = errors.New("opportunity lead public id is invalid")
	ErrOpportunityTitleRequired          = errors.New("opportunity title is required")
	ErrOpportunityOwnerUserIDInvalid     = errors.New("opportunity owner user id is invalid")
	ErrOpportunityAmountCentsInvalid     = errors.New("opportunity amount cents is invalid")
	ErrOpportunityStageInvalid           = errors.New("opportunity stage is invalid")
	ErrOpportunityStageTransitionInvalid = errors.New("opportunity stage transition is invalid")
)

type Opportunity struct {
	PublicID     string
	LeadPublicID string
	Title        string
	Stage        string
	OwnerUserID  string
	AmountCents  int64
}

func NewOpportunity(publicID string, leadPublicID string, title string, ownerUserID string, amountCents int64) (Opportunity, error) {
	return restoreOpportunity(publicID, leadPublicID, title, ownerUserID, amountCents, "qualified")
}

func RestoreOpportunity(publicID string, leadPublicID string, title string, ownerUserID string, amountCents int64, stage string) (Opportunity, error) {
	return restoreOpportunity(publicID, leadPublicID, title, ownerUserID, amountCents, stage)
}

func (opportunity Opportunity) Revise(title string, ownerUserID string, amountCents int64) (Opportunity, error) {
	return restoreOpportunity(opportunity.PublicID, opportunity.LeadPublicID, title, ownerUserID, amountCents, opportunity.Stage)
}

func (opportunity Opportunity) TransitionTo(stage string) (Opportunity, error) {
	targetStage := normalizeOpportunityStage(stage)
	if !isValidOpportunityStage(targetStage) {
		return Opportunity{}, ErrOpportunityStageInvalid
	}

	currentStage := normalizeOpportunityStage(opportunity.Stage)
	if currentStage == targetStage {
		return opportunity, nil
	}

	if !canTransitionOpportunityStage(currentStage, targetStage) {
		return Opportunity{}, ErrOpportunityStageTransitionInvalid
	}

	opportunity.Stage = targetStage
	return opportunity, nil
}

func restoreOpportunity(publicID string, leadPublicID string, title string, ownerUserID string, amountCents int64, stage string) (Opportunity, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedLeadPublicID := strings.TrimSpace(leadPublicID)
	normalizedTitle := strings.TrimSpace(title)
	normalizedOwnerUserID := strings.TrimSpace(ownerUserID)
	normalizedStage := normalizeOpportunityStage(stage)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Opportunity{}, ErrOpportunityPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedLeadPublicID); err != nil {
		return Opportunity{}, ErrOpportunityLeadPublicIDInvalid
	}

	if normalizedTitle == "" {
		return Opportunity{}, ErrOpportunityTitleRequired
	}

	if normalizedOwnerUserID != "" {
		if _, err := uuid.Parse(normalizedOwnerUserID); err != nil {
			return Opportunity{}, ErrOpportunityOwnerUserIDInvalid
		}
	}

	if amountCents <= 0 {
		return Opportunity{}, ErrOpportunityAmountCentsInvalid
	}

	if !isValidOpportunityStage(normalizedStage) {
		return Opportunity{}, ErrOpportunityStageInvalid
	}

	return Opportunity{
		PublicID:     normalizedPublicID,
		LeadPublicID: normalizedLeadPublicID,
		Title:        normalizedTitle,
		Stage:        normalizedStage,
		OwnerUserID:  normalizedOwnerUserID,
		AmountCents:  amountCents,
	}, nil
}

func normalizeOpportunityStage(stage string) string {
	return strings.ToLower(strings.TrimSpace(stage))
}

func isValidOpportunityStage(stage string) bool {
	switch stage {
	case "qualified", "proposal", "negotiation", "won", "lost":
		return true
	default:
		return false
	}
}

func canTransitionOpportunityStage(currentStage string, targetStage string) bool {
	switch currentStage {
	case "qualified":
		return targetStage == "proposal" || targetStage == "negotiation" || targetStage == "lost"
	case "proposal":
		return targetStage == "negotiation" || targetStage == "won" || targetStage == "lost"
	case "negotiation":
		return targetStage == "won" || targetStage == "lost"
	default:
		return false
	}
}
