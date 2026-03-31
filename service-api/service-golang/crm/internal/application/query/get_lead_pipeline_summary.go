// GetLeadPipelineSummary exposes pipeline totals from the filtered lead set.
package query

import "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"

type LeadPipelineSummary struct {
	Total      int
	Assigned   int
	Unassigned int
	ByStatus   map[string]int
	BySource   map[string]int
}

type GetLeadPipelineSummary struct {
	leadRepository repository.LeadRepository
}

func NewGetLeadPipelineSummary(leadRepository repository.LeadRepository) GetLeadPipelineSummary {
	return GetLeadPipelineSummary{leadRepository: leadRepository}
}

func (useCase GetLeadPipelineSummary) Execute(filters LeadFilters) LeadPipelineSummary {
	leads := applyLeadFilters(useCase.leadRepository.List(), filters)
	summary := LeadPipelineSummary{
		ByStatus: make(map[string]int),
		BySource: make(map[string]int),
	}

	for _, lead := range leads {
		summary.Total++
		summary.ByStatus[lead.Status]++
		summary.BySource[lead.Source]++

		if lead.OwnerUserID == "" {
			summary.Unassigned++
			continue
		}

		summary.Assigned++
	}

	return summary
}
