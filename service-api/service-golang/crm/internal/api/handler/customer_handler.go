package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type CustomerHandler struct {
	repositories repository.TenantRepositoryFactory
}

func NewCustomerHandler(repositories repository.TenantRepositoryFactory) CustomerHandler {
	return CustomerHandler{repositories: repositories}
}

func (handler CustomerHandler) List(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request)
	if bundle.CustomerRepository == nil {
		return
	}

	customers := query.NewListCustomers(bundle.CustomerRepository).Execute()
	response := make([]dto.CustomerResponse, 0, len(customers))
	for _, customer := range customers {
		response = append(response, mapCustomer(customer))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler CustomerHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	bundle, _ := handler.resolveRepositories(writer, request)
	if bundle.CustomerRepository == nil {
		return
	}

	customer := query.NewGetCustomerByPublicID(bundle.CustomerRepository).Execute(request.PathValue("publicId"))
	if customer == nil {
		writeNotFound(writer, "customer_not_found", "Customer was not found.")
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(mapCustomer(*customer))
}

func mapCustomer(customer entity.Customer) dto.CustomerResponse {
	return dto.CustomerResponse{
		PublicID:     customer.PublicID,
		LeadPublicID: customer.LeadPublicID,
		Name:         customer.Name,
		Email:        customer.Email,
		Source:       customer.Source,
		Status:       customer.Status,
		OwnerUserID:  customer.OwnerUserID,
	}
}
