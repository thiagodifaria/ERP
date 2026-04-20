// Lead representa a entrada inicial de relacionamento no CRM.
// Validacoes essenciais do agregado devem nascer aqui.
package entity

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrLeadPublicIDInvalid         = errors.New("lead public id is invalid")
	ErrLeadNameRequired            = errors.New("lead name is required")
	ErrLeadEmailInvalid            = errors.New("lead email is invalid")
	ErrLeadOwnerUserIDInvalid      = errors.New("lead owner user id is invalid")
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
	return restoreLead(publicID, name, email, source, "captured", ownerUserID)
}

func RestoreLead(publicID string, name string, email string, source string, status string, ownerUserID string) (Lead, error) {
	return restoreLead(publicID, name, email, source, status, ownerUserID)
}

func restoreLead(publicID string, name string, email string, source string, status string, ownerUserID string) (Lead, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedName := strings.TrimSpace(name)
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	normalizedSource := strings.TrimSpace(source)
	normalizedStatus := normalizeStatus(status)
	normalizedOwnerUserID := strings.TrimSpace(ownerUserID)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Lead{}, ErrLeadPublicIDInvalid
	}

	if normalizedName == "" {
		return Lead{}, ErrLeadNameRequired
	}

	if _, err := mail.ParseAddress(normalizedEmail); err != nil {
		return Lead{}, ErrLeadEmailInvalid
	}

	if normalizedSource == "" {
		normalizedSource = "manual"
	}

	if normalizedOwnerUserID != "" {
		if _, err := uuid.Parse(normalizedOwnerUserID); err != nil {
			return Lead{}, ErrLeadOwnerUserIDInvalid
		}
	}

	if !isValidStatus(normalizedStatus) {
		return Lead{}, ErrLeadStatusInvalid
	}

	return Lead{
		PublicID:    normalizedPublicID,
		Name:        normalizedName,
		Email:       normalizedEmail,
		Source:      normalizedSource,
		Status:      normalizedStatus,
		OwnerUserID: normalizedOwnerUserID,
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

func (lead Lead) ReviseProfile(name string, email string, source string) (Lead, error) {
	revisedLead, err := restoreLead(lead.PublicID, name, email, source, lead.Status, lead.OwnerUserID)
	if err != nil {
		return Lead{}, err
	}
	return revisedLead, nil
}

func (lead Lead) AssignOwner(ownerUserID string) (Lead, error) {
	normalizedOwnerUserID := strings.TrimSpace(ownerUserID)
	if normalizedOwnerUserID != "" {
		if _, err := uuid.Parse(normalizedOwnerUserID); err != nil {
			return Lead{}, ErrLeadOwnerUserIDInvalid
		}
	}

	lead.OwnerUserID = normalizedOwnerUserID
	return lead, nil
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
