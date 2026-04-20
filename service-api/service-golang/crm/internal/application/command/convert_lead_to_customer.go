package command

import (
	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/crm/internal/domain/repository"
)

type ConvertLeadToCustomer struct {
	leadRepository     repository.LeadRepository
	customerRepository repository.CustomerRepository
}

type ConvertLeadToCustomerInput struct {
	LeadPublicID string
}

type ConvertLeadToCustomerResult struct {
	Lead       *entity.Lead
	Customer   *entity.Customer
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
	Conflict   bool
}

func NewConvertLeadToCustomer(
	leadRepository repository.LeadRepository,
	customerRepository repository.CustomerRepository,
) ConvertLeadToCustomer {
	return ConvertLeadToCustomer{
		leadRepository:     leadRepository,
		customerRepository: customerRepository,
	}
}

func (useCase ConvertLeadToCustomer) Execute(input ConvertLeadToCustomerInput) ConvertLeadToCustomerResult {
	lead := useCase.leadRepository.FindByPublicID(input.LeadPublicID)
	if lead == nil {
		return ConvertLeadToCustomerResult{
			ErrorCode: "lead_not_found",
			ErrorText: "Lead was not found.",
			NotFound:  true,
		}
	}

	if lead.Status == "disqualified" {
		return ConvertLeadToCustomerResult{
			ErrorCode:  "lead_not_convertible",
			ErrorText:  "Disqualified leads cannot be converted to customers.",
			BadRequest: true,
		}
	}

	if existing := useCase.customerRepository.FindByEmail(lead.Email); existing != nil {
		return ConvertLeadToCustomerResult{
			ErrorCode: "customer_already_exists",
			ErrorText: "Customer already exists for this lead email.",
			Conflict:  true,
		}
	}

	qualifiedLead := *lead
	if lead.Status != "qualified" {
		updatedLead, err := lead.TransitionTo("qualified")
		if err != nil {
			return ConvertLeadToCustomerResult{
				ErrorCode:  "lead_not_convertible",
				ErrorText:  "Lead is not eligible for customer conversion.",
				BadRequest: true,
			}
		}

		qualifiedLead = useCase.leadRepository.Update(updatedLead)
	}

	customer, err := entity.NewCustomerFromLead(uuid.NewString(), qualifiedLead)
	if err != nil {
		return ConvertLeadToCustomerResult{
			ErrorCode:  "invalid_customer",
			ErrorText:  "Customer payload is invalid.",
			BadRequest: true,
		}
	}

	createdCustomer := useCase.customerRepository.Save(customer)
	return ConvertLeadToCustomerResult{
		Lead:     &qualifiedLead,
		Customer: &createdCustomer,
	}
}
