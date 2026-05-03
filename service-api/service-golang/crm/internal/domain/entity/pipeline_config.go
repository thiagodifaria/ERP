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

type PipelineConfig struct {
	PublicID    string          `json:"publicId"`
	Name        string          `json:"name"`
	Stages      []PipelineStage `json:"stages"`
	AutoScoring bool            `json:"autoScoring"`
}

func NewPipelineConfig(publicID string, name string, stages []PipelineStage, autoScoring bool) (PipelineConfig, error) {
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
		PublicID:    normalizedPublicID,
		Name:        normalizedName,
		Stages:      normalizedStages,
		AutoScoring: autoScoring,
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
	)
	return config
}
