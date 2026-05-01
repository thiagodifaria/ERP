package persistence

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
      SELECT public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, uploaded_by, file_size_bytes, checksum_sha256, visibility, retention_days, current_version_number, archive_reason, archived_at, created_at
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
	source := strings.TrimSpace(filters.Source)
	if source != "" {
		query += fmt.Sprintf(" AND source = $%d", len(args)+1)
		args = append(args, source)
	}
	visibility := strings.ToLower(strings.TrimSpace(filters.Visibility))
	if visibility != "" {
		query += fmt.Sprintf(" AND visibility = $%d", len(args)+1)
		args = append(args, visibility)
	}
	if strings.EqualFold(strings.TrimSpace(filters.Archived), "true") {
		query += " AND archived_at IS NOT NULL"
	}
	if strings.EqualFold(strings.TrimSpace(filters.Archived), "false") {
		query += " AND archived_at IS NULL"
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
		var fileSizeBytes int64
		var checksumSHA256 string
		var visibility string
		var retentionDays int
		var currentVersion int
		var archiveReason string
		var archivedAt sql.NullTime
		var createdAt sql.NullTime

		if scanErr := rows.Scan(&publicID, &ownerType, &ownerPublicID, &fileName, &contentType, &storageKey, &storageDriver, &source, &uploadedBy, &fileSizeBytes, &checksumSHA256, &visibility, &retentionDays, &currentVersion, &archiveReason, &archivedAt, &createdAt); scanErr != nil {
			continue
		}

		var archivedAtValue *time.Time
		if archivedAt.Valid {
			value := archivedAt.Time.UTC()
			archivedAtValue = &value
		}

		attachment, buildErr := entity.NewAttachment(publicID, tenantSlug, ownerType, ownerPublicID, fileName, contentType, storageKey, storageDriver, source, uploadedBy, fileSizeBytes, checksumSHA256, visibility, retentionDays, archiveReason, archivedAtValue, createdAt.Time)
		if buildErr == nil {
			attachment.CurrentVersion = currentVersion
			attachment.VersionCount = currentVersion
			response = append(response, attachment)
		}
	}

	return response
}

func (repository *PostgresAttachmentRepository) FindByPublicID(tenantSlug string, publicID string) (entity.Attachment, bool) {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.Attachment{}, false
	}

	row := repository.database.QueryRow(
		`
      SELECT public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, uploaded_by, file_size_bytes, checksum_sha256, visibility, retention_days, current_version_number, archive_reason, archived_at, created_at
      FROM documents.attachments
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
    `,
		tenantID,
		strings.TrimSpace(publicID),
	)

	return scanAttachment(row, normalizedTenantSlug)
}

func (repository *PostgresAttachmentRepository) Save(attachment entity.Attachment) entity.Attachment {
	tenantID, _, err := repository.resolveTenant(attachment.TenantSlug)
	if err != nil {
		return attachment
	}

	transaction, err := repository.database.Begin()
	if err != nil {
		return attachment
	}
	defer transaction.Rollback()

	row := transaction.QueryRow(
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
        uploaded_by,
        file_size_bytes,
        checksum_sha256,
        visibility,
        retention_days,
        current_version_number
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 1)
      RETURNING id, public_id::text
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
		attachment.FileSizeBytes,
		attachment.ChecksumSHA256,
		attachment.Visibility,
		attachment.RetentionDays,
	)

	var attachmentID int64
	var publicID string
	if err := row.Scan(&attachmentID, &publicID); err != nil {
		return attachment
	}

	if _, err := transaction.Exec(
		`
      INSERT INTO documents.attachment_versions (
        tenant_id, attachment_id, public_id, version_number, file_name, content_type, storage_key, storage_driver, source, uploaded_by, file_size_bytes, checksum_sha256
      )
      VALUES ($1, $2, $3::uuid, 1, $4, $5, $6, $7, $8, $9, $10, $11)
    `,
		tenantID,
		attachmentID,
		uuid.MustParse(attachment.PublicID),
		attachment.FileName,
		attachment.ContentType,
		attachment.StorageKey,
		attachment.StorageDriver,
		attachment.Source,
		attachment.UploadedBy,
		attachment.FileSizeBytes,
		attachment.ChecksumSHA256,
	); err != nil {
		return attachment
	}

	if err := transaction.Commit(); err != nil {
		return attachment
	}

	attachment.PublicID = publicID
	attachment.CurrentVersion = 1
	attachment.VersionCount = 1
	return attachment
}

