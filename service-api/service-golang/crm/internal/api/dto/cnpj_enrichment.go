package dto

type CNPJEnrichmentCapabilityResponse struct {
	Provider       string   `json:"provider"`
	Scope          string   `json:"scope"`
	Configured     bool     `json:"configured"`
	CredentialKey  string   `json:"credentialKey,omitempty"`
	Mode           string   `json:"mode"`
	Status         string   `json:"status"`
	FallbackViable bool     `json:"fallbackViable"`
	Notes          []string `json:"notes"`
}

type CNPJEnrichmentLookupRequest struct {
	CNPJ        string `json:"cnpj"`
	CompanyName string `json:"companyName"`
}

type CNPJEnrichmentLookupResponse struct {
	Provider          string `json:"provider"`
	CNPJ              string `json:"cnpj"`
	NormalizedCNPJ    string `json:"normalizedCnpj"`
	Status            string `json:"status"`
	CompanyName       string `json:"companyName"`
	TradeName         string `json:"tradeName"`
	LegalNature       string `json:"legalNature"`
	TaxRegimeHint     string `json:"taxRegimeHint"`
	PrimaryActivity   string `json:"primaryActivity"`
	HeadquartersCity  string `json:"headquartersCity"`
	HeadquartersState string `json:"headquartersState"`
	FallbackUsed      bool   `json:"fallbackUsed"`
}
