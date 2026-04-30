package repository

import (
	"time"

	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/entity"
)

type ContractFilters struct {
	TenantSlug       string
	Status           string
	CustomerPublicID string
}

type ContractSummary struct {
	TenantSlug           string
	TotalContracts       int
	ActiveContracts      int
	TerminatedContracts  int
	ScheduledCharges     int
	CancelledCharges     int
	Adjustments          int
	HistoryEvents        int
	PendingOutbox        int
	ScheduledAmountCents int64
	CancelledAmountCents int64
}

type AttachmentRecord struct {
	PublicID      string
	TenantSlug    string
	OwnerType     string
	OwnerPublicID string
	FileName      string
	ContentType   string
	StorageKey    string
	StorageDriver string
	Source        string
	UploadedBy    string
	CreatedAt     time.Time
}

type CreateAttachmentInput struct {
	TenantSlug    string
	OwnerType     string
	OwnerPublicID string
	FileName      string
	ContentType   string
	StorageKey    string
	StorageDriver string
	Source        string
	UploadedBy    string
}

type AttachmentGateway interface {
	List(tenantSlug string, ownerType string, ownerPublicID string) ([]AttachmentRecord, error)
	Create(input CreateAttachmentInput) (*AttachmentRecord, error)
}

type ContractRepository interface {
	List(filters ContractFilters) []entity.Contract
	Summary(tenantSlug string) ContractSummary
	FindByPublicID(tenantSlug string, publicID string) (entity.Contract, bool)
	Create(contract entity.Contract, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) entity.Contract
	ListCharges(tenantSlug string, contractPublicID string, status string) []entity.Charge
	ListEvents(tenantSlug string, contractPublicID string) []entity.Event
	ListAdjustments(tenantSlug string, contractPublicID string) []entity.Adjustment
	SaveAdjustment(tenantSlug string, contract entity.Contract, adjustment entity.Adjustment, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) (entity.Contract, bool)
	SaveTermination(tenantSlug string, contract entity.Contract, charges []entity.Charge, event entity.Event, outbox entity.OutboxEvent) (entity.Contract, bool)
}
