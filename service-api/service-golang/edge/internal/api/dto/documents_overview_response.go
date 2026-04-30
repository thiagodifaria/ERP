// DTO do cockpit documental entregue pelo edge.
// O gateway agrega inventario, uploads e governanca sem replicar calculos do analytics.
package dto

type DocumentsExecutiveSummary struct {
	Status                string `json:"status"`
	AttachmentsTotal      int    `json:"attachmentsTotal"`
	ActiveAttachments     int    `json:"activeAttachments"`
	ArchivedAttachments   int    `json:"archivedAttachments"`
	RestrictedAttachments int    `json:"restrictedAttachments"`
	PendingUploads        int    `json:"pendingUploads"`
	CompletedUploads      int    `json:"completedUploads"`
	ExternalStorageAssets int    `json:"externalStorageAssets"`
	LongTermRetention     int    `json:"longTermRetention"`
}

type DocumentsOverviewResponse struct {
	Service             string                    `json:"service"`
	TenantSlug          string                    `json:"tenantSlug"`
	GeneratedAt         string                    `json:"generatedAt"`
	ExecutiveSummary    DocumentsExecutiveSummary `json:"executiveSummary"`
	ServicePulse        map[string]any            `json:"servicePulse"`
	Tenant360           map[string]any            `json:"tenant360"`
	DocumentGovernance  map[string]any            `json:"documentGovernance"`
}
