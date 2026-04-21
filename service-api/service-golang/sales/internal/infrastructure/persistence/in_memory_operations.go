package persistence

import (
	"strings"
	"sync"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
)

type InMemoryInstallmentRepository struct {
	sync.Mutex
	installments []entity.Installment
}

type InMemoryCommissionRepository struct {
	sync.Mutex
	commissions []entity.Commission
}

type InMemoryPendingItemRepository struct {
	sync.Mutex
	items []entity.PendingItem
}

type InMemoryRenegotiationRepository struct {
	sync.Mutex
	renegotiations []entity.Renegotiation
}

func NewInMemoryInstallmentRepository() *InMemoryInstallmentRepository {
	return &InMemoryInstallmentRepository{installments: []entity.Installment{}}
}

func (repository *InMemoryInstallmentRepository) ListBySalePublicID(salePublicID string) []entity.Installment {
	repository.Lock()
	defer repository.Unlock()

	response := make([]entity.Installment, 0)
	for _, installment := range repository.installments {
		if installment.SalePublicID == strings.TrimSpace(salePublicID) {
			response = append(response, installment)
		}
	}

	return response
}

func (repository *InMemoryInstallmentRepository) Save(installment entity.Installment) entity.Installment {
	repository.Lock()
	defer repository.Unlock()
	repository.installments = append(repository.installments, installment)
	return installment
}

func (repository *InMemoryInstallmentRepository) Update(installment entity.Installment) entity.Installment {
	repository.Lock()
	defer repository.Unlock()
	for index, current := range repository.installments {
		if current.PublicID == installment.PublicID {
			repository.installments[index] = installment
			return installment
		}
	}

	repository.installments = append(repository.installments, installment)
	return installment
}

func NewInMemoryCommissionRepository() *InMemoryCommissionRepository {
	return &InMemoryCommissionRepository{commissions: []entity.Commission{}}
}

func (repository *InMemoryCommissionRepository) ListBySalePublicID(salePublicID string) []entity.Commission {
	repository.Lock()
	defer repository.Unlock()

	response := make([]entity.Commission, 0)
	for _, commission := range repository.commissions {
		if commission.SalePublicID == strings.TrimSpace(salePublicID) {
			response = append(response, commission)
		}
	}

	return response
}

func (repository *InMemoryCommissionRepository) FindByPublicID(publicID string) *entity.Commission {
	repository.Lock()
	defer repository.Unlock()

	for _, commission := range repository.commissions {
		if commission.PublicID == strings.TrimSpace(publicID) {
			copied := commission
			return &copied
		}
	}

	return nil
}

func (repository *InMemoryCommissionRepository) Save(commission entity.Commission) entity.Commission {
	repository.Lock()
	defer repository.Unlock()
	repository.commissions = append(repository.commissions, commission)
	return commission
}

func (repository *InMemoryCommissionRepository) Update(commission entity.Commission) entity.Commission {
	repository.Lock()
	defer repository.Unlock()
	for index, current := range repository.commissions {
		if current.PublicID == commission.PublicID {
			repository.commissions[index] = commission
			return commission
		}
	}

	repository.commissions = append(repository.commissions, commission)
	return commission
}

func NewInMemoryPendingItemRepository() *InMemoryPendingItemRepository {
	return &InMemoryPendingItemRepository{items: []entity.PendingItem{}}
}

func (repository *InMemoryPendingItemRepository) ListBySalePublicID(salePublicID string) []entity.PendingItem {
	repository.Lock()
	defer repository.Unlock()

	response := make([]entity.PendingItem, 0)
	for _, item := range repository.items {
		if item.SalePublicID == strings.TrimSpace(salePublicID) {
			response = append(response, item)
		}
	}

	return response
}

func (repository *InMemoryPendingItemRepository) FindByPublicID(publicID string) *entity.PendingItem {
	repository.Lock()
	defer repository.Unlock()
	for _, item := range repository.items {
		if item.PublicID == strings.TrimSpace(publicID) {
			copied := item
			return &copied
		}
	}
	return nil
}

func (repository *InMemoryPendingItemRepository) Save(item entity.PendingItem) entity.PendingItem {
	repository.Lock()
	defer repository.Unlock()
	repository.items = append(repository.items, item)
	return item
}

func (repository *InMemoryPendingItemRepository) Update(item entity.PendingItem) entity.PendingItem {
	repository.Lock()
	defer repository.Unlock()
	for index, current := range repository.items {
		if current.PublicID == item.PublicID {
			repository.items[index] = item
			return item
		}
	}

	repository.items = append(repository.items, item)
	return item
}

func NewInMemoryRenegotiationRepository() *InMemoryRenegotiationRepository {
	return &InMemoryRenegotiationRepository{renegotiations: []entity.Renegotiation{}}
}

func (repository *InMemoryRenegotiationRepository) ListBySalePublicID(salePublicID string) []entity.Renegotiation {
	repository.Lock()
	defer repository.Unlock()
	response := make([]entity.Renegotiation, 0)
	for _, renegotiation := range repository.renegotiations {
		if renegotiation.SalePublicID == strings.TrimSpace(salePublicID) {
			response = append(response, renegotiation)
		}
	}
	return response
}

func (repository *InMemoryRenegotiationRepository) Save(renegotiation entity.Renegotiation) entity.Renegotiation {
	repository.Lock()
	defer repository.Unlock()
	repository.renegotiations = append(repository.renegotiations, renegotiation)
	return renegotiation
}
