package dto

import "time"

type AccessLinkResponse struct {
	AttachmentPublicID string    `json:"attachmentPublicId"`
	TenantSlug         string    `json:"tenantSlug"`
	StorageDriver      string    `json:"storageDriver"`
	StorageKey         string    `json:"storageKey"`
	AccessURL          string    `json:"accessUrl"`
	ExpiresAt          time.Time `json:"expiresAt"`
	AccessMode         string    `json:"accessMode"`
}
