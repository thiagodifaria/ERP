package entity

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrCustomerPublicIDInvalid     = errors.New("customer public id is invalid")
	ErrCustomerLeadPublicIDInvalid = errors.New("customer lead public id is invalid")
	ErrCustomerNameRequired        = errors.New("customer name is required")
	ErrCustomerEmailInvalid        = errors.New("customer email is invalid")
	ErrCustomerOwnerUserIDInvalid  = errors.New("customer owner user id is invalid")
	ErrCustomerStatusInvalid       = errors.New("customer status is invalid")
	ErrCustomerSourceRequired      = errors.New("customer source is required")
)

type Customer struct {
	PublicID     string
	LeadPublicID string
	Name         string
	Email        string
	Source       string
	Status       string
	OwnerUserID  string
}

func NewCustomerFromLead(publicID string, lead Lead) (Customer, error) {
	return restoreCustomer(publicID, lead.PublicID, lead.Name, lead.Email, lead.Source, "active", lead.OwnerUserID)
}

func RestoreCustomer(publicID string, leadPublicID string, name string, email string, source string, status string, ownerUserID string) (Customer, error) {
	return restoreCustomer(publicID, leadPublicID, name, email, source, status, ownerUserID)
}

func restoreCustomer(publicID string, leadPublicID string, name string, email string, source string, status string, ownerUserID string) (Customer, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedLeadPublicID := strings.TrimSpace(leadPublicID)
	normalizedName := strings.TrimSpace(name)
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	normalizedSource := strings.TrimSpace(source)
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	normalizedOwnerUserID := strings.TrimSpace(ownerUserID)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return Customer{}, ErrCustomerPublicIDInvalid
	}

	if _, err := uuid.Parse(normalizedLeadPublicID); err != nil {
		return Customer{}, ErrCustomerLeadPublicIDInvalid
	}

	if normalizedName == "" {
		return Customer{}, ErrCustomerNameRequired
	}

	if _, err := mail.ParseAddress(normalizedEmail); err != nil {
		return Customer{}, ErrCustomerEmailInvalid
	}

	if normalizedSource == "" {
		return Customer{}, ErrCustomerSourceRequired
	}

	if normalizedOwnerUserID != "" {
		if _, err := uuid.Parse(normalizedOwnerUserID); err != nil {
			return Customer{}, ErrCustomerOwnerUserIDInvalid
		}
	}

	switch normalizedStatus {
	case "active", "inactive":
	default:
		return Customer{}, ErrCustomerStatusInvalid
	}

	return Customer{
		PublicID:     normalizedPublicID,
		LeadPublicID: normalizedLeadPublicID,
		Name:         normalizedName,
		Email:        normalizedEmail,
		Source:       normalizedSource,
		Status:       normalizedStatus,
		OwnerUserID:  normalizedOwnerUserID,
	}, nil
}
