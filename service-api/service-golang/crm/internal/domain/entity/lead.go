// Lead representa a entrada inicial de relacionamento no CRM.
// Validacoes essenciais do agregado devem nascer aqui.
package entity

import (
  "errors"
  "net/mail"
  "strings"
)

var (
  ErrLeadNameRequired = errors.New("lead name is required")
  ErrLeadEmailInvalid = errors.New("lead email is invalid")
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
