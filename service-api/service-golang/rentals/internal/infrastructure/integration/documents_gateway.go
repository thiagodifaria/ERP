package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/rentals/internal/domain/repository"
)

type HTTPDocumentsGateway struct {
	baseURL string
	client  *http.Client
}

type InMemoryDocumentsGateway struct {
	records []repository.AttachmentRecord
}

type documentsAttachmentPayload struct {
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

func NewHTTPDocumentsGateway(baseURL string) *HTTPDocumentsGateway {
	return &HTTPDocumentsGateway{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (gateway *HTTPDocumentsGateway) List(tenantSlug string, ownerType string, ownerPublicID string) ([]repository.AttachmentRecord, error) {
	if gateway == nil || gateway.baseURL == "" {
		return []repository.AttachmentRecord{}, nil
	}

	query := url.Values{}
	query.Set("tenantSlug", tenantSlug)
	query.Set("ownerType", ownerType)
	query.Set("ownerPublicId", ownerPublicID)

	response, err := gateway.client.Get(gateway.baseURL + "/api/documents/attachments?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("documents_list_failed:%d", response.StatusCode)
	}

	var payload []documentsAttachmentPayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	result := make([]repository.AttachmentRecord, 0, len(payload))
	for _, attachment := range payload {
		result = append(result, mapAttachmentRecord(attachment))
	}

	return result, nil
}

func (gateway *HTTPDocumentsGateway) Create(input repository.CreateAttachmentInput) (*repository.AttachmentRecord, error) {
	if gateway == nil || gateway.baseURL == "" {
		return nil, errors.New("documents_gateway_unconfigured")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, gateway.baseURL+"/api/documents/attachments", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := gateway.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("documents_create_failed:%d", response.StatusCode)
	}

	var payload documentsAttachmentPayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	record := mapAttachmentRecord(payload)
	return &record, nil
}

func NewInMemoryDocumentsGateway() *InMemoryDocumentsGateway {
	return &InMemoryDocumentsGateway{records: []repository.AttachmentRecord{}}
}

func (gateway *InMemoryDocumentsGateway) List(tenantSlug string, ownerType string, ownerPublicID string) ([]repository.AttachmentRecord, error) {
	result := make([]repository.AttachmentRecord, 0)
	for _, record := range gateway.records {
		if strings.TrimSpace(tenantSlug) != "" && record.TenantSlug != tenantSlug {
			continue
		}
		if strings.TrimSpace(ownerType) != "" && record.OwnerType != ownerType {
			continue
		}
		if strings.TrimSpace(ownerPublicID) != "" && record.OwnerPublicID != ownerPublicID {
			continue
		}
		result = append(result, record)
	}

	return result, nil
}

func (gateway *InMemoryDocumentsGateway) Create(input repository.CreateAttachmentInput) (*repository.AttachmentRecord, error) {
	record := repository.AttachmentRecord{
		PublicID:      uuid.NewString(),
		TenantSlug:    strings.TrimSpace(input.TenantSlug),
		OwnerType:     strings.TrimSpace(input.OwnerType),
		OwnerPublicID: strings.TrimSpace(input.OwnerPublicID),
		FileName:      strings.TrimSpace(input.FileName),
		ContentType:   strings.TrimSpace(input.ContentType),
		StorageKey:    strings.TrimSpace(input.StorageKey),
		StorageDriver: strings.TrimSpace(input.StorageDriver),
		Source:        strings.TrimSpace(input.Source),
		UploadedBy:    strings.TrimSpace(input.UploadedBy),
		CreatedAt:     time.Now().UTC(),
	}
	gateway.records = append(gateway.records, record)
	return &record, nil
}

func mapAttachmentRecord(payload documentsAttachmentPayload) repository.AttachmentRecord {
	return repository.AttachmentRecord{
		PublicID:      payload.PublicID,
		TenantSlug:    payload.TenantSlug,
		OwnerType:     payload.OwnerType,
		OwnerPublicID: payload.OwnerPublicID,
		FileName:      payload.FileName,
		ContentType:   payload.ContentType,
		StorageKey:    payload.StorageKey,
		StorageDriver: payload.StorageDriver,
		Source:        payload.Source,
		UploadedBy:    payload.UploadedBy,
		CreatedAt:     payload.CreatedAt,
	}
}
