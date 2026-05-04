package dto

type SigningCapabilityResponse struct {
	Provider          string   `json:"provider"`
	Scope             string   `json:"scope"`
	Configured        bool     `json:"configured"`
	CredentialKey     string   `json:"credentialKey,omitempty"`
	Mode              string   `json:"mode"`
	Status            string   `json:"status"`
	FallbackViable    bool     `json:"fallbackViable"`
	SupportsContracts bool     `json:"supportsContracts"`
	SupportsProposals bool     `json:"supportsProposals"`
	SupportsRentals   bool     `json:"supportsRentals"`
	Notes             []string `json:"notes"`
}
