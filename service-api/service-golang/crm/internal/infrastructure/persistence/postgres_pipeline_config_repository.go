package persistence

import (
	"database/sql"
	"encoding/json"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type PostgresPipelineConfigRepository struct {
	database *sql.DB
	tenantID int64
}

func NewPostgresPipelineConfigRepository(database *sql.DB, tenantSlug string) (repository.PipelineConfigRepository, error) {
	tenantID, err := lookupCrmTenantID(database, tenantSlug)
	if err != nil {
		return nil, err
	}

	return &PostgresPipelineConfigRepository{
		database: database,
		tenantID: tenantID,
	}, nil
}

func (repository *PostgresPipelineConfigRepository) Get() *entity.PipelineConfig {
	const statement = `
		SELECT public_id, config_name, stages_json, auto_scoring, territory_rules_json, approval_policies_json
		FROM crm.pipeline_configs
		WHERE tenant_id = $1
		LIMIT 1
	`

	row := repository.database.QueryRow(statement, repository.tenantID)

	var publicID string
	var name string
	var stagesRaw []byte
	var autoScoring bool
	var territoryRulesRaw []byte
	var approvalPoliciesRaw []byte
	if err := row.Scan(&publicID, &name, &stagesRaw, &autoScoring, &territoryRulesRaw, &approvalPoliciesRaw); err != nil {
		if err == sql.ErrNoRows {
			defaultConfig := entity.DefaultPipelineConfig()
			return &defaultConfig
		}
		defaultConfig := entity.DefaultPipelineConfig()
		return &defaultConfig
	}

	var stages []entity.PipelineStage
	if err := json.Unmarshal(stagesRaw, &stages); err != nil {
		defaultConfig := entity.DefaultPipelineConfig()
		return &defaultConfig
	}

	var territoryRules []entity.TerritoryRule
	if len(territoryRulesRaw) > 0 {
		if err := json.Unmarshal(territoryRulesRaw, &territoryRules); err != nil {
			defaultConfig := entity.DefaultPipelineConfig()
			return &defaultConfig
		}
	}

	var approvalPolicies []entity.ApprovalPolicy
	if len(approvalPoliciesRaw) > 0 {
		if err := json.Unmarshal(approvalPoliciesRaw, &approvalPolicies); err != nil {
			defaultConfig := entity.DefaultPipelineConfig()
			return &defaultConfig
		}
	}

	config, err := entity.NewPipelineConfig(publicID, name, stages, autoScoring, territoryRules, approvalPolicies)
	if err != nil {
		defaultConfig := entity.DefaultPipelineConfig()
		return &defaultConfig
	}

	return &config
}

func (repository *PostgresPipelineConfigRepository) Save(config entity.PipelineConfig) entity.PipelineConfig {
	stagesRaw, _ := json.Marshal(config.Stages)
	territoryRulesRaw, _ := json.Marshal(config.TerritoryRules)
	approvalPoliciesRaw, _ := json.Marshal(config.ApprovalPolicies)

	const statement = `
		INSERT INTO crm.pipeline_configs (
			tenant_id,
			public_id,
			config_name,
			stages_json,
			auto_scoring,
			territory_rules_json,
			approval_policies_json
		)
		VALUES ($1, $2, $3, $4::jsonb, $5, $6::jsonb, $7::jsonb)
		ON CONFLICT (tenant_id)
		DO UPDATE SET
			public_id = EXCLUDED.public_id,
			config_name = EXCLUDED.config_name,
			stages_json = EXCLUDED.stages_json,
			auto_scoring = EXCLUDED.auto_scoring,
			territory_rules_json = EXCLUDED.territory_rules_json,
			approval_policies_json = EXCLUDED.approval_policies_json,
			updated_at = NOW()
	`

	_, _ = repository.database.Exec(statement, repository.tenantID, config.PublicID, config.Name, string(stagesRaw), config.AutoScoring, string(territoryRulesRaw), string(approvalPoliciesRaw))
	return config
}
