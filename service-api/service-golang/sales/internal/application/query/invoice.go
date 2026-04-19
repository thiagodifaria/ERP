// Queries de leitura para invoices e cobranca do contexto sales.
package query

import (
	"strings"
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type InvoiceFilters struct {
	Status string
}

type InvoiceSummary struct {
	Total              int
	OpenAmountCents    int64
	PaidAmountCents    int64
	OverdueAmountCents int64
	OverdueCount       int
	ByStatus           map[string]int
}

type ListInvoices struct {
	invoiceRepository repository.InvoiceRepository
}

type GetInvoiceSummary struct {
	invoiceRepository repository.InvoiceRepository
}

type GetInvoiceByPublicID struct {
	invoiceRepository repository.InvoiceRepository
}

func NewListInvoices(invoiceRepository repository.InvoiceRepository) ListInvoices {
	return ListInvoices{invoiceRepository: invoiceRepository}
}

func (useCase ListInvoices) Execute(filters InvoiceFilters) []entity.Invoice {
	return applyInvoiceFilters(useCase.invoiceRepository.List(), filters)
}

func NewGetInvoiceSummary(invoiceRepository repository.InvoiceRepository) GetInvoiceSummary {
	return GetInvoiceSummary{invoiceRepository: invoiceRepository}
}

func (useCase GetInvoiceSummary) Execute(filters InvoiceFilters) InvoiceSummary {
	invoices := applyInvoiceFilters(useCase.invoiceRepository.List(), filters)
	summary := InvoiceSummary{
		ByStatus: make(map[string]int),
	}

	for _, invoice := range invoices {
		summary.Total++
		summary.ByStatus[invoice.Status]++

		switch invoice.Status {
		case "paid":
			summary.PaidAmountCents += invoice.AmountCents
		case "cancelled":
			continue
		default:
			summary.OpenAmountCents += invoice.AmountCents
			if invoice.IsOverdue(time.Now().UTC()) {
				summary.OverdueCount++
				summary.OverdueAmountCents += invoice.AmountCents
			}
		}
	}

	return summary
}

func NewGetInvoiceByPublicID(invoiceRepository repository.InvoiceRepository) GetInvoiceByPublicID {
	return GetInvoiceByPublicID{invoiceRepository: invoiceRepository}
}

func (useCase GetInvoiceByPublicID) Execute(publicID string) *entity.Invoice {
	return useCase.invoiceRepository.FindByPublicID(publicID)
}

func applyInvoiceFilters(invoices []entity.Invoice, filters InvoiceFilters) []entity.Invoice {
	status := strings.ToLower(strings.TrimSpace(filters.Status))
	response := make([]entity.Invoice, 0)

	for _, invoice := range invoices {
		if status != "" && invoice.Status != status {
			continue
		}

		response = append(response, invoice)
	}

	return response
}
