package dto

type CreateSigningRequest struct {
	TenantSlug         string   `json:"tenantSlug"`
	AttachmentPublicID string   `json:"attachmentPublicId"`
	DocumentKind       string   `json:"documentKind"`
	RequestedBy        string   `json:"requestedBy"`
	Signers            []string `json:"signers"`
	Provider           string   `json:"provider"`
	RelatedAggregate   string   `json:"relatedAggregate"`
	RelatedAggregateID string   `json:"relatedAggregatePublicId"`
}
