package dto

type CompleteUploadSessionRequest struct {
	TenantSlug     string `json:"tenantSlug"`
	UploadedBy     string `json:"uploadedBy"`
	FileSizeBytes  int64  `json:"fileSizeBytes"`
	ChecksumSHA256 string `json:"checksumSha256"`
}
