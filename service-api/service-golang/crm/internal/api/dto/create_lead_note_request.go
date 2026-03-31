// CreateLeadNoteRequest carries a new operational note for an existing lead.
package dto

type CreateLeadNoteRequest struct {
	Body     string `json:"body"`
	Category string `json:"category"`
}
