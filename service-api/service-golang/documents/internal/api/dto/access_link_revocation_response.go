package dto

import "time"

type AccessLinkRevocationResponse struct {
	AttachmentPublicID string    `json:"attachmentPublicId"`
	TenantSlug         string    `json:"tenantSlug"`
	Revoked            bool      `json:"revoked"`
	Reason             string    `json:"reason"`
	RevokedBy          string    `json:"revokedBy"`
	RevokedAt          time.Time `json:"revokedAt"`
}
