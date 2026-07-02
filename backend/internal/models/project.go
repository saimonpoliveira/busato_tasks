package models

import "github.com/google/uuid"

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusArchived  ProjectStatus = "archived"
	ProjectStatusCompleted ProjectStatus = "completed"
)

type Project struct {
	BaseModel
	Name        string        `gorm:"not null;size:255" json:"name"`
	Description string        `gorm:"type:text" json:"description"`
	Status      ProjectStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	OwnerID     uuid.UUID     `gorm:"type:uuid;not null" json:"owner_id"`
	Owner       User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}
