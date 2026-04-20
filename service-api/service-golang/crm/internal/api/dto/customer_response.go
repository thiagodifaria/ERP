package dto

type CustomerResponse struct {
	PublicID     string `json:"publicId"`
	LeadPublicID string `json:"leadPublicId"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Source       string `json:"source"`
	Status       string `json:"status"`
	OwnerUserID  string `json:"ownerUserId"`
}
