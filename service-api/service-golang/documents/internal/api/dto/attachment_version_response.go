package dto

import "time"

type AttachmentVersionResponse struct {
	PublicID           string    `json:"publicId"`
	TenantSlug         string    `json:"tenantSlug"`
	AttachmentPublicID string    `json:"attachmentPublicId"`
	VersionNumber      int       `json:"versionNumber"`
	FileName           string    `json:"fileName"`
	ContentType        string    `json:"contentType"`
	StorageKey         string    `json:"storageKey"`
	StorageDriver      string    `json:"storageDriver"`
	Source             string    `json:"source"`
	UploadedBy         string    `json:"uploadedBy"`
	FileSizeBytes      int64     `json:"fileSizeBytes"`
	ChecksumSHA256     string    `json:"checksumSha256"`
	CreatedAt          time.Time `json:"createdAt"`
}
