package dto

type ArchiveAttachmentRequest struct {
	TenantSlug string `json:"tenantSlug"`
	Reason     string `json:"reason"`
}
