// Requests da API descrevem payloads publicos de entrada.
package dto

type UpdateLeadStatusRequest struct {
  Status string `json:"status"`
}
