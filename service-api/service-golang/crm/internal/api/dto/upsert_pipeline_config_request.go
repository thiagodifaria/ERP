package dto

type PipelineStageRequest struct {
	Key              string `json:"key"`
	Name             string `json:"name"`
	RequiresApproval bool   `json:"requiresApproval"`
}

type UpsertPipelineConfigRequest struct {
	TenantSlug  string                 `json:"tenantSlug"`
	Name        string                 `json:"name"`
	AutoScoring bool                   `json:"autoScoring"`
	Stages      []PipelineStageRequest `json:"stages"`
}
