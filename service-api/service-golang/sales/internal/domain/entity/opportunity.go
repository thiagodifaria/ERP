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
	ErrOpportunityCustomerPublicIDInvalid = errors.New("opportunity customer public id is invalid")
	ErrOpportunityTitleRequired          = errors.New("opportunity title is required")
	ErrOpportunityOwnerUserIDInvalid     = errors.New("opportunity owner user id is invalid")
	ErrOpportunityAmountCentsInvalid     = errors.New("opportunity amount cents is invalid")
	ErrOpportunitySaleTypeInvalid        = errors.New("opportunity sale type is invalid")
	ErrOpportunityStageInvalid           = errors.New("opportunity stage is invalid")
	ErrOpportunityStageTransitionInvalid = errors.New("opportunity stage transition is invalid")
)

type Opportunity struct {
	PublicID         string
	LeadPublicID     string
	CustomerPublicID string
	Title            string
	Stage            string
	SaleType         string
	OwnerUserID      string
	AmountCents      int64
}

func NewOpportunity(publicID string, leadPublicID string, customerPublicID string, title string, saleType string, ownerUserID string, amountCents int64) (Opportunity, error) {
	return restoreOpportunity(publicID, leadPublicID, customerPublicID, title, saleType, ownerUserID, amountCents, "qualified")
}

func RestoreOpportunity(publicID string, leadPublicID string, customerPublicID string, title string, saleType string, ownerUserID string, amountCents int64, stage string) (Opportunity, error) {
	return restoreOpportunity(publicID, leadPublicID, customerPublicID, title, saleType, ownerUserID, amountCents, stage)
}

func (opportunity Opportunity) Revise(customerPublicID string, title string, saleType string, ownerUserID string, amountCents int64) (Opportunity, error) {
	return restoreOpportunity(opportunity.PublicID, opportunity.LeadPublicID, customerPublicID, title, saleType, ownerUserID, amountCents, opportunity.Stage)
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

func restoreOpportunity(publicID string, leadPublicID string, customerPublicID string, title string, saleType string, ownerUserID string, amountCents int64, stage string) (Opportunity, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedLeadPublicID := strings.TrimSpace(leadPublicID)
	normalizedCustomerPublicID := strings.TrimSpace(customerPublicID)
	normalizedTitle := strings.TrimSpace(title)
	normalizedSaleType := normalizeSaleType(saleType)
	normalizedOwnerUserID := strings.TrimSpace(ownerUserID)
	normalizedStage := normalizeOpportunityStage(stage)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Opportunity{}, ErrOpportunityPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedLeadPublicID); err != nil {
		return Opportunity{}, ErrOpportunityLeadPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedCustomerPublicID); err != nil {
		return Opportunity{}, ErrOpportunityCustomerPublicIDInvalid
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

	if !isValidSaleType(normalizedSaleType) {
		return Opportunity{}, ErrOpportunitySaleTypeInvalid
	}

	if !isValidOpportunityStage(normalizedStage) {
		return Opportunity{}, ErrOpportunityStageInvalid
	}

	return Opportunity{
		PublicID:         normalizedPublicID,
		LeadPublicID:     normalizedLeadPublicID,
		CustomerPublicID: normalizedCustomerPublicID,
		Title:            normalizedTitle,
		Stage:            normalizedStage,
		SaleType:         normalizedSaleType,
		OwnerUserID:      normalizedOwnerUserID,
		AmountCents:      amountCents,
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

func normalizeSaleType(saleType string) string {
	return strings.ToLower(strings.TrimSpace(saleType))
}

func isValidSaleType(saleType string) bool {
	switch saleType {
	case "new", "upsell", "renewal", "cross_sell":
		return true
	default:
		return false
	}
}
