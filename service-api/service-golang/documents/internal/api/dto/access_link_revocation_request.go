package dto

type AccessLinkRevocationRequest struct {
	TenantSlug  string `json:"tenantSlug"`
	AccessToken string `json:"accessToken"`
	Reason      string `json:"reason"`
	RevokedBy   string `json:"revokedBy"`
}
