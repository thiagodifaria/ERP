package dto

type AccessLinkRequest struct {
	TenantSlug       string `json:"tenantSlug"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
	RequestedBy      string `json:"requestedBy"`
}
