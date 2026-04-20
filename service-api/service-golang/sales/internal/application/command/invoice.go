// Commands do ciclo de faturamento inicial do contexto sales.
package command

import (
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type CreateInvoice struct {
	saleRepository    repository.SaleRepository
	invoiceRepository repository.InvoiceRepository
	eventRepository   repository.CommercialEventRepository
	outboxRepository  repository.OutboxEventRepository
}

type CreateInvoiceInput struct {
	SalePublicID string
	Number       string
	DueDate      string
}

type CreateInvoiceResult struct {
	Invoice    *entity.Invoice
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
	Conflict   bool
}

func NewCreateInvoice(
	saleRepository repository.SaleRepository,
	invoiceRepository repository.InvoiceRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) CreateInvoice {
	return CreateInvoice{
		saleRepository:    saleRepository,
		invoiceRepository: invoiceRepository,
		eventRepository:   eventRepository,
		outboxRepository:  outboxRepository,
	}
}

func (useCase CreateInvoice) Execute(input CreateInvoiceInput) CreateInvoiceResult {
	sale := useCase.saleRepository.FindByPublicID(input.SalePublicID)
	if sale == nil {
		return CreateInvoiceResult{
			ErrorCode: "sale_not_found",
			ErrorText: "Sale was not found.",
			NotFound:  true,
		}
	}

	if existing := useCase.invoiceRepository.FindBySalePublicID(sale.PublicID); existing != nil {
		return CreateInvoiceResult{
			ErrorCode: "invoice_already_exists_for_sale",
			ErrorText: "Sale already has an invoice.",
			Conflict:  true,
		}
	}

	if sale.Status == "cancelled" {
		return CreateInvoiceResult{
			ErrorCode:  "sale_not_billable",
			ErrorText:  "Cancelled sales cannot generate invoices.",
			BadRequest: true,
		}
	}

	invoice, err := entity.NewInvoice(newPublicID(), sale.PublicID, input.Number, sale.AmountCents, input.DueDate)
	if err != nil {
		return mapCreateInvoiceValidationError(err)
	}

	created := useCase.invoiceRepository.Save(invoice)

	if sale.Status == "active" {
		invoicedSale, transitionErr := sale.TransitionTo("invoiced")
		if transitionErr == nil {
			useCase.saleRepository.Update(invoicedSale)
		}
	}

	recordCommercialEvent(useCase.eventRepository, "invoice", created.PublicID, "invoice_created", "sales", "Invoice created for sale.")
	appendOutboxEvent(useCase.outboxRepository, "invoice", created.PublicID, "invoice.created", map[string]any{
		"invoicePublicId": created.PublicID,
		"salePublicId":    created.SalePublicID,
		"number":          created.Number,
		"status":          created.Status,
		"amountCents":     created.AmountCents,
		"dueDate":         created.DueDate,
	})
	return CreateInvoiceResult{Invoice: &created}
}

func mapCreateInvoiceValidationError(err error) CreateInvoiceResult {
	switch err {
	case entity.ErrInvoiceNumberInvalid:
		return CreateInvoiceResult{ErrorCode: "invalid_invoice_number", ErrorText: "Invoice number is invalid.", BadRequest: true}
	case entity.ErrInvoiceDueDateInvalid:
		return CreateInvoiceResult{ErrorCode: "invalid_invoice_due_date", ErrorText: "Invoice due date is invalid.", BadRequest: true}
	case entity.ErrInvoiceAmountCentsInvalid:
		return CreateInvoiceResult{ErrorCode: "invalid_invoice_amount", ErrorText: "Invoice amount is invalid.", BadRequest: true}
	default:
		return CreateInvoiceResult{ErrorCode: "invalid_invoice", ErrorText: "Invoice payload is invalid.", BadRequest: true}
	}
}

type UpdateInvoiceStatus struct {
	invoiceRepository repository.InvoiceRepository
	eventRepository   repository.CommercialEventRepository
	outboxRepository  repository.OutboxEventRepository
}

type UpdateInvoiceStatusInput struct {
	PublicID string
	Status   string
}

type UpdateInvoiceStatusResult struct {
	Invoice    *entity.Invoice
	ErrorCode  string
	ErrorText  string
	BadRequest bool
	NotFound   bool
}

func NewUpdateInvoiceStatus(
	invoiceRepository repository.InvoiceRepository,
	eventRepository repository.CommercialEventRepository,
	outboxRepository repository.OutboxEventRepository,
) UpdateInvoiceStatus {
	return UpdateInvoiceStatus{
		invoiceRepository: invoiceRepository,
		eventRepository:   eventRepository,
		outboxRepository:  outboxRepository,
	}
}

func (useCase UpdateInvoiceStatus) Execute(input UpdateInvoiceStatusInput) UpdateInvoiceStatusResult {
	invoice := useCase.invoiceRepository.FindByPublicID(input.PublicID)
	if invoice == nil {
		return UpdateInvoiceStatusResult{
			ErrorCode: "invoice_not_found",
			ErrorText: "Invoice was not found.",
			NotFound:  true,
		}
	}

	updatedInvoice, err := invoice.TransitionTo(input.Status, time.Now().UTC())
	if err != nil {
		switch err {
		case entity.ErrInvoiceStatusInvalid:
			return UpdateInvoiceStatusResult{
				ErrorCode:  "invalid_invoice_status",
				ErrorText:  "Invoice status is invalid.",
				BadRequest: true,
			}
		default:
			return UpdateInvoiceStatusResult{
				ErrorCode:  "invalid_invoice_status_transition",
				ErrorText:  "Invoice status transition is invalid.",
				BadRequest: true,
			}
		}
	}

	saved := useCase.invoiceRepository.Update(updatedInvoice)
	recordCommercialEvent(useCase.eventRepository, "invoice", saved.PublicID, "invoice_status_changed", "sales", "Invoice status transitioned to "+saved.Status+".")
	appendOutboxEvent(useCase.outboxRepository, "invoice", saved.PublicID, "invoice.status_changed", map[string]any{
		"invoicePublicId": saved.PublicID,
		"salePublicId":    saved.SalePublicID,
		"status":          saved.Status,
		"paidAt":          saved.PaidAt,
		"amountCents":     saved.AmountCents,
	})
	return UpdateInvoiceStatusResult{Invoice: &saved}
}
