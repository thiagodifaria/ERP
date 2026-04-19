package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/api/dto"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

func writeJSON(writer http.ResponseWriter, statusCode int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(payload)
}

func writeError(writer http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(writer, statusCode, dto.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func mapOpportunity(opportunity entity.Opportunity) dto.OpportunityResponse {
	return dto.OpportunityResponse{
		PublicID:     opportunity.PublicID,
		LeadPublicID: opportunity.LeadPublicID,
		Title:        opportunity.Title,
		Stage:        opportunity.Stage,
		OwnerUserID:  opportunity.OwnerUserID,
		AmountCents:  opportunity.AmountCents,
	}
}

func mapProposal(proposal entity.Proposal) dto.ProposalResponse {
	return dto.ProposalResponse{
		PublicID:            proposal.PublicID,
		OpportunityPublicID: proposal.OpportunityPublicID,
		Title:               proposal.Title,
		Status:              proposal.Status,
		AmountCents:         proposal.AmountCents,
	}
}

func mapSale(sale entity.Sale) dto.SaleResponse {
	return dto.SaleResponse{
		PublicID:            sale.PublicID,
		OpportunityPublicID: sale.OpportunityPublicID,
		ProposalPublicID:    sale.ProposalPublicID,
		Status:              sale.Status,
		AmountCents:         sale.AmountCents,
	}
}

func mapInvoice(invoice entity.Invoice) dto.InvoiceResponse {
	return dto.InvoiceResponse{
		PublicID:     invoice.PublicID,
		SalePublicID: invoice.SalePublicID,
		Number:       invoice.Number,
		Status:       invoice.Status,
		AmountCents:  invoice.AmountCents,
		DueDate:      invoice.DueDate,
		PaidAt:       invoice.PaidAt,
	}
}
