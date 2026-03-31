// LeadSummaryResponse describes the public CRM pipeline snapshot.
package dto

type LeadSummaryResponse struct {
	Total      int            `json:"total"`
	Assigned   int            `json:"assigned"`
	Unassigned int            `json:"unassigned"`
	ByStatus   map[string]int `json:"byStatus"`
	BySource   map[string]int `json:"bySource"`
}
