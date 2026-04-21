package dto

type ReadinessResponse struct {
	Service      string               `json:"service"`
	Status       string               `json:"status"`
	Dependencies []DependencyResponse `json:"dependencies"`
}

type DependencyResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
