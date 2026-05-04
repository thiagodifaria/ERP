package dto

import "time"

type SigningRequestResponse struct {
	PublicID                 string    `json:"publicId"`
	TenantSlug               string    `json:"tenantSlug"`
	AttachmentPublicID       string    `json:"attachmentPublicId"`
	Provider                 string    `json:"provider"`
	DocumentKind             string    `json:"documentKind"`
	Status                   string    `json:"status"`
	RequestedBy              string    `json:"requestedBy"`
	Signers                  []string  `json:"signers"`
	RelatedAggregate         string    `json:"relatedAggregate"`
	RelatedAggregatePublicID string    `json:"relatedAggregatePublicId"`
	SigningURL               string    `json:"signingUrl"`
	CreatedAt                time.Time `json:"createdAt"`
}
