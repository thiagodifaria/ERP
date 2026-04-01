// Proposal representa a proposta comercial associada a uma oportunidade.
// Regras de aceite e rejeicao devem nascer aqui.
package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrProposalPublicIDInvalid         = errors.New("proposal public id is invalid")
	ErrProposalOpportunityIDInvalid    = errors.New("proposal opportunity public id is invalid")
	ErrProposalTitleRequired           = errors.New("proposal title is required")
	ErrProposalAmountCentsInvalid      = errors.New("proposal amount cents is invalid")
	ErrProposalStatusInvalid           = errors.New("proposal status is invalid")
	ErrProposalStatusTransitionInvalid = errors.New("proposal status transition is invalid")
)

type Proposal struct {
	PublicID            string
	OpportunityPublicID string
	Title               string
	Status              string
	AmountCents         int64
}

func NewProposal(publicID string, opportunityPublicID string, title string, amountCents int64) (Proposal, error) {
	return restoreProposal(publicID, opportunityPublicID, title, amountCents, "draft")
}

func RestoreProposal(publicID string, opportunityPublicID string, title string, amountCents int64, status string) (Proposal, error) {
	return restoreProposal(publicID, opportunityPublicID, title, amountCents, status)
}

func (proposal Proposal) TransitionTo(status string) (Proposal, error) {
	targetStatus := normalizeProposalStatus(status)
	if !isValidProposalStatus(targetStatus) {
		return Proposal{}, ErrProposalStatusInvalid
	}

	currentStatus := normalizeProposalStatus(proposal.Status)
	if currentStatus == targetStatus {
		return proposal, nil
	}

	if !canTransitionProposalStatus(currentStatus, targetStatus) {
		return Proposal{}, ErrProposalStatusTransitionInvalid
	}

	proposal.Status = targetStatus
	return proposal, nil
}

func restoreProposal(publicID string, opportunityPublicID string, title string, amountCents int64, status string) (Proposal, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedOpportunityPublicID := strings.TrimSpace(opportunityPublicID)
	normalizedTitle := strings.TrimSpace(title)
	normalizedStatus := normalizeProposalStatus(status)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Proposal{}, ErrProposalPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedOpportunityPublicID); err != nil {
		return Proposal{}, ErrProposalOpportunityIDInvalid
	}

	if normalizedTitle == "" {
		return Proposal{}, ErrProposalTitleRequired
	}

	if amountCents <= 0 {
		return Proposal{}, ErrProposalAmountCentsInvalid
	}

	if !isValidProposalStatus(normalizedStatus) {
		return Proposal{}, ErrProposalStatusInvalid
	}

	return Proposal{
		PublicID:            normalizedPublicID,
		OpportunityPublicID: normalizedOpportunityPublicID,
		Title:               normalizedTitle,
		Status:              normalizedStatus,
		AmountCents:         amountCents,
	}, nil
}

func normalizeProposalStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func isValidProposalStatus(status string) bool {
	switch status {
	case "draft", "sent", "accepted", "rejected":
		return true
	default:
		return false
	}
}

func canTransitionProposalStatus(currentStatus string, targetStatus string) bool {
	switch currentStatus {
	case "draft":
		return targetStatus == "sent" || targetStatus == "rejected"
	case "sent":
		return targetStatus == "accepted" || targetStatus == "rejected"
	default:
		return false
	}
}
