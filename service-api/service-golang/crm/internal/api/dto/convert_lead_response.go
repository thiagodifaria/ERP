package dto

type ConvertLeadResponse struct {
	Lead     LeadResponse     `json:"lead"`
	Customer CustomerResponse `json:"customer"`
}
