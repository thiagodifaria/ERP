package dto

type BulkLeadImportRequest struct {
	TenantSlug string              `json:"tenantSlug"`
	Items      []CreateLeadRequest `json:"items"`
}
