// ErrorResponse padroniza erros publicos do bootstrap do CRM.
package dto

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
