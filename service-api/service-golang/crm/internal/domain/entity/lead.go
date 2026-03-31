// Lead representa a entrada inicial de relacionamento no CRM.
// Validacoes essenciais do agregado devem nascer aqui.
package entity

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrLeadNameRequired            = errors.New("lead name is required")
	ErrLeadEmailInvalid            = errors.New("lead email is invalid")
	ErrLeadStatusInvalid           = errors.New("lead status is invalid")
	ErrLeadStatusTransitionInvalid = errors.New("lead status transition is invalid")
)

type Lead struct {
	PublicID    string
	Name        string
	Email       string
	Source      string
	Status      string
	OwnerUserID string
}

func NewLead(publicID string, name string, email string, source string, ownerUserID string) (Lead, error) {
	normalizedName := strings.TrimSpace(name)
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	normalizedSource := strings.TrimSpace(source)

	if normalizedName == "" {
		return Lead{}, ErrLeadNameRequired
	}

	if _, err := mail.ParseAddress(normalizedEmail); err != nil {
		return Lead{}, ErrLeadEmailInvalid
	}

	if normalizedSource == "" {
		normalizedSource = "manual"
	}

	return Lead{
		PublicID:    publicID,
		Name:        normalizedName,
		Email:       normalizedEmail,
		Source:      normalizedSource,
		Status:      "captured",
		OwnerUserID: strings.TrimSpace(ownerUserID),
	}, nil
}

func (lead Lead) TransitionTo(status string) (Lead, error) {
	targetStatus := normalizeStatus(status)
	if !isValidStatus(targetStatus) {
		return Lead{}, ErrLeadStatusInvalid
	}

	currentStatus := normalizeStatus(lead.Status)
	if currentStatus == targetStatus {
		return lead, nil
	}

	if !canTransition(currentStatus, targetStatus) {
		return Lead{}, ErrLeadStatusTransitionInvalid
	}

	lead.Status = targetStatus
	return lead, nil
}

func (lead Lead) AssignOwner(ownerUserID string) Lead {
	lead.OwnerUserID = strings.TrimSpace(ownerUserID)
	return lead
}

func normalizeStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func isValidStatus(status string) bool {
	switch status {
	case "captured", "contacted", "qualified", "disqualified":
		return true
	default:
		return false
	}
}

func canTransition(currentStatus string, targetStatus string) bool {
	switch currentStatus {
	case "captured":
		return targetStatus == "contacted" || targetStatus == "qualified" || targetStatus == "disqualified"
	case "contacted":
		return targetStatus == "qualified" || targetStatus == "disqualified"
	case "qualified":
		return targetStatus == "disqualified"
	default:
		return false
	}
}
