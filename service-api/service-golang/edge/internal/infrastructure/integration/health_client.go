// Package integration concentra chamadas HTTP simples para dependencias externas.
// O edge usa esta camada para observar a disponibilidade dos demais servicos.
package integration

import (
  "context"
  "encoding/json"
  "net/http"
  "strings"
  "time"

  "github.com/thiagodifaria/erp/service-api/service-golang/edge/internal/api/dto"
)

type ServiceEndpoint struct {
  Name    string
  BaseURL string
}

type HealthChecker interface {
  Check(ctx context.Context, endpoint ServiceEndpoint) dto.DependencyResponse
}

type HTTPHealthChecker struct {
  client *http.Client
}

func NewHTTPHealthChecker(timeout time.Duration) *HTTPHealthChecker {
  return &HTTPHealthChecker{
    client: &http.Client{Timeout: timeout},
  }
}

func (checker *HTTPHealthChecker) Check(ctx context.Context, endpoint ServiceEndpoint) dto.DependencyResponse {
  if strings.TrimSpace(endpoint.BaseURL) == "" {
    return dto.DependencyResponse{
      Name:   endpoint.Name,
      Status: "not_configured",
    }
  }

  request, err := http.NewRequestWithContext(
    ctx,
    http.MethodGet,
    strings.TrimRight(endpoint.BaseURL, "/")+"/health/ready",
    nil,
  )
  if err != nil {
    return dto.DependencyResponse{
      Name:   endpoint.Name,
      Status: "not_ready",
    }
  }

  response, err := checker.client.Do(request)
  if err != nil {
    return dto.DependencyResponse{
      Name:   endpoint.Name,
      Status: "not_ready",
    }
  }
  defer response.Body.Close()

  if response.StatusCode != http.StatusOK {
    return dto.DependencyResponse{
      Name:   endpoint.Name,
      Status: "not_ready",
    }
  }

  payload := dto.HealthResponse{}
  if err := json.NewDecoder(response.Body).Decode(&payload); err == nil && payload.Status != "" && payload.Status != "ready" && payload.Status != "live" {
    return dto.DependencyResponse{
      Name:   endpoint.Name,
      Status: "not_ready",
    }
  }

  return dto.DependencyResponse{
    Name:   endpoint.Name,
    Status: "ready",
  }
}
