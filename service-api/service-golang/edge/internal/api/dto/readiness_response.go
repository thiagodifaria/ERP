// ReadinessResponse descreve a disponibilidade detalhada do servico edge.
// Dependencias podem evoluir sem alterar o contrato de health basico.
package dto

type ReadinessResponse struct {
	Service      string               `json:"service"`
	Status       string               `json:"status"`
	Dependencies []DependencyResponse `json:"dependencies"`
}

type DependencyResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
