// LeadFilters groups the public filters accepted by CRM lead reads.
package query

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type LeadFilters struct {
	Status      string
	Source      string
	OwnerUserID string
	Search      string
	Assigned    string
}

func applyLeadFilters(leads []entity.Lead, filters LeadFilters) []entity.Lead {
	normalized := normalizeLeadFilters(filters)
	filtered := make([]entity.Lead, 0, len(leads))

	for _, lead := range leads {
		if matchesLeadFilters(lead, normalized) {
			filtered = append(filtered, lead)
		}
	}

	return filtered
}

func normalizeLeadFilters(filters LeadFilters) LeadFilters {
	return LeadFilters{
		Status:      strings.ToLower(strings.TrimSpace(filters.Status)),
		Source:      strings.ToLower(strings.TrimSpace(filters.Source)),
		OwnerUserID: strings.TrimSpace(filters.OwnerUserID),
		Search:      strings.ToLower(strings.TrimSpace(filters.Search)),
		Assigned:    strings.ToLower(strings.TrimSpace(filters.Assigned)),
	}
}

func matchesLeadFilters(lead entity.Lead, filters LeadFilters) bool {
	if filters.Status != "" && strings.ToLower(strings.TrimSpace(lead.Status)) != filters.Status {
		return false
	}

	if filters.Source != "" && strings.ToLower(strings.TrimSpace(lead.Source)) != filters.Source {
		return false
	}

	if filters.OwnerUserID != "" && strings.TrimSpace(lead.OwnerUserID) != filters.OwnerUserID {
		return false
	}

	if filters.Assigned == "true" && strings.TrimSpace(lead.OwnerUserID) == "" {
		return false
	}

	if filters.Assigned == "false" && strings.TrimSpace(lead.OwnerUserID) != "" {
		return false
	}

	if filters.Search != "" && !matchesLeadSearch(lead, filters.Search) {
		return false
	}

	return true
}

func matchesLeadSearch(lead entity.Lead, search string) bool {
	haystacks := []string{
		strings.ToLower(lead.PublicID),
		strings.ToLower(lead.Name),
		strings.ToLower(lead.Email),
		strings.ToLower(lead.Source),
	}

	for _, haystack := range haystacks {
		if strings.Contains(haystack, search) {
			return true
		}
	}

	return false
}
