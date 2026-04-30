package dto

import "time"

type UploadSessionResponse struct {
	PublicID           string     `json:"publicId"`
	TenantSlug         string     `json:"tenantSlug"`
	OwnerType          string     `json:"ownerType"`
	OwnerPublicID      string     `json:"ownerPublicId"`
	FileName           string     `json:"fileName"`
	ContentType        string     `json:"contentType"`
	StorageKey         string     `json:"storageKey"`
	StorageDriver      string     `json:"storageDriver"`
	Source             string     `json:"source"`
	RequestedBy        string     `json:"requestedBy"`
	Visibility         string     `json:"visibility"`
	RetentionDays      int        `json:"retentionDays"`
	Status             string     `json:"status"`
	AttachmentPublicID string     `json:"attachmentPublicId,omitempty"`
	UploadURL          string     `json:"uploadUrl"`
	ExpiresAt          time.Time  `json:"expiresAt"`
	CompletedAt        *time.Time `json:"completedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
}
