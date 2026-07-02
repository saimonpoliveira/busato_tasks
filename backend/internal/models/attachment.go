package models

import "github.com/google/uuid"

type Attachment struct {
	BaseModel
	EntityType   EntityType `gorm:"type:varchar(20);not null;index" json:"entity_type"`
	EntityID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"entity_id"`
	Filename     string     `gorm:"not null;size:500" json:"filename"`
	OriginalName string     `gorm:"not null;size:500" json:"original_name"`
	FilePath     string     `gorm:"not null;size:1000" json:"file_path"`
	Size         int64      `gorm:"not null" json:"size"`
	MimeType     string     `gorm:"size:255" json:"mime_type"`
	UploadedByID uuid.UUID  `gorm:"type:uuid;not null" json:"uploaded_by_id"`
	UploadedBy   User       `gorm:"foreignKey:UploadedByID" json:"uploaded_by,omitempty"`
}
