package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
)

type UserResponse struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	Role      models.UserRole `json:"role"`
	Active    bool            `json:"active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string          `json:"email" validate:"required,email"`
	Password string          `json:"password" validate:"required,min=6"`
	Name     string          `json:"name" validate:"required,min=2,max=255"`
	Role     models.UserRole `json:"role" validate:"omitempty,oneof=admin member"`
}

type UpdateUserRequest struct {
	Email    *string          `json:"email" validate:"omitempty,email"`
	Password *string          `json:"password" validate:"omitempty,min=6"`
	Name     *string          `json:"name" validate:"omitempty,min=2,max=255"`
	Role     *models.UserRole `json:"role" validate:"omitempty,oneof=admin member"`
	Active   *bool            `json:"active"`
}

type UserFilterQuery struct {
	PaginationQuery
	Role   string `form:"role"`
	Active *bool  `form:"active"`
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponses(users []models.User) []UserResponse {
	result := make([]UserResponse, len(users))
	for i, u := range users {
		result[i] = ToUserResponse(&u)
	}
	return result
}
