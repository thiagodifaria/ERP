// LeadNoteResponse expoe o historico operacional publico de um lead.
package dto

import "time"

type LeadNoteResponse struct {
	PublicID     string    `json:"publicId"`
	LeadPublicID string    `json:"leadPublicId"`
	Body         string    `json:"body"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"createdAt"`
}