func (repository *PostgresAttachmentRepository) Archive(tenantSlug string, publicID string, reason string) (entity.Attachment, bool) {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.Attachment{}, false
	}

	row := repository.database.QueryRow(
		`
      UPDATE documents.attachments
      SET
        archived_at = COALESCE(archived_at, timezone('utc', now())),
        archive_reason = CASE
          WHEN archived_at IS NULL THEN COALESCE(NULLIF($3, ''), 'archived')
          ELSE archive_reason
        END
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      RETURNING public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, uploaded_by, file_size_bytes, checksum_sha256, visibility, retention_days, current_version_number, archive_reason, archived_at, created_at
    `,
		tenantID,
		strings.TrimSpace(publicID),
		strings.TrimSpace(reason),
	)

	return scanAttachment(row, normalizedTenantSlug)
}

func (repository *PostgresAttachmentRepository) ListVersions(tenantSlug string, attachmentPublicID string) []entity.AttachmentVersion {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return []entity.AttachmentVersion{}
	}

	rows, err := repository.database.Query(
		`
      SELECT
        version.public_id::text,
        attachment.public_id::text,
        version.version_number,
        version.file_name,
        version.content_type,
        version.storage_key,
        version.storage_driver,
        version.source,
        version.uploaded_by,
        version.file_size_bytes,
        version.checksum_sha256,
        version.created_at
      FROM documents.attachment_versions AS version
      INNER JOIN documents.attachments AS attachment
        ON attachment.id = version.attachment_id
      WHERE version.tenant_id = $1
        AND attachment.public_id = $2::uuid
      ORDER BY version.version_number DESC, version.id DESC
    `,
		tenantID,
		strings.TrimSpace(attachmentPublicID),
	)
	if err != nil {
		return []entity.AttachmentVersion{}
	}
	defer rows.Close()

	response := make([]entity.AttachmentVersion, 0)
	for rows.Next() {
		var publicID string
		var parentPublicID string
		var versionNumber int
		var fileName string
		var contentType string
		var storageKey string
		var storageDriver string
		var source string
		var uploadedBy string
		var fileSizeBytes int64
		var checksumSHA256 string
		var createdAt time.Time
		if err := rows.Scan(&publicID, &parentPublicID, &versionNumber, &fileName, &contentType, &storageKey, &storageDriver, &source, &uploadedBy, &fileSizeBytes, &checksumSHA256, &createdAt); err != nil {
			continue
		}

		version, buildErr := entity.NewAttachmentVersion(publicID, normalizedTenantSlug, parentPublicID, versionNumber, fileName, contentType, storageKey, storageDriver, source, uploadedBy, fileSizeBytes, checksumSHA256, createdAt)
		if buildErr == nil {
			response = append(response, version)
		}
	}

	return response
}

