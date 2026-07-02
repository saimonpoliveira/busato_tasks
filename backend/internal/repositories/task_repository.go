package repositories

import (
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uuid.UUID) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id uuid.UUID) error
	FindAll(filter dto.TaskFilterQuery) ([]models.Task, int64, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("Assignee").First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Task{}, "id = ?", id).Error
}

var taskSortFields = map[string]string{
	"title":      "title",
	"status":     "status",
	"order":      "\"order\"",
	"created_at": "created_at",
}

func (r *taskRepository) FindAll(filter dto.TaskFilterQuery) ([]models.Task, int64, error) {
	filter.Normalize()

	query := r.db.Model(&models.Task{}).Preload("Assignee")

	if filter.TicketID != "" {
		query = query.Where("ticket_id = ?", filter.TicketID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.AssigneeID != "" {
		query = query.Where("assignee_id = ?", filter.AssigneeID)
	}

	query = utils.ApplySearch(query, filter.Search, "title", "description")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySorting(query, filter.SortBy, filter.SortOrder, taskSortFields)
	query = utils.ApplyPagination(query, filter.Page, filter.PageSize)

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
