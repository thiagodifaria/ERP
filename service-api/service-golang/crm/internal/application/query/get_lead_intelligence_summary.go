package query

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"

func GetLeadIntelligenceSummary(factory repository.TenantRepositoryFactory, tenantSlug string) (map[string]any, error) {
	tenantRepositories, err := factory.ForTenant(tenantSlug)
	if err != nil {
		return nil, err
	}

	leads := tenantRepositories.LeadRepository.List()
	buckets := map[string]int{"cold": 0, "warm": 0, "hot": 0}
	totalScore := 0

	for _, lead := range leads {
		score := scoreLead(lead)
		totalScore += score
		switch {
		case score >= 80:
			buckets["hot"]++
		case score >= 50:
			buckets["warm"]++
		default:
			buckets["cold"]++
		}
	}

	averageScore := 0
	if len(leads) > 0 {
		averageScore = totalScore / len(leads)
	}

	config := tenantRepositories.PipelineConfigRepository.Get()
	stageCount := 0
	autoScoring := false
	pipelineName := ""
	if config != nil {
		stageCount = len(config.Stages)
		autoScoring = config.AutoScoring
		pipelineName = config.Name
	}

	return map[string]any{
		"tenantSlug":    tenantSlug,
		"pipelineName":  pipelineName,
		"stageCount":    stageCount,
		"autoScoring":   autoScoring,
		"averageScore":  averageScore,
		"bucketSummary": buckets,
		"leadCount":     len(leads),
	}, nil
}

func scoreLead(lead entity.Lead) int {
	score := 20
	switch lead.Source {
	case "meta_ads", "whatsapp", "telegram":
		score += 20
	case "referral":
		score += 30
	default:
		score += 10
	}

	switch lead.Status {
	case "contacted":
		score += 15
	case "qualified":
		score += 40
	case "disqualified":
		score = 5
	}

	if lead.OwnerUserID != "" {
		score += 10
	}

	if score > 100 {
		return 100
	}

	return score
}
