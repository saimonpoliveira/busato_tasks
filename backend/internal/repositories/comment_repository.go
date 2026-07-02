package repositories

import (
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindByID(id uuid.UUID) (*models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uuid.UUID) error
	FindAll(filter dto.CommentFilterQuery) ([]models.Comment, int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Comment{}, "id = ?", id).Error
}

var commentSortFields = map[string]string{
	"created_at": "created_at",
}

func (r *commentRepository) FindAll(filter dto.CommentFilterQuery) ([]models.Comment, int64, error) {
	filter.Normalize()

	query := r.db.Model(&models.Comment{}).Preload("User")

	if filter.EntityType != "" {
		query = query.Where("entity_type = ?", filter.EntityType)
	}
	if filter.EntityID != "" {
		query = query.Where("entity_id = ?", filter.EntityID)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}

	query = utils.ApplySearch(query, filter.Search, "content")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySorting(query, filter.SortBy, filter.SortOrder, commentSortFields)
	query = utils.ApplyPagination(query, filter.Page, filter.PageSize)

	var comments []models.Comment
	if err := query.Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
