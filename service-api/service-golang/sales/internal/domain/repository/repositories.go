// Repository define as abstractions minimas de persistencia do contexto sales.
// Regras de negocio devem depender destas interfaces.
package repository

import "github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"

type OpportunityRepository interface {
	List() []entity.Opportunity
	FindByPublicID(publicID string) *entity.Opportunity
	Save(opportunity entity.Opportunity) entity.Opportunity
	Update(opportunity entity.Opportunity) entity.Opportunity
}

type ProposalRepository interface {
	ListByOpportunityPublicID(opportunityPublicID string) []entity.Proposal
	FindByPublicID(publicID string) *entity.Proposal
	Save(proposal entity.Proposal) entity.Proposal
	Update(proposal entity.Proposal) entity.Proposal
}

type SaleRepository interface {
	List() []entity.Sale
	FindByPublicID(publicID string) *entity.Sale
	FindByProposalPublicID(proposalPublicID string) *entity.Sale
	Save(sale entity.Sale) entity.Sale
	Update(sale entity.Sale) entity.Sale
}

type InvoiceRepository interface {
	List() []entity.Invoice
	FindByPublicID(publicID string) *entity.Invoice
	FindBySalePublicID(salePublicID string) *entity.Invoice
	Save(invoice entity.Invoice) entity.Invoice
	Update(invoice entity.Invoice) entity.Invoice
}

type CommercialEventRepository interface {
	ListByAggregate(aggregateType string, aggregatePublicID string) []entity.CommercialEvent
	Save(event entity.CommercialEvent) entity.CommercialEvent
}

type OutboxEventRepository interface {
	ListPending(limit int) []entity.OutboxEvent
	FindByPublicID(publicID string) *entity.OutboxEvent
	Save(event entity.OutboxEvent) entity.OutboxEvent
	Update(event entity.OutboxEvent) entity.OutboxEvent
}
