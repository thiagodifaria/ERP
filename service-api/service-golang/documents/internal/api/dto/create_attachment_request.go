package dto

type CreateAttachmentRequest struct {
	TenantSlug     string `json:"tenantSlug"`
	OwnerType      string `json:"ownerType"`
	OwnerPublicID  string `json:"ownerPublicId"`
	FileName       string `json:"fileName"`
	ContentType    string `json:"contentType"`
	StorageKey     string `json:"storageKey"`
	StorageDriver  string `json:"storageDriver"`
	Source         string `json:"source"`
	UploadedBy     string `json:"uploadedBy"`
	FileSizeBytes  int64  `json:"fileSizeBytes"`
	ChecksumSHA256 string `json:"checksumSha256"`
	Visibility     string `json:"visibility"`
	RetentionDays  int    `json:"retentionDays"`
}
