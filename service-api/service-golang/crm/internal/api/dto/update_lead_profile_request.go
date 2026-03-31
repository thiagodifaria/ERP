// UpdateLeadProfileRequest carries partial profile changes for an existing lead.
package dto

type UpdateLeadProfileRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Source *string `json:"source"`
}
