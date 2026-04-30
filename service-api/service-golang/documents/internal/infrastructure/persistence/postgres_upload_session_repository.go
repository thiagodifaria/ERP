package persistence

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thiagodifaria/erp/service-api/service-golang/documents/internal/domain/entity"
)

type PostgresUploadSessionRepository struct {
	database            *sql.DB
	bootstrapTenantSlug string
	cachedTenantIDs     map[string]int64
}

func NewPostgresUploadSessionRepository(database *sql.DB, bootstrapTenantSlug string) *PostgresUploadSessionRepository {
	return &PostgresUploadSessionRepository{
		database:            database,
		bootstrapTenantSlug: strings.ToLower(strings.TrimSpace(bootstrapTenantSlug)),
		cachedTenantIDs:     map[string]int64{},
	}
}

func (repository *PostgresUploadSessionRepository) Save(session entity.UploadSession) entity.UploadSession {
	tenantID, _, err := repository.resolveTenant(session.TenantSlug)
	if err != nil {
		return session
	}

	row := repository.database.QueryRow(
		`
      INSERT INTO documents.upload_sessions (
        tenant_id,
        public_id,
        owner_type,
        owner_public_id,
        file_name,
        content_type,
        storage_key,
        storage_driver,
        source,
        requested_by,
        visibility,
        retention_days,
        status,
        expires_at
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
      RETURNING public_id::text
    `,
		tenantID,
		uuid.MustParse(session.PublicID),
		session.OwnerType,
		uuid.MustParse(session.OwnerPublicID),
		session.FileName,
		session.ContentType,
		session.StorageKey,
		session.StorageDriver,
		session.Source,
		session.RequestedBy,
		session.Visibility,
		session.RetentionDays,
		session.Status,
		session.ExpiresAt,
	)

	var publicID string
	if err := row.Scan(&publicID); err != nil {
		return session
	}

	session.PublicID = publicID
	return session
}

func (repository *PostgresUploadSessionRepository) FindByPublicID(tenantSlug string, publicID string) (entity.UploadSession, bool) {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.UploadSession{}, false
	}

	row := repository.database.QueryRow(
		`
      SELECT public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, requested_by, visibility, retention_days, status, COALESCE(attachment_public_id::text, ''), expires_at, completed_at, created_at
      FROM documents.upload_sessions
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
    `,
		tenantID,
		strings.TrimSpace(publicID),
	)

	return scanUploadSession(row, normalizedTenantSlug)
}

func (repository *PostgresUploadSessionRepository) Complete(tenantSlug string, publicID string, attachmentPublicID string, completedAt time.Time) (entity.UploadSession, bool) {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.UploadSession{}, false
	}

	value := completedAt.UTC()
	if value.IsZero() {
		value = time.Now().UTC()
	}

	row := repository.database.QueryRow(
		`
      UPDATE documents.upload_sessions
      SET
        status = 'completed',
        attachment_public_id = $3::uuid,
        completed_at = $4,
        updated_at = timezone('utc', now())
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      RETURNING public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, requested_by, visibility, retention_days, status, COALESCE(attachment_public_id::text, ''), expires_at, completed_at, created_at
    `,
		tenantID,
		strings.TrimSpace(publicID),
		strings.TrimSpace(attachmentPublicID),
		value,
	)

	return scanUploadSession(row, normalizedTenantSlug)
}

func (repository *PostgresUploadSessionRepository) resolveTenant(tenantSlug string) (int64, string, error) {
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

func scanUploadSession(scanner attachmentScanner, tenantSlug string) (entity.UploadSession, bool) {
	var publicID string
	var ownerType string
	var ownerPublicID string
	var fileName string
	var contentType string
	var storageKey string
	var storageDriver string
	var source string
	var requestedBy string
	var visibility string
	var retentionDays int
	var status string
	var attachmentPublicID string
	var expiresAt sql.NullTime
	var completedAt sql.NullTime
	var createdAt sql.NullTime

	if err := scanner.Scan(
		&publicID,
		&ownerType,
		&ownerPublicID,
		&fileName,
		&contentType,
		&storageKey,
		&storageDriver,
		&source,
		&requestedBy,
		&visibility,
		&retentionDays,
		&status,
		&attachmentPublicID,
		&expiresAt,
		&completedAt,
		&createdAt,
	); err != nil {
		return entity.UploadSession{}, false
	}

	var completedAtValue *time.Time
	if completedAt.Valid {
		value := completedAt.Time.UTC()
		completedAtValue = &value
	}

	session, buildErr := entity.NewUploadSession(
		publicID,
		tenantSlug,
		ownerType,
		ownerPublicID,
		fileName,
		contentType,
		storageKey,
		storageDriver,
		source,
		requestedBy,
		visibility,
		retentionDays,
		status,
		attachmentPublicID,
		expiresAt.Time,
		completedAtValue,
		createdAt.Time,
	)
	if buildErr != nil {
		return entity.UploadSession{}, false
	}

	return session, true
}
