package entity

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrPipelineConfigPublicIDInvalid = errors.New("pipeline config public id is invalid")
	ErrPipelineConfigNameRequired    = errors.New("pipeline config name is required")
	ErrPipelineConfigStagesRequired  = errors.New("pipeline config stages are required")
)

type PipelineStage struct {
	Key              string `json:"key"`
	Name             string `json:"name"`
	Position         int    `json:"position"`
	RequiresApproval bool   `json:"requiresApproval"`
}

type TerritoryRule struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	AssignmentMode string `json:"assignmentMode"`
}

type ApprovalPolicy struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	ApprovalScope string `json:"approvalScope"`
	RequiredRole  string `json:"requiredRole"`
}

type PipelineConfig struct {
	PublicID         string           `json:"publicId"`
	Name             string           `json:"name"`
	Stages           []PipelineStage  `json:"stages"`
	AutoScoring      bool             `json:"autoScoring"`
	TerritoryRules   []TerritoryRule  `json:"territoryRules"`
	ApprovalPolicies []ApprovalPolicy `json:"approvalPolicies"`
}

func NewPipelineConfig(
	publicID string,
	name string,
	stages []PipelineStage,
	autoScoring bool,
	territoryRules []TerritoryRule,
	approvalPolicies []ApprovalPolicy,
) (PipelineConfig, error) {
	normalizedPublicID := strings.TrimSpace(publicID)
	normalizedName := strings.TrimSpace(name)

	if _, err := uuid.Parse(normalizedPublicID); err != nil {
		return PipelineConfig{}, ErrPipelineConfigPublicIDInvalid
	}

	if normalizedName == "" {
		return PipelineConfig{}, ErrPipelineConfigNameRequired
	}

	if len(stages) == 0 {
		return PipelineConfig{}, ErrPipelineConfigStagesRequired
	}

	normalizedStages := make([]PipelineStage, 0, len(stages))
	for index, stage := range stages {
		key := strings.ToLower(strings.TrimSpace(stage.Key))
		name := strings.TrimSpace(stage.Name)
		if key == "" || name == "" {
			return PipelineConfig{}, ErrPipelineConfigStagesRequired
		}

		normalizedStages = append(normalizedStages, PipelineStage{
			Key:              key,
			Name:             name,
			Position:         index + 1,
			RequiresApproval: stage.RequiresApproval,
		})
	}

	return PipelineConfig{
		PublicID:         normalizedPublicID,
		Name:             normalizedName,
		Stages:           normalizedStages,
		AutoScoring:      autoScoring,
		TerritoryRules:   normalizeTerritoryRules(territoryRules),
		ApprovalPolicies: normalizeApprovalPolicies(approvalPolicies),
	}, nil
}

func DefaultPipelineConfig() PipelineConfig {
	config, _ := NewPipelineConfig(
		"11111111-1111-4111-8111-111111111111",
		"Default Revenue Pipeline",
		[]PipelineStage{
			{Key: "captured", Name: "Captured"},
			{Key: "contacted", Name: "Contacted"},
			{Key: "qualified", Name: "Qualified", RequiresApproval: true},
			{Key: "won", Name: "Won"},
		},
		true,
		[]TerritoryRule{
			{Key: "default-brazil", Name: "Default Brazil", AssignmentMode: "owner_or_region"},
		},
		[]ApprovalPolicy{
			{Key: "discount-approval", Name: "Discount Approval", ApprovalScope: "discount", RequiredRole: "manager"},
		},
	)
	return config
}

func normalizeTerritoryRules(rules []TerritoryRule) []TerritoryRule {
	if len(rules) == 0 {
		return nil
	}

	normalized := make([]TerritoryRule, 0, len(rules))
	for _, rule := range rules {
		key := strings.ToLower(strings.TrimSpace(rule.Key))
		name := strings.TrimSpace(rule.Name)
		mode := strings.ToLower(strings.TrimSpace(rule.AssignmentMode))
		if key == "" || name == "" || mode == "" {
			continue
		}
		normalized = append(normalized, TerritoryRule{
			Key:            key,
			Name:           name,
			AssignmentMode: mode,
		})
	}

	return normalized
}

func normalizeApprovalPolicies(policies []ApprovalPolicy) []ApprovalPolicy {
	if len(policies) == 0 {
		return nil
	}

	normalized := make([]ApprovalPolicy, 0, len(policies))
	for _, policy := range policies {
		key := strings.ToLower(strings.TrimSpace(policy.Key))
		name := strings.TrimSpace(policy.Name)
		scope := strings.ToLower(strings.TrimSpace(policy.ApprovalScope))
		requiredRole := strings.ToLower(strings.TrimSpace(policy.RequiredRole))
		if key == "" || name == "" || scope == "" || requiredRole == "" {
			continue
		}
		normalized = append(normalized, ApprovalPolicy{
			Key:           key,
			Name:          name,
			ApprovalScope: scope,
			RequiredRole:  requiredRole,
		})
	}

	return normalized
}
