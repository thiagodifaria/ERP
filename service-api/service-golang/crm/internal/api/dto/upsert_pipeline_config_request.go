package dto

type PipelineStageRequest struct {
	Key              string `json:"key"`
	Name             string `json:"name"`
	RequiresApproval bool   `json:"requiresApproval"`
}

type TerritoryRuleRequest struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	AssignmentMode string `json:"assignmentMode"`
}

type ApprovalPolicyRequest struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	ApprovalScope string `json:"approvalScope"`
	RequiredRole  string `json:"requiredRole"`
}

type UpsertPipelineConfigRequest struct {
	TenantSlug       string                  `json:"tenantSlug"`
	Name             string                  `json:"name"`
	AutoScoring      bool                    `json:"autoScoring"`
	Stages           []PipelineStageRequest  `json:"stages"`
	TerritoryRules   []TerritoryRuleRequest  `json:"territoryRules"`
	ApprovalPolicies []ApprovalPolicyRequest `json:"approvalPolicies"`
}
