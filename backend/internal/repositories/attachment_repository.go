package repositories

import (
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(attachment *models.Attachment) error
	FindByID(id uuid.UUID) (*models.Attachment, error)
	Delete(id uuid.UUID) error
	FindAll(filter dto.AttachmentFilterQuery) ([]models.Attachment, int64, error)
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(attachment *models.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) FindByID(id uuid.UUID) (*models.Attachment, error) {
	var attachment models.Attachment
	err := r.db.Preload("UploadedBy").First(&attachment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *attachmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Attachment{}, "id = ?", id).Error
}

var attachmentSortFields = map[string]string{
	"filename":   "filename",
	"size":       "size",
	"created_at": "created_at",
}

func (r *attachmentRepository) FindAll(filter dto.AttachmentFilterQuery) ([]models.Attachment, int64, error) {
	filter.Normalize()

	query := r.db.Model(&models.Attachment{}).Preload("UploadedBy")

	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.EntityID != "" {
		query = query.Where("entity_id = ?", filter.EntityID)
	}

	query = utils.ApplySearch(query, filter.Search, "original_name", "filename")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySorting(query, filter.SortBy, filter.SortOrder, attachmentSortFields)
	query = utils.ApplyPagination(query, filter.Page, filter.PageSize)

	var attachments []models.Attachment
	if err := query.Find(&attachments).Error; err != nil {
		return nil, 0, err
	}

	return attachments, total, nil
}
