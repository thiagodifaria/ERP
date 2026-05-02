package dto

type StorageCapabilityResponse struct {
	Provider       string   `json:"provider"`
	Scope          string   `json:"scope"`
	Configured     bool     `json:"configured"`
	CredentialKey  string   `json:"credentialKey,omitempty"`
	Mode           string   `json:"mode"`
	Status         string   `json:"status"`
	FallbackViable bool     `json:"fallbackViable"`
	SupportsLinks  bool     `json:"supportsLinks"`
	SupportsUpload bool     `json:"supportsUpload"`
	Notes          []string `json:"notes"`
}
