package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
)

type ProjectResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Status      models.ProjectStatus `json:"status"`
	OwnerID     uuid.UUID            `json:"owner_id"`
	Owner       *UserResponse        `json:"owner,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type CreateProjectRequest struct {
	Name        string               `json:"name" validate:"required,min=2,max=255"`
	Description string               `json:"description" validate:"max=5000"`
	Status      models.ProjectStatus `json:"status" validate:"omitempty,oneof=active archived completed"`
}

type UpdateProjectRequest struct {
	Name        *string               `json:"name" validate:"omitempty,min=2,max=255"`
	Description *string               `json:"description" validate:"omitempty,max=5000"`
	Status      *models.ProjectStatus `json:"status" validate:"omitempty,oneof=active archived completed"`
}

type ProjectFilterQuery struct {
	PaginationQuery
	Status  string `form:"status"`
	OwnerID string `form:"owner_id"`
}

func ToProjectResponse(project *models.Project) ProjectResponse {
	resp := ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		Status:      project.Status,
		OwnerID:     project.OwnerID,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
	if project.Owner.ID != uuid.Nil {
		owner := ToUserResponse(&project.Owner)
		resp.Owner = &owner
	}
	return resp
}

func ToProjectResponses(projects []models.Project) []ProjectResponse {
	result := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		result[i] = ToProjectResponse(&p)
	}
	return result
}
