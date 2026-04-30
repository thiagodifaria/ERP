package persistence

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/entity"
	repo "github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
)

type InMemoryContractRepository struct {
	contracts   []entity.Contract
	charges     []entity.Charge
	events      []entity.Event
	adjustments []entity.Adjustment
	outbox      []entity.OutboxEvent
}

func NewInMemoryContractRepository() *InMemoryContractRepository {
	return &InMemoryContractRepository{
		contracts:   []entity.Contract{},
		charges:     []entity.Charge{},
		events:      []entity.Event{},
		adjustments: []entity.Adjustment{},
		outbox:      []entity.OutboxEvent{},
	}
}

func (repository *InMemoryContractRepository) List(filters repo.ContractFilters) []entity.Contract {
	response := make([]entity.Contract, 0)
	tenantSlug := strings.ToLower(strings.TrimSpace(filters.TenantSlug))
	status := strings.ToLower(strings.TrimSpace(filters.Status))
	customerPublicID := strings.TrimSpace(filters.CustomerPublicID)

	for _, contract := range repository.contracts {
		if tenantSlug != "" && contract.TenantSlug != tenantSlug {
			continue
		}
		if status != "" && contract.Status != status {
			continue
		}
		if customerPublicID != "" && contract.CustomerPublicID != customerPublicID {
			continue
		}
		response = append(response, contract)
	}

	return response
}

func (repository *InMemoryContractRepository) Summary(tenantSlug string) repo.ContractSummary {
	summary := repo.ContractSummary{TenantSlug: strings.ToLower(strings.TrimSpace(tenantSlug))}

	for _, contract := range repository.contracts {
		if summary.TenantSlug != "" && contract.TenantSlug != summary.TenantSlug {
			continue
		}
		summary.TotalContracts++
		switch contract.Status {
		case "active":
			summary.ActiveContracts++
		case "terminated":
			summary.TerminatedContracts++
		}
	}

	for _, charge := range repository.charges {
		contract, ok := repository.FindByPublicID(summary.TenantSlug, charge.ContractPublicID)
		if !ok || (summary.TenantSlug != "" && contract.TenantSlug != summary.TenantSlug) {
			continue
		}
		if charge.Status == "scheduled" {
			summary.ScheduledCharges++
			summary.ScheduledAmountCents += charge.AmountCents
		}
		if charge.Status == "cancelled" {
			summary.CancelledCharges++
			summary.CancelledAmountCents += charge.AmountCents
		}
	}

	for _, adjustment := range repository.adjustments {
		contract, ok := repository.FindByPublicID(summary.TenantSlug, adjustment.ContractPublicID)
		if ok && (summary.TenantSlug == "" || contract.TenantSlug == summary.TenantSlug) {
			summary.Adjustments++
		}
	}

	for _, event := range repository.events {
		contract, ok := repository.FindByPublicID(summary.TenantSlug, event.ContractPublicID)
		if ok && (summary.TenantSlug == "" || contract.TenantSlug == summary.TenantSlug) {
			summary.HistoryEvents++
		}
	}

	for _, outbox := range repository.outbox {
		contract, ok := repository.FindByPublicID(summary.TenantSlug, outbox.AggregatePublicID)
		if ok && (summary.TenantSlug == "" || contract.TenantSlug == summary.TenantSlug) && outbox.Status == "pending" {
			summary.PendingOutbox++
		}
	}

	return summary
}

func (repository *InMemoryContractRepository) FindByPublicID(tenantSlug string, publicID string) (entity.Contract, bool) {
	normalizedTenant := strings.ToLower(strings.TrimSpace(tenantSlug))
	normalizedPublicID := strings.TrimSpace(publicID)

	for _, contract := range repository.contracts {
		if normalizedTenant != "" && contract.TenantSlug != normalizedTenant {
			continue
		}
		if contract.PublicID == normalizedPublicID {
			return contract, true
		}
	}

	return entity.Contract{}, false
}

func (repository *InMemoryContractRepository) Create(contract entity.Contract, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) entity.Contract {
	repository.contracts = append(repository.contracts, contract)
	repository.charges = append(repository.charges, charges...)
	repository.events = append(repository.events, event)
	repository.outbox = append(repository.outbox, outbox)
	return contract
}

func (repository *InMemoryContractRepository) ListCharges(tenantSlug string, contractPublicID string, status string) []entity.Charge {
	contract, ok := repository.FindByPublicID(tenantSlug, contractPublicID)
	if !ok {
		return []entity.Charge{}
	}

	response := make([]entity.Charge, 0)
	normalizedStatus := strings.ToLower(strings.TrimSpace(status))
	for _, charge := range repository.charges {
		if charge.ContractPublicID != contract.PublicID {
			continue
		}
		if normalizedStatus != "" && charge.Status != normalizedStatus {
			continue
		}
		response = append(response, charge)
	}

	return response
}

func (repository *InMemoryContractRepository) ListEvents(tenantSlug string, contractPublicID string) []entity.Event {
	contract, ok := repository.FindByPublicID(tenantSlug, contractPublicID)
	if !ok {
		return []entity.Event{}
	}

	response := make([]entity.Event, 0)
	for _, event := range repository.events {
		if event.ContractPublicID == contract.PublicID {
			response = append(response, event)
		}
	}

	return response
}

func (repository *InMemoryContractRepository) ListAdjustments(tenantSlug string, contractPublicID string) []entity.Adjustment {
	contract, ok := repository.FindByPublicID(tenantSlug, contractPublicID)
	if !ok {
		return []entity.Adjustment{}
	}

	response := make([]entity.Adjustment, 0)
	for _, adjustment := range repository.adjustments {
		if adjustment.ContractPublicID == contract.PublicID {
			response = append(response, adjustment)
		}
	}

	return response
}

func (repository *InMemoryContractRepository) SaveAdjustment(tenantSlug string, contract entity.Contract, adjustment entity.Adjustment, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) (entity.Contract, bool) {
	for index := range repository.contracts {
		if repository.contracts[index].PublicID != contract.PublicID {
			continue
		}
		if strings.ToLower(strings.TrimSpace(tenantSlug)) != "" && repository.contracts[index].TenantSlug != strings.ToLower(strings.TrimSpace(tenantSlug)) {
			continue
		}
		repository.contracts[index] = contract
		repository.adjustments = append(repository.adjustments, adjustment)
		repository.events = append(repository.events, event)
		repository.outbox = append(repository.outbox, outbox)
		for chargeIndex := range repository.charges {
			for _, updatedCharge := range charges {
				if repository.charges[chargeIndex].PublicID == updatedCharge.PublicID {
					repository.charges[chargeIndex] = updatedCharge
				}
			}
		}
		return contract, true
	}

	return entity.Contract{}, false
}

func (repository *InMemoryContractRepository) SaveTermination(tenantSlug string, contract entity.Contract, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) (entity.Contract, bool) {
	for index := range repository.contracts {
		if repository.contracts[index].PublicID != contract.PublicID {
			continue
		}
		if strings.ToLower(strings.TrimSpace(tenantSlug)) != "" && repository.contracts[index].TenantSlug != strings.ToLower(strings.TrimSpace(tenantSlug)) {
			continue
		}
		repository.contracts[index] = contract
		repository.events = append(repository.events, event)
		repository.outbox = append(repository.outbox, outbox)
		for chargeIndex := range repository.charges {
			for _, updatedCharge := range charges {
				if repository.charges[chargeIndex].PublicID == updatedCharge.PublicID {
					repository.charges[chargeIndex] = updatedCharge
				}
			}
		}
		return contract, true
	}

	return entity.Contract{}, false
}
