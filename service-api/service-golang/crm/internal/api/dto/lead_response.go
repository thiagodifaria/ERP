// LeadResponse descreve a saida publica minima do agregado de lead.
package dto

type LeadResponse struct {
	PublicID    string `json:"publicId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Source      string `json:"source"`
	Status      string `json:"status"`
	OwnerUserID string `json:"ownerUserId"`
}
