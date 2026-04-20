// Este cliente resolve contexto de acesso no identity antes de servir rotas protegidas do edge.
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessResolution struct {
	TenantSlug    string   `json:"tenantSlug"`
	SessionPublic string   `json:"sessionPublicId"`
	UserPublicID  string   `json:"userPublicId"`
	Email         string   `json:"email"`
	DisplayName   string   `json:"displayName"`
	RoleCodes     []string `json:"roleCodes"`
	MFAEnabled    bool     `json:"mfaEnabled"`
	Authorized    bool     `json:"authorized"`
	Status        string   `json:"status"`
}

type AccessResolutionError struct {
	StatusCode int
	Code       string
	Message    string
}

func (err AccessResolutionError) Error() string {
	return fmt.Sprintf("%s:%s", err.Code, err.Message)
}

type TenantAccessResolver interface {
	ResolveTenantAccess(ctx context.Context, identityBaseURL string, tenantSlug string, sessionToken string) (AccessResolution, error)
}

type HTTPIdentityAccessResolver struct {
	client *http.Client
}

func NewHTTPIdentityAccessResolver(timeout time.Duration) HTTPIdentityAccessResolver {
	return HTTPIdentityAccessResolver{
		client: &http.Client{Timeout: timeout},
	}
}

func (resolver HTTPIdentityAccessResolver) ResolveTenantAccess(ctx context.Context, identityBaseURL string, tenantSlug string, sessionToken string) (AccessResolution, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		strings.TrimRight(identityBaseURL, "/")+"/api/identity/tenants/"+url.PathEscape(tenantSlug)+"/access",
		nil,
	)
	if err != nil {
		return AccessResolution{}, err
	}

	request.Header.Set("Authorization", "Bearer "+sessionToken)
	response, err := resolver.client.Do(request)
	if err != nil {
		return AccessResolution{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var payload map[string]string
		_ = json.NewDecoder(response.Body).Decode(&payload)
		return AccessResolution{}, AccessResolutionError{
			StatusCode: response.StatusCode,
			Code:       payload["code"],
			Message:    payload["message"],
		}
	}

	var resolution AccessResolution
	if err := json.NewDecoder(response.Body).Decode(&resolution); err != nil {
		return AccessResolution{}, err
	}

	return resolution, nil
}
