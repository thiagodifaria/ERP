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
		PublicID:         opportunity.PublicID,
		LeadPublicID:     opportunity.LeadPublicID,
		CustomerPublicID: opportunity.CustomerPublicID,
		Title:            opportunity.Title,
		Stage:            opportunity.Stage,
		SaleType:         opportunity.SaleType,
		OwnerUserID:      opportunity.OwnerUserID,
		AmountCents:      opportunity.AmountCents,
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
		CustomerPublicID:    sale.CustomerPublicID,
		OwnerUserID:         sale.OwnerUserID,
		SaleType:            sale.SaleType,
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

func mapInstallment(installment entity.Installment) dto.InstallmentResponse {
	return dto.InstallmentResponse{
		PublicID:       installment.PublicID,
		SalePublicID:   installment.SalePublicID,
		SequenceNumber: installment.SequenceNumber,
		AmountCents:    installment.AmountCents,
		DueDate:        installment.DueDate,
		Status:         installment.Status,
	}
}

func mapCommission(commission entity.Commission) dto.CommissionResponse {
	return dto.CommissionResponse{
		PublicID:        commission.PublicID,
		SalePublicID:    commission.SalePublicID,
		RecipientUserID: commission.RecipientUserID,
		RoleCode:        commission.RoleCode,
		RateBps:         commission.RateBps,
		AmountCents:     commission.AmountCents,
		Status:          commission.Status,
	}
}

func mapPendingItem(item entity.PendingItem) dto.PendingItemResponse {
	return dto.PendingItemResponse{
		PublicID:     item.PublicID,
		SalePublicID: item.SalePublicID,
		Code:         item.Code,
		Summary:      item.Summary,
		Status:       item.Status,
		ResolvedAt:   item.ResolvedAt,
	}
}

func mapRenegotiation(renegotiation entity.Renegotiation) dto.RenegotiationResponse {
	return dto.RenegotiationResponse{
		PublicID:            renegotiation.PublicID,
		SalePublicID:        renegotiation.SalePublicID,
		Reason:              renegotiation.Reason,
		PreviousAmountCents: renegotiation.PreviousAmountCents,
		NewAmountCents:      renegotiation.NewAmountCents,
		Status:              renegotiation.Status,
		AppliedAt:           renegotiation.AppliedAt,
	}
}
