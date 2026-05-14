package dto

import "time"

type DocumentAuditEventResponse struct {
	PublicID           string    `json:"publicId"`
	TenantSlug         string    `json:"tenantSlug"`
	AttachmentPublicID string    `json:"attachmentPublicId"`
	EventCode          string    `json:"eventCode"`
	Actor              string    `json:"actor"`
	Reason             string    `json:"reason"`
	CorrelationID      string    `json:"correlationId"`
	CreatedAt          time.Time `json:"createdAt"`
}
