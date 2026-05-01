package dto

type CreateAttachmentVersionRequest struct {
	TenantSlug     string `json:"tenantSlug"`
	FileName       string `json:"fileName"`
	ContentType    string `json:"contentType"`
	StorageKey     string `json:"storageKey"`
	StorageDriver  string `json:"storageDriver"`
	Source         string `json:"source"`
	UploadedBy     string `json:"uploadedBy"`
	FileSizeBytes  int64  `json:"fileSizeBytes"`
	ChecksumSHA256 string `json:"checksumSha256"`
}
