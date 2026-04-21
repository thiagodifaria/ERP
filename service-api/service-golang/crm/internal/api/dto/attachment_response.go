package dto

import "time"

type AttachmentResponse struct {
	PublicID      string    `json:"publicId"`
	TenantSlug    string    `json:"tenantSlug"`
	OwnerType     string    `json:"ownerType"`
	OwnerPublicID string    `json:"ownerPublicId"`
	FileName      string    `json:"fileName"`
	ContentType   string    `json:"contentType"`
	StorageKey    string    `json:"storageKey"`
	StorageDriver string    `json:"storageDriver"`
	Source        string    `json:"source"`
	UploadedBy    string    `json:"uploadedBy"`
	CreatedAt     time.Time `json:"createdAt"`
}
