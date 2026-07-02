package repositories

import (
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	FindByID(id uuid.UUID) (*models.Project, error)
	Update(project *models.Project) error
	Delete(id uuid.UUID) error
	FindAll(filter dto.ProjectFilterQuery) ([]models.Project, int64, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) FindByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Owner").First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

var projectSortFields = map[string]string{
	"name":       "name",
	"status":     "status",
	"created_at": "created_at",
}

func (r *projectRepository) FindAll(filter dto.ProjectFilterQuery) ([]models.Project, int64, error) {
	filter.Normalize()

	query := r.db.Model(&models.Project{}).Preload("Owner")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.OwnerID != "" {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}

	query = utils.ApplySearch(query, filter.Search, "name", "description")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySorting(query, filter.SortBy, filter.SortOrder, projectSortFields)
	query = utils.ApplyPagination(query, filter.Page, filter.PageSize)

	var projects []models.Project
	if err := query.Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}
