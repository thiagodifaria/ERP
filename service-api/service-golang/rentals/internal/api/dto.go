package api

import "time"

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type DependencyResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ReadinessResponse struct {
	Service      string               `json:"service"`
	Status       string               `json:"status"`
	Dependencies []DependencyResponse `json:"dependencies"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CreateContractRequest struct {
	TenantSlug       string `json:"tenantSlug"`
	CustomerPublicID string `json:"customerPublicId"`
	Title            string `json:"title"`
	PropertyCode     string `json:"propertyCode"`
	CurrencyCode     string `json:"currencyCode"`
	AmountCents      int64  `json:"amountCents"`
	BillingDay       int    `json:"billingDay"`
	StartsAt         string `json:"startsAt"`
	EndsAt           string `json:"endsAt"`
	RecordedBy       string `json:"recordedBy"`
}

type CreateAdjustmentRequest struct {
	TenantSlug     string `json:"tenantSlug"`
	EffectiveAt    string `json:"effectiveAt"`
	NewAmountCents int64  `json:"newAmountCents"`
	Reason         string `json:"reason"`
	RecordedBy     string `json:"recordedBy"`
}

type TerminateContractRequest struct {
	TenantSlug  string `json:"tenantSlug"`
	EffectiveAt string `json:"effectiveAt"`
	Reason      string `json:"reason"`
	RecordedBy  string `json:"recordedBy"`
}

type CreateAttachmentRequest struct {
	FileName      string `json:"fileName"`
	ContentType   string `json:"contentType"`
	StorageKey    string `json:"storageKey"`
	StorageDriver string `json:"storageDriver"`
	Source        string `json:"source"`
	UploadedBy    string `json:"uploadedBy"`
}

type UpdateChargeStatusRequest struct {
	TenantSlug       string `json:"tenantSlug"`
	Status           string `json:"status"`
	RecordedBy       string `json:"recordedBy"`
	PaidAt           string `json:"paidAt"`
	PaymentReference string `json:"paymentReference"`
}

type ContractResponse struct {
	PublicID          string     `json:"publicId"`
	TenantSlug        string     `json:"tenantSlug"`
	CustomerPublicID  string     `json:"customerPublicId"`
	Title             string     `json:"title"`
	PropertyCode      string     `json:"propertyCode"`
	CurrencyCode      string     `json:"currencyCode"`
	AmountCents       int64      `json:"amountCents"`
	BillingDay        int        `json:"billingDay"`
	StartsAt          string     `json:"startsAt"`
	EndsAt            string     `json:"endsAt"`
	Status            string     `json:"status"`
	TerminatedAt      *time.Time `json:"terminatedAt,omitempty"`
	TerminationReason string     `json:"terminationReason,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

type ChargeResponse struct {
	PublicID         string     `json:"publicId"`
	ContractPublicID string     `json:"contractPublicId"`
	DueDate          string     `json:"dueDate"`
	AmountCents      int64      `json:"amountCents"`
	Status           string     `json:"status"`
	PaidAt           *time.Time `json:"paidAt,omitempty"`
	PaymentReference string     `json:"paymentReference,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type AdjustmentResponse struct {
	PublicID            string    `json:"publicId"`
	ContractPublicID    string    `json:"contractPublicId"`
	EffectiveAt         string    `json:"effectiveAt"`
	PreviousAmountCents int64     `json:"previousAmountCents"`
	NewAmountCents      int64     `json:"newAmountCents"`
	Reason              string    `json:"reason"`
	RecordedBy          string    `json:"recordedBy"`
	CreatedAt           time.Time `json:"createdAt"`
}

type EventResponse struct {
	PublicID         string    `json:"publicId"`
	ContractPublicID string    `json:"contractPublicId"`
	EventCode        string    `json:"eventCode"`
	Summary          string    `json:"summary"`
	RecordedBy       string    `json:"recordedBy"`
	CreatedAt        time.Time `json:"createdAt"`
}

type AttachmentResponse struct {
	PublicID      string    `json:"publicId"`
	TenantSlug    string    `json:"tenantSlug"`
	OwnerType     string    `json:"ownerType"`
	OwnerPublicID string    `json:"ownerPublicId"`
	FileName      string    `json:"fileName"`
	ContentType   string    `json:"contentType"`
	StorageKey    string    `json:"storageKey"`
	StorageDriver string    `json:"storageDriver"`
	Source        string    `json:"source"`
	UploadedBy    string    `json:"uploadedBy"`
	CreatedAt     time.Time `json:"createdAt"`
}

type SummaryResponse struct {
	TenantSlug           string `json:"tenantSlug"`
	TotalContracts       int    `json:"totalContracts"`
	ActiveContracts      int    `json:"activeContracts"`
	TerminatedContracts  int    `json:"terminatedContracts"`
	ScheduledCharges     int    `json:"scheduledCharges"`
	PaidCharges          int    `json:"paidCharges"`
	CancelledCharges     int    `json:"cancelledCharges"`
	Adjustments          int    `json:"adjustments"`
	HistoryEvents        int    `json:"historyEvents"`
	PendingOutbox        int    `json:"pendingOutbox"`
	ScheduledAmountCents int64  `json:"scheduledAmountCents"`
	PaidAmountCents      int64  `json:"paidAmountCents"`
	CancelledAmountCents int64  `json:"cancelledAmountCents"`
}
