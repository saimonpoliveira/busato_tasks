package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
)

type TicketResponse struct {
	ID          uuid.UUID             `json:"id"`
	ProjectID   uuid.UUID             `json:"project_id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Status      models.TicketStatus   `json:"status"`
	Priority    models.TicketPriority `json:"priority"`
	AssigneeID  *uuid.UUID            `json:"assignee_id,omitempty"`
	Assignee    *UserResponse         `json:"assignee,omitempty"`
	ReporterID  uuid.UUID             `json:"reporter_id"`
	Reporter    *UserResponse         `json:"reporter,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type CreateTicketRequest struct {
	ProjectID   uuid.UUID             `json:"project_id" validate:"required"`
	Title       string                `json:"title" validate:"required,min=2,max=500"`
	Description string                `json:"description" validate:"max=10000"`
	Status      models.TicketStatus   `json:"status" validate:"omitempty,oneof=open in_progress resolved closed"`
	Priority    models.TicketPriority `json:"priority" validate:"omitempty,oneof=low medium high critical"`
	AssigneeID  *uuid.UUID            `json:"assignee_id"`
}

type UpdateTicketRequest struct {
	Title       *string                `json:"title" validate:"omitempty,min=2,max=500"`
	Description *string                `json:"description" validate:"omitempty,max=10000"`
	Status      *models.TicketStatus   `json:"status" validate:"omitempty,oneof=open in_progress resolved closed"`
	Priority    *models.TicketPriority `json:"priority" validate:"omitempty,oneof=low medium high critical"`
	AssigneeID  *uuid.UUID             `json:"assignee_id"`
}

type TicketFilterQuery struct {
	PaginationQuery
	ProjectID  string `form:"project_id"`
	Status     string `form:"status"`
	Priority   string `form:"priority"`
	AssigneeID string `form:"assignee_id"`
	ReporterID string `form:"reporter_id"`
}

func ToTicketResponse(ticket *models.Ticket) TicketResponse {
	resp := TicketResponse{
		ID:          ticket.ID,
		ProjectID:   ticket.ProjectID,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		Priority:    ticket.Priority,
		AssigneeID:  ticket.AssigneeID,
		ReporterID:  ticket.ReporterID,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}
	if ticket.Assignee != nil && ticket.Assignee.ID != uuid.Nil {
		assignee := ToUserResponse(ticket.Assignee)
		resp.Assignee = &assignee
	}
	if ticket.Reporter.ID != uuid.Nil {
		reporter := ToUserResponse(&ticket.Reporter)
		resp.Reporter = &reporter
	}
	return resp
}

func ToTicketResponses(tickets []models.Ticket) []TicketResponse {
	result := make([]TicketResponse, len(tickets))
	for i, t := range tickets {
		result[i] = ToTicketResponse(&t)
	}
	return result
}
