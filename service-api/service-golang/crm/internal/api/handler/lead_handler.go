// LeadHandler expõe a API minima de leads do CRM.
// Regras de negocio devem ficar nos casos de uso.
package handler

import (
  "encoding/json"
  "net/http"

  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/command"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
  "github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type LeadHandler struct {
  listLeads  query.ListLeads
  createLead command.CreateLead
}

func NewLeadHandler(listLeads query.ListLeads, createLead command.CreateLead) LeadHandler {
  return LeadHandler{
    listLeads:  listLeads,
    createLead: createLead,
  }
}

func (handler LeadHandler) List(writer http.ResponseWriter, request *http.Request) {
  leads := handler.listLeads.Execute()
  response := make([]dto.LeadResponse, 0, len(leads))

  for _, lead := range leads {
    response = append(response, mapLead(lead))
  }

  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  _ = json.NewEncoder(writer).Encode(response)
}

func (handler LeadHandler) Create(writer http.ResponseWriter, request *http.Request) {
  var payload dto.CreateLeadRequest
  if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(http.StatusBadRequest)
    _ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
      Code:    "invalid_json",
      Message: "Request body is invalid.",
    })
    return
  }

  result := handler.createLead.Execute(command.CreateLeadInput{
    Name:        payload.Name,
    Email:       payload.Email,
    Source:      payload.Source,
    OwnerUserID: payload.OwnerUserID,
  })

  writer.Header().Set("Content-Type", "application/json")

  switch {
  case result.BadRequest:
    writer.WriteHeader(http.StatusBadRequest)
    _ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
      Code:    result.ErrorCode,
      Message: result.ErrorText,
    })
  case result.Conflict:
    writer.WriteHeader(http.StatusConflict)
    _ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
      Code:    result.ErrorCode,
      Message: result.ErrorText,
    })
  default:
    writer.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(writer).Encode(mapLead(*result.Lead))
  }
}

func mapLead(lead entity.Lead) dto.LeadResponse {
  return dto.LeadResponse{
    PublicID:    lead.PublicID,
    Name:        lead.Name,
    Email:       lead.Email,
    Source:      lead.Source,
    Status:      lead.Status,
    OwnerUserID: lead.OwnerUserID,
  }
}
