// Adapters em memoria sustentam o bootstrap local ate o driver relacional entrar em cena.
package persistence

import (
	"strings"
	"sync"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

const (
	BootstrapLeadPublicID        = "0195e7a0-7a9c-7c1f-8a44-4a6e70000001"
	BootstrapOwnerUserPublicID   = "0195e7a0-7a9c-7c1f-8a44-4a6e70000011"
	BootstrapOpportunityPublicID = "0195e7a0-7a9c-7c1f-8a44-4a6e71000001"
	BootstrapProposalPublicID    = "0195e7a0-7a9c-7c1f-8a44-4a6e71000002"
	BootstrapSalePublicID        = "0195e7a0-7a9c-7c1f-8a44-4a6e71000003"
	BootstrapInvoicePublicID     = "0195e7a0-7a9c-7c1f-8a44-4a6e71000004"
)

type InMemoryOpportunityRepository struct {
	sync.Mutex
	opportunities []entity.Opportunity
}

type InMemoryProposalRepository struct {
	sync.Mutex
	proposals []entity.Proposal
}

type InMemorySaleRepository struct {
	sync.Mutex
	sales []entity.Sale
}

type InMemoryInvoiceRepository struct {
	sync.Mutex
	invoices []entity.Invoice
}

func NewInMemoryOpportunityRepository() *InMemoryOpportunityRepository {
	opportunity, _ := entity.RestoreOpportunity(
		BootstrapOpportunityPublicID,
		BootstrapLeadPublicID,
		"Bootstrap Ops Annual Plan",
		BootstrapOwnerUserPublicID,
		125000,
		"won",
	)

	return &InMemoryOpportunityRepository{
		opportunities: []entity.Opportunity{opportunity},
	}
}

func (repository *InMemoryOpportunityRepository) List() []entity.Opportunity {
	repository.Lock()
	defer repository.Unlock()

	copied := make([]entity.Opportunity, len(repository.opportunities))
	copy(copied, repository.opportunities)
	return copied
}

func (repository *InMemoryOpportunityRepository) FindByPublicID(publicID string) *entity.Opportunity {
	repository.Lock()
	defer repository.Unlock()

	for _, opportunity := range repository.opportunities {
		if opportunity.PublicID == strings.TrimSpace(publicID) {
			copied := opportunity
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryOpportunityRepository) Save(opportunity entity.Opportunity) entity.Opportunity {
	repository.Lock()
	defer repository.Unlock()

	repository.opportunities = append(repository.opportunities, opportunity)
	return opportunity
}

func (repository *InMemoryOpportunityRepository) Update(opportunity entity.Opportunity) entity.Opportunity {
	repository.Lock()
	defer repository.Unlock()

	for index, current := range repository.opportunities {
		if current.PublicID == opportunity.PublicID {
			repository.opportunities[index] = opportunity
			return opportunity
		}
	}

	repository.opportunities = append(repository.opportunities, opportunity)
	return opportunity
}

func NewInMemoryProposalRepository() *InMemoryProposalRepository {
	proposal, _ := entity.RestoreProposal(
		BootstrapProposalPublicID,
		BootstrapOpportunityPublicID,
		"Bootstrap Ops Annual Proposal",
		125000,
		"accepted",
	)

	return &InMemoryProposalRepository{
		proposals: []entity.Proposal{proposal},
	}
}

func (repository *InMemoryProposalRepository) ListByOpportunityPublicID(opportunityPublicID string) []entity.Proposal {
	repository.Lock()
	defer repository.Unlock()

	response := make([]entity.Proposal, 0)
	for _, proposal := range repository.proposals {
		if proposal.OpportunityPublicID == strings.TrimSpace(opportunityPublicID) {
			response = append(response, proposal)
		}
	}

	return response
}

func (repository *InMemoryProposalRepository) FindByPublicID(publicID string) *entity.Proposal {
	repository.Lock()
	defer repository.Unlock()

	for _, proposal := range repository.proposals {
		if proposal.PublicID == strings.TrimSpace(publicID) {
			copied := proposal
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryProposalRepository) Save(proposal entity.Proposal) entity.Proposal {
	repository.Lock()
	defer repository.Unlock()

	repository.proposals = append(repository.proposals, proposal)
	return proposal
}

func (repository *InMemoryProposalRepository) Update(proposal entity.Proposal) entity.Proposal {
	repository.Lock()
	defer repository.Unlock()

	for index, current := range repository.proposals {
		if current.PublicID == proposal.PublicID {
			repository.proposals[index] = proposal
			return proposal
		}
	}

	repository.proposals = append(repository.proposals, proposal)
	return proposal
}

func NewInMemorySaleRepository() *InMemorySaleRepository {
	sale, _ := entity.RestoreSale(
		BootstrapSalePublicID,
		BootstrapOpportunityPublicID,
		BootstrapProposalPublicID,
		125000,
		"active",
	)

	return &InMemorySaleRepository{
		sales: []entity.Sale{sale},
	}
}

func (repository *InMemorySaleRepository) List() []entity.Sale {
	repository.Lock()
	defer repository.Unlock()

	copied := make([]entity.Sale, len(repository.sales))
	copy(copied, repository.sales)
	return copied
}

func (repository *InMemorySaleRepository) FindByPublicID(publicID string) *entity.Sale {
	repository.Lock()
	defer repository.Unlock()

	for _, sale := range repository.sales {
		if sale.PublicID == strings.TrimSpace(publicID) {
			copied := sale
			return &copied
		}
	}

	return nil
}

func (repository *InMemorySaleRepository) FindByProposalPublicID(proposalPublicID string) *entity.Sale {
	repository.Lock()
	defer repository.Unlock()

	for _, sale := range repository.sales {
		if sale.ProposalPublicID == strings.TrimSpace(proposalPublicID) {
			copied := sale
			return &copied
		}
	}

	return nil
}

func (repository *InMemorySaleRepository) Save(sale entity.Sale) entity.Sale {
	repository.Lock()
	defer repository.Unlock()

	repository.sales = append(repository.sales, sale)
	return sale
}

func NewInMemoryInvoiceRepository() *InMemoryInvoiceRepository {
	invoice, _ := entity.RestoreInvoice(
		BootstrapInvoicePublicID,
		BootstrapSalePublicID,
		"BOOTSTRAP-OPS-INV-0001",
		125000,
		"2026-05-10",
		"sent",
		"",
	)

	return &InMemoryInvoiceRepository{
		invoices: []entity.Invoice{invoice},
	}
}

func (repository *InMemoryInvoiceRepository) List() []entity.Invoice {
	repository.Lock()
	defer repository.Unlock()

	copied := make([]entity.Invoice, len(repository.invoices))
	copy(copied, repository.invoices)
	return copied
}

func (repository *InMemoryInvoiceRepository) FindByPublicID(publicID string) *entity.Invoice {
	repository.Lock()
	defer repository.Unlock()

	for _, invoice := range repository.invoices {
		if invoice.PublicID == strings.TrimSpace(publicID) {
			copied := invoice
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryInvoiceRepository) FindBySalePublicID(salePublicID string) *entity.Invoice {
	repository.Lock()
	defer repository.Unlock()

	for _, invoice := range repository.invoices {
		if invoice.SalePublicID == strings.TrimSpace(salePublicID) {
			copied := invoice
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryInvoiceRepository) Save(invoice entity.Invoice) entity.Invoice {
	repository.Lock()
	defer repository.Unlock()

	repository.invoices = append(repository.invoices, invoice)
	return invoice
}

func (repository *InMemoryInvoiceRepository) Update(invoice entity.Invoice) entity.Invoice {
	repository.Lock()
	defer repository.Unlock()

	for index, current := range repository.invoices {
		if current.PublicID == invoice.PublicID {
			repository.invoices[index] = invoice
			return invoice
		}
	}

	repository.invoices = append(repository.invoices, invoice)
	return invoice
}

func (repository *InMemorySaleRepository) Update(sale entity.Sale) entity.Sale {
	repository.Lock()
	defer repository.Unlock()

	for index, current := range repository.sales {
		if current.PublicID == sale.PublicID {
			repository.sales[index] = sale
			return sale
		}
	}

	repository.sales = append(repository.sales, sale)
	return sale
}
