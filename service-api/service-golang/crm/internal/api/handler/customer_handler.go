package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/application/query"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
)

type CustomerHandler struct {
	listCustomers        query.ListCustomers
	getCustomerByPublicID query.GetCustomerByPublicID
}

func NewCustomerHandler(
	listCustomers query.ListCustomers,
	getCustomerByPublicID query.GetCustomerByPublicID,
) CustomerHandler {
	return CustomerHandler{
		listCustomers:        listCustomers,
		getCustomerByPublicID: getCustomerByPublicID,
	}
}

func (handler CustomerHandler) List(writer http.ResponseWriter, _ *http.Request) {
	customers := handler.listCustomers.Execute()
	response := make([]dto.CustomerResponse, 0, len(customers))
	for _, customer := range customers {
		response = append(response, mapCustomer(customer))
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(response)
}

func (handler CustomerHandler) GetByPublicID(writer http.ResponseWriter, request *http.Request) {
	customer := handler.getCustomerByPublicID.Execute(request.PathValue("publicId"))
	if customer == nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(writer).Encode(dto.ErrorResponse{
			Code:    "customer_not_found",
			Message: "Customer was not found.",
		})
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
