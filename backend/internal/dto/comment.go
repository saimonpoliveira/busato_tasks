package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
)

type CommentResponse struct {
	ID         uuid.UUID         `json:"id"`
	EntityType models.EntityType `json:"entity_type"`
	EntityID   uuid.UUID         `json:"entity_id"`
	UserID     uuid.UUID         `json:"user_id"`
	User       *UserResponse     `json:"user,omitempty"`
	Content    string            `json:"content"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type CreateCommentRequest struct {
	EntityType models.EntityType `json:"entity_type" validate:"required,oneof=ticket task"`
	EntityID   uuid.UUID         `json:"entity_id" validate:"required"`
	Content    string            `json:"content" validate:"required,min=1,max=10000"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=10000"`
}

type CommentFilterQuery struct {
	PaginationQuery
	EntityType string `form:"entity_type"`
	EntityID   string `form:"entity_id"`
	UserID     string `form:"user_id"`
}

func ToCommentResponse(comment *models.Comment) CommentResponse {
	resp := CommentResponse{
		ID:         comment.ID,
		EntityType: comment.EntityType,
		EntityID:   comment.EntityID,
		UserID:     comment.UserID,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}
	if comment.User.ID != uuid.Nil {
		user := ToUserResponse(&comment.User)
		resp.User = &user
	}
	return resp
}

func ToCommentResponses(comments []models.Comment) []CommentResponse {
	result := make([]CommentResponse, len(comments))
	for i, c := range comments {
		result[i] = ToCommentResponse(&c)
	}
	return result
}
