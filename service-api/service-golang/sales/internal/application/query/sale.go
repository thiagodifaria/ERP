// Queries de leitura para vendas do contexto sales.
package query

import (
	"strings"

	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/sales/internal/domain/repository"
)

type SaleFilters struct {
	Status string
}

type SaleSummary struct {
	Total              int
	BookedRevenueCents int64
	ByStatus           map[string]int
}

type ListSales struct {
	saleRepository repository.SaleRepository
}

type GetSaleSummary struct {
	saleRepository repository.SaleRepository
}

type GetSaleByPublicID struct {
	saleRepository repository.SaleRepository
}

func NewListSales(saleRepository repository.SaleRepository) ListSales {
	return ListSales{saleRepository: saleRepository}
}

func (useCase ListSales) Execute(filters SaleFilters) []entity.Sale {
	return applySaleFilters(useCase.saleRepository.List(), filters)
}

func NewGetSaleSummary(saleRepository repository.SaleRepository) GetSaleSummary {
	return GetSaleSummary{saleRepository: saleRepository}
}

func (useCase GetSaleSummary) Execute(filters SaleFilters) SaleSummary {
	sales := applySaleFilters(useCase.saleRepository.List(), filters)
	summary := SaleSummary{
		ByStatus: make(map[string]int),
	}

	for _, sale := range sales {
		summary.Total++
		summary.ByStatus[sale.Status]++
		if sale.Status != "cancelled" {
			summary.BookedRevenueCents += sale.AmountCents
		}
	}

	return summary
}

func NewGetSaleByPublicID(saleRepository repository.SaleRepository) GetSaleByPublicID {
	return GetSaleByPublicID{saleRepository: saleRepository}
}

func (useCase GetSaleByPublicID) Execute(publicID string) *entity.Sale {
	return useCase.saleRepository.FindByPublicID(publicID)
}

func applySaleFilters(sales []entity.Sale, filters SaleFilters) []entity.Sale {
	status := strings.ToLower(strings.TrimSpace(filters.Status))
	response := make([]entity.Sale, 0)

	for _, sale := range sales {
		if status != "" && sale.Status != status {
			continue
		}

		response = append(response, sale)
	}

	return response
}