func (repository *PostgresAttachmentRepository) CreateVersion(tenantSlug string, attachmentPublicID string, version entity.AttachmentVersion) (entity.Attachment, entity.AttachmentVersion, bool) {
	tenantID, normalizedTenantSlug, err := repository.resolveTenant(tenantSlug)
	if err != nil {
		return entity.Attachment{}, entity.AttachmentVersion{}, false
	}

	transaction, err := repository.database.Begin()
	if err != nil {
		return entity.Attachment{}, entity.AttachmentVersion{}, false
	}
	defer transaction.Rollback()

	var attachmentID int64
	if err := transaction.QueryRow(
		`
      SELECT id
      FROM documents.attachments
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      LIMIT 1
    `,
		tenantID,
		strings.TrimSpace(attachmentPublicID),
	).Scan(&attachmentID); err != nil {
		return entity.Attachment{}, entity.AttachmentVersion{}, false
	}

	if _, err := transaction.Exec(
		`
      INSERT INTO documents.attachment_versions (
        tenant_id, attachment_id, public_id, version_number, file_name, content_type, storage_key, storage_driver, source, uploaded_by, file_size_bytes, checksum_sha256
      )
      VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `,
		tenantID,
		attachmentID,
		uuid.MustParse(version.PublicID),
		version.VersionNumber,
		version.FileName,
		version.ContentType,
		version.StorageKey,
		version.StorageDriver,
		version.Source,
		version.UploadedBy,
		version.FileSizeBytes,
		version.ChecksumSHA256,
	); err != nil {
		return entity.Attachment{}, entity.AttachmentVersion{}, false
	}

	row := transaction.QueryRow(
		`
      UPDATE documents.attachments
      SET
        file_name = $3,
        content_type = $4,
        storage_key = $5,
        storage_driver = $6,
        source = $7,
        uploaded_by = $8,
        file_size_bytes = $9,
        checksum_sha256 = $10,
        current_version_number = $11
      WHERE tenant_id = $1
        AND public_id = $2::uuid
      RETURNING public_id::text, owner_type, owner_public_id::text, file_name, content_type, storage_key, storage_driver, source, uploaded_by, file_size_bytes, checksum_sha256, visibility, retention_days, current_version_number, archive_reason, archived_at, created_at
    `,
		tenantID,
		strings.TrimSpace(attachmentPublicID),
		version.FileName,
		version.ContentType,
		version.StorageKey,
		version.StorageDriver,
		version.Source,
		version.UploadedBy,
		version.FileSizeBytes,
		version.ChecksumSHA256,
		version.VersionNumber,
	)

	updatedAttachment, ok := scanAttachment(row, normalizedTenantSlug)
	if !ok {
		return entity.Attachment{}, entity.AttachmentVersion{}, false
	}

	if err := transaction.Commit(); err != nil {
		return entity.Attachment{}, entity.AttachmentVersion{}, false
	}

	return updatedAttachment, version, true
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

type attachmentScanner interface {
	Scan(dest ...any) error
}

func scanAttachment(scanner attachmentScanner, tenantSlug string) (entity.Attachment, bool) {
	var publicID string
	var ownerType string
	var ownerPublicID string
	var fileName string
	var contentType string
	var storageKey string
	var storageDriver string
	var source string
	var uploadedBy string
	var fileSizeBytes int64
	var checksumSHA256 string
	var visibility string
	var retentionDays int
	var currentVersion int
	var archiveReason string
	var archivedAt sql.NullTime
	var createdAt sql.NullTime

	if err := scanner.Scan(&publicID, &ownerType, &ownerPublicID, &fileName, &contentType, &storageKey, &storageDriver, &source, &uploadedBy, &fileSizeBytes, &checksumSHA256, &visibility, &retentionDays, &currentVersion, &archiveReason, &archivedAt, &createdAt); err != nil {
		return entity.Attachment{}, false
	}

	var archivedAtValue *time.Time
	if archivedAt.Valid {
		value := archivedAt.Time.UTC()
		archivedAtValue = &value
	}

	attachment, err := entity.NewAttachment(publicID, tenantSlug, ownerType, ownerPublicID, fileName, contentType, storageKey, storageDriver, source, uploadedBy, fileSizeBytes, checksumSHA256, visibility, retentionDays, archiveReason, archivedAtValue, createdAt.Time)
	if err != nil {
		return entity.Attachment{}, false
	}
	attachment.CurrentVersion = currentVersion
	attachment.VersionCount = currentVersion

	return attachment, true
}
