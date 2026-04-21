package dto

type InstallmentResponse struct {
	PublicID       string `json:"publicId"`
	SalePublicID   string `json:"salePublicId"`
	SequenceNumber int    `json:"sequenceNumber"`
	AmountCents    int64  `json:"amountCents"`
	DueDate        string `json:"dueDate"`
	Status         string `json:"status"`
}

type CreateInstallmentsRequest struct {
	Installments []CreateInstallmentItem `json:"installments"`
}

type CreateInstallmentItem struct {
	AmountCents int64  `json:"amountCents"`
	DueDate     string `json:"dueDate"`
}

type CommissionResponse struct {
	PublicID        string `json:"publicId"`
	SalePublicID    string `json:"salePublicId"`
	RecipientUserID string `json:"recipientUserId"`
	RoleCode        string `json:"roleCode"`
	RateBps         int    `json:"rateBps"`
	AmountCents     int64  `json:"amountCents"`
	Status          string `json:"status"`
}

type CreateCommissionRequest struct {
	RecipientUserID string `json:"recipientUserId"`
	RoleCode        string `json:"roleCode"`
	RateBps         int    `json:"rateBps"`
}

type PendingItemResponse struct {
	PublicID     string `json:"publicId"`
	SalePublicID string `json:"salePublicId"`
	Code         string `json:"code"`
	Summary      string `json:"summary"`
	Status       string `json:"status"`
	ResolvedAt   string `json:"resolvedAt"`
}

type CreatePendingItemRequest struct {
	Code    string `json:"code"`
	Summary string `json:"summary"`
}

type RenegotiationResponse struct {
	PublicID            string `json:"publicId"`
	SalePublicID        string `json:"salePublicId"`
	Reason              string `json:"reason"`
	PreviousAmountCents int64  `json:"previousAmountCents"`
	NewAmountCents      int64  `json:"newAmountCents"`
	Status              string `json:"status"`
	AppliedAt           string `json:"appliedAt"`
}

type CreateRenegotiationRequest struct {
	Reason         string `json:"reason"`
	NewAmountCents int64  `json:"newAmountCents"`
}

type CancelSaleRequest struct {
	Reason string `json:"reason"`
}

type SaleCancellationResponse struct {
	Sale   SaleResponse `json:"sale"`
	Reason string       `json:"reason"`
}
