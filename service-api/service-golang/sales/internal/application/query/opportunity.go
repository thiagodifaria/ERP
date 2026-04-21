// Queries de leitura para oportunidades do contexto sales.
package query

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type OpportunityFilters struct {
	Stage            string
	LeadPublicID     string
	CustomerPublicID string
	SaleType         string
	OwnerUserID      string
	Search           string
}

type OpportunitySummary struct {
	Total            int
	TotalAmountCents int64
	ByStage          map[string]int
}

type ListOpportunities struct {
	opportunityRepository repository.OpportunityRepository
}

type GetOpportunitySummary struct {
	opportunityRepository repository.OpportunityRepository
}

type GetOpportunityByPublicID struct {
	opportunityRepository repository.OpportunityRepository
}

func NewListOpportunities(opportunityRepository repository.OpportunityRepository) ListOpportunities {
	return ListOpportunities{opportunityRepository: opportunityRepository}
}

func (useCase ListOpportunities) Execute(filters OpportunityFilters) []entity.Opportunity {
	return applyOpportunityFilters(useCase.opportunityRepository.List(), filters)
}

func NewGetOpportunitySummary(opportunityRepository repository.OpportunityRepository) GetOpportunitySummary {
	return GetOpportunitySummary{opportunityRepository: opportunityRepository}
}

func (useCase GetOpportunitySummary) Execute(filters OpportunityFilters) OpportunitySummary {
	opportunities := applyOpportunityFilters(useCase.opportunityRepository.List(), filters)
	summary := OpportunitySummary{
		ByStage: make(map[string]int),
	}

	for _, opportunity := range opportunities {
		summary.Total++
		summary.TotalAmountCents += opportunity.AmountCents
		summary.ByStage[opportunity.Stage]++
	}

	return summary
}

func NewGetOpportunityByPublicID(opportunityRepository repository.OpportunityRepository) GetOpportunityByPublicID {
	return GetOpportunityByPublicID{opportunityRepository: opportunityRepository}
}

func (useCase GetOpportunityByPublicID) Execute(publicID string) *entity.Opportunity {
	return useCase.opportunityRepository.FindByPublicID(publicID)
}

func applyOpportunityFilters(opportunities []entity.Opportunity, filters OpportunityFilters) []entity.Opportunity {
	response := make([]entity.Opportunity, 0)
	stage := strings.ToLower(strings.TrimSpace(filters.Stage))
	leadPublicID := strings.TrimSpace(filters.LeadPublicID)
	customerPublicID := strings.TrimSpace(filters.CustomerPublicID)
	saleType := strings.ToLower(strings.TrimSpace(filters.SaleType))
	ownerUserID := strings.TrimSpace(filters.OwnerUserID)
	search := strings.ToLower(strings.TrimSpace(filters.Search))

	for _, opportunity := range opportunities {
		if stage != "" && opportunity.Stage != stage {
			continue
		}

		if leadPublicID != "" && opportunity.LeadPublicID != leadPublicID {
			continue
		}

		if customerPublicID != "" && opportunity.CustomerPublicID != customerPublicID {
			continue
		}

		if saleType != "" && opportunity.SaleType != saleType {
			continue
		}

		if ownerUserID != "" && opportunity.OwnerUserID != ownerUserID {
			continue
		}

		if search != "" {
			searchable := strings.ToLower(opportunity.Title + " " + opportunity.LeadPublicID + " " + opportunity.CustomerPublicID + " " + opportunity.SaleType)
			if !strings.Contains(searchable, search) {
				continue
			}
		}

		response = append(response, opportunity)
	}

	return response
}
