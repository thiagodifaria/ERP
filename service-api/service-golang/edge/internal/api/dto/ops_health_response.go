// DTOs da visao operacional do edge sobre os demais servicos.
// Esta resposta prioriza leitura rapida para diagnostico do gateway.
package dto

type OpsHealthSummary struct {
  Total    int `json:"total"`
  Ready    int `json:"ready"`
  Degraded int `json:"degraded"`
}

type ServiceHealthSnapshot struct {
  Name         string               `json:"name"`
  Status       string               `json:"status"`
  Dependencies []DependencyResponse `json:"dependencies"`
}

type OpsHealthResponse struct {
  Service     string                  `json:"service"`
  Status      string                  `json:"status"`
  GeneratedAt string                  `json:"generatedAt"`
  Summary     OpsHealthSummary        `json:"summary"`
  Services    []ServiceHealthSnapshot `json:"services"`
}
