package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/models"
)

type AttachmentResponse struct {
	ID           uuid.UUID         `json:"id"`
	EntityType   models.EntityType `json:"entity_type"`
	EntityID     uuid.UUID         `json:"entity_id"`
	Filename     string            `json:"filename"`
	OriginalName string            `json:"original_name"`
	Size         int64             `json:"size"`
	MimeType     string            `json:"mime_type"`
	UploadedByID uuid.UUID         `json:"uploaded_by_id"`
	UploadedBy   *UserResponse     `json:"uploaded_by,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
}

type AttachmentFilterQuery struct {
	PaginationQuery
	EntityType string `form:"entity_type"`
	EntityID   string `form:"entity_id"`
}

func ToAttachmentResponse(attachment *models.Attachment) AttachmentResponse {
	resp := AttachmentResponse{
		ID:           attachment.ID,
		EntityType:   attachment.EntityType,
		EntityID:     attachment.EntityID,
		Filename:     attachment.Filename,
		OriginalName: attachment.OriginalName,
		Size:         attachment.Size,
		MimeType:     attachment.MimeType,
		UploadedByID: attachment.UploadedByID,
		CreatedAt:    attachment.CreatedAt,
	}
	if attachment.UploadedBy.ID != uuid.Nil {
		user := ToUserResponse(&attachment.UploadedBy)
		resp.UploadedBy = &user
	}
	return resp
}

func ToAttachmentResponses(attachments []models.Attachment) []AttachmentResponse {
	result := make([]AttachmentResponse, len(attachments))
	for i, a := range attachments {
		result[i] = ToAttachmentResponse(&a)
	}
	return result
}
