// DTOs da API descrevem contratos de entrada e saida.
// Objetos de dominio nao devem ser expostos diretamente aqui.
package dto

type HealthResponse struct {
  Service string `json:"service"`
  Status  string `json:"status"`
}
