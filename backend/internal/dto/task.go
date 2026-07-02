package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
)

type TaskResponse struct {
	ID          uuid.UUID        `json:"id"`
	TicketID    uuid.UUID        `json:"ticket_id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Status      models.TaskStatus `json:"status"`
	AssigneeID  *uuid.UUID       `json:"assignee_id,omitempty"`
	Assignee    *UserResponse    `json:"assignee,omitempty"`
	Order       int              `json:"order"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type CreateTaskRequest struct {
	TicketID    uuid.UUID         `json:"ticket_id" validate:"required"`
	Title       string            `json:"title" validate:"required,min=2,max=500"`
	Description string            `json:"description" validate:"max=10000"`
	Status      models.TaskStatus `json:"status" validate:"omitempty,oneof=todo in_progress done cancelled"`
	AssigneeID  *uuid.UUID        `json:"assignee_id"`
	Order       *int              `json:"order"`
}

type UpdateTaskRequest struct {
	Title       *string            `json:"title" validate:"omitempty,min=2,max=500"`
	Description *string            `json:"description" validate:"omitempty,max=10000"`
	Status      *models.TaskStatus `json:"status" validate:"omitempty,oneof=todo in_progress done cancelled"`
	AssigneeID  *uuid.UUID         `json:"assignee_id"`
	Order       *int               `json:"order"`
}

type TaskFilterQuery struct {
	PaginationQuery
	TicketID   string `form:"ticket_id"`
	Status     string `form:"status"`
	AssigneeID string `form:"assignee_id"`
}

func ToTaskResponse(task *models.Task) TaskResponse {
	resp := TaskResponse{
		ID:          task.ID,
		TicketID:    task.TicketID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		AssigneeID:  task.AssigneeID,
		Order:       task.Order,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if task.Assignee != nil && task.Assignee.ID != uuid.Nil {
		assignee := ToUserResponse(task.Assignee)
		resp.Assignee = &assignee
	}
	return resp
}

func ToTaskResponses(tasks []models.Task) []TaskResponse {
	result := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		result[i] = ToTaskResponse(&t)
	}
	return result
}
