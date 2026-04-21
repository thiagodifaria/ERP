package query

import (
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type ListInstallmentsBySale struct {
	repository repository.InstallmentRepository
}

type ListCommissionsBySale struct {
	repository repository.CommissionRepository
}

type ListPendingItemsBySale struct {
	repository repository.PendingItemRepository
}

type ListRenegotiationsBySale struct {
	repository repository.RenegotiationRepository
}

func NewListInstallmentsBySale(repo repository.InstallmentRepository) ListInstallmentsBySale {
	return ListInstallmentsBySale{repository: repo}
}

func (useCase ListInstallmentsBySale) Execute(salePublicID string) []entity.Installment {
	return useCase.repository.ListBySalePublicID(salePublicID)
}

func NewListCommissionsBySale(repo repository.CommissionRepository) ListCommissionsBySale {
	return ListCommissionsBySale{repository: repo}
}

func (useCase ListCommissionsBySale) Execute(salePublicID string) []entity.Commission {
	return useCase.repository.ListBySalePublicID(salePublicID)
}

func NewListPendingItemsBySale(repo repository.PendingItemRepository) ListPendingItemsBySale {
	return ListPendingItemsBySale{repository: repo}
}

func (useCase ListPendingItemsBySale) Execute(salePublicID string) []entity.PendingItem {
	return useCase.repository.ListBySalePublicID(salePublicID)
}

func NewListRenegotiationsBySale(repo repository.RenegotiationRepository) ListRenegotiationsBySale {
	return ListRenegotiationsBySale{repository: repo}
}

func (useCase ListRenegotiationsBySale) Execute(salePublicID string) []entity.Renegotiation {
	return useCase.repository.ListBySalePublicID(salePublicID)
}
