package dto

type ContractsExecutiveSummary struct {
	Status              string `json:"status"`
	HTTPSpecs           int    `json:"httpSpecs"`
	EventSchemas        int    `json:"eventSchemas"`
	ADRCount            int    `json:"adrCount"`
	ContractArtifacts   int    `json:"contractArtifacts"`
	NavigableAPIReady   bool   `json:"navigableApiReady"`
	SchemaRegistryReady bool   `json:"schemaRegistryReady"`
}

type ContractsOverviewResponse struct {
	Service              string                    `json:"service"`
	TenantSlug           string                    `json:"tenantSlug"`
	GeneratedAt          string                    `json:"generatedAt"`
	ExecutiveSummary     ContractsExecutiveSummary `json:"executiveSummary"`
	IntegrationReadiness map[string]any            `json:"integrationReadiness"`
	ContractGovernance   map[string]any            `json:"contractGovernance"`
	HardeningReview      map[string]any            `json:"hardeningReview"`
}
