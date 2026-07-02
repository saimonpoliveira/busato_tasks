package repositories

import (
	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/dto"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
	"gorm.io/gorm"
)

type TicketRepository interface {
	Create(ticket *models.Ticket) error
	FindByID(id uuid.UUID) (*models.Ticket, error)
	Update(ticket *models.Ticket) error
	Delete(id uuid.UUID) error
	FindAll(filter dto.TicketFilterQuery) ([]models.Ticket, int64, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) FindByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("Assignee").Preload("Reporter").First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *ticketRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Ticket{}, "id = ?", id).Error
}

var ticketSortFields = map[string]string{
	"title":      "title",
	"status":     "status",
	"priority":   "priority",
	"created_at": "created_at",
}

func (r *ticketRepository) FindAll(filter dto.TicketFilterQuery) ([]models.Ticket, int64, error) {
	filter.Normalize()

	query := r.db.Model(&models.Ticket{}).Preload("Assignee").Preload("Reporter")

	if filter.ProjectID != "" {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.AssigneeID != "" {
		query = query.Where("assignee_id = ?", filter.AssigneeID)
	}
	if filter.ReporterID != "" {
		query = query.Where("reporter_id = ?", filter.ReporterID)
	}

	query = utils.ApplySearch(query, filter.Search, "title", "description")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = utils.ApplySorting(query, filter.SortBy, filter.SortOrder, ticketSortFields)
	query = utils.ApplyPagination(query, filter.Page, filter.PageSize)

	var tickets []models.Ticket
	if err := query.Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}
