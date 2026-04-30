package dto

type CreateUploadSessionRequest struct {
	TenantSlug       string `json:"tenantSlug"`
	OwnerType        string `json:"ownerType"`
	OwnerPublicID    string `json:"ownerPublicId"`
	FileName         string `json:"fileName"`
	ContentType      string `json:"contentType"`
	StorageKey       string `json:"storageKey"`
	StorageDriver    string `json:"storageDriver"`
	Source           string `json:"source"`
	RequestedBy      string `json:"requestedBy"`
	Visibility       string `json:"visibility"`
	RetentionDays    int    `json:"retentionDays"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
}
