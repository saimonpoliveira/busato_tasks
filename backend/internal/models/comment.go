package models

import "github.com/google/uuid"

type EntityType string

const (
	EntityTypeTicket EntityType = "ticket"
	EntityTypeTask   EntityType = "task"
)

type Comment struct {
	BaseModel
	EntityType EntityType `gorm:"type:varchar(20);not null;index" json:"entity_type"`
	EntityID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"entity_id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content    string     `gorm:"type:text;not null" json:"content"`
}
