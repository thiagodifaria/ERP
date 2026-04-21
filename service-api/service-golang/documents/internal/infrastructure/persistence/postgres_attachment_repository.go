package persistence

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/repository"
)

type PostgresAttachmentRepository struct {
	database            *sql.DB
	bootstrapTenantSlug string
	cachedTenantIDs     map[string]int64
}

func NewPostgresAttachmentRepository(database *sql.DB, bootstrapTenantSlug string) *PostgresAttachmentRepository {
	return &PostgresAttachmentRepository{
		database:            database,
		bootstrapTenantSlug: strings.ToLower(strings.TrimSpace(bootstrapTenantSlug)),
		cachedTenantIDs:     map[string]int64{},
	}
}

func (repository *PostgresAttachmentRepository) List(filters repository.AttachmentFilters) []entity.Attachment {
	tenantID, tenantSlug, err := repository.resolveTenant(filters.TenantSlug)
	if err != nil {
		return []entity.Attachment{}
	}

	query := `
      SELECT public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, uploaded_by, created_at
      FROM documents.attachments
      WHERE tenant_id = $1
    `
	args := []any{tenantID}

	ownerType := strings.TrimSpace(filters.OwnerType)
	if ownerType != "" {
		query += fmt.Sprintf(" AND owner_type = $%d", len(args)+1)
		args = append(args, ownerType)
	}

	ownerPublicID := strings.TrimSpace(filters.OwnerPublicID)
	if ownerPublicID != "" {
		query += fmt.Sprintf(" AND owner_public_id = $%d::uuid", len(args)+1)
		args = append(args, ownerPublicID)
	}

	query += " ORDER BY created_at, id"

	rows, err := repository.database.Query(query, args...)
	if err != nil {
		return []entity.Attachment{}
	}
	defer rows.Close()

	response := make([]entity.Attachment, 0)
	for rows.Next() {
		var publicID string
		var ownerType string
		var ownerPublicID string
		var fileName string
		var contentType string
		var storageKey string
		var storageDriver string
		var source string
		var uploadedBy string
		var createdAt sql.NullTime

		if scanErr := rows.Scan(&publicID, &ownerType, &ownerPublicID, &fileName, &contentType, &storageKey, &storageDriver, &source, &uploadedBy, &createdAt); scanErr != nil {
			continue
		}

		attachment, buildErr := entity.NewAttachment(publicID, tenantSlug, ownerType, ownerPublicID, fileName, contentType, storageKey, storageDriver, source, uploadedBy, createdAt.Time)
		if buildErr == nil {
			response = append(response, attachment)
		}
	}

	return response
}

func (repository *PostgresAttachmentRepository) Save(attachment entity.Attachment) entity.Attachment {
	tenantID, _, err := repository.resolveTenant(attachment.TenantSlug)
	if err != nil {
		return attachment
	}

	row := repository.database.QueryRow(
		`
      INSERT INTO documents.attachments (
        tenant_id,
        public_id,
        owner_type,
        owner_public_id,
        file_name,
        content_type,
        storage_key,
        storage_driver,
        source,
        uploaded_by
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
      RETURNING public_id::text
    `,
		tenantID,
		uuid.MustParse(attachment.PublicID),
		attachment.OwnerType,
		uuid.MustParse(attachment.OwnerPublicID),
		attachment.FileName,
		attachment.ContentType,
		attachment.StorageKey,
		attachment.StorageDriver,
		attachment.Source,
		attachment.UploadedBy,
	)

	var publicID string
	if err := row.Scan(&publicID); err != nil {
		return attachment
	}

	attachment.PublicID = publicID
	return attachment
}

func (repository *PostgresAttachmentRepository) resolveTenant(tenantSlug string) (int64, string, error) {
	slug := strings.ToLower(strings.TrimSpace(tenantSlug))
	if slug == "" {
		slug = repository.bootstrapTenantSlug
	}

	if tenantID, ok := repository.cachedTenantIDs[slug]; ok {
		return tenantID, slug, nil
	}

	var tenantID int64
	if err := repository.database.QueryRow(
		`
      SELECT id
      FROM identity.tenants
      WHERE slug = $1
      LIMIT 1
    `,
		slug,
	).Scan(&tenantID); err != nil {
		return 0, slug, err
	}

	repository.cachedTenantIDs[slug] = tenantID
	return tenantID, slug, nil
}
