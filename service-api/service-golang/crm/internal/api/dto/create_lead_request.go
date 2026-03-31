// Requests da API descrevem payloads publicos de entrada.
package dto

type CreateLeadRequest struct {
  Name        string `json:"name"`
  Email       string `json:"email"`
  Source      string `json:"source"`
  OwnerUserID string `json:"ownerUserId"`
}
